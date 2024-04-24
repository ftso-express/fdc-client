package main

import (
	"context"
	"flare-common/restServer"
	"fmt"
	"local/fdc/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	fmt.Println("In development")

	// Empty context
	context := context.Background()

	// Create Mux router
	muxRouter := mux.NewRouter()

	router := restServer.NewSwaggerRouter(muxRouter, "FDC protocol data provider API", "0.0.0")

	// Register routes

	server.RegisterFDCProviderRoutes(router, context)
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

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Starting server on %s", port)
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Printf("Server error: %v", err)
		}
	}()

	<-cancelChan
	fmt.Printf("Shutting down server")

}
