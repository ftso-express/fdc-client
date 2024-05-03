package collector

import (
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/attestation"
	"local/fdc/client/timing"
	"time"

	"gorm.io/gorm"
)

const (
	bitVoteBufferSize             = 10
	bitVoteOffChainTriggerSeconds = 15
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
	roundManager          *attestation.Manager
}

func (r *Runner) Run() {

	chooseTrigger := make(chan uint64)

	r.roundManager.BitVotes = BitVoteInitializedListener(r.DB, r.FdcContractAddress, r.submit1Sig, r.Protocol, bitVoteBufferSize, chooseTrigger)

	state, _ := database.FetchState(r.DB)
	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhaseEnd(state.BlockTimestamp)

	for {
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

	if (now - 15) > *nextChoosePhaseEndTimestamp {
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

			txs, _ := database.FetchTransactionsByAddressAndSelectorTimestamp(db, submitContractAddress, funcSig, int64(timing.GetChooseStartTimestamp(int(roundID))), int64(timing.GetChooseEndTimestamp(int(roundID))))

			bitVotes := []payload.Message{}
			for _, tx := range txs {
				payloads, _ := payload.ExtractPayloads(tx)
				bitVote := payloads[protocol]
				bitVotes = append(bitVotes, bitVote)

			}

			out <- payload.Round{bitVotes, roundID}
		}

	}()

	return out

}
