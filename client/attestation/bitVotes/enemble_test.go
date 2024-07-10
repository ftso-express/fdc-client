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
	numAttestations := 10
	numVoters := 30
	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
	prob := 0.8

	totalWeight := uint16(0)
	for j := 0; j < numVoters; j++ {
		bitVote := randomBitVotes(numAttestations, prob)
		weightedBitvotes[j] = bitVote

		totalWeight += bitVote.Weight
	}

	fees := make([]*big.Int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)
	}

	solution := bitvotes.Ensemble(weightedBitvotes, fees, 100000000, time.Now().Unix())
	require.Equal(t, numVoters, len(solution.Participants))
	require.Equal(t, numAttestations, len(solution.Solution))

	solutionCheck := bitvotes.BranchAndBound(weightedBitvotes, fees, 0, totalWeight, 100000000, time.Now().Unix())

	require.Equal(t, solutionCheck.Value, solution.Value)
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
	solution := bitvotes.Ensemble(weightedBitvotes, fees, 100000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	// fmt.Println("solution", solution)

	require.Equal(t, bitvotes.Value{big.NewInt(2 * 71), big.NewInt(2 * 71)}, solution.Value)
	require.Equal(t, []bool{false, true, false, false, true, false, false, false}, solution.Solution)
	for j := 0; j < numVoters; j++ {
		if 0.30*float64(numVoters) > float64(j) {
			require.Equal(t, true, solution.Participants[j])
		} else if 0.61*float64(numVoters) > float64(j) {
			require.Equal(t, true, solution.Participants[j])
		} else if 0.90*float64(numVoters) > float64(j) {
			require.Equal(t, false, solution.Participants[j])
		} else {
			require.Equal(t, true, solution.Participants[j])
		}
	}
}
