package queue

import (
	"context"
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

func (q *PriorityQueue[T]) Enqueue(ctx context.Context, item T) error {
	select {
	case q.regular <- item:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *PriorityQueue[T]) EnqueuePriority(ctx context.Context, item T) error {
	select {
	case q.priority <- item:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *PriorityQueue[T]) Dequeue(ctx context.Context) (result T, err error) {
	if q.minDequeueDelta > 0 {
		if err = q.enforceRateLimit(ctx); err != nil {
			return result, err
		}

		defer func() {
			if err == nil {
				q.lastDequeue = time.Now()
			}
		}()
	}

	select {
	case result = <-q.priority:
		return result, nil

	default:
		select {
		case result = <-q.priority:
			return result, nil

		case result = <-q.regular:
			return result, nil

		case <-ctx.Done():
			// Set the err variable so the deferred function can read it.
			err = ctx.Err()
			return result, err
		}
	}
}

func (q *PriorityQueue[T]) enforceRateLimit(ctx context.Context) error {
	now := time.Now()
	delta := now.Sub(q.lastDequeue)
	if delta >= q.minDequeueDelta {
		return nil
	}

	sleepDuration := q.minDequeueDelta - delta
	log.Debugf("enforcing rate limit - sleeping for %s", sleepDuration)

	select {
	case <-time.After(sleepDuration):
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
