package helpers

import (
	"fmt"
	"github.com/mdanialr/webhook/internal/config"
	"os/exec"
	"strings"
)

// pullRepo pull from remote repo if repo name found in config
// file.
func pullRepo(repo string) {
	r, err := lookupRepo(repo, config.Conf.Service)
	if err != nil {
		NzLogErr.Println(err)
	}

	cmd := parsePullCommand(r)
	res, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		NzLogErr.Println("failed to execute git pull from remote repo:", err)
	}

	NzLogInf.Println("\n" + string(res))
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
