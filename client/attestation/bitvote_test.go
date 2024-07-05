package attestation_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"local/fdc/client/attestation"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setAttestations(n int, rules []int) []*attestation.Attestation {

	atts := []*attestation.Attestation{}

	for j := 0; j < n; j++ {

		att := new(attestation.Attestation)

		att.Fee = big.NewInt(10)

		att.Status = attestation.ProcessError

		index := attestation.IndexLog{uint64(j), uint64(j % 2)}

		att.Indexes = append(att.Indexes, index)

		for i := range rules {

			if j%rules[i] == 0 {
				att.Status = attestation.Success
			}
		}

		atts = append(atts, att)
	}

	return atts
}

func TestAndBitwise(t *testing.T) {

	b1, _ := new(big.Int).SetString("01100101", 2)

	b2, _ := new(big.Int).SetString("1100011", 2)

	andb, _ := new(big.Int).SetString("01100001", 2)

	bitvote1 := attestation.BitVote{9, b1}
	bitvote2 := attestation.BitVote{8, b2}

	andBitvote := attestation.AndBitwise(bitvote1, bitvote2)

	if andBitvote.BitVector.Cmp(andb) != 0 {
		t.Error("wrong and vector")
	}

	if andBitvote.Length != 9 {
		t.Error("wrong and length")
	}

}

func TestBitVoteFromAttestationsEmpty(t *testing.T) {

	bitVote, err := attestation.BitVoteFromAttestations([]*attestation.Attestation{})

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if bitVote.Length != 0 {
		t.Errorf("wrong length")
	}

	expected := big.NewInt(0)

	if bitVote.BitVector.Cmp(expected) != 0 {
		t.Error("wrong bitvector")
	}

}

func TestBitVoteFromAttestations(t *testing.T) {

	atts := setAttestations(10, []int{3})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if bitVote.Length != 10 {
		t.Errorf("wrong length")
	}

	expected, _ := big.NewInt(0).SetString("1001001001", 2)

	if bitVote.BitVector.Cmp(expected) != 0 {
		t.Error("wrong bitvector")
	}

}

func TestBitVoteForSetAllEqual(t *testing.T) {

	atts := setAttestations(12, []int{2, 3, 5})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(j%10 + 1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := attestation.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	require.Equal(t, totalWeight, support)

	require.Equal(t, bitVote.Length, bv.Length)

	require.Equal(t, bitVote.BitVector, bv.BitVector)

}

func TestBitVoteForSetHalfMixed(t *testing.T) {

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {

		var atts []*attestation.Attestation

		if j%2 == 1 {

			atts = setAttestations(12, []int{2, 3})
		} else {

			atts = setAttestations(12, []int{2, 5})
		}

		bitVote, err := attestation.BitVoteFromAttestations(atts)

		require.NoError(t, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(j%10 + 1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := attestation.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	atts := setAttestations(12, []int{2, 15})

	expectedBitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	require.Equal(t, totalWeight, support)

	require.Equal(t, expectedBitVote.Length, bv.Length)

	require.Equal(t, expectedBitVote.BitVector, bv.BitVector)

}

func TestBitVoteForSetHalfFirst(t *testing.T) {

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {

		var atts []*attestation.Attestation

		if 51 > j {

			atts = setAttestations(12, []int{2, 3})
		} else {

			atts = setAttestations(12, []int{2, 5})
		}

		bitVote, err := attestation.BitVoteFromAttestations(atts)

		require.NoError(t, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := attestation.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	atts := setAttestations(12, []int{2, 3})

	expectedBitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	require.Equal(t, totalWeight/2+1, support)

	require.Equal(t, expectedBitVote.Length, bv.Length)

	require.Equal(t, expectedBitVote.BitVector, bv.BitVector)

}

func TestBitVoteForSetZero(t *testing.T) {

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {

		var atts []*attestation.Attestation

		if j == 10 {

			atts = setAttestations(12, []int{})
		} else {

			atts = setAttestations(12, []int{2, 5})
		}

		bitVote, err := attestation.BitVoteFromAttestations(atts)

		require.NoError(t, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	shuffle := make([]uint64, 100)

	for j := range shuffle {
		shuffle[j] = uint64(j)
	}

	bv, support := attestation.BitVoteForSet(weightedBitvotes, totalWeight, shuffle)

	expectedBitVote := attestation.BitVote{BitVector: big.NewInt(0), Length: 12}

	require.Equal(t, totalWeight, support)

	require.Equal(t, expectedBitVote.Length, bv.Length)

	require.Equal(t, expectedBitVote.BitVector, bv.BitVector)

}

func TestConsensusAllEqual(t *testing.T) {

	atts := setAttestations(12, []int{2, 3, 5})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {
		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(j%10 + 1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	bv, err := attestation.ConsensusBitVote(&attestation.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      totalWeight,
		Attestations:     atts,
	},
		300)

	require.NoError(t, err)

	require.Equal(t, bitVote, bv)
}

func TestConsensusNotMoreThanHalf(t *testing.T) {

	atts := setAttestations(12, []int{2, 3, 5})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 50; j++ {
		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	_, err = attestation.ConsensusBitVote(&attestation.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      2 * totalWeight,
		Attestations:     atts,
	},
		300,
	)

	require.Error(t, err)

}

func TestConsensusMissingAttestation(t *testing.T) {

	atts := setAttestations(12, []int{2, 3, 5})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 50; j++ {
		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	_, err = attestation.ConsensusBitVote(&attestation.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      totalWeight,
		Attestations:     atts[:3],
	},
		300,
	)

	t.Log("err:", err)

	require.Error(t, err)

}

func createWeightedBitVotes(honestWeightPct, noOfAttestations int) (*attestation.ConsensusBitVoteInput, error) {

	if honestWeightPct > 100 {
		return nil, errors.New("honestWeightPct to high")
	}

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	attsHonest := setAttestations(noOfAttestations, []int{2, 3})
	for i := 0; i < honestWeightPct; i++ {

		bitVote, err := attestation.BitVoteFromAttestations(attsHonest)

		if err != nil {
			return nil, err

		}

		c := new(attestation.WeightedBitVote)
		c.Index = i
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	for i := honestWeightPct; i < 100; i++ {

		attsFaulty := setAttestations(noOfAttestations, []int{2, i})

		bitVote, err := attestation.BitVoteFromAttestations(attsFaulty)

		if err != nil {
			return nil, err

		}

		c := new(attestation.WeightedBitVote)
		c.Index = i
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	return &attestation.ConsensusBitVoteInput{
		RoundID:          1,
		WeightedBitVotes: weightedBitvotes,
		TotalWeight:      totalWeight,
		Attestations:     attsHonest,
	}, nil

}

func TestConsensusMixed(t *testing.T) {

	data, err := createWeightedBitVotes(75, 100)

	require.NoError(t, err)

	start := time.Now()
	bv, err := attestation.ConsensusBitVote(data, 300)

	fmt.Printf("time passed: %v\n", time.Since(start).Seconds())

	require.NoError(t, err)

	bitVote, err := attestation.BitVoteFromAttestations(data.Attestations)

	require.NoError(t, err)

	require.Equal(t, bitVote, bv)

}

func TestEncodeDecodeBitVote(t *testing.T) {

	atts := setAttestations(5, []int{2, 3})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	encoded := bitVote.EncodeBitVoteHex(257)

	require.Equal(t, "0100051d", encoded)

	byteEncoded, err := hex.DecodeString(encoded)

	require.NoError(t, err)

	decoded, roundCheck, err := attestation.DecodeBitVoteBytes(byteEncoded)

	require.NoError(t, err)

	require.Equal(t, bitVote, decoded)

	require.Equal(t, uint8(1), roundCheck)

}

func TestEncodeDecodeZero(t *testing.T) {

	atts := setAttestations(5, []int{})

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	require.NoError(t, err)

	encoded := bitVote.EncodeBitVoteHex(257)

	require.Equal(t, "010005", encoded)

	byteEncoded, err := hex.DecodeString(encoded)

	require.NoError(t, err)

	decoded, roundCheck, err := attestation.DecodeBitVoteBytes(byteEncoded)

	require.NoError(t, err)

	require.Equal(t, bitVote, decoded)

	require.Equal(t, uint8(1), roundCheck)

}

func TestEncodeDecodeNoAttestations(t *testing.T) {

	bitVote, err := attestation.BitVoteFromAttestations([]*attestation.Attestation{})

	require.NoError(t, err)

	encoded := bitVote.EncodeBitVoteHex(257)

	require.Equal(t, "010000", encoded)

	byteEncoded, err := hex.DecodeString(encoded)

	require.NoError(t, err)

	decoded, roundCheck, err := attestation.DecodeBitVoteBytes(byteEncoded)

	require.NoError(t, err)

	require.Equal(t, bitVote, decoded)

	require.Equal(t, uint8(1), roundCheck)

}

func TestDecodeFail(t *testing.T) {

	_, _, err := attestation.DecodeBitVoteBytes([]byte{})

	require.Error(t, err)

	byteEncoded, err := hex.DecodeString("0100")

	require.NoError(t, err)

	_, _, err = attestation.DecodeBitVoteBytes(byteEncoded)

	require.Error(t, err)

	byteEncoded, err = hex.DecodeString("01000000aa")

	require.NoError(t, err)

	_, _, err = attestation.DecodeBitVoteBytes(byteEncoded)

	require.Error(t, err)

}

func BenchmarkConsensusMixed(b *testing.B) {

	weightedBitvotes := []*attestation.WeightedBitVote{}

	totalWeight := uint16(0)

	for j := 0; j < 100; j++ {

		var atts []*attestation.Attestation

		if 65 > j {

			atts = setAttestations(20, []int{2, 3})
		} else {

			atts = setAttestations(20, []int{2, 7})
		}

		bitVote, err := attestation.BitVoteFromAttestations(atts)

		require.NoError(b, err)

		c := new(attestation.WeightedBitVote)
		c.Index = j
		c.Weight = uint16(1)
		totalWeight += c.Weight
		c.BitVote = bitVote

		weightedBitvotes = append(weightedBitvotes, c)

	}

	atts := setAttestations(100, []int{2, 3})

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		_, _ = attestation.ConsensusBitVote(&attestation.ConsensusBitVoteInput{
			RoundID:          1,
			WeightedBitVotes: weightedBitvotes,
			TotalWeight:      totalWeight,
			Attestations:     atts,
		},
			300,
		)
	}

}
