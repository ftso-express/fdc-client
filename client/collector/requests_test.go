package collector_test

import (
	"context"
	"local/fdc/client/collector"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const lastQueriedBlock = 123

var (
	fdcContractAddr  = common.HexToAddress("0xf26Be97eB0d7a9fBf8d67f813D3Be411445885ce")
	listenerInterval = time.Millisecond
)

func TestAttestationRequestListener(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := newMockCollectorDB()
	require.NoError(t, err)

	out := collector.AttestationRequestListener(
		ctx,
		db,
		fdcContractAddr,
		bufferSize,
		listenerInterval,
	)

	for i := 0; i < 2; i++ {
		select {
		case logs := <-out:
			require.Len(t, logs, 1)
			require.Equal(t, db.logs[0], logs[0])

		case <-ctx.Done():
			t.Fatal("context cancelled")
		}
	}
}
