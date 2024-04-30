package attestation

import (
	"errors"
	"flare-common/database"
	"flare-common/payload"
	"local/fdc/client/epoch"
	hub "local/fdc/contracts/FDC"
)

type Manager struct {
	Rounds         map[uint64]*Round
	Timestamps     chan uint64
	Requests       chan []database.Log
	BitVotes       chan payload.Message
	RoundInCollect uint64
	epochManager   epoch.Manager
	hub            *hub.Hub
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

func (m *Manager) OnRequests(requests []database.Log) error {

	for i := range requests {

		round, err := getRoundIdForTimestamp(requests[i].Timestamp)

		if err != nil {
			break
		}

		attestation := Attestation{}

		attestation.RoundID = round

		data, err := ParseAttestationRequestLog(m.hub, requests[i])

		if err != nil {
			break
		}

		attestation.Request = data.Data

		attestation.Fee = data.Fee

		attestation.Status = Waiting

		m.Rounds[round].Attestations = append(m.Rounds[round].Attestations, &attestation)
	}
	return nil

}

func getRoundIdForTimestamp(timestamp uint64) (uint64, error) {
	return 0, errors.New("Wip")
}
