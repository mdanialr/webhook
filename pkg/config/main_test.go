package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

var (
	viperWithoutSecret *viper.Viper
	viperComplete      *viper.Viper
	viperInit          *viper.Viper
)

func TestMain(m *testing.M) {
	setUp()
	out := m.Run()
	cleanUp()
	os.Exit(out)
}

func setUp() {
	viperWithoutSecret = setupWithoutSecret()
	viperComplete = setupComplete()

	os.Create("/tmp/app.yml")
}

func cleanUp() {
	os.Remove("/tmp/app.yml")
}

func setupWithoutSecret() *viper.Viper {
	v := viper.New()
	v.SetConfigType("json")
	jsonTest := `{"repo":"hello-world"}`
	v.ReadConfig(strings.NewReader(jsonTest))
	return v
}

func setupComplete() *viper.Viper {
	v := viper.New()
	v.SetConfigType("json")
	jsonTest := `{"secret":"abc"}`
	v.ReadConfig(strings.NewReader(jsonTest))
	return v
}
