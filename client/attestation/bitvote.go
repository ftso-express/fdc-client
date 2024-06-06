package attestation

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flare-common/payload"
	"fmt"
	"local/fdc/client/shuffle"
	"math"
	"math/big"
	"sync"
)

const (
	NumOfSamples int     = 100000    // the actual number of samples is x2 (one normal and one optimistic for each seed)
	divOpt       uint16  = 5         // totalWeight/divOpt is the weight of the optimistic samples
	valueCap     float64 = 4.0 / 5.0 // bitVote support cap in factor of totalWeight
	protocolId           = 300       // TODO get protocol id (300) from somewhere // read from config/toml
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
	index      int
	bitVote    BitVote
	valueCaped *big.Int
	value      *big.Int // support multiplied with fees
	empty      bool
	err        error
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

	supportingWeight := uint16(0)

	auxBigInt := new(big.Int)

	for j := range shuffled {
		if supportingWeight < totalWeight/2 {
			bitVote = andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote)
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		} else if auxBigInt.And(bitVote.BitVector, weightedBitVotes[shuffled[j]].BitVote.BitVector).Cmp(bitVote.BitVector) == 0 {
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		}

	}

	return bitVote, supportingWeight

}

// bitVoteForSet calculates bitwise and of the WeightedBitVote in the order defined by shuffled
// until the added weight does not exceed x% of the total weight.
// Then it adds the weight of the rest of WeightedBitVote that support the calculated BitVote.
// Returns the BitVote that is the result of the bitwise and, and supportingWeight.
// It is assumed that len(weightedBitVotes) = len(shuffled) > 0, shuffled contains each integer from 0 to len(shuffled)- 1.
func bitVoteForSetOptimistic(weightedBitVotes []*WeightedBitVote, totalWeight uint16, shuffled []uint64, div uint16) (BitVote, uint16) {

	bitVote := (weightedBitVotes)[shuffled[0]].BitVote

	supportingWeight := uint16(0)

	auxBigInt := new(big.Int)

	for j := range shuffled {
		if supportingWeight < totalWeight/div {
			bitVote = andBitwise(bitVote, weightedBitVotes[shuffled[j]].BitVote)
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		} else if auxBigInt.And(bitVote.BitVector, weightedBitVotes[shuffled[j]].BitVote.BitVector).Cmp(bitVote.BitVector) == 0 {
			supportingWeight += weightedBitVotes[shuffled[j]].Weight
		}

	}

	return bitVote, supportingWeight

}

// andBitwise returns the BitVote that has 1s at the places where both a and b have 1s and 0s elsewhere.
// If one BitVote is longer, the resulting BitVote has larger length with 0s on the excess places.
func andBitwise(a, b BitVote) BitVote {

	maxLen := max(a.Length, b.Length)

	bitVector := big.NewInt(0)

	bitVector.And(a.BitVector, b.BitVector)

	return BitVote{maxLen, bitVector}

}

// Value calculates the cappedValue and Value of the BitVote, which is the product of the fees and supportingWeight.
func Values(bitVote BitVote, supportingWeight uint16, totalWeight uint16, attestations []*Attestation, capPercentage float64) (*big.Int, *big.Int, error) {
	fees, err := bitVote.fees(attestations)

	if err != nil {
		return nil, nil, fmt.Errorf("cannot compute value : %s", err)
	}

	auxBigInt := new(big.Int)

	cap := uint16(math.Ceil(float64(totalWeight) * capPercentage))

	if cap < supportingWeight {

		return auxBigInt.Mul(fees, big.NewInt(int64(cap))), auxBigInt.Mul(fees, big.NewInt(int64(supportingWeight))), nil

	}

	return auxBigInt.Mul(fees, big.NewInt(int64(supportingWeight))), auxBigInt, nil
}

type safeTempResult struct {
	mu            *sync.Mutex
	index         int
	maxValueCaped *big.Int
	maxValue      *big.Int
	bitVote       BitVote
	err           error
}

func sumWeightedBitVotes(weightedBitVotes []*WeightedBitVote) (weightVoted uint16) {
	for j := range weightedBitVotes {
		weightVoted += weightedBitVotes[j].Weight
	}

	return weightVoted
}

// ConsensusBitVote calculates the ConsensusBitVote for roundId given the weightedBitVotes.
func ConsensusBitVote(roundId uint64, weightedBitVotes []*WeightedBitVote, totalWeight uint16, attestations []*Attestation) (BitVote, error) {

	noOfVoters := len(weightedBitVotes)

	weightVoted := sumWeightedBitVotes(weightedBitVotes)

	if weightVoted <= totalWeight/2 {

		percentage := (weightVoted * 100) / totalWeight
		return BitVote{}, fmt.Errorf("only %d%% bitVoted", percentage)
	}

	ch := make(chan bitVoteWithValue)

	go func() {

		for i := 0; i < NumOfSamples; i++ {
			go func(j int) {

				seed := shuffle.Seed(roundId, uint64(j), protocolId)
				shuffled := shuffle.FisherYates(uint64(noOfVoters), seed)

				go func() {
					tempBitVote, supportingWeight := bitVoteForSet(weightedBitVotes, totalWeight, shuffled)
					valueCapped, value, err := Values(tempBitVote, supportingWeight, totalWeight, attestations, valueCap)

					ch <- bitVoteWithValue{j, tempBitVote, valueCapped, value, false, err}
				}()

				go func() {

					tempBitVoteOpt, supportingWeightOpt := bitVoteForSetOptimistic(weightedBitVotes, totalWeight, shuffled, divOpt)

					if supportingWeightOpt > totalWeight/2 {

						valueOptimisticCapped, valueOptimistic, err := Values(tempBitVoteOpt, supportingWeightOpt, totalWeight, attestations, valueCap)

						ch <- bitVoteWithValue{j, tempBitVoteOpt, valueOptimisticCapped, valueOptimistic, false, err}
					} else {
						ch <- bitVoteWithValue{empty: true}
					}

				}()

			}(i)
		}
	}()

	tempResult := &safeTempResult{}
	tempResult.mu = new(sync.Mutex)
	tempResult.maxValue = big.NewInt(0)
	tempResult.maxValueCaped = big.NewInt(0)

	tempResult.index = 0

	var wg sync.WaitGroup

	for i := 0; i < 2*NumOfSamples; i++ {
		result := <-ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			if !result.empty {
				if result.err != nil {
					tempResult.mu.Lock()

					tempResult.err = result.err

					tempResult.mu.Unlock()
				} else {

					isBetterCappedUnsafe := result.valueCaped.Cmp(tempResult.maxValueCaped)

					if isBetterCappedUnsafe == 1 {

						tempResult.mu.Lock()

						if result.valueCaped.Cmp(tempResult.maxValueCaped) == 1 {

							tempResult.bitVote = result.bitVote
							tempResult.index = result.index
							tempResult.maxValueCaped = result.valueCaped
							tempResult.maxValue = result.value

						}
						tempResult.mu.Unlock()

					} else if isBetterCappedUnsafe == 0 {

						isBetterUncappedUnsafe := result.value.Cmp(tempResult.maxValue)

						if isBetterUncappedUnsafe == 1 {
							tempResult.mu.Lock()

							if result.valueCaped.Cmp(tempResult.maxValueCaped) == 0 && result.value.Cmp(tempResult.maxValue) == 1 {

								tempResult.bitVote = result.bitVote
								tempResult.index = result.index
								tempResult.maxValue = result.value

							}
							tempResult.mu.Unlock()

						} else if isBetterUncappedUnsafe == 0 && tempResult.index > result.index {

							tempResult.mu.Lock()

							if result.valueCaped.Cmp(tempResult.maxValueCaped) == 0 && result.value.Cmp(tempResult.maxValue) == 0 && tempResult.index > result.index {
								tempResult.bitVote = result.bitVote
								tempResult.index = result.index
							}
							tempResult.mu.Unlock()

						}

					}

				}
			}

		}()

	}
	wg.Wait()

	if tempResult.err != nil {
		return BitVote{}, tempResult.err
	}

	return tempResult.bitVote, nil
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
