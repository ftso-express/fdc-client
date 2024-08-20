package server

import "github.com/ethereum/go-ethereum/common"

type PDPResponseStatus string

const (
	OK            PDPResponseStatus = "OK"
	NOT_AVAILABLE PDPResponseStatus = "NOT_AVAILABLE"
)

type PDPResponse struct {
	Status         PDPResponseStatus `json:"status"`
	Data           string            `json:"data"`
	AdditionalData string            `json:"additionalData"`
}

type merkleRootStorageObject struct {
	message          string
	merkleRoot       common.Hash
	randomNum        common.Hash
	consensusBitVote []byte
}

type RootsByAddress map[string]merkleRootStorageObject
