package helpers

import (
	"crypto/sha1"
	"github.com/mdanialr/webhook/internal/config"
	"gopkg.in/yaml.v3"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReloadConfig(t *testing.T) {
	var tConf *config.Config
	configFilePath = "./../../config.yaml"

	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		require.NoError(t, err)
	}

	if err := yaml.Unmarshal(yamlFile, &tConf); err != nil {
		require.NoError(t, err)
	}

	if err := tConf.CheckConfigFile(); err != nil {
		require.NoError(t, err)
	}

	if err := tConf.SetupLogFile(); err != nil {
		require.NoError(t, err)
	}

	// hash the file to check if maybe there is new config in
	// the future.
	tConf.SHA1 = sha1.Sum(yamlFile)

	err = ReloadConfig(*tConf)
	require.NoError(t, err)
}
