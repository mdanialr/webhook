package helpers

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mdanialr/webhook/internal/config"
)

type pullRepoT struct {
	Name    string
	Service config.Service
}

// pullRepo pull from remote repo if repo name found in config
// file.
func pullRepo(repo pullRepoT) (string, error) {
	r, err := lookupRepo(repo.Name, repo.Service)
	if err != nil {
		return "", fmt.Errorf("lookup repo failed: %v\n", err)
	}

	cmd := parsePullCommand(r)
	res, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute git pull from remote repo: %v\n", err)
	}

	return fmt.Sprintf("\n%v", res), nil
}

// lookupRepo lookup for the repo that match given repo name then
// return it.
func lookupRepo(repo string, srv config.Service) (config.Repos, error) {
	for _, s := range srv {
		if s.Repos.Name == repo {
			return s.Repos, nil
		}
	}

	return config.Repos{}, fmt.Errorf("repo name not found in config file")
}

// parsePullCommand parse all necessary git command and append optional
// command then return it.
func parsePullCommand(repo config.Repos) string {
	cmdSeries := []string{
		"cd " + repo.RootPath,
		"git stash",
		"git pull",
		"git stash clear",
		repo.Cmd,
	}
	return strings.Join(cmdSeries, ";")
}
