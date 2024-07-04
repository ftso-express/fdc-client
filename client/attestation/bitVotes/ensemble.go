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

	return expandedSolution
}
