package bitvotes

import (
	"math/big"
)

type ConsensusSolution struct {
	Votes   map[int]bool // set of votes that support the solution
	Bits    map[int]bool // set of bits that are confirmed
	Value   Value
	Optimal bool
}

func ensemble(allBitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, maxOperations int, seed int64) ([]*AggregatedVote, []*AggregatedFee, *FilterResults, *ConsensusSolution) {
	participationWeight := uint16(0)
	for _, bitVote := range allBitVotes {
		participationWeight += bitVote.Weight
	}

	aggregatedVotes, aggregatedFees, filterResults := FilterAndAggregate(allBitVotes, fees, totalWeight)

	var firstMethod, secondMethod func([]*AggregatedVote, []*AggregatedFee, uint16, uint16, *big.Int, int, int64, Value) *ConsensusSolution
	if len(allBitVotes) < len(fees) {
		firstMethod = BranchAndBoundProviders
		secondMethod = BranchAndBound
	} else {
		firstMethod = BranchAndBound
		secondMethod = BranchAndBoundProviders
	}

	solution := firstMethod(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, seed, Value{big.NewInt(0), big.NewInt(0)})
	if !solution.Optimal {
		solution2 := secondMethod(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, seed, solution.Value)
		if solution2 != nil && solution2.Value.Cmp(solution.Value) == 1 {
			solution = solution2
		}
	}

	return aggregatedVotes, aggregatedFees, filterResults, solution
}

func EnsembleFull(allBitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, maxOperations int, seed int64) Solution {

	aggregatedVotes, aggregadedFees, filterResults, filterSolution := ensemble(allBitVotes, fees, totalWeight, maxOperations, seed)

	return AssembleSolutionFull(filterResults, *filterSolution, aggregadedFees, aggregatedVotes)
}

func EnsembleConsensulBitVote(allBitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, maxOperations int, seed int64) *big.Int {

	_, aggregadedFees, filterResults, filterSolution := ensemble(allBitVotes, fees, totalWeight, maxOperations, seed)

	return AssembleSolution(filterResults, *filterSolution, aggregadedFees)
}

func (solution *ConsensusSolution) CalcValueFromFees(allBitVotes []*AggregatedVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) Value {
	feeSum := big.NewInt(0).Set(assumedFees)
	for i, attestation := range solution.Bits {
		if attestation {
			feeSum.Add(feeSum, fees[i].Fee)
		}
	}
	weight := assumedWeight
	for j, voter := range solution.Votes {
		if voter {
			weight += allBitVotes[j].Weight
		}
	}

	return CalcValue(feeSum, weight, totalWeight)
}
