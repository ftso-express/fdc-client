package bitvotes_test

import (
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBranchAndBoundRandom(t *testing.T) {
	numAttestations := 70
	numVoters := 100
	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
	prob := 1.0

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

	start := time.Now()
	solution := bitvotes.BranchAndBound(weightedBitvotes, fees, 0, totalWeight, 100000000, time.Now().Unix())

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println("value", solution.Value)
	count := 0
	for _, e := range solution.Solution {
		if e {
			count += 1
		}
	}

	fmt.Println("num attestations", count)
}

func TestBranchAndBound65(t *testing.T) {
	numAttestations := 100
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

	fees := make([]*big.Int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)
	}

	start := time.Now()

	solution := bitvotes.BranchAndBound(weightedBitvotes, fees, 0, totalWeight, 50000000, time.Now().Unix())

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
	numAttestations := 30
	numVoters := 30
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

	fees := make([]*big.Int, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fees[j] = big.NewInt(1)
	}

	start := time.Now()
	solution := bitvotes.BranchAndBound(weightedBitvotes, fees, 0, totalWeight, 50000000, time.Now().Unix())

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

func TestCalcValue(t *testing.T) {
	const totalWeight = 100

	tests := []struct {
		feeSum        *big.Int
		weight        uint16
		uncappedValue *big.Int
		cappedValue   *big.Int
	}{
		{
			big.NewInt(0),
			0,
			big.NewInt(0),
			big.NewInt(0),
		},
		{
			big.NewInt(0),
			15,
			big.NewInt(0),
			big.NewInt(0),
		},
		{
			big.NewInt(1),
			0,
			big.NewInt(0),
			big.NewInt(0),
		},
		{
			big.NewInt(1),
			1,
			big.NewInt(1),
			big.NewInt(1),
		},
		{
			big.NewInt(1),
			90,
			big.NewInt(80),
			big.NewInt(90),
		},
	}

	for i, test := range tests {
		value := bitvotes.CalcValue(test.feeSum, test.weight, totalWeight)

		require.Equal(t, test.uncappedValue, value.UncappedValue, fmt.Sprintf("error in test %d", i))
	}

}
