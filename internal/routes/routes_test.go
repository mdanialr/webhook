package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/worker"
)

type fakeLogger struct{}

func (f *fakeLogger) Println(_ ...interface{}) {}

func TestSetupRoutes(t *testing.T) {
	var fakeServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))
	bags := worker.BagOfChannels{
		GithubActionChan: &worker.Channel{JobC: make(chan string)},
	}

	t.Run("1# Success test", func(t *testing.T) {
		conf := config.Model{Secret: "1"}
		app := fiber.New()

		SetupRoutes(app, &conf, &fakeLogger{}, bags, make(chan string), make(chan string), fakeServer.Client())
	})

	t.Run("2# Success test", func(t *testing.T) {
		conf := config.Model{Secret: "1", EnvIsProd: true}
		app := fiber.New()

		SetupRoutes(app, &conf, &fakeLogger{}, bags, make(chan string), make(chan string), fakeServer.Client())
	})
}
