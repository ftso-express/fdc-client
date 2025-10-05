import pytest
from eth_abi import encode as abi_encode, decode as abi_decode

from src.attestation.verification import (
    get_attestation_type,
    get_source,
    get_mic,
    valid_lut,
    Request,
    Response,
    add_round,
    compute_mic,
    get_hash,
    get_lut,
    is_static_type,
)

@pytest.fixture
def sample_request() -> Request:
    """Provides a sample 96-byte request for testing."""
    att_type = b'\x01' * 32
    source = b'\x02' * 32
    mic = b'\x03' * 32
    return Request(att_type + source + mic)

@pytest.fixture
def static_response() -> Response:
    """ABI-encoded response for a static type: (uint256, bytes32, uint32, uint64)"""
    types = ['uint256', 'bytes32', 'uint32', 'uint64']
    values = [1, b'\xAA'*32, 0, 12345] # roundID is 0 initially
    return Response(abi_encode(types, values))

@pytest.fixture
def dynamic_response() -> Response:
    """
    ABI-encoded response for a dynamic type: (string, bytes32, uint32, uint64).
    This is encoded as a single struct type.
    """
    struct_type = '(string,bytes32,uint32,uint64)'
    struct_value = ("hello world", b'\xBB'*32, 0, 54321) # roundID is 0 initially
    return Response(abi_encode([struct_type], [struct_value]))

def test_get_attestation_type(sample_request: Request):
    assert get_attestation_type(sample_request) == b'\x01' * 32

def test_get_source(sample_request: Request):
    assert get_source(sample_request) == b'\x02' * 32

def test_get_mic(sample_request: Request):
    assert get_mic(sample_request) == b'\x03' * 32

def test_request_parsing_too_short():
    short_request = Request(b'\x00' * 95)
    with pytest.raises(ValueError, match="Request is too short"):
        get_attestation_type(short_request)

def test_valid_lut():
    round_start = 1000
    lut_limit = 100
    assert valid_lut(lut=950, lut_limit=lut_limit, round_start=round_start)
    assert not valid_lut(lut=900, lut_limit=lut_limit, round_start=round_start)

def test_is_static_type(static_response, dynamic_response):
    assert is_static_type(static_response)
    assert not is_static_type(dynamic_response)

def test_get_lut(static_response, dynamic_response):
    assert get_lut(static_response) == 12345
    assert get_lut(dynamic_response) == 54321

def test_add_round(static_response):
    new_round_id = 999
    response_with_round = add_round(static_response, new_round_id)
    decoded = abi_decode(['uint256', 'bytes32', 'uint32', 'uint64'], response_with_round)
    assert decoded[2] == new_round_id

def test_get_hash(static_response):
    h = get_hash(static_response, 1)
    assert isinstance(h, bytes) and len(h) == 32
    h2 = get_hash(static_response, 2)
    assert h != h2

def test_compute_mic(static_response):
    abi_types = ['uint256', 'bytes32', 'uint32', 'uint64']
    mic = compute_mic(static_response, abi_types)
    assert isinstance(mic, bytes) and len(mic) == 32

    types = ['uint256', 'bytes32', 'uint32', 'uint64']
    values = [2, b'\xAA'*32, 0, 12345]
    different_response = Response(abi_encode(types, values))
    mic2 = compute_mic(different_response, abi_types)
    assert mic != mic2