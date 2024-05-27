package server

import (
	"encoding/hex"
	"errors"
	"flare-common/logger"
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

const cyclicBufferSize uint64 = 10

var log = logger.GetLogger()

func calculateMaskedRoot(real_root string, random_num string, address string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(real_root), []byte(random_num), []byte(address)))
}

func (storage rootStorage) saveRoot(address string, roundId uint64, root string, random string) {

	roundMod := roundId % cyclicBufferSize

	roundStorage, ok := storage[roundMod]

	if !ok || roundStorage.roundId != roundId {

		rootsF := RootsForRound{roundId: roundId, roots: make(map[string]merkleRootStorageObject)}

		storage[roundMod] = rootsF
	}
	storage[roundMod].roots[address] = merkleRootStorageObject{merkleRoot: root, randomNum: random, roundId: roundId}
}

func (storage rootStorage) Root(roundId uint64, address string) (merkleRootStorageObject, error) {

	roundMod := roundId % cyclicBufferSize

	roundStorage, ok := storage[roundMod]
	if !ok || roundStorage.roundId != roundId {
		return merkleRootStorageObject{}, errors.New("no root for") // edit error
	}

	object, ok := roundStorage.roots[address]

	if !ok {
		return merkleRootStorageObject{}, errors.New("no root for") // edit error
	}

	return object, nil

}

func (controller *FDCProtocolProviderController) submit1Service(round uint64, address string) (string, error) {
	votingRound, exists := controller.manager.Round(round)
	if !exists {
		return "", errors.New("round does not exist")
	}
	bitVoteString, err := votingRound.BitVoteHex()
	if err != nil {
		return "", err
	}
	return bitVoteString, nil
}

func (controller *FDCProtocolProviderController) submit2Service(roundID uint64, address string) (string, error) {
	// Get merkle tree root from attestation client from controller

	// votingRound, exists := controller.manager.Round(roundID)

	// if !exists {
	// 	return "", errors.New("round does not exist")
	// }

	// root, err := votingRound.GetMerkleRootCachedHex()

	// if err != nil {
	// 	return "", errors.New("root does not exist")
	// }

	r1 := rand.New(rand.NewSource(int64(roundID)))
	real_root := hex.EncodeToString(crypto.Keccak256([]byte(fmt.Sprintf("%X", r1.Int63()))))

	r2 := rand.New(rand.NewSource(PROVIDER_RANDOM_SEED))
	random_num := hex.EncodeToString(crypto.Keccak256([]byte(fmt.Sprintf("%X", r2.Int63()))))

	// save root to storage
	//controller.saveRoot(address, roundID, root, random_num)
	controller.storage.saveRoot(address, roundID, real_root, random_num)

	//masked := calculateMaskedRoot(root, random_num)

	masked := calculateMaskedRoot(real_root, random_num, address)

	return masked, nil
}

func (controller *FDCProtocolProviderController) submitSignaturesService(round uint64, address string) (addData, error) {
	// check storage if root was saved

	savedRoot, err := controller.storage.Root(round, address)
	if err != nil {
		return addData{}, fmt.Errorf("round for address not in storage")
	}

	log.Info("SubmitSignaturesHandler")
	log.Logf(zapcore.DebugLevel, "round: %s\n", fmt.Sprint(round))
	log.Logf(zapcore.DebugLevel, "address: %s\n", address)
	log.Logf(zapcore.DebugLevel, "root: %s\n", savedRoot.merkleRoot)
	log.Logf(zapcore.DebugLevel, "random: %s\n", savedRoot.randomNum)

	return addData{data: savedRoot.merkleRoot, additional: savedRoot.randomNum}, nil
}
