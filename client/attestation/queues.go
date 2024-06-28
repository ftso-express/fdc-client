package attestation

import (
	"context"
	"flare-common/queue"
	"local/fdc/client/config"
)

type attestationQueue = queue.PriorityQueue[*Attestation]

type priorityQueues map[string]*attestationQueue

func buildQueues(queuesParams config.Queues) priorityQueues {

	queues := make(priorityQueues)

	for k := range queuesParams {

		params := queuesParams[k]
		queue := queue.NewPriority[*Attestation](&params)

		queues[k] = &queue
	}

	return queues

}

func handler(_ context.Context, at *Attestation) error {

	err := at.handle()

	return err

}

func runQueues(ctx context.Context, queues priorityQueues) {
	for k := range queues {
		go run(ctx, queues[k])
	}
}

func run(ctx context.Context, queue *attestationQueue) {

	for {

		err := queue.Dequeue(ctx, handler)

		if err != nil {
			log.Error(err)
		}
	}

}
