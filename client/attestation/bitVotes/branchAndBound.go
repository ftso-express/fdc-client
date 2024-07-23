package bitvotes

import (
	"math"
	"math/big"
	"math/rand"
	"slices"
)

type SharedStatus struct {
	CurrentBound  Value
	NumOperations int
	RandGen       rand.Source
}

type ProcessInfo struct {
	TotalWeight      uint16
	LowerBoundWeight uint16
	BitVotes         []*AggregatedVote
	Fees             []*AggregatedFee
	NumAttestations  int
	NumProviders     int

	MaxOperations int
	Strategy      func(...interface{}) bool
}

type BranchAndBoundPartialSolution struct {
	Votes map[int]bool
	Bits  map[int]bool
	Value Value
}

type Value struct {
	CappedValue   *big.Int
	UncappedValue *big.Int
}

// Cmp compares Values lexicographically
//
//	-1 if v0 <  v1
//	 0 if v0 == v1
//	+1 if v0 >  v1
func (v0 Value) Cmp(v1 Value) int {
	firstCompare := v0.CappedValue.Cmp(v1.CappedValue)

	if firstCompare != 0 {
		return firstCompare
	}

	return v0.UncappedValue.Cmp(v1.UncappedValue)

}

func CalcValue(feeSum *big.Int, weight, totalWeight uint16) Value {
	weightCaped := min(int64(math.Ceil(float64(totalWeight)*valueCap)), int64(weight))

	return Value{CappedValue: new(big.Int).Mul(feeSum, big.NewInt(weightCaped)),
		UncappedValue: new(big.Int).Mul(feeSum, big.NewInt(int64(weight))),
	}
}

func cmpVal(totalWeight uint16, sign int) func(*AggregatedFee, *AggregatedFee) int {

	return func(fee0, fee1 *AggregatedFee) int {

		val0 := fee0.Value(totalWeight, true)

		val1 := fee1.Value(totalWeight, true)

		cmp := val0.Cmp(val1)

		if cmp != 0 {
			return sign * cmp
		}

		if fee0.Indexes[0] < fee1.Indexes[0] {
			return -1
		}
		if fee0.Indexes[0] > fee1.Indexes[0] {
			return 1
		}

		return 0

	}

}

func cmpValAsc(totalWeight uint16) func(*AggregatedFee, *AggregatedFee) int {
	return cmpVal(totalWeight, 1)
}
func cmpValDsc(totalWeight uint16) func(*AggregatedFee, *AggregatedFee) int {
	return cmpVal(totalWeight, -1)
}

func sortFees(fees []*AggregatedFee, sortFunc func(*AggregatedFee, *AggregatedFee) int) []*AggregatedFee {

	sortedFees := make([]*AggregatedFee, len(fees))
	copy(sortedFees, fees)
	slices.SortStableFunc(sortedFees, sortFunc)

	return sortedFees
}

// BranchAndBoundBitsDouble runs two branch and bound strategies on bits in parallel and returns the better result.
//
// The first strategy sorts the aggregated fees by the descending value (cappedSupport * fee) and at depth k the branch in which k-th bit is included is explored first.
// The second strategy sorts the aggregated fees by the ascending value (cappedSupport * fee) and at depth k the branch in which does not include k-th bit is explored first.
//
// If both strategies find an optimal but different solutions, the solution of the first strategy is returned.
func BranchAndBoundBitsDouble(bitVotes []*AggregatedVote, fees []*AggregatedFee, assumedWeight, absoluteTotalWeight uint16, assumedFees *big.Int, maxOperations int, initialBound Value) *ConsensusSolution {

	solutions := make([]*ConsensusSolution, 2)

	firstDone := make(chan bool, 1)
	secondDone := make(chan bool, 1)

	feesAscVal := sortFees(fees, cmpValAsc(absoluteTotalWeight))

	feesDscVal := sortFees(fees, cmpValDsc(absoluteTotalWeight))

	go func() {

		solution := BranchAndBoundBits(bitVotes, feesAscVal, assumedWeight, absoluteTotalWeight, assumedFees, maxOperations, initialBound, func(...interface{}) bool { return true })

		solutions[0] = solution

		firstDone <- true

		if solution.Optimal {
			solutions[1] = nil
			secondDone <- true // do not wait on the other solution
		}

	}()

	go func() {

		solution := BranchAndBoundBits(bitVotes, feesDscVal, assumedWeight, absoluteTotalWeight, assumedFees, maxOperations, initialBound, func(...interface{}) bool { return false })

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

// BranchAndBoundBits takes a set of aggregated bitVotes and a list of aggregated fees and
// tries to get an optimal subset of votes with the weight more than the half of the total weight.
// The algorithm executes a branch and bound strategy on the space of subsets of attestations, hence
// it is particularly useful when there are not too many attestations. If the algorithm is able search
// through the entire solution space before reaching the given max operations counter, the algorithm
// gives an optimal solution. If solution space is too big, the algorithm gives a
// the best solution it finds. If not solution exceeding initialBound is found, no solution (nil) is returned.
func BranchAndBoundBits(bitVotes []*AggregatedVote, fees []*AggregatedFee, assumedWeight, absoluteTotalWeight uint16, assumedFees *big.Int, maxOperations int, initialBound Value, strategy func(...interface{}) bool) *ConsensusSolution {

	weight := assumedWeight

	votes := make(map[int]bool)
	for i, vote := range bitVotes {
		weight += vote.Weight
		votes[i] = true
	}

	totalFee := big.NewInt(0).Set(assumedFees)
	for _, fee := range fees {
		totalFee.Add(totalFee, fee.Fee)
	}

	// randGen := rand.NewSource(seed)
	// randPerm := RandPerm(numAttestations, randGen)
	// permBitVotes := PermuteBits(bitVotes, randPerm)

	currentStatus := &SharedStatus{
		CurrentBound:  initialBound,
		NumOperations: 0,
		// RandGen:       randGen,
	}

	processInfo := &ProcessInfo{
		TotalWeight:      absoluteTotalWeight,
		LowerBoundWeight: absoluteTotalWeight / 2,
		BitVotes:         bitVotes,
		Fees:             fees,
		NumAttestations:  len(fees),
		NumProviders:     len(bitVotes),
		MaxOperations:    maxOperations,
		Strategy:         strategy,
	}

	permResult := BranchBits(processInfo, currentStatus, 0, votes, weight, totalFee)

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

	if !isOptimal {
		permResult.MaximizeBits(bitVotes, fees, assumedFees, assumedWeight, absoluteTotalWeight)
	}

	result := ConsensusSolution{
		Votes:   make([]*AggregatedVote, 0),
		Bits:    make([]*AggregatedFee, 0),
		Optimal: isOptimal,
		Value:   permResult.Value,
	}

	for key := range permResult.Votes {
		result.Votes = append(result.Votes, bitVotes[key])
	}
	for key := range permResult.Bits {
		result.Bits = append(result.Bits, fees[key])
	}

	return &result
}

func BranchBits(processInfo *ProcessInfo, currentStatus *SharedStatus, branch int, participants map[int]bool, currentWeight uint16, feeSum *big.Int) *BranchAndBoundPartialSolution {

	currentStatus.NumOperations++

	// end of recursion
	if branch == processInfo.NumAttestations {

		value := CalcValue(feeSum, currentWeight, processInfo.TotalWeight)

		if value.Cmp(currentStatus.CurrentBound) == 1 {
			currentStatus.CurrentBound = value

			return &BranchAndBoundPartialSolution{
				Votes: participants,
				Bits:  make(map[int]bool),
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
	if CalcValue(feeSum, currentWeight, processInfo.TotalWeight).Cmp(currentStatus.CurrentBound) != 1 {
		return nil
	}

	var result0 *BranchAndBoundPartialSolution
	var result1 *BranchAndBoundPartialSolution

	if processInfo.Strategy() {

		result0 = BranchBits(processInfo, currentStatus, branch+1, participants, currentWeight, new(big.Int).Sub(feeSum, processInfo.Fees[branch].Fee))
	}

	// prepare and check if a branch is possible
	newParticipants, newCurrentWeight := prepareDataForBranchWithOne(processInfo, currentStatus, participants, currentWeight, branch)

	if newCurrentWeight > processInfo.LowerBoundWeight {

		result1 = BranchBits(processInfo, currentStatus, branch+1, newParticipants, newCurrentWeight, feeSum)
	}

	if !processInfo.Strategy() {
		result0 = BranchBits(processInfo, currentStatus, branch+1, participants, currentWeight, new(big.Int).Sub(feeSum, processInfo.Fees[branch].Fee))
	}

	// max result
	return joinResultsAttestations(result0, result1, branch)
}

// prepareDataForBranchWithOne prepares data for branch in which the bit on bitIndex place is included.
func prepareDataForBranchWithOne(processInfo *ProcessInfo, currentStatus *SharedStatus, votes map[int]bool, currentWeight uint16, bitIndex int) (map[int]bool, uint16) {

	bit := processInfo.Fees[bitIndex].Indexes[0]
	newParticipants := make(map[int]bool)
	newCurrentWeight := currentWeight

	for i := range votes {
		if processInfo.BitVotes[i].BitVector.Bit(bit) == 1 {

			newParticipants[i] = true
		} else {

			newCurrentWeight -= processInfo.BitVotes[i].Weight

		}
		currentStatus.NumOperations++
	}

	return newParticipants, newCurrentWeight

}

// joinResultsAttestations compares two solutions.
// If result1 (branch in which the bit in place branch is included) is greater, result1 with updated bits is returned.
// Otherwise, result0 without the bit it the branch place is returned.
func joinResultsAttestations(result0, result1 *BranchAndBoundPartialSolution, branch int) *BranchAndBoundPartialSolution {

	if result0 == nil && result1 == nil {
		return nil
	} else if result0 != nil && result1 == nil {
		delete(result0.Bits, branch)
		return result0
	} else if result0 == nil || result0.Value.Cmp(result1.Value) == -1 {

		result1.Bits[branch] = true

		return result1
	} else {
		delete(result0.Bits, branch)
		return result0
	}

}

// MaximizeBits adds all bits that are supported by all the votes in the solution and updates the Value.
func (solution *BranchAndBoundPartialSolution) MaximizeBits(votes []*AggregatedVote, fees []*AggregatedFee, assumedFees *big.Int, assumedWeight, totalWeight uint16) {
	for i := range fees {

		if _, isConfirmed := solution.Bits[i]; !isConfirmed {
			check := true
			for j := range solution.Votes {
				if votes[j].BitVector.Bit(fees[i].Indexes[0]) == 0 {
					check = false
					break
				}
			}
			if check {
				solution.Bits[i] = true
			}
		}
	}

	solution.Value = solution.CalcValueFromFees(votes, fees, assumedFees, assumedWeight, totalWeight)
}
