package config

import (
	"crypto/sha1"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

var Conf *Config

// LoadConfigFromFile load config.yml file and assign it to Conf.
func LoadConfigFromFile() error {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		errMsg := "error when load yaml file: " + err.Error()
		return errors.New(errMsg)
	}

	if err := yaml.Unmarshal(yamlFile, &Conf); err != nil {
		errMsg := "failed to unmarshal: " + err.Error()
		return errors.New(errMsg)
	}

	if err := Conf.CheckConfigFile(); err != nil {
		errMsg := "failed to check and sanitize config file: " + err.Error()
		return errors.New(errMsg)
	}

	if err := Conf.SetupLogFile(); err != nil {
		errMsg := "failed to setup log fil: " + err.Error()
		return errors.New(errMsg)
	}

	// hash the file to check if maybe there is new config in
	// the future.
	Conf.SHA1 = sha1.Sum(yamlFile)

	return nil
}
