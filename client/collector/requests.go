package collector

import (
	"context"
	"time"

	"github.com/flare-foundation/go-flare-common/pkg/database"
	"github.com/flare-foundation/go-flare-common/pkg/logger"

	"github.com/flare-foundation/fdc-client/client/timing"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// AttestationRequestListener initiates a channel that serves attestation request events emitted by fdcHub.
func AttestationRequestListener(
	ctx context.Context,
	db *gorm.DB,
	fdcHub common.Address,
	listenerInterval time.Duration,
	logChan chan<- []database.Log,
) {
	trigger := time.NewTicker(listenerInterval)

	// initial query
	_, startTimestamp, err := timing.LastCollectPhaseStart(uint64(time.Now().Unix()))
	if err != nil {
		logger.Panic("time:", err)
	}

	state, err := database.FetchState(ctx, db, nil)
	if err != nil {
		logger.Panic("fetch initial state error:", err)
	}

	lastQueriedBlock := state.Index

	params := database.LogsParams{
		Address: fdcHub,
		Topic0:  AttestationRequestEventSel,
		From:    int64(startTimestamp),
		To:      int64(state.Index),
	}

	logs, err := database.FetchLogsByAddressAndTopic0FromTimestampToBlockNumber(
		ctx, db, params,
	)
	if err != nil {
		logger.Panic("fetch initial logs error")
	}

	// add requests to the channel
	if len(logs) > 0 {
		select {
		case logChan <- logs:
		case <-ctx.Done():
			logger.Info("AttestationRequestListener exiting:", ctx.Err())
			return
		}
	}

	// infinite loop, making query once per listenerInterval from last queried block to the latest confirmed block in indexer db
	for {
		select {
		case <-trigger.C:
		case <-ctx.Done():
			logger.Info("AttestationRequestListener exiting:", ctx.Err())
			return
		}

		state, err = database.FetchState(ctx, db, nil)
		if err != nil {
			logger.Error("fetch state error:", err)
			continue
		}

		params := database.LogsParams{
			Address: fdcHub,
			Topic0:  AttestationRequestEventSel,
			From:    int64(lastQueriedBlock),
			To:      int64(state.Index),
		}

		logs, err := database.FetchLogsByAddressAndTopic0BlockNumber(
			ctx, db, params,
		)
		if err != nil {
			logger.Error("fetch logs error:", err)
			continue
		}
		logger.Debugf("Received %d requests from blocks (%d,%d]", len(logs), lastQueriedBlock, state.Index)

		lastQueriedBlock = state.Index

		// add requests to the channel
		if len(logs) > 0 {
			select {
			case logChan <- logs:
			case <-ctx.Done():
				logger.Info("AttestationRequestListener exiting:", ctx.Err())
				return
			}
		}
	}
}
