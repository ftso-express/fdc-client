import pytest
from src.round import Round, MockVoter, MockVoterSet
from src.attestation.structures import (
    Attestation,
    AttestationStatus,
    IndexLog,
    RoundStatusMutex,
    RoundStatus,
)
from src.attestation.bit_votes.structures import BitVote, WeightedBitVote, IndexTx

# --- Test Fixtures ---

@pytest.fixture
def mock_voter_set():
    """Creates a mock voter set for testing."""
    voters = {
        "signer1": MockVoter(index=0, weight=40),
        "signer2": MockVoter(index=1, weight=30),
        "signer3": MockVoter(index=2, weight=30),
    }
    return MockVoterSet(voters=voters, total_weight=100)

@pytest.fixture
def initial_round(mock_voter_set):
    """Creates a new Round and populates it with attestations."""
    r = Round(round_id=123, voter_set=mock_voter_set)

    # Attestation 0 (Success)
    att0 = Attestation(indexes=[IndexLog(1,0)], round_id=123, round_status=r.status, request=b'mic0', fee=100)
    att0.status = AttestationStatus.SUCCESS
    att0.hash = b'\x00'*32

    # Attestation 1 (Success)
    att1 = Attestation(indexes=[IndexLog(1,1)], round_id=123, round_status=r.status, request=b'mic1', fee=200)
    att1.status = AttestationStatus.SUCCESS
    att1.hash = b'\x01'*32

    # Attestation 2 (Failed)
    att2 = Attestation(indexes=[IndexLog(1,2)], round_id=123, round_status=r.status, request=b'mic2', fee=50)
    att2.status = AttestationStatus.WRONG_MIC

    r.add_attestation(att0)
    r.add_attestation(att1)
    r.add_attestation(att2)

    return r

# --- Test Cases ---

def test_add_attestation(mock_voter_set):
    """Tests adding new and duplicate attestations."""
    r = Round(round_id=1, voter_set=mock_voter_set)
    att = Attestation(indexes=[IndexLog(1,0)], round_id=1, round_status=r.status, request=b'mic_A', fee=10)

    assert r.add_attestation(att) # First time should be True
    assert len(r.attestations) == 1

    att_dup = Attestation(indexes=[IndexLog(1,1)], round_id=1, round_status=r.status, request=b'mic_A', fee=20)
    assert not r.add_attestation(att_dup) # Duplicate should be False
    assert len(r.attestations) == 1
    assert r.attestations[0].fee == 30 # Fee should be summed

def test_get_local_bit_vote(initial_round):
    """Tests the generation of the local node's bit vote."""
    # Attestations 0 and 1 are SUCCESS. Bit vector should be 0b011 = 3.
    bv = initial_round.get_local_bit_vote()
    assert bv.length == 3
    assert bv.bit_vector == 0b011

def test_process_incoming_bit_votes(initial_round):
    """Tests processing valid bit votes from other participants."""
    # Voter 2 (weight=30) votes for bits 1 and 2 (110)
    msg1 = {
        "from": "submit_signer2", "payload": b'\x00\x03\x06', # len=3, vec=6
        "block_number": 10, "tx_index": 1,
    }
    # Voter 3 (weight=30) votes for bit 0 (001)
    msg2 = {
        "from": "submit_signer3", "payload": b'\x00\x03\x01', # len=3, vec=1
        "block_number": 10, "tx_index": 2,
    }

    initial_round.process_incoming_bit_vote(msg1)
    initial_round.process_incoming_bit_vote(msg2)

    assert len(initial_round.bit_votes) == 2

    wbv2 = next(v for v in initial_round.bit_votes if v.index == 1)
    assert wbv2.weight == 30
    assert wbv2.bit_vote.bit_vector == 6

def test_full_round_consensus_and_merkle(initial_round):
    """
    An integration test for a full round, from processing votes to Merkle root.
    """
    # Our local node (voter 1, weight=40) votes for bits 0, 1 (011)
    local_bv = initial_round.get_local_bit_vote()
    initial_round.bit_votes.append(
        WeightedBitVote(index=0, weight=40, bit_vote=local_bv, index_tx=IndexTx(10,0))
    )

    # Voter 2 (weight=30) votes for bits 1, 2 (110)
    initial_round.process_incoming_bit_vote({
        "from": "submit_signer2", "payload": b'\x00\x03\x06',
        "block_number": 10, "tx_index": 1,
    })

    # Voter 3 (weight=30) votes for bit 0 (001)
    initial_round.process_incoming_bit_vote({
        "from": "submit_signer3", "payload": b'\x00\x03\x01',
        "block_number": 10, "tx_index": 2,
    })

    # --- Run Consensus ---
    # Bit 0 support: v1(40) + v3(30) = 70. Value = 7000.
    # Bit 1 support: v1(40) + v2(30) = 70. Value = 14000.
    # Bit 2 support: v2(30) = 30. (INVALID)
    # Optimal solution is to include only bit 1.
    initial_round.compute_consensus()

    # --- Verify Results ---
    assert initial_round.consensus_calculation_finished
    assert initial_round.consensus_bit_vote is not None

    # Expected final vector is bit 1 set -> 0b010 = 2
    expected_bv = 0b010
    assert initial_round.consensus_bit_vote.bit_vector == expected_bv

    # Verify consensus flags on attestations
    assert initial_round.attestations[0].consensus is False # Bit 0
    assert initial_round.attestations[1].consensus is True  # Bit 1
    assert initial_round.attestations[2].consensus is False # Bit 2

    # Verify Merkle root generation
    merkle_root = initial_round.get_merkle_root()
    assert isinstance(merkle_root, bytes)
    assert len(merkle_root) == 32
    assert initial_round.status.value == RoundStatus.DONE # Status updated