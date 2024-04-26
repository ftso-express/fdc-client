package attestation

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"local/fdc/client/shuffle"
	"math/big"
)

const (
	NumOfSamples int = 100 // read from config/toml
)

type BitVote struct {
	Length    uint16 //number attestations
	BitVector *big.Int
}

type WeightedBitVote struct {
	weight  uint64
	bitVote BitVote
}

type bitVoteWithValue struct {
	index   int64
	bitVote BitVote
	value   *big.Int
	err     error
}

// BitVoteFromAttestations calculates BitVote for an array of attestations.
// For i-th attestation in array, i-th bit in BitVote(from the right) is 1 if and only if i-th attestation status is Success.
func BitVoteFromAttestations(attestations []*Attestation) (BitVote, error) {
	bitVector := big.NewInt(0)

	if len(attestations) > 65535 {
		return BitVote{}, errors.New("more than 65536 attestations")
	}

	for i, a := range attestations {
		if a.Status == Success {
			bitVector.SetBit(bitVector, i, 1)
		}

	}
	return BitVote{uint16(len(attestations)), bitVector}, nil
}

// fees sums the fees of the attestation requests indicated in BitVote
func (bv BitVote) fees(attestations []*Attestation) (*big.Int, error) {

	if bv.BitVector.BitLen() > len(attestations) {
		return nil, errors.New("attestations do not hold a confirmed instance")
	}

	fees := big.NewInt(0)

	for i, a := range attestations {

		if bv.BitVector.Bit(i) == 1 {
			fees.Add(fees, a.Fee)
		}
	}
	return fees, nil
}

// BitVoteForSet calculates bitwise and of the WeightedBitVote in the order defined by shuffled
// until the added weight does not exceed 50% of the total weight.
// Then it adds the weight of the rest of WeightedBitVote that support the calculated BitVote.
// Returns the BitVote that is the result of the bitwise and, and supportingWeight.
func bitVoteForSet(weightedBitVotes []WeightedBitVote, totalWeight uint64, shuffled []uint64) (BitVote, uint64) {

	bitVote := (weightedBitVotes)[shuffled[0]].bitVote

	halfWeight := (totalWeight + 1) / 2

	supportingWeight := uint64(0)

	for _, v := range shuffled {
		if supportingWeight < halfWeight {
			bitVote = andBitwise(bitVote, weightedBitVotes[v].bitVote)
			supportingWeight += weightedBitVotes[v].weight
		} else if andBitwise(bitVote, weightedBitVotes[v].bitVote).BitVector == weightedBitVotes[v].bitVote.BitVector {
			supportingWeight += weightedBitVotes[v].weight
		}

	}

	return bitVote, supportingWeight

}

// andBitwise returns the BitVote that has 1s at the places where both a and b have 1s and 0s elsewhere.
// If one BitVote is longer the resulting BitVote has larger length with 0s on the excess places.
func andBitwise(a, b BitVote) BitVote {

	maxLen := max(a.Length, b.Length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

// Value calculates the value of the BitVote, which is the product of the fees and supportingWeight.
func value(bitVote BitVote, supportingWeight uint64, attestations []*Attestation) (*big.Int, error) {
	fees, err := bitVote.fees(attestations)

	if err != nil {
		return nil, errors.New("cannot compute fees")
	}

	return fees.Mul(fees, big.NewInt(int64(supportingWeight))), nil
}

// ConsensusBitVote calculates the ConsensusBitVote for roundId given the weightedBitVotes.
func ConsensusBitVote(roundId uint64, weightedBitVotes []WeightedBitVote, totalWeight uint64, attestations []*Attestation) (BitVote, error) {

	var bitVote BitVote
	maxValue := big.NewInt(0)
	noOfVoters := len(weightedBitVotes)
	index := int64(0)

	ch := make(chan bitVoteWithValue)

	for i := 0; i < NumOfSamples; i++ {
		go func(j int64) {
			seed := shuffle.Seed(int64(roundId), j)
			shuffled := shuffle.FisherYates(uint64(noOfVoters), seed)
			tempBitVote, supportingWeight := bitVoteForSet(weightedBitVotes, totalWeight, shuffled)
			value, err := value(tempBitVote, supportingWeight, attestations)

			ch <- bitVoteWithValue{j, tempBitVote, value, err}
		}(int64(i))
	}

	for i := 0; i < NumOfSamples; i++ {
		result := <-ch

		if result.err != nil {
			return BitVote{}, errors.New("missing attestations. cannot compute consensus bitvote")
		}

		if result.value.Cmp(maxValue) == 1 {
			bitVote = result.bitVote
			index = result.index
			maxValue = result.value
		} else if result.value.Cmp(maxValue) == 0 && index > result.index {
			bitVote = result.bitVote
			index = result.index
		}
	}

	return bitVote, nil
}

// SetBitVoteStatus sets the Consensus status of attestations to true for the attestations chosen by BitVote
func SetBitVoteStatus(attestations []*Attestation, bitVote BitVote) error {

	if bitVote.BitVector.BitLen() > len(attestations) {
		return errors.New("chosen attestation does not exist")
	}

	for i := range attestations {
		attestations[i].Consensus = bitVote.BitVector.Bit(int(attestations[i].Index)) == 1
	}

	return nil

}

// EncodeHEx encodes BitVote with roundCheck to be published on chain
func (b BitVote) EncodeHex(roundId uint64) string {

	var encoding []byte
	roundCheck := byte(roundId % 256)

	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, b.Length)

	encoding = append(encoding, roundCheck)
	encoding = append(encoding, length...)
	encoding = append(encoding, b.BitVector.Bytes()...)

	str := hex.EncodeToString(encoding)

	return str

}

// DecodeHex decodes hex encoded BitVote and returns roundCheck
func DecodeHex(bitVoteHex string) (BitVote, uint8, error) {

	roundCheckStr := bitVoteHex[:2]
	lengthStr := bitVoteHex[2:6]
	bitVectorStr := bitVoteHex[6:]

	roundCheckBytes, err := hex.DecodeString(roundCheckStr)
	if err != nil || len(roundCheckBytes) != 1 {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	roundCheck := uint8(roundCheckBytes[0])

	lengthBytes, err := hex.DecodeString(lengthStr)
	if err != nil || len(lengthBytes) != 2 {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	length := binary.BigEndian.Uint16(lengthBytes)

	bitVector := big.NewInt(0)
	_, success := bitVector.SetString(bitVectorStr, 16)

	if !success {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	if bitVector.BitLen() > int(length) {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	return BitVote{length, bitVector}, roundCheck, nil

}
