package config

import (
	"flare-common/database"
	"time"
)

const (
	FILE string = "config.toml"
)

type FDCConfig struct {
	DBCfg     database.DBConfig
	LoggerCfg LoggerConfig
}

type LoggerConfig struct {
	File string
}

type ListenersConfig struct {
	Protocol uint64

	SubmitContractAddress       string
	SubmitBuffer                int
	OffChainTriggerDelaySeconds int

	RelayContractAddress string
	RelayBuffer          int

	FdcContractAddress string
	RequestBuffer      int
	RequestInterval    time.Duration
}
