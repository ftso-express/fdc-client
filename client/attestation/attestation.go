package attestation

import "math/big"

type Status int

const (
	Unprocessed Status = iota
	Processing
	Success
	WrongMIC
	WrongLUT
	Retrying
)

type Attestation struct {
	Index     uint64
	Request   string
	Response  string
	Fee       *big.Int
	Status    Status
	Consensus bool
}

func (at *Attestation) Verify() {
	at.Status = Unprocessed
}
