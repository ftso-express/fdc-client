package attestation

import (
	"errors"
	"flare-common/merkle"
	"flare-common/policy"
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
	status           RoundStatus
	Attestations     []*Attestation
	bitVotes         []*WeightedBitVote
	bitVoteCheckList map[string]*WeightedBitVote
	ConsensusBitVote BitVote
	voterSet         *policy.VoterSet
	merkletree       merkle.Tree
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
		return lessLog(r.Attestations[i].Index, r.Attestations[j].Index)
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
		return "", errors.New("cannot get bitvote")
	}

	return bitVote.EncodeBitVoteHex(r.roundId), nil
}

// ComputeConsensusBitVote computes the consensus BitVote according to the collected bitVotes.
func (r *Round) ComputeConsensusBitVote() error {

	r.sortBitVotes()

	consensus, err := ConsensusBitVote(r.roundId, r.bitVotes, r.voterSet.TotalWeight, r.Attestations)

	if err != nil {
		return err

	}
	r.ConsensusBitVote = consensus

	return nil
}

func (r *Round) GetConsensusBitVote() BitVote {

	return r.ConsensusBitVote
}

func (r *Round) SetConsensusStatus() error {

	r.sortAttestations()

	// handle no bitVote or chosen request that is not registered

	for i := range r.Attestations {
		r.Attestations[i].Consensus = r.ConsensusBitVote.BitVector.Bit(i) == 1
	}

	return nil

}

// GetMerkleTree computes Merkle tree from sorted hashes of attestations chosen by the consensus bitVote.
// The computed tree is stored in the round.
// If any of the hash of the chosen attestations is missing, the tree is not computed.
func (r *Round) GetMerkleTree() (merkle.Tree, error) {

	r.sortAttestations()

	hashes := []common.Hash{}

	for i := range r.Attestations {
		if r.ConsensusBitVote.BitVector.Bit(i) == 1 {
			if r.Attestations[i].Status != Success {
				return merkle.Tree{}, errors.New("cannot build merkle tree")
			}
			hashes = append(hashes, r.Attestations[i].Hash)
		}
	}

	// sort.Slice(hashes, func(i, j int) bool { return compareHash(hashes[i], hashes[j]) })

	merkleTree := merkle.Build(hashes, false)

	r.merkletree = merkleTree

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

	if len(r.merkletree) != 0 {
		return r.merkletree, nil
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
