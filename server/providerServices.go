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

func calculateMaskedRoot(real_root string, random_num string, address string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(real_root), []byte(random_num), []byte(address)))
}

func saveRoot(storage storage.Cyclic[RootsByAddress], roundId uint64, address string, root string, random string) {

	forRound, exists := storage.Get(roundId)

	if exists {
		forRound[address] = merkleRootStorageObject{root, random}
	} else {
		forRound = make(RootsByAddress)
		forRound[address] = merkleRootStorageObject{root, random}
		storage.Store(roundId, forRound)
	}

}

func (controller *FDCProtocolProviderController) submit1Service(roundId uint64, _ string) (string, bool, error) {
	votingRound, exists := controller.manager.Rounds.Get(roundId)
	if !exists {
		return "", false, nil
	}
	bitVoteString, err := votingRound.BitVoteHex()
	if err != nil {
		return "", false, err
	}
	return bitVoteString, true, nil
}

func (controller *FDCProtocolProviderController) submit2Service(roundId uint64, address string) (string, bool, error) {
	// Get merkle tree root from attestation client from controller

	votingRound, exists := controller.manager.Rounds.Get(roundId)

	if !exists {
		return "", false, nil
	}

	root, err := votingRound.GetMerkleRootCachedHex()

	if err != nil {
		return "", false, nil //decide what to do with empty round
	}

	r2 := rand.New(rand.NewSource(PROVIDER_RANDOM_SEED))
	random_num := hex.EncodeToString(crypto.Keccak256([]byte(fmt.Sprintf("%X", r2.Int63()))))

	// save root to storage
	//controller.saveRoot(address, roundID, root, random_num)
	saveRoot(controller.storage, roundId, address, root, random_num)

	masked := calculateMaskedRoot(root, random_num, address)

	return masked, true, nil
}

func (controller *FDCProtocolProviderController) submitSignaturesService(roundId uint64, address string) (addData, error) {
	// check storage if root was saved

	savedRoots, exists := controller.storage.Get(roundId)
	if !exists {
		return addData{}, fmt.Errorf("roots for round %d not stored", roundId)
	}

	rootData, exists := savedRoots[address]

	if !exists {
		return addData{}, fmt.Errorf("root for %s in round %d not in root storage", address, roundId)
	}

	log.Info("SubmitSignaturesHandler")
	log.Logf(zapcore.DebugLevel, "round: %s\n", fmt.Sprint(roundId))
	log.Logf(zapcore.DebugLevel, "address: %s\n", address)
	log.Logf(zapcore.DebugLevel, "root: %s\n", rootData.merkleRoot)
	log.Logf(zapcore.DebugLevel, "random: %s\n", rootData.randomNum)

	return addData{data: rootData.merkleRoot, additional: rootData.randomNum}, nil
}
