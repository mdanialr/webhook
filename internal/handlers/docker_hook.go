package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/docker"
	"github.com/mdanialr/webhook/internal/docker/client"
)

// DockerHubWebhook handle incoming POST request from docker hub's webhook then send job to workers
// if there is match in config then send back `success` status to the CallbackUrl to verify that
// the webhook is complete and success.
func DockerHubWebhook(jobC chan string, hCl *http.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqHook docker.RequestPayload
		if err := c.BodyParser(&reqHook); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(docker.StdResponse{
				State:   "error",
				Context: "ERROR occurred before sending job to workers",
				Desc:    fmt.Sprintf("failed parsing request body: %v", err),
			})
		}
		reqHook.CreateId()

		go func() {
			jobC <- reqHook.Id
		}()

		// prepare payload to send `success` request and verify this webhook chain process.
		reqPayload := docker.StdResponse{
			State:   "success",
			Context: fmt.Sprintf("Continuous Deployment for %s", reqHook.Id),
			Desc:    "CD successfully triggered by docker hub's webhook",
		}

		// setup and prepare http client
		cl := client.Instance{Cl: hCl, Url: reqHook.CallbackUrl, Ctx: context.Background()}
		if err := cl.DispatchPOST(reqPayload); err != nil {
			c.Status(fiber.StatusBadGateway)
			return c.JSON(docker.StdResponse{
				State:   "error",
				Context: "ERROR occurred after sending job to workers",
				Desc:    fmt.Sprintf("failed when dispatching http client to send POST request to callback_url: %v", err),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
