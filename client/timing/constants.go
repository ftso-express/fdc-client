package timing

import (
	"flare-common/logger"
	"flare-common/paths"
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
		log.Error("timing: error reading .env: %w", err)

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
