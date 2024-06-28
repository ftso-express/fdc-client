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
	workersSem      chan struct{}
}

// PriorityQueueParams values are used to construct a new PriorityQueue.
type PriorityQueueParams struct {
	Size                 int `toml:"size"`
	MaxDequeuesPerSecond int `toml:"max_dequeues_per_second"` // Set to 0 to disable rate-limiting
	MaxWorkers           int `toml:"max_workers"`             // Set to 0 for unlimited workers
}

// NewPriority constructs a new PriorityQueue.
func NewPriority[T any](input *PriorityQueueParams) PriorityQueue[T] {
	if input == nil {
		input = new(PriorityQueueParams)
	}

	q := PriorityQueue[T]{
		regular:  make(chan T, input.Size),
		priority: make(chan T, input.Size),
	}

	if input.MaxDequeuesPerSecond > 0 {
		q.minDequeueDelta = time.Second / time.Duration(input.MaxDequeuesPerSecond)
		log.Info("minDequeueDelta:", q.minDequeueDelta)
	}

	if input.MaxWorkers > 0 {
		q.workersSem = make(chan struct{}, input.MaxWorkers)
	}

	return q
}

// Enqueue adds an item to the queue with regular priority.
func (q *PriorityQueue[T]) Enqueue(ctx context.Context, item T) error {
	select {
	case q.regular <- item:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

// EnqueuePriority adds an item to the queue with high priority.
func (q *PriorityQueue[T]) EnqueuePriority(ctx context.Context, item T) error {
	select {
	case q.priority <- item:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

// Dequeue removes an item from the queue and processes it using the provided handler
// function. If configured, rate limits and concurrent worker limits will be enforced.
// This function will block if an item is not immediately available for
// processing or if necessary to enforce limits.
func (q *PriorityQueue[T]) Dequeue(ctx context.Context, handler func(context.Context, T) error) error {
	result, err := q.dequeueWithRateLimit(ctx)
	if err != nil {
		return err
	}

	if q.workersSem != nil {
		if err := q.incrementWorkers(ctx); err != nil {
			return err
		}
		defer q.decrementWorkers()
	}

	// Avoid panic if the handler is nil - could be used to pop an item without processing.
	if handler == nil {
		return nil
	}

	err = handler(ctx, result)

	// If there was any error we re-queue the item for processing again.
	if err != nil {
		if enqueueErr := q.Enqueue(ctx, result); enqueueErr != nil {
			return enqueueErr
		}

		return err
	}

	return nil
}

func (q *PriorityQueue[T]) dequeueWithRateLimit(ctx context.Context) (result T, err error) {
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

	// Set the err variable so that the deferred function can read it.
	result, err = q.dequeue(ctx)
	return result, err
}

func (q *PriorityQueue[T]) dequeue(ctx context.Context) (T, error) {
	var result T

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
			return result, ctx.Err()
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

// This operation may block until a worker slot is available.
func (q *PriorityQueue[T]) incrementWorkers(ctx context.Context) error {
	log.Debugf("incrementing workers")

	select {
	case q.workersSem <- struct{}{}:
		return nil

	case <-ctx.Done():
		return ctx.Err()

	default:
		log.Debug("enforcing workers limit")

		select {
		case q.workersSem <- struct{}{}:
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// This operation should never block - if it does that indicates that decrement
// has been called too many times.
func (q *PriorityQueue[T]) decrementWorkers() {
	log.Debugf("decrementing workers")

	select {
	case <-q.workersSem:
		return

	default:
		log.Panic("should never block")
	}
}
