package attestation

import (
	"errors"
	"local/fdc/client/verification"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Status int

const (
	Unprocessed Status = iota
	UnsupportedPair
	Waiting
	Processing
	Success
	WrongMIC
	WrongLUT
	Retrying
	ProcessError
)

type Attestation struct {
	Index     uint64
	RoundID   uint64
	Request   verification.Request
	Response  verification.Response
	Fee       *big.Int
	Status    Status
	Consensus bool
	Hash      common.Hash
}

func (a *Attestation) VerifyResponse(r verification.Response) error {

	a.Response = r

	micReq, err := a.Request.GetMic()

	if err != nil {
		a.Status = ProcessError

		return errors.New("no mic in request")
	}
	micRes, err := r.ComputeMic()

	if err != nil {
		a.Status = ProcessError

		return errors.New("cannot compute mic")
	}

	if micReq != micRes {
		a.Status = WrongMIC
	}

	r.AddRound(a.RoundID)

	a.Hash, err = r.ComputeHash()

	if err != nil {
		a.Status = ProcessError

		return errors.New("cannot compute hash")
	}

	a.Status = Success

	return nil
}
