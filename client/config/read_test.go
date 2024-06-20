package config_test

import (
	"local/fdc/client/config"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	USER_FILE   = "../../testFiles/configs/userConfig.toml"   //relative to test
	SYSTEM_FILE = "../../testFiles/configs/systemConfig.toml" //relative to test
)

func TestReadUserRaw(t *testing.T) {

	cfg, err := config.ReadUserRaw(USER_FILE)

	require.NoError(t, err)
	_, err = config.ParseAbi(cfg.Abis)

	require.NoError(t, err)

	_, err = config.ParseVerifiers(cfg.Verifiers)

	require.NoError(t, err)

	require.Equal(t, "coston", cfg.Chain)

}

func TestReadSystem(t *testing.T) {

	_, err := config.ReadSystem(SYSTEM_FILE)

	require.NoError(t, err)

}
