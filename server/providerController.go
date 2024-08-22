package server

import (
	"flare-common/restServer"
	"flare-common/storage"
	"fmt"
	"local/fdc/client/round"
	"strconv"
)

type FDCProtocolProviderController struct {
	rounds     *storage.Cyclic[*round.Round]
	storage    storage.Cyclic[merkleRootStorageObject]
	protocolID uint64
}

type submitXParams struct {
	votingRoundID uint64
	submitAddress string
}

const storageSize = 10

func newFDCProtocolProviderController(rounds *storage.Cyclic[*round.Round], protocolID uint64) *FDCProtocolProviderController {
	storage := storage.NewCyclic[merkleRootStorageObject](storageSize)

	return &FDCProtocolProviderController{rounds: rounds, storage: storage, protocolID: protocolID}
}

func validateSubmitXParams(params map[string]string) (submitXParams, error) {
	if _, ok := params["votingRoundID"]; !ok {
		return submitXParams{}, fmt.Errorf("missing votingRound param")
	}
	votingRoundID, err := strconv.ParseUint(params["votingRoundID"], 10, 64)
	if err != nil {
		return submitXParams{}, fmt.Errorf("votingRound param is not a number")
	}
	if _, ok := params["submitAddress"]; !ok {
		return submitXParams{}, fmt.Errorf("missing submitAddress param")
	}
	submitAddress := params["submitAddress"]
	if !restServer.ValidateEVMAddressString(submitAddress) {
		return submitXParams{}, fmt.Errorf("submitAddress param is not a valid EVM address")
	}
	return submitXParams{votingRoundID: votingRoundID, submitAddress: submitAddress}, nil
}

func (controller *FDCProtocolProviderController) submit1Controller(
	params map[string]string,
	queryParams interface{},
	body interface{}) (PDPResponse, *restServer.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)

	if err != nil {
		log.Error(err)
		return PDPResponse{}, restServer.BadParamsErrorHandler(err)
	}
	rsp, exists, err := controller.submit1Service(pathParams.votingRoundID, pathParams.submitAddress)
	if err != nil {
		log.Error(err)
		return PDPResponse{}, restServer.InternalServerErrorHandler(err)
	}
	if !exists {
		return PDPResponse{Data: rsp, Status: NotAvailable}, nil
	}

	response := PDPResponse{Data: rsp, Status: Ok}

	return response, nil
}

func (controller *FDCProtocolProviderController) submit2Controller(
	params map[string]string,
	queryParams interface{},
	body interface{}) (PDPResponse, *restServer.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)
	if err != nil {
		log.Error(err)
		return PDPResponse{}, restServer.BadParamsErrorHandler(err)
	}
	rsp, exists, err := controller.submit2Service(pathParams.votingRoundID, pathParams.submitAddress)
	if err != nil {
		log.Error(err)
		return PDPResponse{}, restServer.InternalServerErrorHandler(err)
	}

	if !exists {
		return PDPResponse{Data: rsp, Status: NotAvailable}, nil
	}

	response := PDPResponse{Data: rsp, Status: Ok}

	return response, nil
}

func (controller *FDCProtocolProviderController) submitSignaturesController(
	params map[string]string,
	queryParams interface{},
	body interface{}) (PDPResponse, *restServer.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)
	if err != nil {
		log.Error(err)
		return PDPResponse{}, restServer.BadParamsErrorHandler(err)
	}
	message, addData, exists := controller.submitSignaturesService(pathParams.votingRoundID, pathParams.submitAddress)
	if !exists {
		return PDPResponse{Status: NotAvailable}, nil
	}
	response := PDPResponse{Data: message, AdditionalData: addData, Status: Ok}

	return response, nil
}
