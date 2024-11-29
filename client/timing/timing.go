package timing

import (
	"fmt"
)

// RoundIDForTimestamp calculates roundID that is active at timestamp.
//
// j-th round is active in [T0 + j * CollectDurationSec, T0 + (j+1)* CollectDurationSec)
func RoundIDForTimestamp(t uint64) (uint32, error) {
	if t < Chain.T0 {
		return 0, fmt.Errorf("timestamp: %d before first round : %d", t, Chain.T0)
	}

	roundID := (t - Chain.T0) / Chain.CollectDurationSec

	return uint32(roundID), nil
}

func RoundStartTime(n uint32) uint64 {
	return Chain.T0 + uint64(n)*Chain.CollectDurationSec
}

func ChooseStartTimestamp(n uint32) uint64 {
	return RoundStartTime(n + 1)
}

func ChooseEndTimestamp(n uint32) uint64 {
	return ChooseStartTimestamp(n) + Chain.ChooseDurationSec
}

// NextChoosePhaseEnd returns the roundID of the round whose choose phase is next in line to end and the timestamp of the end.
// If t is right at the end of choose phase, the returned round is current and the timestamp is t.
func NextChooseEnd(t uint64) (uint32, uint64) {
	if t < Chain.T0+Chain.ChooseDurationSec+1 {
		return 0, ChooseEndTimestamp(0)
	}

	roundID := (t - Chain.T0 - Chain.ChooseDurationSec - 1) / Chain.CollectDurationSec

	endTimestamp := ChooseEndTimestamp(uint32(roundID))

	return uint32(roundID), endTimestamp
}

// LastCollectPhaseStart returns roundID and start timestamp of the latest round.
func LastCollectPhaseStart(t uint64) (uint32, uint64, error) {
	roundID, err := RoundIDForTimestamp(t)

	if err != nil {
		return 0, 0, err
	}

	startTimestamp := RoundStartTime(roundID)

	return roundID, startTimestamp, nil
}

// ExpectedRewardEpochStartTimestamp returns the expected start timestamp of the rewardEpoch with rewardEpochID.
func ExpectedRewardEpochStartTimestamp(rewardEpochID uint64) uint64 {
	return Chain.T0 + Chain.T0RewardDelay + Chain.RewardEpochLength*Chain.CollectDurationSec*rewardEpochID
}
