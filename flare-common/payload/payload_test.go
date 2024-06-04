package payload_test

import (
	"flare-common/database"
	"flare-common/payload"
	"testing"
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

func TestExtractPayloads(t *testing.T) {

	payloadsOne, err := payload.ExtractPayloads(tx)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(payloadsOne) != 1 {
		t.Errorf("to many payloads")
	}

	payloadFTSO, ok := payloadsOne[100]

	if !ok {
		t.Error("no payload for FTSO")
	}

	if payloadFTSO.Protocol != 100 {
		t.Error("mismatching protocol")
	}

	if payloadFTSO.VotingRound != 655419 {
		t.Errorf("wrong voting round %d", payloadFTSO.VotingRound)

	}

	if payloadFTSO.Length != 32 {
		t.Errorf("wrong  length %d", payloadFTSO.Length)

	}

}

var txError1 = &database.Transaction{
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c7327",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac660",
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
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6",
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
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac664000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6",
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
	Hash:             "8dd67e88aa6f863aeb5cd62874530efd7dafef2d4a8cdf7fbf71844dab1c732z",
	FunctionSig:      "6c532fae",
	Input:            "6c532fae64000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac664000a003b002043a94d3c612d7f5cfd65e53a06d55bac77abbd2a6eb4dff766f51092db36ac6",
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

	_, err := payload.ExtractPayloads(txError1)

	if err == nil {
		t.Errorf("no error 1: %s", err)
	}

	_, err = payload.ExtractPayloads(txError2)

	if err == nil {
		t.Errorf("no error 2: %s", err)
	}

	_, err = payload.ExtractPayloads(txError3)

	if err == nil {
		t.Errorf("no error 3: %s", err)
	}

	_, err = payload.ExtractPayloads(txError4)

	if err == nil {
		t.Errorf("no error 4: %s", err)
	}
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

func TestExtractPayloadsMultiple(t *testing.T) {

	payloadsMultiple, err := payload.ExtractPayloads(txMultiple)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(payloadsMultiple) != 2 {
		t.Errorf("to many payloads")
	}

	payloadFTSO, ok := payloadsMultiple[101]

	if !ok {
		t.Error("no payload for FTSO")
	}

	if payloadFTSO.Protocol != 101 {
		t.Error("mismatching protocol")
	}

	if payloadFTSO.VotingRound != 655419 {
		t.Errorf("wrong voting round %d", payloadFTSO.VotingRound)

	}

	if payloadFTSO.Length != 32 {
		t.Errorf("wrong  length %d", payloadFTSO.Length)

	}

}
