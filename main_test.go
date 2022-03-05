package main

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeReader struct{}

func (f fakeReader) Read(_ []byte) (int, error) {
	return 0, errors.New("this should trigger error")
}

func TestSetup(t *testing.T) {
	fakeConfigFile :=
		`
env: prod
port: 5454
secret: secret
log: /tmp
service:
  - repo:
      name: fiber-ln
      root: /home/nzk/dir/Fiber/light_novel/
      opt_cmd: "go build -o bin/fiber-ln main.go && systemctl restart fiber-ln"
  - repo:
      name: cd_test
      root: /home/nzk/dir/Laravel/cd_test/
      opt_cmd: pwd
`
	var appConf config.Model

	t.Run("Using fake interface should return error", func(t *testing.T) {
		_, err := setup(&appConf, fakeReader{})
		require.Error(t, err)
	})

	t.Run("Success must exactly the same as in config file", func(t *testing.T) {
		_, err := setup(&appConf, bytes.NewBufferString(fakeConfigFile))
		require.NoError(t, err)

		assert.Equal(t, "localhost", appConf.Host)
		assert.Equal(t, strconv.Itoa(int(uint16(5454))), strconv.Itoa(int(appConf.PortNum)))
		assert.Equal(t, "/tmp/", appConf.LogDir)
	})

	fakeConfigFile =
		`
env: prod
port: 5050
log: /tmp
`

	t.Run("Secret does not exist on config file should return error", func(t *testing.T) {
		_, err := setup(&appConf, bytes.NewBufferString(fakeConfigFile))
		require.Error(t, err)
	})

	fakeConfigFile =
		`
env: prod
port: 5050
log: /fake/dir
`

	t.Run("Log dir does not exist should return error", func(t *testing.T) {
		_, err := setup(&appConf, bytes.NewBufferString(fakeConfigFile))
		require.Error(t, err)
	})
}

func TestLogWriterFromChannels(t *testing.T) {
	ch := &worker.Channel{
		JobC: make(chan string, 10),
		InfC: make(chan string, 10),
		ErrC: make(chan string, 10),
	}
	dCh := &worker.DockerChannel{
		JobC: make(chan string, 10),
		InfC: make(chan string, 10),
		ErrC: make(chan string, 10),
	}

	logWriterFromChannel(ch, dCh)
}
