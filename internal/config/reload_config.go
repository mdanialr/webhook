package config

import (
	"crypto/sha1"
	"github.com/mdanialr/webhook/internal/helpers"
	"os"
)

// ReloadConfig reload and repopulate config when there is new
// value in config file by checking old hash against new hash.
func ReloadConfig() error {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	sha1Res := sha1.Sum(f)
	if Conf.SHA1 != sha1Res {
		helpers.NzLogInf.Println("config file has different hash. repopulate config will be initiated")
		if err := LoadConfigFromFile(); err != nil {
			return err
		}
	}

	return nil
}
