package middlewares

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchSHA256Signature(t *testing.T) {
	testCases := []struct {
		name   string
		sample map[string]string
		expect string
		isErr  bool
	}{
		{
			name:   "Signature found should return no error",
			sample: map[string]string{"X-Hub-Signature-256": "exists"},
			expect: "exists",
			isErr:  false,
		},
		{
			name:   "Signature not found should return error",
			sample: map[string]string{},
			expect: "",
			isErr:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := searchSHA256Signature(tt.sample)
			switch tt.isErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
				assert.Equal(t, tt.expect, res)
			}
		})
	}
}
