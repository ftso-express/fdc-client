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
	rounds *storage.Cyclic[*round.Round, uint32],
	protocolID uint8,
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
		KeyName: serverConfig.APIKeyName,
		Keys:    serverConfig.APIKeys,
	}
	keyMiddleware.Init()

	router := restServer.NewSwaggerRouter(muxRouter, restServer.SwaggerRouterConfig{
		Title:           serverConfig.Title,
		Version:         serverConfig.Version,
		SwaggerBasePath: serverConfig.SwaggerPath,
		SecuritySchemes: keyMiddleware.SecuritySchemes(),
	})

	// create FSP sub router
	fspSubRouter := router.WithPrefix(serverConfig.FSPSubpath, serverConfig.FSPTitle)
	// Register routes for FSP
	RegisterFDCProviderRoutes(fspSubRouter, protocolID, rounds, []string{serverConfig.APIKeyName})
	fspSubRouter.AddMiddleware(keyMiddleware.Middleware)

	// create DA sub router
	daSubRouter := router.WithPrefix(serverConfig.DAPSubpath, serverConfig.DATitle)
	// Register routes for DA
	RegisterDARoutes(daSubRouter, rounds, []string{serverConfig.APIKeyName})
	daSubRouter.AddMiddleware(keyMiddleware.Middleware)

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
func RegisterFDCProviderRoutes(router restServer.Router, protocolID uint8, rounds *storage.Cyclic[*round.Round, uint32], securities []string) {
	// Prepare service controller
	controller := newFDCProtocolProviderController(rounds, protocolID)
	paramMap := map[string]string{"votingRoundID": "Voting round ID", "submitAddress": "Submit address"}

	submit1Handler := restServer.GeneralRouteHandler(controller.submit1Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submit1/{votingRoundID}/{submitAddress}", submit1Handler, "Submit1")

	submit2Handler := restServer.GeneralRouteHandler(controller.submit2Controller, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submit2/{votingRoundID}/{submitAddress}", submit2Handler, "Submit2")

	submitSignaturesHandler := restServer.GeneralRouteHandler(controller.submitSignaturesController, http.MethodGet, http.StatusOK, paramMap, nil, nil, PDPResponse{}, securities)
	router.AddRoute("/submitSignatures/{votingRoundID}/{submitAddress}", submitSignaturesHandler, "SubmitSignatures")
}

// Registration of routes for the DA layer WIP
func RegisterDARoutes(router restServer.Router, rounds *storage.Cyclic[*round.Round, uint32], securities []string) {
	// Prepare service controller
	controller := DAController{Rounds: rounds}
	paramMap := map[string]string{"votingRoundID": "Voting round ID"}

	getRequests := restServer.GeneralRouteHandler(controller.getRequestController, http.MethodGet, http.StatusOK, paramMap, nil, nil, RequestsResponse{}, securities)
	router.AddRoute("/getRequests/{votingRoundID}", getRequests, "GetRequests")

	getAttestations := restServer.GeneralRouteHandler(controller.getAttestationController, http.MethodGet, http.StatusOK, paramMap, nil, nil, AttestationResponse{}, securities)
	router.AddRoute("/getAttestations/{votingRoundID}", getAttestations, "GetAttestations")
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
