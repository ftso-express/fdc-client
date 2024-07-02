package attestation

import (
	"math/big"
	"math/rand"
)

type SharedStatus struct {
	CurrentBound    int
	NumOperations   int
	MaxOperations   int
	TotalWeight     uint16
	BitVotes        []*WeightedBitVote
	Fees            []int
	RandGen         rand.Source
	NumAttestations int
	NumProviders    int
}

type BranchAndBoundPartialSolution struct {
	Participants map[int]bool
	Solution     map[int]bool
	Value        int
}

type BranchAndBoundSolution struct {
	Participants []bool
	Solution     []bool
	Value        int
	Optimal      bool
}

func CalcValue(feeSum int, weight, totalWeight uint16) int {
	weightCaped := min(int(float64(totalWeight)*valueCap), int(weight))

	return feeSum * weightCaped
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

func PermuteBits(allBitVotes []*WeightedBitVote, randPerm []int) []*WeightedBitVote {
	permBitVotes := make([]*WeightedBitVote, len(allBitVotes))
	for i, e := range allBitVotes {
		permBitVotes[i] = &WeightedBitVote{Weight: e.Weight, BitVote: BitVote{Length: e.BitVote.Length, BitVector: big.NewInt(0)}}
		for key, val := range randPerm {
			if e.BitVote.BitVector.Bit(key) == 1 {
				permBitVotes[i].BitVote.BitVector.Add(permBitVotes[i].BitVote.BitVector, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(val)), nil))
			}
		}
	}

	return permBitVotes
}

// BranchAndBound is a function that takes a set of weighted bit votes and a list of fees and
// tries to get an optimal subset of votes with the weight more than the half of the total weight.
// The algorithm executes a branch and bound strategy on the space of subsets of attestations, hence
// it is particularly useful when there are not too many attestations. In the case the algorithm is able search
// through the entire solution space before reaching the given max operations counter, the algorithm
// gives an optimal solution. In the case that the solution space is too big, the algorithm gives a
// the best solution it finds. The search strategy is pseudo-randomized, where the pseudo-random
// function is controlled by the given seed.
func BranchAndBound(allBitVotes []*WeightedBitVote, fees []int, absoluteTotalWeight uint16, maxOperations int, seed int64) *BranchAndBoundSolution {
	numAttestations := len(fees)
	numVoters := len(allBitVotes)
	totalWeight := uint16(0)

	participants := make(map[int]bool)
	for i, vote := range allBitVotes {
		totalWeight += vote.Weight
		participants[i] = true
	}

	totalFee := 0
	for _, fee := range fees {
		totalFee += fee
	}
	randGen := rand.NewSource(seed)
	randPerm := RandPerm(numAttestations, randGen)
	permBitVotes := PermuteBits(allBitVotes, randPerm)

	currentBound := &SharedStatus{CurrentBound: 0, NumOperations: 0, MaxOperations: maxOperations,
		TotalWeight: absoluteTotalWeight, BitVotes: permBitVotes, Fees: fees, RandGen: randGen, NumAttestations: numAttestations}

	permResult := Branch(participants, currentBound, 0, totalWeight, totalFee)

	result := BranchAndBoundSolution{Participants: make([]bool, numVoters),
		Solution: make([]bool, numAttestations), Value: permResult.Value}
	for key, val := range permResult.Participants {
		result.Participants[key] = val
	}
	for key, val := range randPerm {
		result.Solution[key] = permResult.Solution[val]
	}
	if currentBound.NumOperations < maxOperations {
		result.Optimal = true
	} else {
		result.Maximize(allBitVotes, fees)
	}

	return &result
}

func Branch(participants map[int]bool, currentStatus *SharedStatus, branch int, currentWeight uint16, feeSum int) *BranchAndBoundPartialSolution {
	currentStatus.NumOperations++

	// end of recursion
	if branch == currentStatus.NumAttestations {
		value := CalcValue(feeSum, currentWeight, currentStatus.TotalWeight)
		if value > currentStatus.CurrentBound {
			currentStatus.CurrentBound = value
		}

		return &BranchAndBoundPartialSolution{Participants: participants, Solution: make(map[int]bool), Value: value}
	}

	// check if we already reached the maximal search space or if we exceeded the bound of the maximal possible value of the solution
	if currentStatus.NumOperations >= currentStatus.MaxOperations || CalcValue(feeSum, currentWeight, currentStatus.TotalWeight) < currentStatus.CurrentBound {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {
		result0 = Branch(participants, currentStatus, branch+1, currentWeight, feeSum-currentStatus.Fees[branch])
	}

	// prepare and check if a branch is possible
	newParticipants := make(map[int]bool)
	newCurrentWeight := currentWeight
	for participant := range participants {
		if currentStatus.BitVotes[participant].BitVote.BitVector.Bit(branch) == 1 {
			newParticipants[participant] = true
		} else {
			newCurrentWeight -= currentStatus.BitVotes[participant].Weight
		}
		currentStatus.NumOperations++
	}
	if newCurrentWeight > currentStatus.TotalWeight/2 {
		result1 = Branch(newParticipants, currentStatus, branch+1, newCurrentWeight, feeSum)
	}

	if randBit == 1 {
		result0 = Branch(participants, currentStatus, branch+1, currentWeight, feeSum-currentStatus.Fees[branch])
	}

	// max result
	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		result0.Solution[branch] = false
		return result0
	} else if result0 == nil || result0.Value < result1.Value {
		result1.Solution[branch] = true
		return result1
	} else {
		result0.Solution[branch] = false
		return result0
	}
}

func (solution *BranchAndBoundSolution) Maximize(allBitVotes []*WeightedBitVote, fees []int) {
	for i, attestation := range solution.Solution {
		if !attestation {
			check := true
			for j, voter := range solution.Participants {
				if voter && allBitVotes[j].BitVote.BitVector.Bit(i) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Solution[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(allBitVotes, fees)
}

func (solution *BranchAndBoundSolution) CalcValueFromFees(allBitVotes []*WeightedBitVote, fees []int) int {
	val := 0
	for i, attestation := range solution.Solution {
		if attestation {
			val += fees[i]
		}
	}
	weight := 0
	totalWeight := 0
	for j, voter := range solution.Participants {
		if voter {
			weight += int(allBitVotes[j].Weight)
		}
		totalWeight += int(allBitVotes[j].Weight)
	}

	weightCaped := min(int(float64(totalWeight)*valueCap), weight)

	return val * weightCaped
}
