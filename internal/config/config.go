package config

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mdanialr/webhook/internal/service"
	"gopkg.in/yaml.v3"
)

// Model holds data from config file.
type Model struct {
	EnvIsProd bool
	Env       string        `yaml:"env"`
	Host      string        `yaml:"host"`
	PortNum   uint16        `yaml:"port"`
	Secret    string        `yaml:"secret"`
	Keyword   string        `yaml:"keyword"`
	Usr       string        `yaml:"username"`
	LogDir    string        `yaml:"log"`
	MaxWorker int           `yaml:"max_worker"`
	Service   service.Model `yaml:"service"`
	LogFile   *os.File
}

// NewConfig read io.Reader then map and load the value to the returned Model.
func NewConfig(fileBuf io.Reader) (mod *Model, err error) {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(fileBuf); err != nil {
		return mod, fmt.Errorf("failed to read from file buffer: %v", err)
	}

	if err := yaml.Unmarshal(buf.Bytes(), &mod); err != nil {
		return mod, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return
}

// IsDifferentHash hash (md5) two input then compare it. res will always be
// false if there are errors. Otherwise, depend on the hash compare result
func IsDifferentHash(x io.Reader, y io.Reader) (res bool, err error) {
	xH, yH := md5.New(), md5.New()

	if _, err := io.Copy(xH, x); err != nil {
		return true, fmt.Errorf("failed copying first file: %v", err)
	}
	if _, err := io.Copy(yH, y); err != nil {
		return true, fmt.Errorf("failed copying second file: %v", err)
	}

	return bytes.Compare(xH.Sum(nil), yH.Sum(nil)) == 0, nil
}

// ReloadConfig reload and repopulate config from given io.Reader.
func (m *Model) ReloadConfig(fileBuf io.Reader) error {
	newM, err := NewConfig(fileBuf)
	if err != nil {
		return fmt.Errorf("failed to create new config from input file buffer: %v", err)
	}

	m = newM

	return nil
}

// Sanitization check and sanitize config Model's instance.
func (m *Model) Sanitization() error {
	if m.Env == "" || (m.Env != "dev" && m.Env != "prod") {
		m.Env = "dev"
	}
	if strings.HasPrefix(m.Env, "prod") {
		m.EnvIsProd = true
	}

	if m.Host == "" {
		m.Host = "localhost"
	}

	if m.PortNum == 0 {
		m.PortNum = 5050
	}

	if m.Keyword == "" {
		m.Keyword = "empty"
	}

	if m.Usr == "" {
		m.Usr = "empty"
	}

	if m.MaxWorker <= 0 {
		m.MaxWorker = 1
	}

	if m.Secret == "" {
		return fmt.Errorf("`secret` is required")
	}

	return nil
}

// SanitizationLog check and sanitize things related to log.
func (m *Model) SanitizationLog() {
	if m.LogDir == "" {
		m.LogDir = "./log/"
	}

	if !strings.HasSuffix(m.LogDir, "/") {
		m.LogDir += "/"
	}
}

// GetSHA256Signature get hmac hash from combination of Model's secret and
// input bytes.
func (m *Model) GetSHA256Signature(in []byte) []byte {
	secret := []byte(m.Secret)
	mac := hmac.New(sha256.New, secret)
	mac.Write(in)

	return mac.Sum(nil)
}
