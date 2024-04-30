package epoch

import "errors"

type Epoch struct {
	EpochID     uint64
	Weights     map[string]uint64
	TotalWeight uint64
	Start       uint64
}

type Manager struct {
	Epochs map[uint64]Epoch
}

func (m *Manager) GetEpochForRound(roundId uint64) (Epoch, error) {

	return Epoch{}, errors.New("wip")
}
