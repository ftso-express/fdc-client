package queue_test

import (
	"context"
	"errors"
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

var defaultTimeout = 10 * time.Second

func TestEnqueueDequeue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	q := queue.NewPriority[int](size, 0)

	for i := 0; i < size; i++ {
		err := q.Enqueue(ctx, i)
		require.NoError(t, err)
	}

	for i := 0; i < size; i++ {
		item, err := q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, i, item)
	}
}

func TestEnqueuePriority(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	q := queue.NewPriority[int](size, 0)

	err := q.Enqueue(ctx, 1)
	require.NoError(t, err)

	err = q.EnqueuePriority(ctx, 42)
	require.NoError(t, err)

	item, err := q.Dequeue(ctx)
	require.NoError(t, err)
	require.Equal(t, 42, item)

	item, err = q.Dequeue(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, item)
}

func TestBlockingDequeue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	q := queue.NewPriority[int](size, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		item, err := q.Dequeue(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 42, item)
	}()

	err := q.Enqueue(ctx, 42)
	require.NoError(t, err)
	wg.Wait()
}

func TestBlockingDequeuePriority(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	q := queue.NewPriority[int](size, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		item, err := q.Dequeue(ctx)
		require.NoError(t, err)
		assert.Equal(t, 42, item)
	}()

	err := q.EnqueuePriority(ctx, 42)
	require.NoError(t, err)
	wg.Wait()
}

func TestRateLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	minDelta := 10 * time.Millisecond
	maxRate := int(time.Second / minDelta)
	t.Log("maxRate:", maxRate)

	q := queue.NewPriority[int](size, maxRate)

	for i := 0; i < size; i++ {
		err := q.Enqueue(ctx, i)
		require.NoError(t, err)
	}

	var lastDequeueTime *time.Time
	for i := 0; i < size; i++ {
		item, err := q.Dequeue(ctx)
		require.NoError(t, err)
		require.Equal(t, i, item)

		now := time.Now()
		if lastDequeueTime != nil {
			delta := now.Sub(*lastDequeueTime)
			require.GreaterOrEqual(t, delta, minDelta, fmt.Sprintf("failed iteration %d", i))
		}

		lastDequeueTime = &now
	}
}

func TestEnqueueTimeout(t *testing.T) {
	q := queue.NewPriority[int](0, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	err := q.Enqueue(ctx, 1)
	require.Error(t, err)
	require.True(t, errors.Is(err, ctx.Err()))
}

func TestEnqueuePriorityTimeout(t *testing.T) {
	q := queue.NewPriority[int](0, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	err := q.EnqueuePriority(ctx, 1)
	require.Error(t, err)
	require.True(t, errors.Is(err, ctx.Err()))
}

func TestDequeueTimeout(t *testing.T) {
	q := queue.NewPriority[int](0, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	_, err := q.Dequeue(ctx)
	require.Error(t, err)
	require.True(t, errors.Is(err, ctx.Err()))
}

func TestDequeueRateLimitTimeout(t *testing.T) {
	ctx := context.Background()

	q := queue.NewPriority[int](2, 1)

	for i := 0; i < 2; i++ {
		err := q.Enqueue(ctx, i)
		require.NoError(t, err)
	}

	item, err := q.Dequeue(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, item)

	ctx, cancel := context.WithTimeout(ctx, time.Nanosecond)
	defer cancel()

	// Since we are immediately attempting to dequeue the second item, this
	// should start to block for around a second to enforce the 1-per-second
	// rate limit. However, the context will cancel long before the rate limit
	// elapses.
	start := time.Now()

	_, err = q.Dequeue(ctx)
	require.Error(t, err)
	require.True(t, errors.Is(err, ctx.Err()))

	// Check that we wait much less than a full second before exiting.
	delta := time.Since(start)
	require.Less(t, delta, 100*time.Millisecond)
}

func BenchmarkPriorityQueue(b *testing.B) {
	ctx := context.Background()

	q := queue.NewPriority[int](size, 0)

	for n := 0; n < b.N; n++ {
		_ = q.Enqueue(ctx, 1)
		_ = q.EnqueuePriority(ctx, 2)
		_, _ = q.Dequeue(ctx)
		_, _ = q.Dequeue(ctx)
	}
}
