package attestation

import (
	"errors"
	"flare-common/contacts/relay"
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

type IndexLog struct {
	BlockNumber uint64
	LogIndex    uint64
}

// lessLog returns true if a has lower blockNumber as b or the same blockNumber and lower LogIndex. Otherwise, it returns false.
func lessLog(a, b IndexLog) bool {
	if a.BlockNumber < b.BlockNumber {
		return true
	}
	if a.BlockNumber == b.BlockNumber && a.LogIndex < b.LogIndex {
		return true
	}

	return false

}

type Attestation struct {
	Index     IndexLog
	RoundID   uint64
	Request   verification.Request
	Response  verification.Response
	Fee       *big.Int
	Status    Status
	Consensus bool
	Hash      common.Hash
}

// validateResponse check the MIC of the response against the MIC of the request. If the check is successful, attestation status is set to success and attestation hash is computed and set.
func (a *Attestation) validateResponse() error {

	micReq, err := a.Request.Mic()

	if err != nil {
		a.Status = ProcessError

		return errors.New("no mic in request")
	}
	micRes, err := a.Response.ComputeMic()

	if err != nil {
		a.Status = ProcessError

		return errors.New("cannot compute mic")
	}

	if micReq != micRes {
		a.Status = WrongMIC
	}

	a.Response.AddRound(a.RoundID)

	a.Hash, err = a.Response.Hash()

	if err != nil {
		a.Status = ProcessError

		return errors.New("cannot compute hash")
	}

	a.Status = Success

	return nil
}

// ParseAttestationRequestLog tries to parse AttestationRequest log as stored in the database.
func ParseAttestationRequestLog(dbLog database.Log) (*hub.HubAttestationRequest, error) {
	contractLog, err := events.ConvertDatabaseLogToChainLog(dbLog)
	if err != nil {
		return nil, err
	}
	return hubFilterer.ParseAttestationRequest(*contractLog)
}

// ParseSigningPolicyInitializedLog tries to parse SigningPolicyInitialized log as stored in the database.
func ParseSigningPolicyInitializedLog(dbLog database.Log) (*relay.RelaySigningPolicyInitialized, error) {
	contractLog, err := events.ConvertDatabaseLogToChainLog(dbLog)
	if err != nil {
		return nil, err
	}
	return relayFilterer.ParseSigningPolicyInitialized(*contractLog)
}
