package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupDefault(t *testing.T) {
	testCases := []struct {
		name    string
		sample  *viper.Viper
		expect  string
		wantErr bool
	}{
		{
			name:    "Should error without secret",
			sample:  viperWithoutSecret,
			expect:  "secret is required",
			wantErr: true,
		},
		{
			name:   "Should pass and other fields has default value as expected",
			sample: viperComplete,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetupDefault(tc.sample)

			switch tc.wantErr {
			case true:
				require.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			case false:
				require.NoError(t, err)
				assert.Equal(t, "dev", tc.sample.GetString("env"))
				assert.Equal(t, "127.0.0.1", tc.sample.GetString("host"))
				assert.Equal(t, 7575, tc.sample.GetInt("port"))
				assert.Equal(t, "/tmp", tc.sample.GetString("log"))
				assert.Equal(t, 1, tc.sample.GetInt("max_worker"))
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	testCases := []struct {
		name    string
		sample  string
		wantErr bool
	}{
		{
			name:    "Should error when config file not found",
			sample:  "/fake/path",
			wantErr: true,
		},
		{
			name:   "Should pass when config file found",
			sample: "/tmp/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := InitConfig(tc.sample)

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}
