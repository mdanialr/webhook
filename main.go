package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdanialr/webhook/internal/app"
	"github.com/mdanialr/webhook/internal/model"
	"github.com/mdanialr/webhook/internal/routes"
	"github.com/mdanialr/webhook/internal/worker"
	"github.com/mdanialr/webhook/pkg/config"
	"github.com/mdanialr/webhook/pkg/logger"
)

func main() {
	conf, err := config.InitConfig(".")
	if err != nil {
		log.Fatalln("failed to init config file:", err)
	}
	if err = config.SetupDefault(conf); err != nil {
		log.Fatalln("failed to sanitize and setup default config:", err)
	}

	var svc model.Service
	if err = conf.Unmarshal(&svc); err != nil {
		log.Fatalln("failed to unmarshal config to service model:", err)
	}
	// run parsing ID and validate it to make sure required fields are set.
	if err = svc.ValidateAndParseID(); err != nil {
		log.Fatalln("failed to validate and parse github config:", err)
	}

	infoLogger, err := logger.InitInfoLogger(conf)
	if err != nil {
		log.Fatalln("failed to init info logger:", err)
	}
	errLogger, err := logger.InitErrorLogger(conf)
	if err != nil {
		log.Fatalln("failed to init error logger:", err)
	}
	appConfig := &config.AppConfig{
		Config: conf,
		InfL:   infoLogger,
		ErrL:   errLogger,
	}

	// prepare common worker channel to store in the bag that contain all different worker channels
	bagOfChannels := worker.BagOfChannels{
		GithubActionChan: &worker.Channel{JobC: make(chan string, 10), InfC: make(chan string, 10), ErrC: make(chan string, 10)},
	}
	// spawn worker pool with max number based on config's max worker
	for w := 1; w <= conf.GetInt("max_worker"); w++ {
		go worker.GithubActionWebhookWorker(bagOfChannels, &svc)
	}
	// spawn worker to write internal logger from Hook Handler
	go logWriterFromChannel(bagOfChannels, appConfig)

	// init custom fiber app logger
	appLog := strings.TrimSuffix(conf.GetString("log"), "/")
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s", appLog, "log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		log.Fatalln("failed to open|create log file:", err)
	}

	fiberApp := app.InitApp(appConfig, logFile)
	routes.SetupRoutes(fiberApp, appConfig, bagOfChannels, &svc)

	host := conf.GetString("host")
	port := conf.GetInt("port")
	appConfig.InfL.Printf("listening on %s:%d\n", host, port)
	appConfig.ErrL.Fatalln(fiberApp.Listen(fmt.Sprintf("%s:%d", host, port)))
}

// logWriterFromChannel listen to channels and write every message
// to internal logger.
func logWriterFromChannel(bag worker.BagOfChannels, conf *config.AppConfig) {
	go func() {
		for inf := range bag.GithubActionChan.InfC {
			conf.InfL.Printf(inf)
		}
	}()
	go func() {
		for err := range bag.GithubActionChan.ErrC {
			conf.ErrL.Printf(err)
		}
	}()
}
