package bitvotes

import (
	"math/big"
)

type ConsensusSolution struct {
	Votes   []*AggregatedVote // set of votes that support the solution
	Bits    []*AggregatedFee  // set of bits that are confirmed
	Value   Value
	Optimal bool
}

func ensemble(allBitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, maxOperations int) (*FilterResults, *ConsensusSolution) {

	aggregatedVotes, aggregatedFees, filterResults := FilterAndAggregate(allBitVotes, fees, totalWeight)

	var solution *ConsensusSolution

	if len(allBitVotes) < len(fees) {

		solution = BranchAndBoundVotesDouble(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, Value{big.NewInt(0), big.NewInt(0)})
		if !solution.Optimal {
			solution2 := BranchAndBoundBitsDouble(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, solution.Value)
			if solution2 != nil && solution2.Value.Cmp(solution.Value) == 1 {
				solution = solution2
			}
		}
	} else {

		solution = BranchAndBoundBitsDouble(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, Value{big.NewInt(0), big.NewInt(0)})

		if !solution.Optimal {

			solution2 := BranchAndBoundVotesDouble(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, solution.Value)
			if solution2 != nil && solution2.Value.Cmp(solution.Value) == 1 {
				solution = solution2
			}
		}
	}

	return filterResults, solution
}

func EnsembleFull(allBitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, maxOperations int) Solution {

	filterResults, filterSolution := ensemble(allBitVotes, fees, totalWeight, maxOperations)

	return AssembleSolutionFull(filterResults, *filterSolution)
}

func EnsembleConsensusBitVote(allBitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, maxOperations int) BitVote {

	filterResults, filterSolution := ensemble(allBitVotes, fees, totalWeight, maxOperations)

	return AssembleSolution(filterResults, *filterSolution, uint16(len(fees)))
}

func (solution *branchAndBoundPartialSolution) CalcValueFromFees(allBitVotes []*AggregatedVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) Value {
	feeSum := big.NewInt(0).Set(assumedFees)
	for i := range solution.Bits {
		feeSum.Add(feeSum, fees[i].Fee)

	}
	weight := assumedWeight
	for j := range solution.Votes {
		weight += allBitVotes[j].Weight
	}

	return CalcValue(feeSum, weight, totalWeight)
}
