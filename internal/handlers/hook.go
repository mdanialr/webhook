package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/helpers"
	"github.com/mdanialr/webhook/internal/models"
	"log"
	"strings"
)

func Hook(c *fiber.Ctx) error {
	// check if there is new config value by checking the new hash
	// file value against old hash.
	if err := helpers.ReloadConfigFile(); err != nil {
		log.Println("failed when try to reload and repopulate config file:", err)
	}

	repo := c.Params("repo")

	reqHook := new(models.Request)
	if err := c.BodyParser(reqHook); err != nil {
		return err
	}

	confMsgKey := config.Conf.Keyword
	confUsrKey := config.Conf.Usr

	remoteUsr := reqHook.Commits[0].Committer.Username
	message := reqHook.Commits[0].Message
	branch := strings.Split(reqHook.Branch, "/")
	branchName := strings.Join(branch[len(branch)-1:], "")

	var isReload bool
	// Validate config's username against incoming committer, only if config's username
	// is not 'empty' or using wildcard '*'.
	if confUsrKey != "empty" && confUsrKey != "*" && confUsrKey != remoteUsr {
		isReload = false
	}
	// Otherwise, no need to validate username.
	if confUsrKey == "empty" || confUsrKey == "*" {
		isReload = true
	}

	// Validate config's message against incoming committer message.
	switch confMsgKey {
	case "empty":
		isReload = true
	default:
		isReload = strings.HasPrefix(message, confMsgKey)
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
