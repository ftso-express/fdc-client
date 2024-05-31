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

func FetchLatestLogsByAddressAndTopic0(
	ctx context.Context, db *gorm.DB, address common.Address, topic0 common.Hash, number int,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLatestLogsByAddressAndTopic0(ctx, db, address, topic0, number)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func fetchLatestLogsByAddressAndTopic0(
	ctx context.Context, db *gorm.DB, address common.Address, topic0 common.Hash, number int,
) ([]Log, error) {
	var logs []Log

	err := db.WithContext(ctx).Where("address = ? AND topic0 = ?",
		hex.EncodeToString(address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(topic0[:]),
	).Order("timestamp DESC").Limit(number).Find(&logs).Error

	return logs, err
}

// Fetch all logs matching address and topic0 from timestamp range (from, to], order by timestamp
func FetchLogsByAddressAndTopic0Timestamp(
	ctx context.Context,
	db *gorm.DB,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLogsByAddressAndTopic0Timestamp(
				ctx, db, address, topic0, from, to,
			)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func fetchLogsByAddressAndTopic0Timestamp(
	ctx context.Context,
	db *gorm.DB,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]Log, error) {
	var logs []Log
	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND timestamp > ? AND timestamp <= ?",
		hex.EncodeToString(address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(topic0[:]),
		from, to,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Fetch all logs matching address and topic0 from timestamp to block number, order by timestamp
func FetchLogsByAddressAndTopic0TimestampToBlockNumber(
	ctx context.Context,
	db *gorm.DB,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLogsByAddressAndTopic0TimestampToBlockNumber(
				ctx, db, address, topic0, from, to,
			)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func fetchLogsByAddressAndTopic0TimestampToBlockNumber(
	ctx context.Context,
	db *gorm.DB,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]Log, error) {
	var logs []Log

	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND timestamp >= ? AND block_number <= ?",
		hex.EncodeToString(address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(topic0[:]),
		from, to,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Fetch all logs matching address and topic0 from block range (from, to], order by timestamp
func FetchLogsByAddressAndTopic0BlockNumber(
	ctx context.Context,
	db *gorm.DB,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLogsByAddressAndTopic0BlockNumber(
				ctx, db, address, topic0, from, to,
			)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func fetchLogsByAddressAndTopic0BlockNumber(
	ctx context.Context,
	db *gorm.DB,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]Log, error) {
	var logs []Log

	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND block_number > ? AND block_number <= ?",
		hex.EncodeToString(address[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(topic0[:]),
		from, to,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Fetch all transactions matching toAddress and functionSel from timestamp range (from, to], order by timestamp
func FetchTransactionsByAddressAndSelectorTimestamp(
	ctx context.Context,
	db *gorm.DB,
	toAddress common.Address,
	functionSel [4]byte,
	from, to int64,
) ([]Transaction, error) {
	var txs []Transaction

	err := backoff.RetryNotify(
		func() error {
			var err error
			txs, err = fetchTransactionsByAddressAndSelectorTimestamp(
				ctx, db, toAddress, functionSel, from, to,
			)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching transactions: %v, retrying after %v", err, duration)
		},
	)

	return txs, err
}

func fetchTransactionsByAddressAndSelectorTimestamp(
	ctx context.Context,
	db *gorm.DB,
	toAddress common.Address,
	functionSel [4]byte,
	from, to int64,
) ([]Transaction, error) {
	var transactions []Transaction

	err := db.WithContext(ctx).Where(
		"to_address = ? AND function_sig = ? AND timestamp > ? AND timestamp <= ?",
		hex.EncodeToString(toAddress[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(functionSel[:]),
		from, to,
	).Order("timestamp").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// Fetch all transactions matching toAddress and functionSel from block number range (from, to], order by timestamp
func FetchTransactionsByAddressAndSelectorBlockNumber(
	ctx context.Context,
	db *gorm.DB,
	toAddress common.Address,
	functionSel [4]byte,
	from, to int64,
) ([]Transaction, error) {
	var txs []Transaction

	err := backoff.RetryNotify(
		func() error {
			var err error
			txs, err = fetchTransactionsByAddressAndSelectorBlockNumber(
				ctx, db, toAddress, functionSel, from, to,
			)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching transactions: %v, retrying after %v", err, duration)
		},
	)

	return txs, err
}

func fetchTransactionsByAddressAndSelectorBlockNumber(
	ctx context.Context,
	db *gorm.DB,
	toAddress common.Address,
	functionSel [4]byte,
	from, to int64,
) ([]Transaction, error) {
	var transactions []Transaction

	err := db.WithContext(ctx).Where(
		"to_address = ? AND function_sig = ? AND block_number > ? AND block_number <= ?",
		hex.EncodeToString(toAddress[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(functionSel[:]),
		from, to,
	).Order("timestamp").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func FetchState(ctx context.Context, db *gorm.DB) (State, error) {
	var state State

	err := backoff.RetryNotify(
		func() error {
			var err error
			state, err = fetchState(ctx, db)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching state: %v, retrying after %v", err, duration)
		},
	)

	return state, err
}

func fetchState(ctx context.Context, db *gorm.DB) (State, error) {
	var states []State

	err := db.WithContext(ctx).Order("block_timestamp DESC").Find(&states).Error

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
