package attestation

import "errors"

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
	totalWeight      uint64
}

func Create(roundId, totalWeight uint64) *Round {
	var r Round

	r.roundId = roundId

	r.status = Collecting

	r.totalWeight = totalWeight

	return &r
}

func (r *Round) GetBitVote() (BitVote, error) {
	return BitVoteFromAttestations(r.attestations)
}

func (r *Round) ComputeConsensusBitVote() error {

	consensus, err := ConsensusBitVote(r.roundId, r.bitVotes, r.totalWeight, r.attestations)

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

	return bitVote.EncodeHex(r.roundId), nil
}

func (r *Round) GetMerkleRoot() {}

func (r *Round) GetMerkleTree() {

}

func (r *Round) GetConsensusBitVote() BitVote {

	return r.ConsensusBitVote
}
