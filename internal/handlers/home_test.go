package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	testCases := []struct {
		name         string
		route        string
		expectedCode int
	}{
		{
			name:         "GET / : Home handler should return 200",
			route:        "/",
			expectedCode: 200,
		},
	}

	app := fiber.New()
	app.Get("/", Home)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)

			res, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, res.StatusCode)
		})
	}
}
