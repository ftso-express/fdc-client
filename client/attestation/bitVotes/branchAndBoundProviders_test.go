package bitvotes_test

import (
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBranchAndBoundProvidersFix(t *testing.T) {
	numAttestations := 100
	numVoters := 30
	weightedBitVotes := make([]*bitvotes.WeightedBitVote, numVoters)

	totalWeight := uint16(0)
	for j := 0; j < numVoters; j++ {
		var bitVote *bitvotes.WeightedBitVote

		if 0.30*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 4})
		} else if 0.60*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3})
		} else if 0.90*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 2})
		} else {
			bitVote = setBitVoteFromPositions(numAttestations, []int{1, 3})
		}

		bitVote.Index = j
		weightedBitVotes[j] = bitVote

		totalWeight += bitVote.Weight
	}

	fees := make([]*big.Int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)
	}

	aggregatedVotes, aggregatedFees, filterResults := bitvotes.FilterAndAggregate(weightedBitVotes, fees, totalWeight)

	initialBound := bitvotes.Value{big.NewInt(0), big.NewInt(0)}

	start := time.Now()
	solution := bitvotes.BranchAndBoundVotes(
		aggregatedVotes,
		aggregatedFees,
		filterResults.GuaranteedWeight,
		totalWeight,
		filterResults.GuaranteedFees,
		50000000,
		initialBound,
		true,
	)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)

	finalSolution := bitvotes.AssembleSolutionFull(filterResults, *solution)

	fmt.Printf("finalSolution.Bits: %v\n", finalSolution.Bits)

	initialBound2 := bitvotes.Value{big.NewInt(0), big.NewInt(0)}

	solution2 := bitvotes.BranchAndBoundBits(
		aggregatedVotes,
		aggregatedFees,
		filterResults.GuaranteedWeight,
		totalWeight,
		filterResults.GuaranteedFees,
		50000000,
		initialBound2,
		true,
	)

	finalSolution2 := bitvotes.AssembleSolutionFull(filterResults, *solution2)

	fmt.Printf("finalSolution2.Bits: %v\n", finalSolution2.Bits)

}

func TestBranchAndBoundProvidersRandom(t *testing.T) {
	numAttestations := 1000
	numVoters := 100
	weightedBitVotes := make([]*bitvotes.WeightedBitVote, numVoters)
	prob := 0.86

	totalWeight := uint16(0)
	for j := 0; j < numVoters; j++ {
		bitVote := randomBitVotes(numAttestations, prob)
		weightedBitVotes[j] = bitVote
		totalWeight += bitVote.Weight
	}

	fees := make([]*big.Int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)
	}

	aggregatedVotes, aggregatedFees, filterResults := bitvotes.FilterAndAggregate(weightedBitVotes, fees, totalWeight)

	initialBound := bitvotes.Value{big.NewInt(0), big.NewInt(0)}

	start := time.Now()
	solution := bitvotes.BranchAndBoundVotes(
		aggregatedVotes,
		aggregatedFees,
		filterResults.GuaranteedWeight,
		totalWeight,
		filterResults.GuaranteedFees,
		200000000,
		initialBound,
		true,
	)
	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)

	initialBound2 := bitvotes.Value{big.NewInt(0), big.NewInt(0)}

	start2 := time.Now()

	solution2 := bitvotes.BranchAndBoundBits(
		aggregatedVotes,
		aggregatedFees,
		filterResults.GuaranteedWeight,
		totalWeight,
		filterResults.GuaranteedFees,
		100000000,
		initialBound2,
		false,
	)
	fmt.Println("time passed:", time.Since(start2).Seconds())

	fmt.Println("solution2", solution2)

	finalSolution := bitvotes.AssembleSolutionFull(filterResults, *solution)

	fmt.Printf("finalSolution.Bits: %v\n", finalSolution.Bits)

	finalSolution2 := bitvotes.AssembleSolutionFull(filterResults, *solution2)

	fmt.Printf("finalSolution2.Bits: %v\n", finalSolution2.Bits)

	require.Equal(t, solution.Value, solution2.Value)

}
