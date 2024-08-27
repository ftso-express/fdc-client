package manager

import (
	"context"
	"flare-common/database"
	"flare-common/logger"
	"flare-common/payload"
	"flare-common/policy"
	"flare-common/storage"
	"fmt"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	"local/fdc/client/round"
	"local/fdc/client/shared"
	"local/fdc/client/timing"
)

var log = logger.GetLogger()

type Manager struct {
	Rounds                storage.Cyclic[uint32, *round.Round] // cyclically cached rounds with buffer roundBuffer.
	lastRoundCreated      uint32
	requests              <-chan []database.Log
	bitVotes              <-chan payload.Round
	signingPolicies       <-chan []shared.VotersData
	signingPolicyStorage  *policy.SigningPolicyStorage
	attestationTypeConfig config.AttestationTypes
	queues                priorityQueues
}

// New initializes attestation round manager from raw user configurations.
func New(configs *config.UserRaw, attestationTypeConfig config.AttestationTypes, sharedDataPipes *shared.DataPipes) (*Manager, error) {
	signingPolicyStorage := policy.NewStorage()

	queues := buildQueues(configs.Queues)

	return &Manager{
			Rounds:                sharedDataPipes.Rounds,
			signingPolicyStorage:  signingPolicyStorage,
			attestationTypeConfig: attestationTypeConfig,
			queues:                queues,
			signingPolicies:       sharedDataPipes.Voters,
			bitVotes:              sharedDataPipes.BitVotes,
			requests:              sharedDataPipes.Requests,
		},
		nil
}

// Run starts processing data received through the manager's channels.
func (m *Manager) Run(ctx context.Context) {
	// Get signing policy first as we cannot process any other message types
	// without a signing policy.
	var signingPolicies []shared.VotersData

	go runQueues(ctx, m.queues)

	select {
	case signingPolicies = <-m.signingPolicies:
		log.Info("Initial signing policies received")

	case <-ctx.Done():
		log.Info("Manager exiting:", ctx.Err())
		return
	}

	for i := range signingPolicies {
		if err := m.OnSigningPolicy(signingPolicies[i]); err != nil {
			log.Panic("signing policy error:", err)
		}
	}

	for {
		select {
		case signingPolicies := <-m.signingPolicies:
			log.Debug("New signing policy received")

			for i := range signingPolicies {
				err := m.OnSigningPolicy(signingPolicies[i])
				if err != nil {
					log.Error("signing policy error:", err)
				}
			}
			deleted := m.signingPolicyStorage.RemoveBefore(uint32(m.lastRoundCreated)) // delete all signing policies that have already ended

			for j := range deleted {
				log.Infof("deleted signing policy for epoch %d", deleted[j])
			}

		case bitVotesForRound := <-m.bitVotes:
			log.Debugf("Received %d bitVotes for round %d", len(bitVotesForRound.Messages), bitVotesForRound.ID)

			for i := range bitVotesForRound.Messages {
				err := m.OnBitVote(bitVotesForRound.Messages[i])
				if err != nil {
					log.Errorf("bit vote error: %s", err)
				}
			}
			r, ok := m.Rounds.Get(bitVotesForRound.ID)
			if !ok {
				break
			}

			err := r.ComputeConsensusBitVote()
			if err != nil {
				log.Warnf("Failed bitVote in round %d: %s", bitVotesForRound.ID, err)
			} else {
				log.Debugf("Consensus bitVote %s for round %d computed.", r.ConsensusBitVote.EncodeBitVoteHex(bitVotesForRound.ID), bitVotesForRound.ID)

				noOfRetried, err := m.retryUnsuccessfulChosen(ctx, r)
				if err != nil {
					log.Warnf("error retrying round %d: s", r.ID, err)
				} else if noOfRetried > 0 {
					log.Debugf("retrying %d attestations in round %d", noOfRetried, r.ID)
				}

			}

		case requests := <-m.requests:
			log.Debugf("Received %d requests.", len(requests))

			for i := range requests {
				err := m.OnRequest(ctx, requests[i])
				if err != nil {
					log.Error(err)
				}
			}

		case <-ctx.Done():
			log.Info("Manager exiting:", ctx.Err())
			return
		}
	}
}

// GetOrCreateRound returns a round for roundID either from manager if a round is already stored or creates a new one and stores it.
func (m *Manager) GetOrCreateRound(roundID uint32) (*round.Round, error) {
	roundForID, ok := m.Rounds.Get(roundID)
	if ok {
		return roundForID, nil
	}

	policy, _ := m.signingPolicyStorage.ForVotingRound(roundID)
	if policy == nil {
		return nil, fmt.Errorf("creating round: no signing policy for round %d", roundID)
	}

	roundForID = round.New(roundID, policy.Voters)
	m.lastRoundCreated = roundID
	log.Debugf("Round %d created", roundID)

	m.Rounds.Store(roundID, roundForID)
	return roundForID, nil
}

// OnBitVote process payload message that is assumed to be a bitVote and adds it to the correct round.
func (m *Manager) OnBitVote(message payload.Message) error {
	if message.Timestamp < timing.ChooseStartTimestamp(message.VotingRound) {
		return fmt.Errorf("bitVote from %s for voting round %d too soon", message.From, message.VotingRound)
	}

	if message.Timestamp >= timing.ChooseEndTimestamp(message.VotingRound) {
		return fmt.Errorf("bitVote from %s for voting round %d too late", message.From, message.VotingRound)
	}

	round, err := m.GetOrCreateRound(message.VotingRound)
	if err != nil {
		return err
	}

	err = round.ProcessBitVote(message)
	if err != nil {
		return fmt.Errorf("error processing bitVote from %s for voting round %d: %s", message.From, message.VotingRound, err)
	}

	return nil
}

// OnRequest process the attestation request.
// The request parsed into an Attestation that is assigned to an attestation round according to the timestamp.
// The request is added to verifier queue.
func (m *Manager) OnRequest(ctx context.Context, request database.Log) error {
	attestation, err := attestation.AttestationFromDatabaseLog(request)
	if err != nil {
		return fmt.Errorf("OnRequest: %s", err)
	}

	round, err := m.GetOrCreateRound(attestation.RoundID)
	if err != nil {
		return fmt.Errorf("OnRequest: %s", err)
	}

	added := round.AddAttestation(&attestation)
	if added {
		if err := m.AddToQueue(ctx, &attestation); err != nil {
			return err
		}
	}

	return nil
}

// OnSigningPolicy parses SigningPolicyInitialized log and stores it into the signingPolicyStorage.
func (m *Manager) OnSigningPolicy(data shared.VotersData) error {
	parsedPolicy := policy.NewSigningPolicy(data.Policy, data.SubmitToSigningAddress)
	log.Infof("Processing signing policy for rewardEpoch %s", data.Policy.RewardEpochId.String())

	err := m.signingPolicyStorage.Add(parsedPolicy)

	return err
}

// retryUnsuccessfulChosen adds the request that were chosen by the consensus bitVote but were not confirmed to the priority verifier queues.
func (m *Manager) retryUnsuccessfulChosen(ctx context.Context, round *round.Round) (int, error) {
	count := 0 //only for logging

	for i := range round.Attestations {
		if round.Attestations[i].Consensus && round.Attestations[i].Status != attestation.Success {
			queueName := round.Attestations[i].QueueName

			queue, ok := m.queues[queueName]
			if !ok {
				return 0, fmt.Errorf("retry: no queue: %s", queueName)
			}

			err := queue.EnqueuePriority(ctx, round.Attestations[i])
			if err != nil {
				return 0, err
			}

			count++
		}
	}

	return count, nil
}

// AddToQueue adds the attestation to the correct verifier queue.
func (m *Manager) AddToQueue(ctx context.Context, attestation *attestation.Attestation) error {
	err := attestation.PrepareRequest(m.attestationTypeConfig)
	if err != nil {
		return fmt.Errorf("preparing request: %s", err)
	}

	queue, ok := m.queues[attestation.QueueName]
	if !ok {
		return fmt.Errorf("queue %s does not exist", attestation.QueueName)
	}

	err = queue.Enqueue(ctx, attestation)

	return err
}
