package attestation

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

var base1 = float64(1.1)
var base2 = float64(20)

type Solution struct {
	CurrentBitVoteCount []int
	SubsetDict          map[int]bool
	Value               int
	Weight              uint16
	ModValue            float64
}

func (sol *Solution) Copy() *Solution {
	newSol := &Solution{CurrentBitVoteCount: make([]int, len(sol.CurrentBitVoteCount)), SubsetDict: map[int]bool{}, Value: sol.Value, Weight: sol.Weight, ModValue: sol.ModValue}

	for i, e := range sol.CurrentBitVoteCount {
		newSol.CurrentBitVoteCount[i] = e
	}

	for key, val := range sol.SubsetDict {
		newSol.SubsetDict[key] = val
	}

	return newSol
}

func Value(currentBitVoteCount []int) int {
	val := 0

	for _, e := range currentBitVoteCount {
		if e == 0 {
			val++
		}
	}

	return val
}

func NewValue(currentBitVoteCount []int, bitVote *BitVote) int {
	val := 0

	for i, e := range currentBitVoteCount {
		if e == 0 && bitVote.BitVector.Bit(i) == 1 {
			val++
		}
	}

	return val
}

func ModValue(currentBitVoteCount []int) float64 {
	val := float64(0)

	for _, e := range currentBitVoteCount {
		val += math.Pow(base1, -float64(e))
	}

	return math.Pow(base2, val)
}

func NewModValue(currentBitVoteCount []int, bitVote *BitVote) float64 {
	val := float64(0)

	for i, e := range currentBitVoteCount {
		exponent := e
		if bitVote.BitVector.Bit(i) == 0 {
			exponent++
		}
		val += math.Pow(base1, -float64(exponent))
	}

	return math.Pow(base2, val)
}

func ModValue2(currentBitVoteCount []int) float64 {
	val := float64(0)

	for _, e := range currentBitVoteCount {
		if e == 0 {
			val += 10
		}
		val += math.Pow(base1, -float64(e))
	}

	return math.Pow(base2, val)
}

func NewModValue2(currentBitVoteCount []int, bitVote *BitVote) float64 {
	val := float64(0)

	for i, e := range currentBitVoteCount {
		exponent := e
		if bitVote.BitVector.Bit(i) == 0 {
			exponent++
		}
		if exponent == 0 {
			val += 10
		}
		val += math.Pow(base1, -float64(exponent))
	}

	return math.Pow(base2, val)
}

func ModValue3(currentBitVoteCount []int, n float64) float64 {
	val := float64(0)

	for _, e := range currentBitVoteCount {
		val -= math.Pow(base1, -float64(e))
	}

	return math.Pow(base2, val*n)
}

func NewModValue3(currentBitVoteCount []int, bitVote *BitVote, n float64) float64 {
	val := float64(0)

	for i, e := range currentBitVoteCount {
		exponent := e
		if bitVote.BitVector.Bit(i) == 0 {
			exponent++
		}
		val -= math.Pow(base1, -float64(exponent))
	}
	fmt.Println(n, math.Pow(base2, val*n))
	return math.Pow(base2, val*n)
}

func BitVotesAdd(currentSolution *Solution, bitVote *WeightedBitVote, chosen int, newModValue float64) {
	for i := range currentSolution.CurrentBitVoteCount {
		if bitVote.BitVote.BitVector.Bit(i) == 0 {
			currentSolution.CurrentBitVoteCount[i] += 1
		}
	}
	currentSolution.Weight += bitVote.Weight
	currentSolution.SubsetDict[chosen] = true
	currentSolution.Value = Value(currentSolution.CurrentBitVoteCount)
	currentSolution.ModValue = newModValue
}

func BitVotesRemove(currentSolution *Solution, bitVote *WeightedBitVote, chosen int) {
	for i := range currentSolution.CurrentBitVoteCount {
		if bitVote.BitVote.BitVector.Bit(i) == 0 {
			currentSolution.CurrentBitVoteCount[i] -= 1
		}
	}
	currentSolution.Weight -= bitVote.Weight
	currentSolution.SubsetDict[chosen] = false
	currentSolution.Value = Value(currentSolution.CurrentBitVoteCount)
	// currentSolution.ModValue = ModValue3(currentSolution.CurrentBitVoteCount, (float64(numIterations)-float64(i))/float64(numIterations))
	currentSolution.ModValue = ModValue(currentSolution.CurrentBitVoteCount)
}

func toString(m map[int]bool) string {
	s := ""
	for i := 0; i < len(m); i++ {
		if m[i] {
			s = s + strconv.Itoa(i)
		}
	}

	return s
}

func MetropolisHastingsSampling(allBitVotes []*WeightedBitVote, fees []int, numIterations int) *Solution {
	var currentSolution *Solution
	var maxSolution *Solution

	numVoters := len(allBitVotes)
	numAttestations := len(fees)
	currentBitVoteCount := make([]int, numAttestations)
	subsetMap := make(map[int]bool)
	totalWeight := uint16(0)

	currentSolution = &Solution{CurrentBitVoteCount: currentBitVoteCount, SubsetDict: subsetMap, Weight: 0}

	for i := 0; i < numVoters; i++ {
		newModValue := NewModValue(currentSolution.CurrentBitVoteCount, &allBitVotes[i].BitVote)
		BitVotesAdd(currentSolution, allBitVotes[i], i, newModValue)
		subsetMap[i] = true
		totalWeight += allBitVotes[i].Weight
	}
	fmt.Println(currentSolution)
	// fmt.Println(subsetMap)

	maxSolution = currentSolution.Copy()

	for i := 0; i < numIterations; i++ {
		for {
			chosen := rand.Intn(numVoters)
			if ok := currentSolution.SubsetDict[chosen]; ok {
				newWeight := currentSolution.Weight - allBitVotes[chosen].Weight

				if newWeight < totalWeight/2 {
					continue
				}
				BitVotesRemove(currentSolution, allBitVotes[chosen], chosen)

				if currentSolution.Value > maxSolution.Value {
					maxSolution = currentSolution.Copy()
				}
				break
			} else {
				newModValue := NewModValue(currentSolution.CurrentBitVoteCount, &allBitVotes[chosen].BitVote)

				r := rand.Float64()
				// fmt.Println("prob", newModValue/currentSolution.ModValue, newModValue, currentSolution.CurrentBitVoteCount)
				if r > newModValue/currentSolution.ModValue {
					continue
				}

				BitVotesAdd(currentSolution, allBitVotes[chosen], chosen, newModValue)

				break
			}
		}
	}

	return maxSolution
}
