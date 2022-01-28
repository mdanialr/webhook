package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
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

type fakeConf struct {
	Secret string
}

func (f *fakeConf) GetSHA256Signature(in []byte) []byte {
	secret := []byte(f.Secret)
	mac := hmac.New(sha256.New, secret)
	mac.Write(in)

	return mac.Sum(nil)
}

func (f *fakeConf) ReloadConfig(_ io.Reader) error { return nil }

func (f *fakeConf) SanitizationLog() {}

func (f *fakeConf) Sanitization() error { return nil }

func TestSecretToken(t *testing.T) {
	type resJSON struct {
		Response string `json:"response"`
	}
	hMacSample := "b94c8a8c984cc521020717c5203fe4cb9fa83d8b0815fc0de8ae99cd9a0a914b" // secret123

	var fConf fakeConf

	app := fiber.New()
	app.Post("/",
		SecretToken(&fConf),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)

	t.Run("1# request that has no any request body or header should not pass and return 400", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var rJSON resJSON
		if err := json.NewDecoder(res.Body).Decode(&rJSON); err != nil {
			t.Fatalf("failed to decode json in test #1: %v", err)
		}
		fmt.Printf("1# response: %v\n", rJSON.Response)
	})

	t.Run("#2 request that has SHA256 signature but not valid hmac should not pass and return 400", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		req.Header.Set("X-Hub-Signature-256", "sha256=secret123")
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var rJSON resJSON
		if err := json.NewDecoder(res.Body).Decode(&rJSON); err != nil {
			t.Fatalf("failed to decode json in test #2: %v", err)
		}
		fmt.Printf("2# response: %v\n", rJSON.Response)
	})

	t.Run("#3 request that has SHA256 signature but doesn't match with the config should not pass and return 400", func(t *testing.T) {
		fConf.Secret = "secret321"

		rBody := strings.NewReader(`{"key": "value"}`)
		req := httptest.NewRequest(fiber.MethodPost, "/", rBody)
		req.Header.Set("X-Hub-Signature-256", fmt.Sprintf("sha256=%s", hMacSample))
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		var rJSON resJSON
		if err := json.NewDecoder(res.Body).Decode(&rJSON); err != nil {
			t.Fatalf("failed to decode json in test #3: %v", err)
		}
		fmt.Printf("3# response: %v\n", rJSON.Response)
	})

	t.Run("#4 the only request that should pass must has SHA256 signature w valid hmac and match with the config", func(t *testing.T) {
		fConf.Secret = "secret123"

		rBody := strings.NewReader(`{"key": "value"}`)
		req := httptest.NewRequest(fiber.MethodPost, "/", rBody)
		req.Header.Set("X-Hub-Signature-256", fmt.Sprintf("sha256=%s", hMacSample))
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}
