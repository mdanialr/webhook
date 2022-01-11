package config

import (
	"crypto/sha1"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Conf *Config

// LoadConfigFromFile load config.yml file and assign it to Conf.
func LoadConfigFromFile() error {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("error when load yaml file: %v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &Conf); err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	if err := Conf.CheckConfigFile(); err != nil {
		return fmt.Errorf("failed to check and sanitize config file: %v", err)
	}

	if err := Conf.SetupLogFile(); err != nil {
		return fmt.Errorf("failed to setup log file: %v", err)
	}

	// hash the file to check if maybe there is new config in
	// the future.
	Conf.SHA1 = sha1.Sum(yamlFile)

	return nil
}
