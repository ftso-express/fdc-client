package config

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/BurntSushi/toml"
)

func ReadConfigs(userFilePath, systemFilePath string) (*UserRaw, *System, error) {
	userConfigRaw, err := ReadUserRaw(userFilePath)
	if err != nil {
		return nil, nil, err
	}

	systemConfig, err := ReadSystem(systemFilePath, userConfigRaw.Chain, userConfigRaw.ProtocolID)
	if err != nil {
		return nil, nil, err
	}

	return &userConfigRaw, &systemConfig, nil
}

func ReadUserRaw(filePath string) (UserRaw, error) {
	return readToml[UserRaw](filePath)
}

func ReadSystem(directory, chain string, protocolID uint8) (System, error) {
	chain = chain + ".toml"
	protocolStr := strconv.FormatUint(uint64(protocolID), 10)
	filePath := path.Join(directory, protocolStr, chain)

	return readToml[System](filePath)
}

func readToml[C any](filePath string) (C, error) {
	var config C

	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed reading file %s with: %s", filePath, err)
	}

	err = toml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("failed marshaling file %s with: %s", filePath, err)
	}

	return config, nil
}
