package queue

import (
	"flare-common/logger"
	"time"
)

var log = logger.GetLogger()

// PriorityQueue is made up of two sub-queues - one regular and one with
// higher priority. Items can be enqueued in either queue and when dequeueing
// items from the priority queue are returned first.
type PriorityQueue[T any] struct {
	regular         chan T
	priority        chan T
	minDequeueDelta time.Duration
	lastDequeue     time.Time
}

func NewPriority[T any](size int, maxDequeuesPerSecond int) PriorityQueue[T] {
	q := PriorityQueue[T]{
		regular:  make(chan T, size),
		priority: make(chan T, size),
	}

	if maxDequeuesPerSecond > 0 {
		q.minDequeueDelta = time.Second / time.Duration(maxDequeuesPerSecond)
		log.Info("minDequeueDelta:", q.minDequeueDelta)
	}

	return q
}

func (q *PriorityQueue[T]) Enqueue(item T) {
	q.regular <- item
}

func (q *PriorityQueue[T]) EnqueuePriority(item T) {
	q.priority <- item
}

func (q *PriorityQueue[T]) Dequeue() T {
	if q.minDequeueDelta > 0 {
		q.enforceRateLimit()

		defer func() {
			q.lastDequeue = time.Now()
		}()
	}

	select {
	case item := <-q.priority:
		return item

	default:
		select {
		case item := <-q.priority:
			return item

		case item := <-q.regular:
			return item
		}
	}
}

func (q *PriorityQueue[T]) enforceRateLimit() {
	now := time.Now()
	delta := now.Sub(q.lastDequeue)
	if delta < q.minDequeueDelta {
		sleepDuration := q.minDequeueDelta - delta
		log.Infof("enforcing rate limit - sleeping for %s", sleepDuration)

		time.Sleep(sleepDuration)
	}
}
