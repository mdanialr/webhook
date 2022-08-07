package app

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mdanialr/webhook/internal/handlers"
	"github.com/mdanialr/webhook/pkg/config"
)

// InitApp setup and initialize the fiber app.
func InitApp(conf *config.AppConfig, logFile *os.File) *fiber.App {
	isProd := strings.HasPrefix(conf.Config.GetString("env"), "prod")

	// if app in production use hostname from Nginx instead.
	var proxyHeader string
	if isProd {
		proxyHeader = "X-Real-Ip"
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: isProd,
		ErrorHandler:          handlers.DefaultError,
		ProxyHeader:           proxyHeader,
	})
	app.Use(recover.New())
	switch isProd {
	case true:
		fConf := logger.Config{
			Format:     "[${time}] ${status} | ${method} - ${latency} - ${ip} | ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			Output:     logFile,
		}
		app.Use(logger.New(fConf))
	case false:
		app.Use(logger.New())
	}

	return app
}
