import pytest
from src.attestation.bit_votes.structures import (
    AggregatedVote,
    AggregatedBit,
    Value,
)
from src.attestation.bit_votes.consensus import branch_and_bound_bits_double

# --- Test Fixture ---

@pytest.fixture
def consensus_scenario():
    """
    Provides a scenario for testing the branch-and-bound algorithm.
    - 3 aggregated bits to choose from.
    - 2 aggregated voters.
    - Total weight of all original voters is 100.
    """
    # Voter 0 (w=40) supports bits 0, 1. Vector: 011
    # Voter 1 (w=30) supports bits 0, 2. Vector: 101
    agg_votes = [
        AggregatedVote(bit_vector=0b011, weight=40, indexes=[0], fees=0),
        AggregatedVote(bit_vector=0b101, weight=30, indexes=[1], fees=0),
    ]

    # Bit 0 (fee=100), supported by both v0, v1. Support = 70.
    # Bit 1 (fee=300), supported by v0. Support = 40.
    # Bit 2 (fee=300), supported by v1. Support = 30.
    agg_bits = [
        AggregatedBit(fee=100, indexes=[0], support=70, value=None),
        AggregatedBit(fee=300, indexes=[1], support=40, value=None),
        AggregatedBit(fee=300, indexes=[2], support=30, value=None),
    ]

    # Base case: no pre-filtered items
    assumed_weight = 0
    assumed_fees = 0

    # Total weight of all voters in the system
    total_weight = 100

    # A very high number for max operations to ensure we find the optimal solution
    max_ops = 1000

    # Initial bound is zero
    initial_bound = Value(0, 0)

    return {
        "bit_votes": agg_votes,
        "bits": agg_bits,
        "assumed_weight": assumed_weight,
        "total_weight": total_weight,
        "assumed_fees": assumed_fees,
        "max_ops": max_ops,
        "initial_bound": initial_bound,
    }

# --- Test Case ---

def test_branch_and_bound_consensus(consensus_scenario):
    """
    Tests the full branch-and-bound algorithm with a manually traceable scenario.

    Possible solutions (must have weight > 50):
    1. Include Bit 0:
       - Voters supporting: v0, v1. Weight = 70.
       - Fees: 100.
       - Value = 100 * 70 = 7000.
    2. Include Bit 1:
       - Voters supporting: v0. Weight = 40. (INVALID, weight <= 50)
    3. Include Bit 2:
       - Voters supporting: v1. Weight = 30. (INVALID, weight <= 50)
    4. Include Bits 0, 1:
       - Voters supporting: v0. Weight = 40. (INVALID)
    5. Include Bits 0, 2:
       - Voters supporting: v1. Weight = 30. (INVALID)
    6. Include Bits 1, 2:
       - Voters supporting: none. Weight = 0. (INVALID)
    7. Include Bits 0, 1, 2:
       - Voters supporting: none. Weight = 0. (INVALID)

    The only valid solution is to include only Bit 0. The final solution
    should contain this bit.
    """

    solution = branch_and_bound_bits_double(**consensus_scenario)

    assert solution.optimal
    assert solution.value.capped_value == 7000

    # The solution should contain only the aggregated bit with index 0.
    assert len(solution.bits) == 1
    assert solution.bits[0].indexes[0] == 0

    # The solution should contain the voters that support this bit.
    assert len(solution.votes) == 2
    vote_indexes = {v.indexes[0] for v in solution.votes}
    assert vote_indexes == {0, 1}