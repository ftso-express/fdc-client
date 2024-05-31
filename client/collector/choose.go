package collector

import (
	"context"
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// BitVoteListener returns a channel that servers payload data submitted do submitContractAddress to method with funcSig for protocol.
// Payload for roundId is served whenever a trigger provides a roundId.
func BitVoteListener(
	ctx context.Context,
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
			var roundId uint64

			select {
			case roundId = <-trigger:
				log.Debug("starting next BitVoteListener iteration")

			case <-ctx.Done():
				log.Info("BitVoteListener exiting:", ctx.Err())
				return
			}

			txs, err := database.FetchTransactionsByAddressAndSelectorTimestamp(
				ctx,
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

// prepareChooseTriggers tracks chain timestamps and makes sure that roundId of the round whose choose phase has just ended to the trigger chanel.
func prepareChooseTriggers(ctx context.Context, trigger chan uint64, db *gorm.DB) {
	state, err := database.FetchState(ctx, db)
	if err != nil {
		log.Panic("database error:", err)
	}

	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhaseEndPointers(state.BlockTimestamp)

	bitVoteTicker := time.NewTicker(time.Hour) // timer will be reset to 90 seconds

	go configureTicker(ctx, bitVoteTicker, time.Unix(int64(*nextChoosePhaseEndTimestamp), 0), bitVoteHeadStart)

	for {

		ticker := time.NewTicker(databasePollTime)

		for {
			state, err := database.FetchState(ctx, db)
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

			select {
			case <-ticker.C:
				log.Debug("starting next prepareChooseTriggers inner iteration")

			case <-ctx.Done():
				log.Info("prepareChooseTriggers exiting:", ctx.Err())
				return
			}

		}

		select {
		case <-bitVoteTicker.C:
			log.Debug("starting next prepareChooseTriggers outer iteration")

		case <-ctx.Done():
			log.Info("prepareChooseTriggers exiting:", ctx.Err())
		}
	}

}

// configureTicker resets the ticker at headStart before start to roundLength
func configureTicker(ctx context.Context, ticker *time.Ticker, start time.Time, headStart time.Duration) {
	select {
	case <-time.After(time.Until(start) - headStart):
		ticker.Reset(roundLength)

	case <-ctx.Done():
		return
	}
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
