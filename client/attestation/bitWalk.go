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
	// CurrentBitVote      *BitVote
	CurrentBitVoteCount []int
	// BitVotes       []*BitVote
	// Subset     []int
	SubsetDict map[int]bool
	Value      int
	Weight     int
	ModValue   float64
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

// func (bitVote *BitVote) Value() int {

// 	fees := 0

// 	for i := 0; i < bitVote.BitVector.BitLen(); i++ {
// 		// fmt.Println(i, e)
// 		if bitVote.BitVector.Bit(i) == 1 {
// 			fees += 1
// 		}
// 	}

// 	return fees
// }

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

func BitVotesAdd(currentBitVoteCount []int, bitVote *BitVote) {
	for i := range currentBitVoteCount {
		if bitVote.BitVector.Bit(i) == 0 {
			currentBitVoteCount[i] += 1
		}
	}
}

func BitVotesRemove(currentBitVoteCount []int, bitVote *BitVote) {
	for i := range currentBitVoteCount {
		if bitVote.BitVector.Bit(i) == 0 {
			currentBitVoteCount[i] -= 1
		}
	}
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

func RandomWalk(allBitVotes []*BitVote, numIterations, numBits int) *Solution {
	setCount := make(map[string]bool)

	var currentSolution *Solution
	// var nextSolution *Solution
	var maxSolution *Solution

	// currentBitVote := allBitVotes[0]
	numVoters := len(allBitVotes)

	currentBitVoteCount := make([]int, numBits)
	// subset := make([]int, numBitvotes)
	subsetDict := make(map[int]bool)
	for i := 0; i < numVoters; i++ {
		// andBitVote := andBitwise(*currentBitVote, *allBitVotes[i])
		// currentBitVote = &andBitVote
		// subset[i] = i
		BitVotesAdd(currentBitVoteCount, allBitVotes[i])
		subsetDict[i] = true
	}
	fmt.Println(currentBitVoteCount)
	fmt.Println(subsetDict)

	currentSolution = &Solution{CurrentBitVoteCount: currentBitVoteCount, SubsetDict: subsetDict, Value: Value(currentBitVoteCount), Weight: numVoters, ModValue: ModValue(currentBitVoteCount)}
	maxSolution = currentSolution.Copy()
	count := 0
	count2 := 0
	count50 := 0

	for i := 0; i < numIterations; i++ {
		for {
			// fmt.Println(currentSolution.ModValue)
			// fmt.Println(currentSolution.CurrentBitVoteCount)

			// for j := numBitvotes - 1; j >= 0; j-- {
			// 	if currentSolution.SubsetDict[j] {
			// 		fmt.Println(" ")

			// 		fmt.Println(j)
			// 		fmt.Println(currentSolution.SubsetDict)
			// 		fmt.Println(currentSolution.CurrentBitVoteCount)
			// 		fmt.Println(currentSolution.Value)

			// 		break
			// 	}
			// }
			chosen := rand.Intn(numVoters)
			if ok := currentSolution.SubsetDict[chosen]; ok {
				if currentSolution.Weight == numVoters/2 { // todo
					continue
				}
				BitVotesRemove(currentSolution.CurrentBitVoteCount, allBitVotes[chosen])
				currentSolution.SubsetDict[chosen] = false
				currentSolution.Value = Value(currentSolution.CurrentBitVoteCount)
				// currentSolution.ModValue = ModValue3(currentSolution.CurrentBitVoteCount, (float64(numIterations)-float64(i))/float64(numIterations))
				currentSolution.ModValue = ModValue(currentSolution.CurrentBitVoteCount)
				currentSolution.Weight--

				if currentSolution.Value > maxSolution.Value {
					maxSolution = currentSolution.Copy()
				}

				setCount[toString(currentSolution.SubsetDict)] = true

				break
				// fmt.Println("uuu")

			} else {
				// fmt.Println("iii")
				// nextBitVote := andBitwise(*allBitVotes[i], *currentSolution.CurrentBitVote)
				// currentSolution.CurrentBitVote = &nextBitVote
				newValue := NewValue(currentSolution.CurrentBitVoteCount, allBitVotes[chosen])
				// newModValue := NewModValue3(currentSolution.CurrentBitVoteCount, allBitVotes[chosen], (float64(numIterations)-float64(i))/float64(numIterations))
				newModValue := NewModValue(currentSolution.CurrentBitVoteCount, allBitVotes[chosen])
				r := rand.Float64()
				// r := 0.5
				// fmt.Println("prob", newModValue/currentSolution.ModValue, newModValue, currentSolution.CurrentBitVoteCount)
				if r > newModValue/currentSolution.ModValue {
					continue
				}
				if chosen < 65 {
					count++
				} else {
					count2++
				}
				if currentSolution.Weight == numVoters/2 {
					count50++
				}

				BitVotesAdd(currentSolution.CurrentBitVoteCount, allBitVotes[chosen])
				currentSolution.SubsetDict[chosen] = true
				currentSolution.Value = newValue
				currentSolution.ModValue = newModValue
				currentSolution.Weight++

				break
			}
		}
	}

	fmt.Println("arrow", float64(count)/float64(count+count2))
	fmt.Println("50level", float64(count50)/float64(numIterations))

	fmt.Println(float64(len(setCount)) / float64(numIterations))

	return maxSolution
}
