package server_test

import (
	"encoding/hex"
	"flare-common/policy"
	"flare-common/storage"
	"fmt"
	"local/fdc/client/attestation"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"local/fdc/client/config"
	"local/fdc/client/round"
	"local/fdc/server"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func makeController(t *testing.T) server.DAController {

	rounds := storage.NewCyclic[*round.Round](10)

	controller := server.DAController{Rounds: &rounds}

	hash := common.HexToHash("0x232")

	request, err := hex.DecodeString(requestEVM)

	require.NoError(t, err)

	response, err := hex.DecodeString(responseEVM)

	require.NoError(t, err)

	abi, abiString, err := config.GetAbi("../tests/configs/abis/EVMTransaction.json")

	require.NoError(t, err)

	round := round.CreateRound(1, policy.NewVoterSet(nil, nil, nil))
	round.Attestations = append(round.Attestations, &attestation.Attestation{
		Request:   request,
		Response:  response,
		RoundId:   1,
		Consensus: true,
		Status:    attestation.Success,
		Hash:      hash,
		Abi:       &abi,
		AbiString: &abiString,
	})
	rounds.Store(votingRoundID, round)

	bitVote := bitvotes.BitVote{Length: 1, BitVector: big.NewInt(1)}

	round.ConsensusBitVote = bitVote

	return controller
}

func TestGetRequests(t *testing.T) {
	controller := makeController(t)

	requests, ok := controller.GetRequests(1)

	require.True(t, ok)

	require.Len(t, requests, 1)

	requests, ok = controller.GetRequests(2)

	require.True(t, !ok)

	require.Nil(t, requests)

}

func TestGetAttestations(t *testing.T) {
	controller := makeController(t)

	attestations, ok := controller.GetAttestations(1)

	require.True(t, ok)

	require.Len(t, attestations, 1)

	fmt.Printf("attestations: %v\n", attestations[0].Abi)
}
