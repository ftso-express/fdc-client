package attestation_test

import (
	"fmt"
	"local/fdc/client/attestation"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSkipDuplicates(t *testing.T) {

	tests := []struct {
		toAdd    []common.Hash
		atTheEnd []common.Hash
	}{
		{
			toAdd:    []common.Hash{common.HexToHash("1"), common.HexToHash("1"), common.HexToHash("2")},
			atTheEnd: []common.Hash{common.HexToHash("1"), common.HexToHash("2")},
		},
		{
			toAdd:    []common.Hash{},
			atTheEnd: []common.Hash(nil),
		},

		{
			toAdd:    []common.Hash{common.HexToHash("1")},
			atTheEnd: []common.Hash{common.HexToHash("1")},
		},

		{
			toAdd:    []common.Hash{common.HexToHash("2"), common.HexToHash("1")},
			atTheEnd: []common.Hash{common.HexToHash("2"), common.HexToHash("1")},
		},
		{
			toAdd:    []common.Hash{common.HexToHash("2"), common.HexToHash("1"), common.HexToHash("2")},
			atTheEnd: []common.Hash{common.HexToHash("2"), common.HexToHash("1")},
		},
	}

	for i, test := range tests {

		hashes := new([]common.Hash)

		added := make(attestation.CheckList)

		for j := range test.toAdd {

			attestation.SkipDuplicates(added, test.toAdd[j], hashes)

		}

		require.Equal(t, test.atTheEnd, *hashes, fmt.Sprintf("error in test %d", i))

	}

}
