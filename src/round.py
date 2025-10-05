import hashlib
from typing import List, Dict, Optional, Any

from src.attestation.structures import Attestation, AttestationStatus, RoundStatusMutex, RoundStatus
from src.attestation.bit_votes.structures import BitVote, WeightedBitVote, IndexTx, Value
from src.attestation.bit_votes.coding import decode_bit_vote_bytes
from src.attestation.bit_votes.creation import bit_vote_from_attestations
from src.attestation.bit_votes.consensus import branch_and_bound_bits_double
from src.attestation.bit_votes.preprocessing import filter_and_aggregate, assemble_solution
from src.attestation.verification import Request

# --- Mock Objects for External Dependencies ---

class MockVoter:
    """A mock to stand in for the complex VoterData object."""
    def __init__(self, index: int, weight: int):
        self.index = index
        self.weight = weight

class MockVoterSet:
    """A mock to stand in for the complex VoterSet object."""
    def __init__(self, voters: Dict[str, MockVoter], total_weight: int):
        # Maps signing address to voter data
        self.voter_data_map = voters
        self.total_weight = total_weight
        # Maps submit address to signing address
        self.submit_to_signing_address = {f"submit_{k}": k for k in voters.keys()}

# --- Round Class ---

class Round:
    """
    Represents all data and logic for a single voting round.
    This is a Python port of the `Round` struct and its methods from `round.go`.
    """
    def __init__(self, round_id: int, voter_set: MockVoterSet):
        self.id = round_id
        self.status = RoundStatusMutex(value=RoundStatus.PRE_CONSENSUS)
        self.voter_set = voter_set

        self.attestations: List[Attestation] = []
        self.attestation_map: Dict[bytes, Attestation] = {}

        self.bit_votes: List[WeightedBitVote] = []
        self.bit_vote_checklist: Dict[str, WeightedBitVote] = {} # Keyed by submit address

        self.consensus_calculation_finished = False
        self.consensus_bit_vote: Optional[BitVote] = None
        self.merkle_tree: Optional[Any] = None # Placeholder for a Merkle tree object

    def add_attestation(self, attestation: Attestation) -> bool:
        """
        Adds an attestation to the round, handling duplicates.
        Returns True if the attestation was new, False if it was a duplicate.
        """
        # Using a hash of the request as a unique identifier
        identifier = hashlib.sha256(attestation.request).digest()

        if identifier in self.attestation_map:
            existing_att = self.attestation_map[identifier]
            existing_att.fee += attestation.fee
            # This logic for prepending/appending indexes can be simplified
            # for the port, as long as sorting happens before use.
            existing_att.indexes.extend(attestation.indexes)
            return False

        self.attestation_map[identifier] = attestation
        self.attestations.append(attestation)
        attestation.round_status = self.status
        return True

    def _sort_attestations(self):
        """Sorts attestations by their index log (block number, then log index)."""
        self.attestations.sort(key=lambda a: (a.get_first_index().block_number, a.get_first_index().log_index))

    def get_local_bit_vote(self) -> BitVote:
        """Generates the bit vote from this node's perspective."""
        self._sort_attestations()
        return bit_vote_from_attestations(self.attestations)

    def process_incoming_bit_vote(self, message: Dict) -> None:
        """
        Processes a bit vote message from another participant.
        A port of `ProcessBitVote`.
        """
        # In Go, message is a complex struct. Here we use a simple dict.
        submit_address = message["from"]

        # Look up voter info
        signing_address = self.voter_set.submit_to_signing_address.get(submit_address)
        if not signing_address:
            raise ValueError(f"No signing address for submit address {submit_address}")

        voter = self.voter_set.voter_data_map.get(signing_address)
        if not voter or voter.weight <= 0:
            raise ValueError(f"Invalid or zero-weight voter for address {signing_address}")

        # Decode the bit vote from the payload
        bit_vote = decode_bit_vote_bytes(message["payload"])
        if bit_vote.length != len(self.attestations):
            raise ValueError("Bit vote length mismatch")

        # Create the weighted bit vote
        new_wbv = WeightedBitVote(
            index=voter.index,
            weight=voter.weight,
            bit_vote=bit_vote,
            index_tx=IndexTx(message["block_number"], message["tx_index"])
        )

        # Handle re-submissions
        if submit_address not in self.bit_vote_checklist:
            self.bit_votes.append(new_wbv)
            self.bit_vote_checklist[submit_address] = new_wbv
        else:
            # Overwrite if the new one is later
            existing_wbv = self.bit_vote_checklist[submit_address]
            if not (existing_wbv.index_tx.block_number > new_wbv.index_tx.block_number or \
               (existing_wbv.index_tx.block_number == new_wbv.index_tx.block_number and \
                existing_wbv.index_tx.transaction_index > new_wbv.index_tx.transaction_index)):

                # Find and replace in the main list
                for i, bv in enumerate(self.bit_votes):
                    if bv.index == existing_wbv.index:
                        self.bit_votes[i] = new_wbv
                        break
                self.bit_vote_checklist[submit_address] = new_wbv

    def compute_consensus(self, max_ops: int = 20_000_000):
        """
        Runs the full consensus algorithm on the collected bit votes.
        A port of `ComputeConsensusBitVote`.
        """
        self._sort_attestations()

        fees = [a.fee for a in self.attestations]

        # 1. Preprocessing
        agg_votes, agg_bits, fr = filter_and_aggregate(
            self.bit_votes, fees, self.voter_set.total_weight
        )

        # 2. Branch and Bound
        solution = branch_and_bound_bits_double(
            agg_votes, agg_bits, fr.guaranteed_weight,
            self.voter_set.total_weight, fr.guaranteed_fees, max_ops, Value(0,0)
        )

        # 3. Assemble the final bit vote
        self.consensus_bit_vote = assemble_solution(fr, solution.bits, len(self.attestations))
        self.consensus_calculation_finished = True

        # 4. Update the status of all attestations
        self.status.value = RoundStatus.CONSENSUS
        self._set_consensus_status()

    def _set_consensus_status(self):
        """Updates the 'consensus' flag on each attestation based on the result."""
        if self.consensus_bit_vote is None:
            return

        for i, att in enumerate(self.attestations):
            if (self.consensus_bit_vote.bit_vector >> i) & 1:
                att.consensus = True
            else:
                att.consensus = False

    def get_merkle_root(self) -> Optional[bytes]:
        """
        Builds a Merkle tree from confirmed attestations and returns the root.
        This is a simplified version of the Go implementation.
        """
        if self.status.value != RoundStatus.CONSENSUS:
            raise RuntimeError("Cannot build Merkle tree before consensus is reached.")

        hashes = []
        for att in self.attestations:
            if att.consensus:
                if att.status != AttestationStatus.SUCCESS or att.hash is None:
                    # In a real system, this is a critical error.
                    raise ValueError(f"Attestation in consensus but not confirmed: {att}")
                hashes.append(att.hash)

        if not hashes:
            self.merkle_tree = b'\x00' * 32 # Placeholder
            return self.merkle_tree

        # A real implementation would use a proper Merkle library.
        # For this test, we'll just hash the concatenation of hashes.
        import hashlib
        m = hashlib.sha256()
        for h in sorted(hashes): # Sort for determinism
            m.update(h)

        self.merkle_tree = m.digest()
        self.status.value = RoundStatus.DONE
        return self.merkle_tree