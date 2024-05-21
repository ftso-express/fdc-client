package main

import (
	"context"
	"fmt"
	"local/fdc/client/collector"
	"local/fdc/client/config"
	"local/fdc/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Start attestation client collector\

	userConfigRaw, err := config.ReadUserRaw()
	if err != nil {
		log.Panicf("cannot read user config: %s", err)
	}

	systemConfig, err := config.ReadSystem()

	if err != nil {
		log.Panicf("cannot read system config: %s", err)
	}

	collector := collector.New(userConfigRaw, systemConfig)
	go collector.Run()

	// Prepare context
	// Empty context
	baseCtx := context.Background()
	context := server.AttestationServerContext{Context: baseCtx, Manager: collector.RoundManager}

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)

	// run server
	log.Println("Running server")
	go server.RunProviderServer(context)

	<-cancelChan
	fmt.Printf("Shutting down server")

}

// func printStructFields(s interface{}) {
// 	// Get the type of the struct
// 	t := reflect.TypeOf(s)

// 	// Iterate over the fields of the struct
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		fmt.Printf("%s: %v\n", field.Name, reflect.ValueOf(s).Field(i))
// 	}
// }
