from src.attestation.bit_votes.structures import BitVote

def encode_bit_vote(bv: BitVote) -> bytes:
    """
    Encodes a BitVote object into a byte string.

    The encoding consists of the length (as a 2-byte big-endian integer)
    followed by the byte representation of the bit vector.

    Args:
        bv: The BitVote object to encode.

    Returns:
        The encoded byte string.
    """
    # Pack the length as a 2-byte big-endian integer
    length_bytes = bv.length.to_bytes(2, 'big')

    # Get the byte representation of the bit vector
    # The bit_length() + 7 // 8 formula correctly calculates the number of bytes
    bit_vector_bytes = bv.bit_vector.to_bytes((bv.bit_vector.bit_length() + 7) // 8, 'big')

    return length_bytes + bit_vector_bytes

def encode_bit_vote_hex(bv: BitVote) -> str:
    """
    Encodes a BitVote and returns its hexadecimal representation (without '0x').
    """
    return encode_bit_vote(bv).hex()

def decode_bit_vote_bytes(data: bytes) -> BitVote:
    """
    Decodes a byte string into a BitVote object.

    Args:
        data: The byte string to decode.

    Returns:
        The decoded BitVote object.

    Raises:
        ValueError: If the byte string is too short or represents an
                    invalid bit vote (e.g., bit length exceeds specified length).
    """
    if len(data) < 2:
        raise ValueError("Bit vote data is too short to contain a length.")

    length_bytes = data[0:2]
    bit_vector_bytes = data[2:]

    length = int.from_bytes(length_bytes, 'big')
    bit_vector = int.from_bytes(bit_vector_bytes, 'big')

    # The bit_length of the decoded integer cannot be greater than the
    # specified length of the bit vector.
    if bit_vector.bit_length() > length:
        raise ValueError("Invalid bit vote: bit vector is longer than specified length.")

    return BitVote(length=length, bit_vector=bit_vector)