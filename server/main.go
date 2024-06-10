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
	systemServerConfig config.SystemRestServerConfig,
	userServerConfig config.UserRestServerConfig,
) Server {
	// Create Mux router
	muxRouter := mux.NewRouter()

	// create api auth middleware
	keyMiddleware := &restServer.AipKeyAuthMiddleware{
		KeyName: userServerConfig.ApiKeyName,
		Keys:    userServerConfig.ApiKeys,
	}
	keyMiddleware.Init()

	router := restServer.NewSwaggerRouter(muxRouter, restServer.SwaggerRouterConfig{
		Title:           systemServerConfig.Title,
		Version:         systemServerConfig.Version,
		SwaggerBasePath: systemServerConfig.SwaggerPath,
		SecuritySchemes: keyMiddleware.SecuritySchemes(),
	})

	// create fsp sub router
	fspSubRouter := router.WithPrefix(systemServerConfig.FSPSubpath, systemServerConfig.FSPTitle)
	// Register routes for FSP
	RegisterFDCProviderRoutes(rounds, fspSubRouter, []string{userServerConfig.ApiKeyName})
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
		Addr:    userServerConfig.Addr,
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
