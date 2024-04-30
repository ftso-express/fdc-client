package attestation

import (
	"errors"
	"flare-common/merkle"
	"local/fdc/client/epoch"

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
	attestations     []*Attestation
	bitVotes         []WeightedBitVote
	ConsensusBitVote BitVote
	epoch            epoch.Epoch
	merkletree       merkle.Tree
}

func CreateRound(r *Round, roundId uint64, epoch epoch.Epoch, status RoundStatus) *Round {

	r.roundId = roundId

	r.status = status

	r.epoch = epoch

	return r
}

func (r *Round) GetBitVote() (BitVote, error) {
	return BitVoteFromAttestations(r.attestations)
}

func (r *Round) ComputeConsensusBitVote() error {

	consensus, err := ConsensusBitVote(r.roundId, r.bitVotes, r.epoch.TotalWeight, r.attestations)

	if err != nil {
		return err

	}
	r.ConsensusBitVote = consensus

	return nil
}

func (r *Round) GetBitVoteHex() (string, error) {

	bitVote, err := BitVoteFromAttestations(r.attestations)

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

	// handle no bitVote or chosen request that is not registered

	for i := range r.attestations {
		r.attestations[i].Consensus = r.ConsensusBitVote.BitVector.Bit(int(r.attestations[i].Index)) == 1
	}

	return nil

}

func (r *Round) GetMerkleTree() (merkle.Tree, error) {

	hashes := []common.Hash{}

	for i := range r.attestations {
		if r.ConsensusBitVote.BitVector.Bit(int(r.attestations[i].Index)) == 1 {
			if r.attestations[i].Status != Success {
				return merkle.Tree{}, errors.New("cannot build merkle tree")
			}
			hashes = append(hashes, r.attestations[i].Hash)
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
