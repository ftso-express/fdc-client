package server

// wip

import (
	"errors"
	"flare-common/restServer"
	"flare-common/storage"
	"local/fdc/client/round"
	"strconv"
)

type DAController struct {
	Rounds *storage.Cyclic[uint32, *round.Round]
}

type RequestsResponse struct {
	Status   ResponseStatus
	Requests []DARequest
}

type AttestationResponse struct {
	Status       ResponseStatus
	Attestations []DAAttestation
}

func validateRoundIDParam(params map[string]string) (uint32, error) {

	votingRoundIDStr, exists := params["votingRoundID"]

	if !exists {
		return 0, errors.New("missing votingRound param")
	}

	votingRoundID, err := strconv.ParseUint(votingRoundIDStr, 10, 32)

	if err != nil {
		return 0, errors.New("votingRound param is not a 32 bit number")
	}

	return uint32(votingRoundID), nil

}

func (controller *DAController) getRequestController(
	params map[string]string,
	_ interface{},
	_ interface{},
) (RequestsResponse, *restServer.ErrorHandler) {
	votingRoundID, err := validateRoundIDParam(params)

	if err != nil {
		log.Error(err)
		return RequestsResponse{}, restServer.BadParamsErrorHandler(err)
	}

	requests, exists := controller.GetRequests(votingRoundID)

	if !exists {
		return RequestsResponse{Status: NotAvailable}, nil
	}

	return RequestsResponse{Status: Ok, Requests: requests}, nil
}

func (controller *DAController) getAttestationController(
	params map[string]string,
	_ interface{},
	_ interface{}) (AttestationResponse, *restServer.ErrorHandler) {
	votingRoundID, err := validateRoundIDParam(params)

	if err != nil {
		log.Error(err)
		return AttestationResponse{}, restServer.BadParamsErrorHandler(err)
	}

	attestations, exists := controller.GetAttestations(votingRoundID)

	if !exists {
		return AttestationResponse{Status: NotAvailable}, nil
	}

	return AttestationResponse{Status: Ok, Attestations: attestations}, nil
}
