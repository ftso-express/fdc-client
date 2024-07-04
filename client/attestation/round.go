package attestation

import (
	"errors"
	"flare-common/merkle"
	"flare-common/payload"
	"flare-common/policy"
	"fmt"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

type RoundStatus int

const (
	Collecting RoundStatus = iota
	Choosing
	Retry
	Committed
	Done
	Failed
)

type Round struct {
	roundId          uint64
	Attestations     []*Attestation
	bitVotes         []*bitvotes.WeightedBitVote
	bitVoteCheckList map[common.Address]*bitvotes.WeightedBitVote
	ConsensusBitVote bitvotes.BitVote
	voterSet         *policy.VoterSet
	merkleTree       merkle.Tree
}

// CreateRound returns a pointer to a new round with roundId and voterSet.
func CreateRound(roundId uint64, voterSet *policy.VoterSet) *Round {

	r := &Round{}

	r.roundId = roundId

	r.voterSet = voterSet

	return r
}

// sortAttestations sorts round's attestations according to their IndexLog.
func (r *Round) sortAttestations() {

	sort.Slice(r.Attestations, func(i, j int) bool {
		return earlierLog(r.Attestations[i].Index, r.Attestations[j].Index)
	})
}

// sortBitVotes sorts rounds' bitVotes according to the signingPolicy Index of their providers.
func (r *Round) sortBitVotes() {

	sort.Slice(r.bitVotes, func(i, j int) bool {
		return r.bitVotes[i].Index < r.bitVotes[j].Index
	})
}

// BitVote returns the BitVote for the round according to the current status of Attestations.
func (r *Round) BitVote() (bitvotes.BitVote, error) {

	r.sortAttestations()
	return BitVoteFromAttestations(r.Attestations)
}

// BitVoteHex returns the hex string encoded BitVote for the round according to the current status of Attestations.
func (r *Round) BitVoteHex() (string, error) {

	r.sortAttestations()

	bitVote, err := BitVoteFromAttestations(r.Attestations)

	if err != nil {
		return "", fmt.Errorf("cannot get bitvote for round %d", r.roundId)
	}

	return bitVote.EncodeBitVoteHex(r.roundId), nil
}

// ComputeConsensusBitVote computes the consensus BitVote according to the collected bitVotes.
func (r *Round) ComputeConsensusBitVote() error {

	r.sortBitVotes()

	r.sortAttestations()

	fees := make([]int, len(r.Attestations))
	for i, e := range r.Attestations {
		fees[i] = int(e.Fee.Int64())
	}

	consensus, err := bitvotes.ConsensusBitVote(&bitvotes.ConsensusBitVoteInput{
		RoundID:          r.roundId,
		WeightedBitVotes: r.bitVotes,
		TotalWeight:      r.voterSet.TotalWeight,
		Fees:             fees,
	})
	if err != nil {
		return err

	}
	r.ConsensusBitVote = consensus

	return r.SetConsensusStatus()
}

func (r *Round) GetConsensusBitVote() (bitvotes.BitVote, error) {

	if r.ConsensusBitVote.BitVector == nil {
		return bitvotes.BitVote{}, errors.New("no consensus bitVote")
	}

	return r.ConsensusBitVote, nil
}

// SetConsensusStatus sets consensus status of the attestations.
// The scenario where a chosen attestation is missing is not possible as in such case, it is not possible to compute the consensus bitVote.
// It is assumed that the Attestations are already ordered.
func (r *Round) SetConsensusStatus() error {

	consensusBitVote, err := r.GetConsensusBitVote()

	if err != nil {
		return err
	}

	for i := range r.Attestations {
		r.Attestations[i].Consensus = consensusBitVote.BitVector.Bit(i) == 1
	}

	return nil

}

// GetMerkleTree computes Merkle tree from sorted hashes of attestations chosen by the consensus bitVote.
// The computed tree is stored in the round.
// If any of the hash of the chosen attestations is missing, the tree is not computed.
func (r *Round) GetMerkleTree() (merkle.Tree, error) {

	r.sortAttestations()

	hashes := []common.Hash{}

	added := make(map[common.Hash]bool)

	for i := range r.Attestations {
		if r.Attestations[i].Consensus {
			if r.Attestations[i].Status != Success {
				return merkle.Tree{}, errors.New("cannot build merkle tree")
			}

			//skip duplicates
			if _, alreadyAdded := added[r.Attestations[i].Hash]; !alreadyAdded {

				hashes = append(hashes, r.Attestations[i].Hash)
				added[r.Attestations[i].Hash] = true

			}

		}
	}

	// sort.Slice(hashes, func(i, j int) bool { return compareHash(hashes[i], hashes[j]) })

	merkleTree := merkle.Build(hashes, false)

	r.merkleTree = merkleTree

	return merkleTree, nil

}

// func compareHash(a, b common.Hash) bool {

//		for i := range a {
//			if a[i] < b[i] {
//				return true
//			}
//		}
//		return false
//	}

func (r *Round) GetMerkleTreeCached() (merkle.Tree, error) {

	if len(r.merkleTree) != 0 {
		return r.merkleTree, nil
	}

	return r.GetMerkleTree()

}

func (r *Round) GetMerkleRootCached() (common.Hash, error) {

	tree, err := r.GetMerkleTreeCached()

	if err != nil {
		return common.Hash{}, err
	}

	return tree.Root()

}

func (r *Round) GetMerkleRootCachedHex() (string, error) {

	root, err := r.GetMerkleRootCached()

	if err != nil {
		return "", err
	}

	return root.Hex(), nil

}

// ProcessBitVote decodes bitVote message, checks roundCheck, adds voter weight and index, and stores bitVote to the round.
// If the voter is invalid, or has zero weight, the bitVote is ignored.
// If a voter already submitted a valid bitVote for the round, the bitVote is overwritten.
func (r *Round) ProcessBitVote(message payload.Message) error {

	bitVote, roundCheck, err := bitvotes.DecodeBitVoteBytes(message.Payload)

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

		weightedBitVote = &bitvotes.WeightedBitVote{}
		r.bitVotes = append(r.bitVotes, weightedBitVote)

		weightedBitVote.BitVote = bitVote
		weightedBitVote.Weight = weight
		weightedBitVote.Index = voter.Index
		weightedBitVote.IndexTx = bitvotes.IndexTx{message.BlockNumber, message.TransactionIndex}
	} else if exists && bitvotes.EarlierTx(weightedBitVote.IndexTx, bitvotes.IndexTx{message.BlockNumber, message.TransactionIndex}) {
		// more than one submission. The later submission is considered to be valid.

		weightedBitVote.BitVote = bitVote
		weightedBitVote.Weight = weight
		weightedBitVote.Index = voter.Index
		weightedBitVote.IndexTx = bitvotes.IndexTx{message.BlockNumber, message.TransactionIndex}

	}

	return nil
}
