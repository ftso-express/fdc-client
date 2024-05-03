package attestation

import (
	"errors"
	"flare-common/database"
	"flare-common/events"
	"local/fdc/client/verification"
	hub "local/fdc/contracts/FDC"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Status int

const (
	Unprocessed     Status = iota
	UnsupportedPair Status = 2
	Waiting         Status = 3
	Processing      Status = 4
	Success         Status = 5
	WrongMIC        Status = 6
	WrongLUT        Status = 7
	Retrying        Status = 8
	ProcessError    Status = 9
)

type Index struct {
	BlockNumber uint64
	LogIndex    uint64
}

func Less(a, b Index) bool {
	if a.BlockNumber < b.BlockNumber {
		return true
	}
	if a.BlockNumber == b.BlockNumber && a.LogIndex < b.LogIndex {
		return true
	}

	return false

}

type Attestation struct {
	Index     Index
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

func ParseAttestationRequestLog(hub *hub.Hub, dbLog database.Log) (*hub.HubAttestationRequest, error) {
	contractLog, err := events.ConvertDatabaseLogToChainLog(dbLog)
	if err != nil {
		return nil, err
	}
	return hub.HubFilterer.ParseAttestationRequest(*contractLog)
}
