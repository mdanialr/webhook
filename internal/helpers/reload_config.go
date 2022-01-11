package helpers

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/mdanialr/webhook/internal/config"
)

// ReloadConfig reload and repopulate config when there is new
// value in config file by checking old hash against new hash.
func ReloadConfig() error {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to open and read config.yaml file: %v", err)
	}

	sha1Res := sha1.Sum(f)
	if config.Conf.SHA1 != sha1Res {
		NzLogInf.Println("config file has different hash. repopulate config will be initiated")
		if err := config.LoadConfigFromFile(); err != nil {
			return fmt.Errorf("failed to load config file triggered by reload config: %v", err)
		}
	}

	return nil
}
