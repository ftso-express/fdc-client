package attestation

import (
	"errors"
	"flare-common/contacts/relay"
	"flare-common/database"
	"flare-common/payload"
	"flare-common/policy"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"

	"github.com/ethereum/go-ethereum/common"
)

// capacity of rounds cache
const roundBuffer uint64 = 256

// hubFilterer is only used for Attestation Requests logs parsing. Set in init().
var hubFilterer *hub.HubFilterer

// relayFilterer is only used for SigningPolicyInitialized logs parsing. Set in init()
var relayFilterer *relay.RelayFilterer

// init sets the hubFilterer and relayFilterer.
func init() {

	hubFilterer, _ = hub.NewHubFilterer(common.Address{}, nil)

	relayFilterer, _ = relay.NewRelayFilterer(common.Address{}, nil)
}

type Manager struct {
	rounds               map[uint64]*Round //cyclically cached rounds with buffer roundBuffer.
	Requests             <-chan []database.Log
	BitVotes             <-chan payload.Round
	SigningPolicies      <-chan []database.Log
	signingPolicyStorage policy.SigningPolicyStorage
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Run() {

	for {

		select {
		case round := <-m.BitVotes:

			for i := range round.Messages {

				m.OnBitVote(round.Messages[i])
			}

			r, ok := m.Round(round.ID)

			if !ok {
				break
			}
			r.ComputeConsensusBitVote()

		case requests := <-m.Requests:

			for i := range requests {

				m.OnRequest(requests[i])

			}

		case signingPolicies := <-m.SigningPolicies:

			for i := range signingPolicies {

				m.OnSigningPolicy(signingPolicies[i])

			}
		}
	}
}

// GetOrCreateRound returns a round for roundId either from manager if a round is already stored or creates a new one and stores it.
func (m *Manager) GetOrCreateRound(roundId uint64) (*Round, error) {

	round, ok := m.Round(roundId)

	if ok && round.roundId == roundId {
		return round, nil
	}

	policy, _ := m.signingPolicyStorage.GetForVotingRound(uint32(roundId))

	if policy == nil {
		return nil, errors.New("no signing policy")
	}
	round = CreateRound(roundId, policy.Voters)

	m.Store(round)
	return round, nil
}

// Round returns a round for roundId stored by the Manager. If round is not stored, false is returned.
func (m *Manager) Round(roundId uint64) (*Round, bool) {

	roundReminder := roundId / roundBuffer

	round, ok := m.rounds[roundReminder]

	if ok && round.roundId == roundId {
		return round, true
	}

	return nil, false
}

// Store stores round in to the cyclic cache
func (m *Manager) Store(round *Round) {

	roundReminder := round.roundId / roundBuffer

	m.rounds[roundReminder] = round
}

// OnBitVote process payload message that is assumed to be a bitVote and adds it to the correct round.
func (m *Manager) OnBitVote(message payload.Message) error {

	if message.Timestamp < timing.ChooseStartTimestamp(int(message.VotingRound)) {
		return errors.New("bitvote too soon")
	}

	if message.Timestamp > timing.ChooseEndTimestamp(int(message.VotingRound)) {
		return errors.New("bitvote too late")
	}

	round, err := m.GetOrCreateRound(message.VotingRound)

	if err != nil {
		return err
	}

	err = round.ProcessBitVote(message)

	if err != nil {
		return err
	}

	return nil
}

// OnRequest process the attestation request.
// The request parsed into an Attestation that is assigned to an attestation round according to the timestamp.
// The request is sent to verifier server and the verifier's response is validated.
func (m *Manager) OnRequest(request database.Log) error {

	roundID := timing.RoundIDForTimestamp(request.Timestamp)

	attestation := Attestation{}

	attestation.RoundID = roundID

	data, err := ParseAttestationRequestLog(request)

	if err != nil {
		return err
	}

	attestation.Request = data.Data

	attestation.Fee = data.Fee

	attestation.Status = Waiting

	attestation.Index.BlockNumber = request.BlockNumber
	attestation.Index.LogIndex = request.LogIndex

	round, err := m.GetOrCreateRound(roundID)

	if err != nil {
		return err
	}

	round.Attestations = append(round.Attestations, &attestation)

	go func() {
		source, _ := attestation.Request.Source()
		attType, _ := attestation.Request.AttestationType()

		url, key := VerifierServer(attType, source)

		attestation.Status = Processing

		err = ResolveAttestationRequest(&attestation, url, key)

		if err != nil {
			attestation.Status = ProcessError
		} else {

			attestation.validateResponse()
		}
	}()

	return nil

}

// OnSigningPolicy parsed SigningPolicyInitialized log and stores it into the signingPolicyStorage.
func (m *Manager) OnSigningPolicy(initializedPolicy database.Log) error {

	data, err := ParseSigningPolicyInitializedLog(initializedPolicy)

	if err != nil {
		return err
	}

	parsedPolicy := policy.NewSigningPolicy(data)

	err = m.signingPolicyStorage.Add(parsedPolicy)

	return err

}

// VerifierServer retrieves url and credentials for the verifier's server for the pair of attType and source.
func VerifierServer(attType, source []byte) (string, string) {
	return "url", "key"
}
