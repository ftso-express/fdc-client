package queue

import (
	"context"
	"flare-common/logger"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const defaultMaxAttempts uint64 = 10

var log = logger.GetLogger()

// PriorityQueue is made up of two sub-queues - one regular and one with
// higher priority. Items can be enqueued in either queue and when dequeueing
// items from the priority queue are returned first.
type PriorityQueue[T any] struct {
	regular         chan priorityQueueItem[T]
	priority        chan priorityQueueItem[T]
	minDequeueDelta time.Duration
	lastDequeue     time.Time
	workersSem      chan struct{}
	maxAttempts     uint64
	deadLetterQueue chan T
	backoff         func() backoff.BackOff
	timeOff         time.Duration
}

type priorityQueueItem[T any] struct {
	value    T
	backoff  backoff.BackOff
	priority bool
}

// PriorityQueueInput values are used to construct a new PriorityQueue.
type PriorityQueueParams struct {
	Size                 int    `toml:"size"`
	MaxDequeuesPerSecond int    `toml:"max_dequeues_per_second"` // Set to 0 to disable rate-limiting.
	MaxWorkers           int    `toml:"max_workers"`             // Set to 0 for unlimited workers.
	MaxAttempts          int32  `toml:"max_attempts"`            // Set to negative for unlimited retry attempts. If unset or set to 0, the default value (10) is applied.
	TimeOff              uint32 `toml:"time_off"`                // In seconds. Only relevant if Backoff is not set.

	// Pass a callback to specify the backoff policy which affects when items
	// are returned to the queue after an error. By default, items are
	// re-queued after TimeOff.
	Backoff func() backoff.BackOff
}

// NewPriority constructs a new PriorityQueue.
func NewPriority[T any](input *PriorityQueueParams) PriorityQueue[T] {
	if input == nil {
		input = new(PriorityQueueParams)
	}

	q := PriorityQueue[T]{
		regular:  make(chan priorityQueueItem[T], input.Size),
		priority: make(chan priorityQueueItem[T], input.Size),
		backoff:  input.Backoff,
		timeOff:  time.Duration(input.TimeOff) * time.Second,
	}

	if input.MaxDequeuesPerSecond > 0 {
		q.minDequeueDelta = time.Second / time.Duration(input.MaxDequeuesPerSecond)
		log.Info("minDequeueDelta:", q.minDequeueDelta)
	}

	if input.MaxWorkers > 0 {
		q.workersSem = make(chan struct{}, input.MaxWorkers)
	}

	if input.MaxAttempts > -1 {

		q.maxAttempts = uint64(input.MaxAttempts)

		if input.MaxAttempts == 0 {
			q.maxAttempts = defaultMaxAttempts
		}

		q.deadLetterQueue = make(chan T, input.Size)
	}

	return q
}

func (q PriorityQueue[T]) DeadLetterQueue() <-chan T {
	return q.deadLetterQueue
}

// Enqueue adds an item to the queue with regular priority.
func (q *PriorityQueue[T]) Enqueue(ctx context.Context, item T) error {
	return q.enqueue(ctx, priorityQueueItem[T]{value: item, backoff: q.newBackoff(), priority: false})
}

func (q *PriorityQueue[T]) enqueue(ctx context.Context, item priorityQueueItem[T]) error {
	select {
	case q.regular <- item:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *PriorityQueue[T]) enqueuePriority(ctx context.Context, item priorityQueueItem[T]) error {
	select {
	case q.priority <- item:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

// EnqueuePriority adds an item to the queue with high priority.
func (q *PriorityQueue[T]) EnqueuePriority(ctx context.Context, item T) error {
	return q.enqueuePriority(ctx, priorityQueueItem[T]{value: item, backoff: q.newBackoff(), priority: true})

}

func (q *PriorityQueue[T]) newBackoff() (bOff backoff.BackOff) {
	if q.backoff == nil {
		bOff = backoff.NewConstantBackOff(q.timeOff)
	} else {
		bOff = q.backoff()
	}

	if q.maxAttempts > 0 {
		bOff = backoff.WithMaxRetries(bOff, q.maxAttempts-1)
	}

	return bOff
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

	err = handler(ctx, result.value)

	// If there was any error we re-queue the item for processing again.
	if err != nil {
		q.handleError(ctx, result)
		return err
	}

	return nil
}

func (q *PriorityQueue[T]) handleError(ctx context.Context, item priorityQueueItem[T]) {
	waitDuration := item.backoff.NextBackOff()

	if waitDuration == backoff.Stop {
		// Attempt to send the item to the dead letter queue, but do not block if it is full -
		// in that case the item will be discarded.
		select {
		case q.deadLetterQueue <- item.value:
			log.Errorf("max retry attempts reached, sent item to dead letter queue: %v", item.value)

		default:
			log.Errorf("max retry attempts reached, discarding as dead letter queue is full: %v", item.value)
		}
	}

	go func() {
		if waitDuration > 0 {
			log.Debugf("sleeping for %v before retrying", waitDuration)

			select {
			case <-time.After(waitDuration):

			case <-ctx.Done():
				log.Errorf("context cancelled while waiting to retry item %v", item.value)
				return
			}
		}

		if item.priority {
			err := q.enqueuePriority(ctx, item)
			if err != nil {
				log.Errorf("error enqueing priority item %v for retry: %v", item.value, err)

			}
		} else if err := q.enqueue(ctx, item); err != nil {
			log.Errorf("error enqueing item %v for retry: %v", item.value, err)
		}
	}()

}

func (q *PriorityQueue[T]) dequeueWithRateLimit(ctx context.Context) (result priorityQueueItem[T], err error) {
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

func (q *PriorityQueue[T]) dequeue(ctx context.Context) (priorityQueueItem[T], error) {
	var result priorityQueueItem[T]

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
