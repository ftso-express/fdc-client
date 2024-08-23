package attestation

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"flare-common/database"
	"flare-common/events"
	"flare-common/logger"
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	"local/fdc/contracts/fdchub"
	"math"
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

var log = logger.GetLogger()

// fdcFilterer is only used for Attestation Requests logs parsing. Set in init().
var fdcFilterer *fdchub.FdcHubFilterer

// init sets the fdcFilterer
func init() {

	var err error

	fdcFilterer, err = fdchub.NewFdcHubFilterer(common.Address{}, nil)

	if err != nil {
		log.Panic("cannot get fdc contract:", err)
	}

}

type IndexLog struct {
	BlockNumber uint64
	LogIndex    uint64 // consecutive number of log in block
}

type Attestation struct {
	Indexes           []IndexLog // indexLogs of all logs in the round with the Request. The earliest is in the first place.
	RoundID           uint32
	Request           Request
	Response          Response
	Fee               *big.Int // sum of fees of all logs in the round with the Request
	Status            Status
	Consensus         bool
	Hash              common.Hash
	ResponseABI       *abi.Arguments
	ResponseABIString *string
	LUTLimit          uint64
	QueueName         string
	Credentials       *VerifierCredentials
}

// earlierLog returns true if a has lower blockNumber then b or has the same blockNumber and lower LogIndex. Otherwise, it returns false.
func EarlierLog(a, b IndexLog) bool {
	if a.BlockNumber < b.BlockNumber {
		return true
	}
	if a.BlockNumber == b.BlockNumber && a.LogIndex < b.LogIndex {
		return true
	}

	return false
}

// attestationFromDatabaseLog creates an Attestation from a request event log.
func AttestationFromDatabaseLog(request database.Log) (Attestation, error) {
	requestLog, err := ParseAttestationRequestLog(request)
	if err != nil {
		return Attestation{}, fmt.Errorf("parsing log: %s", err)
	}

	roundD, err := timing.RoundIDForTimestamp(request.Timestamp)
	if err != nil {
		return Attestation{}, fmt.Errorf("parsing log, roundID: %s", err)
	}

	indexes := []IndexLog{{request.BlockNumber, request.LogIndex}}

	attestation := Attestation{
		Indexes: indexes,
		RoundID: roundD,
		Request: requestLog.Data,
		Fee:     requestLog.Fee,
		Status:  Waiting,
	}

	return attestation, nil
}

// Handle sends the attestation request to the correct verifier server and validates the response.
// The response is saved in the struct.
func (a *Attestation) Handle(ctx context.Context) error {
	responseBytes, confirmed, err := ResolveAttestationRequest(ctx, a)
	if err != nil {
		a.Status = ProcessError
		return fmt.Errorf("handle, resolve request: %s", err)
	}
	if !confirmed {
		a.Status = Unconfirmed
		log.Debugf("unconfirmed request: ")
		return nil
	}

	a.Response = responseBytes
	err = a.validateResponse()
	if err != nil {
		return fmt.Errorf("handle, validate response: %s", err)
	}

	return nil
}

// prepareRequest adds response ABI, LUT limit and verifierCredentials to the Attestation.
func (a *Attestation) PrepareRequest(attestationTypesConfigs config.AttestationTypes) error {
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
	a.ResponseABI = &attestationTypeConfig.ResponseArguments

	a.ResponseABIString = &attestationTypeConfig.ResponseABIString

	sourceConfig, ok := attestationTypeConfig.SourcesConfig[source]
	if !ok {
		a.Status = UnsupportedPair
		return fmt.Errorf("prepare request: no configs for: %s, %s", string(bytes.Trim(attType[:], "\x00")), string(bytes.Trim(source[:], "\x00")))
	}

	a.LUTLimit = sourceConfig.LUTLimit
	a.Status = Processing
	a.Credentials = new(VerifierCredentials)
	a.Credentials.URL = sourceConfig.URL
	a.Credentials.apiKey = sourceConfig.APIKey
	a.QueueName = sourceConfig.QueueName

	return nil

}

// validateResponse checks the MIC and LUT of the attestation. If both conditions pass, hash is computed and added to the attestation.
func (a *Attestation) validateResponse() error {
	// MIC
	micReq, err := a.Request.Mic()
	if err != nil {
		a.Status = ProcessError
		return fmt.Errorf("reading mic in request: %s, %s ", hex.EncodeToString(a.Request), err)
	}

	micRes, err := a.Response.ComputeMic(a.ResponseABI)
	if err != nil {
		a.Status = ProcessError
		return fmt.Errorf("cannot compute mic for request: %s, %s", hex.EncodeToString(a.Request), err)
	}

	if micReq != micRes {
		a.Status = WrongMIC
		return fmt.Errorf("wrong mic in request: %s", hex.EncodeToString(a.Request))
	}

	// LUT
	lut, err := a.Response.LUT()
	if err != nil {
		a.Status = ProcessError
		return fmt.Errorf("cannot read lut from request: %s, %s", hex.EncodeToString(a.Request), err)
	}

	roundStart := timing.ChooseStartTimestamp(a.RoundID)
	if !validLUT(lut, a.LUTLimit, roundStart) {
		a.Status = InvalidLUT
		return fmt.Errorf("lot too old in request: %s", hex.EncodeToString(a.Request))
	}

	// HASH
	a.Hash, err = a.Response.Hash(a.RoundID)
	if err != nil {
		a.Status = ProcessError
		return fmt.Errorf("cannot compute hash for request: %s", hex.EncodeToString(a.Request))
	}

	a.Status = Success

	return nil
}

// ParseAttestationRequestLog tries to parse AttestationRequest log as stored in the database.
func ParseAttestationRequestLog(dbLog database.Log) (*fdchub.FdcHubAttestationRequest, error) {
	contractLog, err := events.ConvertDatabaseLogToChainLog(dbLog)
	if err != nil {
		return nil, err
	}
	return fdcFilterer.ParseAttestationRequest(*contractLog)
}

// index is used to safely retrieve Index for sorting purposes.
func (a *Attestation) Index() IndexLog {
	if len(a.Indexes) > 0 {
		return a.Indexes[0]
	}
	log.Panicf("attestation without index in round %d with request %s", a.RoundID, hex.EncodeToString(a.Request)) // this should never happen

	return IndexLog{math.MaxUint64, math.MaxUint64}
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
