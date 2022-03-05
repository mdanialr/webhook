package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleDockerStdResponses = []docker.StdResponse{
	{"ok", "ctx", "description"},
	{"error", "testing", "this is a testing"},
}

func TestInstance_DispatchPOST(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fiber.MIMEApplicationJSON, req.Header.Get("content-type"))
	}))

	t.Run("Should has correct content-type which is json", func(t *testing.T) {
		uri := fmt.Sprintf("%s/webhook", fakeServer.URL)
		api := Instance{Cl: fakeServer.Client(), Url: uri, Ctx: context.Background()}
		require.NoError(t, api.DispatchPOST(sampleDockerStdResponses[0]))
	})

	fakeServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("failed to read request body:", err)
		}

		var stdResponse docker.StdResponse
		err = json.Unmarshal(body, &stdResponse)
		if err != nil {
			log.Println("failed unmarshalling request:", err)
		}

		assert.Equal(t, sampleDockerStdResponses[0].State, stdResponse.State)
		assert.Equal(t, sampleDockerStdResponses[0].Context, stdResponse.Context)
		assert.Equal(t, sampleDockerStdResponses[0].Desc, stdResponse.Desc)

		rw.WriteHeader(fiber.StatusOK)
	}))

	t.Run("1# Should has correct json payload", func(t *testing.T) {
		uri := fmt.Sprintf("%s/webhook", fakeServer.URL)
		api := Instance{Cl: fakeServer.Client(), Url: uri, Ctx: context.Background()}
		require.NoError(t, api.DispatchPOST(sampleDockerStdResponses[0]))
	})

	fakeServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("failed to read request body:", err)
		}

		var stdResponse docker.StdResponse
		err = json.Unmarshal(body, &stdResponse)
		if err != nil {
			log.Println("failed unmarshalling request:", err)
		}

		assert.Equal(t, sampleDockerStdResponses[1].State, stdResponse.State)
		assert.Equal(t, sampleDockerStdResponses[1].Context, stdResponse.Context)
		assert.Equal(t, sampleDockerStdResponses[1].Desc, stdResponse.Desc)

		rw.WriteHeader(fiber.StatusOK)
	}))

	t.Run("2# Should has correct json payload", func(t *testing.T) {
		uri := fmt.Sprintf("%s/webhook", fakeServer.URL)
		api := Instance{Cl: fakeServer.Client(), Url: uri, Ctx: context.Background()}
		require.NoError(t, api.DispatchPOST(sampleDockerStdResponses[1]))
	})
}

func TestInstance_DispatchPOST_ErrorPaths(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))

	t.Run("Should error when using invalid host url", func(t *testing.T) {
		uri := fmt.Sprintf("%s/webhook", "http://localhost:9090")
		api := Instance{Cl: fakeServer.Client(), Url: uri, Ctx: context.Background()}
		require.Error(t, api.DispatchPOST(sampleDockerStdResponses[0]))
	})

	t.Run("Should error when injecting non struct type", func(t *testing.T) {
		uri := fmt.Sprintf("%s/webhook", fakeServer.URL)
		api := Instance{Cl: fakeServer.Client(), Url: uri, Ctx: context.Background()}
		require.Error(t, api.DispatchPOST(make(chan int)))
	})

	t.Run("Should error when using nil context", func(t *testing.T) {
		uri := fmt.Sprintf("%s/webhook", fakeServer.URL)
		api := Instance{Cl: fakeServer.Client(), Url: uri, Ctx: nil}
		require.Error(t, api.DispatchPOST(sampleDockerStdResponses[0]))
	})
}
