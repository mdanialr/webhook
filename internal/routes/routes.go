package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/handlers"
	"github.com/mdanialr/webhook/internal/middlewares"
	"github.com/mdanialr/webhook/internal/model"
	"github.com/mdanialr/webhook/internal/worker"
	"github.com/mdanialr/webhook/pkg/config"
)

func SetupRoutes(app *fiber.App, appConf *config.AppConfig, bag worker.BagOfChannels, svc *model.Service) {
	app.Get("/", handlers.Home)
	app.Post("/github/webhook",
		middlewares.Auth(appConf),
		handlers.GithubAction(bag.GithubActionChan.JobC, svc),
	)

	app.Use(handlers.DefaultRouteNotFound)
}
