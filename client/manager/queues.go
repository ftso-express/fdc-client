package manager

import (
	"context"

	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/flare-foundation/go-flare-common/pkg/queue"

	"github.com/flare-foundation/fdc-client/client/attestation"
	"github.com/flare-foundation/fdc-client/client/config"
)

type attestationQueue = queue.PriorityQueue[*attestation.Attestation]

type attestationQueues map[string]*attestationQueue

// buildQueues builds attestation queues from configurations.
func buildQueues(queuesConfigs config.Queues) attestationQueues {
	queues := make(attestationQueues)

	for k := range queuesConfigs {
		params := queuesConfigs[k]
		queue := queue.NewPriority[*attestation.Attestation](&params)
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
func run(ctx context.Context, queue *attestationQueue) {
	for {
		err := queue.DequeueAsync(ctx, handler, discard)
		if err != nil {
			logger.Warn(err)
		}

		if err := ctx.Err(); err != nil {
			logger.Infof("queue worker exiting: %v", err)
			return
		}
	}
}
