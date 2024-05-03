package timing

import "time"

const (
	collectTime = 90 * time.Second
	chooseTime  = 30 * time.Second
	commitTime  = 20 * time.Second
	offset      = 30 * time.Second
	t0          = 1704250616
)

func GetRoundIDForTimestamp(t uint64) uint64 {

	roundID := uint64((t - t0 + 30) / 90)

	return roundID
}

func GetRoundStartTime(n int) time.Time {
	return time.Unix(t0, 0).Add(collectTime*time.Duration(n) - offset)
}

func GetRoundStartTimestamp(n int) uint64 {
	return uint64(GetRoundStartTime(n).Unix())
}

func GetChooseStartTimestamp(n int) uint64 {
	return uint64(GetRoundStartTime(n).Add(collectTime).Unix())
}

func GetChooseEndTimestamp(n int) uint64 {
	return uint64(GetRoundStartTime(n).Add(collectTime + chooseTime).Unix())
}

func NextChoosePhaseEnd(t uint64) (*int, *uint64) {
	roundID := int((t - t0) / 90)
	endTimestamp := uint64(t0 + (roundID+1)*90)

	return &roundID, &endTimestamp
}
