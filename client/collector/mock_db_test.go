package collector_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"flare-common/contracts/relay"
	"flare-common/database"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type mockCollectorDB struct {
	logs []database.Log
	txs  []database.Transaction
}

func (c mockCollectorDB) FetchState(ctx context.Context) (database.State, error) {
	return database.State{Index: lastQueriedBlock}, nil
}

func (c mockCollectorDB) FetchLatestLogsByAddressAndTopic0(
	ctx context.Context, addr common.Address, topic0 common.Hash, num int,
) ([]database.Log, error) {
	if addr != relayContractAddr {
		return nil, errors.New("unknown address")
	}

	return c.logs, nil
}

func (c mockCollectorDB) FetchLogsByAddressAndTopic0Timestamp(
	ctx context.Context, addr common.Address, topic0 common.Hash, from, to int64,
) ([]database.Log, error) {
	if addr != relayContractAddr {
		return nil, errors.New("unknown address")
	}

	return c.logs, nil
}

func (c mockCollectorDB) FetchLogsByAddressAndTopic0BlockNumber(
	ctx context.Context,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]database.Log, error) {
	return c.logs, nil
}

func (c mockCollectorDB) FetchLogsByAddressAndTopic0TimestampToBlockNumber(
	ctx context.Context,
	address common.Address,
	topic0 common.Hash,
	from, to int64,
) ([]database.Log, error) {
	return c.logs, nil
}

func (c mockCollectorDB) FetchTransactionsByAddressAndSelectorTimestamp(
	ctx context.Context,
	toAddress common.Address,
	functionSel [4]byte,
	from, to int64,
) ([]database.Transaction, error) {
	if toAddress != submitContractAddr {
		return nil, errors.New("unknown address")
	}

	if functionSel != funcSel {
		return nil, errors.New("unknown function selector")
	}

	return c.txs, nil
}

func newMockCollectorDB() (*mockCollectorDB, error) {
	log, err := newTestLog()
	if err != nil {
		return nil, err
	}

	tx, err := newTestTx()
	if err != nil {
		return nil, err
	}

	return &mockCollectorDB{logs: []database.Log{*log}, txs: []database.Transaction{*tx}}, nil
}

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

func newTestTx() (*database.Transaction, error) {
	var b bytes.Buffer

	b.Write(funcSel[:])
	b.WriteByte(protocol)

	if err := binary.Write(&b, binary.BigEndian, uint32(roundID)); err != nil {
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, uint16(len(payload))); err != nil {
		return nil, err
	}

	b.Write(payload)

	return &database.Transaction{Input: hex.EncodeToString(b.Bytes())}, nil
}
