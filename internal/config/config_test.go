package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mdanialr/webhook/internal/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeReader just fake type to satisfies io.Reader interfaces
// so it could trigger error buffer read from.
type fakeReader struct{}

func (_ fakeReader) Read(_ []byte) (_ int, _ error) {
	return 0, fmt.Errorf("this should trigger error in test")
}

func TestNewConfig(t *testing.T) {
	fakeConfigFile := `
env: prod
port: 5005
secret: mM0B9VclhC2tH51RKUonN2NlQOPOp5RpqroyrO7n68hnVSvli8
log: /home/nzk/test-app/webhook/log
`
	buf := bytes.NewBufferString(fakeConfigFile)
	t.Run("Using valid value and left-out optional params should be pass", func(t *testing.T) {
		mod, err := NewConfig(buf)
		require.NoError(t, err)

		assert.Equal(t, "prod", mod.Env)
		assert.Equal(t, "", mod.Host)
		assert.Equal(t, uint16(5005), mod.PortNum)
		assert.Equal(t, "mM0B9VclhC2tH51RKUonN2NlQOPOp5RpqroyrO7n68hnVSvli8", mod.Secret)
		assert.Equal(t, "/home/nzk/test-app/webhook/log", mod.LogDir)
	})

	fakeConfigFile = `
env: dev
max_worker: 1
service:
  - repo:
      user: user
      name: fiber-ln
      path: /home/nzk/dir/Fiber/light_novel/
      opt_cmd: "go build -o bin/fiber-ln main.go && systemctl restart fiber-ln"
  - repo:
      user: user
      name: cd_test
      path: /home/nzk/dir/Laravel/cd_test/
      opt_cmd: pwd
`
	buf = bytes.NewBufferString(fakeConfigFile)
	t.Run("Using valid value should be pass", func(t *testing.T) {
		mod, err := NewConfig(buf)
		require.NoError(t, mod.Service.Sanitization())
		require.NoError(t, err)

		assert.Equal(t, "dev", mod.Env)
		assert.Equal(t, 1, mod.MaxWorker)
		assert.Equal(t, "fiber-ln", mod.Service[0].Repo.Name)
		assert.Equal(t, "/home/nzk/dir/Fiber/light_novel/", mod.Service[0].Repo.Path)
		assert.Equal(t, "go build -o bin/fiber-ln main.go && systemctl restart fiber-ln", mod.Service[0].Repo.Cmd)
		assert.Equal(t, "cd_test", mod.Service[1].Repo.Name)
		assert.Equal(t, "pwd", mod.Service[1].Repo.Cmd)
	})

	fakeConfigFile = `
env: 2134
host: 313
port: "number"
keyword: &k0
max_worker: "lol"
	`
	buf = bytes.NewBufferString(fakeConfigFile)
	t.Run("Using mismatch type should be error in yaml unmarshalling", func(t *testing.T) {
		_, err := NewConfig(buf)
		require.Error(t, err)
	})

	t.Run("Injecting fake reader should be error in buffer read from", func(t *testing.T) {
		_, err := NewConfig(fakeReader{})
		require.Error(t, err)
	})

	fakeConfigFile = `port: 123`
	buf = bytes.NewBufferString(fakeConfigFile)
	t.Run("Port w 1234 should be no error", func(t *testing.T) {
		_, err := NewConfig(buf)
		require.NoError(t, err)
	})
}

func TestIsDifferentHash(t *testing.T) {
	fakeConfigOne := `this is the first file`
	fakeConfigTwo := `this is the first file`
	bufOne := bytes.NewBufferString(fakeConfigOne)
	bufTwo := bytes.NewBufferString(fakeConfigTwo)

	t.Run("Using same file should be equal", func(t *testing.T) {
		out, err := IsDifferentHash(bufOne, bufTwo)
		require.NoError(t, err)
		assert.True(t, out)
	})

	fakeConfigTwo = `this is the real second file`
	bufTwo = bytes.NewBufferString(fakeConfigTwo)
	t.Run("Using different file should not be equal", func(t *testing.T) {
		out, err := IsDifferentHash(bufOne, bufTwo)
		require.NoError(t, err)
		assert.False(t, out)
	})

	t.Run("Injecting fake reader should be error in copying first file", func(t *testing.T) {
		_, err := IsDifferentHash(fakeReader{}, bufTwo)
		require.Error(t, err)
	})

	t.Run("Injecting fake reader should be error in copying second file", func(t *testing.T) {
		_, err := IsDifferentHash(bufOne, fakeReader{})
		require.Error(t, err)
	})
}

func TestReloadConfig(t *testing.T) {
	oldMod := Model{
		Env:     "dev",
		PortNum: 5005,
		Secret:  "secret",
		LogDir:  "/home/nzk/test-app/webhook/log",
		Service: github.Service{
			{Repo: github.Model{
				Name: "fiber-ln",
				Path: "/home/nzk/dir/Fiber/light_novel/",
				Cmd:  "go build -o bin/fiber-ln main.go && systemctl restart fiber-ln",
			}},
		},
	}

	// Pretend we already have one model. Then pretend that we will
	// load new one config file and repopulate old model with the
	// newly loaded. So we can compare the old model vs new model
	// to make sure that the old model successfully reloaded.

	newFakeConfigFile := `
env: prod
port: 5050
secret: secret
log: /var/log/webhook/log
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
	buf := bytes.NewBufferString(newFakeConfigFile)

	newMod, err := NewConfig(buf)
	require.NoError(t, err)

	t.Run("Should be no error. Then compare old v new", func(t *testing.T) {
		err := oldMod.ReloadConfig(buf)
		require.NoError(t, err)

		assert.NotEqual(t, newMod.Env, oldMod.Env)
		assert.NotEqual(t, newMod.PortNum, oldMod.PortNum)
		assert.Equal(t, newMod.Secret, oldMod.Secret)
		assert.NotEqual(t, newMod.LogDir, oldMod.LogDir)
		assert.NotEqual(t, len(newMod.Service), len(oldMod.Service))
		assert.Equal(t, newMod.Service[0].Repo.Cmd, oldMod.Service[0].Repo.Cmd)
	})

	t.Run("Injecting fake reader should be error", func(t *testing.T) {
		err := oldMod.ReloadConfig(fakeReader{})
		require.Error(t, err)
	})
}

func TestSanitization_Env(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Env w dev should be dev",
			sample: Model{Env: "dev", Secret: "lol"},
			expect: "dev",
		},
		{
			name:   "Env w prod should be prod",
			sample: Model{Env: "prod", Secret: "lol"},
			expect: "prod",
		},
		{
			name:   "Env w/o value should be dev",
			sample: Model{Secret: "lol"},
			expect: "dev",
		},
		{
			name:   "Env w value not match either dev or prod should be dev",
			sample: Model{Env: "uu", Secret: "lol"},
			expect: "dev",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.Env)
		})
	}
}

func TestSanitization_EnvIsProd(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect bool
	}{
		{
			name:   "Env w dev should be false",
			sample: Model{Env: "dev", Secret: "lol"},
			expect: false,
		},
		{
			name:   "Env w value not match either dev or prod should be false",
			sample: Model{Env: "lol", Secret: "lol"},
			expect: false,
		},
		{
			name:   "Env w/o should be false",
			sample: Model{Secret: "lol"},
			expect: false,
		},
		{
			name:   "Env w prod should be true",
			sample: Model{Env: "prod", Secret: "lol"},
			expect: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.EnvIsProd)
		})
	}
}

func TestSanitization_Host(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Host w localhost should be localhost",
			sample: Model{Host: "localhost", Secret: "lol"},
			expect: "localhost",
		},
		{
			name:   "Host w 120.222.46.23 should be 120.222.46.23",
			sample: Model{Host: "120.222.46.23", Secret: "lol"},
			expect: "120.222.46.23",
		},
		{
			name:   "Host w/o value should be localhost",
			sample: Model{Secret: "lol"},
			expect: "localhost",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.Host)
		})
	}
}

func TestSanitization_Port(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect uint16
	}{
		{
			name:   "Port w 2003 should be 2003",
			sample: Model{PortNum: 2003, Secret: "lol"},
			expect: 2003,
		},
		{
			name:   "Port w 44444 should be 44444",
			sample: Model{PortNum: 44444, Secret: "lol"},
			expect: 44444,
		},
		{
			name:   "Port w/o value should be 5050",
			sample: Model{Secret: "lol"},
			expect: 5050,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.PortNum)
		})
	}
}

func TestSanitization_MaxWorker(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect int
	}{
		{
			name:   "MaxWorker w 4 should be 4",
			sample: Model{MaxWorker: 4, Secret: "lol"},
			expect: 4,
		},
		{
			name:   "MaxWorker w/o value should be 1",
			sample: Model{Secret: "lol"},
			expect: 1,
		},
		{
			name:   "MaxWorker w less than or equal to 0 should be 1",
			sample: Model{MaxWorker: -2, Secret: "lol"},
			expect: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			require.NoError(t, err)
			assert.Equal(t, tt.expect, tt.sample.MaxWorker)
		})
	}
}

func TestSanitization_Secret(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
		isErr  bool
	}{
		{
			name:   "Secret w '12345678' should be '12345678' and has no error",
			sample: Model{Secret: "12345678"},
			expect: "12345678",
			isErr:  false,
		},
		{
			name:   "Secret w '#*&($@NfaSfCSwadw' should be '#*&($@NfaSfCSwadw' and has no error",
			sample: Model{Secret: "#*&($@NfaSfCSwadw"},
			expect: "#*&($@NfaSfCSwadw",
			isErr:  false,
		},
		{
			name:  "Secret w/o value should be error because its required",
			isErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sample.Sanitization()
			switch tt.isErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tt.expect, tt.sample.Secret)
			case true:
				require.Error(t, err)
			}
		})
	}
}

func TestGetSHA256Signature(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		input  []byte
		expect []byte
	}{
		{
			name:   "1# Should always return byte w 32 length",
			sample: Model{Secret: "123"},
			input:  []byte{51, 52, 53},
			expect: []byte{72, 206, 203, 202, 192, 235, 177, 216, 188, 12, 57, 93, 92, 199, 66, 200, 240, 234, 245, 229, 150, 150, 168, 224, 70, 47, 250, 117, 153, 7, 129, 223},
		},
		{
			name:   "2# Should always return byte w 32 length",
			sample: Model{Secret: "321"},
			expect: []byte{134, 103, 135, 124, 79, 17, 222, 100, 44, 154, 116, 140, 56, 163, 126, 78, 138, 197, 89, 119, 86, 227, 21, 43, 235, 187, 140, 83, 103, 162, 181, 236},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.sample.GetSHA256Signature(tt.input)
			assert.Equal(t, tt.expect, out)
		})
	}
}

func TestSanitizationLog(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Log w /var/www/log/ should be /var/www/log/",
			sample: Model{Secret: "lol", LogDir: "/var/www/log/"},
			expect: "/var/www/log/",
		},
		{
			name:   "Log w /var/www/log should be /var/www/log/",
			sample: Model{Secret: "lol", LogDir: "/var/www/log"},
			expect: "/var/www/log/",
		},
		{
			name:   "Log w/o value should be ./log/",
			sample: Model{Secret: "lol"},
			expect: "./log/",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.sample.SanitizationLog()
			assert.Equal(t, tt.expect, tt.sample.LogDir)
		})
	}
}
