package timing

import (
	"fmt"
)

func RoundIDForTimestamp(t uint64) (uint64, error) {
	if t+Chain.OffsetSec < Chain.T0 {
		return 0, fmt.Errorf("timestamp: %d before first round : %d", t, Chain.T0-Chain.OffsetSec)
	}

	roundID := (t + Chain.OffsetSec - Chain.T0) / Chain.CollectDurationSec

	return roundID, nil
}

func RoundStartTime(n uint64) uint64 {
	return Chain.T0 + n*Chain.CollectDurationSec - Chain.OffsetSec
}

func ChooseStartTimestamp(n uint64) uint64 {
	return RoundStartTime(n + 1)
}

func ChooseEndTimestamp(n uint64) uint64 {
	return ChooseStartTimestamp(n) + Chain.ChooseDurationSec
}

// NextChoosePhaseEnd returns the roundID of the round whose choose phase is next in line to end and the timestamp of the end.
// If t is right at the end of choose phase, the returned round is current and the timestamp is t.
func NextChooseEnd(t uint64) (uint64, uint64) {

	if t+Chain.OffsetSec < Chain.T0+Chain.ChooseDurationSec+1 {
		return 0, ChooseEndTimestamp(0)
	}

	roundID := (t - Chain.T0 + Chain.OffsetSec - Chain.ChooseDurationSec - 1) / Chain.CollectDurationSec

	endTimestamp := ChooseEndTimestamp(roundID)

	return roundID, endTimestamp
}

// LastCollectPhaseStart returns roundID and start timestamp of the latest round.
func LastCollectPhaseStart(t uint64) (uint64, uint64, error) {
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
