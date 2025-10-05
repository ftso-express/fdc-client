import pytest
from src.attestation.structures import Attestation, AttestationStatus, IndexLog, RoundStatusMutex
from src.attestation.bit_votes.creation import bit_vote_from_attestations

# Helper to create a mock Attestation with a specific status
def create_mock_attestation(status: AttestationStatus) -> Attestation:
    att = Attestation(
        indexes=[IndexLog(1,1)],
        round_id=1,
        round_status=RoundStatusMutex(),
        request=b''
    )
    att.status = status
    return att

def test_bit_vote_from_attestations_mixed():
    """Tests creating a bit vote from a mixed list of attestations."""
    attestations = [
        create_mock_attestation(AttestationStatus.SUCCESS),      # bit 0
        create_mock_attestation(AttestationStatus.WRONG_MIC),     # bit 1 (fail)
        create_mock_attestation(AttestationStatus.SUCCESS),      # bit 2
        create_mock_attestation(AttestationStatus.INVALID_LUT),   # bit 3 (fail)
        create_mock_attestation(AttestationStatus.SUCCESS),      # bit 4
    ]

    bv = bit_vote_from_attestations(attestations)

    assert bv.length == 5
    # Expected bit vector: 10101 in binary
    expected_vector = (1 << 0) | (1 << 2) | (1 << 4)
    assert bv.bit_vector == expected_vector

def test_bit_vote_from_attestations_all_success():
    """Tests the case where all attestations are successful."""
    attestations = [create_mock_attestation(AttestationStatus.SUCCESS) for _ in range(5)]

    bv = bit_vote_from_attestations(attestations)

    assert bv.length == 5
    # Expected: 11111 in binary
    assert bv.bit_vector == 0b11111

def test_bit_vote_from_attestations_none_success():
    """Tests the case where no attestations are successful."""
    attestations = [create_mock_attestation(AttestationStatus.PROCESS_ERROR) for _ in range(5)]

    bv = bit_vote_from_attestations(attestations)

    assert bv.length == 5
    assert bv.bit_vector == 0

def test_bit_vote_from_attestations_empty_list():
    """Tests that an empty list of attestations produces an empty bit vote."""
    bv = bit_vote_from_attestations([])

    assert bv.length == 0
    assert bv.bit_vector == 0

def test_bit_vote_from_attestations_too_many():
    """Tests that the function raises an error if there are too many attestations."""
    # Create a list-like object that reports a large length to avoid creating a huge list
    class MockLargeList:
        def __len__(self):
            return 65536

    with pytest.raises(ValueError, match="Cannot create a bit vote for more than 65535 attestations."):
        bit_vote_from_attestations(MockLargeList())