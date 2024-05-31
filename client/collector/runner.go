package collector

import (
	"encoding/hex"
	"errors"
	"flare-common/contracts/relay"
	"flare-common/database"
	"flare-common/logger"
	"flare-common/payload"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

const (
	bitVoteBufferSize             = 10
	bitVoteOffChainTriggerSeconds = 15
	requestsBufferSize            = 10
	requestListenerInterval       = 2 * time.Second
	signingPolicyBufferSize       = 3
)

const (
	roundLength      = 90 * time.Second
	databasePollTime = 2 * time.Second
	bitVoteHeadStart = 5 * time.Second
)

var signingPolicyInitializedEventSel common.Hash
var attestationRequestEventSel common.Hash
var log = logger.GetLogger()

func init() {
	relayAbi, err := relay.RelayMetaData.GetAbi()
	if err != nil {
		log.Panic("cannot get relayAby:", err)
	}

	signingPolicyEvent, ok := relayAbi.Events["SigningPolicyInitialized"]

	if !ok {
		log.Panic("cannot get SigningPolicyInitialized event:", err)
	}

	signingPolicyInitializedEventSel = signingPolicyEvent.ID

	fdcAbi, err := hub.HubMetaData.GetAbi()

	if err != nil {
		log.Panic("cannot get fdcAbi:", err)
	}

	requestEvent, ok := fdcAbi.Events["AttestationRequest"]

	if !ok {
		log.Panic("cannot get AttestationRequest event:", err)
	}

	attestationRequestEventSel = requestEvent.ID

}

type Runner struct {
	Protocol              uint64
	SubmitContractAddress common.Address
	FdcContractAddress    common.Address
	RelayContractAddress  common.Address
	DB                    *gorm.DB
	submit1Sel            [4]byte
	RoundManager          *attestation.Manager
}

const submit1FuncSelHex = "6c532fae"

func New(user config.UserConfigRaw, system config.SystemConfig) *Runner {
	// TODO: Luka - get these from config
	// CONSTANTS
	// requestEventSignature := "251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9"
	// signingPolicySignature := "91d0280e969157fc6c5b8f952f237b03d934b18534dafcac839075bbc33522f8"

	roundManager := attestation.NewManager(user)

	db, err := database.Connect(&user.DB)
	if err != nil {
		log.Panic("Could not connect to database:", err)
	}

	submit1FuncSel, err := parseFuncSel(submit1FuncSelHex)
	if err != nil {
		log.Panic("Could not parse submit1FuncSel:", err)
	}

	runner := Runner{
		Protocol:              system.Listener.Protocol,
		SubmitContractAddress: system.Listener.SubmitContractAddress,
		FdcContractAddress:    system.Listener.FdcContractAddress,
		RelayContractAddress:  system.Listener.RelayContractAddress,
		DB:                    db,
		submit1Sel:            submit1FuncSel,
		RoundManager:          roundManager,
	}

	return &runner
}

func parseFuncSel(sigInput string) ([4]byte, error) {
	var ret [4]byte
	inputBytes := []byte(sigInput)

	if hex.DecodedLen(len(inputBytes)) != 4 {
		return ret, errors.New("invalid length for function selector")
	}

	_, err := hex.Decode(ret[:], inputBytes)
	return ret, err
}

func (r *Runner) Run() {

	chooseTrigger := make(chan uint64)

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(r.DB, r.RelayContractAddress, 3)

	r.RoundManager.BitVotes = BitVoteListener(r.DB, r.FdcContractAddress, r.submit1Sel, r.Protocol, bitVoteBufferSize, chooseTrigger)

	r.RoundManager.Requests = AttestationRequestListener(r.DB, r.FdcContractAddress, requestsBufferSize, requestListenerInterval)

	go r.RoundManager.Run()

	prepareChooseTriggers(chooseTrigger, r.DB)

}

// prepareChooseTriggers tracks chain timestamps and makes sure that roundId of the round whose choose phase has just ended to the trigger chanel.
func prepareChooseTriggers(trigger chan uint64, DB *gorm.DB) {

	state, err := database.FetchState(DB)
	if err != nil {
		log.Panic("database error:", err)
	}

	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhaseEndPointers(state.BlockTimestamp)

	bitVoteTicker := time.NewTicker(time.Hour) // timer will be reset to 90 seconds

	go configureBitVoteTicker(bitVoteTicker, time.Unix(int64(*nextChoosePhaseEndTimestamp), 0), bitVoteHeadStart)

	for {

		ticker := time.NewTicker(databasePollTime)

		for {

			state, err := database.FetchState(DB)

			if err != nil {
				log.Error("database error:", err)
			} else {

				done := tryTriggerBitVote(
					nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp, state.BlockTimestamp, trigger,
				)

				if done {

					break
				}
			}
			<-ticker.C

		}

		<-bitVoteTicker.C
	}

}

// configureBitVoteTicker resets the ticker at headStart before start to roundLength
func configureBitVoteTicker(ticker *time.Ticker, start time.Time, headStart time.Duration) {

	time.Sleep(time.Until(start) - headStart)

	ticker.Reset(roundLength) // get this from config or constant

}

// tryTriggerBitVote checks whether the blockchain timestamp has surpassed end of choose phase or local time has surpassed it for more than bitVoteOffChainTriggerSeconds.
// If conditions are met, roundId is passed to the chanel c.
func tryTriggerBitVote(nextChoosePhaseRoundIDEnd *int, nextChoosePhaseEndTimestamp *uint64, currentBlockTime uint64, c chan uint64) bool {

	now := uint64(time.Now().Unix())

	if currentBlockTime > *nextChoosePhaseEndTimestamp {
		c <- uint64(*nextChoosePhaseRoundIDEnd)

		log.Infof("bitVote for round %d started with on-chain time", *nextChoosePhaseRoundIDEnd)

		*nextChoosePhaseRoundIDEnd, *nextChoosePhaseEndTimestamp = timing.NextChoosePhaseEnd(currentBlockTime)

		return true
	}

	if (now - bitVoteOffChainTriggerSeconds) > *nextChoosePhaseEndTimestamp {
		c <- uint64(*nextChoosePhaseRoundIDEnd)
		log.Infof("bitVote for round %d started with off-chain time", *nextChoosePhaseRoundIDEnd)

		*nextChoosePhaseRoundIDEnd++
		*nextChoosePhaseEndTimestamp = *nextChoosePhaseEndTimestamp + 90

		return true

	}

	return false
}

// BitVoteListener returns a channel that servers payload data submitted do submitContractAddress to method with funcSig for protocol.
// Payload for roundId is served whenever a trigger provides a roundId.
func BitVoteListener(
	db *gorm.DB,
	submitContractAddress common.Address,
	funcSel [4]byte,
	protocol uint64,
	bufferSize int,
	trigger <-chan uint64,
) <-chan payload.Round {

	out := make(chan payload.Round, bufferSize)

	go func() {

		for {
			roundId := <-trigger

			txs, err := database.FetchTransactionsByAddressAndSelectorTimestamp(
				db,
				submitContractAddress,
				funcSel,
				int64(timing.ChooseStartTimestamp(int(roundId))),
				int64(timing.ChooseEndTimestamp(int(roundId))),
			)
			if err != nil {
				log.Error("fetch txs error:", err)
				continue
			}

			bitVotes := []payload.Message{}
			for i := range txs {
				tx := &txs[i]
				payloads, err := payload.ExtractPayloads(tx)
				if err != nil {
					log.Error("extract payload error:", err)
					continue
				}

				bitVote, ok := payloads[protocol]
				if ok {
					bitVotes = append(bitVotes, bitVote)
				}

			}

			if len(bitVotes) > 0 {

				log.Infof("Received %d for round %d", len(bitVotes), roundId)

				out <- payload.Round{Messages: bitVotes, Id: roundId}
			} else {
				log.Infof("No bitVotes for round %d", roundId)
			}

		}

	}()

	return out

}

// AttestationRequestListener returns a channel that serves attestation requests events emitted by fdcContractAddress.
func AttestationRequestListener(
	db *gorm.DB, fdcContractAddress common.Address, bufferSize int, ListenerInterval time.Duration,
) <-chan []database.Log {

	out := make(chan []database.Log, bufferSize)

	go func() {

		trigger := time.NewTicker(ListenerInterval)

		_, startTimestamp := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))

		state, err := database.FetchState(db)
		if err != nil {
			log.Panic("fetch initial state error:", err)
		}

		lastQueriedBlock := state.Index

		logs, err := database.FetchLogsByAddressAndTopic0TimestampToBlockNumber(
			db, fdcContractAddress, attestationRequestEventSel, int64(startTimestamp), int64(state.Index),
		)
		if err != nil {
			log.Panic("fetch initial logs error")
		}

		if len(logs) > 0 {
			out <- logs
		}

		for {
			<-trigger.C

			state, err = database.FetchState(db)
			if err != nil {
				log.Error("fetch state error:", err)
				continue
			}

			logs, err := database.FetchLogsByAddressAndTopic0BlockNumber(
				db, fdcContractAddress, attestationRequestEventSel, int64(lastQueriedBlock), int64(state.Index),
			)
			if err != nil {
				log.Error("fetch logs error:", err)
				continue
			}

			lastQueriedBlock = state.Index

			if len(logs) > 0 {
				log.Debugf("Adding %d request logs to channel", len(logs))
				out <- logs
			}

		}

	}()

	return out
}

// SigningPolicyInitializedListener returns a channel that serves signingPolicyInitialized events emitted by relayContractAddress.
func SigningPolicyInitializedListener(db *gorm.DB, relayContractAddress common.Address, bufferSize int) <-chan []database.Log {
	out := make(chan []database.Log, bufferSize)

	go func() {
		logs, err := database.FetchLatestLogsByAddressAndTopic0(
			db, relayContractAddress, signingPolicyInitializedEventSel, 3,
		)

		latestQuery := time.Now()

		if err != nil {
			log.Panic("error fetching initial logs:", err)
		}

		log.Debug("Logs length:", len(logs))

		if len(logs) == 0 {
			log.Panic("No initial signing policies found:", err)
		}

		sorted := []database.Log{}

		// signingPolicyStorage expects policies in increasing order
		for i := range logs {
			sorted = append(sorted, logs[len(logs)-i-1])
		}

		out <- sorted

		spiTargetedListener(db, relayContractAddress, logs[0], latestQuery, out)

	}()

	return out

}

// spiTargetedListener that only starts aggressive queries for new signingPolicyInitialized events a bit before the expected emission and stops once it get one and waits until the next window.
func spiTargetedListener(
	db *gorm.DB,
	relayContractAddress common.Address,
	lastLog database.Log,
	latestQuery time.Time,
	out chan<- []database.Log,
) {
	lastSigningPolicy, err := attestation.ParseSigningPolicyInitializedLog(lastLog)
	if err != nil {
		log.Panic("error parsing initial logs:", err)
	}

	lastInitializedRewardEpochID := lastSigningPolicy.RewardEpochId.Uint64()

	for {
		expectedStartOfTheNextSigningPolicyInitialized := timing.ExpectedRewardEpochStartTimestamp(lastInitializedRewardEpochID + 1)

		log.Info(expectedStartOfTheNextSigningPolicyInitialized)

		untilStart := time.Until(time.Unix(int64(expectedStartOfTheNextSigningPolicyInitialized)-90*15, 0)) //use const for headStart 90*15

		log.Infof("next signing policy expected in %.1fh", untilStart.Hours())

		timer := time.NewTimer(untilStart)

		<-timer.C

		if err := queryNextSPI(db, relayContractAddress, latestQuery, out); err != nil {
			log.Error("error querying next SPI event:", err)
			continue
		}

		lastInitializedRewardEpochID++
	}
}

func queryNextSPI(
	db *gorm.DB,
	relayContractAddress common.Address,
	latestQuery time.Time,
	out chan<- []database.Log,
) error {
	ticker := time.NewTicker(89 * time.Second) // ticker that is guaranteed to tick at least once per SystemVotingRound

	for {
		now := time.Now()

		logs, err := database.FetchLogsByAddressAndTopic0Timestamp(
			db, relayContractAddress, signingPolicyInitializedEventSel, latestQuery.Unix(), now.Unix(),
		)

		latestQuery = now

		if err != nil {
			return err
		}

		if len(logs) > 0 {
			log.Debug("Adding signing policy to channel")
			out <- logs

			ticker.Stop()
			return nil
		}

		<-ticker.C
	}
}
