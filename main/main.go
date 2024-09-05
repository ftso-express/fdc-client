package main

import (
	"context"
	"flag"
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
	SYSTEM_DIRECTORY string = "configs/systemConfigs" // relative to project root
)

var CfgFlag = flag.String("config", "configs/userConfig.toml", "Configuration file (toml format)")

func main() {
	flag.Parse()
	userConfigRaw, systemConfig, err := config.ReadConfigs(*CfgFlag, SYSTEM_DIRECTORY)
	if err != nil {
		logger.Panicf("cannot read configs: %s", err)
	}
	logger.SetLogger(userConfigRaw.Logging)

	err = timing.Set(systemConfig.Timing)
	if err != nil {
		logger.Panicf("cannot set timing: %s", err)
	}
	attestationTypeConfig, err := config.ParseAttestationTypes(userConfigRaw.AttestationTypeConfig)
	if err != nil {
		logger.Panicf("att types config: %s", err)
	}

	// Prepare context, that can cancel all goroutines
	ctx, cancel := context.WithCancel(context.Background())

	// Prepare shared data connections that collector, manager and server will use
	sharedDataPipes := shared.NewDataPipes()

	// Start attestation client collector
	collector := collector.New(userConfigRaw, systemConfig, sharedDataPipes)
	go collector.Run(ctx)

	// Start attestation client manager
	manager, err := manager.New(userConfigRaw, attestationTypeConfig, sharedDataPipes)
	if err != nil {
		logger.Panicf("failed to create the manager: %s", err)
	}
	go manager.Run(ctx)

	// Run attestation client server
	srv := server.New(&sharedDataPipes.Rounds, userConfigRaw.ProtocolID, userConfigRaw.RestServer)
	go srv.Run(ctx)
	logger.Info("Running server")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM)
	// Block until a termination signal is received.
	select {
	case <-cancelChan:
		logger.Info("Received an interrupt signal, shutting down...")
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down...")
	}

	cancel()
	srv.Shutdown()
}
