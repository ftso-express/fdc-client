package bitvotes_test

import (
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"testing"
	"time"
)

func TestBitWalk65(t *testing.T) {
	numAttestations := 40
	numVoters := 100
	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
	totalWeight := uint16(0)

	for j := 0; j < numVoters; j++ {
		var bitVote *bitvotes.WeightedBitVote
		if 0.65*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromRules(numAttestations, []int{2, 3})
		} else {
			bitVote = setBitVoteFromRules(numAttestations, []int{3, 7})
		}
		weightedBitvotes[j] = bitVote

		totalWeight += bitVote.Weight
	}

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := bitvotes.MetropolisHastingsSampling(weightedBitvotes, fees, 100000)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
}

func TestBitWalkRandom(t *testing.T) {
	numAttestations := 100
	numVoters := 100
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
	solution := bitvotes.MetropolisHastingsSampling(weightedBitvotes, fees, 100000)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
}

func TestBitWalkFix(t *testing.T) {
	numAttestations := 5
	numVoters := 100
	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
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
		weightedBitvotes[j] = bitVote

		totalWeight += bitVote.Weight
	}

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := bitvotes.MetropolisHastingsSampling(weightedBitvotes, fees, 100000)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
}

