package collector

import (
	"context"
	"encoding/hex"
	"flare-common/database"
	"flare-common/policy"
	"local/fdc/client/shared"
	"math/big"

	"fmt"

	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

func fetchVoterRegisteredEventsForRewardEpoch(ctx context.Context, db *gorm.DB, registryContractAddress common.Address, rewardEpochID uint64) ([]database.Log, error) {
	var logs []database.Log

	epochIDBig := new(big.Int).SetUint64(rewardEpochID)

	epochID := common.BigToHash(epochIDBig)
	err := db.WithContext(ctx).Where(
		"address = ? AND topic0 = ? AND topic2 = ?",
		hex.EncodeToString(registryContractAddress[:]), // encodes without 0x prefix and without checksum
		hex.EncodeToString(voterRegisteredEventSel[:]),
		hex.EncodeToString(epochID[:]),
	).Find(&logs).Error

	return logs, err

}

// FetchVoterRegisteredEventsForRewardEpoch gets all VoterRegisteredEvents emitted for rewardEpochID.
func FetchVoterRegisteredEventsForRewardEpoch(ctx context.Context, db *gorm.DB, registryContractAddress common.Address, rewardEpochID uint64) ([]database.Log, error) {
	var logs []database.Log

	err := backoff.RetryNotify(
		func() error {
			var err error
			logs, err = fetchVoterRegisteredEventsForRewardEpoch(ctx, db, registryContractAddress, rewardEpochID)
			return err
		},
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
		func(err error, duration time.Duration) {
			log.Errorf("error fetching logs: %v, retrying after %v", err, duration)
		},
	)

	return logs, err
}

func BuildSubmitToSigningPolicyAddress(registryEvents []database.Log) (map[common.Address]common.Address, error) {
	submitToSigning := make(map[common.Address]common.Address)

	for i := range registryEvents {
		event, err := policy.ParseVoterRegisteredEvent(registryEvents[i])

		if err != nil {
			return nil, err
		}

		submitToSigning[event.SubmitAddress] = event.SigningPolicyAddress
	}

	return submitToSigning, nil

}

func SubmitToSigningPolicyAddress(ctx context.Context, db *gorm.DB, registryContractAddress common.Address, rewardEpochID uint64) (map[common.Address]common.Address, error) {
	logs, err := FetchVoterRegisteredEventsForRewardEpoch(ctx, db, registryContractAddress, rewardEpochID)
	if err != nil {
		return nil, fmt.Errorf("error fetching registered events: %s", err)
	}

	submitToSigning, err := BuildSubmitToSigningPolicyAddress(logs)
	if err != nil {
		return nil, fmt.Errorf("error building submitToSigning map: %s", err)
	}

	return submitToSigning, nil
}

func AddSubmitAddressesToSigningPolicy(ctx context.Context, db *gorm.DB, registryContractAddress common.Address, signingPolicyLog database.Log) (shared.VotersData, error) {
	data, err := policy.ParseSigningPolicyInitializedEvent(signingPolicyLog)
	if err != nil {
		return shared.VotersData{}, err
	}

	ok := data.RewardEpochId.IsUint64()

	if !ok {
		return shared.VotersData{}, fmt.Errorf("reward epoch %v too high", data.RewardEpochId)
	}

	rewardEpochID := data.RewardEpochId.Uint64()
	submitToSigning, err := SubmitToSigningPolicyAddress(ctx, db, registryContractAddress, rewardEpochID)

	log.Debugf("received %d registered submit addresses", len(submitToSigning))

	if err != nil {
		return shared.VotersData{}, fmt.Errorf("error adding submit addresses: %s", err)
	}

	return shared.VotersData{Policy: data, SubmitToSigningAddress: submitToSigning}, nil

}
