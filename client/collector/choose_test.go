package collector_test

import (
	"context"
	"local/fdc/client/collector"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

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

	db, err := newMockCollectorDB()
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
