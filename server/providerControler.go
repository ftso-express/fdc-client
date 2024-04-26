package server

import (
	"context"
	"flare-common/restServer"
	"fmt"
	"net/http"
	"strconv"
)

type FDCProtocolProviderController struct {
	someValue string
}

type submitXParams struct {
	votingRoundId uint64
	submitAddress string
}

func validateSubmitXParams(params map[string]string) (submitXParams, error) {
	if _, ok := params["votingRoundId"]; !ok {
		return submitXParams{}, fmt.Errorf("missing votingRound param")
	}
	votingRoundId, err := strconv.ParseUint(params["votingRoundId"], 10, 64)
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
	return submitXParams{votingRoundId: votingRoundId, submitAddress: submitAddress}, nil
}

func NewFDCProtocolProviderController(ctx context.Context) *FDCProtocolProviderController {
	return &FDCProtocolProviderController{someValue: "initial value"}
}

func (controller *FDCProtocolProviderController) Submit1(
	params map[string]string,
	queryParams interface{},
	body interface{}) (PDPResponse, *restServer.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)
	if err != nil {
		return PDPResponse{}, restServer.BadParamsErrorHandler(err)
	}
	rsp := submit1Handler(pathParams.votingRoundId, pathParams.submitAddress)
	fmt.Printf("previous value: %s\n", controller.someValue)
	controller.someValue = "Submit1"
	return rsp, nil
}

func (controller *FDCProtocolProviderController) Submit2(
	params map[string]string,
	queryParams interface{},
	body interface{}) (PDPResponse, *restServer.ErrorHandler) {
	pathParams, err := validateSubmitXParams(params)
	if err != nil {
		return PDPResponse{}, restServer.BadParamsErrorHandler(err)
	}
	// TODO: update this method call
	rsp := submit1Handler(pathParams.votingRoundId, pathParams.submitAddress)
	fmt.Printf("previous value: %s\n", controller.someValue)
	controller.someValue = "Submit2"
	return rsp, nil
}

func RegisterFDCProviderRoutes(router restServer.Router, ctx context.Context) {
	ctrl := NewFDCProtocolProviderController(ctx)
	paramMap := map[string]string{"votingRoundId": "Voting round ID", "submitAddress": "Submit address"}
	submit1Handler := restServer.GeneralRouteHandler(ctrl.Submit1, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submit1/{votingRoundId}/{submitAddress}", submit1Handler, "Submit1")
	submit2Handler := restServer.GeneralRouteHandler(ctrl.Submit2, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submit2/{votingRoundId}/{submitAddress}", submit2Handler, "Submit2")
}
