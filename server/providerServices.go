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

// storeRoot stores root, random, and consensusBitVote for roundId to be used in submitSignatures.
//
// It only stores the merkleRootStorageObject the first time it is called.
func storeRoot(storage storage.Cyclic[merkleRootStorageObject], roundId uint64, message string, root, random common.Hash, consensusBitVote []byte) {
	_, exists := storage.Get(roundId)

	if exists {
		log.Debugf("root for round %d already stored", roundId)
	}
	if !exists {

		object := merkleRootStorageObject{
			message:          message,
			merkleRoot:       root,
			randomNum:        random,
			consensusBitVote: consensusBitVote,
		}
		storage.Store(roundId, object)
	}

}

// submit1Service returns 0x prefixed hex encoded bitVote for roundId and a boolean indicating its existence.
func (controller *FDCProtocolProviderController) submit1Service(roundId uint64, _ string) (string, bool, error) {
	votingRound, exists := controller.rounds.Get(roundId)
	if !exists {
		log.Infof("submit1 round %d not stored", roundId)
		return "", false, nil
	}
	bitVoteString, err := votingRound.BitVoteHex(false)
	if err != nil {
		log.Errorf("submit1: error for bitVote %s", err)

		return "", false, err
	}

	payloadMsg, err := payload.BuildMessage(controller.protocolId, roundId, bitVoteString)

	if err != nil {
		log.Errorf("submit1: error building payload %s", err)

		return "", false, err
	}

	log.Debugf("submit1: for round %d: %s", roundId, payloadMsg)

	return payloadMsg, true, nil
}

// submit2Service returns 0x prefixed commit data for roundId and address and an indicator of success.
// commit data is a hash of merkleRoot, roundId, address, and encodedBitVote.
// The data is stored to be used in the reveal.
func (controller *FDCProtocolProviderController) submit2Service(roundId uint64, address string) (string, bool, error) {

	commit, exists := controller.storage.Get(roundId)

	if exists {
		return commit.message, true, nil
	}

	votingRound, exists := controller.rounds.Get(roundId)

	if !exists {
		log.Infof("submit2: round %d not stored", roundId)

		return "", false, nil
	}

	consensusBitVote, err := votingRound.GetConsensusBitVote()

	if err != nil {
		log.Infof("submit2: consensus bitVote for round %d not available: %s", roundId, err)

		return "", false, nil
	}

	encodedBitVote := consensusBitVote.EncodeBitVote(roundId)
	root, err := votingRound.MerkleRoot()

	if err != nil {
		log.Infof("submit2: Merkle root for round %d not available: %s", roundId, err)

		return "", false, nil
	}

	randomBig, err := rand.Int(rand.Reader, maxRandom)

	if err != nil {
		log.Errorf("submit2: getting random number for round %d: %s", roundId, err)

		return "", false, nil
	}

	random := common.BigToHash(randomBig)

	masked := calculateMaskedRoot(root, random, common.HexToAddress(address), encodedBitVote)

	payloadMsg, err := payload.BuildMessage(controller.protocolId, roundId, masked)

	if err != nil {
		log.Errorf("submit2: error building payload for round %d: %s", roundId, err)

		return "", false, nil
	}

	storeRoot(controller.storage, roundId, payloadMsg, root, random, encodedBitVote)

	log.Debugf("submit2: for round %d and address %s: %s", roundId, address, masked)

	return payloadMsg, true, nil
}

// submitSignaturesService returns merkleRoot encoded in to payload for signing, additionalData, and an indicator of success for roundId.
// Additional data is concatenation of stored randomNumber and consensusBitVote.
func (controller *FDCProtocolProviderController) submitSignaturesService(roundId uint64, address string) (string, string, bool) {
	savedRoot, exists := controller.storage.Get(roundId)
	if !exists {
		log.Infof("submitSignatures: data for round %d not stored", roundId)
		return "", "", false
	}

	message := buildMessageForSigning(uint8(controller.protocolId), uint32(roundId), savedRoot.merkleRoot)

	log.Info("SubmitSignaturesHandler")
	log.Logf(zapcore.DebugLevel, "round: %s", fmt.Sprint(roundId))
	log.Logf(zapcore.DebugLevel, "address: %s", address)
	log.Logf(zapcore.DebugLevel, "root: %s", savedRoot.merkleRoot.Hex())
	log.Logf(zapcore.DebugLevel, "random: %s", savedRoot.randomNum.String())
	log.Logf(zapcore.DebugLevel, "consensus: %s", "0x"+hex.EncodeToString(savedRoot.consensusBitVote))

	additionalData := savedRoot.randomNum.Hex() + hex.EncodeToString(savedRoot.consensusBitVote)

	return message, additionalData, true
}

// buildMessageForSigning builds payload message for submitSignatures.
//
// protocolId (1 byte)
// roundId (4 bytes)
// randomQualityScore (1 byte)
// merkleRoot (32 bytes)
func buildMessageForSigning(protocolId uint8, roundId uint32, merkleRoot common.Hash) string {

	data := make([]byte, 38)

	data[0] = uint8(protocolId)                            // protocolId (1 byte)
	binary.BigEndian.PutUint32(data[1:5], uint32(roundId)) // roundId (4 bytes)
	data[5] = 1                                            // random quality score (1 byte)
	copy(data[6:38], merkleRoot[:])

	return "0x" + hex.EncodeToString(data)

}
