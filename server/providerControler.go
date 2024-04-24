package server

import (
	"context"
	"flare-common/restServer"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type FDCProtocolProviderController struct {
	someValue string
}

func NewFDCProtocolProviderController(ctx context.Context) *FDCProtocolProviderController {
	return &FDCProtocolProviderController{someValue: "initial value"}
}

func (controller *FDCProtocolProviderController) Submit1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Submit1Handler")
	fmt.Printf("%+v\n", r)
	fmt.Printf("%+v\n", vars)
	// TODO: validate input vars
	rsp := submit1Handler(vars["votingRound"], vars["submitAddress"])
	restServer.WriteApiResponseOk(w, rsp)
	fmt.Printf("previous value: %s\n", controller.someValue)
	controller.someValue = "Submit1"
}

func (controller *FDCProtocolProviderController) Submit2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Submit2Handler")
	fmt.Printf("%+v\n", r)
	fmt.Printf("%+v\n", vars)
	// TODO: validate input vars
	rsp := submit1Handler(vars["votingRound"], vars["submitAddress"])
	restServer.WriteApiResponseOk(w, rsp)
	fmt.Printf("previous value: %s\n", controller.someValue)
	controller.someValue = "Submit2"
}

func RegisterFDCProviderRoutes(router restServer.Router, ctx context.Context) {
	ctrl := NewFDCProtocolProviderController(ctx)
	paramMap := map[string]string{"votingRoundId": "Voting round ID", "submitAddress": "Submit address"}
	submit1Handler := restServer.WrappedRouteHandler(ctrl.Submit1, http.MethodGet, http.StatusOK, paramMap, PDPResponse{})
	router.AddRoute("/submit1/{votingRoundId}/{submitAddress}", submit1Handler, "Submit1")
	submit2Handler := restServer.WrappedRouteHandler(ctrl.Submit2, http.MethodGet, http.StatusOK, paramMap, PDPResponse{})
	router.AddRoute("/submit2/{votingRoundId}/{submitAddress}", submit2Handler, "Submit2")
}
