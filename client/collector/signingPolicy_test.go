package collector_test

import (
	"context"
	"flare-common/database"
	"local/fdc/client/collector"
	"local/fdc/tests/mocks"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	relayContractAddrHex = "0x26A90DA287264E2E20a45d8c2c79Ca98439c5aa8"
)

var relayContractAddr = common.HexToAddress(relayContractAddrHex)

func TestSPIListener(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mocks.NewMockCollectorDB()
	require.NoError(t, err)

	signingPolicyChan := make(chan []database.Log, 2)
	go collector.SigningPolicyInitializedListener(ctx, db, relayContractAddr, signingPolicyChan)

	// Two iterations: first when the listener is first initialised and second
	// from the queryNextSPI iteration.
	for i := 0; i < 2; i++ {
		select {
		case logs := <-signingPolicyChan:
			require.Len(t, logs, 1)
			require.Equal(t, db.Logs[0], logs[0])

		case <-ctx.Done():
			t.Fatal("timed out")
		}
	}
}
