package bitvotes

import (
	"math/big"
	"math/rand"
)

func PermuteBitVotes(bitVotes []*WeightedBitVote, randPerm []int) []*WeightedBitVote {
	permBitVotes := make([]*WeightedBitVote, len(bitVotes))
	for i, e := range bitVotes {
		permBitVotes[randPerm[i]] = e
	}

	return permBitVotes
}

// BranchAndBoundProviders is similar than BranchAndBound, the difference is that it
// executes a branch and bound strategy on the space of subsets of attestation providers, hence
// it is particularly useful when there are not too many distinct providers.
func BranchAndBoundProviders(bitVotes []*WeightedBitVote, fees []*big.Int, assumedWeight, absoluteTotalWeight uint16, maxOperations int, seed int64) *ConsensusSolution {
	numAttestations := len(fees)
	numProviders := len(bitVotes)
	totalWeight := assumedWeight

	for _, vote := range bitVotes {
		totalWeight += vote.Weight
	}

	totalFee := big.NewInt(0)
	startingSolution := make(map[int]bool)
	for i, fee := range fees {
		totalFee.Add(totalFee, fee)
		startingSolution[i] = true
	}
	randGen := rand.NewSource(seed)
	randPerm := RandPerm(numProviders, randGen)
	permBitVotes := PermuteBitVotes(bitVotes, randPerm)

	currentBound := &SharedStatus{
		CurrentBound:     Value{CappedValue: big.NewInt(0), UncappedValue: big.NewInt(0)},
		NumOperations:    0,
		MaxOperations:    maxOperations,
		TotalWeight:      absoluteTotalWeight,
		LowerBoundWeight: absoluteTotalWeight / 2,
		BitVotes:         permBitVotes,
		Fees:             fees,
		RandGen:          randGen,
		NumProviders:     numProviders,
	}

	permResult := BranchProviders(startingSolution, totalFee, currentBound, 0, totalWeight)

	result := orderResultProviders(permResult, numProviders, numAttestations, randPerm)

	if currentBound.NumOperations < maxOperations {
		result.Optimal = true
	} else {
		result.MaximizeProviders(bitVotes, fees, assumedWeight, absoluteTotalWeight)
	}

	return result
}

func orderResultProviders(permResult *BranchAndBoundPartialSolution, numProviders, numAttestations int, permutation []int) *ConsensusSolution {

	result := ConsensusSolution{
		Participants: make([]bool, numProviders),
		Solution:     make([]bool, numAttestations),
		Value:        permResult.Value,
	}

	for key, val := range permutation {
		result.Participants[key] = permResult.Participants[val]
	}
	for key, val := range permResult.Solution {
		result.Solution[key] = val
	}

	return &result
}

func BranchProviders(solution map[int]bool, feeSum *big.Int, currentStatus *SharedStatus, branch int, currentMaxWeight uint16) *BranchAndBoundPartialSolution {
	currentStatus.NumOperations++

	// end of recursion
	if branch == currentStatus.NumProviders {
		value := CalcValue(feeSum, currentMaxWeight, currentStatus.TotalWeight)
		if value.Cmp(currentStatus.CurrentBound) == 1 {
			currentStatus.CurrentBound = value
		}

		return &BranchAndBoundPartialSolution{
			Participants: make(map[int]bool),
			Solution:     solution,
			Value:        value,
		}
	}

	// check if we already reached the maximal search space or if we exceeded the bound of the maximal possible value of the solution
	if currentStatus.NumOperations >= currentStatus.MaxOperations || CalcValue(feeSum, currentMaxWeight, currentStatus.TotalWeight).Cmp(currentStatus.CurrentBound) == -1 {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	newCurrentMaxWeight := currentMaxWeight - currentStatus.BitVotes[branch].Weight

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {
		// check if a branch is possible
		if newCurrentMaxWeight > currentStatus.LowerBoundWeight {
			result0 = BranchProviders(solution, feeSum, currentStatus, branch+1, newCurrentMaxWeight)
		}
	}

	// prepare a new branch
	newSolution, newFeeSum := prepareDataForBranchWithProvider(solution, feeSum, currentStatus, branch)

	result1 = BranchProviders(newSolution, newFeeSum, currentStatus, branch+1, currentMaxWeight)

	if randBit == 1 {
		if newCurrentMaxWeight > currentStatus.TotalWeight/2 {
			result0 = BranchProviders(solution, feeSum, currentStatus, branch+1, newCurrentMaxWeight)
		}
	}

	// max result
	return joinResultsProviders(result0, result1, branch)
}

func prepareDataForBranchWithProvider(solution map[int]bool, feeSum *big.Int, currentStatus *SharedStatus, branch int) (map[int]bool, *big.Int) {

	newSolution := make(map[int]bool)
	newFeeSum := new(big.Int).Set(feeSum)
	for sol := range solution {
		if currentStatus.BitVotes[branch].BitVote.BitVector.Bit(sol) == 0 {
			newFeeSum.Sub(newFeeSum, currentStatus.Fees[sol])
		} else {
			newSolution[sol] = true
		}
		currentStatus.NumOperations++
	}

	return newSolution, newFeeSum

}

func joinResultsProviders(result0, result1 *BranchAndBoundPartialSolution, branch int) *BranchAndBoundPartialSolution {
	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		result0.Participants[branch] = false
		return result0
	} else if result0 == nil || result0.Value.Cmp(result1.Value) == -1 {
		result1.Participants[branch] = true
		return result1
	} else {
		result0.Participants[branch] = false
		return result0
	}

}
