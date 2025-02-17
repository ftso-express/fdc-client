package manager

import (
	"context"

	"github.com/flare-foundation/fdc-client/client/attestation"
	"github.com/flare-foundation/fdc-client/client/config"
	"github.com/flare-foundation/go-flare-common/pkg/priority"
)

// weight implements priority.Weight[wTup]
type weight struct {
	Round uint32
}

func (x weight) Self() weight {
	return x
}

// Less returns true if x represents lower priority than y
//
//   - ">" later rounds have lower priority
//   - "=" implementation detail (if two items have the same priority, we do not want the later to have priority)
func (x weight) Less(y weight) bool {
	return x.Round >= y.Round
}

type attestationQueue = priority.PriorityQueue[*attestation.Attestation, weight]

type attestationQueues map[string]*attestationQueue

// buildQueues builds attestation queues from configurations.
func buildQueues(queuesConfigs config.Queues) attestationQueues {
	queues := make(attestationQueues)

	for k := range queuesConfigs {
		params := queuesConfigs[k]
		queue := priority.New[*attestation.Attestation, weight](params, k)
		queues[k] = &queue
	}

	return queues
}

// handler handles dequeued attestation.
func handler(ctx context.Context, at *attestation.Attestation) error {
	return at.Handle(ctx)
}

// discard discards requests that do not need to be handled
func discard(ctx context.Context, at *attestation.Attestation) bool {
	return at.Discard(ctx)
}

// runQueues runs all attestation queues at once.
func runQueues(ctx context.Context, queues attestationQueues) {
	for k := range queues {
		go func(k string) {
			run(ctx, queues[k])
		}(k)
	}
}

// run tracks and handles all dequeued attestations from a queue.
func run(ctx context.Context, q *attestationQueue) {
	q.InitiateAndRun(ctx)
	for {
		q.Dequeue(ctx, handler, discard)

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
