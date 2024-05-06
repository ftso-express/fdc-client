package attestation

import (
	"errors"
	"flare-common/database"
	"flare-common/payload"
	"flare-common/policy"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"
)

type Manager struct {
	Rounds               map[uint64]*Round
	Timestamps           chan uint64
	Requests             <-chan []database.Log
	BitVotes             <-chan payload.Round
	RoundInCollect       uint64
	signingPolicyStorage policy.SigningPolicyStorage
	hub                  *hub.Hub
}

func (m *Manager) Run() {

	for {

		select {
		case round := <-m.BitVotes:

			for i := range round.Messages {

				m.OnBitVote(round.Messages[i])
			}

			m.Rounds[round.ID].ComputeConsensusBitVote()

		case requests := <-m.Requests:

			for i := range requests {

				m.OnRequest(requests[i])

			}

		}
	}
}

func (m *Manager) GetOrCreateRound(roundId uint64, status RoundStatus) (*Round, error) {

	round, ok := m.Rounds[roundId]

	if ok && round.roundId == roundId {
		round.status = status
		return round, nil
	}

	policy, _ := m.signingPolicyStorage.GetForVotingRound(uint32(roundId))

	if policy == nil {
		return nil, errors.New("no signing policy")
	}
	round = CreateRound(roundId, policy.Voters, status)

	m.Rounds[roundId] = round
	return round, nil
}

// OnBitVote process message that is assumed to be a bitVote
func (m *Manager) OnBitVote(message payload.Message) error {

	if message.Timestamp < timing.ChooseStartTimestamp(int(message.VotingRound)) {
		return errors.New("bitvote too soon")
	}

	if message.Timestamp > timing.ChooseEndTimestamp(int(message.VotingRound)) {
		return errors.New("bitvote too late")
	}

	round, err := m.GetOrCreateRound(message.VotingRound, Choosing)

	if err != nil {
		return err
	}

	err = round.ProcessBitVote(message)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) OnRequest(request database.Log) error {

	roundID := timing.RoundIDForTimestamp(request.Timestamp)

	attestation := Attestation{}

	attestation.RoundID = roundID

	data, err := ParseAttestationRequestLog(m.hub, request)

	if err != nil {
		return err
	}

	attestation.Request = data.Data

	attestation.Fee = data.Fee

	attestation.Status = Waiting

	attestation.Index.BlockNumber = request.BlockNumber
	attestation.Index.LogIndex = request.LogIndex

	round, err := m.GetOrCreateRound(roundID, Collecting)

	if err != nil {
		return err
	}

	round.Attestations = append(round.Attestations, &attestation)

	go func() {
		source, _ := attestation.Request.Source()
		attType, _ := attestation.Request.AttestationType()

		url, key := VerifierServer(attType, source)

		err = ResolveAttestationRequest(&attestation, url, key)

		(&attestation).VerifyResponse()
	}()

	return nil

}

func (m *Manager) OnSigningPolicy() {}

func VerifierServer(attType, source []byte) (string, string) {
	return "url", "key"
}
