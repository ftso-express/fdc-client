package timing

import (
	"fmt"
)

func RoundIdForTimestamp(t uint64) (uint64, error) {
	if t+Chain.OffsetSec < Chain.T0 {
		return 0, fmt.Errorf("timestamp: %d before first round : %d", t, Chain.T0-Chain.OffsetSec)
	}

	roundId := (t + Chain.OffsetSec - Chain.T0) / Chain.CollectDurationSec

	return roundId, nil
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

// NextChoosePhaseEnd returns the roundId of the round whose choose phase is next in line to end and the timestamp of the end.
// If t is right at the end of choose phase, the returned round is current and the timestamp is t.
func NextChooseEnd(t uint64) (uint64, uint64) {

	if t+Chain.OffsetSec < Chain.T0+Chain.ChooseDurationSec+1 {
		return 0, ChooseEndTimestamp(0)
	}

	roundId := (t - Chain.T0 + Chain.OffsetSec - Chain.ChooseDurationSec - 1) / Chain.CollectDurationSec

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
	return Chain.T0 + Chain.T0RewardDelay + Chain.RewardEpochLength*Chain.CollectDurationSec*rewardEpochId
}

// ExpectedRewardEpochStartTimestamp returns the expected timestamp of the rewardEpoch with rewardEpochId.
func ExpectedVotingEpochStartTimestamp(votingEpoch uint64) uint64 {
	return Chain.T0 + Chain.CollectDurationSec*votingEpoch
}
