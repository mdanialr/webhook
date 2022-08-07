package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/model"
)

// GithubAction handler that would handle incoming POST request from GitHub action workflow
// that would trigger this webhook request.
func GithubAction(jobC chan<- string, svc *model.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var reqHook model.RequestPayload
		c.BodyParser(&reqHook)
		reqHook.CreateID()

		if repo, _ := svc.LookupGithub(reqHook.ID); repo != nil {
			if repo.Event != reqHook.Event {
				return c.SendStatus(200) // immediately return if event is not match.
			}
		}

		go func() {
			jobC <- reqHook.ID
		}()

		return c.SendStatus(200)
	}
}
