package config

import (
	"flare-common/errorf"
	"os"

	"github.com/BurntSushi/toml"
)

func ReadUserRaw(filePath string) (UserConfigRaw, error) {

	return readToml[UserConfigRaw](filePath)

}

func ReadSystem(filePath string) (SystemConfig, error) {

	return readToml[SystemConfig](filePath)

}

func readToml[C any](filePath string) (C, error) {
	var config C

	file, err := os.ReadFile(filePath)

	if err != nil {
		return config, errorf.ReadingFile(filePath, err)
	}

	err = toml.Unmarshal(file, &config)

	if err != nil {
		return config, errorf.Unmarshal(filePath, err)
	}

	return config, nil
}
