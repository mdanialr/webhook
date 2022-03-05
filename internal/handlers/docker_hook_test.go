package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleReqDockerWebhookModel = []docker.RequestPayload{
	{
		CallbackUrl: "",
		PushData:    docker.PushData{Pusher: "user5", Tag: "test-one"},
		Repository:  docker.Repository{RepoName: "user1/repo3"},
	},
	{
		CallbackUrl: "",
		PushData:    docker.PushData{Pusher: "user4", Tag: "test-one"},
		Repository:  docker.Repository{RepoName: "user1/repo6"},
	},
}

var fakeConfigDockerHubWebhook = `
dockers:
  - docker:
      user: user1
      pass: password
      repo: repo4
      tag: testing
      args: -p 5050:5000
  - docker:
      user: user1
      pass: password
      repo: repo7
      tag: newest
      args: -p 5600:4000
`

var fakeDockerHubServer = func(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		reqBody, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		var stdResponse docker.StdResponse
		require.NoError(t, json.Unmarshal(reqBody, &stdResponse))

		assert.Equal(t, "success", stdResponse.State, "`state` should contain 'success' when webhook payload received successfully")
		assert.Contains(t, stdResponse.Context, "Continuous Deployment", "`context` field should contain 'Continuous Deployment'")
		assert.Contains(t, stdResponse.Desc, "CD success", "`description` field should contain 'CD success'")

		rw.WriteHeader(fiber.StatusOK)
	}))
}

func TestDockerHubWebhook(t *testing.T) {
	buf := bytes.NewBufferString(fakeConfigDockerHubWebhook)
	sampleReqDockerWebhookModel[0].CallbackUrl = fakeDockerHubServer(t).URL
	sampleReqDockerWebhookModel[1].CallbackUrl = fakeDockerHubServer(t).URL

	conf, err := config.NewConfig(buf)
	require.NoError(t, err)
	require.NoError(t, conf.Dockers.Sanitization())

	const route = "/docker/hook"

	app := fiber.New()
	app.Post(route, DockerHubWebhook(fakeChan, fakeDockerHubServer(t).Client()))

	t.Run("Should pass when sending correct JSON format and payload to verify webhook chain", func(t *testing.T) {
		js, err := json.Marshal(sampleReqDockerWebhookModel[0])
		require.NoError(t, err)

		req := httptest.NewRequest(fiber.MethodPost, route, bytes.NewBuffer(js))
		req.Header.Set("content-type", fiber.MIMEApplicationJSON)
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

func TestDockerHubWebhook_ErrorPaths(t *testing.T) {
	buf := bytes.NewBufferString(fakeConfigDockerHubWebhook)

	conf, err := config.NewConfig(buf)
	require.NoError(t, err)
	require.NoError(t, conf.Dockers.Sanitization())

	const (
		route              = "/docker/hook"
		expectedStatusCode = fiber.StatusBadRequest
		expectedMIME       = fiber.MIMEApplicationJSON
	)

	app := fiber.New()
	app.Post(route, DockerHubWebhook(fakeChan, fakeDockerHubServer(t).Client()))

	t.Run("Should error and return 400 when sending format anything other than JSON", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, route, nil)
		res, _ := app.Test(req)
		assert.Equal(t, expectedStatusCode, res.StatusCode)
		assert.Equal(t, expectedMIME, res.Header.Get("content-type"))
	})

	t.Run("Error should return standard response and not empty", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, route, nil)
		res, _ := app.Test(req)
		assert.Equal(t, expectedStatusCode, res.StatusCode)
		assert.Equal(t, expectedMIME, res.Header.Get("content-type"))

		var js docker.StdResponse
		if err := json.NewDecoder(res.Body).Decode(&js); err != nil {
			t.Fatalf("failed to decode json in test: %v", err)
		}
		assert.NotEmpty(t, js.State, "`state` field in json response should no empty when there is error")
		assert.Equal(t, "error", js.State, "`state` field must contain 'error' when there is error")
		assert.NotEmpty(t, js.Context, "`context` field in json response should not empty")
		assert.NotEmpty(t, js.Desc, "`description` field in json response should not empty")
	})

	t.Run("Should error when using invalid callback_url", func(t *testing.T) {
		sampleReqDockerWebhookModel[1].CallbackUrl = "http://localhost:9000"

		jsPayload, err := json.Marshal(sampleReqDockerWebhookModel[1])
		require.NoError(t, err)

		req := httptest.NewRequest(fiber.MethodPost, route, bytes.NewBuffer(jsPayload))
		req.Header.Set("content-type", fiber.MIMEApplicationJSON)
		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadGateway, res.StatusCode)
		assert.Equal(t, expectedMIME, res.Header.Get("content-type"))

		var js docker.StdResponse
		if err := json.NewDecoder(res.Body).Decode(&js); err != nil {
			t.Fatalf("failed to decode json in test: %v", err)
		}
		assert.NotEmpty(t, js.State, "`state` field in json response should no empty when")
		assert.Equal(t, "error", js.State, "`state` field must contain 'error' when there is error")
		assert.NotEmpty(t, js.Context, "`context` field in json response should not empty")
		assert.NotEmpty(t, js.Desc, "`description` field in json response should not empty")
	})
}
