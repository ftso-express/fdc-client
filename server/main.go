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

	router := restServer.NewSwaggerRouter(muxRouter, "FDC protocol data provider API", "0.0.0")

	// Register routes

	RegisterFDCProviderRoutes(router, rounds)
	router.Finalize()

	// Bind to a port and pass our router in

	// fmt.Printf("Listening on port %s\n", port)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

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
