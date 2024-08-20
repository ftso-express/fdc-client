package policy

import (
	relayContract "flare-common/contracts/relay"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type VoterData struct {
	Index  int
	Weight uint16
}

type VoterSet struct {
	voters      []common.Address //signingPolicyAddress
	weights     []uint16
	TotalWeight uint16
	thresholds  []uint16

	VoterDataMap           map[common.Address]VoterData //signingPolicyAddressToWeight
	SubmitToSigningAddress map[common.Address]common.Address
}

func NewVoterSet(voters []common.Address, weights []uint16, SubmitToSigningAddress map[common.Address]common.Address) *VoterSet {
	vs := VoterSet{
		voters:     voters,
		weights:    weights,
		thresholds: make([]uint16, len(weights)),
	}
	// sum does not exceed uint16, guaranteed by the smart contract
	for i, w := range weights {
		vs.thresholds[i] = vs.TotalWeight
		vs.TotalWeight += w
	}

	vMap := make(map[common.Address]VoterData)
	for i, voter := range vs.voters {
		log.Debugf("New voter: %v", voter)

		if _, ok := vMap[voter]; !ok {
			vMap[voter] = VoterData{
				Index:  i,
				Weight: vs.weights[i],
			}
		}
	}
	vs.VoterDataMap = vMap

	vs.SubmitToSigningAddress = SubmitToSigningAddress
	return &vs
}

type SigningPolicy struct {
	rewardEpochId      int64
	startVotingRoundId uint32
	threshold          uint16
	seed               *big.Int
	rawBytes           []byte
	blockTimestamp     uint64

	// The set of all Voters and their weights
	Voters *VoterSet
}

func NewSigningPolicy(r *relayContract.RelaySigningPolicyInitialized, submitToSigning map[common.Address]common.Address) *SigningPolicy {
	return &SigningPolicy{
		rewardEpochId:      r.RewardEpochId.Int64(),
		startVotingRoundId: r.StartVotingRoundId,
		threshold:          r.Threshold,
		seed:               r.Seed,
		rawBytes:           r.SigningPolicyBytes,
		blockTimestamp:     r.Timestamp,
		Voters:             NewVoterSet(r.Voters, r.Weights, submitToSigning),
	}
}

type SigningPolicyStorage struct {

	// sorted list of signing policies, sorted by rewardEpochId (and also by startVotingRoundId)
	spList []*SigningPolicy

	// mutex
	sync.Mutex
}

func NewSigningPolicyStorage() *SigningPolicyStorage {
	return &SigningPolicyStorage{
		spList: make([]*SigningPolicy, 0, 10),
	}
}
