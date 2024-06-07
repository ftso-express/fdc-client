package main

import (
	"context"
	"flare-common/logger"
	"local/fdc/client/collector"
	"local/fdc/client/config"
	"local/fdc/server"
	"os"
	"os/signal"
	"syscall"
)

const (
	USER_FILE   string = "configs/userConfig.toml"   //relative to project root
	SYSTEM_FILE string = "configs/systemConfig.toml" //relative to project root
)

var log = logger.GetLogger()

func main() {
	// Start attestation client collector\

	userConfigRaw, err := config.ReadUserRaw(USER_FILE)
	if err != nil {
		log.Panicf("cannot read user config: %s", err)
	}

	systemConfig, err := config.ReadSystem(SYSTEM_FILE)

	if err != nil {
		log.Panicf("cannot read system config: %s", err)
	}

	// Prepare context
	ctx, cancel := context.WithCancel(context.Background())

	collector := collector.New(userConfigRaw, systemConfig)
	go collector.Run(ctx)

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)

	// run server
	log.Info("Running server")
	srv := server.New(&collector.RoundManager.Rounds, systemConfig.RestServer, userConfigRaw.RestServer)
	go srv.Run(ctx)

	// Block until a termination signal is received.
	select {
	case <-cancelChan:
		log.Info("Received an interrupt signal, shutting down...")
	case <-ctx.Done():
		log.Info("Context cancelled, shutting down...")
	}

	cancel()
	srv.Shutdown()
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
