package attestation

import (
	"errors"
	"flare-common/merkle"
	"flare-common/policy"
	"fmt"
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
	bitVotes         []*WeightedBitVote
	bitVoteCheckList map[common.Address]*WeightedBitVote
	ConsensusBitVote BitVote
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
func (r *Round) BitVote() (BitVote, error) {

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

	consensus, err := ConsensusBitVote(&ConsensusBitVoteInput{
		RoundID:          r.roundId,
		WeightedBitVotes: r.bitVotes,
		TotalWeight:      r.voterSet.TotalWeight,
		Attestations:     r.Attestations,
	})
	if err != nil {
		return err

	}
	r.ConsensusBitVote = consensus

	return r.SetConsensusStatus()
}

func (r *Round) GetConsensusBitVote() (BitVote, error) {

	if r.ConsensusBitVote.BitVector == nil {
		return BitVote{}, errors.New("no consensus bitVote")
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
