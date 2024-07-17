package bitvotes_test

import (
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEnsembleRandom(t *testing.T) {
	numAttestations := 40
	numVoters := 100
	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
	aggregatedBitvotes := make([]*bitvotes.AggregatedVote, numVoters)

	prob := 0.8

	totalWeight := uint16(0)
	for j := 0; j < numVoters; j++ {
		bitVote := randomBitVotes(numAttestations, prob)
		weightedBitvotes[j] = bitVote

		agg := bitvotes.AggregatedVote{BitVector: bitVote.BitVote.BitVector, Weight: bitVote.Weight, Indexes: []int{j}}

		aggregatedBitvotes[j] = &agg

		totalWeight += bitVote.Weight
	}

	fees := make([]*big.Int, numAttestations)
	aggFees := make([]*bitvotes.AggregatedFee, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)

		aggFee := bitvotes.AggregatedFee{big.NewInt(1), []int{j}}

		aggFees[j] = &aggFee
	}

	solution := bitvotes.EnsembleFull(weightedBitvotes, fees, totalWeight, 100000000, time.Now().Unix())
	// require.Equal(t, numVoters, len(solution.Participants))
	// require.Equal(t, numAttestations, len(solution.Solution))

	solutionCheck := bitvotes.BranchAndBoundBits(aggregatedBitvotes, aggFees, 0, totalWeight, big.NewInt(0), 100000000, time.Now().Unix(), bitvotes.Value{big.NewInt(0), big.NewInt(0)})

	require.Equal(t, solutionCheck.Value, solution.Value)

	fmt.Printf("solution: %v\n", solution)
}

func TestEnsembleFixed(t *testing.T) {
	numAttestations := 8
	numVoters := 100
	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)

	totalWeight := uint16(0)
	for j := 0; j < numVoters; j++ {
		var bitVote *bitvotes.WeightedBitVote
		if 0.30*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3, 4})
		} else if 0.61*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositions(numAttestations, []int{1, 4})
		} else if 0.90*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositions(numAttestations, []int{3, 4})
		} else {
			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3, 4, 5, 6, 7})
		}
		weightedBitvotes[j] = bitVote
		totalWeight += bitVote.Weight
	}

	fees := make([]*big.Int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)
	}

	start := time.Now()
	solution := bitvotes.EnsembleFull(weightedBitvotes, fees, totalWeight, 100000000, time.Now().Unix())

	fmt.Printf("solution: %v\n", solution.Bits)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)

	require.Equal(t, bitvotes.Value{big.NewInt(2 * 71), big.NewInt(2 * 71)}, solution.Value)

	fmt.Println(solution.Bits)
	require.ElementsMatch(t, []int{1, 4}, solution.Bits)
	for j := 0; j < numVoters; j++ {
		if 0.30*float64(numVoters) > float64(j) {
			require.Contains(t, solution.Votes, j)
		} else if 0.61*float64(numVoters) > float64(j) {
			require.Contains(t, solution.Votes, j)
		} else if 0.90*float64(numVoters) > float64(j) {
			require.NotContains(t, solution.Votes, j)
		} else {
			require.Contains(t, solution.Votes, j)
		}
	}
}
