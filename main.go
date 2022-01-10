package main

import (
	"fmt"
	"github.com/mdanialr/webhook/internal/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/routes"
)

func main() {
	appConfig := config.Conf
	var proxyHeader string
	switch appConfig.EnvIsProd {
	case true:
		proxyHeader = "X-Real-Ip"
	case false:
		proxyHeader = ""
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: appConfig.EnvIsProd,
		ErrorHandler:          handlers.DefaultError,
		ProxyHeader:           proxyHeader,
	})

	routes.SetupRoutes(app)

	log.Printf("listening on %s:%v\n", appConfig.Host, appConfig.PortNum)
	log.Fatalln(app.Listen(fmt.Sprintf("%s:%v", appConfig.Host, appConfig.PortNum)))
}

func init() {
	// First thing first, init and load the config file
	if err := config.LoadConfigFromFile(); err != nil {
		log.Fatalln("failed to load config file:", err)
	}
}
