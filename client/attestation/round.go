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

func CreateRound(roundId uint64, voterSet *policy.VoterSet, status RoundStatus) *Round {

	r := &Round{}

	r.roundId = roundId

	r.status = status

	r.voterSet = voterSet

	return r
}

func (r *Round) SortAttestations() {

	sort.Slice(r.Attestations, func(i, j int) bool {
		return Less(r.Attestations[i].Index, r.Attestations[j].Index)
	})
}

func (r *Round) SortBitVotes() {

	sort.Slice(r.bitVotes, func(i, j int) bool {
		return r.bitVotes[i].Index < r.bitVotes[j].Index
	})
}

func (r *Round) GetBitVote() (BitVote, error) {

	r.SortAttestations()
	return BitVoteFromAttestations(r.Attestations)
}

func (r *Round) ComputeConsensusBitVote() error {

	r.SortBitVotes()

	consensus, err := ConsensusBitVote(r.roundId, r.bitVotes, r.voterSet.TotalWeight, r.Attestations)

	if err != nil {
		return err

	}
	r.ConsensusBitVote = consensus

	return nil
}

func (r *Round) GetBitVoteHex() (string, error) {

	r.SortAttestations()

	bitVote, err := BitVoteFromAttestations(r.Attestations)

	if err != nil {
		return "", errors.New("cannot get bitvote")
	}

	return bitVote.EncodeBitVoteHex(r.roundId), nil
}

func (r *Round) GetMerkleRoot() {}

func (r *Round) GetConsensusBitVote() BitVote {

	return r.ConsensusBitVote
}

func (r *Round) SetConsensusStatus() error {

	r.SortAttestations()

	// handle no bitVote or chosen request that is not registered

	for i := range r.Attestations {
		r.Attestations[i].Consensus = r.ConsensusBitVote.BitVector.Bit(i) == 1
	}

	return nil

}

func (r *Round) GetMerkleTree() (merkle.Tree, error) {

	r.SortAttestations()

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
