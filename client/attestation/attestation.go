package attestation

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"flare-common/database"
	"flare-common/events"
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	InvalidLUT
	Retrying
	ProcessError
	Unconfirmed
)

type IndexLog struct {
	BlockNumber uint64
	LogIndex    uint64 // consecutive number of log in block
}

// earlierLog returns true if a has lower blockNumber then b or has the same blockNumber and lower LogIndex. Otherwise, it returns false.
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
	Indexes     []IndexLog // indexLogs of all logs in the round with the Request. The earliest is in the first place.
	RoundId     uint64
	Request     Request
	Response    Response
	Fee         *big.Int // sum of fees of all logs in the round with the Request
	Status      Status
	Consensus   bool
	Hash        common.Hash
	Abi         *abi.Arguments
	LutLimit    uint64
	queueName   string
	Credentials *VerifierCredentials
}

// attestationFromDatabaseLog creates an Attestation from a request event log.
func attestationFromDatabaseLog(request database.Log) (Attestation, error) {

	requestLog, err := ParseAttestationRequestLog(request)

	if err != nil {
		return Attestation{}, fmt.Errorf("parsing log: %w", err)
	}

	roundId, err := timing.RoundIdForTimestamp(request.Timestamp)

	if err != nil {
		return Attestation{}, fmt.Errorf("parsing log, roundId: %w", err)
	}

	indexes := []IndexLog{{request.BlockNumber, request.LogIndex}}

	attestation := Attestation{
		Indexes: indexes,
		RoundId: roundId,
		Request: requestLog.Data,
		Fee:     requestLog.Fee,
		Status:  Waiting,
	}

	return attestation, nil
}

// AddToQueue adds the attestation to the correct verifier queue.
func (m *Manager) AddToQueue(ctx context.Context, attestation *Attestation) error {

	err := attestation.prepareRequest(m.attestationTypeConfig)

	if err != nil {
		return fmt.Errorf("preparing request: %w", err)
	}

	queue, ok := m.queues[attestation.queueName]

	if !ok {
		return fmt.Errorf("queue %s does not exist", attestation.queueName)
	}

	err = queue.Enqueue(ctx, attestation)

	return err
}

// handle sends the attestation request to the correct verifier server and validates the response.
func (a *Attestation) handle(ctx context.Context) error {

	confirmed, err := ResolveAttestationRequest(ctx, a)

	if err != nil {
		a.Status = ProcessError

		return fmt.Errorf("handle, resolve request: %w", err)
	}

	if !confirmed {

		a.Status = Unconfirmed

		log.Debugf("unconfirmed request: ")

		return nil
	}

	err = a.validateResponse()

	if err != nil {

		return fmt.Errorf("handle, validate response: %w", err)

	}

	return nil

}

// prepareRequest adds response ABI, LUT limit and verifierCredentials to the Attestation and returns the verifier credentials.
func (a *Attestation) prepareRequest(attestationTypesConfigs config.AttestationTypes) error {

	attType, err := a.Request.AttestationType()
	if err != nil {
		a.Status = ProcessError
		return err
	}

	source, err := a.Request.Source()

	if err != nil {
		a.Status = ProcessError
		return err
	}

	attestationTypeConfig, ok := attestationTypesConfigs[attType]

	if !ok {
		a.Status = UnsupportedPair
		return fmt.Errorf("prepare request: no configs for: %s", string(bytes.Trim(attType[:], "\x00")))

	}

	a.Abi = &attestationTypeConfig.ResponseArguments

	sourceConfig, ok := attestationTypeConfig.SourcesConfig[source]

	if !ok {
		a.Status = UnsupportedPair
		return fmt.Errorf("prepare request: no configs for: %s, %s", string(bytes.Trim(attType[:], "\x00")), string(bytes.Trim(source[:], "\x00")))

	}

	a.LutLimit = sourceConfig.LutLimit
	a.Status = Processing

	a.Credentials = new(VerifierCredentials)

	a.Credentials.Url = sourceConfig.Url
	a.Credentials.apiKey = sourceConfig.ApiKey

	a.queueName = sourceConfig.QueueName

	return nil

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

// index is used to safely retrieve Index for sorting purposes.
func (a *Attestation) index() IndexLog {
	if len(a.Indexes) > 0 {
		return a.Indexes[0]
	}
	log.Panicf("attestation without index in round %d with request %s", a.RoundId, hex.EncodeToString(a.Request)) // this should never happen
	return IndexLog{18446744073709551615, 18446744073709551615}
}

// BitVoteFromAttestations calculates BitVote for an array of attestations.
// For i-th attestation in array, i-th bit in BitVote(from the right) is 1 if and only if i-th attestation status is Success.
// Sorting of attestation must be done prior.
func BitVoteFromAttestations(attestations []*Attestation) (bitvotes.BitVote, error) {
	bitVector := big.NewInt(0)

	if len(attestations) > 65535 {
		return bitvotes.BitVote{}, errors.New("more than 65536 attestations")
	}

	for i, a := range attestations {
		if a.Status == Success {
			bitVector.SetBit(bitVector, i, 1)
		}

	}

	return bitvotes.BitVote{Length: uint16(len(attestations)), BitVector: bitVector}, nil
}
