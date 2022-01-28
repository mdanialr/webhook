package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeConfwError struct{}

func (f *fakeConfwError) GetSHA256Signature(_ []byte) []byte { return nil }

func (f *fakeConfwError) ReloadConfig(_ io.Reader) error { return fmt.Errorf("error") }

func (f *fakeConfwError) SanitizationLog() {}

func (f *fakeConfwError) Sanitization() error { return nil }

type fakeLogger struct{}

func (f fakeLogger) Println(_ ...interface{}) {}

func TestReloadConfig(t *testing.T) {
	const tmpConfigPath = "/tmp/test-config-file.yaml"

	f, err := os.Create(tmpConfigPath)
	require.NoError(t, err)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln("failed closing tmp config file:", err)
		}
	}()

	someRandomFile := bytes.Buffer{}
	someRandomFile.WriteString(`some random file's content'`)
	_, err = someRandomFile.WriteTo(f)
	require.NoError(t, err, "failed to write content to file")

	t.Cleanup(func() {
		if err := os.Remove(tmpConfigPath); err != nil {
			log.Fatalln("failed cleaning tmp config file:", err)
		}
	})

	t.Run("1# every request should pass and continue to next handler which return 200", func(t *testing.T) {
		var fConf fakeConf
		configFilePath = tmpConfigPath

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
		configFilePath = "/tmp/no_exist.yaml"

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
