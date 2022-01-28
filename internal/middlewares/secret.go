package middlewares

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
)

// SecretToken middleware check incoming request signature, to make
// sure incoming request come from intended source then response
// with error if signature doesn't match.
func SecretToken(conf config.Interface) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s, err := searchSHA256Signature(c.GetReqHeaders())
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"response": fmt.Sprintf("failed read request signature: %v", err),
			})
		}

		rSig, err := hex.DecodeString(s)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"response": "failed hex-decoding signature",
			})
		}

		cSig := conf.GetSHA256Signature(c.Request().Body())

		if !hmac.Equal(rSig, cSig) {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"response": "signature doesn't match",
			})
		}

		return c.Next()
	}
}

// searchSHA256Signature search for sha256 signature in given map,
//return err if not found.
func searchSHA256Signature(m map[string]string) (string, error) {
	s := strings.TrimPrefix(m["X-Hub-Signature-256"], "sha256=")
	if s == "" {
		return "", errors.New("SHA256 signature not found")
	}

	return s, nil
}
