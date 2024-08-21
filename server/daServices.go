package server

import (
	"encoding/hex"
	"flare-common/merkle"
	"fmt"
	"local/fdc/client/attestation"
)

//  wip

func (controller *DAController) GetRequests(roundId uint64) ([]DARequest, bool) {
	round, exists := controller.Rounds.Get(roundId)

	if !exists {
		return nil, false
	}

	requests := make([]DARequest, len(round.Attestations))

	for i := range round.Attestations {
		requests[i] = AttestationToDARequest(round.Attestations[i])
	}

	return requests, true
}

func AttestationToDARequest(att *attestation.Attestation) DARequest {
	var status AttestationStatus

	// TODO: granulate options
	switch att.Status {
	case attestation.Success:
		status = Valid
	default:
		status = Failed

	}

	dARequest := DARequest{
		Request:   hex.EncodeToString(att.Request),
		Response:  hex.EncodeToString(att.Response),
		Status:    status,
		Consensus: att.Consensus,
		Indexes:   att.Indexes,
	}

	return dARequest
}

func (controller *DAController) GetAttestations(roundId uint64) ([]DAAttestation, bool) {
	round, exists := controller.Rounds.Get(roundId)

	if !exists {
		return nil, false
	}

	merkleTree, err := round.MerkleTree()

	if err != nil {
		return nil, false
	}

	attestations := make([]DAAttestation, 0)

	for i := range round.Attestations {
		att, ok, err := AttestationToDAAttestation(round.Attestations[i])

		if err != nil {
			return nil, false
		}

		if ok {
			err := att.addProof(merkleTree)

			if err != nil {
				return nil, false
			}

			attestations = append(attestations, att)

		}

	}
	return attestations, true

}

func AttestationToDAAttestation(att *attestation.Attestation) (DAAttestation, bool, error) {

	isConfirmed := att.Status == attestation.Success
	isSelected := att.Consensus

	if !isConfirmed && isSelected {
		return DAAttestation{}, false, fmt.Errorf("request %s in round %d is in consensus but not confirmed", hex.EncodeToString(att.Request), att.RoundId)
	}

	if !isConfirmed || !isSelected {
		return DAAttestation{}, false, nil
	}

	dAAttestation := DAAttestation{
		RoundId:  att.RoundId,
		Request:  hex.EncodeToString(att.Request),
		Response: hex.EncodeToString(att.Response),
		Abi:      *att.AbiString,
		hash:     att.Hash,
	}

	return dAAttestation, true, nil

}

func (DAAtt *DAAttestation) addProof(tree merkle.Tree) error {

	proofCommon, err := tree.GetProofFromHash(DAAtt.hash)

	if err != nil {
		return fmt.Errorf("no proof for request %s in round %d", DAAtt.Request, DAAtt.RoundId)
	}

	proof := make([]string, len(proofCommon))

	for i := range proofCommon {
		proof[i] = proofCommon[i].Hex()
	}

	DAAtt.Proof = proof

	return nil

}
