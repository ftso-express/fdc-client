package shared

import (
	"flare-common/database"
	"flare-common/payload"
	"flare-common/storage"
	"local/fdc/client/round"
)

const (
	bitVoteBufferSize              = 2
	requestsBufferSize             = 10
	signingPolicyBufferSize        = 3
	roundBuffer             uint64 = 256
)

type SharedDataPipes struct {
	Rounds          storage.Cyclic[*round.Round] // cyclically cached rounds with buffer roundBuffer.
	Requests        chan []database.Log
	BitVotes        chan payload.Round
	SigningPolicies chan []database.Log
}

func NewSharedDataPipes() *SharedDataPipes {
	return &SharedDataPipes{
		Rounds:          storage.NewCyclic[*round.Round](roundBuffer),
		SigningPolicies: make(chan []database.Log, signingPolicyBufferSize),
		BitVotes:        make(chan payload.Round, bitVoteBufferSize),
		Requests:        make(chan []database.Log, requestsBufferSize),
	}
}
