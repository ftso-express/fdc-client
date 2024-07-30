package collector

import (
	"context"
	"errors"
	"flare-common/database"
	"flare-common/policy"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// SigningPolicyInitializedListener initiates a channel that serves signingPolicyInitialized events emitted by relayContractAddress.
func SigningPolicyInitializedListener(
	ctx context.Context,
	db collectorDB,
	relayContractAddress common.Address,
	logChan chan<- []database.Log,
) {
	logs, err := db.FetchLatestLogsByAddressAndTopic0(
		ctx, relayContractAddress, signingPolicyInitializedEventSel, 3,
	)
	if err != nil {
		log.Panic("error fetching initial logs:", err)
	}
	latestQuery := time.Now()
	log.Debug("Logs length:", len(logs))
	if len(logs) == 0 {
		log.Panic("No initial signing policies found:", err)
	}

	// signingPolicyStorage expects policies in increasing order
	var sorted []database.Log
	for i := range logs {
		sorted = append(sorted, logs[len(logs)-i-1])
	}

	select {
	case logChan <- sorted:
	case <-ctx.Done():
		log.Info("SigningPolicyInitializedListener exiting:", ctx.Err())
	}

	spiTargetedListener(ctx, db, relayContractAddress, logs[0], latestQuery, logChan)
}

// spiTargetedListener that only starts aggressive queries for new signingPolicyInitialized events a bit before the expected emission and stops once it gets one and waits until the next window.
func spiTargetedListener(
	ctx context.Context,
	db collectorDB,
	relayContractAddress common.Address,
	lastLog database.Log,
	latestQuery time.Time,
	out chan<- []database.Log,
) {
	lastSigningPolicy, err := policy.ParseSigningPolicyInitializedEvent(lastLog)
	if err != nil {
		log.Panic("error parsing initial logs:", err)
	}

	lastInitializedRewardEpochID := lastSigningPolicy.RewardEpochId.Uint64()

	for {
		expectedStartOfTheNextSigningPolicyInitialized := timing.ExpectedRewardEpochStartTimestamp(lastInitializedRewardEpochID + 1)

		untilStart := time.Until(time.Unix(int64(expectedStartOfTheNextSigningPolicyInitialized)-90*15, 0)) //use const for headStart 90*15

		log.Infof("next signing policy expected in %.1fh", untilStart.Hours())

		timer := time.NewTimer(untilStart)

		select {
		case <-timer.C:
			log.Debug("querying for next signing policy")

		case <-ctx.Done():
			log.Info("spiTargetedListener exiting:", ctx.Err())
			return
		}

		if err := queryNextSPI(ctx, db, relayContractAddress, latestQuery, out); err != nil {
			if errors.Is(err, ctx.Err()) {
				log.Info("spiTargetedListener exiting:", err)
				return
			}

			log.Error("error querying next SPI event:", err)
			continue
		}

		lastInitializedRewardEpochID++
	}
}

func queryNextSPI(
	ctx context.Context,
	db collectorDB,
	relayContractAddress common.Address,
	latestQuery time.Time,
	out chan<- []database.Log,
) error {
	ticker := time.NewTicker(89 * time.Second) // ticker that is guaranteed to tick at least once per SystemVotingRound

	for {
		now := time.Now()

		logs, err := db.FetchLogsByAddressAndTopic0Timestamp(
			ctx, relayContractAddress, signingPolicyInitializedEventSel, latestQuery.Unix(), now.Unix(),
		)

		latestQuery = now
		if err != nil {
			return err
		}

		if len(logs) > 0 {
			log.Debug("Adding signing policy to channel")

			select {
			case out <- logs:

			case <-ctx.Done():
				return ctx.Err()
			}

			ticker.Stop()
			return nil
		}

		select {
		case <-ticker.C:
			log.Debug("starting next queryNextSPI iteration")

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
