package collector

import (
	"flare-common/database"
	"local/fdc/client/attestation"
	"local/fdc/client/timing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// SigningPolicyInitializedListener returns a channel that serves signingPolicyInitialized events emitted by relayContractAddress.
func SigningPolicyInitializedListener(db *gorm.DB, relayContractAddress common.Address, bufferSize int) <-chan []database.Log {
	out := make(chan []database.Log, bufferSize)

	go func() {
		logs, err := database.FetchLatestLogsByAddressAndTopic0(
			db, relayContractAddress, signingPolicyInitializedEventSel, 3,
		)

		latestQuery := time.Now()

		if err != nil {
			log.Panic("error fetching initial logs:", err)
		}

		log.Debug("Logs length:", len(logs))

		if len(logs) == 0 {
			log.Panic("No initial signing policies found:", err)
		}

		sorted := []database.Log{}

		// signingPolicyStorage expects policies in increasing order
		for i := range logs {
			sorted = append(sorted, logs[len(logs)-i-1])
		}

		out <- sorted

		spiTargetedListener(db, relayContractAddress, logs[0], latestQuery, out)

	}()

	return out

}

// spiTargetedListener that only starts aggressive queries for new signingPolicyInitialized events a bit before the expected emission and stops once it get one and waits until the next window.
func spiTargetedListener(
	db *gorm.DB,
	relayContractAddress common.Address,
	lastLog database.Log,
	latestQuery time.Time,
	out chan<- []database.Log,
) {
	lastSigningPolicy, err := attestation.ParseSigningPolicyInitializedLog(lastLog)
	if err != nil {
		log.Panic("error parsing initial logs:", err)
	}

	lastInitializedRewardEpochID := lastSigningPolicy.RewardEpochId.Uint64()

	for {
		expectedStartOfTheNextSigningPolicyInitialized := timing.ExpectedRewardEpochStartTimestamp(lastInitializedRewardEpochID + 1)

		log.Info(expectedStartOfTheNextSigningPolicyInitialized)

		untilStart := time.Until(time.Unix(int64(expectedStartOfTheNextSigningPolicyInitialized)-90*15, 0)) //use const for headStart 90*15

		log.Infof("next signing policy expected in %.1fh", untilStart.Hours())

		timer := time.NewTimer(untilStart)

		<-timer.C

		if err := queryNextSPI(db, relayContractAddress, latestQuery, out); err != nil {
			log.Error("error querying next SPI event:", err)
			continue
		}

		lastInitializedRewardEpochID++
	}
}

func queryNextSPI(
	db *gorm.DB,
	relayContractAddress common.Address,
	latestQuery time.Time,
	out chan<- []database.Log,
) error {
	ticker := time.NewTicker(89 * time.Second) // ticker that is guaranteed to tick at least once per SystemVotingRound

	for {
		now := time.Now()

		logs, err := database.FetchLogsByAddressAndTopic0Timestamp(
			db, relayContractAddress, signingPolicyInitializedEventSel, latestQuery.Unix(), now.Unix(),
		)

		latestQuery = now

		if err != nil {
			return err
		}

		if len(logs) > 0 {
			log.Debug("Adding signing policy to channel")
			out <- logs

			ticker.Stop()
			return nil
		}

		<-ticker.C
	}
}
