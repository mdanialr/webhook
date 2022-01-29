package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"testing"
)

type fakeLogger struct{}

func (f *fakeLogger) Println(_ ...interface{}) {}

func TestSetupRoutes(t *testing.T) {

	t.Run("1# Success test", func(t *testing.T) {
		conf := config.Model{Secret: "1"}
		app := fiber.New()

		SetupRoutes(app, &conf, &fakeLogger{}, make(chan string))
	})

	t.Run("2# Success test", func(t *testing.T) {
		conf := config.Model{Secret: "1", EnvIsProd: true}
		app := fiber.New()

		SetupRoutes(app, &conf, &fakeLogger{}, make(chan string))
	})
}
