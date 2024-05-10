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

// isStaticType checks whether bytes that represent abi.encoded request or response that encodes an instance of static type.
// abi.encode(X) = enc((X)) of X of type T is encoding of tuple (X) of type (T). By specification, enc((X)) = head(X)tail(X).
// If T is static, head(X) = enc(X) and tail(X) is empty. If T is dynamic, head(X) = bytes32(len(head(X))) = bytes32(32) and tail = enc(X).
// See https://docs.soliditylang.org/en/latest/abi-spec.html for detailed specification.
func isStaticType(bytes []byte) (bool, error) {

	if len(bytes) < 96 {
		return false, errors.New("bytes are to short")
	}

	first32Bytes := [32]byte(bytes[:32])

	d := [32]byte{}

	d[31] = byte(32)

	return d != first32Bytes, nil

}

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

	static, err := isStaticType(r)

	if err != nil {
		return common.Hash{}, err
	}

	if static {

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

	} else {

		if len(r) < 160 {
			return common.Hash{}, errors.New("response is to short")
		}

		d := make([]byte, 32)

		slices.Replace(d, 0, 32, r[96:128]...)

		zero32bytes := make([]byte, 32)

		slices.Replace(r, 96, 128, zero32bytes...)

		mic := crypto.Keccak256Hash(r)

		defer slices.Replace(r, 96, 128, d...)

		return mic, nil

	}

}

// AddRound sets the roundId in the response (third 32 bytes).
func (r Response) AddRound(roundId uint64) (Response, error) {

	static, err := isStaticType(r)

	if err != nil {
		return Response{}, err
	}

	if static {
		if len(r) < 128 {
			return []byte{}, errors.New("response is to short")
		}

		end := make([]byte, 32)
		buf := make([]byte, 8)

		binary.BigEndian.PutUint64(buf, roundId)

		slices.Replace(end, 24, 32, buf...)

		slices.Replace(r, 64, 96, end...)

		return r, nil
	} else if len(r) < 160 {
		return []byte{}, errors.New("response is to short")
	}

	end := make([]byte, 32)
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf, roundId)

	slices.Replace(end, 24, 32, buf...)

	slices.Replace(r, 96, 128, end...)

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
