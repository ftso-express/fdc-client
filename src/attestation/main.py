from src.attestation.structures import Attestation, AttestationStatus
from src.attestation.verification import (
    compute_mic,
    get_mic as get_mic_from_request,
    get_lut,
    valid_lut,
    get_hash,
)
from src.timing import Timing

def validate_response(attestation: Attestation, timing: Timing) -> None:
    """
    Validates the response within an Attestation object.

    This function orchestrates the validation process by checking the Message
    Integrity Code (MIC) and the Latest Used Timestamp (LUT). If both are
    valid, it computes the final hash and updates the attestation's status
    to SUCCESS. Otherwise, it sets the appropriate error status.

    Args:
        attestation: The Attestation object to validate. It is modified in place.
        timing: The Timing configuration object needed for LUT validation.
    """
    if attestation.response is None:
        attestation.status = AttestationStatus.PROCESS_ERROR
        raise ValueError("Attestation has no response to validate.")

    if attestation.response_abi is None or attestation.response_abi_string is None:
        attestation.status = AttestationStatus.PROCESS_ERROR
        raise ValueError("Attestation is not prepared with ABI information.")

    # 1. Validate the Message Integrity Code (MIC)
    try:
        mic_req = get_mic_from_request(attestation.request)
        # Note: The Go code assumes a list of types, but eth-abi needs the raw string.
        # We will need to adjust how we store/pass the ABI info. For now, we assume
        # the `response_abi` attribute holds the list of type strings.
        mic_res = compute_mic(attestation.response, attestation.response_abi)
    except Exception:
        attestation.status = AttestationStatus.PROCESS_ERROR
        raise

    if mic_req != mic_res:
        attestation.status = AttestationStatus.WRONG_MIC
        return

    # 2. Validate the Latest Used Timestamp (LUT)
    try:
        lut = get_lut(attestation.response)
        round_start = timing.choose_start_timestamp(attestation.round_id)
    except Exception:
        attestation.status = AttestationStatus.PROCESS_ERROR
        raise

    if not valid_lut(lut, attestation.lut_limit, round_start):
        attestation.status = AttestationStatus.INVALID_LUT
        return

    # 3. All checks passed, compute final hash and set status to Success
    try:
        attestation.hash = get_hash(attestation.response, attestation.round_id)
        attestation.status = AttestationStatus.SUCCESS
    except Exception:
        attestation.status = AttestationStatus.PROCESS_ERROR
        raise