package server

import (
	"encoding/hex"
	"flare-common/logger"
	"flare-common/storage"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap/zapcore"
)

type addData struct {
	data       string
	additional string
}

// TODO: Luka - Get this from config
const PROVIDER_RANDOM_SEED = 42

var log = logger.GetLogger()

// calculateMaskedRoot masks the root with random number and address.
func calculateMaskedRoot(root string, random string, address string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(root), []byte(random), []byte(address)))
}

// storeRoot stores root and random for roundId and address for use in submitSignatures.
func storeRoot(storage storage.Cyclic[RootsByAddress], roundId uint64, address string, root string, random string) {

	forRound, exists := storage.Get(roundId)

	if exists {
		forRound[address] = merkleRootStorageObject{root, random} //store root for another address
	} else { // store for a new round
		forRound = make(RootsByAddress)
		forRound[address] = merkleRootStorageObject{root, random}
		storage.Store(roundId, forRound)
	}

}

func (controller *FDCProtocolProviderController) submit1Service(roundId uint64, _ string) (string, bool, error) {
	votingRound, exists := controller.manager.Rounds.Get(roundId)
	if !exists {
		log.Infof("submit1 round %d not stored", roundId)
		return "", false, nil
	}
	bitVoteString, err := votingRound.BitVoteHex()
	if err != nil {
		log.Errorf("submit1: error for bitVote %s", err)

		return "", false, err
	}

	log.Debugf("submit1: for round %d: %s", roundId, bitVoteString)
	return bitVoteString, true, nil
}

func (controller *FDCProtocolProviderController) submit2Service(roundId uint64, address string) (string, bool, error) {
	// Get merkle tree root from attestation client from controller

	votingRound, exists := controller.manager.Rounds.Get(roundId)

	if !exists {
		log.Infof("submit2: round %d not stored", roundId)

		return "", false, nil
	}

	root, err := votingRound.GetMerkleRootCachedHex()

	if err != nil {
		log.Infof("submit2: Merkle root for round %d not available", roundId)

		return "", false, nil //decide what to do with empty round
	}

	r2 := rand.New(rand.NewSource(PROVIDER_RANDOM_SEED))
	random := hex.EncodeToString(crypto.Keccak256([]byte(fmt.Sprintf("%X", r2.Int63()))))

	// storeRoot saves root together with random number and address to storage
	//controller.saveRoot(address, roundID, root, random_num)
	storeRoot(controller.storage, roundId, address, root, random)

	masked := calculateMaskedRoot(root, random, address)

	log.Debugf("submit2: for round %d and address %s: %s", roundId, address, masked)
	return masked, true, nil
}

func (controller *FDCProtocolProviderController) submitSignaturesService(roundId uint64, address string) (addData, bool) {

	savedRoots, exists := controller.storage.Get(roundId)
	if !exists {
		log.Infof("submitSigantures: data for round %d not stored", roundId)
		return addData{}, false
	}

	rootData, exists := savedRoots[address]

	if !exists {
		log.Infof("submitSigantures: root for address %s not stored for round %d", address, roundId)

		return addData{}, false
	}

	log.Info("SubmitSignaturesHandler")
	log.Logf(zapcore.DebugLevel, "round: %s\n", fmt.Sprint(roundId))
	log.Logf(zapcore.DebugLevel, "address: %s\n", address)
	log.Logf(zapcore.DebugLevel, "root: %s\n", rootData.merkleRoot)
	log.Logf(zapcore.DebugLevel, "random: %s\n", rootData.randomNum)

	return addData{data: rootData.merkleRoot, additional: rootData.randomNum}, true
}
