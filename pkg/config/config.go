package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// AppConfig a bag containing all necessary things for this app.
type AppConfig struct {
	Config *viper.Viper
	InfL   *log.Logger
	ErrL   *log.Logger
}

// GetSHA256Signature retrieve hmac from combination of config's secret and the given bytes.
func (a *AppConfig) GetSHA256Signature(in []byte) []byte {
	secret := []byte(a.Config.GetString("secret"))
	mac := hmac.New(sha256.New, secret)
	mac.Write(in)

	return mac.Sum(nil)
}

// InitConfig init config and return preconfigured viper instance.
func InitConfig(filePath string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(filePath)
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

// SetupDefault setup default value and return error if required fields is not present.
func SetupDefault(v *viper.Viper) error {
	if !v.IsSet("secret") {
		return fmt.Errorf("secret is required")
	}
	v.SetDefault("env", "dev")
	v.SetDefault("host", "127.0.0.1")
	v.SetDefault("port", 7575)
	v.SetDefault("log", "/tmp")
	v.SetDefault("max_worker", 1)

	return nil
}
