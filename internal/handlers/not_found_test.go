package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRouteNotFound(t *testing.T) {
	testCases := []struct {
		name         string
		route        string
		method       string
		expectedCode int
		expectedMIME string
	}{
		{
			name:         "GET / : Not found route should return 404 and JSON response",
			route:        "/",
			method:       fiber.MethodGet,
			expectedCode: 404,
			expectedMIME: fiber.MIMEApplicationJSON,
		},
	}

	app := fiber.New()
	app.Use(DefaultRouteNotFound)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.route, nil)

			res, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedMIME, res.Header.Get("Content-Type"))
		})
	}
}
