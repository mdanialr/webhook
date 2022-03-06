package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/handlers"
	nzLog "github.com/mdanialr/webhook/internal/logger"
	"github.com/mdanialr/webhook/internal/middlewares"
)

func SetupRoutes(app *fiber.App, conf *config.Model, l nzLog.Interface, jobC chan string, dC chan string, hCl *http.Client) {
	// Built-in fiber middlewares
	app.Use(recover.New())
	// Use log file only in production
	switch conf.EnvIsProd {
	case true:
		fConf := logger.Config{
			Format:     "[${time}] ${status} | ${method} - ${latency} - ${ip} | ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			Output:     conf.LogFile,
		}
		app.Use(logger.New(fConf))
	case false:
		app.Use(logger.New())
	}

	// This app's endpoints
	app.Get("/", handlers.Home)
	app.Post("/hook/:repo",
		middlewares.ReloadConfig(conf, l),
		middlewares.SecretToken(conf),
		handlers.Hook(conf, jobC),
	)
	app.Post("/docker/webhook",
		middlewares.ReloadConfig(conf, l),
		handlers.DockerHubWebhook(dC, hCl),
	)

	// Custom middlewares AFTER endpoints
	app.Use(handlers.DefaultRouteNotFound)
}
