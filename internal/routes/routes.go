package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/handlers"
	"github.com/mdanialr/webhook/internal/middlewares"
)

func SetupRoutes(app *fiber.App) {
	// Built-in fiber middlewares
	app.Use(recover.New())
	// Use log file only in production
	var fConf logger.Config
	switch config.Conf.EnvIsProd {
	case true:
		fConf = logger.Config{
			Format:     "[${time}] ${status} | ${method} - ${latency} - ${ip} | ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			Output:     config.Conf.LogFile,
		}
	case false:
		fConf = logger.ConfigDefault
	}
	app.Use(logger.New(fConf))

	// This app's endpoints
	app.Get("/", handlers.Home)
	app.Post("/hook/:repo",
		middlewares.SecretToken,
		handlers.Hook,
	)

	// Custom middlewares AFTER endpoints
	app.Use(handlers.DefaultRouteFound)
}
