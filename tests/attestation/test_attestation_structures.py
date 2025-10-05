import pytest
from src.attestation.structures import (
    Attestation,
    AttestationStatus,
    IndexLog,
    RoundStatus,
    RoundStatusMutex,
    earlier_log,
)

def test_attestation_initialization():
    """Tests the basic initialization of the Attestation dataclass."""
    index = IndexLog(block_number=100, log_index=1)
    round_status = RoundStatusMutex()
    att = Attestation(
        indexes=[index],
        round_id=123,
        round_status=round_status,
        request=b'\xde\xad\xbe\xef',
    )
    assert att.round_id == 123
    assert att.status == AttestationStatus.WAITING
    assert att.request == b'\xde\xad\xbe\xef'
    assert att.get_first_index() == index

def test_attestation_invalid_request_type():
    """Tests that a TypeError is raised for an invalid request type."""
    with pytest.raises(TypeError, match="Request must be a bytes object."):
        Attestation(
            indexes=[IndexLog(1,1)],
            round_id=1,
            round_status=RoundStatusMutex(),
            request="this is not bytes"
        )

def test_earlier_log_logic():
    """Tests the earlier_log function with various cases."""
    log1 = IndexLog(100, 1)
    log2 = IndexLog(100, 2)
    log3 = IndexLog(101, 0)

    assert earlier_log(log1, log2)
    assert not earlier_log(log2, log1)
    assert earlier_log(log1, log3)
    assert not earlier_log(log3, log1)
    assert earlier_log(log2, log3)
    assert not earlier_log(log3, log2)
    assert not earlier_log(log1, log1)

def test_round_status_mutex():
    """Tests the RoundStatusMutex dataclass."""
    rs_mutex = RoundStatusMutex()
    assert rs_mutex.value == RoundStatus.UNASSIGNED
    rs_mutex.value = RoundStatus.CONSENSUS
    assert rs_mutex.value == RoundStatus.CONSENSUS

def test_get_first_index_empty():
    """Tests that get_first_index raises an error if indexes are empty."""
    att = Attestation([], 1, RoundStatusMutex(), b'')
    with pytest.raises(IndexError, match="Attestation has no indexes."):
        att.get_first_index()