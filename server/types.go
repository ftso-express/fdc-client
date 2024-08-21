package server

import (
	"local/fdc/client/attestation"

	"github.com/ethereum/go-ethereum/common"
)

type ResponseStatus string

const (
	Ok           ResponseStatus = "OK"
	NotAvailable ResponseStatus = "NOT_AVAILABLE"
)

type AttestationStatus string

const (
	Valid  AttestationStatus = "OK"
	Failed AttestationStatus = "FAILED"
	Error  AttestationStatus = "ERROR"
)

type PDPResponse struct {
	Status         ResponseStatus `json:"status"`
	Data           string         `json:"data"`
	AdditionalData string         `json:"additionalData"`
}

type merkleRootStorageObject struct {
	message          string
	merkleRoot       common.Hash
	randomNum        common.Hash
	consensusBitVote []byte
}

type RootsByAddress map[string]merkleRootStorageObject

type DARequest struct {
	Request   string                 `json:"request"`
	Response  string                 `json:"response"`
	Status    AttestationStatus      `json:"status"`
	Consensus bool                   `json:"consensus"`
	Indexes   []attestation.IndexLog `json:"indexes"`
}

type DAAttestation struct {
	RoundId  uint64   `json:"roundId"`
	Request  string   `json:"request"`
	Response string   `json:"response"`
	Proof    []string `json:"proof"`
	hash     common.Hash
}
