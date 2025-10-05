import pytest
from src.attestation.structures import Attestation, AttestationStatus, IndexLog, RoundStatusMutex
from src.attestation.preparation import (
    prepare_request,
    AttestationTypeConfig,
    SourceConfig,
    VerifierCredentials,
)

# Mock data for testing
ATTESTATION_TYPE_PAYMENT = b'Pay '
SOURCE_BTC = b'BTC '

@pytest.fixture
def mock_configs():
    """Provides a mock configuration dictionary for tests."""
    return {
        ATTESTATION_TYPE_PAYMENT: AttestationTypeConfig(
            response_abi="mock_abi",
            response_abi_string="mock_abi_string",
            sources={
                SOURCE_BTC: SourceConfig(
                    lut_limit=1209600,
                    url="http://btc-verifier.com",
                    api_key="btc-key",
                    queue_name="btc_queue",
                )
            },
        )
    }

@pytest.fixture
def base_attestation():
    """Provides a base attestation object for tests."""
    request_data = ATTESTATION_TYPE_PAYMENT + SOURCE_BTC + b"some_payload"
    return Attestation(
        indexes=[IndexLog(1, 1)],
        round_id=1,
        round_status=RoundStatusMutex(),
        request=request_data,
    )

def test_prepare_request_success(base_attestation, mock_configs):
    """Tests the successful preparation of an attestation."""
    prepare_request(base_attestation, mock_configs)

    assert base_attestation.status == AttestationStatus.PROCESSING
    assert base_attestation.response_abi == "mock_abi"
    assert base_attestation.response_abi_string == "mock_abi_string"
    assert base_attestation.lut_limit == 1209600
    assert base_attestation.queue_name == "btc_queue"
    assert isinstance(base_attestation.credentials, VerifierCredentials)
    assert base_attestation.credentials.url == "http://btc-verifier.com"
    assert base_attestation.credentials.api_key == "btc-key"

def test_prepare_request_unsupported_type(base_attestation, mock_configs):
    """Tests failure when the attestation type is not in the config."""
    base_attestation.request = b'UNKN' + SOURCE_BTC + b"payload"

    with pytest.raises(ValueError, match="Unsupported attestation type"):
        prepare_request(base_attestation, mock_configs)

    assert base_attestation.status == AttestationStatus.UNSUPPORTED_PAIR

def test_prepare_request_unsupported_source(base_attestation, mock_configs):
    """Tests failure when the source is not in the config for a valid type."""
    base_attestation.request = ATTESTATION_TYPE_PAYMENT + b'ETH ' + b"payload"

    with pytest.raises(ValueError, match="Unsupported source"):
        prepare_request(base_attestation, mock_configs)

    assert base_attestation.status == AttestationStatus.UNSUPPORTED_PAIR

def test_prepare_request_invalid_format(base_attestation, mock_configs):
    """Tests failure when the request data is too short."""
    base_attestation.request = b'short'

    with pytest.raises(ValueError, match="Invalid request format"):
        prepare_request(base_attestation, mock_configs)

    assert base_attestation.status == AttestationStatus.PROCESS_ERROR