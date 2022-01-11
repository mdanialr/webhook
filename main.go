package main

import (
	"fmt"
	"github.com/mdanialr/webhook/internal/handlers"
	"github.com/mdanialr/webhook/internal/helpers"
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
	// Init internal logging
	if err := helpers.InitNzLog(); err != nil {
		log.Fatalln("failed to init internal logging:", err)
	}

	WorkerChan := make(chan string)
	// Assign the chan to global var in helpers, so it can be accessed by handlers
	// to send the job later
	helpers.WorkerChan = WorkerChan
	// Spawn workers as many as on max worker in config
	for w := 1; w <= config.Conf.MaxWorker; w++ {
		go helpers.CDWorker(WorkerChan)
	}
}
