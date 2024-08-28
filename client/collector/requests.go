package collector

import (
	"context"
	"flare-common/database"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// AttestationRequestListener initiates a channel that serves attestation requests events emitted by fdcContractAddress.
func AttestationRequestListener(
	ctx context.Context,
	db *gorm.DB,
	fdcContractAddress common.Address,
	listenerInterval time.Duration,
	logChan chan<- []database.Log,
) {
	trigger := time.NewTicker(listenerInterval)

	_, startTimestamp, err := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))
	if err != nil {
		log.Panic("time:", err)
	}

	state, err := database.FetchState(ctx, db, nil)
	if err != nil {
		log.Panic("fetch initial state error:", err)
	}

	lastQueriedBlock := state.Index

	params := database.LogsParams{
		Address: fdcContractAddress,
		Topic0:  AttestationRequestEventSel,
		From:    int64(startTimestamp),
		To:      int64(state.Index),
	}

	logs, err := database.FetchLogsByAddressAndTopic0FromTimestampToBlockNumber(
		ctx, db, params,
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

		state, err = database.FetchState(ctx, db, nil)
		if err != nil {
			log.Error("fetch state error:", err)
			continue
		}

		params := database.LogsParams{
			Address: fdcContractAddress,
			Topic0:  AttestationRequestEventSel,
			From:    int64(lastQueriedBlock),
			To:      int64(state.Index),
		}

		logs, err := database.FetchLogsByAddressAndTopic0BlockNumber(
			ctx, db, params,
		)
		if err != nil {
			log.Error("fetch logs error:", err)
			continue
		}

		lastQueriedBlock = state.Index

		if len(logs) > 0 {
			select {
			case logChan <- logs:
			case <-ctx.Done():
				log.Info("AttestationRequestListener exiting:", ctx.Err())
				return
			}
		}

	}
}
