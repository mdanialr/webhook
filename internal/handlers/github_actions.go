package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/github"
)

// GithubAction handler that would handle incoming POST request from GitHub action workflow
// that would trigger this webhook request.
func GithubAction(jobC chan string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqHook github.RequestPayload
		if err := c.BodyParser(&reqHook); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed parsing request body: %v", err),
			})
		}
		reqHook.CreateId()

		go func() {
			jobC <- reqHook.Id
		}()

		return c.SendStatus(200)
	}
}
