package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/handlers"
	"github.com/mdanialr/webhook/internal/logger"
	"github.com/mdanialr/webhook/internal/routes"
	"github.com/mdanialr/webhook/internal/worker"
)

func main() {
	f, err := os.ReadFile("app-config.yaml")
	if err != nil {
		log.Fatalln("failed to read config file:", err)
	}

	var appConfig config.Model
	app, err := setup(&appConfig, bytes.NewReader(f))
	if err != nil {
		log.Fatalln("failed setup the app:", err)
	}

	// prepare common worker channel to store in the bag that contain all different worker channels
	bagOfChannels := worker.BagOfChannels{
		GithubActionChan:  &worker.Channel{JobC: make(chan string, 10), InfC: make(chan string, 10), ErrC: make(chan string, 10)},
		GithubWebhookChan: &worker.Channel{JobC: make(chan string, 10), InfC: make(chan string, 10), ErrC: make(chan string, 10)},
		DockerWebhookChan: &worker.Channel{JobC: make(chan string, 10), InfC: make(chan string, 10), ErrC: make(chan string, 10)},
	}
	// spawn worker pool with max number based on config's max worker
	for w := 1; w <= appConfig.MaxWorker; w++ {
		go worker.GithubCDWorker(bagOfChannels, &appConfig)
		go worker.DockerCDWorker(bagOfChannels, &appConfig)
		go worker.GithubActionWebhookWorker(bagOfChannels, &appConfig)
	}
	// spawn worker to write internal logger from Hook Handler
	go logWriterFromChannel(bagOfChannels)

	// init custom app logger
	appConfig.LogFile, err = os.OpenFile(appConfig.LogDir+"log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		log.Fatalln("failed to open|create log file:", err)
	}

	cl := &http.Client{}
	routes.SetupRoutes(app, &appConfig, logger.InfL, bagOfChannels, cl)

	logger.InfL.Printf("listening on %s:%v\n", appConfig.Host, appConfig.PortNum)
	logger.ErrL.Fatalln(app.Listen(fmt.Sprintf("%s:%v", appConfig.Host, appConfig.PortNum)))
}

// setup prepare everything that necessary before starting this app.
func setup(conf *config.Model, fBuf io.Reader) (*fiber.App, error) {
	// init and load the config file.
	newConf, err := config.NewConfig(fBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v\n", err)
	}
	*conf = *newConf
	if err = conf.Sanitization(); err != nil {
		return nil, fmt.Errorf("failed sanitizing config: %v\n", err)
	}
	conf.SanitizationLog()
	if err = conf.Dockers.Sanitization(); err != nil {
		return nil, fmt.Errorf("failed sanitizing docker config: %s\n", err)
	}

	// Init internal logging.
	if err := logger.InitLogger(conf); err != nil {
		return nil, fmt.Errorf("failed to init internal logging: %v\n", err)
	}

	// if app in production use hostname from Nginx instead.
	var proxyHeader string
	if conf.EnvIsProd {
		proxyHeader = "X-Real-Ip"
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: conf.EnvIsProd,
		ErrorHandler:          handlers.DefaultError,
		ProxyHeader:           proxyHeader,
	})

	return app, nil
}

// logWriterFromChannel listen to channels and write every message
// to internal logger.
func logWriterFromChannel(bag worker.BagOfChannels) {
	go func() {
		for inf := range bag.GithubActionChan.InfC {
			logger.InfL.Printf(inf)
		}
	}()
	go func() {
		for err := range bag.GithubActionChan.ErrC {
			logger.ErrL.Printf(err)
		}
	}()

	go func() {
		for inf := range bag.GithubWebhookChan.InfC {
			logger.InfL.Printf(inf)
		}
	}()
	go func() {
		for err := range bag.GithubWebhookChan.ErrC {
			logger.ErrL.Printf(err)
		}
	}()

	go func() {
		for inf := range bag.DockerWebhookChan.InfC {
			logger.InfL.Printf(inf)
		}
	}()
	go func() {
		for err := range bag.DockerWebhookChan.ErrC {
			logger.ErrL.Printf(err)
		}
	}()
}
