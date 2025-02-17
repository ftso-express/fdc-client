package manager

import (
	"context"

	"github.com/flare-foundation/fdc-client/client/attestation"
	"github.com/flare-foundation/fdc-client/client/config"
	"github.com/flare-foundation/go-flare-common/pkg/priority"
)

// Weight implements priority.Weight[wTup]
type Weight struct {
	Index attestation.IndexLog
}

func (x Weight) Self() Weight {
	return x
}

// Less returns true if x represents lower priority than y
func (x Weight) Less(y Weight) bool {
	return attestation.EarlierLog(y.Index, x.Index)
}

type attestationQueue = priority.PriorityQueue[*attestation.Attestation, Weight]

type attestationQueues map[string]*attestationQueue

// buildQueues builds attestation queues from configurations.
func buildQueues(queuesConfigs config.Queues) attestationQueues {
	queues := make(attestationQueues)

	for k := range queuesConfigs {
		params := queuesConfigs[k]
		queue := priority.New[*attestation.Attestation, Weight](params, k)
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
