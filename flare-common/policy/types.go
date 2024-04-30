package policy

import (
	relayContract "flare-common/contacts/relay"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type VoterData struct {
	index  int
	weight uint16
}

type VoterSet struct {
	voters      []common.Address
	weights     []uint16
	totalWeight uint16
	thresholds  []uint16

	voterDataMap map[common.Address]VoterData
}

func NewVoterSet(voters []common.Address, weights []uint16) *VoterSet {
	vs := VoterSet{
		voters:     voters,
		weights:    weights,
		thresholds: make([]uint16, len(weights)),
	}
	// sum does not exceed uint16, guaranteed by the smart contract
	for i, w := range weights {
		vs.thresholds[i] = vs.totalWeight
		vs.totalWeight += w
	}

	vMap := make(map[common.Address]VoterData)
	for i, voter := range vs.voters {
		if _, ok := vMap[voter]; !ok {
			vMap[voter] = VoterData{
				index:  i,
				weight: vs.weights[i],
			}
		}
	}
	vs.voterDataMap = vMap
	return &vs
}

type signingPolicy struct {
	rewardEpochId      int64
	startVotingRoundId uint32
	threshold          uint16
	seed               *big.Int
	rawBytes           []byte
	blockTimestamp     uint64

	// The set of all voters and their weights
	voters *VoterSet
}

func newSigningPolicy(r *relayContract.RelaySigningPolicyInitialized) *signingPolicy {
	return &signingPolicy{
		rewardEpochId:      r.RewardEpochId.Int64(),
		startVotingRoundId: r.StartVotingRoundId,
		threshold:          r.Threshold,
		seed:               r.Seed,
		rawBytes:           r.SigningPolicyBytes,
		blockTimestamp:     r.Timestamp,

		voters: NewVoterSet(r.Voters, r.Weights),
	}
}

type signingPolicyStorage struct {

	// sorted list of signing policies, sorted by rewardEpochId (and also by startVotingRoundId)
	spList []*signingPolicy

	// mutex
	sync.Mutex
}

func newSigningPolicyStorage() *signingPolicyStorage {
	return &signingPolicyStorage{
		spList: make([]*signingPolicy, 0, 10),
	}
}
