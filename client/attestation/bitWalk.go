package attestation

import (
	"math"
	"math/rand"
	"strconv"
)

// magic numbers for calculating modified value of the solution, we need to figure them out
var base1 = float64(1.2)
var base2 = float64(20)
var constantFactor = float64(100)

type RandomWalkState struct {
	CurrentBitVoteCount []int
	SubsetDict          map[int]bool
	Value               int
	Weight              uint16
	ModValue            float64
	Fees                []int
	TotalWeight         uint16
	TotalFee            int
}

func (sol *RandomWalkState) Copy() *RandomWalkState {
	newSol := &RandomWalkState{CurrentBitVoteCount: make([]int, len(sol.CurrentBitVoteCount)), SubsetDict: map[int]bool{},
		Value: sol.Value, Weight: sol.Weight, ModValue: sol.ModValue, Fees: sol.Fees, TotalWeight: sol.TotalWeight,
		TotalFee: sol.TotalFee}

	for i, e := range sol.CurrentBitVoteCount {
		newSol.CurrentBitVoteCount[i] = e
	}

	for key, val := range sol.SubsetDict {
		newSol.SubsetDict[key] = val
	}

	return newSol
}

func (sol *RandomWalkState) CalcValue() int {
	val := 0

	for i, e := range sol.CurrentBitVoteCount {
		if e == 0 {
			val += sol.Fees[i]
		}
	}

	weightCaped := min(int(float64(sol.TotalWeight)*valueCap), int(sol.Weight))

	return val * weightCaped
}

func (sol *RandomWalkState) CalcModValue() float64 {
	val := float64(0)

	for i, e := range sol.CurrentBitVoteCount {
		val += math.Pow(base1, -float64(e)) * float64(sol.Fees[i])
	}
	val *= constantFactor / float64(sol.TotalFee)

	return math.Pow(base2, val)
}

func (sol *RandomWalkState) CalcNewModValue(bitVote *BitVote) float64 {
	val := float64(0)

	for i, e := range sol.CurrentBitVoteCount {
		exponent := e
		if bitVote.BitVector.Bit(i) == 0 {
			exponent++
		}
		val += math.Pow(base1, -float64(exponent)) * float64(sol.Fees[i])
	}
	val *= constantFactor / float64(sol.TotalFee)

	return math.Pow(base2, val)
}

func (currentSolution *RandomWalkState) BitVotesAdd(bitVote *WeightedBitVote, chosen int, newModValue float64) {
	for i := range currentSolution.CurrentBitVoteCount {
		if bitVote.BitVote.BitVector.Bit(i) == 0 {
			currentSolution.CurrentBitVoteCount[i] += 1
		}
	}
	currentSolution.Weight += bitVote.Weight
	currentSolution.SubsetDict[chosen] = true
	currentSolution.Value = currentSolution.CalcValue()
	currentSolution.ModValue = newModValue
}

func (currentSolution *RandomWalkState) BitVotesRemove(bitVote *WeightedBitVote, chosen int) {
	for i := range currentSolution.CurrentBitVoteCount {
		if bitVote.BitVote.BitVector.Bit(i) == 0 {
			currentSolution.CurrentBitVoteCount[i] -= 1
		}
	}
	currentSolution.Weight -= bitVote.Weight
	currentSolution.SubsetDict[chosen] = false
	currentSolution.Value = currentSolution.CalcValue()
	currentSolution.ModValue = currentSolution.CalcModValue()
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

func MetropolisHastingsSampling(allBitVotes []*WeightedBitVote, fees []int, numIterations int) *RandomWalkState {
	var currentSolution *RandomWalkState
	var maxSolution *RandomWalkState

	numVoters := len(allBitVotes)
	numAttestations := len(fees)
	currentBitVoteCount := make([]int, numAttestations)
	subsetMap := make(map[int]bool)
	totalWeight := uint16(0)
	totalFee := 0
	for _, fee := range fees {
		totalFee += fee
	}

	currentSolution = &RandomWalkState{CurrentBitVoteCount: currentBitVoteCount, SubsetDict: subsetMap, Weight: 0, Fees: fees, TotalFee: totalFee}

	for i := 0; i < numVoters; i++ {
		newModValue := currentSolution.CalcNewModValue(&allBitVotes[i].BitVote)
		currentSolution.BitVotesAdd(allBitVotes[i], i, newModValue)
		subsetMap[i] = true
		totalWeight += allBitVotes[i].Weight
	}
	currentSolution.TotalWeight = totalWeight

	maxSolution = currentSolution.Copy()

	for i := 0; i < numIterations; i++ {
		for {
			chosen := rand.Intn(numVoters)
			if ok := currentSolution.SubsetDict[chosen]; ok {
				newWeight := currentSolution.Weight - allBitVotes[chosen].Weight

				if newWeight < totalWeight/2 {
					continue
				}
				currentSolution.BitVotesRemove(allBitVotes[chosen], chosen)

				if currentSolution.Value > maxSolution.Value {
					maxSolution = currentSolution.Copy()
				}
				break
			} else {
				newModValue := currentSolution.CalcNewModValue(&allBitVotes[chosen].BitVote)
				r := rand.Float64()
				// fmt.Println("prob", newModValue/currentSolution.ModValue)
				if r > newModValue/currentSolution.ModValue {
					continue
				}

				currentSolution.BitVotesAdd(allBitVotes[chosen], chosen, newModValue)
				break
			}
		}
	}

	return maxSolution
}
