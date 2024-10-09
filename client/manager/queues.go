package manager

import (
	"context"

	"gitlab.com/flarenetwork/libs/go-flare-common/pkg/logger"
	"gitlab.com/flarenetwork/libs/go-flare-common/pkg/queue"

	"gitlab.com/flarenetwork/fdc/fdc-client/client/attestation"
	"gitlab.com/flarenetwork/fdc/fdc-client/client/config"
)

type attestationQueue = queue.PriorityQueue[*attestation.Attestation]

type priorityQueues map[string]*attestationQueue

// buildQueues builds queues from configurations
func buildQueues(queuesConfigs config.Queues) priorityQueues {

	queues := make(priorityQueues)

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

// runQueues runs all queues at once.
func runQueues(ctx context.Context, queues priorityQueues) {
	for k := range queues {
		go func(k string) {
			run(ctx, queues[k])
		}(k)
	}
}

// run tracks and handles all dequeued attestations from queue.
func run(ctx context.Context, queue *attestationQueue) {
	for {
		err := queue.Dequeue(ctx, handler)
		if err != nil {
			logger.Error(err)
		}

		if err := ctx.Err(); err != nil {
			logger.Infof("queue worker exiting: %v", err)
			return
		}
	}

}
