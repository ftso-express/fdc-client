package shuffle

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var uint64Ty, _ = abi.NewType("uint64", "uint64", nil)

var args = abi.Arguments{{Type: uint64Ty}, {Type: uint64Ty}, {Type: uint64Ty}}

// Seed returns the initial seed for consecutive sample in roundId.
// The Seed in Solidity is computed by keccak256(abi.encode(roundId,sample));, where roundId and sample are of type uint64.
func Seed(roundId, sample, protocolId uint64) common.Hash {

	encoded, _ := args.Pack(roundId, sample, protocolId)

	return crypto.Keccak256Hash(encoded)
}

// nextSeed returns the next seed in sequence for computation in Fisher Yates shuffle
func nextSeed(seed common.Hash) common.Hash {

	return crypto.Keccak256Hash(seed.Bytes())
}

// seedToNumberMod considers seed as a (hex) number (big endian) and computes it mod.
// It is assumed that mod is a positive number.
func seedToNumberMod(seed common.Hash, mod uint64) uint64 {
	num := seed.Big()

	num.Mod(num, big.NewInt(int64(mod)))

	return num.Uint64()
}

// FisherYates returns a shuffled list of first n nonnegative integers
func FisherYates(n uint64, seed common.Hash) []uint64 {
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
