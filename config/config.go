package config

import (
	"errors"
	"flare-common/database"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

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

type FDCConfig struct {
	Logger   LoggerConfig
	Timing   TimingConfig
	Listener ListenersConfig
}

type LoggerConfig struct {
	File string
}

type TimingConfig struct {
	CollectTime uint64
	ChooseTime  uint64
	CommitTime  uint64
	Offset      uint64
	T0          uint64
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

type AbiConfigUnparsed map[string]string // map from attestation type to abi of its Response struct

type AbiConfig struct {
	ResponseArguments  map[[32]byte]abi.Arguments // map from attestation type to abi of its Response struct
	ResponseAbisString map[[32]byte]string
}

func WhiteSpaceStrip(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func ParseAbiConfig(data AbiConfigUnparsed) (AbiConfig, error) {

	arguments := make(map[[32]byte]abi.Arguments)
	abis := make(map[[32]byte]string)

	for k, v := range data {

		if len(k) > 32 {
			return AbiConfig{}, errors.New("attestation type name too long")
		}

		var nameBytes [32]byte

		copy(nameBytes[:], []byte(k)[0:32])

		var arg abi.Argument

		file, err := os.ReadFile(v)

		if err != nil {
			return AbiConfig{}, fmt.Errorf("error reading file %s", v)
		}

		err = arg.UnmarshalJSON(file)

		if err != nil {
			return AbiConfig{}, errors.New("error parsing abi")
		}

		args := abi.Arguments{arg}

		arguments[nameBytes] = args

		abis[nameBytes] = WhiteSpaceStrip(string(file))

	}

	return AbiConfig{arguments, abis}, nil

}

type VerifierCredentials struct {
	Url    string `toml:"url"`
	ApiKey string `toml:"api_key"`
}

type VerifierConfigUnparsed map[string]map[string]VerifierCredentials // map from hash of attestation type and source Id to verifier credentials

type VerifierConfig struct {
	VerifiersCredentials map[string]VerifierCredentials // map from hash of attestation type and source Id to verifier credentials
}
