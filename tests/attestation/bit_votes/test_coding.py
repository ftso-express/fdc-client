import pytest
from src.attestation.bit_votes.structures import BitVote
from src.attestation.bit_votes.coding import (
    encode_bit_vote,
    encode_bit_vote_hex,
    decode_bit_vote_bytes,
)

@pytest.mark.parametrize("length, vector, expected_bytes", [
    (10, 0b10101, b'\x00\n\x15'),  # length 10, vector 21
    (8, 0b11111111, b'\x00\x08\xff'), # length 8, vector 255
    (16, 0b1000000000000000, b'\x00\x10\x80\x00'), # length 16, vector 32768
    (0, 0, b'\x00\x00'), # Empty bit vote
])
def test_encode_decode_roundtrip(length, vector, expected_bytes):
    """
    Tests that encoding and decoding a BitVote results in the original object.
    """
    # Create the initial BitVote
    bv = BitVote(length=length, bit_vector=vector)

    # Test encoding
    encoded = encode_bit_vote(bv)
    assert encoded == expected_bytes

    # Test hex encoding
    assert encode_bit_vote_hex(bv) == expected_bytes.hex()

    # Test decoding
    decoded_bv = decode_bit_vote_bytes(encoded)

    # Verify the decoded object matches the original
    assert decoded_bv.length == bv.length
    assert decoded_bv.bit_vector == bv.bit_vector

def test_decode_error_too_short():
    """
    Tests that decoding fails if the byte string is too short.
    """
    with pytest.raises(ValueError, match="Bit vote data is too short"):
        decode_bit_vote_bytes(b'\x00')

def test_decode_error_invalid_vector():
    """
    Tests that decoding fails if the bit vector is longer than the
    specified length.
    """
    # Length is 4, but the vector 0b10000 requires 5 bits.
    invalid_data = b'\x00\x04\x10'
    with pytest.raises(ValueError, match="Invalid bit vote: bit vector is longer than specified length."):
        decode_bit_vote_bytes(invalid_data)

def test_encode_empty_vector():
    """
    Tests encoding a BitVote with a zero vector but non-zero length.
    """
    bv = BitVote(length=10, bit_vector=0)
    # The vector is 0, so its byte representation is empty.
    # The encoding should just be the length.
    assert encode_bit_vote(bv) == b'\x00\n'

    # Test roundtrip
    decoded = decode_bit_vote_bytes(b'\x00\n')
    assert decoded.length == 10
    assert decoded.bit_vector == 0