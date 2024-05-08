package server

import (
	"flare-common/restServer"
	"local/fdc/client/attestation"
	"net/http"
)

type merkleRootStorageObject struct {
	merkleRoot string
	randomNum  string
}

type FDCProtocolProviderController struct {
	manager   *attestation.Manager
	someValue string
	// Mapper from submit address to voting round ID to merkle root storage object
	// Submit address -> Voting round ID -> Merkle root storage object
	rootStorage map[string]map[uint64]merkleRootStorageObject
}

func newFDCProtocolProviderController(ctx AttestationServerContext) *FDCProtocolProviderController {
	return &FDCProtocolProviderController{manager: ctx.Manager, someValue: "initial value"}
}

// Registration of routes for the FDC protocol provider
func RegisterFDCProviderRoutes(router restServer.Router, ctx AttestationServerContext) {
	// Prepare service controller

	ctrl := newFDCProtocolProviderController(ctx)
	paramMap := map[string]string{"votingRoundId": "Voting round ID", "submitAddress": "Submit address"}
	submit1Handler := restServer.GeneralRouteHandler(ctrl.Submit1Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submit1/{votingRoundId}/{submitAddress}", submit1Handler, "Submit1")
	submit2Handler := restServer.GeneralRouteHandler(ctrl.submit2Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submit2/{votingRoundId}/{submitAddress}", submit2Handler, "Submit2")
	submitSignaturesHandler := restServer.GeneralRouteHandler(ctrl.submitSignaturesController, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submitSignatures/{votingRoundId}/{submitAddress}", submitSignaturesHandler, "SubmitSignatures")
}
