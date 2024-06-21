package queue_test

import (
	"flare-common/queue"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	size          = 10
	benchmarkSize = 100000
)

func TestEnqueueDequeue(t *testing.T) {
	q := queue.NewPriority[int](size, 0)

	for i := 0; i < size; i++ {
		q.Enqueue(i)
	}

	for i := 0; i < size; i++ {
		item := q.Dequeue()
		require.Equal(t, i, item)
	}
}

func TestEnqueuePriority(t *testing.T) {
	q := queue.NewPriority[int](size, 0)

	q.Enqueue(1)
	q.EnqueuePriority(42)

	require.Equal(t, 42, q.Dequeue())
	require.Equal(t, 1, q.Dequeue())
}

func TestBlockingDequeue(t *testing.T) {
	q := queue.NewPriority[int](size, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		item := q.Dequeue()
		assert.Equal(t, 42, item)
	}()

	q.Enqueue(42)
	wg.Wait()
}

func TestBlockingDequeuePriority(t *testing.T) {
	q := queue.NewPriority[int](size, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		item := q.Dequeue()
		assert.Equal(t, 42, item)
	}()

	q.EnqueuePriority(42)
	wg.Wait()
}

func TestRateLimit(t *testing.T) {
	minDelta := 10 * time.Millisecond
	maxRate := int(time.Second / minDelta)
	t.Log("maxRate:", maxRate)

	q := queue.NewPriority[int](size, maxRate)

	for i := 0; i < size; i++ {
		q.Enqueue(i)
	}

	var lastDequeueTime *time.Time
	for i := 0; i < size; i++ {
		q.Dequeue()
		now := time.Now()
		if lastDequeueTime != nil {
			delta := now.Sub(*lastDequeueTime)
			require.GreaterOrEqual(t, delta, minDelta, fmt.Sprintf("failed iteration %d", i))
		}

		lastDequeueTime = &now
	}
}

func BenchmarkPriorityQueue(b *testing.B) {
	q := queue.NewPriority[int](size, 0)

	for n := 0; n < b.N; n++ {
		q.Enqueue(1)
		q.EnqueuePriority(2)
		q.Dequeue()
		q.Dequeue()
	}
}
