package bitvote

import (
	"errors"
	"local/fdc/client/attestation"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type BitVote struct {
	Length    uint16 //number of relevant bits
	BitVector *big.Int
}

type WeightedBitVote struct {
	weight  uint64
	bitVote BitVote
}

const (
	NumOfSamples int = 100
)

func ForRound(attestations []attestation.Attestation) (BitVote, error) {
	bitVector := big.NewInt(0)

	if len(attestations) > 65535 {
		return BitVote{}, errors.New("more than 65536 attestations")
	}

	for i, a := range attestations {
		if a.Status == attestation.Success {
			bitVector.SetBit(bitVector, i, 1)
		}

	}
	return BitVote{uint16(len(attestations)), bitVector}, nil
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

func BitVoteForSet(weightedBitVotes []WeightedBitVote, totalWeight uint64, shuffled []uint64) (BitVote, uint64) {

	bitVote := (weightedBitVotes)[shuffled[0]].bitVote

	halfWeight := (totalWeight + 1) / 2

	supportingWeight := uint64(0)

	for _, v := range shuffled {
		if supportingWeight < halfWeight {
			bitVote = ANDbitwise(bitVote, weightedBitVotes[v].bitVote)
			supportingWeight += weightedBitVotes[v].weight
		} else if ANDbitwise(bitVote, weightedBitVotes[v].bitVote).BitVector == weightedBitVotes[v].bitVote.BitVector {
			supportingWeight += weightedBitVotes[v].weight
		}

	}

	return bitVote, supportingWeight

}

func Value(bitVote BitVote, supportingWeight uint64, attestations []attestation.Attestation) *big.Int {
	fees := bitVote.fees(attestations)

	return fees.Mul(fees, big.NewInt(int64(supportingWeight)))
}

func ANDbitwise(a, b BitVote) BitVote {

	maxLen := max(a.Length, b.Length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

type bitVoteWithValue struct {
	index   int
	bitVote BitVote
	value   *big.Int
}

func ConsensusBitVote(roundId uint64, weightedBitVotes []WeightedBitVote, totalWeight uint64, attestations []attestation.Attestation) BitVote {

	var bitVote BitVote
	maxValue := big.NewInt(0)
	noOfVoters := len(weightedBitVotes)
	index := 0

	ch := make(chan bitVoteWithValue)

	for i := 0; i < NumOfSamples; i++ {
		go func(j int) {
			seed := crypto.Keccak256([]byte{byte(roundId)}, []byte{byte(j)})
			shuffled := FisherYates(uint64(noOfVoters), seed)
			tempBitVote, supportingWeight := BitVoteForSet(weightedBitVotes, totalWeight, shuffled)
			value := Value(tempBitVote, supportingWeight, attestations)

			ch <- bitVoteWithValue{j, tempBitVote, value}
		}(i)
	}

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

func SetBitVoteStatus(attestations []attestation.Attestation, bitVote BitVote) {

	for i := range attestations {
		attestations[i].Consensus = bitVote.BitVector.Bit(int(attestations[i].Index)) == 1
	}

}
