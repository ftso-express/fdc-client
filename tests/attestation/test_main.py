import pytest
from eth_abi import encode as abi_encode

from src.attestation.structures import Attestation, AttestationStatus, IndexLog, RoundStatusMutex
from src.attestation.main import validate_response
from src.attestation.verification import compute_mic, get_mic, Request, Response
from src.timing import Timing

# --- Test Fixtures ---

@pytest.fixture
def default_timing():
    """Provides a default Timing instance."""
    return Timing({})

@pytest.fixture
def sample_abi_types():
    """Provides a sample list of ABI types for a response."""
    return ['uint256', 'bytes32', 'uint32', 'uint64']

@pytest.fixture
def valid_attestation(default_timing, sample_abi_types):
    """Creates a complete and valid Attestation object for testing."""

    # Create a valid response
    response_values = [1, b'\xAA'*32, 0, default_timing.t0]
    response = Response(abi_encode(sample_abi_types, response_values))

    # Create a request with the correct MIC for that response
    mic = compute_mic(response, sample_abi_types)
    request_data = b'\x01'*32 + b'\x02'*32 + mic

    att = Attestation(
        indexes=[IndexLog(1, 1)],
        round_id=0,
        round_status=RoundStatusMutex(),
        request=Request(request_data),
        response=response,
        response_abi=sample_abi_types,
        response_abi_string="mock",
        lut_limit=1000,
    )
    return att

# --- Test Cases ---

def test_validate_response_success(valid_attestation, default_timing):
    """Tests the full validation process for a successful case."""
    validate_response(valid_attestation, default_timing)

    assert valid_attestation.status == AttestationStatus.SUCCESS
    assert valid_attestation.hash is not None
    assert isinstance(valid_attestation.hash, bytes)
    assert len(valid_attestation.hash) == 32

def test_validate_response_wrong_mic(valid_attestation, default_timing):
    """Tests that validation fails if the MIC does not match."""
    # Tamper with the response so the MIC will be incorrect
    valid_attestation.response = Response(valid_attestation.response[:-1] + b'\x00')

    validate_response(valid_attestation, default_timing)

    assert valid_attestation.status == AttestationStatus.WRONG_MIC
    assert valid_attestation.hash is None

def test_validate_response_invalid_lut(valid_attestation, default_timing, sample_abi_types):
    """Tests that validation fails if the LUT is outside the valid window."""
    # Create a response with an old LUT
    lut_too_old = default_timing.t0 - 2000
    response_values = [1, b'\xAA'*32, 0, lut_too_old]
    response = Response(abi_encode(sample_abi_types, response_values))

    # Update the attestation with the new response and a matching MIC
    mic = compute_mic(response, sample_abi_types)
    valid_attestation.request = Request(b'\x01'*32 + b'\x02'*32 + mic)
    valid_attestation.response = response

    validate_response(valid_attestation, default_timing)

    assert valid_attestation.status == AttestationStatus.INVALID_LUT
    assert valid_attestation.hash is None

def test_validate_response_missing_response(valid_attestation, default_timing):
    """Tests that validation raises an error if the response is missing."""
    valid_attestation.response = None
    with pytest.raises(ValueError, match="Attestation has no response to validate."):
        validate_response(valid_attestation, default_timing)
    assert valid_attestation.status == AttestationStatus.PROCESS_ERROR

def test_validate_response_not_prepared(valid_attestation, default_timing):
    """Tests that validation raises an error if the ABI info is missing."""
    valid_attestation.response_abi = None
    with pytest.raises(ValueError, match="Attestation is not prepared with ABI information."):
        validate_response(valid_attestation, default_timing)
    assert valid_attestation.status == AttestationStatus.PROCESS_ERROR