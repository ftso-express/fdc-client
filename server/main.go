package server

import (
	"context"
	"flare-common/restServer"
	"fmt"
	"local/fdc/client/attestation"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type AttestationServerContext struct {
	Manager *attestation.Manager
	context.Context
}

func RunProviderServer(ctx AttestationServerContext) {

	fmt.Println("In development")

	// Create Mux router
	muxRouter := mux.NewRouter()

	router := restServer.NewSwaggerRouter(muxRouter, "FDC protocol data provider API", "0.0.0")

	// Register routes

	RegisterFDCProviderRoutes(router, ctx)
	router.Finalize()

	// Bind to a port and pass our router in

	var port = "8008"
	// fmt.Printf("Listening on port %s\n", port)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	corsMuxRouter := cors.Handler(muxRouter)
	srv := &http.Server{
		Handler: corsMuxRouter,
		Addr:    fmt.Sprintf(":%s", port),
		// Good practice: enforce timeouts for servers you create -- config?
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}

	go func() {
		fmt.Printf("Starting server on %s", port)
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Printf("Server error: %v", err)
		}
	}()

}
