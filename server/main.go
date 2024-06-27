package server

import (
	"context"
	"flare-common/restServer"
	"flare-common/storage"
	"local/fdc/client/attestation"
	"local/fdc/client/config"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	srv *http.Server
}

func New(
	rounds *storage.Cyclic[*attestation.Round],
	serverConfig config.RestServer,
) Server {
	// Create Mux router
	muxRouter := mux.NewRouter()

	// Register a healthcheck endpoint at the top level.
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
	RegisterFDCProviderRoutes(rounds, fspSubRouter, []string{serverConfig.ApiKeyName})
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
