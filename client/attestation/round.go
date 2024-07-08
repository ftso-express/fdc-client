package attestation

import (
	"errors"
	"flare-common/merkle"
	"flare-common/policy"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	attestationMap   map[common.Hash]*Attestation
	bitVotes         []*WeightedBitVote
	bitVoteCheckList map[common.Address]*WeightedBitVote
	ConsensusBitVote BitVote
	voterSet         *policy.VoterSet
	merkleTree       merkle.Tree
}

// CreateRound returns a pointer to a new round with roundId and voterSet.
func CreateRound(roundId uint64, voterSet *policy.VoterSet) *Round {

	r := &Round{roundId: roundId, voterSet: voterSet, attestationMap: make(map[common.Hash]*Attestation)}

	return r
}

// addAttestation checks whether an attestation with such request is already in the round.
// If not it is added to the round, if yes the fee is added to the existent attestation
// and Index is set to the earlier one.
func (r *Round) addAttestation(attestation *Attestation) bool {
	identifier := crypto.Keccak256Hash(attestation.Request)

	att, exists := r.attestationMap[identifier]

	if exists {

		att.Fee = att.Fee.Add(att.Fee, attestation.Fee)

		if earlierLog(attestation.index(), att.index()) {

			att.Indexes = prepend(att.Indexes, attestation.index())
		}

		att.Indexes = append(att.Indexes, attestation.index())

		return false
	}

	r.attestationMap[identifier] = attestation

	r.Attestations = append(r.Attestations, attestation)

	return true
}

// prepend places the element at the beginning of the slice and moves the potentially replaced element to the end.
func prepend[T any](slice []T, element T) []T {

	if len(slice) == 0 {
		slice = append(slice, element)
		return slice
	}

	slice = append(slice, slice[0])

	slice[0] = element

	return slice

}

// sortAttestations sorts round's attestations according to their IndexLog.
// we assume that attestations have at least one index.
func (r *Round) sortAttestations() {

	sort.Slice(r.Attestations, func(i, j int) bool {
		return earlierLog(r.Attestations[i].index(), r.Attestations[j].index())
	})
}

// sortBitVotes sorts round's bitVotes according to the signingPolicy Index of their providers.
func (r *Round) sortBitVotes() {

	sort.Slice(r.bitVotes, func(i, j int) bool {
		return r.bitVotes[i].Index < r.bitVotes[j].Index
	})
}

// BitVoteHex returns the hex string encoded BitVote for the round according to the current status of Attestations.
func (r *Round) BitVoteHex() (string, error) {

	r.sortAttestations()

	bitVote, err := BitVoteFromAttestations(r.Attestations)

	if err != nil {
		return "", fmt.Errorf("cannot get bitVote for round %d: %w", r.roundId, err)
	}

	return bitVote.EncodeBitVoteHex(r.roundId), nil
}

// ComputeConsensusBitVote computes the consensus BitVote according to the collected bitVotes and sets consensus status to the attestations.
func (r *Round) ComputeConsensusBitVote(protocolId uint64) error {

	r.sortBitVotes()

	r.sortAttestations()

	consensus, err := ConsensusBitVote(&ConsensusBitVoteInput{
		RoundID:          r.roundId,
		WeightedBitVotes: r.bitVotes,
		TotalWeight:      r.voterSet.TotalWeight,
		Attestations:     r.Attestations,
	},
		protocolId)
	if err != nil {
		return err

	}
	r.ConsensusBitVote = consensus

	return r.setConsensusStatus(consensus)
}

// GetConsensusBitVote returns consensus BitVote if it is already computed.
func (r *Round) GetConsensusBitVote() (BitVote, error) {

	if r.ConsensusBitVote.BitVector == nil {
		return BitVote{}, errors.New("no consensus bitVote")
	}

	return r.ConsensusBitVote, nil
}

// setConsensusStatus sets consensus status of the attestations.
// The scenario where a chosen attestation is missing is not possible as in such case, it is not possible to compute the consensus bitVote.
// It is assumed that the Attestations are already ordered.
func (r *Round) setConsensusStatus(consensusBitVote BitVote) error {

	// sanity check
	if consensusBitVote.BitVector.BitLen() > len(r.Attestations) {
		return fmt.Errorf("missing attestation for round %d", r.roundId)
	}

	for i := range r.Attestations {
		r.Attestations[i].Consensus = consensusBitVote.BitVector.Bit(i) == 1
	}

	return nil

}

// MerkleTree computes Merkle tree from sorted hashes of attestations chosen by the consensus bitVote.
// The computed tree is stored in the round.
// If any of the hash of the chosen attestations is not successfully verified, the tree is not computed.
func (r *Round) MerkleTree() (merkle.Tree, error) {

	r.sortAttestations()

	hashes := []common.Hash{}

	for i := range r.Attestations {
		if r.Attestations[i].Consensus {
			if r.Attestations[i].Status != Success {
				return merkle.Tree{}, errors.New("cannot build merkle tree")
			}

			hashes = append(hashes, r.Attestations[i].Hash)

		}
	}

	merkleTree := merkle.Build(hashes, false)

	r.merkleTree = merkleTree

	return merkleTree, nil

}

// MerkleTreeCached gets Merkle tree from cache if it is already computed or computes it.
func (r *Round) MerkleTreeCached() (merkle.Tree, error) {

	if len(r.merkleTree) != 0 {
		return r.merkleTree, nil
	}

	return r.MerkleTree()

}

// MerkleRoot returns Merkle root for a round if there is one.
func (r *Round) MerkleRoot() (common.Hash, error) {

	tree, err := r.MerkleTreeCached()

	if err != nil {
		return common.Hash{}, err
	}

	return tree.Root()

}

// MerkleRootHex returns Merkle root for a round as a hex string.
func (r *Round) MerkleRootHex() (string, error) {

	root, err := r.MerkleRoot()

	if err != nil {
		return "", err
	}

	return root.Hex(), nil

}
