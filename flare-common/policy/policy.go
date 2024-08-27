package policy

import (
	relayContract "flare-common/contracts/relay"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SigningPolicy struct {
	rewardEpochID      int64
	startVotingRoundID uint32
	threshold          uint16
	seed               *big.Int
	rawBytes           []byte
	blockTimestamp     uint64

	// The set of all Voters and their weights
	Voters *VoterSet
}

func NewSigningPolicy(r *relayContract.RelaySigningPolicyInitialized, submitToSigning map[common.Address]common.Address) *SigningPolicy {
	return &SigningPolicy{
		rewardEpochID:      r.RewardEpochId.Int64(),
		startVotingRoundID: r.StartVotingRoundId,
		threshold:          r.Threshold,
		seed:               r.Seed,
		rawBytes:           r.SigningPolicyBytes,
		blockTimestamp:     r.Timestamp,
		Voters:             NewVoterSet(r.Voters, r.Weights, submitToSigning),
	}
}
