package config

import (
	"flare-common/database"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	USER_FILE   string = "userConfig.toml"
	SYSTEM_FILE string = "systemConfig.toml"
)

type UserConfigRaw struct {
	Abis      AbiConfigUnparsed      `toml:"abis"`
	Verifiers VerifierConfigUnparsed `toml:"verifiers"`
	DB        database.DBConfig      `toml:"db"`
}

type SystemConfig struct {
	Logger   LoggerConfig    `toml:"logger"`
	Timing   TimingConfig    `toml:"timing"`
	Listener ListenersConfig `toml:"listener"`
}

type LoggerConfig struct {
	File   string `toml:"file"`
	Prefix string `toml:"prefix"`
	Flag   int    `toml:"flag"`
}

type TimingConfig struct {
	CollectTime uint64
	ChooseTime  uint64
	CommitTime  uint64
	Offset      uint64
	T0          uint64
}

type ListenersConfig struct {
	Protocol uint64 `toml:"protocol"`

	SubmitContractAddress       string `toml:"submit_contract_address"`
	SubmitBuffer                int    `toml:"submit_buffer"`
	OffChainTriggerDelaySeconds int    `toml:"off_chain_trigger_delay_seconds"`

	RelayContractAddress string `toml:"relay_contract_address"`
	RelayBuffer          int    `toml:"relay_buffer"`

	FdcContractAddress     string `toml:"fdc_contract_address"`
	RequestBuffer          int    `toml:"request_buffer"`
	RequestIntervalSeconds int    `toml:"request_interval_seconds"`
}

type AbiConfigUnparsed map[string]string // map from attestation type to abi of its Response struct

type AbiConfig struct {
	ResponseArguments  map[[32]byte]abi.Arguments // map from attestation type to abi of its Response struct
	ResponseAbisString map[[32]byte]string
}

type VerifierCredentials struct {
	Url    string `toml:"url"`
	ApiKey string `toml:"api_key"`
}

type VerifierConfigUnparsed map[string]map[string]VerifierCredentials // map from attestation type and source Id to verifier credentials

type VerifierConfig map[[64]byte]VerifierCredentials // map from attestation type and source Id to verifier credentials
