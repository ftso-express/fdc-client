package main

import (
	"context"
	"flare-common/logger"
	"fmt"
	"local/fdc/client/collector"
	"local/fdc/client/config"
	"local/fdc/server"
	"os"
	"os/signal"
	"syscall"
)

const (
	USER_FILE   string = "configs/userConfig.toml"   //relative to main
	SYSTEM_FILE string = "configs/systemConfig.toml" //relative to main
)

var log = logger.GetLogger()

func main() {
	// Start attestation client collector\

	userConfigRaw, err := config.ReadUserRaw(USER_FILE)
	if err != nil {
		log.Panic("cannot read user config:", err)
	}

	systemConfig, err := config.ReadSystem(SYSTEM_FILE)

	if err != nil {
		log.Panic("cannot read system config:", err)
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
	log.Info("Running server")
	go server.RunProviderServer(context)

	<-cancelChan
	fmt.Println("\nShutting down server")

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
