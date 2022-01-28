package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
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

func Hook(conf *config.Model) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		repo := c.Params("repo")

		var reqHook request
		if err := c.BodyParser(&reqHook); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed parsing request body: %v", err),
			})
		}

		confKeyword := conf.Keyword
		confUsr := conf.Usr

		remoteUsr := reqHook.Commits[0].Committer.Username
		message := reqHook.Commits[0].Message
		branch := strings.Split(reqHook.Branch, "/")
		branchName := strings.Join(branch[len(branch)-1:], "")

		// 1# Validate config's username against incoming committer.
		// 2# Validate config's keyword against incoming message's prefix.
		// 3# If both config's username or keyword is 'empty' no need to
		// validate username.

		var isReload bool
		switch {
		case confUsr == remoteUsr && confKeyword == "empty":
			isReload = true
		case confUsr == "empty" && strings.HasPrefix(message, confKeyword):
			isReload = true
		case confUsr == "empty" && confKeyword == "empty":
			isReload = true
		case confUsr == remoteUsr && strings.HasPrefix(message, confKeyword):
			isReload = true
		case confUsr != remoteUsr && confKeyword == "empty":
			isReload = false
		case confUsr == remoteUsr && !strings.HasPrefix(message, confKeyword):
			isReload = false
		case confUsr == "empty" && !strings.HasPrefix(message, confKeyword):
			isReload = false
		}

	if isReload {
		helpers.WorkerChan <- repo
	}

		return c.JSON(fiber.Map{
			"committer": remoteUsr,
			"message":   message,
			"branch":    branchName,
			"reload":    isReload,
			"repo":      repo,
		})
	}
}
