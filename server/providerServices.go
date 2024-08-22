package server

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"flare-common/logger"
	"flare-common/payload"
	"flare-common/storage"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap/zapcore"
)

var maxRandom *big.Int // max uint256

var log = logger.GetLogger()

func init() {
	// set maxRandom to max uint256
	maxUint256Plus1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	maxRandom = new(big.Int).Sub(maxUint256Plus1, big.NewInt(1))
}

// calculateMaskedRoot masks the root with random number and address.
func calculateMaskedRoot(root common.Hash, random common.Hash, address common.Address, bitVote []byte) string {
	return hex.EncodeToString(crypto.Keccak256(root.Bytes(), random.Bytes(), address.Bytes(), bitVote))
}

// storeRoot stores root, random, and consensusBitVote for roundID to be used in submitSignatures.
//
// It only stores the merkleRootStorageObject the first time it is called.
func storeRoot(storage storage.Cyclic[merkleRootStorageObject], roundID uint64, message string, root, random common.Hash, consensusBitVote []byte) {
	_, exists := storage.Get(roundID)

	if exists {
		log.Debugf("root for round %d already stored", roundID)
	}
	if !exists {

		object := merkleRootStorageObject{
			message:          message,
			merkleRoot:       root,
			randomNum:        random,
			consensusBitVote: consensusBitVote,
		}
		storage.Store(roundID, object)
	}

}

// submit1Service returns 0x prefixed hex encoded bitVote for roundID and a boolean indicating its existence.
func (controller *FDCProtocolProviderController) submit1Service(roundID uint64, _ string) (string, bool, error) {
	votingRound, exists := controller.rounds.Get(roundID)
	if !exists {
		log.Infof("submit1 round %d not stored", roundID)
		return "", false, nil
	}
	bitVoteString, err := votingRound.BitVoteHex(false)
	if err != nil {
		log.Errorf("submit1: error for bitVote %s", err)

		return "", false, err
	}

	payloadMsg, err := payload.BuildMessage(controller.protocolID, roundID, bitVoteString)

	if err != nil {
		log.Errorf("submit1: error building payload %s", err)

		return "", false, err
	}

	log.Debugf("submit1: for round %d: %s", roundID, payloadMsg)

	return payloadMsg, true, nil
}

// submit2Service returns 0x prefixed commit data for roundID and address and an indicator of success.
// commit data is a hash of merkleRoot, roundID, address, and encodedBitVote.
// The data is stored to be used in the reveal.
func (controller *FDCProtocolProviderController) submit2Service(roundID uint64, address string) (string, bool, error) {

	commit, exists := controller.storage.Get(roundID)

	if exists {
		return commit.message, true, nil
	}

	votingRound, exists := controller.rounds.Get(roundID)

	if !exists {
		log.Infof("submit2: round %d not stored", roundID)

		return "", false, nil
	}

	consensusBitVote, err := votingRound.GetConsensusBitVote()

	if err != nil {
		log.Infof("submit2: consensus bitVote for round %d not available: %s", roundID, err)

		return "", false, nil
	}

	encodedBitVote := consensusBitVote.EncodeBitVote(roundID)
	root, err := votingRound.MerkleRoot()

	if err != nil {
		log.Infof("submit2: Merkle root for round %d not available: %s", roundID, err)

		return "", false, nil
	}

	randomBig, err := rand.Int(rand.Reader, maxRandom)

	if err != nil {
		log.Errorf("submit2: getting random number for round %d: %s", roundID, err)

		return "", false, nil
	}

	random := common.BigToHash(randomBig)

	masked := calculateMaskedRoot(root, random, common.HexToAddress(address), encodedBitVote)

	payloadMsg, err := payload.BuildMessage(controller.protocolID, roundID, masked)

	if err != nil {
		log.Errorf("submit2: error building payload for round %d: %s", roundID, err)

		return "", false, nil
	}

	storeRoot(controller.storage, roundID, payloadMsg, root, random, encodedBitVote)

	log.Debugf("submit2: for round %d and address %s: %s", roundID, address, masked)

	return payloadMsg, true, nil
}

// submitSignaturesService returns merkleRoot encoded in to payload for signing, additionalData, and an indicator of success for roundID.
// Additional data is concatenation of stored randomNumber and consensusBitVote.
func (controller *FDCProtocolProviderController) submitSignaturesService(roundID uint64, address string) (string, string, bool) {
	savedRoot, exists := controller.storage.Get(roundID)
	if !exists {
		log.Infof("submitSignatures: data for round %d not stored", roundID)
		return "", "", false
	}

	message := buildMessageForSigning(uint8(controller.protocolID), uint32(roundID), savedRoot.merkleRoot)

	log.Info("SubmitSignaturesHandler")
	log.Logf(zapcore.DebugLevel, "round: %s", fmt.Sprint(roundID))
	log.Logf(zapcore.DebugLevel, "address: %s", address)
	log.Logf(zapcore.DebugLevel, "root: %s", savedRoot.merkleRoot.Hex())
	log.Logf(zapcore.DebugLevel, "random: %s", savedRoot.randomNum.String())
	log.Logf(zapcore.DebugLevel, "consensus: %s", "0x"+hex.EncodeToString(savedRoot.consensusBitVote))

	additionalData := savedRoot.randomNum.Hex() + hex.EncodeToString(savedRoot.consensusBitVote)

	return message, additionalData, true
}

// buildMessageForSigning builds payload message for submitSignatures.
//
// protocolID (1 byte)
// roundID (4 bytes)
// randomQualityScore (1 byte)
// merkleRoot (32 bytes)
func buildMessageForSigning(protocolID uint8, roundID uint32, merkleRoot common.Hash) string {

	data := make([]byte, 38)

	data[0] = uint8(protocolID)
	binary.BigEndian.PutUint32(data[1:5], uint32(roundID))
	data[5] = 1
	copy(data[6:38], merkleRoot[:])

	return "0x" + hex.EncodeToString(data)

}
