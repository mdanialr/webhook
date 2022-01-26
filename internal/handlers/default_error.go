package handlers

import "github.com/gofiber/fiber/v2"

// DefaultError catch all thrown error by this app then
// return json message instead http.
func DefaultError(ctx *fiber.Ctx, err error) error {
	fibErr, ok := err.(*fiber.Error)
	if !ok {
		fibErr = fiber.ErrInternalServerError
	}

	data := struct {
		Message string        `json:"message"`
		Detail  []interface{} `json:"detail"`
		Status  int           `json:"status"`
	}{
		Message: fibErr.Message,
		Detail:  []interface{}{fibErr.Error()},
		Status:  fibErr.Code,
	}

	return ctx.Status(fibErr.Code).JSON(data)
}
