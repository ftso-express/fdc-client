package attestation

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"flare-common/payload"
	"fmt"
	"local/fdc/client/shuffle"
	"math"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	NumOfSamples int = 100000 // read from config/toml
)

type BitVote struct {
	Length    uint16 //number attestations
	BitVector *big.Int
}

type IndexTx struct {
	BlockNumber      uint64
	TransactionIndex uint64
}

// earlierTx compares IndexTxs a,b. Returns true if a has lower BlockNumber than b or has the same BlockNumber and lower TransactionIndex.
func earlierTx(a, b IndexTx) bool {
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
	index       uint64
	bitVote     BitVote
	valueCapped *big.Int
	value       *big.Int // support multiplied with fees
}

// BitVoteFromAttestations calculates BitVote for an array of attestations.
// For i-th attestation in array, i-th bit in BitVote(from the right) is 1 if and only if i-th attestation status is Success.
// Sorting of attestation must be done prior.
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

	for i := range attestations {

		if bv.BitVector.Bit(i) == 1 {
			fees.Add(fees, attestations[i].Fee)
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

	supportingWeight := uint32(0) //supporting weight always fits into uint16. We use uint16 safe comparison.

	for j := range shuffled {
		if 2*supportingWeight < uint32(totalWeight) {
			bitVote = andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote)
			supportingWeight += uint32(weightedBitVotes[shuffled[j]].Weight)
		} else if andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote).BitVector.Cmp(bitVote.BitVector) == 0 {
			supportingWeight += uint32(weightedBitVotes[shuffled[j]].Weight)
		}

	}

	return bitVote, uint16(supportingWeight)

}

// bitVoteForSet calculates bitwise and of the WeightedBitVote in the order defined by shuffled
// until the added weight does not exceed 20% of the total weight.
// Then it adds the weight of the rest of WeightedBitVote that support the calculated BitVote.
// Returns the BitVote that is the result of the bitwise and, and supportingWeight.
//
// We assume that the sum of the weights of the WeightedBitVotes is more than 50% of th totalWeight.
func bitVoteForSetOptimistic(weightedBitVotes []*WeightedBitVote, totalWeight uint16, shuffled []uint64) (BitVote, uint16) {

	bitVote := (weightedBitVotes)[shuffled[0]].BitVote

	supportingWeight := uint32(0) //supporting weight always fits into uint16. We use uint16 safe comparison.

	for j := range shuffled {
		if 5*supportingWeight < uint32(totalWeight) {
			bitVote = andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote)
			supportingWeight += uint32(weightedBitVotes[shuffled[j]].Weight)
		} else if andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote).BitVector.Cmp(bitVote.BitVector) == 0 {
			supportingWeight += uint32(weightedBitVotes[shuffled[j]].Weight)
		}

	}

	return bitVote, uint16(supportingWeight)

}

// andBitwise returns the BitVote that has 1s at the places where both a and b have 1s and 0s elsewhere.
// If one BitVote is longer, the resulting BitVote has larger length with 0s on the excess places.
func andBitwise(a, b BitVote) BitVote {

	maxLen := max(a.Length, b.Length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

// Value calculates the Value of the BitVote, which is the product of the fees and supportingWeight.
func Value(bitVote BitVote, supportingWeight uint16, attestations []*Attestation) (*big.Int, error) {
	fees, err := bitVote.fees(attestations)

	if err != nil {
		return nil, fmt.Errorf("cannot compute value : %s", err)
	}

	return fees.Mul(fees, big.NewInt(int64(supportingWeight))), nil
}

// Value calculates the Value of the BitVote, which is the product of the fees and caped supportingWeight.
func ValueCapped(bitVote BitVote, supportingWeight, totalWeight uint16, attestations []*Attestation, capPercentage float64) (*big.Int, error) {
	fees, err := bitVote.fees(attestations)

	if err != nil {
		return nil, fmt.Errorf("cannot compute value : %s", err)
	}

	cap := uint16(math.Ceil(float64(totalWeight) * capPercentage))

	if cap < supportingWeight {

		return fees.Mul(fees, big.NewInt(int64(cap))), nil

	}

	return fees.Mul(fees, big.NewInt(int64(supportingWeight))), nil
}

type tempBitVoteResult struct {
	index          uint64
	maxValueCapped *big.Int
	maxValue       *big.Int
	bitVote        BitVote
}

type ConsensusBitVoteInput struct {
	RoundID          uint64
	WeightedBitVotes []*WeightedBitVote
	TotalWeight      uint16
	Attestations     []*Attestation
}

// ConsensusBitVote calculates the ConsensusBitVote for roundId given the weightedBitVotes.
func ConsensusBitVote(input *ConsensusBitVoteInput) (BitVote, error) {
	weightVoted := sumWeightedBitVotes(input.WeightedBitVotes)

	if weightVoted <= input.TotalWeight/2 {
		percentage := (float32(weightVoted) * 100.0) / float32(input.TotalWeight)
		return BitVote{}, errors.Errorf("only %.1f%% bitVoted", percentage)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	tmpResult := tempBitVoteResult{
		maxValueCapped: big.NewInt(0),
		maxValue:       big.NewInt(0),
	}

	for i := 0; i < NumOfSamples; i++ {
		index := uint64(i)
		eg.Go(func() error {
			bitVoteVals, err := calcBitVoteVals(ctx, input, index)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			for j := range bitVoteVals {
				updateTmpResult(&tmpResult, &bitVoteVals[j])
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return BitVote{}, err
	}

	return tmpResult.bitVote, nil
}

func sumWeightedBitVotes(weightedBitVotes []*WeightedBitVote) (weightVoted uint16) {
	for j := range weightedBitVotes {
		weightVoted += (weightedBitVotes[j].Weight)
	}

	return weightVoted
}

const (
	valueCap   = 4.0 / 5.0
	protocolID = 300 // TODO get protocol id (300) from somewhere
)

func calcBitVoteVals(ctx context.Context, input *ConsensusBitVoteInput, index uint64) ([]bitVoteWithValue, error) {
	results := make([]bitVoteWithValue, 0, 2)

	seed := shuffle.Seed(input.RoundID, index, protocolID)
	shuffled := shuffle.FisherYates(uint64(len(input.WeightedBitVotes)), seed)

	tempBitVote, supportingWeight := bitVoteForSet(input.WeightedBitVotes, input.TotalWeight, shuffled)
	value, err := Value(tempBitVote, supportingWeight, input.Attestations)
	if err != nil {
		return nil, err
	}

	valueCapped, err := ValueCapped(
		tempBitVote, supportingWeight, input.TotalWeight, input.Attestations, valueCap,
	)
	if err != nil {
		return nil, err
	}

	results = append(results, bitVoteWithValue{
		index:       index,
		bitVote:     tempBitVote,
		valueCapped: valueCapped,
		value:       value,
	})

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	tempBitVoteOpt, supportingWeightOpt := bitVoteForSetOptimistic(
		input.WeightedBitVotes, input.TotalWeight, shuffled,
	)

	if supportingWeightOpt > input.TotalWeight/2 {
		valueOptimistic, err := Value(tempBitVoteOpt, supportingWeightOpt, input.Attestations)
		if err != nil {
			return nil, err
		}

		valueOptimisticCapped, err := ValueCapped(
			tempBitVoteOpt, supportingWeightOpt, input.TotalWeight, input.Attestations, valueCap,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, bitVoteWithValue{
			index:       index,
			bitVote:     tempBitVoteOpt,
			valueCapped: valueOptimisticCapped,
			value:       valueOptimistic,
		})
	}

	return results, nil
}

func updateTmpResult(tmpResult *tempBitVoteResult, result *bitVoteWithValue) {
	if result.valueCapped.Cmp(tmpResult.maxValueCapped) == 1 {
		tmpResult.bitVote = result.bitVote
		tmpResult.index = result.index
		tmpResult.maxValueCapped = result.valueCapped
		tmpResult.maxValue = result.value
	} else if result.valueCapped.Cmp(tmpResult.maxValueCapped) == 0 && result.value.Cmp(tmpResult.maxValue) == 1 {
		tmpResult.bitVote = result.bitVote
		tmpResult.index = result.index
		tmpResult.maxValue = result.value
	} else if result.value.Cmp(tmpResult.maxValue) == 0 && tmpResult.index > result.index {
		tmpResult.bitVote = result.bitVote
		tmpResult.index = result.index
	}
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

// DecodeBitVoteBytes decodes bytes encoded BitVote and returns roundCheck
func DecodeBitVoteBytes(bitVoteByte []byte) (BitVote, uint8, error) {

	if len(bitVoteByte) < 3 {
		return BitVote{}, 0, errors.New("bitVote too short")
	}

	roundCheck := bitVoteByte[0]
	lengthBytes := bitVoteByte[1:3]
	bitVector := bitVoteByte[3:]

	length := binary.BigEndian.Uint16(lengthBytes)

	bigBitVector := new(big.Int).SetBytes(bitVector)

	// TODO: decide whether leading zeros are legal

	if bigBitVector.BitLen() > int(length) {
		return BitVote{}, 0, errors.New("bad bitvote")

	}

	return BitVote{length, bigBitVector}, roundCheck, nil

}

// ProcessBitVote decodes bitVote message, checks roundCheck, adds voter weight and index, and stores bitVote to the round.
// If the voter is invalid, or has zero weight, the bitVote is ignored.
// If a voter already submitted a valid bitVote for the round, the bitVote is overwritten.
func (r *Round) ProcessBitVote(message payload.Message) error {

	bitVote, roundCheck, err := DecodeBitVoteBytes(message.Payload)

	if err != nil {
		return err
	}

	if roundCheck != uint8(message.VotingRound%256) {
		return fmt.Errorf("wrong round check from %s", message.From)
	}

	voter, exists := r.voterSet.VoterDataMap[message.From]

	if !exists {
		return fmt.Errorf("invalid voter %s", message.From)
	}

	weight := voter.Weight

	if weight <= 0 {
		return fmt.Errorf("zero weight voter %s ", message.From)
	}

	// check if a bitVote was already submitted by the sender
	weightedBitVote, exists := r.bitVoteCheckList[message.From]

	if !exists {
		// first submission

		weightedBitVote = &WeightedBitVote{}
		r.bitVotes = append(r.bitVotes, weightedBitVote)

		weightedBitVote.BitVote = bitVote
		weightedBitVote.Weight = weight
		weightedBitVote.Index = voter.Index
		weightedBitVote.indexTx = IndexTx{message.BlockNumber, message.TransactionIndex}
	} else if exists && earlierTx(weightedBitVote.indexTx, IndexTx{message.BlockNumber, message.TransactionIndex}) {
		// more than one submission. The later submission is considered to be valid.

		weightedBitVote.BitVote = bitVote
		weightedBitVote.Weight = weight
		weightedBitVote.Index = voter.Index
		weightedBitVote.indexTx = IndexTx{message.BlockNumber, message.TransactionIndex}

	}

	return nil
}
