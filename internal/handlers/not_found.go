package handlers

import "github.com/gofiber/fiber/v2"

// DefaultRouteFound default handler to catch
// all not found error then return json message.
func DefaultRouteFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "oops!! route not found",
	})
}
