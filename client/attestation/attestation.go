package attestation

import (
	"errors"
	"flare-common/database"
	"flare-common/events"
	"fmt"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	InvalidLUT      Status = 7
	Retrying        Status = 8
	ProcessError    Status = 9
)

type IndexLog struct {
	BlockNumber uint64
	LogIndex    uint64 // consecutive number of log in block
}

// earlierLog returns true if a has lower blockNumber as b or the same blockNumber and lower LogIndex. Otherwise, it returns false.
func earlierLog(a, b IndexLog) bool {
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
	RoundId   uint32
	Request   Request
	Response  Response
	Fee       *big.Int
	Status    Status
	Consensus bool
	Hash      common.Hash
	abi       abi.Arguments
	lutLimit  uint64
}

// handleAttestation sends the attestation request to the correct verifier server and validates the response.
func (m *Manager) handleAttestation(attestation *Attestation) error {
	attTypeAndSource, err := attestation.Request.AttestationTypeAndSource()
	if err != nil {
		attestation.Status = ProcessError
		return err
	}

	attType, err := attestation.Request.AttestationType()
	if err != nil {
		attestation.Status = ProcessError
		return err
	}

	var ok bool

	attestation.abi, ok = m.abiConfig.ResponseArguments[attType]

	if !ok {
		attestation.Status = UnsupportedPair
		return fmt.Errorf("handle attestation: no abi for: %s", string(attType[:]))

	}

	typeSourceConfig, ok := m.verifierServers[attTypeAndSource]

	if !ok {
		attestation.Status = UnsupportedPair
		return fmt.Errorf("handle attestation: no verifier for pair %s %s", string(attTypeAndSource[0:32]), string(attTypeAndSource[32:64]))

	}

	attestation.lutLimit = typeSourceConfig.LutLimit
	attestation.Status = Processing

	err = ResolveAttestationRequest(attestation, typeSourceConfig)

	if err != nil {
		attestation.Status = ProcessError

		return fmt.Errorf("handleAttestation, resolve request: %w", err)
	} else {
		err := attestation.validateResponse()

		if err != nil {

			return fmt.Errorf("handelAttestation, validate response: %w", err)

		}

		return nil
	}
}

// validateResponse checks the MIC and LUT of the attestation. If both conditions pass, hash is computed and added to the attestation.
func (a *Attestation) validateResponse() error {

	// MIC
	micReq, err := a.Request.Mic()

	if err != nil {
		a.Status = ProcessError

		return errors.New("no mic in request")
	}

	micRes, err := a.Response.ComputeMic(a.abi)

	if err != nil {
		a.Status = ProcessError

		return errors.New("cannot compute mic")
	}

	if micReq != micRes {
		a.Status = WrongMIC
		return errors.New("wrong mic")
	}

	// LUT
	lut, err := a.Response.LUT()

	if err != nil {
		a.Status = ProcessError

		return errors.New("cannot read lut")
	}

	roundStart := timing.ChooseStartTimestamp(int(a.RoundId))

	if !validLUT(lut, a.lutLimit, roundStart) {
		a.Status = InvalidLUT
		return errors.New("lut too old")
	}

	// HASH
	a.Hash, err = a.Response.Hash(a.RoundId)

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
