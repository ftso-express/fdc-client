package collector_test

import (
	"context"
	"local/fdc/client/collector"
	"sync"
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
	t0                    = 1658429955
	roundLengthSeconds    = 90
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

func TestPrepareChooseTriggers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := newMockCollectorDB()
	require.NoError(t, err)

	db.state.BlockTimestamp = t0 + roundLengthSeconds

	trigger := make(chan uint64)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		collector.PrepareChooseTriggers(ctx, trigger, db)
		wg.Done()
	}()

	select {
	case roundID := <-trigger:
		require.Equal(t, uint64(1), roundID)

	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}
