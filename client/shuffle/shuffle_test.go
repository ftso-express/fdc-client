package shuffle_test

import (
	"local/fdc/client/shuffle"
	"slices"
	"testing"
)

func TestSeed(t *testing.T) {

	seed1 := shuffle.Seed(1, 1, 300)

	seed2 := shuffle.Seed(0, 0, 300)

	if seed1 == seed2 {
		t.Error("not random seeds")
	}

}

func TestFisherYates(t *testing.T) {

	seed1 := shuffle.Seed(1, 1, 300)

	shuffle1 := shuffle.FisherYates(0, seed1)

	if len(shuffle1) != 0 {
		t.Error("non empty array")
	}

	shuffle2 := shuffle.FisherYates(100, seed1)

	if len(shuffle2) != 100 {
		t.Error("wrong length array")
	}

	for j := 0; j < 100; j++ {

		ok := slices.Contains[[]uint64, uint64](shuffle2, uint64(j))

		if !ok {
			t.Errorf("missing %d", j)
		}
	}

	seed2 := shuffle.Seed(1, 0, 300)

	shuffle3 := shuffle.FisherYates(100, seed2)

	equal := true
	for j := range shuffle3 {
		if shuffle3[j] != shuffle2[j] {
			equal = false
			break
		}
	}

	if equal {
		t.Error("non random")
	}

}

func BenchmarkFisherYates(b *testing.B) {

	for i := 0; i < b.N; i++ {

		seed := shuffle.Seed(1, uint64(i), 300)

		shuffle.FisherYates(100, seed)

	}

}
