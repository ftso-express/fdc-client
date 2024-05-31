package server

import (
	"flare-common/restServer"
	"flare-common/storage"
	"local/fdc/client/attestation"
	"net/http"
)

const storageSize = 10

type merkleRootStorageObject struct {
	merkleRoot string
	randomNum  string
}

type RootsByAddress map[string]merkleRootStorageObject

type FDCProtocolProviderController struct {
	manager *attestation.Manager
	storage storage.Cyclic[RootsByAddress]
}

func newFDCProtocolProviderController(manager *attestation.Manager) *FDCProtocolProviderController {
	storage := storage.NewCyclic[RootsByAddress](storageSize)

	return &FDCProtocolProviderController{manager: manager, storage: storage}
}

// Registration of routes for the FDC protocol provider
func RegisterFDCProviderRoutes(manager *attestation.Manager, router restServer.Router) {
	// Prepare service controller

	ctrl := newFDCProtocolProviderController(manager)
	paramMap := map[string]string{"votingRoundId": "Voting round ID", "submitAddress": "Submit address"}
	submit1Handler := restServer.GeneralRouteHandler(ctrl.submit1Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submit1/{votingRoundId}/{submitAddress}", submit1Handler, "Submit1")
	submit2Handler := restServer.GeneralRouteHandler(ctrl.submit2Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submit2/{votingRoundId}/{submitAddress}", submit2Handler, "Submit2")
	submitSignaturesHandler := restServer.GeneralRouteHandler(ctrl.submitSignaturesController, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{})
	router.AddRoute("/submitSignatures/{votingRoundId}/{submitAddress}", submitSignaturesHandler, "SubmitSignatures")
}
