package attestation

import (
	"errors"
	"flare-common/contacts/relay"
	"flare-common/database"
	"flare-common/payload"
	"flare-common/policy"
	"local/fdc/client/timing"
	"local/fdc/config"
	hub "local/fdc/contracts/FDC"
	"log"

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
	signingPolicyStorage *policy.SigningPolicyStorage
	verifierServers      map[[64]byte]config.VerifierCredentials // the keys are crypto.Keccak256Hash(AttestationTypeAndSource)
	abiConfig            config.AbiConfig
}

// NewManager initializes attestation round manager
func NewManager() *Manager {
	rounds := make(map[uint64]*Round)
	signingPolicyStorage := policy.NewSigningPolicyStorage()
	return &Manager{rounds: rounds, signingPolicyStorage: signingPolicyStorage}
}

// Run starts processing data received through the manager's channels.
// SigningPolicies channel is prioritized.
func (m *Manager) Run() {

	for {

		select {
		case signingPolicies := <-m.SigningPolicies:

			log.Println("New signing policy received.")

			for i := range signingPolicies {

				m.OnSigningPolicy(signingPolicies[i])

			}

		default:
			{
				select {
				case signingPolicies := <-m.SigningPolicies:

					log.Println("New signing policy received.")

					for i := range signingPolicies {

						m.OnSigningPolicy(signingPolicies[i])

					}
				case round := <-m.BitVotes:

					log.Printf("Received %d bitVotes for round %d.", len(round.Messages), round.ID)

					for i := range round.Messages {

						m.OnBitVote(round.Messages[i])
					}

					r, ok := m.Round(round.ID)

					if !ok {
						break
					}
					err := r.ComputeConsensusBitVote()

					if err != nil {
						log.Printf("Failed bitVote for round %d", round.ID)
					} else {
						log.Printf("Consensus bitVote for round %d computed.", round.ID)
					}

				case requests := <-m.Requests:

					log.Printf("Received %d requests.", len(requests))

					for i := range requests {

						m.OnRequest(requests[i])

					}
				}
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
		log.Printf("No signing policy for round %d.", roundId)
		return nil, errors.New("no signing policy")
	}

	round = CreateRound(roundId, policy.Voters)
	log.Printf("Round %d created.", roundId)

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

	log.Println("Processing request")
	// log.Println(request)

	roundID := timing.RoundIDForTimestamp(request.Timestamp)

	attestation := Attestation{}

	attestation.RoundID = roundID

	data, err := ParseAttestationRequestLog(request)

	if err != nil {
		log.Println("Error parsing attestation request")
		return err
	}

	attestation.Request = data.Data

	attestation.Fee = data.Fee

	attestation.Status = Waiting

	attestation.Index.BlockNumber = request.BlockNumber
	attestation.Index.LogIndex = request.LogIndex

	round, err := m.GetOrCreateRound(roundID)

	if err != nil {
		log.Println("Error getting or creating round")
		return err
	}

	round.Attestations = append(round.Attestations, &attestation)

	go func() error {
		attTypeAndSource, err := attestation.Request.AttestationTypeAndSource()

		if err != nil {

			attestation.Status = ProcessError
			return err

		}

		attType, err := attestation.Request.AttestationType()

		if err != nil {

			attestation.Status = ProcessError
			return err

		}

		attestation.abi = m.abiConfig.ResponseArguments[attType]

		verifier, ok := m.VerifierServer(attTypeAndSource)

		if !ok {
			attestation.Status = UnsupportedPair
			return errors.New("unsupported pair")

		}

		attestation.Status = Processing

		err = ResolveAttestationRequest(&attestation, verifier)

		if err != nil {
			log.Println("Error resolving attestation request")
			attestation.Status = ProcessError

			return err
		} else {
			log.Println("Response received, validating...")
			err := attestation.validateResponse()
			log.Println(attestation.Status, attestation.RoundID)
			return err
		}
	}()

	return nil

}

// OnSigningPolicy parsed SigningPolicyInitialized log and stores it into the signingPolicyStorage.
func (m *Manager) OnSigningPolicy(initializedPolicy database.Log) error {
	log.Println("Processing signing policy")

	data, err := ParseSigningPolicyInitializedLog(initializedPolicy)

	if err != nil {
		log.Println("Error parsing signing policy")
		return err
	}

	parsedPolicy := policy.NewSigningPolicy(data)

	err = m.signingPolicyStorage.Add(parsedPolicy)

	return err

}
