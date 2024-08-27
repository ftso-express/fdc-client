package policy

import (
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
