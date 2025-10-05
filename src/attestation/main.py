from src.attestation.structures import Attestation, AttestationStatus
from src.attestation.verification import (
    compute_mic,
    get_mic as get_mic_from_request,
    get_lut,
    valid_lut,
    get_hash,
)
from src.timing import Timing
from src.logger import log

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
    log.info(f"Validating response for round {attestation.round_id}")

    if attestation.response is None:
        attestation.status = AttestationStatus.PROCESS_ERROR
        log.error("Attestation has no response to validate.")
        raise ValueError("Attestation has no response to validate.")

    if attestation.response_abi is None or attestation.response_abi_string is None:
        attestation.status = AttestationStatus.PROCESS_ERROR
        log.error("Attestation is not prepared with ABI information.")
        raise ValueError("Attestation is not prepared with ABI information.")

    # 1. Validate the Message Integrity Code (MIC)
    try:
        log.debug("Validating Message Integrity Code (MIC)...")
        mic_req = get_mic_from_request(attestation.request)
        mic_res = compute_mic(attestation.response, attestation.response_abi)
        log.debug(f"Request MIC: {mic_req}, Response MIC: {mic_res}")
    except Exception as e:
        attestation.status = AttestationStatus.PROCESS_ERROR
        log.error(f"Error during MIC validation: {e}")
        raise

    if mic_req != mic_res:
        attestation.status = AttestationStatus.WRONG_MIC
        log.warning("MIC validation failed: Mismatch between request and response.")
        return

    log.info("MIC validation successful.")

    # 2. Validate the Latest Used Timestamp (LUT)
    try:
        log.debug("Validating Latest Used Timestamp (LUT)...")
        lut = get_lut(attestation.response)
        round_start = timing.choose_start_timestamp(attestation.round_id)
        log.debug(f"LUT: {lut}, LUT Limit: {attestation.lut_limit}, Round Start: {round_start}")
    except Exception as e:
        attestation.status = AttestationStatus.PROCESS_ERROR
        log.error(f"Error during LUT validation setup: {e}")
        raise

    if not valid_lut(lut, attestation.lut_limit, round_start):
        attestation.status = AttestationStatus.INVALID_LUT
        log.warning("LUT validation failed: Invalid timestamp.")
        return

    log.info("LUT validation successful.")

    # 3. All checks passed, compute final hash and set status to Success
    try:
        log.debug("Computing final hash...")
        attestation.hash = get_hash(attestation.response, attestation.round_id)
        attestation.status = AttestationStatus.SUCCESS
        log.info(f"Response validation successful. Final hash: {attestation.hash}")
    except Exception as e:
        attestation.status = AttestationStatus.PROCESS_ERROR
        log.error(f"Error computing final hash: {e}")
        raise