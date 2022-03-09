package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// request hold the most outer scope of incoming JSON from GitHub Webhook
type request struct {
	Commits []commit `json:"commits"`
	Branch  string   `json:"ref"`
}

// commit hold message that identified whether it contains the necessary char or not
type commit struct {
	Message   string    `json:"message"`
	Committer committer `json:"committer"`
}

// committer hold who did the commit
type committer struct {
	Username string `json:"username"`
}

func Hook(jobC chan string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		repo := c.Params("repo")

		var reqHook request
		if err := c.BodyParser(&reqHook); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed parsing request body: %v", err),
			})
		}

		remoteUsr := reqHook.Commits[0].Committer.Username
		message := reqHook.Commits[0].Message
		branch := strings.Split(reqHook.Branch, "/")
		branchName := strings.Join(branch[len(branch)-1:], "")

		go func() {
			jobC <- repo
		}()

		return c.JSON(fiber.Map{
			"committer": remoteUsr,
			"message":   message,
			"branch":    branchName,
			"repo":      repo,
		})
	}
}
