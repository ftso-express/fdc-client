package bitvotes

import (
	"math/big"
	"strconv"
)

// AggregateBitVotes joins WeightedBitVotes with the same BitVote and sums their weight.
func AggregateBitVotes(bitVotes []*WeightedBitVote) ([]*WeightedBitVote, map[int][]int) {
	aggregatorsArray := []*WeightedBitVote{}
	aggregationMap := make(map[int][]int)
	wordToIndex := make(map[string]int)

	counter := 0
	for j := range bitVotes {
		bitVoteString := bitVotes[j].BitVote.BitVector.String()
		index, ok := wordToIndex[bitVoteString]

		if ok {
			aggregatorsArray[index].Weight += bitVotes[j].Weight
			aggregationMap[index] = append(aggregationMap[index], j)
		} else {
			newAggregator := &WeightedBitVote{BitVote: bitVotes[j].BitVote, Index: bitVotes[j].Index, Weight: bitVotes[j].Weight}
			wordToIndex[bitVoteString] = counter
			aggregationMap[counter] = []int{j}
			aggregatorsArray = append(aggregatorsArray, newAggregator)
			counter++
		}
	}

	return aggregatorsArray, aggregationMap
}

// AggregateAttestations joins attestations among WeightedBitVotes with the same BitVote and sums their fees.
func AggregateAttestations(bitVotes []*WeightedBitVote, fees []int) ([]*WeightedBitVote, []int, map[int][]int) {
	if len(fees) == 0 {
		return bitVotes, fees, nil
	}

	wordToIndex := make(map[string]int)
	aggregationMap := make(map[int][]int)

	newFees := make([]int, 0)
	counter := 0
	for i := 0; i < len(fees); i++ {
		word := ""
		for _, e := range bitVotes {
			word += strconv.Itoa(int(e.BitVote.BitVector.Bit(i)))
		}

		if index, ok := wordToIndex[word]; ok {
			aggregationMap[index] = append(aggregationMap[index], i)
			newFees[index] += fees[i]
		} else {
			wordToIndex[word] = counter
			aggregationMap[counter] = []int{i}
			newFees = append(newFees, fees[i])
			counter++
		}
	}

	newLength := uint16(counter)

	newBitVotes := make([]*WeightedBitVote, len(bitVotes))
	for i, e := range bitVotes {
		newBitVotes[i] = &WeightedBitVote{Index: e.Index, IndexTx: e.IndexTx, Weight: e.Weight, BitVote: BitVote{Length: newLength, BitVector: big.NewInt(0)}}
		for j := 0; j < int(newLength); j++ {
			newBit := bitVotes[i].BitVote.BitVector.Bit(aggregationMap[j][0])
			if newBit == 1 {
				newBitVotes[i].BitVote.BitVector.Add(newBitVotes[i].BitVote.BitVector, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(j)), nil))
			}
		}
	}

	return newBitVotes, newFees, aggregationMap
}

// FilterBitVotes filters out those bit votes whose bit vector is all ones or all zeros
func FilterBitVotes(bitVotes []*WeightedBitVote) ([]*WeightedBitVote, []int, uint16, []int, uint16) {
	removedOnes := make([]int, 0)
	removedOnesWeight := uint16(0)
	removedZeros := make([]int, 0)
	removedZerosWeight := uint16(0)

	newBitVotes := make([]*WeightedBitVote, 0)

	ones := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bitVotes[0].BitVote.Length)), nil)
	ones.Sub(ones, big.NewInt(1))
	zeros := big.NewInt(0)

	for j, bitVote := range bitVotes {
		if bitVote.BitVote.BitVector.Cmp(ones) == 0 {
			removedOnes = append(removedOnes, j)
			removedOnesWeight += bitVote.Weight
		} else if bitVote.BitVote.BitVector.Cmp(zeros) == 0 {
			removedZeros = append(removedZeros, j)
			removedZerosWeight += bitVote.Weight
		} else {
			newBitVotes = append(newBitVotes, bitVote)
		}
	}

	return newBitVotes, removedOnes, removedOnesWeight, removedZeros, removedZerosWeight
}

// FilterAttestations filters out those attestations that are confirmed by all or by less than half
func FilterAttestations(bitVotes []*WeightedBitVote, fees []int, totalWeight uint16) ([]*WeightedBitVote, []int, int, []int, []int) {
	removedOnes := make([]int, 0)
	removedLowWeight := make([]int, 0)
	remains := make([]int, 0)
	removedOnesFees := 0

	for i := 0; i < int(bitVotes[0].BitVote.Length); i++ {
		checkOnes := true
		weight := uint16(0)
		for _, bitVote := range bitVotes {
			if bitVote.BitVote.BitVector.Bit(i) == 0 {
				checkOnes = false
			} else {
				weight += bitVote.Weight
			}
		}
		if checkOnes {
			removedOnes = append(removedOnes, i)
			removedOnesFees += fees[i]
		} else if weight < totalWeight/2 {
			removedLowWeight = append(removedLowWeight, i)
		} else {
			remains = append(remains, i)
		}
	}

	newLength := uint16(len(remains))
	newBitVotes := make([]*WeightedBitVote, len(bitVotes))
	for i, e := range bitVotes {
		newBitVotes[i] = &WeightedBitVote{Index: e.Index, IndexTx: e.IndexTx, Weight: e.Weight, BitVote: BitVote{Length: newLength, BitVector: big.NewInt(0)}}
		for j := 0; j < int(newLength); j++ {
			newBit := bitVotes[i].BitVote.BitVector.Bit(remains[j])
			if newBit == 1 {
				newBitVotes[i].BitVote.BitVector.Add(newBitVotes[i].BitVote.BitVector, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(j)), nil))
			}
		}
	}

	newFees := make([]int, newLength)
	for j := 0; j < int(newLength); j++ {
		newFees[j] = fees[remains[j]]
	}

	return newBitVotes, newFees, removedOnesFees, removedOnes, removedLowWeight
}

type PreProcessInfo struct {
	RemovedBitVotesOnes    []int
	RemovedBitVotesZeros   []int
	BitVotesAggregationMap map[int][]int

	RemovedAttestationsOnes      []int
	RemovedAttestationsLowWeight []int
	AttestationsAggregationMap   map[int][]int

	RemovedZerosWeight uint16
	RemovedOnesWeight  uint16

	RemovedOnesFees int

	NumAttestationsBeforeAggregation int
	NumBitVotesBeforeAggregation     int
}

func PreProcess(bitVotes []*WeightedBitVote, fees []int) ([]*WeightedBitVote, []int, *PreProcessInfo) {
	totalWeight := uint16(0)
	for _, bitVote := range bitVotes {
		totalWeight += bitVote.Weight
	}

	filteredBitVotes1, filteredFees, removedOnesFee, removedAttestationsOnes, removedLowWeight := FilterAttestations(bitVotes, fees, totalWeight)
	filteredBitVotes2, removedBitVotesOnes, removedOnesWeight, removedBitVotesZeros, removedZerosWeight := FilterBitVotes(filteredBitVotes1)

	aggregatedBitVotes1, aggregatedFees, attestationsAggregationMap := AggregateAttestations(filteredBitVotes2, filteredFees)
	aggregateBitVotes2, bitVotesAggregationMap := AggregateBitVotes(aggregatedBitVotes1)

	info := &PreProcessInfo{
		RemovedBitVotesOnes: removedBitVotesOnes, RemovedBitVotesZeros: removedBitVotesZeros,
		BitVotesAggregationMap: bitVotesAggregationMap, RemovedAttestationsOnes: removedAttestationsOnes,
		RemovedAttestationsLowWeight: removedLowWeight, AttestationsAggregationMap: attestationsAggregationMap,
		RemovedZerosWeight: removedZerosWeight, RemovedOnesWeight: removedOnesWeight,
		RemovedOnesFees: removedOnesFee, NumAttestationsBeforeAggregation: len(filteredFees), NumBitVotesBeforeAggregation: len(filteredBitVotes2),
	}

	return aggregateBitVotes2, aggregatedFees, info
}

func ExpandSolution(compressedSolution *ConsensusSolution, preProcessInfo *PreProcessInfo) *ConsensusSolution {
	unAggregatedSolution := &ConsensusSolution{
		Participants: make([]bool, preProcessInfo.NumBitVotesBeforeAggregation),
		Solution:     make([]bool, preProcessInfo.NumAttestationsBeforeAggregation),
		Value:        compressedSolution.Value,
		Optimal:      compressedSolution.Optimal,
	}

	for key, val := range preProcessInfo.BitVotesAggregationMap {
		for _, pos := range val {
			unAggregatedSolution.Participants[pos] = compressedSolution.Participants[key]
		}
	}
	for key, val := range preProcessInfo.AttestationsAggregationMap {
		for _, pos := range val {
			unAggregatedSolution.Solution[pos] = compressedSolution.Solution[key]
		}
	}

	newSolution := &ConsensusSolution{
		Participants: make([]bool, preProcessInfo.NumBitVotesBeforeAggregation+len(preProcessInfo.RemovedBitVotesOnes)+len(preProcessInfo.RemovedBitVotesZeros)),
		Solution:     make([]bool, preProcessInfo.NumAttestationsBeforeAggregation+len(preProcessInfo.RemovedAttestationsOnes)+len(preProcessInfo.RemovedAttestationsLowWeight)),
		Value:        compressedSolution.Value,
		Optimal:      compressedSolution.Optimal,
	}

	participantCounter := 0
	for i := range newSolution.Participants {
		checkRemoved := false
		for _, removed := range preProcessInfo.RemovedBitVotesOnes {
			if i == removed {
				newSolution.Participants[i] = true
				checkRemoved = true
				break
			}
		}
		if checkRemoved {
			continue
		}
		for _, removed := range preProcessInfo.RemovedBitVotesZeros {
			if i == removed {
				newSolution.Participants[i] = false
				checkRemoved = true
				break
			}
		}
		if checkRemoved {
			continue
		} else {
			newSolution.Participants[i] = unAggregatedSolution.Participants[participantCounter]
			participantCounter++
		}
	}

	solutionsCounter := 0
	for i := range newSolution.Solution {
		checkRemoved := false
		for _, removed := range preProcessInfo.RemovedAttestationsOnes {
			if i == removed {
				newSolution.Solution[i] = true
				checkRemoved = true
				break
			}
		}
		if checkRemoved {
			continue
		}
		for _, removed := range preProcessInfo.RemovedAttestationsLowWeight {
			if i == removed {
				newSolution.Solution[i] = false
				checkRemoved = true
				break
			}
		}
		if checkRemoved {
			continue
		} else {
			newSolution.Solution[i] = unAggregatedSolution.Solution[solutionsCounter]
			solutionsCounter++
		}
	}

	return newSolution
}
