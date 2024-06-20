package queue

// PriorityQueue is made up of two sub-queues - one regular and one with
// higher priority. Items can be enqueued in either queue and when dequeueing
// items from the priority queue are returned first.
type PriorityQueue[T any] struct {
	regular  chan T
	priority chan T
}

func NewPriority[T any](size int) PriorityQueue[T] {
	return PriorityQueue[T]{
		regular:  make(chan T, size),
		priority: make(chan T, size),
	}
}

func (q *PriorityQueue[T]) Enqueue(item T) {
	q.regular <- item
}

func (q *PriorityQueue[T]) EnqueuePriority(item T) {
	q.priority <- item
}

func (q *PriorityQueue[T]) Dequeue() T {
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
