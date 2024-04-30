package attestation

import (
	"flare-common/payload"
	"local/fdc/client/epoch"
)

type Manager struct {
	Rounds         map[uint64]*Round
	Timestamps     chan uint64
	Requests       chan uint64
	BitVotes       chan payload.Message
	RoundInCollect uint64
	epochManager   epoch.Manager
}

func (m *Manager) GetOrCreateRound(roundId uint64, status RoundStatus) (*Round, error) {

	round := m.Rounds[roundId]

	if round.roundId == roundId {
		round.status = status
		return round, nil
	}

	epoch, err := m.epochManager.GetEpochForRound(roundId)

	if err != nil {
		return nil, err
	}
	round = CreateRound(round, roundId, epoch, status)

	m.Rounds[roundId] = round
	return round, nil
}

func (m *Manager) OnBitVote(message payload.Message) error {
	// TODO check timestamp
	round, err := m.GetOrCreateRound(message.VotingRound, Choosing)

	if err != nil {
		return err
	}

	weightedBitVote, err := ProcessBitVote(message, round.epoch)

	if err != nil {
		return err
	}

	round.bitVotes = append(round.bitVotes, weightedBitVote)

	return nil
}
