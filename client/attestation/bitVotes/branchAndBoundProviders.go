package bitvotes

import (
	"math/big"
	"slices"
)

func PermuteBitVotes(bitVotes []*WeightedBitVote, randPerm []int) []*WeightedBitVote {
	permBitVotes := make([]*WeightedBitVote, len(bitVotes))
	for i, e := range bitVotes {
		permBitVotes[randPerm[i]] = e
	}

	return permBitVotes
}

func CalcValueVote(feeSum *big.Int, weight uint16) *big.Int {

	return new(big.Int).Mul(feeSum, big.NewInt(int64(weight)))

}

func cmpValVote(sign int) func(*AggregatedVote, *AggregatedVote) int {

	return func(vote0, vote1 *AggregatedVote) int {

		val0 := CalcValueVote(vote0.Fees, vote0.Weight)

		val1 := CalcValueVote(vote1.Fees, vote1.Weight)

		cmp := val0.Cmp(val1)

		if cmp != 0 {
			return sign * cmp
		}

		if vote0.Indexes[0] < vote1.Indexes[0] {
			return -1
		}
		if vote0.Indexes[0] > vote1.Indexes[0] {
			return 1
		}

		return 0

	}
}

func cmpValVoteAsc() func(*AggregatedVote, *AggregatedVote) int {
	return cmpValVote(1)
}
func cmpValVoteDsc() func(*AggregatedVote, *AggregatedVote) int {
	return cmpValVote(-1)
}

func sortVotes(votes []*AggregatedVote, sortFunc func(*AggregatedVote, *AggregatedVote) int) []*AggregatedVote {

	sortedFees := make([]*AggregatedVote, len(votes))
	copy(sortedFees, votes)
	slices.SortStableFunc(sortedFees, sortFunc)

	return sortedFees
}

// BranchAndBoundProvidersDouble runs two branch and bound strategies on votes in parallel and returns the better result.
//
// The first strategy sorts the aggregated votes by the descending value (weight * fee) and at depth k the branch in which k-th vote is included is explored first.
// The second strategy sorts the aggregated votes by the ascending value (weight * fee) and at depth k the branch in which does not include k-th vote is explored first.
//
// If both strategies find an optimal but different solutions, the solution of the first strategy is returned.
func BranchAndBoundProvidersDouble(bitVotes []*AggregatedVote, fees []*AggregatedFee, assumedWeight, absoluteTotalWeight uint16, assumedFees *big.Int, maxOperations int, initialBound Value) *ConsensusSolution {

	solutions := make([]*ConsensusSolution, 2)

	firstDone := make(chan bool, 1)
	secondDone := make(chan bool, 1)

	votesAscVal := sortVotes(bitVotes, cmpValVoteAsc())

	votesDscVal := sortVotes(bitVotes, cmpValVoteDsc())

	go func() {

		solution := BranchAndBoundBits(votesDscVal, fees, assumedWeight, absoluteTotalWeight, assumedFees, maxOperations, initialBound, func(...interface{}) bool { return true })

		solutions[0] = solution

		firstDone <- true

		if solution.Optimal {
			solutions[1] = nil
			secondDone <- true // do not wait on the other solution
		}

	}()

	go func() {

		solution := BranchAndBoundBits(votesAscVal, fees, assumedWeight, absoluteTotalWeight, assumedFees, maxOperations, initialBound, func(...interface{}) bool { return false })

		solutions[1] = solution

		secondDone <- true

	}()

	<-firstDone
	<-secondDone

	if solutions[0] == nil {
		return solutions[1]
	}
	if solutions[1] == nil {
		return solutions[0]
	}

	if solutions[0].Value.Cmp(solutions[1].Value) == -1 {

		return solutions[1]
	}

	return solutions[0]

}

// BranchAndBoundVotes is similar to BranchAndBound, the difference is that it
// executes a branch and bound strategy on the space of subsets of bitVotes, hence
// it is particularly useful when there are not too many distinct bitVotes.
func BranchAndBoundVotes(bitVotes []*AggregatedVote, fees []*AggregatedFee, assumedWeight, absoluteTotalWeight uint16, assumedFees *big.Int, maxOperations int, initialBound Value) *ConsensusSolution {
	totalWeight := assumedWeight

	for _, vote := range bitVotes {
		totalWeight += vote.Weight
	}

	totalFee := big.NewInt(0).Set(assumedFees)
	bits := make(map[int]bool)
	for i, fee := range fees {
		totalFee.Add(totalFee, fee.Fee)
		bits[i] = true
	}

	currentStatus := &SharedStatus{
		CurrentBound:  initialBound,
		NumOperations: 0,
	}

	processInfo := &ProcessInfo{
		TotalWeight:      absoluteTotalWeight,
		LowerBoundWeight: absoluteTotalWeight / 2,
		BitVotes:         bitVotes,
		Fees:             fees,
		NumAttestations:  len(fees),
		NumProviders:     len(bitVotes),
		MaxOperations:    maxOperations,
	}

	permResult := BranchVotes(processInfo, currentStatus, 0, bits, totalFee, totalWeight)

	isOptimal := currentStatus.NumOperations < maxOperations

	// empty solution
	if permResult == nil {

		return &ConsensusSolution{
			Votes:   bitVotes,
			Bits:    []*AggregatedFee{},
			Value:   Value{big.NewInt(0), big.NewInt(0)},
			Optimal: isOptimal,
		}

	}

	result := ConsensusSolution{
		Votes:   make([]*AggregatedVote, 0),
		Bits:    make([]*AggregatedFee, 0),
		Optimal: isOptimal,
	}

	if !isOptimal {
		permResult.MaximizeVotes(bitVotes, fees, assumedFees, assumedWeight, absoluteTotalWeight)
	}

	result.Value = permResult.Value

	for key := range permResult.Votes {
		result.Votes = append(result.Votes, bitVotes[key])
	}
	for key := range permResult.Bits {
		result.Bits = append(result.Bits, fees[key])
	}

	return &result
}

func BranchVotes(processInfo *ProcessInfo, currentStatus *SharedStatus, branch int, bits map[int]bool, feeSum *big.Int, weight uint16) *BranchAndBoundPartialSolution {
	currentStatus.NumOperations++

	// end of recursion
	if branch == processInfo.NumProviders {
		value := CalcValue(feeSum, weight, processInfo.TotalWeight)

		if value.Cmp(currentStatus.CurrentBound) == 1 {

			currentStatus.CurrentBound = value

			return &BranchAndBoundPartialSolution{
				Votes: make(map[int]bool),
				Bits:  bits,
				Value: value,
			}

		}

		return nil

	}

	// check if we already reached the maximal search space
	if currentStatus.NumOperations >= processInfo.MaxOperations {
		return nil
	}

	// check if the potential maximal value of a branch in not higher than the current bound
	if CalcValue(feeSum, weight, processInfo.TotalWeight).Cmp(currentStatus.CurrentBound) != 1 {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	newWeight := weight - processInfo.BitVotes[branch].Weight

	// decide randomly which branch is first
	randBit := currentStatus.RandGen.Int63() % 2
	if randBit == 0 {
		// check if a branch is possible
		if newWeight > processInfo.LowerBoundWeight {
			result0 = BranchVotes(processInfo, currentStatus, branch+1, bits, feeSum, newWeight)
		}
	}

	// prepare a new branch
	newBits, newFeeSum := prepareDataForBranchWithVote(processInfo, currentStatus, bits, feeSum, processInfo.BitVotes[branch].Indexes[0])

	result1 = BranchVotes(processInfo, currentStatus, branch+1, newBits, newFeeSum, weight)

	if randBit == 1 {
		if newWeight > processInfo.LowerBoundWeight {
			result0 = BranchVotes(processInfo, currentStatus, branch+1, bits, feeSum, newWeight)
		}
	}

	// max result
	return joinResultsVotes(result0, result1, branch)
}

// prepareDataForBranchWithVote prepares data for branch in which the vote on voteIndex is included.
func prepareDataForBranchWithVote(processInfo *ProcessInfo, currentStatus *SharedStatus, bits map[int]bool, feeSum *big.Int, voteIndex int) (map[int]bool, *big.Int) {

	newBits := make(map[int]bool)
	newFeeSum := new(big.Int).Set(feeSum)
	for sol := range bits {
		if processInfo.BitVotes[voteIndex].BitVector.Bit(sol) == 0 {
			newFeeSum.Sub(newFeeSum, processInfo.Fees[sol].Fee)
		} else {
			newBits[sol] = true
		}
		currentStatus.NumOperations++
	}

	return newBits, newFeeSum

}

// joinResultsVotes compares two solutions.
// If result1 (the result in which the vote on place branch is included) is greater, result1 with updated votes is returned.
// Otherwise, result0 without the vote on branch place is returned.
func joinResultsVotes(result0, result1 *BranchAndBoundPartialSolution, branch int) *BranchAndBoundPartialSolution {
	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		delete(result0.Votes, branch)
		return result0
	} else if result0 == nil || result0.Value.Cmp(result1.Value) == -1 {
		result1.Votes[branch] = true
		return result1
	} else {
		delete(result0.Votes, branch)
		return result0
	}

}

// MaximizeVotes adds all votes that confirm all bits in the solution and updates the value.
func (solution *BranchAndBoundPartialSolution) MaximizeVotes(votes []*AggregatedVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) {
	for i := range votes {

		if _, isIncluded := solution.Votes[i]; !isIncluded {
			check := true
			for j := range solution.Bits {
				if votes[i].BitVector.Bit(fees[i].Indexes[j]) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Votes[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(votes, fees, assumedFees, assumedWeight, totalWeight)
}
