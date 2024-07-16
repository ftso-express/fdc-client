package bitvotes

import (
	"local/fdc/client/utils"
	"math/big"
	"slices"
	"strconv"
)

type FilterResults struct {
	AlwaysInBits   []int
	AlwaysOutBits  []int
	RemainingBits  map[int]bool
	GuaranteedFees *big.Int

	AlwaysInVotes    []int
	AlwaysOutVotes   []int
	RemainingVotes   map[int]bool
	GuaranteedWeight uint16
	RemainingWeight  uint16
}

// FilterBits identifies the bits that are supported by all or are supported by none
// and moves them from RemainingBits to AlwaysInBits or AlwaysOutBits, respectively.
// Fees of bits moved to AlwaysInBits are added to guaranteedFees.
func FilterBits(bitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, results *FilterResults) *FilterResults {

	for i := range results.RemainingBits {
		support := results.GuaranteedWeight

		for j := range results.RemainingVotes {
			if bitVotes[j].BitVote.BitVector.Bit(i) == 1 {
				support += bitVotes[j].Weight
			}
		}

		if support == results.RemainingWeight {

			results.AlwaysInBits = append(results.AlwaysInBits, i)
			results.GuaranteedFees.Add(results.GuaranteedFees, fees[i])

			delete(results.RemainingBits, i)

		} else if support <= totalWeight/2 {
			results.AlwaysOutBits = append(results.AlwaysOutBits, i)

			delete(results.RemainingBits, i)

		}

	}

	return results

}

// FilterBitsOnes moves bits that are supported by all RemainingVotes to AlwaysInBits and updates GuaranteedFees.
func FilterBitsOnes(bitVotes []*WeightedBitVote, fees []*big.Int, results *FilterResults) *FilterResults {

bits:
	for i := range results.RemainingBits {

		for j := range results.RemainingVotes {
			if bitVotes[j].BitVote.BitVector.Bit(i) == 0 {
				continue bits
			}
		}

		results.AlwaysInBits = append(results.AlwaysInBits, i)
		results.GuaranteedFees.Add(results.GuaranteedFees, fees[i])

		delete(results.RemainingBits, i)

	}

	return results

}

// FilterVoters moves votes from RemainingVotes that are all zero or all one on RemainingBits to AlwaysOutVotes
// or AlwaysInVotes, respectively.
//
// GuraranteedWight and RemainingWeight are updated accrodingly.
func FilterVotes(bitVotes []*WeightedBitVote, totalWeight uint16, results *FilterResults) (*FilterResults, bool) {

	somethingChanged := false

votes:
	for i := range results.RemainingVotes {

		allOnes := true
		allZeros := true

		for j := range results.RemainingBits {

			if !allOnes && !allZeros {

				continue votes

			}

			if allOnes && bitVotes[i].BitVote.BitVector.Bit(j) == 0 {

				allOnes = false

			}

			if allZeros && bitVotes[i].BitVote.BitVector.Bit(j) == 1 {
				allZeros = false
			}

		}

		if allZeros && len(results.AlwaysInBits) == 0 {

			somethingChanged = true

			results.AlwaysOutVotes = append(results.AlwaysOutVotes, i)

			delete(results.RemainingVotes, i)

		} else if allOnes {

			somethingChanged = true

			results.AlwaysInVotes = append(results.AlwaysInVotes, i)

			results.GuaranteedWeight += bitVotes[i].Weight

			delete(results.RemainingVotes, i)

		}

	}

	return results, somethingChanged

}

// Filter identifies the bits and votes that are guaranteed to be included in the selcection of the consensus bitVote.
func Filter(bitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16) *FilterResults {

	remainingBits := make(map[int]bool)
	for i := range fees {
		remainingBits[i] = true
	}

	remainingWeight := uint16(0)

	remainingVotes := make(map[int]bool)
	for i := range bitVotes {
		remainingVotes[i] = true
		remainingWeight += bitVotes[i].Weight
	}

	resultsValue := FilterResults{
		AlwaysInBits:   []int{},
		AlwaysOutBits:  []int{},
		RemainingBits:  remainingBits,
		GuaranteedFees: big.NewInt(0),

		AlwaysInVotes:    []int{},
		AlwaysOutVotes:   []int{},
		RemainingVotes:   remainingVotes,
		GuaranteedWeight: 0,
		RemainingWeight:  remainingWeight,
	}

	results := &resultsValue

	results = FilterBits(bitVotes, fees, totalWeight, results)

	results, changed := FilterVotes(bitVotes, totalWeight, results)

	if changed {
		results = FilterBitsOnes(bitVotes, fees, results)
	}

	return results

}

type AggregatedFee struct {
	Fee     *big.Int
	Indexes []int
}

// AggreagteBits aggregates fees of the bits that agree on all the RemainingVotes.
func AggregateBits(bitVotes []*WeightedBitVote, fees []*big.Int, filterResults *FilterResults) []*AggregatedFee {

	aggregator := map[string]*AggregatedFee{}

	remainingBitsSorted := utils.Keys(filterResults.RemainingBits)

	slices.Sort(remainingBitsSorted)

	remainingVotesSorted := utils.Keys(filterResults.RemainingVotes)

	slices.Sort(remainingVotesSorted)

	for _, i := range remainingBitsSorted {

		state := ""

		for _, j := range remainingVotesSorted {

			bit := bitVotes[j].BitVote.BitVector.Bit(i)

			state += strconv.FormatUint(uint64(bit), 10)

		}

		aggFee, exists := aggregator[state]

		if !exists {

			newIndexedFee := AggregatedFee{Fee: new(big.Int).Set(fees[i]), Indexes: []int{i}}

			aggregator[state] = &newIndexedFee

		} else {
			aggFee.Fee.Add(aggFee.Fee, fees[i])

			if i < aggFee.Indexes[0] {
				aggFee.Indexes = utils.Prepend(aggFee.Indexes, i)
			} else {
				aggFee.Indexes = append(aggFee.Indexes, i)

			}
		}
	}

	return utils.Values(aggregator)

}

type AggregatedVote struct {
	BitVector *big.Int
	Weight    uint16
	Indexes   []int //places of the voted in initial []*WeightedBitVotes
}

func AggregateVotes(bitVotes []*WeightedBitVote, fees []*big.Int, filterResults *FilterResults) []*AggregatedVote {
	aggregator := map[string]*AggregatedVote{} // maps identifer to the place in the aggregatedVotes

	remainingBitsSorted := utils.Keys(filterResults.RemainingBits)

	slices.Sort(remainingBitsSorted)

	remainingVotesSorted := utils.Keys(filterResults.RemainingVotes)

	slices.Sort(remainingVotesSorted)

	for _, i := range remainingVotesSorted {

		identifier := ""

		for _, j := range remainingBitsSorted {

			bit := bitVotes[i].BitVote.BitVector.Bit(j)

			identifier += strconv.FormatUint(uint64(bit), 10)

		}

		aggVote, exists := aggregator[identifier]

		if !exists {
			newAggregatedBitVote := AggregatedVote{
				BitVector: new(big.Int).Set(bitVotes[i].BitVote.BitVector),
				Weight:    bitVotes[i].Weight,
				Indexes:   []int{i},
			}

			aggregator[identifier] = &newAggregatedBitVote

		} else {

			aggVote.Weight += bitVotes[i].Weight

			if bitVotes[i].Index < aggVote.Indexes[0] { // indexes[0] always exists!
				aggVote.Indexes = utils.Prepend(aggVote.Indexes, i)
			} else {
				aggVote.Indexes = append(aggVote.Indexes, i)

			}
		}

	}

	return utils.Values(aggregator)

}

func FilterAndAggregate(bitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16) ([]*AggregatedVote, []*AggregatedFee, *FilterResults) {

	filterResults := Filter(bitVotes, fees, totalWeight)

	aggregatedVotes := AggregateVotes(bitVotes, fees, filterResults)

	aggregatedFees := AggregateBits(bitVotes, fees, filterResults)

	return aggregatedVotes, aggregatedFees, filterResults
}

func AssembleSolution(filterResults *FilterResults, filteredSolution ConsensusSolution, aggregatedFees []*AggregatedFee) *big.Int {

	consensusBitVote := big.NewInt(0)

	for _, i := range filterResults.AlwaysInBits {

		consensusBitVote.SetBit(consensusBitVote, i, 1)

	}

	for k := range filteredSolution.Solution {
		indexes := aggregatedFees[k].Indexes

		for _, i := range indexes {
			consensusBitVote.SetBit(consensusBitVote, i, 1)

		}
	}

	return consensusBitVote

}

type Solution struct {
	Bits    []int
	Votes   []int
	Value   Value
	Optimal bool
}

func AssembleSolutionFull(filterResults *FilterResults, filteredSolution ConsensusSolution, aggregatedFees []*AggregatedFee, aggregatedVotes []*AggregatedVote) Solution {

	bits := []int{}

	bits = append(bits, filterResults.AlwaysInBits...)

	for k := range filteredSolution.Solution {

		indexes := aggregatedFees[k].Indexes

		bits = append(bits, indexes...)

	}

	voters := []int{}

	voters = append(voters, filterResults.AlwaysInVotes...)

	for k := range filteredSolution.Participants {
		indexes := aggregatedVotes[k].Indexes

		voters = append(voters, indexes...)

	}

	return Solution{
		Bits:    bits,
		Votes:   voters,
		Value:   filteredSolution.Value,
		Optimal: filteredSolution.Optimal,
	}

}
