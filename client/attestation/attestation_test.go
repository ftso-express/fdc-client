package attestation_test

import (
	"local/fdc/client/attestation"
	"math/big"
	"testing"
)

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
