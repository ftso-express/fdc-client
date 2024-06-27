package attestation

import (
	"errors"
	"flare-common/database"
	"flare-common/events"
	"fmt"
	"local/fdc/client/config"
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
	RoundId   uint64
	Request   Request
	Response  Response
	Fee       *big.Int
	Status    Status
	Consensus bool
	Hash      common.Hash
	Abi       abi.Arguments
	LutLimit  uint64
}

// attestationFromDatabaseLog creates an Attestation from a request event log.
func attestationFromDatabaseLog(request database.Log) (Attestation, error) {

	requestLog, err := ParseAttestationRequestLog(request)

	if err != nil {
		return Attestation{}, fmt.Errorf("parsing attestation, parsing log: %w", err)
	}

	roundId, err := timing.RoundIdForTimestamp(request.Timestamp)

	if err != nil {
		return Attestation{}, fmt.Errorf("parsing attestation: %w", err)
	}

	index := IndexLog{request.BlockNumber, request.LogIndex}

	attestation := Attestation{
		Index:   index,
		RoundId: roundId,
		Request: requestLog.Data,
		Fee:     requestLog.Fee,
		Status:  Waiting,
	}

	return attestation, nil
}

// handleAttestation sends the attestation request to the correct verifier server and validates the response.
func (m *Manager) handleAttestation(attestation *Attestation) error {

	typeSourceConfig, err := attestation.prepareRequest(m.attestationTypeConfig)

	if err != nil {
		return fmt.Errorf("handleAttestation: %w", err)

	}

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

// prepareRequest adds response abi and lutLimit to the Attestation and returns the verifier credentials.
func (a *Attestation) prepareRequest(attestationTypesConfigs config.AttestationTypes) (config.Source, error) {

	attType, err := a.Request.AttestationType()
	if err != nil {
		a.Status = ProcessError
		return config.Source{}, err
	}

	source, err := a.Request.Source()

	if err != nil {
		a.Status = ProcessError
		return config.Source{}, err
	}

	attestationTypeConfig, ok := attestationTypesConfigs[attType]

	if !ok {
		a.Status = UnsupportedPair
		return config.Source{}, fmt.Errorf("prepare request: no configs for: %s", string(attType[:]))

	}

	a.Abi = attestationTypeConfig.ResponseArguments

	sourceConfig, ok := attestationTypeConfig.SourcesConfig[source]

	if !ok {
		a.Status = UnsupportedPair
		return config.Source{}, fmt.Errorf("prepare request: no configs for: %s, %s", string(attType[:]), string(source[:]))

	}

	a.LutLimit = sourceConfig.LutLimit
	a.Status = Processing

	return sourceConfig, nil

}

// validateResponse checks the MIC and LUT of the attestation. If both conditions pass, hash is computed and added to the attestation.
func (a *Attestation) validateResponse() error {

	// MIC
	micReq, err := a.Request.Mic()

	if err != nil {
		a.Status = ProcessError

		return errors.New("no mic in request")
	}

	micRes, err := a.Response.ComputeMic(a.Abi)

	if err != nil {
		a.Status = ProcessError

		return fmt.Errorf("cannot compute mic %w", err)
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

	roundStart := timing.ChooseStartTimestamp(a.RoundId)

	if !validLUT(lut, a.LutLimit, roundStart) {
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
