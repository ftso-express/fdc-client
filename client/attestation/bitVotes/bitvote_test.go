package bitvotes_test

import (
	"encoding/hex"
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func randomBitVotes(numAttest int, prob float64) *bitvotes.WeightedBitVote {
	weight := uint16(1)
	bitVector := big.NewInt(0)

	for j := 0; j < numAttest; j++ {
		if rand.Float64() < prob {
			bitVector.SetBit(bitVector, j, 1)
		}
	}

	return &bitvotes.WeightedBitVote{Weight: weight, BitVote: bitvotes.BitVote{uint16(numAttest), bitVector}}
}

func setBitVoteFromPositions(numAttest int, rules []int) *bitvotes.WeightedBitVote {
	weight := uint16(1)
	bitVector := big.NewInt(0)

	for j := 0; j < numAttest; j++ {
		for i := range rules {
			if j == rules[i] {
				bitVector.SetBit(bitVector, j, 1)
			}
		}
	}

	return &bitvotes.WeightedBitVote{Weight: weight, BitVote: bitvotes.BitVote{uint16(numAttest), bitVector}}
}

func setBitVoteFromRules(n int, rules []int) *bitvotes.WeightedBitVote {
	weight := uint16(1)
	bitVector := big.NewInt(0)

	for j := 0; j < n; j++ {
		for i := range rules {
			if j%rules[i] == 0 {
				bitVector.SetBit(bitVector, j, 1)
			}
		}
	}

	return &bitvotes.WeightedBitVote{Weight: weight, BitVote: bitvotes.BitVote{uint16(n), bitVector}}
}

func TestAndBitwise(t *testing.T) {

	b1, _ := new(big.Int).SetString("01100101", 2)

	b2, _ := new(big.Int).SetString("1100011", 2)

	andb, _ := new(big.Int).SetString("01100001", 2)

	bitvote1 := bitvotes.BitVote{9, b1}
	bitvote2 := bitvotes.BitVote{8, b2}

	andBitvote := bitvotes.AndBitwise(bitvote1, bitvote2)

	if andBitvote.BitVector.Cmp(andb) != 0 {
		t.Error("wrong and vector")
	}

	if andBitvote.Length != 9 {
		t.Error("wrong and length")
	}

}

func TestBitVoteForSetAllEqual(t *testing.T) {
	bitVote := setBitVoteFromRules(12, []int{2, 3, 5})

	weightedBitvotes := []*bitvotes.WeightedBitVote{}
	totalWeight := uint16(0)
	for j := 0; j < 100; j++ {
		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(j%10 + 1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := bitvotes.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	require.Equal(t, totalWeight, support)
	require.Equal(t, bitVote.BitVote.Length, bv.Length)
	require.Equal(t, bitVote.BitVote.BitVector, bv.BitVector)

}

func TestBitVoteForSetHalfMixed(t *testing.T) {
	weightedBitvotes := []*bitvotes.WeightedBitVote{}
	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		var bitVote *bitvotes.WeightedBitVote

		if j%2 == 1 {
			bitVote = setBitVoteFromRules(12, []int{2, 3})
		} else {
			bitVote = setBitVoteFromRules(12, []int{2, 5})
		}

		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(j%10 + 1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := bitvotes.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	expectedBitVote := setBitVoteFromRules(12, []int{2, 15})
	require.Equal(t, totalWeight, support)
	require.Equal(t, expectedBitVote.BitVote.Length, bv.Length)
	require.Equal(t, expectedBitVote.BitVote.BitVector, bv.BitVector)

}

func TestBitVoteForSetHalfFirst(t *testing.T) {

	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		var bitVote *bitvotes.WeightedBitVote

		if 51 > j {
			bitVote = setBitVoteFromRules(12, []int{2, 3})
		} else {
			bitVote = setBitVoteFromRules(12, []int{2, 5})
		}

		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := bitvotes.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	expectedBitVote := setBitVoteFromRules(12, []int{2, 3})

	require.Equal(t, totalWeight/2+1, support)
	require.Equal(t, expectedBitVote.BitVote.Length, bv.Length)

	require.Equal(t, expectedBitVote.BitVote.BitVector, bv.BitVector)

}

func TestBitVoteForSetZero(t *testing.T) {

	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		var bitVote *bitvotes.WeightedBitVote
		if j == 10 {
			bitVote = setBitVoteFromRules(12, []int{})
		} else {
			bitVote = setBitVoteFromRules(12, []int{2, 5})
		}

		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := bitvotes.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	expectedBitVote := bitvotes.BitVote{BitVector: big.NewInt(0), Length: 12}

	require.Equal(t, totalWeight, support)

	require.Equal(t, expectedBitVote.Length, bv.Length)

	require.Equal(t, expectedBitVote.BitVector, bv.BitVector)

}

func TestConsensusAllEqual(t *testing.T) {
	bitVote := setBitVoteFromRules(12, []int{2, 3, 5})

	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(j%10 + 1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	fees := make([]int, 12)
	for i := range fees {
		fees[i] = 1
	}

	bv, err := bitvotes.ConsensusBitVote(&bitvotes.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      totalWeight,
		Fees:             fees,
	})

	require.NoError(t, err)

	require.Equal(t, bitVote.BitVote, bv)
}

func TestConsensusNotMoreThanHalf(t *testing.T) {
	bitVote := setBitVoteFromRules(12, []int{2, 3, 5})
	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 50; j++ {
		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}
	fees := make([]int, 12)
	for i := range fees {
		fees[i] = 1
	}

	_, err := bitvotes.ConsensusBitVote(&bitvotes.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      2 * totalWeight,
		Fees:             fees,
	})

	require.Error(t, err)

}

func TestConsensusMissingAttestation(t *testing.T) {
	bitVote := setBitVoteFromRules(12, []int{2, 3, 5})

	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 50; j++ {
		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}
	fees := make([]int, 12)
	for i := range fees {
		fees[i] = 1
	}

	_, err := bitvotes.ConsensusBitVote(&bitvotes.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      totalWeight,
		Fees:             fees[:3],
	})

	t.Log("err:", err)

	require.Error(t, err)
}

func TestConsensusMixed(t *testing.T) {

	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		var bitVote *bitvotes.WeightedBitVote

		if 65 > j {
			bitVote = setBitVoteFromRules(100, []int{2, 3})
		} else {
			bitVote = setBitVoteFromRules(100, []int{2, 7})
		}

		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	fees := make([]int, 100)
	for i := range fees {
		fees[i] = 1
	}
	start := time.Now()
	bv, err := bitvotes.ConsensusBitVote(&bitvotes.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      totalWeight,
		Fees:             fees,
	})

	fmt.Printf("time passed: %v\n", time.Since(start).Seconds())

	require.NoError(t, err)

	fmt.Printf("bitVote 1: %v\n", bv.BitVector.Text(10))

	bitVote := setBitVoteFromRules(100, []int{2, 3})

	fmt.Printf("bitVote 2: %v\n", bitVote.BitVote.BitVector.Text(10))

	require.NoError(t, err)

}

func TestEncodeDecodeBitVote(t *testing.T) {
	bitVote := setBitVoteFromRules(5, []int{2, 3})

	encoded := bitVote.BitVote.EncodeBitVoteHex(257)

	require.Equal(t, "0100051d", encoded)

	byteEncoded, err := hex.DecodeString(encoded)

	require.NoError(t, err)

	decoded, roundCheck, err := bitvotes.DecodeBitVoteBytes(byteEncoded)

	require.NoError(t, err)

	require.Equal(t, bitVote.BitVote, decoded)

	require.Equal(t, uint8(1), roundCheck)

}

func TestEncodeDecodeZero(t *testing.T) {
	bitVote := setBitVoteFromRules(5, []int{})

	encoded := bitVote.BitVote.EncodeBitVoteHex(257)

	require.Equal(t, "010005", encoded)

	byteEncoded, err := hex.DecodeString(encoded)

	require.NoError(t, err)

	decoded, roundCheck, err := bitvotes.DecodeBitVoteBytes(byteEncoded)

	require.NoError(t, err)

	require.Equal(t, bitVote.BitVote, decoded)

	require.Equal(t, uint8(1), roundCheck)

}

func TestEncodeDecodeNoAttestations(t *testing.T) {
	bitVote := setBitVoteFromRules(0, []int{})
	encoded := bitVote.BitVote.EncodeBitVoteHex(257)

	require.Equal(t, "010000", encoded)

	byteEncoded, err := hex.DecodeString(encoded)

	require.NoError(t, err)

	decoded, roundCheck, err := bitvotes.DecodeBitVoteBytes(byteEncoded)

	require.NoError(t, err)

	require.Equal(t, bitVote.BitVote, decoded)

	require.Equal(t, uint8(1), roundCheck)

}

func TestDecodeFail(t *testing.T) {

	_, _, err := bitvotes.DecodeBitVoteBytes([]byte{})

	require.Error(t, err)

	byteEncoded, err := hex.DecodeString("0100")

	require.NoError(t, err)

	_, _, err = bitvotes.DecodeBitVoteBytes(byteEncoded)

	require.Error(t, err)

	byteEncoded, err = hex.DecodeString("01000000aa")

	require.NoError(t, err)

	_, _, err = bitvotes.DecodeBitVoteBytes(byteEncoded)

	require.Error(t, err)

}

func BenchmarkConsensusMixed(b *testing.B) {

	weightedBitvotes := []*bitvotes.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		var bitVote *bitvotes.WeightedBitVote

		if 65 > j {
			bitVote = setBitVoteFromRules(100, []int{2, 3})
		} else {
			bitVote = setBitVoteFromRules(100, []int{2, 7})
		}

		c := new(bitvotes.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote.BitVote

		weightedBitvotes = append(weightedBitvotes, c)
	}

	fees := make([]int, 100)
	for i := range fees {
		fees[i] = 1
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_, _ = bitvotes.ConsensusBitVote(&bitvotes.ConsensusBitVoteInput{
			RoundID:          1,
			WeightedBitVotes: weightedBitvotes,
			TotalWeight:      totalWeight,
			Fees:             fees,
		})
	}

}
