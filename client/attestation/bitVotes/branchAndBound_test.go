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
	aggBitVotes := make([]*bitvotes.AggregatedVote, numVoters)
	prob := 0.8

	totalWeight := uint16(0)
	for j := 0; j < numVoters; j++ {
		bitVote := randomBitVoteAggregated(numAttestations, prob, j)
		aggBitVotes[j] = bitVote

		totalWeight += bitVote.Weight
	}

	fees := make([]*bitvotes.AggregatedFee, numAttestations)
	for j := 0; j < numAttestations; j++ {
		fee := bitvotes.AggregatedFee{Fee: big.NewInt(1), Indexes: []int{j}}

		fees[j] = &fee
	}

	initialBound := bitvotes.Value{big.NewInt(0), big.NewInt(0)}

	start := time.Now()
	solution := bitvotes.BranchAndBoundBits(
		aggBitVotes,
		fees,
		0,
		totalWeight,
		big.NewInt(0),
		100000000,
		initialBound,
		true,
	)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println("value", solution.Value)

	fmt.Println("num attestations", len(solution.Bits))
}

// func TestBranchAndBound65(t *testing.T) {
// 	numAttestations := 100
// 	numVoters := 100
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
// 	totalWeight := uint16(0)

// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote
// 		if 0.65*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{2, 3})
// 		} else {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{3, 7})
// 		}
// 		weightedBitvotes[j] = bitVote

// 		totalWeight += bitVote.Weight
// 	}

// 	fees := make([]*big.Int, numAttestations)
// 	for j := 0; j < numAttestations; j++ {
// 		fees[j] = big.NewInt(1)
// 	}

// 	start := time.Now()

// 	solution := bitvotes.BranchAndBound(weightedBitvotes, fees, 0, totalWeight, 50000000, time.Now().Unix())

// 	fmt.Println("time passed:", time.Since(start).Seconds())
// 	fmt.Println("solution", solution)
// 	fmt.Println(solution.Value)
// 	count := 0
// 	for _, e := range solution.Solution {
// 		if e {
// 			count += 1
// 		}
// 	}

// 	fmt.Println("num attestations", count)
// }

func TestBranchAndBoundFix(t *testing.T) {
	numAttestations := 5
	numVoters := 100
	weightedBitvotes := make([]*bitvotes.AggregatedVote, numVoters)
	totalWeight := uint16(0)

	for j := 0; j < numVoters; j++ {
		var bitVote *bitvotes.AggregatedVote

		if 0.30*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositionAgg(numAttestations, []int{0, 1, 2, 4}, j)
		} else if 0.8*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositionAgg(numAttestations, []int{0, 1, 2, 3}, j)
		} else if 0.90*float64(numVoters) > float64(j) {
			bitVote = setBitVoteFromPositionAgg(numAttestations, []int{0, 2}, j)
		} else {
			bitVote = setBitVoteFromPositionAgg(numAttestations, []int{1, 3}, j)
		}
		weightedBitvotes[j] = bitVote

		totalWeight += bitVote.Weight
	}

	fees := make([]*bitvotes.AggregatedFee, numAttestations)
	for j := 0; j < numAttestations; j++ {

		fee := bitvotes.AggregatedFee{Fee: big.NewInt(1), Indexes: []int{j}, Support: 1}

		fees[j] = &fee
	}

	initialBound := bitvotes.Value{big.NewInt(0), big.NewInt(0)}

	start := time.Now()
	solution := bitvotes.BranchAndBoundBits(weightedBitvotes, fees, 0, totalWeight, big.NewInt(0), 50000000, initialBound, true)

	fmt.Println("time passed:", time.Since(start).Seconds())

	require.Equal(t, 3, len(solution.Bits), "not enough bits in solution")

	for _, j := range []int{0, 1, 2} {
		require.Contains(t, solution.Bits, fees[j], "wrong bits in solution")
	}

	require.Equal(t, bitvotes.Value{big.NewInt(240), big.NewInt(240)}, solution.Value)

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
			big.NewInt(90),
			big.NewInt(80),
		},
	}

	for i, test := range tests {
		value := bitvotes.CalcValue(test.feeSum, test.weight, totalWeight)

		require.Equal(t, test.cappedValue, value.CappedValue, fmt.Sprintf("error in test %d", i))

		require.Equal(t, test.uncappedValue, value.UncappedValue, fmt.Sprintf("error in test %d", i))
	}

}

func TestSort(t *testing.T) {

	fee0 := bitvotes.AggregatedFee{
		Fee:     big.NewInt(1),
		Indexes: []int{0},
		Support: 10,
	}

	fee1 := bitvotes.AggregatedFee{
		Fee:     big.NewInt(3),
		Indexes: []int{1},
		Support: 5,
	}

	fee2 := bitvotes.AggregatedFee{
		Fee:     big.NewInt(1),
		Indexes: []int{2},
		Support: 10,
	}

	fee3 := bitvotes.AggregatedFee{
		Fee:     big.NewInt(1),
		Indexes: []int{2},
		Support: 8,
	}
	tests := []struct {
		totalWeight uint16
		fees        []*bitvotes.AggregatedFee
		asc         []*bitvotes.AggregatedFee
		dsc         []*bitvotes.AggregatedFee
	}{
		{11,
			[]*bitvotes.AggregatedFee{&fee0, &fee1, &fee2, &fee3},
			[]*bitvotes.AggregatedFee{&fee3, &fee0, &fee2, &fee1},
			[]*bitvotes.AggregatedFee{&fee1, &fee0, &fee2, &fee3},
		},
	}

	for _, test := range tests {

		asc := bitvotes.SortFees(test.fees, bitvotes.CmpValAsc(test.totalWeight))

		dsc := bitvotes.SortFees(test.fees, bitvotes.CmpValDsc(test.totalWeight))

		require.Equal(t, test.asc, asc)

		require.Equal(t, test.dsc, dsc)

	}

}
