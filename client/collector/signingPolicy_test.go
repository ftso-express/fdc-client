package collector_test

import (
	"context"
	"encoding/hex"
	"flare-common/contracts/relay"
	"flare-common/database"
	"local/fdc/client/collector"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type mockCollectorDB struct {
	logs []database.Log
}

func (c mockCollectorDB) FetchLatestLogsByAddressAndTopic0(
	ctx context.Context, addr common.Address, topic0 common.Hash, num int,
) ([]database.Log, error) {
	return c.logs, nil
}

func (c mockCollectorDB) FetchLogsByAddressAndTopic0Timestamp(
	ctx context.Context, addr common.Address, topic0 common.Hash, from, to int64,
) ([]database.Log, error) {
	return c.logs, nil
}

func newMockCollectorDB() (*mockCollectorDB, error) {
	log, err := newTestLog()
	if err != nil {
		return nil, err
	}

	return &mockCollectorDB{logs: []database.Log{*log}}, nil
}

const (
	spiLogName         = "SigningPolicyInitialized"
	rewardEpochID      = 1
	startVotingRoundID = 1
	threshold          = 0
	seed               = 0
	voterAddrHex       = "0xac872479e5EFc21989A4183Dc580C8264C9e54f5"
	timestamp          = 0
)

func newTestLog() (*database.Log, error) {
	relayABI, err := relay.RelayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	event, ok := relayABI.Events[spiLogName]
	if !ok {
		return nil, errors.Errorf("event %s not found in ABI", spiLogName)
	}

	var indexedArgs abi.Arguments
	for i := range event.Inputs {
		if event.Inputs[i].Indexed {
			indexedArgs = append(indexedArgs, event.Inputs[i])
		}
	}

	if len(indexedArgs) != 1 {
		return nil, errors.Errorf("unexpected number of indexed args: %d %+v", len(indexedArgs), indexedArgs)
	}

	topic1, err := indexedArgs.Pack(big.NewInt(rewardEpochID))
	if err != nil {
		return nil, errors.Wrap(err, "packing topic1")
	}

	voters := []common.Address{common.HexToAddress(voterAddrHex)}
	weights := []uint16{1}
	signingPolicyBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	eventData, err := event.Inputs.NonIndexed().Pack(
		uint32(startVotingRoundID),
		uint16(threshold),
		big.NewInt(seed),
		voters,
		weights,
		signingPolicyBytes,
		uint64(timestamp),
	)
	if err != nil {
		return nil, errors.Wrap(err, "packing eventData")
	}

	return &database.Log{
		Data:   hex.EncodeToString(eventData),
		Topic0: hex.EncodeToString(event.ID[:]),
		Topic1: hex.EncodeToString(topic1[:]),
		Topic2: "NULL",
		Topic3: "NULL",
	}, nil
}

const (
	bufferSize           = 10
	relayContractAddrHex = "0x26A90DA287264E2E20a45d8c2c79Ca98439c5aa8"
)

func TestSPIListener(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	db, err := newMockCollectorDB()
	require.NoError(t, err)

	relayContractAddr := common.HexToAddress(relayContractAddrHex)
	out := collector.SigningPolicyInitializedListener(ctx, db, relayContractAddr, bufferSize)

	// Two iterations: first when the listener is first initialised and second
	// from the queryNextSPI iteration.
	for i := 0; i < 2; i++ {
		select {
		case logs := <-out:
			require.Len(t, logs, 1)
			require.Equal(t, db.logs[0], logs[0])

		case <-ctx.Done():
			t.Fatal("timed out")
		}
	}
}
