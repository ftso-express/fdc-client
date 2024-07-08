package bitvotes

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"local/fdc/client/shuffle"
	"math"
	"math/big"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	NumOfSamples int     = 100000    // the actual number of samples is x2 (one normal and one optimistic for each seed)
	divOpt       uint16  = 5         // totalWeight/divOpt is the weight of the optimistic samples
	valueCap     float64 = 4.0 / 5.0 // bitVote support cap in factor of totalWeight
)

type BitVote struct {
	Length    uint16 //number of attestations
	BitVector *big.Int
}

type IndexTx struct {
	BlockNumber      uint64
	TransactionIndex uint64
}

// earlierTx compares IndexTxs a,b. Returns true if a has lower BlockNumber than b or has the same BlockNumber and lower TransactionIndex.
func EarlierTx(a, b IndexTx) bool {
	if a.BlockNumber < b.BlockNumber {
		return true
	}
	if a.BlockNumber == b.BlockNumber && a.TransactionIndex < b.TransactionIndex {
		return true
	}

	return false
}

type WeightedBitVote struct {
	Index   int
	IndexTx IndexTx
	Weight  uint16
	BitVote BitVote
}

type bitVoteWithValue struct {
	index       int
	bitVote     BitVote
	valueCapped *big.Int
	value       *big.Int // support multiplied with fees
}

// fees sums the fees of the attestation requests indicated in BitVote
func (bv BitVote) fees(fees []int) (int, error) {

	if bv.BitVector.BitLen() > len(fees) {
		return 0, errors.New("a confirmed instance missing from attestations")
	}

	fee := 0

	for i := range fees {

		if bv.BitVector.Bit(i) == 1 {
			fee += fees[i]
		}
	}
	return fee, nil
}

// bitVoteForSet calculates bitwise and of the WeightedBitVote in the order defined by shuffled
// until the added weight does not exceed 50% of the total weight.
// Then it adds the weight of the rest of WeightedBitVote that support the calculated BitVote.
// Returns the BitVote that is the result of the bitwise and, and supportingWeight.
//
// We assume that the sum of the weights of the WeightedBitVotes is more than 50% of th totalWeight.
func bitVoteForSet(weightedBitVotes []*WeightedBitVote, totalWeight uint16, shuffled []uint64) (BitVote, uint16) {

	bitVote := (weightedBitVotes)[shuffled[0]].BitVote

	supportingWeight := uint16(0)

	auxBigInt := new(big.Int)

	for j := range shuffled {
		if supportingWeight < totalWeight/2 {
			bitVote = andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote)
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		} else if auxBigInt.And(bitVote.BitVector, weightedBitVotes[shuffled[j]].BitVote.BitVector).Cmp(bitVote.BitVector) == 0 {
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		}

	}

	return bitVote, supportingWeight

}

// bitVoteForSet calculates bitwise and of the WeightedBitVote in the order defined by shuffled
// until the added weight does not exceed x% of the total weight.
// Then it adds the weight of the rest of WeightedBitVote that support the calculated BitVote.
// Returns the BitVote that is the result of the bitwise and, and supportingWeight.
// It is assumed that len(weightedBitVotes) = len(shuffled) > 0, shuffled contains each integer from 0 to len(shuffled)- 1.
func bitVoteForSetOptimistic(weightedBitVotes []*WeightedBitVote, totalWeight uint16, shuffled []uint64, div uint16) (BitVote, uint16) {

	bitVote := (weightedBitVotes)[shuffled[0]].BitVote

	supportingWeight := uint16(0)

	auxBigInt := new(big.Int)

	for j := range shuffled {
		if supportingWeight < totalWeight/div {
			bitVote = andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote)
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		} else if auxBigInt.And(bitVote.BitVector, weightedBitVotes[shuffled[j]].BitVote.BitVector).Cmp(bitVote.BitVector) == 0 {
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		}

	}

	return bitVote, supportingWeight

}

// andBitwise returns the BitVote that has 1s at the places where both a and b have 1s and 0s elsewhere.
// If one BitVote is longer, the resulting BitVote has larger length with 0s on the excess places.
func andBitwise(a, b BitVote) BitVote {

	maxLen := max(a.Length, b.Length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

// Value calculates the cappedValue and Value of the BitVote, which is the product of the fees and supportingWeight.
func Values(bitVote BitVote, supportingWeight uint16, totalWeight uint16, fees []int, capPercentage float64) (*big.Int, *big.Int, error) {
	fee, err := bitVote.fees(fees)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot compute value : %s", err)
	}

	feeBig := big.NewInt(int64(fee))

	auxBigInt := new(big.Int)

	cap := uint16(math.Ceil(float64(totalWeight) * capPercentage))

	if cap < supportingWeight {

		return auxBigInt.Mul(feeBig, big.NewInt(int64(cap))), auxBigInt.Mul(feeBig, big.NewInt(int64(supportingWeight))), nil

	}
	return auxBigInt.Mul(feeBig, big.NewInt(int64(supportingWeight))), auxBigInt, nil // todo this cannot be right
}

type tempBitVoteResult struct {
	index          int
	maxValueCapped *big.Int
	maxValue       *big.Int
	bitVote        BitVote
}

type ConsensusBitVoteInput struct {
	RoundID          uint64
	WeightedBitVotes []*WeightedBitVote
	TotalWeight      uint16
	Fees             []int
}

// ConsensusBitVote calculates the ConsensusBitVote for roundId given the weightedBitVotes.
func ConsensusBitVote(input *ConsensusBitVoteInput, protocolId uint64) (BitVote, error) {
	weightVoted := sumWeightedBitVotes(input.WeightedBitVotes)

	if weightVoted <= input.TotalWeight/2 {
		percentage := (float32(weightVoted) * 100.0) / float32(input.TotalWeight)
		return BitVote{}, fmt.Errorf("only %.1f%% bitVoted", percentage)
	}

	var eg errgroup.Group
	var mu sync.Mutex
	tmpResult := tempBitVoteResult{
		maxValueCapped: big.NewInt(0),
		maxValue:       big.NewInt(0),
	}

	for i := 0; i < NumOfSamples; i++ {
		index := i
		eg.Go(func() error {
			bitVoteVals, err := calcBitVoteVals(input, index, protocolId)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			updateTmpResult(&tmpResult, bitVoteVals[0])
			if bv := bitVoteVals[1]; bv != nil {
				updateTmpResult(&tmpResult, bitVoteVals[1])
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return BitVote{}, err
	}

	return tmpResult.bitVote, nil
}

func sumWeightedBitVotes(weightedBitVotes []*WeightedBitVote) (weightVoted uint16) {
	for j := range weightedBitVotes {
		weightVoted += (weightedBitVotes[j].Weight)
	}

	return weightVoted
}

func calcBitVoteVals(input *ConsensusBitVoteInput, index int, protocolId uint64) ([2]*bitVoteWithValue, error) {
	var results [2]*bitVoteWithValue

	seed := shuffle.Seed(input.RoundID, uint64(index), protocolId)
	shuffled := shuffle.FisherYates(uint64(len(input.WeightedBitVotes)), seed)

	var eg errgroup.Group

	eg.Go(func() error {
		tempBitVote, supportingWeight := bitVoteForSet(input.WeightedBitVotes, input.TotalWeight, shuffled)
		valueCapped, value, err := Values(tempBitVote, supportingWeight, input.TotalWeight, input.Fees, valueCap)
		if err != nil {
			return err
		}

		results[0] = &bitVoteWithValue{
			index:       index,
			bitVote:     tempBitVote,
			valueCapped: valueCapped,
			value:       value,
		}

		return nil
	})

	eg.Go(func() error {
		tempBitVoteOpt, supportingWeightOpt := bitVoteForSetOptimistic(
			input.WeightedBitVotes, input.TotalWeight, shuffled, divOpt,
		)

		if supportingWeightOpt > input.TotalWeight/2 {
			valueOptimisticCapped, valueOptimistic, err := Values(
				tempBitVoteOpt, supportingWeightOpt, input.TotalWeight, input.Fees, valueCap,
			)
			if err != nil {
				return err
			}

			results[1] = &bitVoteWithValue{
				index:       index,
				bitVote:     tempBitVoteOpt,
				valueCapped: valueOptimisticCapped,
				value:       valueOptimistic,
			}
		}

		return nil
	})

	err := eg.Wait()

	return results, err
}

func updateTmpResult(tmpResult *tempBitVoteResult, result *bitVoteWithValue) {
	if result.valueCapped.Cmp(tmpResult.maxValueCapped) == 1 {
		tmpResult.bitVote = result.bitVote
		tmpResult.index = result.index
		tmpResult.maxValueCapped = result.valueCapped
		tmpResult.maxValue = result.value
	} else if result.valueCapped.Cmp(tmpResult.maxValueCapped) == 0 && result.value.Cmp(tmpResult.maxValue) == 1 {
		tmpResult.bitVote = result.bitVote
		tmpResult.index = result.index
		tmpResult.maxValue = result.value
	} else if result.valueCapped.Cmp(tmpResult.maxValueCapped) == 0 && result.value.Cmp(tmpResult.maxValue) == 0 && tmpResult.index > result.index {
		tmpResult.bitVote = result.bitVote
		tmpResult.index = result.index
	}
}

// EncodeBitVoteHex encodes BitVote with roundCheck to be published on chain
func (b BitVote) EncodeBitVoteHex(roundId uint64) string {

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

// DecodeBitVoteBytes decodes bytes encoded BitVote and returns roundCheck
func DecodeBitVoteBytes(bitVoteByte []byte) (BitVote, uint8, error) {

	if len(bitVoteByte) < 3 {
		return BitVote{}, 0, errors.New("bitVote too short")
	}

	roundCheck := bitVoteByte[0]
	lengthBytes := bitVoteByte[1:3]
	bitVector := bitVoteByte[3:]

	length := binary.BigEndian.Uint16(lengthBytes)

	bigBitVector := new(big.Int).SetBytes(bitVector)

	// TODO: decide whether leading zeros are legal

	if bigBitVector.BitLen() > int(length) {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	return BitVote{length, bigBitVector}, roundCheck, nil

}
