package attestation_test

import (
	"local/fdc/client/attestation"
	"math/big"
	"testing"
)

func setAttestations() []*attestation.Attestation {

	atts := []*attestation.Attestation{}

	for j := 0; j < 10; j++ {

		att := new(attestation.Attestation)

		if j%3 == 1 {
			att.Status = attestation.Success
		} else {
			att.Status = attestation.WrongMIC

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

func TestBitVoteFromAttestations(t *testing.T) {

	atts := setAttestations()

	bitVote, err := attestation.BitVoteFromAttestations(atts)

	if err != nil {
		t.Errorf("error: %s", err)
	}

	if bitVote.Length != 10 {
		t.Errorf("wrong length")
	}

	expected, _ := big.NewInt(0).SetString("0010010010", 2)

	if bitVote.BitVector.Cmp(expected) != 0 {
		t.Error("wrong bitvector")
	}
}
