package database

import (
	"context"
	"encoding/hex"
	"errors"
	"flare-common/logger"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

var log = logger.GetLogger()

type LatestLogsParams struct {
	Address common.Address
	Topic0  common.Hash
	Number  int
}

func FetchLatestLogsByAddressAndTopic0(
	ctx context.Context, db *gorm.DB, params LatestLogsParams,
) ([]Log, error) {
	return RetryWrapper(fetchLatestLogsByAddressAndTopic0, "fetching logs")(ctx, db, params)
}

func fetchLatestLogsByAddressAndTopic0(
	ctx context.Context, db *gorm.DB, params LatestLogsParams,
) ([]Log, error) {
	var logs []Log

	err := db.WithContext(ctx).Where("address = ? AND topic0 = ?",
		hex.EncodeToString(params.Address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(params.Topic0[:]),
	).Order("timestamp DESC").Limit(params.Number).Find(&logs).Error

	return logs, err
}

type LogsParams struct {
	Address  common.Address
	Topic0   common.Hash
	From, To int64 //blockNumber or timestamp depending on the function
}

// Fetch all logs matching address and topic0 from timestamp range (from, to], order by timestamp
func FetchLogsByAddressAndTopic0Timestamp(ctx context.Context, db *gorm.DB, params LogsParams) ([]Log, error) {
	return RetryWrapper(fetchLogsByAddressAndTopic0Timestamp, "fetching logs")(ctx, db, params)
}

func fetchLogsByAddressAndTopic0Timestamp(ctx context.Context, db *gorm.DB, params LogsParams) ([]Log, error) {
	var logs []Log
	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND timestamp > ? AND timestamp <= ?",
		hex.EncodeToString(params.Address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(params.Topic0[:]),
		params.From,
		params.To,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Fetch all logs matching address and topic0 from timestamp (included) to block number (included), order by timestamp
func FetchLogsByAddressAndTopic0FromTimestampToBlockNumber(ctx context.Context, db *gorm.DB, params LogsParams) ([]Log, error) {
	return RetryWrapper(fetchLogsByAddressAndTopic0FromTimestampToBlockNumber, "fetching logs")(ctx, db, params)
}

func fetchLogsByAddressAndTopic0FromTimestampToBlockNumber(ctx context.Context, db *gorm.DB, params LogsParams) ([]Log, error) {
	var logs []Log

	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND timestamp >= ? AND block_number <= ?",
		hex.EncodeToString(params.Address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(params.Topic0[:]),
		params.From,
		params.To,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Fetch all logs matching address and topic0 from block range (from, to], order by timestamp
func FetchLogsByAddressAndTopic0BlockNumber(ctx context.Context, db *gorm.DB, params LogsParams) ([]Log, error) {
	return RetryWrapper(fetchLogsByAddressAndTopic0BlockNumber, "fetching logs")(ctx, db, params)
}

func fetchLogsByAddressAndTopic0BlockNumber(ctx context.Context, db *gorm.DB, params LogsParams) ([]Log, error) {
	var logs []Log

	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND block_number > ? AND block_number <= ?",
		hex.EncodeToString(params.Address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(params.Topic0[:]),
		params.From,
		params.To,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

type TxParams struct {
	ToAddress   common.Address
	FunctionSel [4]byte
	From, To    int64 //blockNumber or timestamp depending on the function
}

// Fetch all transactions matching toAddress and functionSel from timestamp range (from, to], order by timestamp
func FetchTransactionsByAddressAndSelectorTimestamp(ctx context.Context, db *gorm.DB, params TxParams) ([]Transaction, error) {
	return RetryWrapper(fetchTransactionsByAddressAndSelectorTimestamp, "fetching transactions")(ctx, db, params)
}

func fetchTransactionsByAddressAndSelectorTimestamp(ctx context.Context, db *gorm.DB, params TxParams) ([]Transaction, error) {
	var transactions []Transaction

	err := db.WithContext(ctx).Where(
		"to_address = ? AND function_sig = ? AND timestamp > ? AND timestamp <= ?",
		hex.EncodeToString(params.ToAddress[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(params.FunctionSel[:]),
		params.From,
		params.To,
	).Order("timestamp").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// Fetch all transactions matching toAddress and functionSel from block number range (from, to], order by timestamp
func FetchTransactionsByAddressAndSelectorBlockNumber(ctx context.Context, db *gorm.DB, params TxParams) ([]Transaction, error) {
	return RetryWrapper(fetchTransactionsByAddressAndSelectorBlockNumber, "fetching transactions")(ctx, db, params)
}

func fetchTransactionsByAddressAndSelectorBlockNumber(ctx context.Context, db *gorm.DB, params TxParams) ([]Transaction, error) {
	var transactions []Transaction

	err := db.WithContext(ctx).Where(
		"to_address = ? AND function_sig = ? AND block_number > ? AND block_number <= ?",
		hex.EncodeToString(params.ToAddress[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(params.FunctionSel[:]),
		params.From,
		params.To,
	).Order("timestamp").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// FetchState returns the state of the indexer database
func FetchState(ctx context.Context, db *gorm.DB, _ any) (State, error) {
	return RetryWrapper(fetchState, "fetching, state")(ctx, db, nil)

}

func fetchState(ctx context.Context, db *gorm.DB, _ any) (State, error) {
	var states []State

	err := db.WithContext(ctx).Where(
		"name = ?",
		"last_database_block",
	).Find(&states).Error

	if err != nil {
		var state State
		return state, err
	}

	if len(states) == 0 {
		var state State
		return state, errors.New("no states in database")

	}

	return states[0], nil
}

func RetryWrapper[F any, P any](query func(context.Context, *gorm.DB, P) (F, error), errorMsg string) func(context.Context, *gorm.DB, P) (F, error) {
	wrappedFunc := func(ctx context.Context, db *gorm.DB, params P) (F, error) {
		var returnValue F

		err := backoff.RetryNotify(
			func() error {
				var err error
				returnValue, err = query(ctx, db, params)
				return err
			},
			backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
			func(err error, duration time.Duration) {
				log.Errorf("error %s: %v, retrying after %v", errorMsg, err, duration)
			},
		)

		return returnValue, err
	}

	return wrappedFunc
}
