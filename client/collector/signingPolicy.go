package collector

import (
	"context"
	"errors"
	"flare-common/database"
	"flare-common/policy"
	"local/fdc/client/shared"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// SigningPolicyInitializedListener initiates a channel that serves signingPolicyInitialized events emitted by relayContractAddress.
func SigningPolicyInitializedListener(
	ctx context.Context,
	db *gorm.DB,
	relayContractAddress common.Address,
	registryContractAddress common.Address,
	votersDataChan chan<- []shared.VotersData,
) {
	logs, err := database.FetchLatestLogsByAddressAndTopic0(
		ctx, db, relayContractAddress, signingPolicyInitializedEventSel, 3,
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
	var sorted []shared.VotersData
	for i := range logs {
		votersData, err := AddSubmitAddressesToSigningPolicy(ctx, db, registryContractAddress, logs[len(logs)-i-1])
		if err != nil {
			log.Panic("error fetching initial signing policies with submit addresses:", err)
		}

		sorted = append(sorted, votersData)
		log.Info("fetched initial policy for round ", votersData.Policy.RewardEpochId)
	}

	select {
	case votersDataChan <- sorted:
	case <-ctx.Done():
		log.Info("SigningPolicyInitializedListener exiting:", ctx.Err())
	}

	spiTargetedListener(ctx, db, relayContractAddress, registryContractAddress, logs[0], latestQuery, votersDataChan)
}

// spiTargetedListener that only starts aggressive queries for new signingPolicyInitialized events a bit before the expected emission and stops once it gets one and waits until the next window.
func spiTargetedListener(
	ctx context.Context,
	db *gorm.DB,
	relayContractAddress common.Address,
	registryContractAddress common.Address,
	lastLog database.Log,
	latestQuery time.Time,
	votersDataChan chan<- []shared.VotersData,
) {
	lastSigningPolicy, err := policy.ParseSigningPolicyInitializedEvent(lastLog)
	if err != nil {
		log.Panic("error parsing initial logs:", err)
	}

	lastInitializedRewardEpochID := lastSigningPolicy.RewardEpochId.Uint64()

	for {
		expectedStartOfTheNextSigningPolicyInitialized := timing.ExpectedRewardEpochStartTimestamp(lastInitializedRewardEpochID + 1)
		untilStart := time.Until(time.Unix(int64(expectedStartOfTheNextSigningPolicyInitialized)-int64(timing.Chain.CollectDurationSec)*(int64(timing.Chain.RewardEpochLength)/20+1), 0)) // head start for querying 1/20 of the reward period
		timer := time.NewTimer(untilStart)

		log.Infof("next signing policy expected in %.1f minutes", untilStart.Minutes())
		select {
		case <-timer.C:
			log.Debug("querying for next signing policy")

		case <-ctx.Done():
			log.Info("spiTargetedListener exiting:", ctx.Err())
			return
		}

		logsWithSubmitAddresses, err := queryNextSPI(ctx, db, relayContractAddress, registryContractAddress, latestQuery, lastInitializedRewardEpochID)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				log.Info("spiTargetedListener exiting:", err)
				return
			}

			log.Error("error querying next SPI event:", err)
			continue
		}
		votersDataChan <- logsWithSubmitAddresses

		latestQuery = time.Now()
		lastInitializedRewardEpochID++
	}
}

func queryNextSPI(
	ctx context.Context,
	db *gorm.DB,
	relayContractAddress common.Address,
	registryContractAddress common.Address,
	latestQuery time.Time,
	latestRewardEpoch uint64,
) (
	[]shared.VotersData,
	error,
) {
	ticker := time.NewTicker(time.Duration(timing.Chain.CollectDurationSec-1) * time.Second) // ticker that is guaranteed to tick at least once per SystemVotingRound

	for {
		now := time.Now()

		logs, err := database.FetchLogsByAddressAndTopic0Timestamp(
			ctx, db, relayContractAddress, signingPolicyInitializedEventSel, latestQuery.Unix(), now.Unix(),
		)
		if err != nil {
			return nil, err
		}

		if len(logs) > 0 {
			log.Debug("Adding signing policy to channel")

			votersDataArray := make([]shared.VotersData, 0)
			for i := range logs {
				votersData, err := AddSubmitAddressesToSigningPolicy(ctx, db, registryContractAddress, logs[i])
				if err != nil {
					return nil, err
				}
				if votersData.Policy.RewardEpochId.Uint64() > latestRewardEpoch {
					votersDataArray = append(votersDataArray, votersData)
					log.Info("fetched policy for round ", votersData.Policy.RewardEpochId)
				}
			}
			if len(votersDataArray) > 0 {
				ticker.Stop()
				return votersDataArray, nil
			}
		}

		select {
		case <-ticker.C:
			log.Debug("starting next queryNextSPI iteration")

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
