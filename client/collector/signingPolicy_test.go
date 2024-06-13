package collector_test

import (
	"context"
	"local/fdc/client/collector"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	spiLogName         = "SigningPolicyInitialized"
	rewardEpochID      = 1
	startVotingRoundID = 1
	threshold          = 0
	seed               = 0
	voterAddrHex       = "0xac872479e5EFc21989A4183Dc580C8264C9e54f5"
	timestamp          = 0
)

const (
	bufferSize           = 10
	relayContractAddrHex = "0x26A90DA287264E2E20a45d8c2c79Ca98439c5aa8"
)

var relayContractAddr = common.HexToAddress(relayContractAddrHex)

func TestSPIListener(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := newMockCollectorDB()
	require.NoError(t, err)

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
