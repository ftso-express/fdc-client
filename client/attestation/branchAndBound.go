package attestation

import (
	"fmt"
	"math/big"
	"math/rand"
)

type CurrentStatus struct {
	CurrentBound  int
	NumOperations int
	MaxOperations int
	TotalWeight   uint16
	BitVotes      []*WeightedBitVote
	Fees          []int
	RandGen       rand.Source
	// mu           sync.Mutex
}

type BranchAndBoundPartialSolution struct {
	Participants    map[int]bool
	SolutionReverse []bool
	Value           int
}

type BranchAndBoundSolution struct {
	Participants []bool
	Solution     []bool
	Value        int
}

func CalcValue(feeSum int, weight, totalWeight uint16) int {

	weightCaped := min(int(float64(totalWeight)*valueCap), int(weight))

	return feeSum * weightCaped
}

func RandPerm(n int, randGen rand.Source) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := randGen.Int63() % int64(i+1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

func PermuteBits(allBitVotes []*WeightedBitVote, randPerm []int) []*WeightedBitVote {
	permBitVotes := make([]*WeightedBitVote, len(allBitVotes))
	for i, e := range allBitVotes {
		permBitVotes[i] = &WeightedBitVote{Weight: e.Weight, BitVote: BitVote{Length: e.BitVote.Length, BitVector: big.NewInt(0)}}
		for key, val := range randPerm {
			if e.BitVote.BitVector.Bit(key) == 1 {
				permBitVotes[i].BitVote.BitVector.Add(permBitVotes[i].BitVote.BitVector, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(val)), nil))
			}
		}
	}

	return permBitVotes
}

func BranchAndBound(allBitVotes []*WeightedBitVote, fees []int, numIterations int, seed int64) *BranchAndBoundSolution {
	numAttestations := len(fees)
	numVoters := len(allBitVotes)
	totalWeight := uint16(0)

	participants := make(map[int]bool)
	for i, vote := range allBitVotes {
		totalWeight += vote.Weight
		participants[i] = true
	}

	totalFee := 0
	for _, fee := range fees {
		totalFee += fee
	}
	randGen := rand.NewSource(seed)
	randPerm := RandPerm(numAttestations, randGen)
	permBitVotes := PermuteBits(allBitVotes, randPerm)

	currentBound := &CurrentStatus{CurrentBound: 0, NumOperations: 0, MaxOperations: numIterations, TotalWeight: totalWeight, BitVotes: permBitVotes, Fees: fees, RandGen: randGen}

	permResult := Branch(participants, totalFee, currentBound, 0, numAttestations, totalWeight)

	result := BranchAndBoundSolution{Participants: make([]bool, numVoters), Solution: make([]bool, numAttestations), Value: permResult.Value}
	for key, val := range permResult.Participants {
		result.Participants[key] = val
	}
	for key, val := range randPerm {
		result.Solution[key] = permResult.SolutionReverse[numAttestations-val-1]
	}

	return &result
}

func Branch(participants map[int]bool, feeSum int, currentStatus *CurrentStatus, branch, numAttestations int, currentWeight uint16) *BranchAndBoundPartialSolution {
	currentStatus.NumOperations++

	// end of recursion
	if branch == numAttestations {
		value := CalcValue(feeSum, currentWeight, currentStatus.TotalWeight)
		if value > currentStatus.CurrentBound {
			currentStatus.CurrentBound = value
		}

		return &BranchAndBoundPartialSolution{Participants: participants, SolutionReverse: []bool{}, Value: value}
	}

	// check if we already reached the maximal search space or if we exceeded the bound of the maximal possible value of the solution
	if currentStatus.NumOperations >= currentStatus.MaxOperations || CalcValue(feeSum, currentWeight, currentStatus.TotalWeight) < currentStatus.CurrentBound {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {
		result0 = Branch(participants, feeSum-currentStatus.Fees[branch], currentStatus, branch+1, numAttestations, currentWeight)
	}

	// prepare and check if a branch is possible
	newParticipants := make(map[int]bool)
	newCurrentWeight := currentWeight
	for participant := range participants {
		if currentStatus.BitVotes[participant].BitVote.BitVector.Bit(branch) == 1 {
			newParticipants[participant] = true
		} else {
			newCurrentWeight -= currentStatus.BitVotes[participant].Weight
		}
		currentStatus.NumOperations++
	}
	if newCurrentWeight > currentStatus.TotalWeight/2 {
		result1 = Branch(newParticipants, feeSum, currentStatus, branch+1, numAttestations, newCurrentWeight)
	}

	if randBit == 1 {
		result0 = Branch(participants, feeSum-currentStatus.Fees[branch], currentStatus, branch+1, numAttestations, currentWeight)
	}

	if branch == 0 {
		fmt.Println(currentStatus.NumOperations, randBit)
	}
	// max result
	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		result0.SolutionReverse = append(result0.SolutionReverse, false)
		return result0
	} else if result0 == nil || result0.Value < result1.Value {
		result1.SolutionReverse = append(result1.SolutionReverse, true)
		return result1
	} else {
		result0.SolutionReverse = append(result0.SolutionReverse, false)
		return result0
	}
}
