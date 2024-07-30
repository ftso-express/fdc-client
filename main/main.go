package main

import (
	"context"
	"flare-common/logger"
	"local/fdc/client/collector"
	"local/fdc/client/config"
	"local/fdc/client/manager"
	"local/fdc/client/shared"
	"local/fdc/client/timing"
	"local/fdc/server"
	"os"
	"os/signal"
	"syscall"
)

const (
	USER_FILE        string = "configs/userConfig.toml" // relative to project root
	SYSTEM_DIRECTORY string = "configs/systemConfigs"   // relative to project root
)

var log = logger.GetLogger()

func main() {
	userConfigRaw, systemConfig, err := config.GetConfigs(USER_FILE, SYSTEM_DIRECTORY)
	if err != nil {
		log.Panicf("cannot read configs: %s", err)
	}
	err = timing.Set(systemConfig.Timing)
	if err != nil {
		log.Panicf("cannot set timing: %s", err)
	}

	// Prepare context, that can cancel all goroutines
	ctx, cancel := context.WithCancel(context.Background())

	// Prepare shared data connections that collector, manager and server will use
	sharedDataPipes := shared.NewSharedDataPipes()

	// Start attestation client collector
	collector := collector.NewCollector(userConfigRaw, systemConfig, sharedDataPipes)
	go collector.Run(ctx)

	// Start attestation client manager
	manager, err := manager.NewManager(userConfigRaw, sharedDataPipes)
	if err != nil {
		log.Panicf("failed to create the manager: %s", err)
	}
	go manager.Run(ctx)

	// Run attestation clientserver
	srv := server.New(&sharedDataPipes.Rounds, userConfigRaw.RestServer)
	go srv.Run(ctx)
	log.Info("Running server")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)
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
