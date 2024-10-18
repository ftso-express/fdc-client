package shared

import (
	"github.com/flare-foundation/go-flare-common/pkg/contracts/relay"
	"github.com/flare-foundation/go-flare-common/pkg/database"
	"github.com/flare-foundation/go-flare-common/pkg/payload"
	"github.com/flare-foundation/go-flare-common/pkg/storage"

	"github.com/flare-foundation/fdc-client/client/round"

	"github.com/ethereum/go-ethereum/common"
)

const (
	bitVoteBufferSize           = 2
	requestsBufferSize          = 10
	signingPolicyBufferSize     = 3
	roundBuffer             int = 256
)

type VotersData struct {
	Policy                 *relay.RelaySigningPolicyInitialized
	SubmitToSigningAddress map[common.Address]common.Address
}

type DataPipes struct {
	Rounds   storage.Cyclic[uint32, *round.Round] // cyclically cached rounds with buffer roundBuffer.
	Requests chan []database.Log
	BitVotes chan payload.Round
	Voters   chan []VotersData
}

func NewDataPipes() *DataPipes {
	return &DataPipes{
		Rounds:   storage.NewCyclic[uint32, *round.Round](roundBuffer),
		Voters:   make(chan []VotersData, signingPolicyBufferSize),
		BitVotes: make(chan payload.Round, bitVoteBufferSize),
		Requests: make(chan []database.Log, requestsBufferSize),
	}
}
