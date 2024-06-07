package server

import (
	"context"
	"flare-common/restServer"
	"flare-common/storage"
	"local/fdc/client/attestation"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	srv *http.Server
}

func New(address string, rounds *storage.Cyclic[*attestation.Round]) Server {
	// Create Mux router
	muxRouter := mux.NewRouter()

	// create api auth middleware
	keyMiddleware := &restServer.AipKeyAuthMiddleware{
		KeyName: "X-API-KEY",
		Keys:    []string{"123456"},
	}
	keyMiddleware.Init()

	router := restServer.NewSwaggerRouter(muxRouter, restServer.SwaggerRouterConfig{
		Title:           "FDC protocol data provider API",
		Version:         "0.0.0",
		SwaggerBasePath: "/api-doc",
		SecuritySchemes: keyMiddleware.SecuritySchemes(),
	})

	// create fsp sub router
	fspSubRouter := router.WithPrefix("/fsp", "FDC protocol data provider for FSP client")
	// Register routes for FSP
	RegisterFDCProviderRoutes(rounds, fspSubRouter, []string{"X-API-KEY"})
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
		Addr:    address,
		// Good practice: enforce timeouts for servers you create -- config?
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}

	return Server{srv: srv}
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
