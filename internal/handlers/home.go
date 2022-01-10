package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Home index endpoint handler.
func Home(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"time": time.Now().Format("2006-Jan-02T03:04:05"),
	})
}
