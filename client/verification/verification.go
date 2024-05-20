package verification

import (
	"encoding/binary"
	"errors"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
func (r Request) AttestationType() ([32]byte, error) {

	res := [32]byte{}

	if len(r) < 96 {
		return res, errors.New("request is to short")
	}

	copy(res[:], r[0:32])

	return res, nil

}

// Source returns the source of the request (the second 32 bytes).
func (r Request) Source() ([32]byte, error) {

	res := [32]byte{}

	if len(r) < 96 {
		return res, errors.New("request is to short")
	}

	copy(res[:], r[32:64])

	return res, nil
}

// AttestationTypeAndSource returns byte encoded AttestationType and Source (the first 64 bytes).
func (r Request) AttestationTypeAndSource() ([64]byte, error) {

	res := [64]byte{}

	if len(r) < 96 {
		return [64]byte{}, errors.New("request is to short")
	}

	copy(res[:], r[0:64])

	return res, nil
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
func (r Response) ComputeMicMaybe() (common.Hash, error) {

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

	// store roundId
	d := make([]byte, 32)
	roundIdBytes := r[roundIdStartByte:roundIdEndByte]
	copy(d, roundIdBytes)

	// restore roundId at the end
	defer copy(roundIdBytes, d)

	// set roundId to zero
	zero32bytes := make([]byte, 32)
	copy(roundIdBytes, zero32bytes)

	mic := crypto.Keccak256Hash(r)

	return mic, nil

}

func (r Response) ComputeMic(args abi.Arguments) (common.Hash, error) {

	decoded, err := args.Unpack(r)

	if err != nil {
		return common.Hash{}, err
	}

	stringArgument := abi.Argument{}

	stringArgument.Name = "string"
	stringArgument.Indexed = false

	stringArgument.Type, err = abi.NewType("string", "string", []abi.ArgumentMarshaling{})

	if err != nil {
		return common.Hash{}, err
	}

	args = append(args, stringArgument)

	withSalt, err := args.Pack(decoded[0], "Flare")

	if err != nil {
		return common.Hash{}, err
	}

	mic := crypto.Keccak256Hash(withSalt)

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

	r = slices.Replace(r, roundIdStartByte, roundIdEndByte, roundIdSlot...)

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
