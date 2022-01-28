package middlewares

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type fakeConfwError struct{}

func (f *fakeConfwError) GetSHA256Signature(_ []byte) []byte { return nil }

func (f *fakeConfwError) ReloadConfig(_ io.Reader) error { return fmt.Errorf("error") }

func (f *fakeConfwError) SanitizationLog() {}

func (f *fakeConfwError) Sanitization() error { return nil }

type fakeLogger struct{}

func (f fakeLogger) Println(_ ...interface{}) {}

func TestReloadConfig(t *testing.T) {

	t.Run("1# every request should pass and continue to next handler which return 200", func(t *testing.T) {
		var fConf fakeConf
		configFilePath = "testdata/config.yaml"

		app := fiber.New()
		app.Post("/",
			ReloadConfig(&fConf, fakeLogger{}),
			func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			},
		)

		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("2# should log error but keep continue to next handler", func(t *testing.T) {
		var fConf fakeConfwError
		configFilePath = "testdata/no_exist.yaml"

		app := fiber.New()
		app.Post("/",
			ReloadConfig(&fConf, fakeLogger{}),
			func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			},
		)

		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		res, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}
