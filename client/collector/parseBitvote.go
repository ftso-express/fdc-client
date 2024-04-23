package collector

import (
	"errors"
	"flare-common/payload"
	"local/fdc/client/bitvote"
	"math/big"
	"strconv"
)

func ExtractBitVote(message payload.Message) (bitvote.BitVote, string, error) {

	var bitVote bitvote.BitVote

	roundSig, _ := strconv.ParseUint(message.Payload[:2], 16, 64)

	if message.VotingRound%256 != roundSig {
		return bitVote, "", errors.New("round does not match")
	}

	len, _ := strconv.ParseUint(message.Payload[2:6], 16, 16)

	bitVector := big.NewInt(0)

	bitVector.SetString(message.Payload[2:], 16)

	bitVote.Length = uint16(len)
	bitVote.BitVector = bitVector

	return bitVote, message.From, nil

}
