package collector

import (
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/attestation"
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

func New() *Runner {
	// TODO: Luka - get these from config
	// CONSTANTS
	requestEventSignature := "251377668af6553101c9bb094ba89c0c536783e005e203625e6cd57345918cc9"
	signingPolicySignature := "91d0280e969157fc6c5b8f952f237b03d934b18534dafcac839075bbc33522f8"
	submissionAddress := "2cA6571Daa15ce734Bbd0Bf27D5C9D16787fc33f"
	fdcContractAddress := "Cf6798810Bc8C0B803121405Fee2A5a9cc0CA5E5"
	relayAddress := "32D46A1260BB2D8C9d5Ab1C9bBd7FF7D7CfaabCC"
	submit1FuncSig := "6c532fae"

	roundManager := attestation.NewManager()

	config := database.DBConfig{
		Host:       "localhost",
		Port:       3306,
		Database:   "flare_ftso_indexer",
		Username:   "root",
		Password:   "root",
		LogQueries: false,
	}
	db, err := database.Connect(&config)
	if err != nil {
		log.Println("Could not connect to database")
		panic(err)
	}

	runner := Runner{
		Protocol:              200,
		SubmitContractAddress: submissionAddress,
		RequestEventSig:       requestEventSignature,
		FdcContractAddress:    fdcContractAddress,
		RelayContractAddress:  relayAddress,
		SigningPolicyEventSig: signingPolicySignature,
		DB:                    db,
		submit1Sig:            submit1FuncSig,
		RoundManager:          roundManager,
	}

	return &runner
}

func (r *Runner) Run() {

	chooseTrigger := make(chan uint64)

	r.RoundManager.SigningPolicies = SigningPolicyInitializedListener(r.DB, r.RelayContractAddress, r.SigningPolicyEventSig, 3)

	r.RoundManager.BitVotes = BitVoteInitializedListener(r.DB, r.FdcContractAddress, r.submit1Sig, r.Protocol, bitVoteBufferSize, chooseTrigger)

	r.RoundManager.Requests = RequestsInitializedListener(r.DB, r.FdcContractAddress, r.RequestEventSig, requestsBufferSize, requestListenerInterval)

	state, _ := database.FetchState(r.DB)
	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhaseEnd(state.BlockTimestamp)

	go r.RoundManager.Run()

	for {

		time.Sleep(2 * time.Second)
		state, _ := database.FetchState(r.DB)
		tryTriggerBitVote(nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp, state.BlockTimestamp, chooseTrigger)

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

	// TODO: handle errors

	out := make(chan payload.Round, bufferSize)

	go func() {

		for {
			roundID := <-trigger

			txs, _ := database.FetchTransactionsByAddressAndSelectorTimestamp(db, submitContractAddress, funcSig, int64(timing.ChooseStartTimestamp(int(roundID))), int64(timing.ChooseEndTimestamp(int(roundID))))

			bitVotes := []payload.Message{}
			for _, tx := range txs {
				payloads, _ := payload.ExtractPayloads(tx)
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

func RequestsInitializedListener(db *gorm.DB, fdcContractAddress, requestEventSig string, bufferSize int, ListenerInterval time.Duration) <-chan []database.Log {

	out := make(chan []database.Log, bufferSize)

	go func() {

		trigger := time.NewTicker(ListenerInterval)

		_, startTimestamp := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))

		state, _ := database.FetchState(db)

		lastQueriedBlock := state.Index

		logs, _ := database.FetchLogsByAddressAndTopic0TimestampToBlockNumber(db, fdcContractAddress, requestEventSig, int64(startTimestamp), int64(state.Index))

		if len(logs) > 0 {

			out <- logs
		}

		for {
			<-trigger.C

			state, _ = database.FetchState(db)

			logs, _ := database.FetchLogsByAddressAndTopic0BlockNumber(db, fdcContractAddress, requestEventSig, int64(lastQueriedBlock), int64(state.Index))

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

		logs, _ := database.FetchLogsByAddressAndTopic0Timestamp(db, relayContractAddress, signingPolicyInitializedEventSig, twoWeeksBefore.Unix(), latestQuery.Unix())

		log.Println("Logs length ", len(logs))

		out <- logs

		ticker := time.NewTicker(80 * time.Second) //TODO: optimize to reduce number of queries

		for {
			<-ticker.C

			now := time.Now()

			logs, _ := database.FetchLogsByAddressAndTopic0Timestamp(db, relayContractAddress, signingPolicyInitializedEventSig, latestQuery.Unix(), now.Unix())

			latestQuery = now

			if len(logs) > 0 {
				log.Println("Added signing policy to channel")
				out <- logs
			}

		}

	}()

	return out

}
