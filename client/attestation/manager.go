package attestation

import (
	"errors"
	"flare-common/contracts/relay"
	"flare-common/database"
	"flare-common/logger"
	"flare-common/payload"
	"flare-common/policy"
	"flare-common/storage"
	"fmt"
	"local/fdc/client/config"
	"local/fdc/client/timing"
	hub "local/fdc/contracts/FDC"

	"github.com/ethereum/go-ethereum/common"
)

// capacity of rounds cache
const roundBuffer uint64 = 256

var log = logger.GetLogger()

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
	Rounds               storage.Cyclic[*Round] // cyclically cached rounds with buffer roundBuffer.
	lastRoundCreated     uint64
	Requests             <-chan []database.Log
	BitVotes             <-chan payload.Round
	SigningPolicies      <-chan []database.Log
	signingPolicyStorage *policy.SigningPolicyStorage
	verifierServers      map[[64]byte]config.VerifierCredentials // the keys are AttestationTypeAndSource
	abiConfig            config.AbiConfig
}

// NewManager initializes attestation round manager
func NewManager(configs config.UserConfigRaw) *Manager {
	rounds := storage.NewCyclic[*Round](roundBuffer)
	signingPolicyStorage := policy.NewSigningPolicyStorage()

	abiConfig, err := config.ParseAbi(configs.Abis)

	if err != nil {
		log.Panic("parsing abis:", err)

	}

	verifierServers, err := config.ParseVerifiers(configs.Verifiers)

	if err != nil {
		log.Panic("parsing verifiers:", err)

	}

	return &Manager{Rounds: rounds, signingPolicyStorage: signingPolicyStorage, abiConfig: abiConfig, verifierServers: verifierServers}
}

// Run starts processing data received through the manager's channels.
func (m *Manager) Run() {
	// Get signing policy first as we cannot process any other message types
	// without a signing policy.
	signingPolicies := <-m.SigningPolicies
	log.Info("Initial signing policies received")

	for i := range signingPolicies {
		if err := m.OnSigningPolicy(signingPolicies[i]); err != nil {
			log.Panic("signing policy error:", err)
		}
	}

	for {
		select {
		case signingPolicies := <-m.SigningPolicies:

			log.Debug("New signing policy received")

			for i := range signingPolicies {

				if err := m.OnSigningPolicy(signingPolicies[i]); err != nil {
					log.Error("signing policy error:", err)
				}

			}
			deleted := m.signingPolicyStorage.RemoveBeforeVotingRound(uint32(m.lastRoundCreated)) // delete all signing policies that have already ended

			for j := range deleted {
				log.Infof("deleted signing policy for epoch %d", deleted[j])
			}

		case round := <-m.BitVotes:

			log.Debugf("Received %d bitVotes for round %d", len(round.Messages), round.Id)

			for i := range round.Messages {

				if err := m.OnBitVote(round.Messages[i]); err != nil {
					log.Error("bit vote error", err)
				}
			}

			r, ok := m.Rounds.Get(round.Id)
			if !ok {
				break
			}

			err := r.ComputeConsensusBitVote()

			if err != nil {
				log.Warnf("Failed bitVote in round %d: %s", round.Id, err)
			} else {
				log.Debugf("Consensus bitVote %s for round %d computed.", r.ConsensusBitVote.EncodeBitVoteHex(round.Id), round.Id)
			}

		case requests := <-m.Requests:

			log.Debugf("Received %d requests.", len(requests))

			for i := range requests {

				if err := m.OnRequest(requests[i]); err != nil {
					log.Error("requests error:", err)
				}

			}
		}
	}
}

// GetOrCreateRound returns a round for roundId either from manager if a round is already stored or creates a new one and stores it.
func (m *Manager) GetOrCreateRound(roundId uint64) (*Round, error) {

	round, ok := m.Rounds.Get(roundId)

	if ok {
		return round, nil
	}

	policy, _ := m.signingPolicyStorage.GetForVotingRound(uint32(roundId))

	if policy == nil {
		log.Errorf("No signing policy for round %d.", roundId)
		return nil, errors.New("no signing policy")
	}

	round = CreateRound(roundId, policy.Voters)
	m.lastRoundCreated = roundId
	log.Debugf("Round %d created", roundId)

	m.Rounds.Store(roundId, round)
	return round, nil
}

// Store stores round in to the cyclic cache

// OnBitVote process payload message that is assumed to be a bitVote and adds it to the correct round.
func (m *Manager) OnBitVote(message payload.Message) error {

	if message.Timestamp < timing.ChooseStartTimestamp(int(message.VotingRound)) {
		return fmt.Errorf("bitvote from %s too soon", message.From)
	}

	if message.Timestamp > timing.ChooseEndTimestamp(int(message.VotingRound)) {
		return fmt.Errorf("bitvote from %s too late", message.From)
	}

	round, err := m.GetOrCreateRound(message.VotingRound)

	if err != nil {
		log.Errorf("could not get round %d: %s", message.VotingRound, err)
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

	attestation.RoundId = roundID

	data, err := ParseAttestationRequestLog(request)

	if err != nil {
		log.Error("Error parsing attestation request")
		return err
	}

	attestation.Request = data.Data

	attestation.Fee = data.Fee

	attestation.Status = Waiting

	attestation.Index.BlockNumber = request.BlockNumber
	attestation.Index.LogIndex = request.LogIndex

	round, err := m.GetOrCreateRound(roundID)

	if err != nil {
		log.Error("Error getting or creating round")
		return err
	}

	round.Attestations = append(round.Attestations, &attestation)

	go func() {
		if err := m.handleAttestation(&attestation); err != nil {
			log.Error("Error handling attestation:", err)
		}
	}()

	return nil

}

func (m *Manager) handleAttestation(attestation *Attestation) error {
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

	var ok bool

	attestation.abi, ok = m.abiConfig.ResponseArguments[attType]

	if !ok {
		attestation.Status = UnsupportedPair
		return errors.New("unsupported att type: no abi")

	}

	verifier, ok := m.VerifierServer(attTypeAndSource)

	if !ok {
		attestation.Status = UnsupportedPair
		return errors.New("unsupported pair")

	}

	attestation.lutLimit = verifier.LutLimit

	attestation.Status = Processing

	err = ResolveAttestationRequest(attestation, verifier)

	if err != nil {
		log.Error("Error resolving attestation request")
		attestation.Status = ProcessError

		return err
	} else {
		err := attestation.validateResponse()
		return err
	}
}

// OnSigningPolicy parsed SigningPolicyInitialized log and stores it into the signingPolicyStorage.
func (m *Manager) OnSigningPolicy(initializedPolicy database.Log) error {

	data, err := ParseSigningPolicyInitializedLog(initializedPolicy)

	if err != nil {
		log.Errorf("Error parsing signing policy")
		return err
	}

	parsedPolicy := policy.NewSigningPolicy(data)

	log.Infof("Processing signing policy for rewardEpoch %s", data.RewardEpochId.String())

	err = m.signingPolicyStorage.Add(parsedPolicy)

	return err

}
