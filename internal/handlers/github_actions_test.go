package handlers

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var payloadGithubAction = `
{
    "event": "push",
    "repository": "user/repo-y",
    "commit": "d09eda93a6ce94000f89254cb8e61363501d4117",
    "ref": "refs/heads/stable",
    "head": "",
    "workflow": "CI/CD",
    "requestID": "5a9ea3d9-cb99-4a71-b2ec-9b03606b1727"
}
`

func TestGithubAction(t *testing.T) {
	sample := model.Service{
		Github: []*model.Repo{
			{Event: "push", Branch: "stable", Name: "repo-y", User: "name"},
		},
	}

	testCases := []struct {
		name             string
		reqBody          io.Reader
		contentType      string
		expectStatusCode int
	}{
		{
			name:             "Should pass when sending correct json format and contain correct json structure",
			reqBody:          bytes.NewBufferString(payloadGithubAction),
			contentType:      fiber.MIMEApplicationJSON,
			expectStatusCode: fiber.StatusOK,
		},
	}

	const ROUTE = "/github/webhook"

	app := fiber.New()
	fakeChan := make(chan string)
	app.Post(ROUTE, GithubAction(fakeChan, &sample))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, ROUTE, tc.reqBody)
			req.Header.Set("content-type", tc.contentType)

			res, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.expectStatusCode, res.StatusCode)
		})
	}
}
