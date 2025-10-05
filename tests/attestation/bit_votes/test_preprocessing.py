import pytest
from src.attestation.bit_votes.structures import (
    BitVote,
    WeightedBitVote,
    IndexTx,
    AggregatedBit,
)
from src.attestation.bit_votes.preprocessing import (
    filter_votes_and_bits,
    aggregate_bits,
    aggregate_votes,
    filter_and_aggregate,
    assemble_solution,
)

# --- Test Fixtures ---

@pytest.fixture
def sample_votes_and_fees():
    """
    Provides a sample set of weighted bit votes and fees for testing.
    4 bits, 4 voters. Bit order is right-to-left (LSB is bit 0).
    Voter 0 (w=10): 1100 (12) -> supports bits 2, 3
    Voter 1 (w=10): 1100 (12) -> supports bits 2, 3
    Voter 2 (w=20): 1010 (10) -> supports bits 1, 3
    Voter 3 (w=30): 0011 (3)  -> supports bits 0, 1
    """
    bit_votes = [
        WeightedBitVote(0, IndexTx(1,0), 10, BitVote(4, 0b1100)),
        WeightedBitVote(1, IndexTx(1,1), 10, BitVote(4, 0b1100)),
        WeightedBitVote(2, IndexTx(1,2), 20, BitVote(4, 0b1010)),
        WeightedBitVote(3, IndexTx(1,3), 30, BitVote(4, 0b0011)),
    ]
    fees = [100, 200, 300, 400] # Fees for bits 0, 1, 2, 3
    total_weight = 70
    return bit_votes, fees, total_weight

# --- Corrected Test Cases ---

def test_filter_votes_and_bits(sample_votes_and_fees):
    """
    Tests the initial filtering of always-in and always-out bits/votes.
    - Bit 0 support: 30. 30*2 <= 70. OUT.
    - Bit 2 support: 20. 20*2 <= 70. OUT.
    - Remaining bits: {1, 3}.
    - Voter 2 has pattern '11' on remaining bits -> all-ones. IN.
    - Remaining voters: {0, 1, 3}.
    """
    bit_votes, fees, total_weight = sample_votes_and_fees
    fr = filter_votes_and_bits(bit_votes, fees, total_weight)

    assert sorted(fr.always_out_bits) == [0, 2]
    assert fr.remaining_bits == {1, 3}
    assert fr.guaranteed_fees == 0

    assert fr.always_in_votes == [2]
    assert fr.remaining_votes == {0, 1, 3}
    assert fr.guaranteed_weight == 20
    assert fr.remaining_weight == total_weight

def test_aggregate_votes(sample_votes_and_fees):
    """
    Tests aggregation on remaining voters {0, 1, 3} and bits {1, 3}.
    - Voters 0, 1 on {1,3} -> pattern '10' (bit 3=1, bit 1=0).
    - Voter 3 on {1,3} -> pattern '01' (bit 3=0, bit 1=1).
    Should result in 2 aggregated votes.
    """
    bit_votes, fees, total_weight = sample_votes_and_fees
    fr = filter_votes_and_bits(bit_votes, fees, total_weight)
    agg_votes = aggregate_votes(bit_votes, fees, fr)

    assert len(agg_votes) == 2

    agg_vote_10 = next(v for v in agg_votes if v.weight == 20)
    assert sorted(agg_vote_10.indexes) == [0, 1]

    agg_vote_01 = next(v for v in agg_votes if v.weight == 30)
    assert agg_vote_01.indexes == [3]

def test_aggregate_bits(sample_votes_and_fees):
    """
    Tests aggregation on remaining voters {0, 1, 3} and bits {1, 3}.
    - Bit 1 support from {0,1,3} -> pattern '001'.
    - Bit 3 support from {0,1,3} -> pattern '110'.
    Should result in 2 aggregated bits.
    """
    bit_votes, fees, total_weight = sample_votes_and_fees
    fr = filter_votes_and_bits(bit_votes, fees, total_weight)
    agg_bits = aggregate_bits(bit_votes, fees, fr)

    assert len(agg_bits) == 2

    agg_bit_1 = next(b for b in agg_bits if b.indexes == [1])
    assert agg_bit_1.fee == fees[1]
    assert agg_bit_1.support == fr.guaranteed_weight + 30 # v2(in) + v3

    agg_bit_3 = next(b for b in agg_bits if b.indexes == [3])
    assert agg_bit_3.fee == fees[3]
    assert agg_bit_3.support == fr.guaranteed_weight + 10 + 10 # v2(in) + v0 + v1

def test_filter_and_aggregate_pipeline(sample_votes_and_fees):
    """Tests the full filter_and_aggregate pipeline."""
    bit_votes, fees, total_weight = sample_votes_and_fees
    agg_votes, agg_bits, fr = filter_and_aggregate(bit_votes, fees, total_weight)

    assert len(agg_votes) == 2
    assert len(agg_bits) == 2
    assert fr.remaining_bits == {1, 3}

def test_assemble_solution():
    """Tests the final assembly of the consensus bit vote."""
    class MockFilterResults:
        always_in_bits = [0, 4]

    solution_bits = [
        AggregatedBit(fee=0, indexes=[1, 2], support=0, value=None),
        AggregatedBit(fee=0, indexes=[5], support=0, value=None),
    ]

    bv = assemble_solution(MockFilterResults(), solution_bits, 8)

    assert bv.length == 8
    expected_vector = (1 << 0) | (1 << 4) | (1 << 1) | (1 << 2) | (1 << 5)
    assert bv.bit_vector == expected_vector