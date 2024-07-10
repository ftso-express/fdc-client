package bitvotes

import "math/big"

type ConsensusSolution struct {
	Participants []bool
	Solution     []bool
	Value        Value
	Optimal      bool
}

func Ensemble(allBitVotes []*WeightedBitVote, fees []*big.Int, maxOperations int, seed int64) *ConsensusSolution {
	totalWeight := uint16(0)
	for _, bitVote := range allBitVotes {
		totalWeight += bitVote.Weight
	}

	preProcessedBitVotes, newFees, preProccesInfo := PreProcess(allBitVotes, fees)

	var firstMethod, secondMethod func([]*WeightedBitVote, []*big.Int, uint16, uint16, int, int64) *ConsensusSolution
	if len(allBitVotes) < len(fees) {
		firstMethod = BranchAndBoundProviders
		secondMethod = BranchAndBound
	} else {
		firstMethod = BranchAndBound
		secondMethod = BranchAndBoundProviders
	}

	solution := firstMethod(preProcessedBitVotes, newFees, preProccesInfo.RemovedOnesWeight, totalWeight, maxOperations, seed)
	if !solution.Optimal {
		solution2 := secondMethod(preProcessedBitVotes, newFees, preProccesInfo.RemovedOnesWeight, totalWeight, maxOperations, seed)
		if solution2.Value.Cmp(solution.Value) == 1 {
			solution = solution2
		}
	}

	expandedSolution := ExpandSolution(solution, preProccesInfo)
	expandedSolution.Value = expandedSolution.CalcValueFromFees(allBitVotes, fees, 0, totalWeight)

	return expandedSolution
}

func (solution *ConsensusSolution) CalcValueFromFees(allBitVotes []*WeightedBitVote, fees []*big.Int, assumedWeight, totalWeight uint16) Value {
	feeSum := big.NewInt(0)
	for i, attestation := range solution.Solution {
		if attestation {
			feeSum.Add(feeSum, fees[i])
		}
	}
	weight := assumedWeight
	for j, voter := range solution.Participants {
		if voter {
			weight += allBitVotes[j].Weight
		}
	}

	return CalcValue(feeSum, weight, totalWeight)
}

func (solution *ConsensusSolution) MaximizeSolution(allBitVotes []*WeightedBitVote, fees []*big.Int, assumedWeight, totalWeight uint16) {
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

	solution.Value = solution.CalcValueFromFees(allBitVotes, fees, assumedWeight, totalWeight)
}

func (solution *ConsensusSolution) MaximizeProviders(allBitVotes []*WeightedBitVote, fees []*big.Int, assumedWeight, totalWeight uint16) {
	for i, provider := range solution.Participants {
		if !provider {
			check := true
			for j, solution := range solution.Solution {
				if solution && allBitVotes[i].BitVote.BitVector.Bit(j) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Participants[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(allBitVotes, fees, assumedWeight, totalWeight)
}
