from typing import List
from src.attestation.structures import Attestation, AttestationStatus
from src.attestation.bit_votes.structures import BitVote

def bit_vote_from_attestations(attestations: List[Attestation]) -> BitVote:
    """
    Calculates a BitVote from a list of attestations.

    For the i-th attestation in the list, the i-th bit in the resulting
    BitVote's bit_vector is set to 1 if and only if the attestation's
    status is SUCCESS.

    The list of attestations is assumed to be sorted correctly prior to
    calling this function.

    Args:
        attestations: A list of Attestation objects.

    Returns:
        A BitVote object representing the successful attestations.

    Raises:
        ValueError: If the number of attestations exceeds the maximum
                    limit of 65535.
    """
    num_attestations = len(attestations)

    if num_attestations > 65535:
        raise ValueError("Cannot create a bit vote for more than 65535 attestations.")

    bit_vector = 0
    for i, att in enumerate(attestations):
        if att.status == AttestationStatus.SUCCESS:
            bit_vector |= 1 << i

    return BitVote(length=num_attestations, bit_vector=bit_vector)