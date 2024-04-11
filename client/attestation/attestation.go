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
	Request   string
	Response  string
	Fee       *big.Int
	Status    Status
	Consensus bool
}
