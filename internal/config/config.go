package config

import (
	"crypto/sha1"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Conf Config

// LoadConfigFromFile load config.yml file and assign it to Conf.
func (c *Config) LoadConfigFromFile(yamlFilePath string) error {
	yamlFile, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("error when load yaml file: %v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	if err := c.CheckConfigFile(); err != nil {
		return fmt.Errorf("failed to check and sanitize config file: %v", err)
	}

	if err := c.SetupLogFile(); err != nil {
		return fmt.Errorf("failed to setup log file: %v", err)
	}

	// hash the file to check if maybe there is new config in
	// the future.
	c.SHA1 = sha1.Sum(yamlFile)

	return nil
}
