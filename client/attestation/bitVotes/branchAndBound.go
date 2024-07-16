package bitvotes

import (
	"math"
	"math/big"
	"math/rand"
)

type SharedStatus struct {
	CurrentBound  Value
	NumOperations int
	RandGen       rand.Source
}

type ProcessInfo struct {
	TotalWeight      uint16
	LowerBoundWeight uint16
	BitVotes         []*AggregatedVote
	Fees             []*AggregatedFee
	NumAttestations  int
	NumProviders     int

	MaxOperations int
}

type BranchAndBoundPartialSolution struct {
	Votes map[int]bool
	Bits  map[int]bool
	Value Value
}

type Value struct {
	CappedValue   *big.Int
	UncappedValue *big.Int
}

// Cmp compares Values lexicographically
//
//	-1 if v0 <  v1
//	 0 if v0 == v1
//	+1 if v0 >  v1
func (v0 Value) Cmp(v1 Value) int {
	firstCompare := v0.CappedValue.Cmp(v1.CappedValue)

	if firstCompare != 0 {
		return firstCompare
	}

	return v0.UncappedValue.Cmp(v1.UncappedValue)

}

func CalcValue(feeSum *big.Int, weight, totalWeight uint16) Value {
	weightCaped := min(int64(math.Ceil(float64(totalWeight)*valueCap)), int64(weight))

	return Value{CappedValue: new(big.Int).Mul(feeSum, big.NewInt(weightCaped)),
		UncappedValue: new(big.Int).Mul(feeSum, big.NewInt(weightCaped)),
	}
}

func RandPerm(n int, randGen rand.Source) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := randGen.Int63() % int64(i+1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

// func PermuteBits(bitVotes []*AggregatedBitVote, randPerm []int) []*WeightedBitVote {
// 	permBitVotes := make([]*AggregatedBitVote, len(bitVotes))
// 	for i, e := range bitVotes {
// 		permBitVotes[i] = &AggregatedBitVote{Weight: e.Weight, BitVector: big.NewInt(0)}}
// 		for key, val := range randPerm {
// 			if e.BitVote.BitVector.Bit(key) == 1 {
// 				permBitVotes[i].BitVote.BitVector.Add(permBitVotes[i].BitVote.BitVector, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(val)), nil))
// 			}
// 		}
// 	}

// 	return permBitVotes
// }

// BranchAndBound is a function that takes a set of weighted bit votes and a list of fees and
// tries to get an optimal subset of votes with the weight more than the half of the total weight.
// The algorithm executes a branch and bound strategy on the space of subsets of attestations, hence
// it is particularly useful when there are not too many attestations. In the case the algorithm is able search
// through the entire solution space before reaching the given max operations counter, the algorithm
// gives an optimal solution. In the case that the solution space is too big, the algorithm gives a
// the best solution it finds. The search strategy is pseudo-randomized, where the pseudo-random
// function is controlled by the given seed.
func BranchAndBound(bitVotes []*AggregatedVote, fees []*AggregatedFee, assumedWeight, absoluteTotalWeight uint16, assumedFees *big.Int, maxOperations int, seed int64, initialBound Value) *ConsensusSolution {

	weight := assumedWeight

	votes := make(map[int]bool)
	for i, vote := range bitVotes {
		weight += vote.Weight
		votes[i] = true
	}

	totalFee := big.NewInt(0).Set(assumedFees)
	for _, fee := range fees {
		totalFee.Add(totalFee, fee.Fee)
	}

	randGen := rand.NewSource(seed)
	// randPerm := RandPerm(numAttestations, randGen)
	// permBitVotes := PermuteBits(bitVotes, randPerm)

	currentStatus := &SharedStatus{
		CurrentBound:  initialBound,
		NumOperations: 0,
		RandGen:       randGen,
	}

	processInfo := &ProcessInfo{
		TotalWeight:      absoluteTotalWeight,
		LowerBoundWeight: absoluteTotalWeight / 2,
		BitVotes:         bitVotes,
		Fees:             fees,
		NumAttestations:  len(fees),
		NumProviders:     len(bitVotes),
		MaxOperations:    maxOperations,
	}

	permResult := Branch(processInfo, currentStatus, 0, votes, weight, totalFee)

	if permResult == nil {
		return nil
	}

	result := ConsensusSolution{
		Participants: permResult.Votes,
		Solution:     permResult.Bits,
		Value:        permResult.Value,
	}

	// for key, val := range permResult.Participants {
	// 	result.Participants[key] = val
	// }
	// for key, val := range randPerm {
	// 	result.Solution[key] = permResult.Solution[val]
	// }
	if currentStatus.NumOperations < maxOperations {
		result.Optimal = true
	} else {
		result.MaximizeSolution(bitVotes, fees, assumedFees, assumedWeight, absoluteTotalWeight)
	}

	return &result
}

func Branch(processInfo *ProcessInfo, currentStatus *SharedStatus, branch int, participants map[int]bool, currentWeight uint16, feeSum *big.Int) *BranchAndBoundPartialSolution {

	currentStatus.NumOperations++

	// end of recursion
	if branch == processInfo.NumAttestations {

		value := CalcValue(feeSum, currentWeight, processInfo.TotalWeight)

		if value.Cmp(currentStatus.CurrentBound) == 1 {
			currentStatus.CurrentBound = value
		}

		return &BranchAndBoundPartialSolution{Votes: participants, Bits: make(map[int]bool), Value: value}
	}

	// check if we already reached the maximal search space
	if currentStatus.NumOperations >= processInfo.MaxOperations {
		return nil
	}

	// check if the estimated maximal value of a branch is lower, then the current highest value
	if CalcValue(feeSum, currentWeight, processInfo.TotalWeight).Cmp(currentStatus.CurrentBound) == -1 {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {

		result0 = Branch(processInfo, currentStatus, branch+1, participants, currentWeight, new(big.Int).Sub(feeSum, processInfo.Fees[branch].Fee))
	}

	// prepare and check if a branch is possible
	newParticipants, newCurrentWeight := prepareDataForBranchWithOne(processInfo, currentStatus, participants, processInfo.Fees[branch].Indexes[0], currentWeight)

	if newCurrentWeight > processInfo.LowerBoundWeight {

		result1 = Branch(processInfo, currentStatus, branch+1, newParticipants, newCurrentWeight, feeSum)
	}

	if randBit == 1 {
		result0 = Branch(processInfo, currentStatus, branch+1, participants, currentWeight, new(big.Int).Sub(feeSum, processInfo.Fees[branch].Fee))
	}

	// max result
	return joinResultsAttestations(result0, result1, branch)
}

func prepareDataForBranchWithOne(processInfo *ProcessInfo, currentStatus *SharedStatus, participants map[int]bool, bit int, currentWeight uint16) (map[int]bool, uint16) {

	newParticipants := make(map[int]bool)
	newCurrentWeight := currentWeight

	for participant := range participants {
		if processInfo.BitVotes[participant].BitVector.Bit(bit) == 1 {

			newParticipants[participant] = true
		} else {

			newCurrentWeight -= processInfo.BitVotes[participant].Weight

		}
		currentStatus.NumOperations++
	}

	return newParticipants, newCurrentWeight

}

func joinResultsAttestations(result0, result1 *BranchAndBoundPartialSolution, branch int) *BranchAndBoundPartialSolution {

	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		delete(result0.Bits, branch)
		return result0
	} else if result0 == nil || result0.Value.Cmp(result1.Value) == -1 {
		result1.Bits[branch] = true

		return result1
	} else {
		delete(result0.Bits, branch)
		return result0
	}

}
