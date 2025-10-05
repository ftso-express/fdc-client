from dataclasses import dataclass, field
from typing import Dict, List, Set

@dataclass
class BitVote:
    """
    Represents a bit vector indicating successful attestations.
    A Python port of the Go BitVote struct.
    """
    length: int  # Corresponds to uint16 in Go
    bit_vector: int  # Corresponds to *big.Int in Go

    def __post_init__(self):
        if not 0 <= self.length <= 65535:
            raise ValueError("Length must be a 16-bit unsigned integer.")

@dataclass
class IndexTx:
    """
    Represents the position of a transaction in the blockchain.
    """
    block_number: int
    transaction_index: int

@dataclass
class WeightedBitVote:
    """
    Represents a BitVote submitted by a specific voter with a given weight.
    """
    index: int  # Signing policy index of the voter
    index_tx: IndexTx
    weight: int  # Corresponds to uint16 in Go
    bit_vote: BitVote

from typing import Optional

def earlier_tx(a: IndexTx, b: IndexTx) -> bool:
    """
    Compares two IndexTx objects. Returns True if 'a' occurred before 'b'.
    """
    if a.block_number < b.block_number:
        return True
    if a.block_number == b.block_number and a.transaction_index < b.transaction_index:
        return True
    return False

@dataclass
class Value:
    """
    Represents the value of a consensus solution, typically composed of
    a primary (capped) and secondary (uncapped) value for tie-breaking.
    """
    capped_value: int
    uncapped_value: int

    def copy(self) -> 'Value':
        """Returns a copy of the Value object."""
        return Value(self.capped_value, self.uncapped_value)

    def cmp(self, other: 'Value') -> int:
        """
        Compares this Value with another.
        Returns 1 if self > other, -1 if self < other, 0 if equal.
        """
        if self.capped_value > other.capped_value:
            return 1
        if self.capped_value < other.capped_value:
            return -1
        # Capped values are equal, compare uncapped for tie-breaking
        if self.uncapped_value > other.uncapped_value:
            return 1
        if self.uncapped_value < other.uncapped_value:
            return -1
        return 0

@dataclass
class FilterResults:
    """
    Holds the results of the initial filtering and preprocessing step.
    """
    always_in_bits: List[int] = field(default_factory=list)
    always_out_bits: List[int] = field(default_factory=list)
    remaining_bits: Set[int] = field(default_factory=set)
    guaranteed_fees: int = 0

    always_in_votes: List[int] = field(default_factory=list)
    always_out_votes: List[int] = field(default_factory=list)
    remaining_votes: Set[int] = field(default_factory=set)
    guaranteed_weight: int = 0
    remaining_weight: int = 0

@dataclass
class AggregatedBit:
    """
    Represents a set of bits that have the same support pattern across
    the remaining voters, aggregated into a single unit.
    """
    fee: int
    indexes: List[int]
    support: int
    value: Optional[Value] = None # Cached value

@dataclass
class AggregatedVote:
    """
    Represents a set of voters that have the same bit vector across the
    remaining bits, aggregated into a single unit.
    """
    bit_vector: int
    weight: int
    indexes: List[int]
    fees: int  # Sum of fees of the bits this vote supports, for sorting