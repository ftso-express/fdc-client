from dataclasses import dataclass

from src.attestation.structures import (
    Attestation,
    IndexLog,
    RoundStatus,
    RoundStatusMutex,
)
from src.timing import Timing

@dataclass
class DatabaseLog:
    """A mock dataclass to represent a log entry from the database."""
    timestamp: int
    block_number: int
    log_index: int
    data: bytes  # Raw log data

@dataclass
class ParsedAttestationRequest:
    """
    A mock dataclass to represent the parsed attestation request data.
    This stands in for Go's fdchub.FdcHubAttestationRequest.
    """
    data: bytes  # The core request payload
    fee: int

def parse_attestation_request_log(db_log: DatabaseLog) -> ParsedAttestationRequest:
    """
    A mock function to simulate parsing an attestation request log.
    In a real implementation, this would involve ABI decoding.
    For now, it extracts a mock fee and returns the request data.
    """
    # Simulate extracting a fee from the log data for demonstration
    mock_fee = int.from_bytes(db_log.data[:2], 'big') if len(db_log.data) >= 2 else 0
    # The actual request data would be the rest of the payload
    request_data = db_log.data[2:]

    return ParsedAttestationRequest(data=request_data, fee=mock_fee)

def attestation_from_database_log(db_log: DatabaseLog, timing: Timing) -> Attestation:
    """
    Creates an Attestation from a database log entry.
    This is a Python port of the Go function AttestationFromDatabaseLog.
    """
    # 1. Parse the log
    try:
        parsed_request = parse_attestation_request_log(db_log)
    except Exception as e:
        raise ValueError(f"Failed to parse log: {e}") from e

    # 2. Determine the round ID
    try:
        round_id = timing.round_id_for_timestamp(db_log.timestamp)
    except ValueError as e:
        raise ValueError(f"Failed to determine round ID for log: {e}") from e

    # 3. Create the Attestation object
    indexes = [IndexLog(block_number=db_log.block_number, log_index=db_log.log_index)]

    # Each attestation gets its own round status object
    round_status = RoundStatusMutex(value=RoundStatus.UNASSIGNED)

    attestation = Attestation(
        indexes=indexes,
        round_id=round_id,
        request=parsed_request.data,
        fee=parsed_request.fee,
        round_status=round_status,
        # Default status is WAITING
    )

    return attestation