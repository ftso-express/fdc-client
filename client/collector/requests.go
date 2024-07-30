package collector

import (
	"context"
	"flare-common/database"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// AttestationRequestListener initiates a channel that serves attestation requests events emitted by fdcContractAddress.
func AttestationRequestListener(
	ctx context.Context,
	db collectorDB,
	fdcContractAddress common.Address,
	ListenerInterval time.Duration,
	logChan chan<- []database.Log,
) {
	trigger := time.NewTicker(ListenerInterval)

	_, startTimestamp, err := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))
	if err != nil {
		log.Panic("time:", err)
	}

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
		case logChan <- logs:
		case <-ctx.Done():
			log.Info("AttestationRequestListener exiting:", ctx.Err())
			return
		}
	}

	for {
		select {
		case <-trigger.C:
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
			case logChan <- logs:
				log.Debugf("Added %d request logs to channel", len(logs))
			case <-ctx.Done():
				log.Info("AttestationRequestListener exiting:", ctx.Err())
				return
			}
		}

	}
}
