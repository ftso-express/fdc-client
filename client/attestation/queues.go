package attestation

import (
	"context"
	"flare-common/queue"
	"local/fdc/client/config"
)

type attestationQueue = queue.PriorityQueue[*Attestation]

type priorityQueues map[string]*attestationQueue

// buildQueues builds queues from configurations
func buildQueues(queuesConfigs config.Queues) priorityQueues {

	queues := make(priorityQueues)

	for k := range queuesConfigs {

		params := queuesConfigs[k]
		queue := queue.NewPriority[*Attestation](&params)

		queues[k] = &queue
	}

	return queues

}

// handler handles dequeued attestation.
func handler(_ context.Context, at *Attestation) error {

	err := at.handle()

	return err

}

// runQueues runs all queues at once.
func runQueues(ctx context.Context, queues priorityQueues) {
	for k := range queues {
		go run(ctx, queues[k])
	}
}

// run tracks and handles all dequeued attestations from queue.
func run(ctx context.Context, queue *attestationQueue) {

	for {

		err := queue.Dequeue(ctx, handler)

		if err != nil {
			log.Error(err)
		}
	}

}
