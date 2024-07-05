package bitvotes

type ConsensusSolution struct {
	Participants []bool
	Solution     []bool
	Value        int
	Optimal      bool
}

func Ensemble(allBitVotes []*WeightedBitVote, fees []int, maxOperations int, seed int64) *ConsensusSolution {
	totalWeight := uint16(0)
	for _, bitVote := range allBitVotes {
		totalWeight += bitVote.Weight
	}

	preProcessedBitVotes, newFees, preProccesInfo := PreProcess(allBitVotes, fees)

	var firstMethod, secondMethod func([]*WeightedBitVote, []int, uint16, uint16, int, int64) *ConsensusSolution
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
		if solution2.Value > solution.Value {
			solution = solution2
		}
	}

	expandedSolution := ExpandSolution(solution, preProccesInfo)
	expandedSolution.Value = expandedSolution.CalcValueFromFees(allBitVotes, fees, 0, totalWeight)

	return expandedSolution
}

func (solution *ConsensusSolution) CalcValueFromFees(allBitVotes []*WeightedBitVote, fees []int, assumedWeight, totalWeight uint16) int {
	val := 0
	for i, attestation := range solution.Solution {
		if attestation {
			val += fees[i]
		}
	}
	weight := assumedWeight
	for j, voter := range solution.Participants {
		if voter {
			weight += allBitVotes[j].Weight
		}
	}

	weightCaped := min(int(float64(totalWeight)*valueCap), int(weight))

	return val * weightCaped
}

func (solution *ConsensusSolution) MaximizeSolution(allBitVotes []*WeightedBitVote, fees []int, assumedWeight, totalWeight uint16) {
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

func (solution *ConsensusSolution) MaximizeProviders(allBitVotes []*WeightedBitVote, fees []int, assumedWeight, totalWeight uint16) {
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
