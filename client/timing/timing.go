package timing

import "time"

const (
	collectTime = 90 * time.Second
	chooseTime  = 30 * time.Second
	commitTime  = 20 * time.Second
	offset      = 30 * time.Second
	// TODO: Luka - get this from the config
	t0 = 1658429955
)

func RoundIDForTimestamp(t uint64) uint64 {

	roundID := uint64((t - t0 + 30) / 90)

	return roundID
}

func RoundStartTime(n int) time.Time {
	return time.Unix(t0, 0).Add(collectTime*time.Duration(n) - offset)
}

func RoundStartTimestamp(n int) uint64 {
	return uint64(RoundStartTime(n).Unix())
}

func ChooseStartTimestamp(n int) uint64 {
	return uint64(RoundStartTime(n).Add(collectTime).Unix())
}

func ChooseEndTimestamp(n int) uint64 {
	return uint64(RoundStartTime(n).Add(collectTime + chooseTime).Unix())
}

func NextChoosePhaseEnd(t uint64) (int, uint64) {
	roundID := int((t - t0) / 90)
	endTimestamp := uint64(t0 + (roundID+1)*90)

	return roundID, endTimestamp
}

func NextChoosePhaseEndPointers(t uint64) (*int, *uint64) {
	roundID := int((t - t0) / 90)
	endTimestamp := uint64(t0 + (roundID+1)*90)

	return &roundID, &endTimestamp
}

func LastCollectPhaseStart(t uint64) (int, uint64) {
	roundID := RoundIDForTimestamp(t)

	startTimestamp := t0 + roundID*90 - 30

	return int(roundID), startTimestamp
}

func ExpectedRewardEpochStartTimestamp(rewardEpochId uint64) uint64 {
	return t0 + 240*90*rewardEpochId // TODO get this from config. currently for coston?
}
