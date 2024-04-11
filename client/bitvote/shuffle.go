package bitvote

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func NextSeed(seed []byte) []byte {

	return crypto.Keccak256(seed)
}

func SeedToNumberMod(seed []byte, mod uint64) uint64 {
	temp := big.NewInt(0)
	temp.SetBytes(seed)

	temp.Mod(temp, big.NewInt(int64(mod)))

	return temp.Uint64()

}

func FisherYates(n uint64, seed []byte) []uint64 {
	shuffledList := make([]uint64, n)

	for i := uint64(0); i < n; i++ {
		random := SeedToNumberMod(seed, i+1)
		seed = NextSeed(seed)
		if i != random {
			shuffledList[i] = shuffledList[random]
		}
		shuffledList[random] = i
	}

	return shuffledList

}
