package server

import (
	"errors"
	"flare-common/restServer"
	"flare-common/storage"
	"local/fdc/client/round"
	"strconv"
)

type FDCDAController struct {
	rounds *storage.Cyclic[*round.Round]
}

type RequestsResponse struct {
	Status   ResponseStatus
	Requests []DARequest
}

type AttestationResponse struct {
	Status       ResponseStatus
	Attestations []DAAttestation
}

func validateRoundIdParam(params map[string]string) (uint64, error) {

	votingRoundIdStr, exists := params["votingRoundId"]

	if !exists {
		return 0, errors.New("missing votingRound param")
	}

	votingRoundId, err := strconv.ParseUint(votingRoundIdStr, 10, 64)

	if err != nil {
		return 0, errors.New("votingRound param is not a number")
	}

	return votingRoundId, nil

}

func (controller *FDCDAController) getRequestController(
	params map[string]string,
	_ interface{},
	_ interface{}) (RequestsResponse, *restServer.ErrorHandler) {
	votingRound, err := validateRoundIdParam(params)

	if err != nil {
		log.Error(err)
		return RequestsResponse{}, restServer.BadParamsErrorHandler(err)
	}

	request, exists := controller.GetRequests(votingRound)

	if !exists {
		return RequestsResponse{Status: NOT_AVAILABLE}, nil
	}

	return RequestsResponse{Status: OK, Requests: request}, nil
}

func (controller *FDCDAController) getAttestationController(
	params map[string]string,
	_ interface{},
	_ interface{}) (AttestationResponse, *restServer.ErrorHandler) {
	votingRound, err := validateRoundIdParam(params)

	if err != nil {
		log.Error(err)
		return AttestationResponse{}, restServer.BadParamsErrorHandler(err)
	}

	attestations, exists := controller.GetAttestations(votingRound)

	if !exists {
		return AttestationResponse{Status: NOT_AVAILABLE}, nil
	}

	return AttestationResponse{Status: OK, Attestations: attestations}, nil
}
