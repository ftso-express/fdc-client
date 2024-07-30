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

var (
	fdcContractAddr  = common.HexToAddress("0xf26Be97eB0d7a9fBf8d67f813D3Be411445885ce")
	listenerInterval = time.Millisecond
)

func TestAttestationRequestListener(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mocks.NewMockCollectorDB()
	require.NoError(t, err)

	requestChan := make(chan []database.Log, 10)

	go collector.AttestationRequestListener(
		ctx,
		db,
		fdcContractAddr,
		listenerInterval,
		requestChan,
	)

	for i := 0; i < 2; i++ {
		select {
		case logs := <-requestChan:
			require.Len(t, logs, 1)
			require.Equal(t, db.Logs[0], logs[0])

		case <-ctx.Done():
			t.Fatal("context cancelled")
		}
	}
}
