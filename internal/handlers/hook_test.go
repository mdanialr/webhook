package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var payloadGithubWebhook = `
{
  "ref": "refs/heads/main",
  "commits": [
    {
      "message": "test your ass",
      "committer": {
        "username": "user"
      }
    }
  ]
}
`

var fakeChan = make(chan string)

type responseJSON struct {
	Committer string `json:"committer"`
	Message   string `json:"message"`
	Branch    string `json:"branch"`
	Repo      string `json:"repo"`
}

func TestHook(t *testing.T) {
	const (
		ROUTE   = "/hook/"
		REQUEST = ROUTE + "repo"
	)
	var expectedResponseJson = responseJSON{
		Repo: "repo", Branch: "main",
		Message:   "test your ass",
		Committer: "user",
	}

	testCases := []struct {
		name             string
		reqBody          io.Reader
		contentType      string
		expectStatusCode int
		wantErr          bool
	}{
		{
			name:             "Should error and return `400` status code when sending request other than json format",
			reqBody:          nil,
			expectStatusCode: fiber.StatusBadRequest,
			wantErr:          true,
		},
		{
			name:             "Should error and return `400` status code when sending empty body request even its correct json content-type",
			reqBody:          nil,
			contentType:      fiber.MIMEApplicationJSON,
			expectStatusCode: fiber.StatusBadRequest,
			wantErr:          true,
		},
		{
			name:             "Should error and return `400` status code when sending correct body request but not json content-type",
			reqBody:          bytes.NewBufferString(payloadGithubWebhook),
			expectStatusCode: fiber.StatusBadRequest,
			wantErr:          true,
		},
		{
			name:             "Should pass when sending correct json format and contain correct json structure",
			reqBody:          bytes.NewBufferString(payloadGithubWebhook),
			contentType:      fiber.MIMEApplicationJSON,
			expectStatusCode: fiber.StatusOK,
		},
		{
			name:             "Should return correct json structure if test is pass",
			reqBody:          bytes.NewBufferString(payloadGithubWebhook),
			contentType:      fiber.MIMEApplicationJSON,
			expectStatusCode: fiber.StatusOK,
		},
	}

	app := fiber.New()
	app.Post(ROUTE+":repo", Hook(&config.Model{}, fakeChan))

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, REQUEST, tc.reqBody)
			req.Header.Set("content-type", tc.contentType)

			res, err := app.Test(req)
			require.NoError(t, err)

			switch tc.wantErr {
			case true:
				assert.Equal(t, tc.expectStatusCode, res.StatusCode)
			case false:
				assert.Equal(t, tc.expectStatusCode, res.StatusCode)
				var rJSON responseJSON
				if err := json.NewDecoder(res.Body).Decode(&rJSON); err != nil {
					t.Fatalf("failed to decode json in test #%v: %v", i+1, err)
				}
				assert.Equal(t, expectedResponseJson, rJSON)
			}
			assert.Equal(t, fiber.MIMEApplicationJSON, res.Header.Get("content-type"))
		})
	}
}
