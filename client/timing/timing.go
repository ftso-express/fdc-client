package timing

import "time"

const (
	collectTime = 90 * time.Second
	chooseTime  = 30 * time.Second
	commitTime  = 20 * time.Second
	offset      = 30 * time.Second
	t0          = 1704250616
)

type Round struct {
	Start time.Time
	ID    int
}

func GetRound(n int) *Round {
	return &Round{time.Unix(t0, 0).Add(collectTime*time.Duration(n) - offset), n}

}

func GetRoundForTimestamp(t uint64) *Round {

	round := int((t - t0 + 30) / 90)

	return GetRound(round)

}

func (r *Round) Next() *Round {

	return &Round{r.Start.Add(90 * time.Second), r.ID + 1}
}

func (r *Round) ToStart() time.Duration {
	return time.Until(r.Start)
}
