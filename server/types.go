package server

import (
	"github.com/flare-foundation/fdc-client/client/attestation"

	"github.com/ethereum/go-ethereum/common"
)

type DAResponseStatus string

const (
	Ok           DAResponseStatus = "OK"
	NotAvailable DAResponseStatus = "NOT_AVAILABLE"
)

type AttestationStatus string

const (
	Valid     AttestationStatus = "OK"
	WrongMIC  AttestationStatus = "WrongMIC"
	FailedLUT AttestationStatus = "FailedLUT"
	Failed    AttestationStatus = "FAILED"
	Error     AttestationStatus = "ERROR"
)

type DARequest struct {
	Request   string                 `json:"request"`
	Response  string                 `json:"response"`
	Status    AttestationStatus      `json:"status"`
	Consensus bool                   `json:"consensus"`
	Indexes   []attestation.IndexLog `json:"indexes"`
}

type DAAttestation struct {
	RoundID     uint32   `json:"roundId"`
	Request     string   `json:"request"`
	Response    string   `json:"response"`
	ResponseABI string   `json:"abi"`
	Proof       []string `json:"proof"`
	hash        common.Hash
}
