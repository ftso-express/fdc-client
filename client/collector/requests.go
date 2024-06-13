package collector

import (
	"context"
	"flare-common/database"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// AttestationRequestListener returns a channel that serves attestation requests events emitted by fdcContractAddress.
func AttestationRequestListener(
	ctx context.Context,
	db collectorDB,
	fdcContractAddress common.Address,
	bufferSize int,
	ListenerInterval time.Duration,
) <-chan []database.Log {

	out := make(chan []database.Log, bufferSize)

	go func() {

		trigger := time.NewTicker(ListenerInterval)

		_, startTimestamp := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))

		state, err := db.FetchState(ctx)
		if err != nil {
			log.Panic("fetch initial state error:", err)
		}

		lastQueriedBlock := state.Index

		logs, err := db.FetchLogsByAddressAndTopic0TimestampToBlockNumber(
			ctx, fdcContractAddress, attestationRequestEventSel, int64(startTimestamp), int64(state.Index),
		)
		if err != nil {
			log.Panic("fetch initial logs error")
		}

		if len(logs) > 0 {
			select {
			case out <- logs:

			case <-ctx.Done():
				log.Info("AttestationRequestListener exiting:", ctx.Err())
				return
			}
		}

		for {
			select {
			case <-trigger.C:
				log.Debug("starting next AttestationRequestListener iteration")

			case <-ctx.Done():
				log.Info("AttestationRequestListener exiting:", ctx.Err())
				return
			}

			state, err = db.FetchState(ctx)
			if err != nil {
				log.Error("fetch state error:", err)
				continue
			}

			logs, err := db.FetchLogsByAddressAndTopic0BlockNumber(
				ctx, fdcContractAddress, attestationRequestEventSel, int64(lastQueriedBlock), int64(state.Index),
			)
			if err != nil {
				log.Error("fetch logs error:", err)
				continue
			}

			lastQueriedBlock = state.Index

			if len(logs) > 0 {
				select {
				case out <- logs:
					log.Debugf("Added %d request logs to channel", len(logs))

				case <-ctx.Done():
					log.Info("AttestationRequestListener exiting:", ctx.Err())
					return
				}
			}

		}

	}()

	return out
}
