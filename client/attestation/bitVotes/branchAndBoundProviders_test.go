package bitvotes_test

// import (
// 	"fmt"
// 	bitvotes "local/fdc/client/attestation/bitVotes"
// 	"math/big"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"
// )

// func TestBranchAndBoundProvidersFix(t *testing.T) {
// 	numAttestations := 100
// 	numVoters := 30
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)

// 	totalWeight := uint16(0)
// 	for j := 0; j < numVoters; j++ {
// 		var bitVote *bitvotes.WeightedBitVote

// 		if 0.30*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 4})
// 		} else if 0.60*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 1, 2, 3})
// 		} else if 0.90*float64(numVoters) > float64(j) {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{0, 2})
// 		} else {
// 			bitVote = setBitVoteFromPositions(numAttestations, []int{1, 3})
// 		}
// 		weightedBitvotes[j] = bitVote

// 		totalWeight += bitVote.Weight
// 	}

// 	fees := make([]*big.Int, numAttestations)
// 	for j := 0; j < numAttestations; j++ {
// 		fees[j] = big.NewInt(1)
// 	}

// 	start := time.Now()
// 	solution := bitvotes.BranchAndBoundProviders(weightedBitvotes, fees, 0, totalWeight, 50000000, time.Now().Unix())

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

// func TestBranchAndBoundProvidersRandom(t *testing.T) {
// 	numAttestations := 30
// 	numVoters := 30
// 	weightedBitvotes := make([]*bitvotes.WeightedBitVote, numVoters)
// 	prob := 0.8

// 	totalWeight := uint16(0)
// 	for j := 0; j < numVoters; j++ {
// 		bitVote := randomBitVotes(numAttestations, prob)
// 		weightedBitvotes[j] = bitVote
// 		totalWeight += bitVote.Weight
// 	}

// 	fees := make([]*big.Int, numAttestations)
// 	for j := 0; j < numAttestations; j++ {
// 		fees[j] = big.NewInt(1)
// 	}

// 	start := time.Now()
// 	solution := bitvotes.BranchAndBoundProviders(weightedBitvotes, fees, 0, totalWeight, 100000000, time.Now().Unix())

// 	fmt.Println("time passed:", time.Since(start).Seconds())
// 	fmt.Println("solution", solution)
// 	fmt.Println(solution.Value)
// 	count := 0
// 	for _, e := range solution.Solution {
// 		if e {
// 			count += 1
// 		}
// 	}
// 	count2 := 0
// 	for _, e := range solution.Participants {
// 		if e {
// 			count2 += 1
// 		}
// 	}
// 	fmt.Println("num attestations, providers", count, count2)

// 	solution2 := bitvotes.BranchAndBound(weightedBitvotes, fees, 0, totalWeight, 100000000, time.Now().Unix())
// 	fmt.Println("solution2", solution2)
// 	count = 0
// 	for _, e := range solution2.Solution {
// 		if e {
// 			count += 1
// 		}
// 	}
// 	count2 = 0
// 	for _, e := range solution2.Participants {
// 		if e {
// 			count2 += 1
// 		}
// 	}
// 	fmt.Println("2 num attestations, providers", count, count2)
// 	require.Equal(t, solution.Value, solution2.Value)
// }
