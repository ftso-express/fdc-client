import asyncio
from dataclasses import dataclass, field
from enum import Enum
from typing import Any, List, Optional

class RoundStatus(Enum):
    UNASSIGNED = 0
    PRE_CONSENSUS = 1
    CONSENSUS = 2
    DONE = 3
    FAILED = 4

class AttestationStatus(Enum):
    UNPROCESSED = 0
    UNSUPPORTED_PAIR = 1
    WAITING = 2
    PROCESSING = 3
    SUCCESS = 4
    WRONG_MIC = 5
    INVALID_LUT = 6
    RETRYING = 7
    PROCESS_ERROR = 8
    UNCONFIRMED = 9

@dataclass
class IndexLog:
    """Represents the position of a log in the blockchain."""
    block_number: int
    log_index: int

@dataclass
class RoundStatusMutex:
    """A thread-safe container for the round status."""
    value: RoundStatus = RoundStatus.UNASSIGNED
    lock: asyncio.Lock = field(default_factory=asyncio.Lock)

@dataclass
class Attestation:
    """
    Represents a single attestation request and its lifecycle.
    This is the Python equivalent of the Go Attestation struct.
    """
    indexes: List[IndexLog]
    round_id: int
    round_status: RoundStatusMutex
    request: bytes  # Corresponds to Request type in Go (fdchub.FdcHubAttestationRequest)
    response: Optional[bytes] = None # Corresponds to Response type in Go
    fee: int = 0
    status: AttestationStatus = AttestationStatus.WAITING
    consensus: bool = False
    hash: Optional[bytes] = None # Corresponds to common.Hash

    # Placeholders for more complex types
    response_abi: Optional[Any] = None # Corresponds to *abi.Arguments
    response_abi_string: Optional[str] = None
    lut_limit: int = 0
    queue_name: Optional[str] = None
    credentials: Optional[Any] = None # Corresponds to *VerifierCredentials
    queue_pointer: Optional[Any] = None # Corresponds to *priority.Item

    def __post_init__(self):
        # The request in Go is a complex struct, here we just check it's bytes
        if not isinstance(self.request, bytes):
            raise TypeError("Request must be a bytes object.")

    def get_first_index(self) -> IndexLog:
        """Safely retrieve the first IndexLog for sorting purposes."""
        if not self.indexes:
            # This should not happen in a correctly functioning system
            raise IndexError("Attestation has no indexes.")
        return self.indexes[0]

def earlier_log(a: IndexLog, b: IndexLog) -> bool:
    """
    Returns True if log 'a' occurred before log 'b'.
    """
    if a.block_number < b.block_number:
        return True
    if a.block_number == b.block_number and a.log_index < b.log_index:
        return True
    return False