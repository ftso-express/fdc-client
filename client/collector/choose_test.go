package collector_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flare-common/database"
	"local/fdc/client/collector"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type mockCollectorDBBitVote struct {
	txs []database.Transaction
}

func (c mockCollectorDBBitVote) FetchLatestLogsByAddressAndTopic0(
	ctx context.Context, addr common.Address, topic0 common.Hash, num int,
) ([]database.Log, error) {
	return nil, errors.New("no implemented")
}

func (c mockCollectorDBBitVote) FetchLogsByAddressAndTopic0Timestamp(
	ctx context.Context, addr common.Address, topic0 common.Hash, from, to int64,
) ([]database.Log, error) {
	return nil, errors.New("no implemented")
}

func (c mockCollectorDBBitVote) FetchTransactionsByAddressAndSelectorTimestamp(
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

func newMockCollectorDBBitVote() (*mockCollectorDBBitVote, error) {
	tx, err := newTestTx()
	if err != nil {
		return nil, err
	}

	return &mockCollectorDBBitVote{txs: []database.Transaction{*tx}}, nil
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

const (
	submitContractAddrHex = "0x90C6423ec3Ea40591bAdb177171B64c7e6556028"
	protocol              = 0xff
	roundID               = 1
)

var (
	submitContractAddr = common.HexToAddress(submitContractAddrHex)
	funcSel            = [4]byte{1, 2, 3, 4}
	payload            = []byte{1, 2, 3, 4, 5, 6, 7, 8}
)

func TestBitVoteListener(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := newMockCollectorDBBitVote()
	require.NoError(t, err)

	trigger := make(chan uint64)

	out := collector.BitVoteListener(
		ctx,
		db,
		submitContractAddr,
		funcSel,
		protocol,
		bufferSize,
		trigger,
	)

	select {
	case trigger <- roundID:
		t.Log("sent roundID to trigger")

	case <-ctx.Done():
		t.Fatal("context cancelled")
	}

	select {
	case round := <-out:
		require.Equal(t, uint64(roundID), round.Id)
		require.Len(t, round.Messages, 1)
		cupaloy.SnapshotT(t, round.Messages[0])

	case <-ctx.Done():
		t.Fatal("context cancelled")
	}

}
