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
	rounds  *storage.Cyclic[*attestation.Round]
	storage storage.Cyclic[RootsByAddress]
}

func newFDCProtocolProviderController(rounds *storage.Cyclic[*attestation.Round]) *FDCProtocolProviderController {
	storage := storage.NewCyclic[RootsByAddress](storageSize)

	return &FDCProtocolProviderController{rounds: rounds, storage: storage}
}

// Registration of routes for the FDC protocol provider
func RegisterFDCProviderRoutes(rounds *storage.Cyclic[*attestation.Round], router restServer.Router, securities []string) {
	// Prepare service controller

	ctrl := newFDCProtocolProviderController(rounds)
	paramMap := map[string]string{"votingRoundId": "Voting round ID", "submitAddress": "Submit address"}
	submit1Handler := restServer.GeneralRouteHandler(ctrl.submit1Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submit1/{votingRoundId}/{submitAddress}", submit1Handler, "Submit1")
	submit2Handler := restServer.GeneralRouteHandler(ctrl.submit2Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submit2/{votingRoundId}/{submitAddress}", submit2Handler, "Submit2")
	submitSignaturesHandler := restServer.GeneralRouteHandler(ctrl.submitSignaturesController, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submitSignatures/{votingRoundId}/{submitAddress}", submitSignaturesHandler, "SubmitSignatures")
}
