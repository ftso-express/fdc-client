from dataclasses import dataclass
from typing import Dict, Any

from src.attestation.structures import Attestation, AttestationStatus

# Mock configuration structures to mirror the Go implementation's config.
# In a real application, these would be populated from the main Config object.

@dataclass
class VerifierCredentials:
    """Represents the credentials needed to connect to a verifier."""
    url: str
    api_key: str

@dataclass
class SourceConfig:
    """Represents the configuration for a specific source of an attestation type."""
    lut_limit: int
    url: str
    api_key: str
    queue_name: str

@dataclass
class AttestationTypeConfig:
    """Represents the configuration for a specific attestation type."""
    response_abi: Any  # Placeholder for a proper ABI object
    response_abi_string: str
    sources: Dict[bytes, SourceConfig]  # Keyed by source ID

# Mock functions to stand in for Go's `request.AttestationType()` and `request.Source()`.
# In a full implementation, this logic would be part of a dedicated Request class.
def get_attestation_type_from_request(request: bytes) -> bytes:
    """Mock function to extract the attestation type from a request payload."""
    if len(request) < 4:
        raise IndexError("Request data is too short to extract attestation type.")
    return request[:4]

def get_source_from_request(request: bytes) -> bytes:
    """Mock function to extract the source from a request payload."""
    if len(request) < 8:
        raise IndexError("Request data is too short to extract source.")
    return request[4:8]

def prepare_request(
    attestation: Attestation,
    attestation_types_configs: Dict[bytes, AttestationTypeConfig]
) -> None:
    """
    Prepares an attestation for verification by adding the response ABI, LUT limit,
    and verifier credentials based on the application's configuration.

    This function modifies the attestation object in place and updates its status.
    It raises a ValueError if the configuration for the given attestation
    type or source is not found.
    """
    try:
        attestation_type = get_attestation_type_from_request(attestation.request)
        source = get_source_from_request(attestation.request)
    except IndexError:
        attestation.status = AttestationStatus.PROCESS_ERROR
        raise ValueError("Invalid request format: cannot extract type or source.")

    # Find the configuration for the attestation type
    type_config = attestation_types_configs.get(attestation_type)
    if not type_config:
        attestation.status = AttestationStatus.UNSUPPORTED_PAIR
        raise ValueError(f"Unsupported attestation type: {attestation_type.hex()}")

    # Find the configuration for the source within that type
    source_config = type_config.sources.get(source)
    if not source_config:
        attestation.status = AttestationStatus.UNSUPPORTED_PAIR
        raise ValueError(f"Unsupported source '{source.hex()}' for type '{attestation_type.hex()}'")

    # Populate the attestation object with the retrieved configuration
    attestation.response_abi = type_config.response_abi
    attestation.response_abi_string = type_config.response_abi_string
    attestation.lut_limit = source_config.lut_limit
    attestation.credentials = VerifierCredentials(url=source_config.url, api_key=source_config.api_key)
    attestation.queue_name = source_config.queue_name

    # Update the status to indicate it's ready for processing
    attestation.status = AttestationStatus.PROCESSING