package timing_test

import (
	"fmt"
	"local/fdc/client/timing"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundIdForTimestamp(t *testing.T) {

	_, err := timing.RoundIdForTimestamp(0)

	require.Error(t, err)

	tests := []struct {
		timestamp uint64
		roundId   uint64
	}{
		{
			timestamp: timing.ChainConstants.T0 - timing.OffsetSec,
			roundId:   0,
		},
		{
			timestamp: timing.ChainConstants.T0 + 10000*timing.CollectDurationSec - timing.OffsetSec/2,
			roundId:   10000,
		},
	}

	for i, test := range tests {
		roundId, err := timing.RoundIdForTimestamp(test.timestamp)

		require.NoError(t, err, fmt.Sprintf("unexpected error in test %d: %s", i, err))
		require.Equal(t, test.roundId, roundId, fmt.Sprintf("wrong round in test %d", i))
	}
}
