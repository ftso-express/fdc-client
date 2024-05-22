package shuffle

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// Seed returns the initial seed for consecutive sample in roundId
func Seed(roundId, sample int64) []byte {
	r := big.NewInt(roundId)
	s := big.NewInt(roundId)

	return crypto.Keccak256(r.Bytes(), s.Bytes())

}

// nextSeed returns the next seed in sequence for computation in Fisher Yates shuffle
func nextSeed(seed []byte) []byte {

	return crypto.Keccak256(seed)
}

// seedToNumberMod considers seed as a (hex) number (big endian) computes it mod.
// It is assumed that mod is a positive number.
func seedToNumberMod(seed []byte, mod uint64) uint64 {
	temp := big.NewInt(0)
	temp.SetBytes(seed)

	temp.Mod(temp, big.NewInt(int64(mod)))

	return temp.Uint64()

}

// FisherYates returns a shuffled list of first n nonnegative integers
func FisherYates(n uint64, seed []byte) []uint64 {
	shuffledList := make([]uint64, n)

	for i := uint64(0); i < n; i++ {
		random := seedToNumberMod(seed, i+1) // i+1 > 0
		seed = nextSeed(seed)
		if i != random {
			shuffledList[i] = shuffledList[random]
		}
		shuffledList[random] = i
	}

	return shuffledList

}
