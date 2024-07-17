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

// BranchAndBoundProviders is similar to BranchAndBound, the difference is that it
// executes a branch and bound strategy on the space of subsets of attestation providers, hence
// it is particularly useful when there are not too many distinct providers.
func BranchAndBoundProviders(bitVotes []*AggregatedVote, fees []*AggregatedFee, assumedWeight, absoluteTotalWeight uint16, assumedFees *big.Int, maxOperations int, seed int64, initialBound Value) *ConsensusSolution {
	totalWeight := assumedWeight

	for _, vote := range bitVotes {
		totalWeight += vote.Weight
	}

	totalFee := big.NewInt(0).Set(assumedFees)
	bits := make(map[int]bool)
	for i, fee := range fees {
		totalFee.Add(totalFee, fee.Fee)
		bits[i] = true
	}
	randGen := rand.NewSource(seed)
	// randPerm := RandPerm(numProviders, randGen)
	// permBitVotes := PermuteBitVotes(bitVotes, randPerm)

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

	permResult := BranchProviders(processInfo, currentStatus, 0, bits, totalFee, totalWeight)

	result := ConsensusSolution{
		Votes: permResult.Votes,
		Bits:  permResult.Bits,
		Value: permResult.Value,
	}
	// for key, val := range randPerm {
	// 	result.Participants[key] = permResult.Participants[val]
	// }
	// for key, val := range permResult.Solution {
	// 	result.Solution[key] = val
	// }
	if currentStatus.NumOperations < maxOperations {
		result.Optimal = true
	} else {
		result.MaximizeVotes(bitVotes, fees, assumedFees, assumedWeight, absoluteTotalWeight)
	}

	return &result
}

func BranchProviders(processInfo *ProcessInfo, currentStatus *SharedStatus, branch int, bits map[int]bool, feeSum *big.Int, weight uint16) *BranchAndBoundPartialSolution {
	currentStatus.NumOperations++

	// end of recursion
	if branch == processInfo.NumProviders {
		value := CalcValue(feeSum, weight, processInfo.TotalWeight)
		if value.Cmp(currentStatus.CurrentBound) == 1 {
			currentStatus.CurrentBound = value
		}

		return &BranchAndBoundPartialSolution{
			Votes: make(map[int]bool),
			Bits:  bits,
			Value: value,
		}
	}

	// check if we already reached the maximal search space or if we exceeded the bound of the maximal possible value of the solution
	if currentStatus.NumOperations >= processInfo.MaxOperations || CalcValue(feeSum, weight, processInfo.TotalWeight).Cmp(currentStatus.CurrentBound) == -1 {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	newWeight := weight - processInfo.BitVotes[branch].Weight

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {
		// check if a branch is possible
		if newWeight > processInfo.LowerBoundWeight {
			result0 = BranchProviders(processInfo, currentStatus, branch+1, bits, feeSum, newWeight)
		}
	}

	// prepare a new branch
	newBits, newFeeSum := prepareDataForBranchWithProvider(processInfo, currentStatus, bits, feeSum, processInfo.BitVotes[branch].Indexes[0])

	result1 = BranchProviders(processInfo, currentStatus, branch+1, newBits, newFeeSum, weight)

	if randBit == 1 {
		if newWeight > processInfo.LowerBoundWeight {
			result0 = BranchProviders(processInfo, currentStatus, branch+1, bits, feeSum, newWeight)
		}
	}

	// max result
	return joinResultsProviders(result0, result1, branch)
}

func prepareDataForBranchWithProvider(processInfo *ProcessInfo, currentStatus *SharedStatus, bits map[int]bool, feeSum *big.Int, providerIndex int) (map[int]bool, *big.Int) {

	newSolution := make(map[int]bool)
	newFeeSum := new(big.Int).Set(feeSum)
	for sol := range bits {
		if processInfo.BitVotes[providerIndex].BitVector.Bit(sol) == 0 {
			newFeeSum.Sub(newFeeSum, processInfo.Fees[sol].Fee)
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
		result0.Votes[branch] = false
		return result0
	} else if result0 == nil || result0.Value.Cmp(result1.Value) == -1 {
		result1.Votes[branch] = true
		return result1
	} else {
		result0.Votes[branch] = false
		return result0
	}

}

// MaximizeVotes adds all votes that confirm all bits in the solution and updates the value.
func (solution *ConsensusSolution) MaximizeVotes(votes []*AggregatedVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) {
	for i := range votes {

		if _, isIncluded := solution.Votes[i]; !isIncluded {
			check := true
			for j := range solution.Bits {
				if votes[i].BitVector.Bit(fees[i].Indexes[j]) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Votes[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(votes, fees, assumedFees, assumedWeight, totalWeight)
}
