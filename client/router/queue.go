package router

import (
	"local/fdc/client/attestation"
	"time"
)

type RatedQueue struct {
	queue         chan attestation.Attestation
	priorityQueue chan attestation.Attestation
	perSec        uint64
	lastTime      int64
	limitPerSec   uint64
}

func (q *RatedQueue) Push(a attestation.Attestation) {
	q.queue <- a
}

func (q *RatedQueue) PriorityPush(a attestation.Attestation) {
	q.priorityQueue <- a
}

func (q *RatedQueue) canPop() bool {
	now := time.Now().Unix()
	if q.lastTime == now && q.perSec < q.limitPerSec {
		return true
	}
	if q.lastTime < now {
		return true
	}
	return false
}

func (q *RatedQueue) Pop()
