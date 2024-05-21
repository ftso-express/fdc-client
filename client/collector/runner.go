package collector

import (
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	"local/fdc/client/timing"
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
	requestEventSignature := "251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9"
	signingPolicySignature := "91d0280e969157fc6c5b8f952f237b03d934b18534dafcac839075bbc33522f8"
	submit1FuncSig := "6c532fae"

	roundManager := attestation.NewManager(user)

	db, err := database.Connect(&user.DB)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	runner := Runner{
		Protocol:              system.Listener.Protocol,
		SubmitContractAddress: system.Listener.SubmitContractAddress,
		RequestEventSig:       requestEventSignature,
		FdcContractAddress:    system.Listener.FdcContractAddress,
		RelayContractAddress:  system.Listener.RelayContractAddress,
		SigningPolicyEventSig: signingPolicySignature,
		DB:                    db,
		submit1Sig:            submit1FuncSig,
		RoundManager:          roundManager,
	}

	return &runner
}

const databasePollSeconds = 2

func (r *Runner) Run() {

	chooseTrigger := make(chan uint64)

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(r.DB, r.RelayContractAddress, r.SigningPolicyEventSig, 3)

	r.RoundManager.BitVotes = BitVoteInitializedListener(r.DB, r.FdcContractAddress, r.submit1Sig, r.Protocol, bitVoteBufferSize, chooseTrigger)

	r.RoundManager.Requests = RequestsInitializedListener(r.DB, r.FdcContractAddress, r.RequestEventSig, requestsBufferSize, requestListenerInterval)

	state, err := database.FetchState(r.DB)
	if err != nil {
		log.Fatal("database error:", err)
	}

	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhaseEnd(state.BlockTimestamp)

	go r.RoundManager.Run()

	ticker := time.NewTicker(databasePollSeconds * time.Second)

	for {
		state, err := database.FetchState(r.DB)
		if err != nil {
			log.Println("database error:", err)
		} else {
			tryTriggerBitVote(
				nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp, state.BlockTimestamp, chooseTrigger,
			)
		}

		<-ticker.C
	}

}

func tryTriggerBitVote(nextChoosePhaseRoundIDEnd *int, nextChoosePhaseEndTimestamp *uint64, currentBlockTime uint64, c chan uint64) {

	now := uint64(time.Now().Unix())

	if currentBlockTime > *nextChoosePhaseEndTimestamp {
		c <- uint64(*nextChoosePhaseRoundIDEnd)
		nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp = timing.NextChoosePhaseEnd(currentBlockTime)
	}

	if (now - bitVoteOffChainTriggerSeconds) > *nextChoosePhaseEndTimestamp {
		c <- uint64(*nextChoosePhaseRoundIDEnd)
		*nextChoosePhaseRoundIDEnd++
		*nextChoosePhaseEndTimestamp = *nextChoosePhaseEndTimestamp + 90

	}

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
