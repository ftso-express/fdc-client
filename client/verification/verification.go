package verification

import (
	"encoding/binary"
	"errors"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Request []byte

type Response []byte

// AttestationType returns the attestation type of the request (the first 32 bytes).
func (r Request) AttestationType() ([]byte, error) {

	if len(r) < 96 {
		return []byte{}, errors.New("request is to short")
	}

	return r[:32], nil
}

// Source returns the source of the request (the second 32 bytes).
func (r Request) Source() ([]byte, error) {

	if len(r) < 96 {
		return []byte{}, errors.New("request is to short")
	}

	return r[32:64], nil
}

// AttestationTypeAndSource returns byte encoded AttestationType and Source (the first 64 bytes).
func (r Request) AttestationTypeAndSource() ([]byte, error) {

	if len(r) < 96 {
		return []byte{}, errors.New("request is to short")
	}

	return r[:64], nil
}

// Mic returns Message Integrity code of the request (the third 32 bytes).
func (r Request) Mic() (common.Hash, error) {

	if len(r) < 96 {
		return common.Hash{}, errors.New("request is to short")
	}

	mic := common.Hash{}

	mic.SetBytes(r[64:96])

	return mic, nil
}

// ComputeMic computes Mic from the response.
// Mic is the hash of the response with roundID set to 0.
func (r Response) ComputeMic() (common.Hash, error) {
	if len(r) < 128 {
		return common.Hash{}, errors.New("response is to short")
	}

	d := make([]byte, 32)

	slices.Replace(d, 0, 32, r[64:96]...)

	zero32bytes := make([]byte, 32)

	slices.Replace(r, 64, 96, zero32bytes...)

	mic := crypto.Keccak256Hash(r)

	defer slices.Replace(r, 64, 96, d...)

	return mic, nil

}

// AddRound sets the roundId in the response (third 32 bytes).
func (r Response) AddRound(roundId uint64) ([]byte, error) {
	if len(r) < 128 {
		return []byte{}, errors.New("response is to short")
	}

	end := make([]byte, 32)
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf, roundId)

	slices.Replace(end, 24, 32, buf...)

	slices.Replace(r, 64, 96, end...)

	return r, nil

}

// Hash computes hash of the response.
func (r Response) Hash() (common.Hash, error) {
	if len(r) < 128 {
		return common.Hash{}, errors.New("response is to short")
	}

	hash := crypto.Keccak256Hash(r)

	return hash, nil
}
