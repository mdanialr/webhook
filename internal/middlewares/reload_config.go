package middlewares

import (
	"bytes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/logger"
)

var configFilePath = "config.yaml"

// ReloadConfig reload config instance at every call.
func ReloadConfig(conf config.Interface, l logger.Interface) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		f, err := os.ReadFile(configFilePath)
		if err != nil {
			l.Println("failed to open file while reload config:", err)
		}

		if err := conf.ReloadConfig(bytes.NewReader(f)); err != nil {
			l.Println("failed to reload config:", err)
		}

		return c.Next()
	}
}
