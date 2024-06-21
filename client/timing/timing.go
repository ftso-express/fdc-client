package timing

import (
	"fmt"
)

func RoundIdForTimestamp(t uint64) (uint64, error) {

	if t+OffsetSec < ChainConstants.T0 {
		return 0, fmt.Errorf("timestamp: %d before first round : %d", t, ChainConstants.T0-OffsetSec)
	}

	roundId := (t + OffsetSec - ChainConstants.T0) / CollectDurationSec

	return roundId, nil
}

func RoundStartTime(n uint64) uint64 {

	return ChainConstants.T0 + n*CollectDurationSec - OffsetSec
}

func ChooseStartTimestamp(n uint64) uint64 {
	return RoundStartTime(n + 1)
}

func ChooseEndTimestamp(n uint64) uint64 {
	return ChooseStartTimestamp(n) + ChooseDurationSec
}

// NextChoosePhaseEnd returns the roundId of the round whose choose phase is next in line to end and the timestamp of the end.
func NextChooseEnd(t uint64) (uint64, uint64) {

	if t+OffsetSec < ChainConstants.T0+ChooseDurationSec {
		return 0, ChooseEndTimestamp(0)
	}

	roundId := (t - ChainConstants.T0 + OffsetSec - ChooseDurationSec) / CollectDurationSec

	endTimestamp := ChooseEndTimestamp(roundId)

	return roundId, endTimestamp
}

func NextChoosePhasePtr(t uint64) (*uint64, *uint64) {

	roundId, endTimestamp := NextChooseEnd(t)

	return &roundId, &endTimestamp
}

// LastCollectPhaseStart returns roundId and start timestamp of the latest round.
func LastCollectPhaseStart(t uint64) (uint64, uint64, error) {
	roundId, err := RoundIdForTimestamp(t)

	if err != nil {
		return 0, 0, err
	}

	startTimestamp := RoundStartTime(roundId)

	return roundId, startTimestamp, nil
}

// ExpectedRewardEpochStartTimestamp returns the expected timestamp of the rewardEpoch with rewardEpochId.
func ExpectedRewardEpochStartTimestamp(rewardEpochId uint64) uint64 {
	return ChainConstants.T0 + ChainConstants.RewardEpochLength*ChooseDurationSec*rewardEpochId
}
