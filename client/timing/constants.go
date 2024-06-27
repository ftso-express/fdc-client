package timing

import (
	"errors"
	"local/fdc/client/config"
)

const (
	CollectDurationSec = 90
	ChooseDurationSec  = 30
	CommitDurationSec  = 20
	OffsetSec          = 30
)

const (
	defaultT0                = 1658429955 // SGB and Coston
	defaultRewardEpochLength = 240        //Coston
)

var Chain config.Timing = config.Timing{
	T0:                defaultT0,
	RewardEpochLength: defaultRewardEpochLength,
}

func Set(chainTiming config.Timing) error {

	if chainTiming.RewardEpochLength == 0 {
		return errors.New("illegal reward epoch length")
	}

	Chain = chainTiming

	return nil

}
