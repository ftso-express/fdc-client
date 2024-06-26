package attestation_test

import (
	"fmt"
	"local/fdc/client/attestation"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBranchAndBoundRandom(t *testing.T) {
	numAttestations := 30
	numVoters := 100
	weightedBitvotes := []*attestation.WeightedBitVote{}
	prob := 0.8

	for j := 0; j < numVoters; j++ {
		atts := random_attestations(numAttestations, prob)

		bitVote, err := attestation.BitVoteFromAttestations(atts)
		require.NoError(t, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := attestation.BranchAndBound(weightedBitvotes, fees, 100000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
	count := 0
	for _, e := range solution.Solution {
		if e {
			count += 1
		}
	}

	fmt.Println("num attestations", count)
}

func TestBranchAndBound65(t *testing.T) {
	numAttestations := 40
	numVoters := 100
	weightedBitvotes := []*attestation.WeightedBitVote{}

	for j := 0; j < numVoters; j++ {
		var atts []*attestation.Attestation

		if 0.65*float64(numVoters) > float64(j) {
			atts = setAttestations(numAttestations, []int{2, 3})
		} else {
			atts = setAttestations(numAttestations, []int{3, 7})
		}

		bitVote, err := attestation.BitVoteFromAttestations(atts)
		require.NoError(t, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := attestation.BranchAndBound(weightedBitvotes, fees, 50000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
	count := 0
	for _, e := range solution.Solution {
		if e {
			count += 1
		}
	}

	fmt.Println("num attestations", count)
}

func TestBranchAndBoundFix(t *testing.T) {
	numAttestations := 5
	numVoters := 100
	weightedBitvotes := []*attestation.WeightedBitVote{}

	for j := 0; j < numVoters; j++ {
		var atts []*attestation.Attestation

		if 0.30*float64(numVoters) > float64(j) {
			atts = setAttestationsFix(numAttestations, []int{0, 1, 2, 4})
		} else if 0.60*float64(numVoters) > float64(j) {
			atts = setAttestationsFix(numAttestations, []int{0, 1, 2, 3})
		} else if 0.90*float64(numVoters) > float64(j) {
			atts = setAttestationsFix(numAttestations, []int{0, 2})
		} else {
			atts = setAttestationsFix(numAttestations, []int{1, 3})
		}

		bitVote, err := attestation.BitVoteFromAttestations(atts)
		require.NoError(t, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	fees := make([]int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = 1
	}

	start := time.Now()
	solution := attestation.BranchAndBound(weightedBitvotes, fees, 50000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
	count := 0
	for _, e := range solution.Solution {
		if e {
			count += 1
		}
	}

	fmt.Println("num attestations", count)
}
