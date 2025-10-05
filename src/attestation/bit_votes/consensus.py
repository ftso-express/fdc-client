import math
from dataclasses import dataclass
from typing import List, Dict, Set, Optional

from src.attestation.bit_votes.structures import (
    AggregatedVote,
    AggregatedBit,
    Value,
)
from src.attestation.bit_votes.preprocessing import calc_value

# --- Data Structures for the Algorithm ---

@dataclass
class ConsensusSolution:
    """Represents the final output of the consensus algorithm."""
    votes: List[AggregatedVote]
    bits: List[AggregatedBit]
    value: Value
    optimal: bool

@dataclass
class _PartialSolution:
    """Represents a potential solution found during the recursive search."""
    votes: Set[int]
    bits: Set[int]
    value: Value

@dataclass
class _ProcessInfo:
    """Holds the static information for a run of the algorithm."""
    total_weight: int
    lower_bound_weight: int
    bit_votes: List[AggregatedVote]
    bits: List[AggregatedBit]
    max_operations: int
    exclude_first: bool
    # Pre-calculated total fee for pruning
    total_fee: int
    assumed_fees: int
    assumed_weight: int


class _SharedStatus:
    """Holds the shared state that is updated during the search."""
    def __init__(self, initial_bound: Value):
        self.current_bound = initial_bound
        self.num_operations = 0

# --- Main Algorithm Class ---

class BranchAndBound:
    """
    Encapsulates the state and logic for the branch-and-bound consensus algorithm.
    """
    def __init__(self, process_info: _ProcessInfo, shared_status: _SharedStatus):
        self.info = process_info
        self.status = shared_status

    def run(self) -> Optional[_PartialSolution]:
        """Starts the recursive search."""
        initial_included_votes = set(range(len(self.info.bit_votes)))
        initial_weight = self.info.assumed_weight + sum(v.weight for v in self.info.bit_votes)

        return self._branch(0, initial_included_votes, initial_weight, self.info.total_fee)

    def _branch(
        self,
        branch_idx: int,
        included_votes: Set[int],
        current_weight: int,
        potential_fee: int,
    ) -> Optional[_PartialSolution]:
        """The core recursive function of the branch-and-bound algorithm."""
        self.status.num_operations += 1

        if branch_idx == len(self.info.bits):
            value = calc_value(potential_fee, current_weight, self.info.total_weight)
            if value.cmp(self.status.current_bound) > 0:
                self.status.current_bound = value
                return _PartialSolution(votes=included_votes, bits=set(), value=value)
            return None

        if self.status.num_operations >= self.info.max_operations:
            return None

        potential_value = calc_value(potential_fee, current_weight, self.info.total_weight)
        if potential_value.cmp(self.status.current_bound) <= 0:
            return None

        if self.info.exclude_first:
            res_exclude = self._branch_exclude(branch_idx, included_votes, current_weight, potential_fee)
            res_include = self._branch_include(branch_idx, included_votes, current_weight, potential_fee)
        else:
            res_include = self._branch_include(branch_idx, included_votes, current_weight, potential_fee)
            res_exclude = self._branch_exclude(branch_idx, included_votes, current_weight, potential_fee)

        return self._join_results(res_exclude, res_include, branch_idx)

    def _branch_exclude(self, branch_idx: int, included_votes: Set[int], current_weight: int, potential_fee: int):
        bit_to_exclude = self.info.bits[branch_idx]
        return self._branch(
            branch_idx + 1,
            included_votes,
            current_weight,
            potential_fee - bit_to_exclude.fee
        )

    def _branch_include(self, branch_idx: int, included_votes: Set[int], current_weight: int, potential_fee: int):
        bit_to_include = self.info.bits[branch_idx]
        new_included_votes = set()
        new_current_weight = self.info.assumed_weight

        for vote_idx in included_votes:
            vote = self.info.bit_votes[vote_idx]
            original_bit_index = bit_to_include.indexes[0]
            if (vote.bit_vector >> original_bit_index) & 1:
                new_included_votes.add(vote_idx)
                new_current_weight += vote.weight

        if new_current_weight <= self.info.lower_bound_weight:
            return None

        return self._branch(branch_idx + 1, new_included_votes, new_current_weight, potential_fee)

    @staticmethod
    def _join_results(res0: Optional[_PartialSolution], res1: Optional[_PartialSolution], branch_idx: int):
        if res0 is None and res1 is None: return None
        if res1 is None: return res0
        if res0 is None or res0.value.cmp(res1.value) < 0:
            res1.bits.add(branch_idx)
            return res1
        return res0

# --- Helper and Runner Functions ---

def _sort_bits_by_value(bits: List[AggregatedBit], total_weight: int, reverse: bool):
    def get_value(bit: AggregatedBit):
        if bit.value is None:
            bit.value = calc_value(bit.fee, bit.support, total_weight)
        return bit.value
    return sorted(bits, key=lambda b: (get_value(b).capped_value, get_value(b).uncapped_value, b.indexes[0]), reverse=reverse)

def _maximize_bits(solution: _PartialSolution, info: _ProcessInfo):
    current_bits_indices = {info.bits[i].indexes[0] for i in solution.bits}
    for i, bit in enumerate(info.bits):
        if i in solution.bits: continue
        original_bit_index = bit.indexes[0]
        if original_bit_index in current_bits_indices: continue
        is_supported = all((info.bit_votes[vote_idx].bit_vector >> original_bit_index) & 1 for vote_idx in solution.votes)
        if is_supported:
            solution.bits.add(i)

def branch_and_bound_bits(
    bit_votes: List[AggregatedVote], bits: List[AggregatedBit], assumed_weight: int,
    total_weight: int, assumed_fees: int, max_ops: int,
    initial_bound: Value, exclude_first: bool
) -> ConsensusSolution:
    total_fee = assumed_fees + sum(b.fee for b in bits)
    info = _ProcessInfo(
        total_weight=total_weight, lower_bound_weight=total_weight // 2,
        bit_votes=bit_votes, bits=bits, max_operations=max_ops,
        exclude_first=exclude_first, total_fee=total_fee,
        assumed_fees=assumed_fees, assumed_weight=assumed_weight
    )
    status = _SharedStatus(initial_bound=initial_bound)
    provisional_result = BranchAndBound(info, status).run()
    is_optimal = status.num_operations < max_ops

    if provisional_result is None:
        return ConsensusSolution(votes=[], bits=[], value=Value(0, 0), optimal=is_optimal)

    if not is_optimal:
        _maximize_bits(provisional_result, info)

    final_fee = assumed_fees + sum(info.bits[i].fee for i in provisional_result.bits)
    final_weight = assumed_weight + sum(info.bit_votes[i].weight for i in provisional_result.votes)
    final_value = calc_value(final_fee, final_weight, total_weight)
    provisional_result.value = final_value

    final_votes = [info.bit_votes[i] for i in provisional_result.votes]
    final_bits = [info.bits[i] for i in provisional_result.bits]

    return ConsensusSolution(votes=final_votes, bits=final_bits, value=provisional_result.value, optimal=is_optimal)

def branch_and_bound_bits_double(
    bit_votes: List[AggregatedVote], bits: List[AggregatedBit], assumed_weight: int,
    total_weight: int, assumed_fees: int, max_ops: int, initial_bound: Value
) -> ConsensusSolution:
    bits_dsc = _sort_bits_by_value(list(bits), total_weight, reverse=True)
    solution1 = branch_and_bound_bits(
        bit_votes, bits_dsc, assumed_weight, total_weight,
        assumed_fees, max_ops, initial_bound, exclude_first=False
    )
    if solution1.optimal: return solution1

    new_bound = solution1.value
    bits_asc = _sort_bits_by_value(list(bits), total_weight, reverse=False)
    solution2 = branch_and_bound_bits(
        bit_votes, bits_asc, assumed_weight, total_weight,
        assumed_fees, max_ops, new_bound, exclude_first=True
    )
    return solution1 if solution1.value.cmp(solution2.value) >= 0 else solution2