package github

import (
	"fmt"
	"strings"
)

// RequestPayload payload of the webhook that would be sent by GitHub actions.
type RequestPayload struct {
	Id       string // To identify and lookup in config which structured from User, RepoName and Branch (User_RepoName_Branch).
	User     string // The user which own this repo that could be extracted from Repo.
	RepoName string // The name of the repository that could be extracted from Repo.
	Branch   string // Branch of this repo which extracted from Ref.
	Repo     string `json:"repository"` // Contain owner & repo name where this webhook triggered.
	Ref      string `json:"ref"`        // Contain reference which branch is this webhook triggered.
}

// CreateId create id that would be lookup in config. So the structure of the id should be
// the same as structuring id in `service.Model`.
func (r *RequestPayload) CreateId() {
	splitRepo := strings.Split(r.Repo, "/")
	splitBr := strings.Split(r.Ref, "/")

	r.User = splitRepo[0]
	r.RepoName = splitRepo[len(splitRepo)-1]
	r.Branch = splitBr[len(splitBr)-1]

	if splitBr[1] == "tags" {
		r.Branch = splitBr[1]
	}

	r.Id = fmt.Sprintf(
		"%s_%s_%s",
		r.User, r.RepoName, r.Branch,
	)
}
