package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultError(t *testing.T) {
	testCases := []struct {
		name         string
		route        string
		method       string
		expectedCode int
		expectedMIME string
	}{
		{
			name:         "GET /err : internal server error should return 500 in JSON response",
			route:        "/err",
			method:       fiber.MethodGet,
			expectedCode: 500,
			expectedMIME: fiber.MIMEApplicationJSON,
		},
		{
			name:         "GET /error : internal server error should return 500 in JSON response",
			route:        "/error",
			method:       fiber.MethodGet,
			expectedCode: 500,
			expectedMIME: fiber.MIMEApplicationJSON,
		},
	}

	app := fiber.New(fiber.Config{ErrorHandler: DefaultError})
	app.Get("/err", func(c *fiber.Ctx) error {
		// This should trigger default error handler
		return fiber.ErrInternalServerError
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return errors.New("this should trigger type assertion error")
	})

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.route, nil)

			res, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedMIME, res.Header.Get("Content-Type"))
		})
	}
}
