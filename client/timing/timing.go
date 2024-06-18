package timing

import (
	"flare-common/logger"
	"flare-common/paths"

	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	CollectDurationSec = 90
	ChooseDurationSec  = 30
	CommitDurationSec  = 20
	OffsetSec          = 30

	t0SGB     = 1658429955
	t0Coston2 = 1658430000

	testRewardEpochLength = 240
	rewardEpochLength     = 3360
)

var log = logger.GetLogger()

var ChainConstants Constants

type Constants struct {
	T0                uint64
	RewardEpochLength uint64 //number of voting rounds in one RewardEpoch
}

var costonTiming = Constants{
	RewardEpochLength: testRewardEpochLength,
	T0:                t0SGB,
}

var coston2Timing = Constants{
	RewardEpochLength: testRewardEpochLength,
	T0:                t0Coston2,
}

var songbirdTiming = Constants{
	RewardEpochLength: rewardEpochLength,
	T0:                t0SGB,
}

var flareTiming = Constants{
	RewardEpochLength: rewardEpochLength,
	T0:                t0SGB,
}

func init() {

	envPath, err := paths.LocalToAbsolute("../../.env")

	if err != nil {
		log.Panicf("timing: error getting path to .env: %w", err)

	}

	err = godotenv.Load(envPath)

	if err != nil {
		log.Panicf("timing: error reading .env: %w", err)

	}

	switch chain := os.Getenv("CHAIN"); chain {
	case "COSTON":
		ChainConstants = costonTiming
	case "COSTON2":
		ChainConstants = coston2Timing
	case "SONGBIRD":
		ChainConstants = songbirdTiming
	case "FLARE":
		ChainConstants = flareTiming
	default:
		ChainConstants = costonTiming

	}

}

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

func NextChoosePhaseEndPointers(t uint64) (*uint64, *uint64) {
	roundId := (t - ChainConstants.T0 + OffsetSec - ChooseDurationSec) / CollectDurationSec
	endTimestamp := ChooseEndTimestamp(roundId)

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
