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

	// roundId is encoded in the third 32bytes slot
	roundIdStartByte := 64
	roundIdEndByte := 96
	commonFieldsLength := 128

	// if Response is encoded dynamic struct the first 32 bytes are bytes32(32)
	if !static {
		roundIdStartByte += 32
		roundIdEndByte += 32
		commonFieldsLength += 32
	}

	if len(r) < commonFieldsLength {
		return common.Hash{}, errors.New("response is to short")
	}

	d := make([]byte, 32)

	// store roundId
	slices.Replace(d, 0, 32, r[64:96]...)

	// restore roundId at the end
	defer slices.Replace(r, roundIdStartByte, roundIdEndByte, d...)

	zero32bytes := make([]byte, 32)

	// set roundId to zero
	slices.Replace(r, roundIdStartByte, roundIdEndByte, zero32bytes...)

	mic := crypto.Keccak256Hash(r)

	return mic, nil

}

// AddRound sets the roundId in the response (third 32 bytes).
func (r Response) AddRound(roundId uint64) (Response, error) {

	static, err := isStaticType(r)

	if err != nil {
		return Response{}, err
	}

	// roundId is encoded in the third slot
	roundIdStartByte := 64
	roundIdEndByte := 96
	commonFieldsLength := 128

	// if Response is encoded dynamic struct the first 32 bytes are bytes32(32)
	if !static {
		roundIdStartByte += 32
		roundIdEndByte += 32
		commonFieldsLength += 32
	}

	if len(r) < commonFieldsLength {
		return Response{}, errors.New("response is to short")
	}

	roundIdEncoded := make([]byte, 0)

	binary.BigEndian.AppendUint64(roundIdEncoded, roundId)

	roundIdSlot := make([]byte, 32-len(roundIdEncoded))

	roundIdSlot = append(roundIdSlot, roundIdEncoded...)

	slices.Replace(r, roundIdStartByte, roundIdEndByte, roundIdSlot...)

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
