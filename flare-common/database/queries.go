package database

import (
	"errors"
	"flare-common/logger"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/gorm"
)

var log = logger.GetLogger()

// Fetch all logs matching address and topic0 from timestamp range (from, to], order by timestamp
func FetchLogsByAddressAndTopic0Timestamp(
	db *gorm.DB, address string, topic0 string, from int64, to int64,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLogsByAddressAndTopic0Timestamp(
				db, address, topic0, from, to,
			)
			return err
		},
		backoff.NewExponentialBackOff(),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err

}

func fetchLogsByAddressAndTopic0Timestamp(
	db *gorm.DB, address string, topic0 string, from int64, to int64,
) ([]Log, error) {
	var logs []Log
	err := db.Where(
		"address = ? AND topic0 = ? AND timestamp > ? AND timestamp <= ?",
		strings.ToLower(strings.TrimPrefix(address, "0x")),
		strings.ToLower(strings.TrimPrefix(topic0, "0x")),
		from, to,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Fetch all logs matching address and topic0 from timestamp to block number, order by timestamp
func FetchLogsByAddressAndTopic0TimestampToBlockNumber(
	db *gorm.DB, address string, topic0 string, from int64, to int64,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLogsByAddressAndTopic0TimestampToBlockNumber(
				db, address, topic0, from, to,
			)
			return err
		},
		backoff.NewExponentialBackOff(),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func fetchLogsByAddressAndTopic0TimestampToBlockNumber(
	db *gorm.DB, address string, topic0 string, from int64, to int64,
) ([]Log, error) {
	var logs []Log
	err := db.Where(
		"address = ? AND topic0 = ? AND timestamp >= ? AND block_number <= ?",
		strings.ToLower(strings.TrimPrefix(address, "0x")),
		strings.ToLower(strings.TrimPrefix(topic0, "0x")),
		from, to,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Fetch all logs matching address and topic0 from block range (from, to], order by timestamp
func FetchLogsByAddressAndTopic0BlockNumber(
	db *gorm.DB, address string, topic0 string, from int64, to int64,
) ([]Log, error) {
	var logs []Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchLogsByAddressAndTopic0BlockNumber(
				db, address, topic0, from, to,
			)
			return err
		},
		backoff.NewExponentialBackOff(),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func fetchLogsByAddressAndTopic0BlockNumber(
	db *gorm.DB, address string, topic0 string, from int64, to int64,
) ([]Log, error) {
	var logs []Log

	err := db.Where(
		"address = ? AND topic0 = ? AND block_number > ? AND block_number <= ?",
		strings.ToLower(strings.TrimPrefix(address, "0x")),
		strings.ToLower(strings.TrimPrefix(topic0, "0x")),
		from, to,
	).Order("timestamp").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Fetch all transactions matching toAddress and functionSig from timestamp range (from, to], order by timestamp
func FetchTransactionsByAddressAndSelectorTimestamp(
	db *gorm.DB, toAddress string, functionSig string, from int64, to int64,
) ([]Transaction, error) {
	var txs []Transaction

	err := backoff.RetryNotify(
		func() error {
			var err error
			txs, err = fetchTransactionsByAddressAndSelectorTimestamp(
				db, toAddress, functionSig, from, to,
			)
			return err
		},
		backoff.NewExponentialBackOff(),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching transactions: %v, retrying after %v", err, duration)
		},
	)

	return txs, err
}

func fetchTransactionsByAddressAndSelectorTimestamp(
	db *gorm.DB, toAddress string, functionSig string, from int64, to int64,
) ([]Transaction, error) {
	var transactions []Transaction

	err := db.Where(
		"to_address = ? AND function_sig = ? AND timestamp > ? AND timestamp <= ?",
		strings.ToLower(strings.TrimPrefix(toAddress, "0x")),
		strings.ToLower(strings.TrimPrefix(functionSig, "0x")),
		from, to,
	).Order("timestamp").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// Fetch all transactions matching toAddress and functionSig from block number range (from, to], order by timestamp
func FetchTransactionsByAddressAndSelectorBlockNumber(
	db *gorm.DB, toAddress string, functionSig string, from int64, to int64,
) ([]Transaction, error) {
	var txs []Transaction

	err := backoff.RetryNotify(
		func() error {
			var err error
			txs, err = fetchTransactionsByAddressAndSelectorBlockNumber(
				db, toAddress, functionSig, from, to,
			)
			return err
		},
		backoff.NewExponentialBackOff(),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching transactions: %v, retrying after %v", err, duration)
		},
	)

	return txs, err
}

func fetchTransactionsByAddressAndSelectorBlockNumber(
	db *gorm.DB, toAddress string, functionSig string, from int64, to int64,
) ([]Transaction, error) {
	var transactions []Transaction

	err := db.Where(
		"to_address = ? AND function_sig = ? AND block_number > ? AND block_number <= ?",
		strings.ToLower(strings.TrimPrefix(toAddress, "0x")),
		strings.ToLower(strings.TrimPrefix(functionSig, "0x")),
		from, to,
	).Order("timestamp").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func FetchState(db *gorm.DB) (State, error) {
	var state State

	err := backoff.RetryNotify(
		func() error {
			var err error
			state, err = fetchState(db)
			return err
		},
		backoff.NewExponentialBackOff(),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching state: %v, retrying after %v", err, duration)
		},
	)

	return state, err
}

func fetchState(db *gorm.DB) (State, error) {
	var states []State

	err := db.Order("block_timestamp DESC").Find(&states).Error

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
