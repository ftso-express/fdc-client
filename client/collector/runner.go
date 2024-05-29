package collector

import (
	"flare-common/contracts/relay"
	"flare-common/database"
	"flare-common/logger"
	"flare-common/payload"
	"fmt"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"
	"strings"
	"time"

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

var signingPolicyInitializedEventSig string
var attestationRequestEventSig string
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

	signingPolicyInitializedEventSig = strings.TrimPrefix(signingPolicyEvent.ID.String(), "0x")

	fdcAbi, err := hub.HubMetaData.GetAbi()

	if err != nil {
		log.Panic("cannot get fdcAbi:", err)
	}

	requestEvent, ok := fdcAbi.Events["AttestationRequest"]

	if !ok {
		log.Panic("cannot get AttestationRequest event:", err)
	}

	attestationRequestEventSig = strings.TrimPrefix(requestEvent.ID.String(), "0x")

}

type Runner struct {
	Protocol              uint64
	SubmitContractAddress string
	FdcContractAddress    string
	RelayContractAddress  string
	DB                    *gorm.DB
	submit1Sig            string
	RoundManager          *attestation.Manager
}

func New(user config.UserConfigRaw, system config.SystemConfig) *Runner {
	// TODO: Luka - get these from config
	// CONSTANTS
	// requestEventSignature := "251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9"
	// signingPolicySignature := "91d0280e969157fc6c5b8f952f237b03d934b18534dafcac839075bbc33522f8"
	submit1FuncSig := "6c532fae"

	roundManager := attestation.NewManager(user)

	db, err := database.Connect(&user.DB)
	if err != nil {
		log.Panic("Could not connect to database:", err)
	}

	runner := Runner{
		Protocol:              system.Listener.Protocol,
		SubmitContractAddress: system.Listener.SubmitContractAddress,
		FdcContractAddress:    system.Listener.FdcContractAddress,
		RelayContractAddress:  system.Listener.RelayContractAddress,
		DB:                    db,
		submit1Sig:            submit1FuncSig,
		RoundManager:          roundManager,
	}

	return &runner
}

func (r *Runner) Run() {

	chooseTrigger := make(chan uint64)

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(r.DB, r.RelayContractAddress, 3)

	r.RoundManager.BitVotes = BitVoteListener(r.DB, r.FdcContractAddress, r.submit1Sig, r.Protocol, bitVoteBufferSize, chooseTrigger)

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
// Payload for roundID is served whenever a trigger provides a roundID.
func BitVoteListener(
	db *gorm.DB, submitContractAddress, funcSig string, protocol uint64, bufferSize int, trigger <-chan uint64,
) <-chan payload.Round {

	out := make(chan payload.Round, bufferSize)

	go func() {

		for {
			roundID := <-trigger

			txs, err := database.FetchTransactionsByAddressAndSelectorTimestamp(
				db,
				submitContractAddress,
				funcSig,
				int64(timing.ChooseStartTimestamp(int(roundID))),
				int64(timing.ChooseEndTimestamp(int(roundID))),
			)
			if err != nil {
				// TODO implement backoff/retry
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

				log.Infof("Received %d for round %d", len(bitVotes), roundID)

				out <- payload.Round{Messages: bitVotes, ID: roundID}
			} else {
				log.Infof("No bitVotes for round %d", roundID)
			}

		}

	}()

	return out

}

// AttestationRequestListener returns a channel that serves attestation requests events emitted by fdcContractAddress.
func AttestationRequestListener(
	db *gorm.DB, fdcContractAddress string, bufferSize int, ListenerInterval time.Duration,
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
			db, fdcContractAddress, attestationRequestEventSig, int64(startTimestamp), int64(state.Index),
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
				// TODO implement backoff/retry
				log.Error("fetch state error:", err)
				continue
			}

			logs, err := database.FetchLogsByAddressAndTopic0BlockNumber(
				db, fdcContractAddress, attestationRequestEventSig, int64(lastQueriedBlock), int64(state.Index),
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

// fetchLastSigningPolicyInitializedEvents returns last number of signingPolicyInitialized events emitted by relayContractAddress.
// The events are sorted by the timestamp descending.
func fetchLastSigningPolicyInitializedEvents(db *gorm.DB, relayContractAddress string, number int) ([]database.Log, error) {

	var logs []database.Log

	err := db.Where("address = ? AND topic0 = ?",
		strings.ToLower(strings.TrimPrefix(relayContractAddress, "0x")),
		strings.ToLower(strings.TrimPrefix(signingPolicyInitializedEventSig, "0x")),
	).Order("timestamp DESC").Limit(number).Find(&logs).Error

	if err != nil {
		return logs, fmt.Errorf("error fetching last sig policy logs: %s", err)
	}

	return logs, nil

}

// SigningPolicyInitializedListener returns a channel that serves signingPolicyInitialized events emitted by relayContractAddress.
func SigningPolicyInitializedListener(db *gorm.DB, relayContractAddress string, bufferSize int) <-chan []database.Log {
	out := make(chan []database.Log, bufferSize)

	go func() {

		logs, err := fetchLastSigningPolicyInitializedEvents(db, relayContractAddress, 3)

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
func spiTargetedListener(db *gorm.DB, relayContractAddress string, lastLog database.Log, latestQuery time.Time, out chan<- []database.Log) {

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

		ticker := time.NewTicker(89 * time.Second) // ticker that is guaranteed to tick at least once per SystemVotingRound

	aggressiveQuery:
		for {

			now := time.Now()

			logs, err := database.FetchLogsByAddressAndTopic0Timestamp(
				db, relayContractAddress, signingPolicyInitializedEventSig, latestQuery.Unix(), now.Unix(),
			)

			latestQuery = now

			if err != nil {
				log.Error("fetch logs error:", err)
				continue
			}

			if len(logs) > 0 {
				log.Debug("Adding signing policy to channel")
				out <- logs

				lastInitializedRewardEpochID++
				ticker.Stop()
				timer.Stop()
				break aggressiveQuery

			}

			<-ticker.C

		}

	}

}
