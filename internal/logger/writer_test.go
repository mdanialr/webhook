package logger

import (
	"testing"

	"github.com/mdanialr/webhook/internal/config"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	testCases := []struct {
		name       string
		sampleConf config.Model
		isErr      bool
	}{
		{
			name: "Success w no error",
			sampleConf: config.Model{
				LogDir: "/tmp/",
			},
			isErr: false,
		},
		{name: "Failed w error", isErr: true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.sampleConf.SanitizationLog()

			switch tt.isErr {
			case false:
				require.NoError(t, InitLogger(&tt.sampleConf))
			case true:
				require.Error(t, InitLogger(&tt.sampleConf))
			}
		})
	}
}
