package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"strings"
)

// SecretToken middleware check incoming request signature, to make
// sure incoming request come from intended source then response
// with error if signature doesn't match.
func SecretToken(c *fiber.Ctx) error {
	if config.Conf.Secret != "" {
		reqSig, err := readReqSig(c.GetReqHeaders())
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"response": "failed to read request signature: " + err.Error(),
			})
		}
		confSig := getConfSignature(c.Request().Body())

		if !hmac.Equal(reqSig, confSig) {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"response": "signature doesn't match",
			})
		}
	}

	return c.Next()
}

// readReqSig read and decode incoming request signature
// in request headers and return its []byte.
func readReqSig(reqH map[string]string) ([]byte, error) {
	s := strings.TrimPrefix(reqH["X-Hub-Signature-256"], "sha256=")
	if s == "" {
		msg := "missing X-Hub-Signature-256 header in request"
		return nil, errors.New(msg)
	}

	sig, err := hex.DecodeString(s)
	if err != nil {
		msg := "failed hex-decoding signature " + s + ": " + err.Error()
		return nil, errors.New(msg)
	}

	return sig, nil
}

// getConfSignature get secret from config file.
func getConfSignature(body []byte) []byte {
	secret := []byte(config.Conf.Secret)
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)

	return mac.Sum(nil)
}
