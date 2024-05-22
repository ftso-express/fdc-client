package attestation

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flare-common/payload"
	"fmt"
	"local/fdc/client/shuffle"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	NumOfSamples int = 100 // read from config/toml
)

type BitVote struct {
	Length    uint16 //number attestations
	BitVector *big.Int
}

type IndexTx struct {
	BlockNumber      uint64
	TransactionIndex uint64
}

// LessTx compares IndexTxs a,b. Returns true if a has lower BlockNumber than b or has the same BlockNumber and lower TransactionIndex.
func LessTx(a, b IndexTx) bool {
	if a.BlockNumber < b.BlockNumber {
		return true
	}
	if a.BlockNumber == b.BlockNumber && a.TransactionIndex < b.TransactionIndex {
		return true
	}

	return false

}

type WeightedBitVote struct {
	Index   int
	indexTx IndexTx
	Weight  uint16
	BitVote BitVote
}

type bitVoteWithValue struct {
	index   int64
	bitVote BitVote
	value   *big.Int // support multiplied with fees
	err     error
}

// BitVoteFromAttestations calculates BitVote for an array of attestations.
// For i-th attestation in array, i-th bit in BitVote(from the right) is 1 if and only if i-th attestation status is Success.
func BitVoteFromAttestations(attestations []*Attestation) (BitVote, error) {
	bitVector := big.NewInt(0)

	if len(attestations) > 65535 {
		return BitVote{}, errors.New("more than 65536 attestations")
	}

	for i, a := range attestations {
		if a.Status == Success {
			bitVector.SetBit(bitVector, i, 1)
		}

	}
	return BitVote{uint16(len(attestations)), bitVector}, nil
}

// fees sums the fees of the attestation requests indicated in BitVote
func (bv BitVote) fees(attestations []*Attestation) (*big.Int, error) {

	if bv.BitVector.BitLen() > len(attestations) {
		return nil, errors.New("a confirmed instance missing from attestations")
	}

	fees := big.NewInt(0)

	for i, a := range attestations {

		if bv.BitVector.Bit(i) == 1 {
			fees.Add(fees, a.Fee)
		}
	}
	return fees, nil
}

// bitVoteForSet calculates bitwise and of the WeightedBitVote in the order defined by shuffled
// until the added weight does not exceed 50% of the total weight.
// Then it adds the weight of the rest of WeightedBitVote that support the calculated BitVote.
// Returns the BitVote that is the result of the bitwise and, and supportingWeight.
//
// We assume that the sum of the weights of the WeightedBitVotes is more than 50% of th totalWeight.
func bitVoteForSet(weightedBitVotes []*WeightedBitVote, totalWeight uint16, shuffled []uint64) (BitVote, uint16) {

	bitVote := (weightedBitVotes)[shuffled[0]].BitVote

	halfWeight := (totalWeight + 1) / 2

	supportingWeight := uint16(0)

	for _, v := range shuffled {
		if supportingWeight < halfWeight {
			bitVote = andBitwise(bitVote, weightedBitVotes[v].BitVote)
			supportingWeight += weightedBitVotes[v].Weight
		} else if andBitwise(bitVote, weightedBitVotes[v].BitVote).BitVector == weightedBitVotes[v].BitVote.BitVector {
			supportingWeight += weightedBitVotes[v].Weight
		}

	}

	return bitVote, supportingWeight

}

// andBitwise returns the BitVote that has 1s at the places where both a and b have 1s and 0s elsewhere.
// If one BitVote is longer the resulting BitVote has larger length with 0s on the excess places.
func andBitwise(a, b BitVote) BitVote {

	maxLen := max(a.Length, b.Length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

// Value calculates the value of the BitVote, which is the product of the fees and supportingWeight.
func value(bitVote BitVote, supportingWeight uint16, attestations []*Attestation) (*big.Int, error) {
	fees, err := bitVote.fees(attestations)

	if err != nil {
		return nil, fmt.Errorf("cannot compute value : %s", err)
	}

	return fees.Mul(fees, big.NewInt(int64(supportingWeight))), nil
}

// ConsensusBitVote calculates the ConsensusBitVote for roundId given the weightedBitVotes.
func ConsensusBitVote(roundId uint64, weightedBitVotes []*WeightedBitVote, totalWeight uint16, attestations []*Attestation) (BitVote, error) {

	noOfVoters := len(weightedBitVotes)

	weightVoted := uint16(0)
	for j := range weightedBitVotes {
		weightVoted += weightedBitVotes[j].Weight
	}

	if (totalWeight+1)/2 > weightVoted {

		percentage := (weightVoted * 100) / totalWeight
		log.Printf("Only %d%% voted in round %d.", roundId, percentage)
		return BitVote{}, errors.New("not enough weight bitVoted to get a consensus")
	}

	var bitVote BitVote
	maxValue := big.NewInt(0)

	index := int64(0)

	ch := make(chan bitVoteWithValue)

	for i := 0; i < NumOfSamples; i++ {
		go func(j int64) {
			seed := shuffle.Seed(int64(roundId), j)
			shuffled := shuffle.FisherYates(uint64(noOfVoters), seed)
			tempBitVote, supportingWeight := bitVoteForSet(weightedBitVotes, totalWeight, shuffled)
			value, err := value(tempBitVote, supportingWeight, attestations)

			ch <- bitVoteWithValue{j, tempBitVote, value, err}
		}(int64(i))
	}

	for i := 0; i < NumOfSamples; i++ {
		result := <-ch

		if result.err != nil {
			log.Printf("Cannot compute consensus bitVote round %d because of a missing attestation.", roundId)

			return BitVote{}, errors.New("missing attestations. cannot compute consensus bitvote")
		}

		if result.value.Cmp(maxValue) == 1 {
			bitVote = result.bitVote
			index = result.index
			maxValue = result.value
		} else if result.value.Cmp(maxValue) == 0 && index > result.index {
			bitVote = result.bitVote
			index = result.index
		}
	}

	return bitVote, nil
}

// SetBitVoteStatus sets the Consensus status of attestations to true for the attestations chosen by BitVote
func SetBitVoteStatus(attestations []*Attestation, bitVote BitVote) error {

	if bitVote.BitVector.BitLen() > len(attestations) {
		return errors.New("chosen attestation does not exist")
	}

	for i := range attestations {
		attestations[i].Consensus = bitVote.BitVector.Bit(i) == 1
	}

	return nil

}

// EncodeBitVoteHex encodes BitVote with roundCheck to be published on chain
func (b BitVote) EncodeBitVoteHex(roundId uint64) string {

	var encoding []byte
	roundCheck := byte(roundId % 256)

	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, b.Length)

	encoding = append(encoding, roundCheck)
	encoding = append(encoding, length...)
	encoding = append(encoding, b.BitVector.Bytes()...)

	str := hex.EncodeToString(encoding)

	return str

}

// DecodeBitVoteHex decodes hex encoded BitVote and returns roundCheck
func DecodeBitVoteHex(bitVoteHex string) (BitVote, uint8, error) {

	roundCheckStr := bitVoteHex[:2]
	lengthStr := bitVoteHex[2:6]
	bitVectorStr := bitVoteHex[6:]

	roundCheckBytes, err := hex.DecodeString(roundCheckStr)
	if err != nil || len(roundCheckBytes) != 1 {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	roundCheck := uint8(roundCheckBytes[0])

	lengthBytes, err := hex.DecodeString(lengthStr)
	if err != nil || len(lengthBytes) != 2 {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	length := binary.BigEndian.Uint16(lengthBytes)

	bitVector := big.NewInt(0)
	_, success := bitVector.SetString(bitVectorStr, 16)

	if !success {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	if bitVector.BitLen() > int(length) {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	return BitVote{length, bitVector}, roundCheck, nil

}

// ProcessBitVote decodes bitVote message, checks roundCheck, adds voter weight and index, and stores bitVote to the round.
// If the voter is invalid, or has zero weight, the bitVote is ignored.
// If a voter already submitted a valid bitVote for the round, the bitVote is overwritten.
func (r *Round) ProcessBitVote(message payload.Message) error {

	bitVote, roundCheck, err := DecodeBitVoteHex(message.Payload)

	if err != nil {
		return err
	}

	if roundCheck != uint8(message.VotingRound%256) {
		return errors.New("wrong round check")
	}

	voter, success := r.voterSet.VoterDataMap[common.HexToAddress(message.From)]

	if !success {
		return errors.New("invalid voter")
	}

	weight := voter.Weight

	if weight <= 0 {
		return errors.New("zero weight")
	}

	// check if a bitVote was already submitted by the sender
	weightedBitVote, ok := r.bitVoteCheckList[message.From]

	// first submission
	if !ok {
		weightedBitVote = &WeightedBitVote{}
		r.bitVotes = append(r.bitVotes, weightedBitVote)

		weightedBitVote.BitVote = bitVote
		weightedBitVote.Weight = weight
		weightedBitVote.Index = voter.Index
		weightedBitVote.indexTx = IndexTx{message.BlockNumber, message.TransactionIndex}
	}

	// more than one submission. The later submission is considered to be valid.
	if ok && LessTx(weightedBitVote.indexTx, IndexTx{message.BlockNumber, message.TransactionIndex}) {

		weightedBitVote.BitVote = bitVote
		weightedBitVote.Weight = weight
		weightedBitVote.Index = voter.Index
		weightedBitVote.indexTx = IndexTx{message.BlockNumber, message.TransactionIndex}

	}

	return nil
}
