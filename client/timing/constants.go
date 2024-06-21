package timing

import (
	"flare-common/logger"
	"strings"
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

var ChainConstants Constants = costonTiming

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

func Set(chain string) {

	switch chainLC := strings.ToLower(chain); chainLC {
	case "coston":
		ChainConstants = costonTiming
	case "coston2":
		ChainConstants = coston2Timing
	case "songbird":
		ChainConstants = songbirdTiming
	case "flare":
		ChainConstants = flareTiming
	default:
		log.Warnf("No timing for chain %s. Switching to coston timing", chainLC)
		ChainConstants = costonTiming

	}

}
