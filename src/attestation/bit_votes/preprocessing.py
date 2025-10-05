import math
from typing import List, Dict
from src.attestation.bit_votes.structures import (
    WeightedBitVote,
    FilterResults,
    AggregatedBit,
    AggregatedVote,
    BitVote,
    Value,
)

# VALUE_CAP_FACTOR is the equivalent of valueCap in the Go code.
VALUE_CAP_FACTOR = 4.0 / 5.0

def calc_value(fee_sum: int, weight: int, total_weight: int) -> Value:
    """
    Calculates the value of a potential solution.
    A port of the `CalcValue` function from `branchAndBoundBits.go`.
    """
    if 2 * weight <= total_weight:
        return Value(capped_value=0, uncapped_value=0)

    weight_cap = math.ceil(total_weight * VALUE_CAP_FACTOR)
    capped_weight = min(weight_cap, weight)

    capped_value = fee_sum * int(capped_weight)
    uncapped_value = fee_sum * weight

    return Value(capped_value=capped_value, uncapped_value=uncapped_value)


def _filter_bits(fr: FilterResults, bit_votes: List[WeightedBitVote], fees: List[int], total_weight: int):
    """
    Identifies bits that are supported by all or by not more than 50% of the totalWeight.
    This is a helper for `filter_votes_and_bits` and modifies `fr` in place.
    """
    bits_to_remove = []
    for i in fr.remaining_bits:
        support = fr.guaranteed_weight
        for j in fr.remaining_votes:
            if (bit_votes[j].bit_vote.bit_vector >> i) & 1:
                support += bit_votes[j].weight

        if support == fr.remaining_weight:
            fr.always_in_bits.append(i)
            fr.guaranteed_fees += fees[i]
            bits_to_remove.append(i)
        elif support * 2 <= total_weight:
            fr.always_out_bits.append(i)
            bits_to_remove.append(i)

    for i in bits_to_remove:
        fr.remaining_bits.remove(i)

def _filter_votes(fr: FilterResults, bit_votes: List[WeightedBitVote]) -> bool:
    """
    Moves votes that are all-zero or all-one on the remaining bits.
    Returns True if any votes were moved. Modifies `fr` in place.
    """
    something_changed = False
    votes_to_remove = []
    for i in fr.remaining_votes:
        all_ones = True
        all_zeros = not fr.always_in_bits

        for j in fr.remaining_bits:
            if not all_ones and not all_zeros:
                break
            bit = (bit_votes[i].bit_vote.bit_vector >> j) & 1
            if all_ones and not bit:
                all_ones = False
            if all_zeros and bit:
                all_zeros = False

        if all_ones:
            something_changed = True
            fr.always_in_votes.append(i)
            fr.guaranteed_weight += bit_votes[i].weight
            votes_to_remove.append(i)
        elif all_zeros:
            something_changed = True
            fr.always_out_votes.append(i)
            votes_to_remove.append(i)

    for i in votes_to_remove:
        fr.remaining_votes.remove(i)
    return something_changed

def _filter_bits_ones(fr: FilterResults, bit_votes: List[WeightedBitVote], fees: List[int]):
    """
    Moves bits that are supported by all remaining votes to the 'always in' set.
    Modifies `fr` in place.
    """
    bits_to_remove = []
    for i in fr.remaining_bits:
        supported_by_all = True
        for j in fr.remaining_votes:
            if not ((bit_votes[j].bit_vote.bit_vector >> i) & 1):
                supported_by_all = False
                break
        if supported_by_all:
            fr.always_in_bits.append(i)
            fr.guaranteed_fees += fees[i]
            bits_to_remove.append(i)

    for i in bits_to_remove:
        fr.remaining_bits.remove(i)

def filter_votes_and_bits(bit_votes: List[WeightedBitVote], fees: List[int], total_weight: int) -> FilterResults:
    """
    Identifies bits and votes that are guaranteed to be in or out of the consensus.
    This is a Python port of the `Filter` function in `preprocess.go`.
    """
    fr = FilterResults(
        remaining_bits=set(range(len(fees))),
        remaining_votes=set(range(len(bit_votes))),
        remaining_weight=sum(v.weight for v in bit_votes)
    )

    _filter_bits(fr, bit_votes, fees, total_weight)
    if _filter_votes(fr, bit_votes):
        _filter_bits_ones(fr, bit_votes, fees)

    return fr

def aggregate_bits(bit_votes: List[WeightedBitVote], fees: List[int], fr: FilterResults) -> List[AggregatedBit]:
    """Aggregates bits that have the same support pattern among remaining voters."""
    aggregator: Dict[str, AggregatedBit] = {}
    remaining_bits_sorted = sorted(list(fr.remaining_bits))
    remaining_votes_sorted = sorted(list(fr.remaining_votes))

    for i in remaining_bits_sorted:
        identifier = []
        support = fr.guaranteed_weight
        for j in remaining_votes_sorted:
            bit = (bit_votes[j].bit_vote.bit_vector >> i) & 1
            if bit:
                support += bit_votes[j].weight
            identifier.append(str(bit))

        id_str = "".join(identifier)
        if id_str not in aggregator:
            aggregator[id_str] = AggregatedBit(fee=fees[i], indexes=[i], support=support)
        else:
            aggregator[id_str].fee += fees[i]
            aggregator[id_str].indexes.append(i)

    return list(aggregator.values())

def aggregate_votes(bit_votes: List[WeightedBitVote], fees: List[int], fr: FilterResults) -> List[AggregatedVote]:
    """Aggregates votes that have the same bit pattern on the remaining bits."""
    aggregator: Dict[str, AggregatedVote] = {}
    remaining_bits_sorted = sorted(list(fr.remaining_bits))
    remaining_votes_sorted = sorted(list(fr.remaining_votes))

    for i in remaining_votes_sorted:
        fees_vote = fr.guaranteed_fees
        identifier = []
        for j in remaining_bits_sorted:
            bit = (bit_votes[i].bit_vote.bit_vector >> j) & 1
            if bit:
                fees_vote += fees[j]
            identifier.append(str(bit))

        id_str = "".join(identifier)
        if id_str not in aggregator:
            aggregator[id_str] = AggregatedVote(
                bit_vector=bit_votes[i].bit_vote.bit_vector,
                weight=bit_votes[i].weight,
                indexes=[i],
                fees=fees_vote,
            )
        else:
            aggregator[id_str].weight += bit_votes[i].weight
            aggregator[id_str].indexes.append(i)

    return list(aggregator.values())

def filter_and_aggregate(bit_votes: List[WeightedBitVote], fees: List[int], total_weight: int) -> tuple[List[AggregatedVote], List[AggregatedBit], FilterResults]:
    """Runs the full filtering and aggregation pipeline."""
    filter_results = filter_votes_and_bits(bit_votes, fees, total_weight)
    aggregated_votes = aggregate_votes(bit_votes, fees, filter_results)
    aggregated_bits = aggregate_bits(bit_votes, fees, filter_results)
    return aggregated_votes, aggregated_bits, filter_results

def assemble_solution(fr: FilterResults, solution_bits: List[AggregatedBit], num_attestations: int) -> BitVote:
    """Assembles the final bit vector from the results of the consensus algorithm."""
    consensus_bit_vector = 0
    for i in fr.always_in_bits:
        consensus_bit_vector |= 1 << i

    for agg_bit in solution_bits:
        for i in agg_bit.indexes:
            consensus_bit_vector |= 1 << i

    return BitVote(length=num_attestations, bit_vector=consensus_bit_vector)