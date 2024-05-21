package config

import (
	"flare-common/errorf"
	"os"

	"github.com/BurntSushi/toml"
)

func ReadUserRaw() (UserConfigRaw, error) {

	return readToml[UserConfigRaw](USER_FILE)

}

func ReadSystem() (SystemConfig, error) {

	return readToml[SystemConfig](SYSTEM_FILE)

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
