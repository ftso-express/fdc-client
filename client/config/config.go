package config

import (
	"flare-common/database"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type GlobalConfig struct {
	User   UserConfig
	System SystemConfig
}

type userConfigCommon struct {
	Chain      string               `toml:"chain"`
	DB         database.DBConfig    `toml:"db"`
	RestServer UserRestServerConfig `toml:"rest_server"`
}

type UserConfigRaw struct {
	Abis      AbiConfigUnparsed      `toml:"abis"`
	Verifiers VerifierConfigUnparsed `toml:"verifiers"`
	userConfigCommon
}

type UserConfig struct {
	Abis      AbiConfig
	Verifiers VerifierConfig
	userConfigCommon
}

type SystemConfig struct {
	Logger     LoggerConfig           `toml:"logger"`
	Timing     TimingConfig           `toml:"timing"`
	Listener   ListenersConfig        `toml:"listener"`
	RestServer SystemRestServerConfig `toml:"rest_server"`
}

type SystemRestServerConfig struct {
	Title       string `toml:"title"`
	FSPTitle    string `toml:"fsp_sub_router_title"`
	FSPSubpath  string `toml:"fsp_sub_router_path"`
	Version     string `toml:"version"`
	Addr        string `toml:"addr"`
	SwaggerPath string `toml:"swagger_path"`
}

type UserRestServerConfig struct {
	Addr       string   `toml:"addr"`
	ApiKeyName string   `toml:"api_key_name"`
	ApiKeys    []string `toml:"api_keys"`
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
	Protocol uint8 `toml:"protocol"`

	SubmitContractAddress       common.Address `toml:"submit_contract_address"`
	SubmitBuffer                int            `toml:"submit_buffer"`
	OffChainTriggerDelaySeconds int            `toml:"off_chain_trigger_delay_seconds"`

	RelayContractAddress common.Address `toml:"relay_contract_address"`
	RelayBuffer          int            `toml:"relay_buffer"`

	FdcContractAddress     common.Address `toml:"fdc_contract_address"`
	RequestBuffer          int            `toml:"request_buffer"`
	RequestIntervalSeconds int            `toml:"request_interval_seconds"`
}

type AbiConfigUnparsed map[string]string // map from attestation type to location of abi of its Response struct

type AbiConfig struct {
	ResponseArguments  map[[32]byte]abi.Arguments // map from attestation type to abi of its Response struct
	ResponseAbisString map[[32]byte]string
}

type VerifierCredentialsBig struct {
	Url      string   `toml:"url"`
	ApiKey   string   `toml:"api_key"`
	LutLimit *big.Int `toml:"lut_limit"`
}

type VerifierCredentials struct {
	Url      string
	ApiKey   string
	LutLimit uint64
}

type VerifierConfigUnparsed map[string]map[string]VerifierCredentialsBig // map from attestation type and source Id to verifier credentials

type VerifierConfig map[[64]byte]VerifierCredentials // map from attestation type and source Id to verifier credentials
