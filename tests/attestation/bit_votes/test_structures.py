import pytest
from src.attestation.bit_votes.structures import (
    BitVote,
    IndexTx,
    WeightedBitVote,
    earlier_tx,
)

def test_bitvote_initialization():
    """Tests the initialization of the BitVote dataclass."""
    bv = BitVote(length=10, bit_vector=0b10101)
    assert bv.length == 10
    assert bv.bit_vector == 21

def test_bitvote_invalid_length():
    """Tests that BitVote raises an error for an invalid length."""
    with pytest.raises(ValueError, match="Length must be a 16-bit unsigned integer."):
        BitVote(length=65536, bit_vector=1) # Too large
    with pytest.raises(ValueError, match="Length must be a 16-bit unsigned integer."):
        BitVote(length=-1, bit_vector=1) # Negative

def test_indextx_initialization():
    """Tests the initialization of the IndexTx dataclass."""
    itx = IndexTx(block_number=123, transaction_index=45)
    assert itx.block_number == 123
    assert itx.transaction_index == 45

def test_weightedbitvote_initialization():
    """Tests the initialization of the WeightedBitVote dataclass."""
    bv = BitVote(length=8, bit_vector=0b1100)
    itx = IndexTx(block_number=1, transaction_index=2)
    wbv = WeightedBitVote(index=5, index_tx=itx, weight=100, bit_vote=bv)

    assert wbv.index == 5
    assert wbv.weight == 100
    assert wbv.bit_vote.length == 8
    assert wbv.index_tx.block_number == 1

def test_earlier_tx_logic():
    """Tests the earlier_tx function with various cases."""
    tx1 = IndexTx(100, 1)
    tx2 = IndexTx(100, 2)
    tx3 = IndexTx(101, 0)

    assert earlier_tx(tx1, tx2)
    assert not earlier_tx(tx2, tx1)
    assert earlier_tx(tx1, tx3)
    assert not earlier_tx(tx3, tx1)
    assert earlier_tx(tx2, tx3)
    assert not earlier_tx(tx3, tx2)
    assert not earlier_tx(tx1, tx1)