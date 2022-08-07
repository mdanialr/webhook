package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/model"
	"github.com/mdanialr/webhook/internal/worker"
	"github.com/mdanialr/webhook/pkg/config"
)

func TestSetupRoutes(t *testing.T) {
	bags := worker.BagOfChannels{
		GithubActionChan: &worker.Channel{JobC: make(chan string)},
	}

	t.Run("Should pass", func(t *testing.T) {
		app := fiber.New()

		SetupRoutes(app, &config.AppConfig{}, bags, &model.Service{})
	})
}
