package bitvotes

import (
	"fmt"
	"math/big"
)

type ConsensusSolution struct {
	Participants map[int]bool
	Solution     map[int]bool
	Value        Value
	Optimal      bool
}

func ensemble(allBitVotes []*WeightedBitVote, fees []*big.Int, maxOperations int, seed int64) ([]*AggregatedBitVote, []*AggregatedFee, *FilterResults, *ConsensusSolution) {
	totalWeight := uint16(0)
	for _, bitVote := range allBitVotes {
		totalWeight += bitVote.Weight
	}

	aggregatedVotes, aggregatedFees, filterResults := FilterAndAggregate(allBitVotes, fees, totalWeight)

	var firstMethod, secondMethod func([]*AggregatedBitVote, []*AggregatedFee, uint16, uint16, *big.Int, int, int64) *ConsensusSolution
	if len(allBitVotes) < len(fees) {
		firstMethod = BranchAndBoundProviders
		secondMethod = BranchAndBound
	} else {
		firstMethod = BranchAndBound
		secondMethod = BranchAndBoundProviders
	}

	solution := firstMethod(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, seed)
	if !solution.Optimal {
		solution2 := secondMethod(aggregatedVotes, aggregatedFees, filterResults.GuaranteedWeight, totalWeight, filterResults.GuaranteedFees, maxOperations, seed)
		if solution2.Value.Cmp(solution.Value) == 1 {
			solution = solution2
		}
	}

	fmt.Printf("solution: %v\n", solution)

	return aggregatedVotes, aggregatedFees, filterResults, solution
}

func EnsembleFull(allBitVotes []*WeightedBitVote, fees []*big.Int, maxOperations int, seed int64) Solution {

	aggregatedVotes, aggregadedFees, filterResults, filterSolution := ensemble(allBitVotes, fees, maxOperations, seed)

	return AssembleSolutionFull(filterResults, *filterSolution, aggregadedFees, aggregatedVotes)
}

func EnsembleConsensulBitVote(allBitVotes []*WeightedBitVote, fees []*big.Int, maxOperations int, seed int64) *big.Int {

	_, aggregadedFees, filterResults, filterSolution := ensemble(allBitVotes, fees, maxOperations, seed)

	return AssembleSolution(filterResults, *filterSolution, aggregadedFees)
}

func (solution *ConsensusSolution) CalcValueFromFees(allBitVotes []*AggregatedBitVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) Value {
	feeSum := big.NewInt(0).Set(assumedFees)
	for i, attestation := range solution.Solution {
		if attestation {
			feeSum.Add(feeSum, fees[i].Fee)
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

func (solution *ConsensusSolution) MaximizeSolution(allBitVotes []*AggregatedBitVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) {
	for i, attestation := range solution.Solution {
		if !attestation {
			check := true
			for j, voter := range solution.Participants {
				if voter && allBitVotes[j].BitVector.Bit(i) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Solution[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(allBitVotes, fees, assumedFees, assumedWeight, totalWeight)
}

func (solution *ConsensusSolution) MaximizeProviders(allBitVotes []*AggregatedBitVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) {
	for i, provider := range solution.Participants {
		if !provider {
			check := true
			for j, solution := range solution.Solution {
				if solution && allBitVotes[i].BitVector.Bit(j) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Participants[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(allBitVotes, fees, assumedFees, assumedWeight, totalWeight)
}
