package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultRouteNotFound(t *testing.T) {
	testCases := []struct {
		name         string
		route        string
		expectedCode int
		expectedMIME string
	}{
		{
			name:         "GET / : Not found route should return 404 and JSON response",
			route:        "/",
			expectedCode: 404,
			expectedMIME: fiber.MIMEApplicationJSON,
		},
	}

	app := fiber.New()
	app.Use(DefaultRouteNotFound)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)

			res, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedMIME, res.Header.Get("Content-Type"))
		})
	}
}
