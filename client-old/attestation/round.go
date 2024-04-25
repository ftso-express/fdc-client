package attestation

import (
	"errors"

	"local/fdc/client/bitvote"
)

type RoundStatus int

const (
	Collect RoundStatus = iota
	Choose
	Retry
	Commit
	Reveal
	Finished
	Abandoned
)

type Round struct {
	RoundId          uint64
	Attestations     []attestation.Attestation
	BitVote          bitvote.BitVote
	WeightedBitVotes []bitvote.WeightedBitVote
	ConsensusBitVote bitvote.BitVote
	TotalWeight      uint64
	Status           Status
}

func (r *Round) SetBitVote() (bitvote.BitVote, error) {

	tempVote, err := bitvote.ForRound(r.Attestations)

	if err != nil {
		return bitvote.BitVote{}, errors.New("bitvote can not be computed")
	}

	r.BitVote = tempVote
	return tempVote, nil

}
func (r *Round) GetConsensusBitVote() bitvote.BitVote {

	r.ConsensusBitVote = bitvote.ConsensusBitVote(r.RoundId, r.WeightedBitVotes, r.TotalWeight, r.Attestations)

	return r.ConsensusBitVote
}

func (r *Round) SetConsensusStatus() {

	bitvote.SetBitVoteStatus(r.Attestations, r.ConsensusBitVote)
}

func (r *Round) RetryRequestUnsuccessfulButVoted() {

	for j := range r.Attestations {
		if r.Attestations[j].Consensus && r.Attestations[j].Status != attestation.Success {
			r.Attestations[j].Verify()
		}
	}
}

func (r *Round) MerkleTree() {}

func (r *Round) MerkleRoot() {}
