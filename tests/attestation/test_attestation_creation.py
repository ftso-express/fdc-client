import pytest
from src.attestation.creation import (
    attestation_from_database_log,
    DatabaseLog,
    parse_attestation_request_log,
)
from src.attestation.structures import AttestationStatus, RoundStatus
from src.timing import Timing

@pytest.fixture
def default_timing():
    """Provides a default Timing instance for tests."""
    return Timing({})

def test_parse_attestation_request_log():
    """Tests the mock log parsing function."""
    # Fee is 1, request data is 'hello'
    log_data = b'\x00\x01hello'
    db_log = DatabaseLog(timestamp=1, block_number=1, log_index=1, data=log_data)
    parsed = parse_attestation_request_log(db_log)
    assert parsed.fee == 1
    assert parsed.data == b'hello'

def test_attestation_from_database_log_success(default_timing: Timing):
    """
    Tests the successful creation of an Attestation object from a database log.
    """
    timestamp = default_timing.t0 + 100
    log_data = b'\x00\x0A' + b'test_request_data' # Fee 10
    db_log = DatabaseLog(
        timestamp=timestamp,
        block_number=200,
        log_index=5,
        data=log_data,
    )

    att = attestation_from_database_log(db_log, default_timing)

    expected_round_id = (timestamp - default_timing.t0) // default_timing.collect_duration_sec

    assert att.round_id == expected_round_id
    assert att.request == b'test_request_data'
    assert att.fee == 10
    assert len(att.indexes) == 1
    assert att.indexes[0].block_number == 200
    assert att.indexes[0].log_index == 5
    assert att.status == AttestationStatus.WAITING
    assert att.round_status.value == RoundStatus.UNASSIGNED

def test_attestation_from_database_log_invalid_timestamp(default_timing: Timing):
    """
    Tests that creation fails if the timestamp is before T0.
    """
    db_log = DatabaseLog(
        timestamp=default_timing.t0 - 1, # Invalid timestamp
        block_number=200,
        log_index=5,
        data=b'\x00\x0Atest_request_data',
    )

    with pytest.raises(ValueError, match="Failed to determine round ID"):
        attestation_from_database_log(db_log, default_timing)