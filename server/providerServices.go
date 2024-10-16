package server

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/flare-foundation/go-flare-common/pkg/logger"
	"github.com/flare-foundation/go-flare-common/pkg/payload"

	"github.com/ethereum/go-ethereum/common"
)

// submit1Service returns an empty response.
func (controller *FDCProtocolProviderController) submit1Service(_ uint32, _ string) (string, bool, error) {
	return "0x", true, nil
}

// submit2Service returns 0x prefixed hex encoded bitVote for roundID and a boolean indicating its existence.
func (controller *FDCProtocolProviderController) submit2Service(roundID uint32, _ string) (string, bool, error) {
	votingRound, exists := controller.rounds.Get(roundID)
	if !exists {
		logger.Infof("submit2: round %d not stored", roundID)
		return "", false, nil
	}
	bitVote, err := votingRound.BitVoteBytes()
	if err != nil {
		logger.Errorf("submit2: error for bitVote %s", err)

		return "", false, err
	}

	payloadMsg := payload.BuildMessage(controller.protocolID, roundID, bitVote)

	logger.Debugf("submit2: for round %d: %s", roundID, payloadMsg)

	return payloadMsg, true, nil
}

// submitSignaturesService returns merkleRoot encoded in to payload for signing, additionalData, and an indicator of success for roundID.
// Additional data is concatenation of stored randomNumber and consensusBitVote.
func (controller *FDCProtocolProviderController) submitSignaturesService(roundID uint32, address string) (PDPResponse, bool, error) {
	votingRound, exists := controller.rounds.Get(roundID)
	if !exists {
		logger.Infof("submitSignatures: round %d not stored", roundID)
		return PDPResponse{Status: NotAvailable}, false, nil
	}

	consensusBitVote, exists, computed := votingRound.GetConsensusBitVote()
	if !computed {
		logger.Infof("submitSignatures: consensus bitVote for round %d not computed", roundID)
		return PDPResponse{Status: NotAvailable}, false, nil
	} else if !exists {
		logger.Infof("submitSignatures: consensus bitVote for round %d not available: %s", roundID)
		return PDPResponse{Status: Ok}, false, nil
	}

	encodedBitVote := "0x" + consensusBitVote.EncodeBitVoteHex()

	root, err := votingRound.MerkleRoot()
	if err != nil {
		logger.Infof("submitSignatures: Merkle root for round %d not available: %s", roundID, err)

		return PDPResponse{Status: NotAvailable}, false, nil
	}

	message := buildMessageForSigning(uint8(controller.protocolID), uint32(roundID), root)

	logger.Info("SubmitSignaturesHandler")
	logger.Debugf("round: %s", fmt.Sprint(roundID))
	logger.Debugf("address: %s", address)
	logger.Debugf("root: %v", root)
	logger.Debugf("consensus: %s", encodedBitVote)

	return PDPResponse{Status: Ok, Data: message, AdditionalData: encodedBitVote}, true, nil
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
	data[5] = 0 // is secure random not used in FDC protocol
	copy(data[6:38], merkleRoot[:])

	return "0x" + hex.EncodeToString(data)

}
