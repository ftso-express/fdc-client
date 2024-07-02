package attestation

import (
	"math/rand"
)

func PermuteBitVotes(allBitVotes []*WeightedBitVote, randPerm []int) []*WeightedBitVote {
	permBitVotes := make([]*WeightedBitVote, len(allBitVotes))
	for i, e := range allBitVotes {
		permBitVotes[randPerm[i]] = e
	}

	return permBitVotes
}

// BranchAndBoundProviders is similar than BranchAndBound, the difference is that it
// executes a branch and bound strategy on the space of subsets of attestation providers, hence
// it is particularly useful when there are not too many distinct providers.
func BranchAndBoundProviders(allBitVotes []*WeightedBitVote, fees []int, absoluteTotalWeight uint16, maxOperations int, seed int64) *BranchAndBoundSolution {
	numAttestations := len(fees)
	numProviders := len(allBitVotes)
	totalWeight := uint16(0)

	for _, vote := range allBitVotes {
		totalWeight += vote.Weight
	}

	totalFee := 0
	startingSolution := make(map[int]bool)
	for i, fee := range fees {
		totalFee += fee
		startingSolution[i] = true
	}
	randGen := rand.NewSource(seed)
	randPerm := RandPerm(numProviders, randGen)
	permBitVotes := PermuteBitVotes(allBitVotes, randPerm)

	currentBound := &SharedStatus{CurrentBound: 0, NumOperations: 0, MaxOperations: maxOperations,
		TotalWeight: absoluteTotalWeight, BitVotes: permBitVotes, Fees: fees, RandGen: randGen, NumProviders: numProviders}

	permResult := BranchProviders(startingSolution, totalFee, currentBound, 0, totalWeight)

	result := BranchAndBoundSolution{Participants: make([]bool, numProviders),
		Solution: make([]bool, numAttestations), Value: permResult.Value}
	for key, val := range randPerm {
		result.Participants[key] = permResult.Participants[val]
	}
	for key, val := range permResult.Solution {
		result.Solution[key] = val
	}
	if currentBound.NumOperations < maxOperations {
		result.Optimal = true
	} else {
		result.Maximize(allBitVotes, fees)
	}

	return &result
}

func BranchProviders(solution map[int]bool, feeSum int, currentStatus *SharedStatus, branch int, currentMaxWeight uint16) *BranchAndBoundPartialSolution {
	currentStatus.NumOperations++

	// end of recursion
	if branch == currentStatus.NumProviders {
		value := CalcValue(feeSum, currentMaxWeight, currentStatus.TotalWeight)
		if value > currentStatus.CurrentBound {
			currentStatus.CurrentBound = value
		}

		return &BranchAndBoundPartialSolution{Participants: make(map[int]bool), Solution: solution, Value: value}
	}

	// check if we already reached the maximal search space or if we exceeded the bound of the maximal possible value of the solution
	if currentStatus.NumOperations >= currentStatus.MaxOperations || CalcValue(feeSum, currentMaxWeight, currentStatus.TotalWeight) < currentStatus.CurrentBound {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	newCurrentMaxWeight := currentMaxWeight - currentStatus.BitVotes[branch].Weight

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {
		// check if a branch is possible
		if newCurrentMaxWeight > currentStatus.TotalWeight/2 {
			result0 = BranchProviders(solution, feeSum, currentStatus, branch+1, newCurrentMaxWeight)
		}
	}

	// prepare a new branch
	newSolution := make(map[int]bool)
	newFeeSum := feeSum
	for sol := range solution {
		if currentStatus.BitVotes[branch].BitVote.BitVector.Bit(sol) == 0 {
			newFeeSum -= currentStatus.Fees[sol]
		} else {
			newSolution[sol] = true
		}
		currentStatus.NumOperations++
	}

	result1 = BranchProviders(newSolution, newFeeSum, currentStatus, branch+1, currentMaxWeight)

	if randBit == 1 {
		if newCurrentMaxWeight > currentStatus.TotalWeight/2 {
			result0 = BranchProviders(solution, feeSum, currentStatus, branch+1, newCurrentMaxWeight)
		}
	}

	// max result
	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		result0.Participants[branch] = false
		return result0
	} else if result0 == nil || result0.Value < result1.Value {
		result1.Participants[branch] = true
		return result1
	} else {
		result0.Participants[branch] = false
		return result0
	}
}
