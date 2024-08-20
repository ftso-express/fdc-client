package server

import (
	"context"
	"flare-common/restServer"
	"flare-common/storage"
	"local/fdc/client/config"
	"local/fdc/client/round"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	srv *http.Server
}

func New(
	rounds *storage.Cyclic[*round.Round],
	protocolId uint64,
	serverConfig config.RestServer,
) Server {
	// Create Mux router
	muxRouter := mux.NewRouter()

	// Register a health check endpoint at the top level.
	muxRouter.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// create api auth middleware
	keyMiddleware := &restServer.AipKeyAuthMiddleware{
		KeyName: serverConfig.ApiKeyName,
		Keys:    serverConfig.ApiKeys,
	}
	keyMiddleware.Init()

	router := restServer.NewSwaggerRouter(muxRouter, restServer.SwaggerRouterConfig{
		Title:           serverConfig.Title,
		Version:         serverConfig.Version,
		SwaggerBasePath: serverConfig.SwaggerPath,
		SecuritySchemes: keyMiddleware.SecuritySchemes(),
	})

	// create fsp sub router
	fspSubRouter := router.WithPrefix(serverConfig.FSPSubpath, serverConfig.FSPTitle)
	// Register routes for FSP
	RegisterFDCProviderRoutes(fspSubRouter, protocolId, rounds, []string{serverConfig.ApiKeyName})
	fspSubRouter.AddMiddleware(keyMiddleware.Middleware)

	// Register routes
	router.Finalize()

	// Create CORS handler
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	corsMuxRouter := cors.Handler(muxRouter)
	srv := &http.Server{
		Handler: corsMuxRouter,
		Addr:    serverConfig.Addr,
		// Good practice: enforce timeouts for servers you create -- config?
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}

	return Server{srv: srv}
}

// Registration of routes for the FDC protocol provider
func RegisterFDCProviderRoutes(router restServer.Router, protocolId uint64, rounds *storage.Cyclic[*round.Round], securities []string) {
	// Prepare service controller
	controller := newFDCProtocolProviderController(rounds, protocolId)
	paramMap := map[string]string{"votingRoundId": "Voting round ID", "submitAddress": "Submit address"}

	submit1Handler := restServer.GeneralRouteHandler(controller.submit1Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submit1/{votingRoundId}/{submitAddress}", submit1Handler, "Submit1")

	submit2Handler := restServer.GeneralRouteHandler(controller.submit2Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submit2/{votingRoundId}/{submitAddress}", submit2Handler, "Submit2")

	submitSignaturesHandler := restServer.GeneralRouteHandler(controller.submitSignaturesController, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submitSignatures/{votingRoundId}/{submitAddress}", submitSignaturesHandler, "SubmitSignatures")
}

// Registration of routes for the DA layer
func RegisterDARoutes(router restServer.Router, rounds *storage.Cyclic[*round.Round], securities []string) {
	// Prepare service controller
	controller := FDCDAController{rounds: rounds}
	paramMap := map[string]string{"votingRoundId": "Voting round ID"}

	getRequests := restServer.GeneralRouteHandler(controller.getRequestController, http.MethodGet, http.StatusOK, paramMap, nil, nil, RequestsResponse{}, securities)
	router.AddRoute("/getRequests/{votingRoundId}", getRequests, "GetRequests")

	getAttestations := restServer.GeneralRouteHandler(controller.getAttestationController, http.MethodGet, http.StatusOK, paramMap, nil, nil, AttestationResponse{}, securities)
	router.AddRoute("/getAttestations/{votingRoundId}", getAttestations, "GetAttestations")
}

func (s *Server) Run(ctx context.Context) {
	log.Infof("Starting server on %s", s.srv.Addr)

	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panicf("Server error: %v", err)
	}
}

var shutdownTimeout = 5 * time.Second

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed:", err)
	} else {
		log.Info("Server gracefully stopped")
	}
}
