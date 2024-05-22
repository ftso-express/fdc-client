package config_test

import (
	"local/fdc/client/config"
	"testing"
)

const (
	USER_FILE   string = "../../testFiles/configs/userConfig.toml"   //relative to test
	SYSTEM_FILE string = "../../testFiles/configs/systemConfig.toml" //relative to test
)

func TestReadUserRaw(t *testing.T) {

	cfg, err := config.ReadUserRaw(USER_FILE)

	if err != nil {
		t.Errorf("error: %s", err)
	}

	_, err = config.ParseAbi(cfg.Abis)

	if err != nil {
		t.Errorf("error: %s", err)
	}

	_, err = config.ParseVerifiers(cfg.Verifiers)

	if err != nil {
		t.Errorf("error: %s", err)
	}

}
