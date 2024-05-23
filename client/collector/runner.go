package collector

import (
	"flare-common/contracts/relay"
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"
	"log"
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

var signingPolicyEventSig string
var requestEventSig string

func init() {

	relayAbi, err := relay.RelayMetaData.GetAbi()

	if err != nil {
		log.Panicf("cannot get relayAby: %s", err)
	}

	signingPolicyEventSig = relayAbi.Events["SigningPolicyInitialized"].ID.String()[2:] //remove 0x

	fdcAbi, err := hub.HubMetaData.GetAbi()

	if err != nil {
		log.Panicf("cannot get fdcAbi: %s", err)
	}

	requestEventSig = fdcAbi.Events["AttestationRequest"].ID.String()[2:] //remove 0x

}

type Runner struct {
	Protocol              uint64
	SubmitContractAddress string
	RequestEventSig       string
	FdcContractAddress    string
	RelayContractAddress  string
	SigningPolicyEventSig string
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
		log.Fatal("Could not connect to database:", err)
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

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(r.DB, r.RelayContractAddress, signingPolicyEventSig, 3)

	r.RoundManager.BitVotes = BitVoteInitializedListener(r.DB, r.FdcContractAddress, r.submit1Sig, r.Protocol, bitVoteBufferSize, chooseTrigger)

	r.RoundManager.Requests = RequestsInitializedListener(r.DB, r.FdcContractAddress, requestEventSig, requestsBufferSize, requestListenerInterval)

	go r.RoundManager.Run()

	prepareChooseTriggers(chooseTrigger, r.DB)

}

// prepareChooseTriggers tracks chain timestamps and makes sure that roundId of the round whose choose phase has just ended to the trigger chanel.
func prepareChooseTriggers(trigger chan uint64, DB *gorm.DB) {

	state, err := database.FetchState(DB)
	if err != nil {
		log.Fatal("database error:", err)
	}

	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhaseEnd(state.BlockTimestamp)

	bitVoteTicker := time.NewTicker(time.Hour) // timer will be reset to 90 seconds

	go configureBitVoteTicker(bitVoteTicker, time.Unix(int64(*nextChoosePhaseEndTimestamp), 0), bitVoteHeadStart)

	for {

		ticker := time.NewTicker(databasePollTime)

		for {

			state, err := database.FetchState(DB)

			if err != nil {
				log.Println("database error:", err)
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
		nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp = timing.NextChoosePhaseEnd(currentBlockTime)

		return true
	}

	if (now - bitVoteOffChainTriggerSeconds) > *nextChoosePhaseEndTimestamp {
		c <- uint64(*nextChoosePhaseRoundIDEnd)
		*nextChoosePhaseRoundIDEnd++
		*nextChoosePhaseEndTimestamp = *nextChoosePhaseEndTimestamp + 90

		return true

	}

	return false
}

// BitVoteInitializedListener returns an initialized channel that servers payload data submitted do submitContractAddress to method with funcSig for protocol.
// Payload for roundID is served whenever a trigger provides a roundID
func BitVoteInitializedListener(db *gorm.DB, submitContractAddress, funcSig string, protocol uint64, bufferSize int, trigger <-chan uint64) <-chan payload.Round {

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
				log.Println("fetch txs error:", err)
				continue
			}

			bitVotes := []payload.Message{}
			for i := range txs {
				tx := &txs[i]
				payloads, err := payload.ExtractPayloads(tx)
				if err != nil {
					log.Println("extract payload error:", err)
					continue
				}

				bitVote, ok := payloads[protocol]
				if ok {
					bitVotes = append(bitVotes, bitVote)
				}

			}

			if len(bitVotes) > 0 {
				out <- payload.Round{Messages: bitVotes, ID: roundID}
			}
		}

	}()

	return out

}

func RequestsInitializedListener(
	db *gorm.DB, fdcContractAddress, requestEventSig string, bufferSize int, ListenerInterval time.Duration,
) <-chan []database.Log {

	out := make(chan []database.Log, bufferSize)

	go func() {

		trigger := time.NewTicker(ListenerInterval)

		_, startTimestamp := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))

		state, err := database.FetchState(db)
		if err != nil {
			log.Fatal("fetch initial state error:", err)
		}

		lastQueriedBlock := state.Index

		logs, err := database.FetchLogsByAddressAndTopic0TimestampToBlockNumber(
			db, fdcContractAddress, requestEventSig, int64(startTimestamp), int64(state.Index),
		)
		if err != nil {
			log.Fatal("fetch initial logs error")
		}

		if len(logs) > 0 {
			out <- logs
		}

		for {
			<-trigger.C

			state, err = database.FetchState(db)
			if err != nil {
				// TODO implement backoff/retry
				log.Print("fetch state error:", err)
				continue
			}

			logs, err := database.FetchLogsByAddressAndTopic0BlockNumber(
				db, fdcContractAddress, requestEventSig, int64(lastQueriedBlock), int64(state.Index),
			)
			if err != nil {
				log.Print("fetch logs error:", err)
				continue
			}

			lastQueriedBlock = state.Index

			if len(logs) > 0 {
				log.Println("Added request to channel")
				log.Println(len(logs))
				out <- logs
			}

		}

	}()

	return out
}

func SigningPolicyInitializedListener(db *gorm.DB, relayContractAddress, signingPolicyInitializedEventSig string, bufferSize int) <-chan []database.Log {
	out := make(chan []database.Log, bufferSize)

	go func() {

		latestQuery := time.Now()
		twoWeeksBefore := latestQuery.Add(-2 * 7 * 24 * time.Hour)

		logs, err := database.FetchLogsByAddressAndTopic0Timestamp(
			db, relayContractAddress, signingPolicyInitializedEventSig, twoWeeksBefore.Unix(), latestQuery.Unix(),
		)
		if err != nil {
			log.Fatal("error fetching initial logs:", err)
		}

		log.Println("Logs length ", len(logs))

		out <- logs

		ticker := time.NewTicker(80 * time.Second) //TODO: optimize to reduce number of queries

		for {
			<-ticker.C

			now := time.Now()

			logs, err := database.FetchLogsByAddressAndTopic0Timestamp(
				db, relayContractAddress, signingPolicyInitializedEventSig, latestQuery.Unix(), now.Unix(),
			)
			if err != nil {
				log.Println("fetch logs error:", err)
				continue
			}

			latestQuery = now

			if len(logs) > 0 {
				log.Println("Added signing policy to channel")
				out <- logs
			}

		}

	}()

	return out

}
