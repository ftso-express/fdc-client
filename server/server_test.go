package server_test

import (
	"context"
	"encoding/hex"
	"flare-common/policy"
	"flare-common/storage"
	"local/fdc/client/attestation"
	bitvotes "local/fdc/client/attestation/bitVotes"
	"local/fdc/client/config"
	"local/fdc/client/round"
	"local/fdc/server"
	"local/fdc/tests/mocks"
	"math/big"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	votingRoundID = 1
	submitAddress = "0xf4Bf90cf71F52b4e0369a356D1F871A6237AD0C4"
)

func TestServer(t *testing.T) {
	rounds := storage.NewCyclic[*round.Round](10)
	serverConfig := config.RestServer{
		Title:       "FDC protocol data provider API",
		FSPTitle:    "FDC protocol data provider for FSP client",
		FSPSubpath:  "/fsp",
		Version:     "0.0.0",
		SwaggerPath: "/api-doc",
		Addr:        "localhost:8080",
		ApiKeyName:  "X-API-KEY",
		ApiKeys:     []string{"12345", "123456"},
	}

	s := server.New(&rounds, 200, serverConfig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go s.Run(ctx)
	defer s.Shutdown()

	hash := common.HexToHash("0x232")

	round := round.CreateRound(votingRoundID, policy.NewVoterSet(nil, nil, nil))
	round.Attestations = append(round.Attestations, &attestation.Attestation{
		RoundId:   votingRoundID,
		Consensus: true,
		Status:    attestation.Success,
		Hash:      hash,
	})
	rounds.Store(votingRoundID, round)

	bitVote := bitvotes.BitVote{Length: 1, BitVector: big.NewInt(1)}

	round.ConsensusBitVote = bitVote

	//Wait for the server to be ready.
	u := url.URL{Scheme: "http", Host: "localhost:8080", Path: "/health"}
	healthURL := u.String()

	require.Eventually(
		t,
		func() bool {
			rsp, err := http.Get(healthURL)
			if err != nil {
				return false
			}

			return rsp.StatusCode == http.StatusOK
		},
		10*time.Second,
		100*time.Millisecond,
	)

	t.Run("submit1", func(t *testing.T) {
		rspData, err := mocks.MakeGetRequest("submit1", &serverConfig, votingRoundID, submitAddress)
		require.NoError(t, err)

		t.Log(rspData)
		require.Equal(t, server.OK, rspData.Status)
		cupaloy.SnapshotT(t, rspData)
	})

	var submitString string

	t.Run("submit2", func(t *testing.T) {
		rspData, err := mocks.MakeGetRequest("submit2", &serverConfig, votingRoundID, submitAddress)
		require.NoError(t, err)

		t.Log(rspData)
		require.Equal(t, server.OK, rspData.Status)

		submitString = rspData.Data

		require.Equal(t, "0x", submitString[0:2])

	})

	t.Run("submitSignatures", func(t *testing.T) {
		rspData, err := mocks.MakeGetRequest("submitSignatures", &serverConfig, votingRoundID, submitAddress)
		require.NoError(t, err)

		t.Log(rspData)
		require.Equal(t, server.OK, rspData.Status)

		require.Equal(t, "0xc80000000101", rspData.Data[:14])

		require.Equal(t, hash.Hex()[2:], rspData.Data[14:])

		random := rspData.AdditionalData[2:66]

		consensusBitVote := rspData.AdditionalData[66:]

		consensusBitVoteBytes, err := hex.DecodeString(consensusBitVote)
		require.NoError(t, err)

		commitCheck := server.CalculateMaskedRoot(common.HexToHash(rspData.Data), common.HexToHash(random), common.HexToAddress(submitAddress), consensusBitVoteBytes)

		require.Equal(t, submitString[16:], commitCheck)

	})
}
