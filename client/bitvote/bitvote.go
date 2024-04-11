package bitvote

import (
	"errors"
	"local/fdc/client/attestation"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type BitVote struct {
	length    byte //number of relevant bits
	BitVector *big.Int
}

const (
	NumOfSamples int = 100
)

func bitvote(attestations []attestation.Attestation) (BitVote, error) {
	bitVector := big.NewInt(0)

	if len(attestations) > 255 {
		return BitVote{}, errors.New("more than 255 attestations")
	}

	for i, a := range attestations {
		if a.Status == attestation.Success {
			bitVector.SetBit(bitVector, i, 1)
		}

	}
	return BitVote{byte(len(attestations)), bitVector}, nil
}

func (bv BitVote) fees(attestations []attestation.Attestation) *big.Int {

	fees := big.NewInt(0)

	for i, a := range attestations {

		if bv.BitVector.Bit(i) == 1 {
			fees.Add(fees, a.Fee)
		}
	}
	return fees
}

func BitVoteForSet(bitVotes []BitVote, weights []uint64, totalWeight uint64, shuffled []uint64) (BitVote, uint64) {

	bitVote := (bitVotes)[shuffled[0]]

	halfWeight := (totalWeight + 1) / 2

	supportingWeight := uint64(0)

	for _, v := range shuffled {
		if supportingWeight < halfWeight {
			bitVote = ANDbitwise(bitVote, bitVotes[v])
			supportingWeight += weights[v]
		} else if ANDbitwise(bitVote, bitVotes[v]).BitVector == bitVotes[v].BitVector {
			supportingWeight += weights[v]
		}

	}

	return bitVote, supportingWeight

}

func Value(bitVote BitVote, supportingWeight uint64, attestations []attestation.Attestation) *big.Int {
	fees := bitVote.fees(attestations)

	return fees.Mul(fees, big.NewInt(int64(supportingWeight)))
}

func ANDbitwise(a, b BitVote) BitVote {

	maxLen := max(a.length, b.length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

type bitVoteWithValue struct {
	index   int
	bitVote BitVote
	value   *big.Int
}

func ConsensusBitVote(roundId uint64, bitVotes []BitVote, weights []uint64, totalWeight uint64, attestations []attestation.Attestation) BitVote {

	var bitVote BitVote
	maxValue := big.NewInt(0)
	noOfVoters := len(weights)
	index := 0

	ch := make(chan bitVoteWithValue)

	go func() {
		for i := 0; i < NumOfSamples; i++ {
			seed := crypto.Keccak256([]byte{byte(roundId)}, []byte{byte(i)})
			shuffled := FisherYates(uint64(noOfVoters), seed)
			tempBitVote, supportingWeight := BitVoteForSet(bitVotes, weights, totalWeight, shuffled)
			value := Value(tempBitVote, supportingWeight, attestations)

			ch <- bitVoteWithValue{i, tempBitVote, value}
		}
	}()

	for i := 0; i < NumOfSamples; i++ {
		result := <-ch

		if result.value.Cmp(maxValue) == 1 {
			bitVote = result.bitVote
			index = result.index
			maxValue = result.value
		} else if result.value.Cmp(maxValue) == 0 && index > result.index {
			bitVote = result.bitVote
			index = result.index
		}
	}

	return bitVote
}
