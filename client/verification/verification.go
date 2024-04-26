package verification

import (
	"encoding/binary"
	"errors"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
)

type Request []byte

type Response []byte

func (r Request) GetMic() ([]byte, error) {
	if len(r) < 96 {
		return []byte{}, errors.New("request is to short")
	}

	return r[64:96], nil
}

func (r Response) ComputeMic() ([]byte, error) {
	if len(r) < 128 {
		return []byte{}, errors.New("response is to short")
	}

	d := make([]byte, 32)

	slices.Replace(d, 0, 32, r[64:96]...)

	zero32bytes := make([]byte, 32)

	slices.Replace(r, 64, 96, zero32bytes...)

	mic := crypto.Keccak256(r)

	defer slices.Replace(r, 64, 96, d...)

	return mic, nil

}

func (r Response) AddRound(roundId uint64) ([]byte, error) {
	if len(r) < 128 {
		return []byte{}, errors.New("response is to short")
	}

	end := make([]byte, 32)
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf, roundId)

	slices.Replace(end, 24, 32, buf...)

	slices.Replace(r, 64, 96, end...)

	return r, nil

}
