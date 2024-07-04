package bitvotes_test

import (
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
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

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := bitvotes.Ensemble(weightedBitvotes, fees, 100000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
	count := 0
	for _, e := range solution.Solution {
		if e {
			count += 1
		}
	}

	fmt.Println("num attestations", count, len(solution.Solution))
	fmt.Println("num bitvotes", len(solution.Participants))
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

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := bitvotes.Ensemble(weightedBitvotes, fees, 100000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	// fmt.Println("solution", solution)

	require.Equal(t, 71, solution.Value)
	require.Equal(t, []bool{false, true, false, false, true, false, false, false}, solution.Solution)

}
