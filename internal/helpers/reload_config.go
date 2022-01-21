package helpers

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/mdanialr/webhook/internal/config"
)

var configFilePath = "config.yaml"

// ReloadConfig reload and repopulate config when there is new
// value in config file by checking old hash against new hash.
func ReloadConfig(conf config.Config) error {
	f, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to open and read config.yaml file: %v", err)
	}

	sha1Res := sha1.Sum(f)
	if conf.SHA1 != sha1Res {
		NzLogInf.Println("config file has different hash. repopulate config will be initiated")
		if err := conf.LoadConfigFromFile(configFilePath); err != nil {
			return fmt.Errorf("failed to load config file triggered by reload config: %v", err)
		}
	}

	return nil
}
