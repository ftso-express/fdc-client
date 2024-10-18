package bitvotes_test

import (
	"math/big"
	"testing"

	bitvotes "github.com/flare-foundation/fdc-client/client/attestation/bitVotes"

	"github.com/stretchr/testify/require"
)

// func TestAggregateBitvotes(t *testing.T) {
// 	numAttestations := 100
// 	numVoters := 100
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)

// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote
// 		if 0.65*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{2, 3})
// 		} else {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{3, 7})
// 		}
// 		weightedBitvotes[j] = bitVote
// 	}

// 	aggregateBitVotes, aggregationMap := bitvotes.AggregateBitVotes(weightedBitvotes)

// 	require.Equal(t, len(aggregateBitVotes), 2)
// 	require.Equal(t, len(aggregationMap[0]), 65)
// 	require.Equal(t, len(aggregationMap[1]), 35)
// }

// func TestAggregateAttestations(t *testing.T) {
// 	numAttestations := 100
// 	numVoters := 100
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)

// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote
// 		if 0.65*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{2, 3})
// 		} else {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{3, 7})
// 		}
// 		weightedBitvotes[j] = bitVote
// 	}

// 	fees := make([]*big.Int, numAttestations)
// 	for j := 0; j < numAttestations; j++ {
// 		fees[j] = big.NewInt(1)
// 	}

// 	aggregatedBitVotes, aggregatedFees, _ := bitvotes.AggregateAttestations(weightedBitvotes, fees)

// 	require.Equal(t, len(aggregatedBitVotes), numAttestations)
// 	require.Equal(t, len(aggregatedFees), 4)
// 	require.Equal(t, aggregatedBitVotes[0].BitVote.Length, uint16(4))
// }

// func TestFilterBitVotes(t *testing.T) {
// 	numAttestations := 5
// 	numVoters := 100
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
// 	totalWeight := uint16(0)
// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote
// 		if 0.30*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3, 4})
// 		} else if 0.70*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{})
// 		} else {
// 			bitVote = setBitVoteFromRules(numAttestations, []int{3, 4})
// 		}
// 		weightedBitvotes[j] = bitVote
// 		totalWeight += bitVote.Weight
// 	}

// 	filtered, removedOnes, removedOnesWeight, removedZeros, removedZerosWeight := bitvotes.FilterBitVotes(weightedBitvotes)

// 	require.Equal(t, len(removedOnes), 30)
// 	require.Equal(t, removedOnesWeight, uint16(30))

// 	require.Equal(t, len(removedZeros), 40)
// 	require.Equal(t, removedZerosWeight, uint16(40))

// 	require.Equal(t, len(filtered), 30)
// }

// func TestFilterAttestations(t *testing.T) {
// 	numAttestations := 10
// 	numVoters := 100
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
// 	totalWeight := uint16(0)
// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote
// 		if 0.30*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3, 4})
// 		} else if 0.70*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{1, 4})
// 		} else {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{3, 4})
// 		}
// 		weightedBitvotes[j] = bitVote
// 		totalWeight += bitVote.Weight
// 	}
// 	fees := make([]*big.Int, numAttestations)
// 	for j := 0; j < numAttestations; j++ {
// 		fees[j] = big.NewInt(1)
// 	}

// 	filtered, _, _, removedOnes, removedLowWeight := bitvotes.FilterAttestations(weightedBitvotes, fees, totalWeight)

// 	require.Equal(t, len(removedOnes), 1)
// 	require.Equal(t, len(removedLowWeight), 7)

// 	require.Equal(t, filtered[0].BitVote.Length, uint16(2))
// }

// func TestPreProcess(t *testing.T) {
// 	numAttestations := 8
// 	numVoters := 100
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
// 	totalWeight := uint16(0)
// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote
// 		if 0.30*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3, 4})
// 		} else if 0.61*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{1, 4})
// 		} else if 0.90*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{3, 4})
// 		} else {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3, 4, 5, 6, 7})
// 		}
// 		weightedBitvotes[j] = bitVote
// 		totalWeight += bitVote.Weight
// 	}
// 	fees := make([]*big.Int, numAttestations)
// 	for j := 0; j < numAttestations; j++ {
// 		fees[j] = big.NewInt(1)
// 	}

// 	preProcessedBitVotes, newFees, preProccesInfo := bitvotes.PreProcess(weightedBitvotes, fees)
// 	fmt.Println(preProcessedBitVotes[0], preProcessedBitVotes[1])
// 	fmt.Println(preProcessedBitVotes, newFees, preProccesInfo)

// 	fmt.Printf("newFees: %v\n", newFees)

// 	require.Equal(t, 2, len(preProcessedBitVotes))
// 	require.Equal(t, 2, len(newFees))

// 	require.Equal(t, preProccesInfo.RemovedZerosWeight, uint16(0))
// }

func TestFilter(t *testing.T) {
	tests := []struct {
		vectors        []string
		weights        []uint16
		fees           []*big.Int
		totalWeight    uint16
		AlwaysInBits   []int
		AlwaysOutBits  []int
		RemainingBits  map[int]bool
		GuaranteedFees *big.Int

		AlwaysInVotes    []int
		AlwaysOutVotes   []int
		RemainingVotes   map[int]bool
		GuaranteedWeight uint16
	}{

		{
			vectors:        []string{"1111", "1100", "1001"},
			weights:        []uint16{3, 2, 1},
			fees:           []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)},
			totalWeight:    6,
			AlwaysInBits:   []int{3},
			AlwaysOutBits:  []int{1},
			RemainingBits:  map[int]bool{0: true, 2: true},
			GuaranteedFees: big.NewInt(1),

			AlwaysInVotes:    []int{0},
			AlwaysOutVotes:   []int{},
			RemainingVotes:   map[int]bool{1: true, 2: true},
			GuaranteedWeight: 3,
		},

		{
			vectors:       []string{"1111", "1100", "1001", "0000"},
			weights:       []uint16{3, 2, 1, 7},
			totalWeight:   13,
			fees:          []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)},
			AlwaysInBits:  []int{},
			AlwaysOutBits: []int{0, 1, 2, 3},
			RemainingBits: map[int]bool{},

			GuaranteedFees:   big.NewInt(0),
			AlwaysInVotes:    []int{0, 1, 2, 3},
			AlwaysOutVotes:   []int{},
			RemainingVotes:   map[int]bool{},
			GuaranteedWeight: 13,
		},
		{
			vectors:       []string{"1111", "1100", "1001", "0010"},
			weights:       []uint16{4, 4, 1, 1},
			totalWeight:   10,
			fees:          []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)},
			AlwaysInBits:  []int{3},
			AlwaysOutBits: []int{0, 1},
			RemainingBits: map[int]bool{2: true},

			GuaranteedFees:   big.NewInt(1),
			AlwaysInVotes:    []int{0, 1},
			AlwaysOutVotes:   []int{3},
			RemainingVotes:   map[int]bool{2: true},
			GuaranteedWeight: 8,
		},
	}

	for _, test := range tests {
		weightedBitVotes := make([]*bitvotes.WeightedBitVote, len(test.vectors))

		for i := range test.vectors {
			weightedVote := new(bitvotes.WeightedBitVote)

			weightedVote.Weight = test.weights[i]
			weightedVote.BitVote.BitVector, _ = new(big.Int).SetString(test.vectors[i], 2)

			weightedBitVotes[i] = weightedVote
		}

		results := bitvotes.Filter(weightedBitVotes, test.fees, test.totalWeight)

		require.ElementsMatch(t, test.AlwaysInBits, results.AlwaysInBits)
		require.ElementsMatch(t, test.AlwaysOutBits, results.AlwaysOutBits)

		require.ElementsMatch(t, test.AlwaysInVotes, results.AlwaysInVotes)
		require.ElementsMatch(t, test.AlwaysOutVotes, results.AlwaysOutVotes)

		require.Equal(t, test.RemainingBits, results.RemainingBits)
		require.Equal(t, test.RemainingVotes, results.RemainingVotes)

		require.Equal(t, test.GuaranteedFees, results.GuaranteedFees)

		require.Equal(t, test.GuaranteedWeight, results.GuaranteedWeight)

	}

}

func TestFilterAndAggregate(t *testing.T) {
	tests := []struct {
		vectors        []string
		weights        []uint16
		fees           []*big.Int
		totalWeight    uint16
		AlwaysInBits   []int
		AlwaysOutBits  []int
		RemainingBits  map[int]bool
		GuaranteedFees *big.Int

		AlwaysInVotes    []int
		AlwaysOutVotes   []int
		RemainingVotes   map[int]bool
		GuaranteedWeight uint16

		NuOfAggVotes int
		NuOfAggFees  int

		totalAggWeight uint16
		totalAggFees   *big.Int
	}{

		{
			vectors:        []string{"111111", "111100", "110111", "110100", "111000"},
			weights:        []uint16{2, 1, 1, 1, 1},
			fees:           []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)},
			totalWeight:    6,
			AlwaysInBits:   []int{4, 5},
			AlwaysOutBits:  []int{0, 1},
			RemainingBits:  map[int]bool{2: true, 3: true},
			GuaranteedFees: big.NewInt(2),

			AlwaysInVotes:    []int{0, 1},
			AlwaysOutVotes:   []int{},
			RemainingVotes:   map[int]bool{2: true, 3: true, 4: true},
			GuaranteedWeight: 3,

			NuOfAggVotes: 2,
			NuOfAggFees:  2,

			totalAggWeight: 3,
			totalAggFees:   big.NewInt(2),
		},
	}

	for _, test := range tests {
		weightedBitVotes := make([]*bitvotes.WeightedBitVote, len(test.vectors))

		for i := range test.vectors {
			weightedVote := new(bitvotes.WeightedBitVote)

			weightedVote.Weight = test.weights[i]
			weightedVote.BitVote.BitVector, _ = new(big.Int).SetString(test.vectors[i], 2)

			weightedBitVotes[i] = weightedVote
		}

		aggVotes, aggFees, results := bitvotes.FilterAndAggregate(weightedBitVotes, test.fees, test.totalWeight)

		require.ElementsMatch(t, test.AlwaysInBits, results.AlwaysInBits)
		require.ElementsMatch(t, test.AlwaysOutBits, results.AlwaysOutBits)

		require.ElementsMatch(t, test.AlwaysInVotes, results.AlwaysInVotes)
		require.ElementsMatch(t, test.AlwaysOutVotes, results.AlwaysOutVotes)

		require.Equal(t, test.RemainingBits, results.RemainingBits)
		require.Equal(t, test.RemainingVotes, results.RemainingVotes)

		require.Equal(t, test.GuaranteedFees, results.GuaranteedFees)
		require.Equal(t, test.GuaranteedWeight, results.GuaranteedWeight)

		require.Len(t, aggVotes, test.NuOfAggVotes)
		require.Len(t, aggFees, test.NuOfAggFees)

		sumWeight := uint16(0)
		for _, vote := range aggVotes {
			sumWeight += vote.Weight
		}

		require.Equal(t, test.totalAggWeight, sumWeight)

		sumFees := big.NewInt(0)
		for _, fee := range aggFees {
			sumFees.Add(sumFees, fee.Fee)
		}

		require.Equal(t, test.totalAggFees, sumFees)
	}
}

type Solution struct {
	Bits    []int
	Votes   []int
	Value   bitvotes.Value
	Optimal bool
}

func AssembleSolutionFull(filterResults *bitvotes.FilterResults, filteredSolution *bitvotes.ConsensusSolution) Solution {
	bits := []int{}
	bits = append(bits, filterResults.AlwaysInBits...)

	for k := range filteredSolution.Bits {
		indexes := filteredSolution.Bits[k].Indexes
		bits = append(bits, indexes...)
	}

	voters := []int{}
	voters = append(voters, filterResults.AlwaysInVotes...)

	for k := range filteredSolution.Votes {
		indexes := filteredSolution.Votes[k].Indexes

		voters = append(voters, indexes...)
	}

	return Solution{
		Bits:    bits,
		Votes:   voters,
		Value:   filteredSolution.Value,
		Optimal: filteredSolution.Optimal,
	}
}

func TestFilterAssemble(t *testing.T) {
	filteredResult := &bitvotes.FilterResults{AlwaysInBits: []int{1, 2, 5}}
	filteredSolution := &bitvotes.ConsensusSolution{Bits: []*bitvotes.AggregatedBit{{Indexes: []int{4, 7}}}}

	fullSolution := AssembleSolutionFull(filteredResult, filteredSolution)
	bits := bitvotes.AssembleSolution(filteredResult, filteredSolution, 10)

	require.Equal(t, 10, int(bits.Length))
	for _, e := range fullSolution.Bits {
		require.Equal(t, uint(1), bits.BitVector.Bit(e))
	}
}
