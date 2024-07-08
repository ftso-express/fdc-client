package payload_test

import (
	"flare-common/database"
	"flare-common/payload"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var tx = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac66",
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

var txMultiple = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6665000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac66",
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

func TestExtractPayloads(t *testing.T) {

	tests := []struct {
		tx           *database.Transaction
		protocol     uint8
		nuOfPayloads int
		votingRound  uint64
		length       uint16
	}{
		{
			tx:           tx,
			protocol:     100,
			nuOfPayloads: 1,
			votingRound:  655419,
			length:       32,
		},
		{
			tx:           txMultiple,
			protocol:     101,
			nuOfPayloads: 2,
			votingRound:  655419,
			length:       32,
		},
	}

	for i, test := range tests {

		payloads, err := payload.ExtractPayloads(test.tx)

		require.NoError(t, err, fmt.Sprintf("error in test %d", i))

		require.Equal(t, test.nuOfPayloads, len(payloads), fmt.Sprintf("wrong number of payloads in test %d", i))

		payloadFTSO, ok := payloads[test.protocol]

		require.True(t, ok, fmt.Sprintf("missing payload in test %d", i))

		require.Equal(t, test.protocol, payloadFTSO.Protocol, fmt.Sprintf("wrong protocol id in test %d", i))

		require.Equal(t, test.votingRound, payloadFTSO.VotingRound, fmt.Sprintf("wrong voting round in test %d", i))

		require.Equal(t, test.length, payloadFTSO.Length)

	}

}

var txError1 = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac660", //too long
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

var txError2 = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6", //too short
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

var txError3 = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac664000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6", //too short
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

var txError4 = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6z", //illegal character
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

var txError5 = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b001043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac66", //wrong length
	BlockNumber:      16143116,
	BlockHash:        "40888ee23c4d7da30c42f826ea187386eac4564b02ce801f0b0b91ef1e71da05",
	TransactionIndex: 0,
	FromAddress:      "6bba3b6fb0dc902845666fdad70b2a814a57b6f3",
	ToAddress:        "2ca6571daa15ce734bbd0bf27d5c9d16787fc33f",
	Status:           1,
	Value:            "0",
	GasPrice:         "37500000000",
	Gas:              2500000,
	Timestamp:        1717417740,
}

func TestExtractPayloadsError(t *testing.T) {

	txs := []*database.Transaction{
		txError1,
		txError2,
		txError3,
		txError4,
		txError5,
	}

	for i, tx := range txs {

		_, err := payload.ExtractPayloads(tx)
		require.Error(t, err, fmt.Sprintf("error in test %d", i))

	}
}
