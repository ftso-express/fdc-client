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
	var log = logger.GetLogger()

	flag.Parse()
	userConfigRaw, systemConfig, err := config.GetConfigs(*CfgFlag, SYSTEM_DIRECTORY)
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

	// Run attestation client server
	srv := server.New(&sharedDataPipes.Rounds, uint64(userConfigRaw.ProtocolId), userConfigRaw.RestServer)
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
