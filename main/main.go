package main

import (
	"context"
	"fmt"
	"local/fdc/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Prepare context
	// Empty context
	context := context.Background()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)

	// run server
	go server.RunProviderServer(context)

	<-cancelChan
	fmt.Printf("Shutting down server")

}
