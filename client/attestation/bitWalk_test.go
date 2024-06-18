package attestation_test

import (
	"fmt"
	"local/fdc/client/attestation"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBitWalk(t *testing.T) {
	numAttestations := 100
	numVoters := 100
	weightedBitvotes := []*attestation.WeightedBitVote{}

	for j := 0; j < numVoters; j++ {
		var atts []*attestation.Attestation

		if 0.65*float64(numVoters) > float64(j) {
			atts = setAttestations(numAttestations, []int{2, 3})
		} else {
			atts = setAttestations(numAttestations, []int{2, 7})
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
		fees[j] = 10
	}

	start := time.Now()
	solution := attestation.MetropolisHastingsSampling(weightedBitvotes, fees, 100000)

	fmt.Println("time passed:", time.Since(start).Seconds())
	fmt.Println("solution", solution)
	fmt.Println(solution.Value)
	// fmt.Println(solution.CurrentBitVote.BitVector.Text(2))

	// require.NoError(t, err)

	// fmt.Printf("bitVote 1: %v\n", bv.BitVector.Text(10))

	// bitVote, err := attestation.BitVoteFromAttestations(atts)

	// fmt.Printf("bitVote 2: %v\n", bitVote.BitVector.Text(10))

	// require.NoError(t, err)
}
