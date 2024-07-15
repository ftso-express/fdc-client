package bitvotes

import (
	"local/fdc/client/utils"
	"math/big"
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
	RemovedWeight    uint16
}

func FilterBits(bitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16, results *FilterResults) (*FilterResults, bool) {

	somethingChanged := false

	remainingBits := results.RemainingBits

	remainingVotes := results.RemainingVotes

	for i := range remainingBits {
		support := results.GuaranteedWeight

		for j := range remainingVotes {
			if bitVotes[j].BitVote.BitVector.Bit(i) == 1 {
				support += bitVotes[j].Weight
			}
		}

		if support == results.RemainingWeight {

			results.AlwaysInBits = append(results.AlwaysInBits, i)

			delete(results.RemainingBits, i)

			results.GuaranteedFees.Add(results.GuaranteedFees, fees[i])
			somethingChanged = true
		} else if support <= totalWeight/2 {
			results.AlwaysOutBits = append(results.AlwaysOutBits, i)
			delete(results.RemainingBits, i)

			somethingChanged = true

		}

	}

	return results, somethingChanged

}

func FilterVotes(bitVotes []*WeightedBitVote, totalWeight uint16, results *FilterResults) (*FilterResults, bool) {

	somethingChanged := false

	remainingBits := results.RemainingBits

	remainingVotes := results.RemainingVotes

votes:
	for i := range remainingVotes {

		allOnes := true
		allZeros := true

		for j := range remainingBits {

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

		if allZeros {

			results.AlwaysOutVotes = append(results.AlwaysOutVotes, i)

			somethingChanged = true

			delete(results.RemainingVotes, i)

			results.RemainingWeight -= bitVotes[i].Weight

			results.RemovedWeight += bitVotes[i].Weight

		} else if allOnes {
			results.AlwaysInVotes = append(results.AlwaysInVotes, i)

			somethingChanged = true

			results.GuaranteedWeight += bitVotes[i].Weight

			delete(results.RemainingVotes, i)

		}

	}

	return results, somethingChanged

}

func Filter(bitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16) *FilterResults {

	remainingBits := make(map[int]bool)
	for i := range fees {
		remainingBits[i] = true
	}

	remainingVotes := make(map[int]bool)
	for i := range bitVotes {
		remainingVotes[i] = true
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
		RemainingWeight:  totalWeight,
		RemovedWeight:    0,
	}

	results := &resultsValue

	results, _ = FilterBits(bitVotes, fees, totalWeight, results)

	results, changed := FilterVotes(bitVotes, totalWeight, results)

	for changed {
		results, changed = FilterBits(bitVotes, fees, totalWeight, results)

		if !changed {
			break
		}

		results, changed = FilterVotes(bitVotes, totalWeight, results)
	}

	return results

}

type AggregatedFee struct {
	Fee     *big.Int
	Indexes []int
}

func AggregateBits(bitVotes []*WeightedBitVote, fees []*big.Int, filterResults *FilterResults) []*AggregatedFee {

	aggregator := map[string]int{}

	index := map[int][]int{}

	counter := 0

	aggregatedFees := []*AggregatedFee{}

	for i := range filterResults.RemainingBits {

		state := ""

		for j := range filterResults.RemainingVotes {

			bit := bitVotes[j].BitVote.BitVector.Bit(i)

			state += strconv.FormatUint(uint64(bit), 10)

		}

		k, exists := aggregator[state]

		if !exists {
			aggregator[state] = counter

			newIndexedFee := AggregatedFee{Fee: new(big.Int).Set(fees[i]), Indexes: []int{i}}
			aggregatedFees = append(aggregatedFees, &newIndexedFee)

			index[counter] = []int{i}

			counter++

		} else {
			index[k] = append(index[k], i)
			aggregatedFees[k].Fee.Add(aggregatedFees[k].Fee, fees[i])

			if i < aggregatedFees[k].Indexes[0] {
				aggregatedFees[k].Indexes = utils.Prepend(aggregatedFees[k].Indexes, i)
			} else {
				aggregatedFees[k].Indexes = append(aggregatedFees[k].Indexes, i)

			}

		}

	}

	return aggregatedFees

}

type AggregatedBitVote struct {
	BitVector *big.Int
	Weight    uint16
	Indexes   []int
}

func AggregateVotes(bitVotes []*WeightedBitVote, fees []*big.Int, filterResults *FilterResults) []*AggregatedBitVote {
	aggregator := map[string]int{}

	aggregatedVotes := []*AggregatedBitVote{}

	for i := range filterResults.RemainingVotes {

		state := ""

		for j := range filterResults.RemainingBits {

			bit := bitVotes[i].BitVote.BitVector.Bit(j)

			state += strconv.FormatUint(uint64(bit), 10)

		}

		k, exists := aggregator[state]

		if !exists {
			newAggregatedBitVote := AggregatedBitVote{
				BitVector: new(big.Int).Set(bitVotes[i].BitVote.BitVector),
				Weight:    bitVotes[i].Weight,
				Indexes:   []int{i},
			}
			aggregatedVotes = append(aggregatedVotes, &newAggregatedBitVote)

		} else {

			aggregatedVotes[k].Weight += bitVotes[i].Weight

			if bitVotes[i].Index < aggregatedVotes[k].Indexes[0] { // indexes[0] always exists!
				aggregatedVotes[k].Indexes = utils.Prepend(aggregatedVotes[k].Indexes, i)
			} else {
				aggregatedVotes[k].Indexes =
					append(aggregatedVotes[k].Indexes, i)

			}
		}

	}

	return aggregatedVotes

}

func FilterAndAggregate(bitVotes []*WeightedBitVote, fees []*big.Int, totalWeight uint16) ([]*AggregatedBitVote, []*AggregatedFee, *FilterResults) {

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

func AssembleSolutionFull(filterResults *FilterResults, filteredSolution ConsensusSolution, aggregatedFees []*AggregatedFee, aggregatedVotes []*AggregatedBitVote) Solution {

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
