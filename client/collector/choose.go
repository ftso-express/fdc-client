package collector

import (
	"context"
	"flare-common/payload"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// BitVoteListener returns a channel that servers payloads data submitted do submitContractAddress to method with funcSig for protocol.
// Payloads for roundId are served whenever a trigger provides a roundId.
func BitVoteListener(
	ctx context.Context,
	db collectorDB,
	submitContractAddress common.Address,
	funcSel [4]byte,
	protocol uint8,
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

			txs, err := db.FetchTransactionsByAddressAndSelectorTimestamp(
				ctx,
				submitContractAddress,
				funcSel,
				int64(timing.ChooseStartTimestamp(roundId)),
				int64(timing.ChooseEndTimestamp(roundId)),
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
				log.Infof("Received %d bitVotes for round %d", len(bitVotes), roundId)

				select {
				case out <- payload.Round{Messages: bitVotes, Id: roundId}:
					log.Debugf("sent bitVotes for round %d", roundId)

				case <-ctx.Done():
					log.Info("BitVoteListener exiting")
					return
				}
			} else {
				log.Infof("No bitVotes for round %d", roundId)
			}

		}

	}()

	return out

}

// PrepareChooseTriggers tracks chain timestamps and passes roundId of the round whose choose phase has just ended to the trigger channel.
func PrepareChooseTriggers(ctx context.Context, trigger chan uint64, db collectorDB) {
	state, err := db.FetchState(ctx)
	if err != nil {
		log.Panic("database error:", err)
	}

	nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp := timing.NextChoosePhasePtr(state.BlockTimestamp)

	bitVoteTicker := time.NewTicker(time.Hour) // timer will be reset to 90 seconds

	go configureTicker(ctx, bitVoteTicker, time.Unix(int64(*nextChoosePhaseEndTimestamp), 0), bitVoteHeadStart)

	for {
		ticker := time.NewTicker(databasePollTime)

		for {
			state, err := db.FetchState(ctx)
			if err != nil {
				log.Error("database error:", err)
			} else {

				done := tryTriggerBitVote(
					ctx, nextChoosePhaseRoundIDEnd, nextChoosePhaseEndTimestamp, state.BlockTimestamp, trigger,
				)

				if done {
					break
				}
			}

			select {
			case <-ticker.C:

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

// configureTicker resets the ticker at headStart before start to collect phase duration.
func configureTicker(ctx context.Context, ticker *time.Ticker, start time.Time, headStart time.Duration) {
	select {
	case <-time.After(time.Until(start) - headStart):
		ticker.Reset(timing.CollectDurationSec * time.Second)

	case <-ctx.Done():
		return
	}
}

// tryTriggerBitVote checks whether the blockchain timestamp has surpassed the end of choose phase or local time has surpassed it for more than bitVoteOffChainTriggerSeconds.
// If conditions are met, roundId is passed to the channel c.
func tryTriggerBitVote(
	ctx context.Context,
	nextChoosePhaseRoundIDEnd *uint64,
	nextChoosePhaseEndTimestamp *uint64,
	currentBlockTime uint64,
	c chan uint64,
) bool {
	now := uint64(time.Now().Unix())

	if currentBlockTime > *nextChoosePhaseEndTimestamp {
		select {
		case c <- *nextChoosePhaseRoundIDEnd:
			log.Infof("bitVote for round %d started with on-chain time", *nextChoosePhaseRoundIDEnd)

		case <-ctx.Done():
			log.Info("tryTriggerBitVote exiting:", ctx.Err())
			return false
		}

		*nextChoosePhaseRoundIDEnd, *nextChoosePhaseEndTimestamp = timing.NextChooseEnd(currentBlockTime)

		return true
	}

	if (now - bitVoteOffChainTriggerSeconds) > *nextChoosePhaseEndTimestamp {
		select {
		case c <- uint64(*nextChoosePhaseRoundIDEnd):
			log.Infof("bitVote for round %d started with off-chain time", *nextChoosePhaseRoundIDEnd)

		case <-ctx.Done():
			log.Info("tryTriggerBitVote exiting:", ctx.Err())
			return false
		}

		*nextChoosePhaseRoundIDEnd++
		*nextChoosePhaseEndTimestamp = *nextChoosePhaseEndTimestamp + 90

		return true

	}

	return false
}
