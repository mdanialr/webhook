package handlers

import "github.com/gofiber/fiber/v2"

// DefaultRouteNotFound default handler to catch
// all not found error then return json message.
func DefaultRouteNotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "oops!! route not found",
	})
}
