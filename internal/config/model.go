package config

import (
	"errors"
	"os"
	"strings"
)

// Repos holds data for every repo.
type Repos struct {
	Name     string `yaml:"name"`
	RootPath string `yaml:"root"`
	Cmd      string `yaml:"opt_cmd"`
}

// Service holds list of data for all repos.
type Service []struct {
	Repos Repos `yaml:"repo"`
}

// Config holds data from config file.
type Config struct {
	EnvIsProd bool
	Env       string  `yaml:"env"`
	Host      string  `yaml:"host"`
	PortNum   uint16  `yaml:"port"`
	Secret    string  `yaml:"secret"`
	Keyword   string  `yaml:"keyword"`
	Usr       string  `yaml:"username"`
	LogDir    string  `yaml:"log"`
	Service   Service `yaml:"service"`
	LogFile   *os.File
	SHA1      [20]byte
}

// CheckConfigFile check and sanitize config file.
func (c *Config) CheckConfigFile() error {
	// Set default env to dev. If env is dev then the bool
	// is false otherwise true.
	if c.Env == "" {
		c.Env = "dev"
	}
	if strings.HasPrefix(c.Env, "dev") {
		c.EnvIsProd = false
	}
	if strings.HasPrefix(c.Env, "prod") {
		c.EnvIsProd = true
	}

	// Set default where this app listen to. Host defaults
	// to localhost and Port default to 5050.
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.PortNum == 0 {
		c.PortNum = 5050
	}

	// Set fields to default 'empty' if not set to make validation in
	// hook handler easier.
	if c.Keyword == "" {
		c.Keyword = "empty"
	}
	if c.Usr == "" {
		c.Usr = "empty"
	}

	// Validate required fields.
	if c.Secret == "" {
		return errors.New("`secret` in config file is empty")
	}

	return nil
}

// SetupLogFile sanitize and setup all things related
// to log file.
func (c *Config) SetupLogFile() error {
	// Set default log dir
	if c.LogDir == "" {
		c.LogDir = "./logs/"
	}

	// Make sure log dir has trailing slash
	if !strings.HasSuffix(c.LogDir, "/") {
		c.LogDir += "/"
	}

	// Make sure log dir already exists
	if err := os.MkdirAll(c.LogDir, 0770); err != nil {
		errMsg := "failed to create log path: " + err.Error()
		return errors.New(errMsg)
	}

	// Open output log file
	fl, err := os.OpenFile(c.LogDir+"log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		errMsg := "failed to open|create log file: " + err.Error()
		return errors.New(errMsg)
	}

	// Assign log file writer to config
	c.LogFile = fl

	return nil
}
