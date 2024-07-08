package attestation

import (
	"context"
	"flare-common/queue"
	"local/fdc/client/config"
	"sync"
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
func handler(ctx context.Context, at *Attestation) error {
	return at.handle(ctx)

}

// runQueues runs all queues at once.
func runQueues(ctx context.Context, queues priorityQueues) {
	var wg sync.WaitGroup

	wg.Add(len(queues))

	for k := range queues {
		go func(k string) {
			run(ctx, queues[k])
			wg.Done()
		}(k)
	}

	wg.Wait()
}

// run tracks and handles all dequeued attestations from queue.
func run(ctx context.Context, queue *attestationQueue) {
	for {
		err := queue.Dequeue(ctx, handler)
		if err != nil {
			log.Error(err)
		}

		if err := ctx.Err(); err != nil {
			log.Infof("queue worker exiting: %v", err)
			return
		}
	}

}
