from typing import NewType, List, Tuple
from eth_abi import encode as abi_encode, decode as abi_decode
from Crypto.Hash import keccak
from src.logger import log

# Type aliases to mirror the Go implementation
Request = NewType("Request", bytes)
Response = NewType("Response", bytes)

def get_attestation_type(r: Request) -> bytes:
    """Extracts the attestation type from the request (first 32 bytes)."""
    if len(r) < 96:
        log.error("Request is too short to extract attestation type.")
        raise ValueError("Request is too short")
    return r[0:32]

def get_source(r: Request) -> bytes:
    """Extracts the source from the request (second 32 bytes)."""
    if len(r) < 96:
        log.error("Request is too short to extract source.")
        raise ValueError("Request is too short")
    return r[32:64]

def get_mic(r: Request) -> bytes:
    """Extracts the Message Integrity Code (MIC) from the request (third 32 bytes)."""
    if len(r) < 96:
        log.error("Request is too short to extract MIC.")
        raise ValueError("Request is too short")
    return r[64:96]

def valid_lut(lut: int, lut_limit: int, round_start: int) -> bool:
    """
    Safely checks whether round_start - lut < lut_limit.
    This function uses integer arithmetic and is a direct port of the Go logic.
    """
    log.debug(f"Validating LUT: lut={lut}, lut_limit={lut_limit}, round_start={round_start}")
    if lut > round_start:
        log.debug("LUT is in the future, invalid.")
        return False

    lhs = round_start - lut
    is_valid = lhs < lut_limit
    log.debug(f"round_start - lut = {lhs}. Is valid: {is_valid}")
    return is_valid

def is_static_type(encoded_bytes: bytes) -> bool:
    """
    Checks if ABI-encoded bytes represent a static or dynamic type.
    If the first 32 bytes represent the offset to the data (i.e., 32),
    it's a dynamic type. Otherwise, it's static.
    """
    if len(encoded_bytes) < 32:
        log.error("Encoded bytes are too short to determine type.")
        raise ValueError("Encoded bytes are too short to determine type.")

    offset = int.from_bytes(encoded_bytes[0:32], 'big')
    is_dynamic = offset == 32
    log.debug(f"Offset is {offset}. Is dynamic type: {is_dynamic}")
    return not is_dynamic

def get_lut(response: Response) -> int:
    """
    Returns the LUT (Latest Used Timestamp) from the response.
    LUT is assumed to be a uint64 in the fourth 32-byte slot.
    """
    static = is_static_type(response)
    log.debug(f"Response is a static type: {static}")

    lut_start = 3 * 32
    lut_end = 4 * 32

    if not static:
        lut_start += 32
        lut_end += 32

    if len(response) < lut_end:
        log.error("Response is too short to contain LUT.")
        raise ValueError("Response is too short to contain LUT")

    lut_bytes = response[lut_start:lut_end]
    lut = int.from_bytes(lut_bytes, 'big')
    log.debug(f"Extracted LUT: {lut}")
    return lut

def add_round(response: Response, round_id: int) -> Response:
    """
    Sets the round ID in the response (third 32-byte slot).
    Returns a new Response bytes object with the round ID inserted.
    """
    static = is_static_type(response)
    log.debug(f"Response is a static type: {static}")

    round_id_start = 2 * 32
    round_id_end = 3 * 32

    if not static:
        round_id_start += 32
        round_id_end += 32

    if len(response) < round_id_end:
        log.error("Response is too short to add round ID.")
        raise ValueError("Response is too short to add round ID")

    round_id_bytes = round_id.to_bytes(32, 'big')
    log.debug(f"Inserting round ID {round_id} into response.")

    new_response_bytes = (
        response[:round_id_start] +
        round_id_bytes +
        response[round_id_end:]
    )

    return Response(new_response_bytes)

def get_hash(response: Response, round_id: int) -> bytes:
    """
    Computes the final hash of the response after adding the round ID.
    """
    log.debug(f"Computing hash for response with round ID {round_id}")
    response_with_round = add_round(response, round_id)

    hasher = keccak.new(digest_bits=256)
    hasher.update(response_with_round)
    final_hash = hasher.digest()
    log.debug(f"Computed final hash: {final_hash.hex()}")
    return final_hash

def compute_mic(response: Response, response_abi_types: List[str]) -> bytes:
    """
    Computes the Message Integrity Code (MIC) from the response.
    MIC is keccak256(abi.encode(abi.decode(response), "Flare")).
    """
    log.debug(f"Computing MIC for response with ABI types: {response_abi_types}")
    decoded_data = abi_decode(response_abi_types, response)
    log.debug(f"Decoded data: {decoded_data}")

    types_with_salt = response_abi_types + ['string']
    data_with_salt = list(decoded_data) + ["Flare"]
    log.debug(f"Data with salt: {data_with_salt}")

    encoded_with_salt = abi_encode(types_with_salt, data_with_salt)

    hasher = keccak.new(digest_bits=256)
    hasher.update(encoded_with_salt)
    mic = hasher.digest()
    log.debug(f"Computed MIC: {mic.hex()}")
    return mic