package handlers

import "github.com/gofiber/fiber/v2"

// DefaultError catch all thrown error by this app then
// return json message instead http.
func DefaultError(ctx *fiber.Ctx, err error) error {
	fiberr, ok := err.(*fiber.Error)
	if !ok {
		fiberr = fiber.ErrInternalServerError
	}

	data := struct {
		Message string        `json:"message"`
		Detail  []interface{} `json:"detail"`
		Status  int           `json:"status"`
	}{
		Message: fiberr.Message,
		Detail:  []interface{}{fiberr.Error()},
		Status:  fiberr.Code,
	}

	return ctx.Status(fiberr.Code).JSON(data)
}
