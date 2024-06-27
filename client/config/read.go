package config

import (
	"flare-common/errorf"
	"os"
	"path"
	"strconv"

	"github.com/BurntSushi/toml"
)

func ReadUserRaw(filePath string) (UserRaw, error) {

	return readToml[UserRaw](filePath)

}

func ReadSystem(directory, chain string, protocolId uint8) (System, error) {

	chain = chain + ".toml"

	protocolStr := strconv.FormatUint(uint64(protocolId), 10)

	filePath := path.Join(directory, protocolStr, chain)

	return readToml[System](filePath)

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
