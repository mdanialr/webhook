package model

import (
	"fmt"
	"strings"
)

// RequestPayload payload of the webhook that would be sent by GitHub actions.
type RequestPayload struct {
	// ID additional key to make it easier when lookup for a repo which structured from User, Name and Branch as (User_Name_Branch).
	ID       string
	User     string // The user that own this repo that could be extracted from Repo.
	RepoName string // The name of the repository that could be extracted from Repo.
	Branch   string // Branch of this repo which extracted from Ref.
	Tags     bool   // whether the given request payload is normal branch or a tags.
	Repo     string `json:"repository"` // Contain owner & repo name where this webhook triggered.
	Ref      string `json:"ref"`        // Contain reference which branch is this webhook triggered.
	Event    string `json:"event"`      // Contain what type of event that triggered this webhook.
}

// CreateID create ID.
func (r *RequestPayload) CreateID() {
	splitRepo := strings.Split(r.Repo, "/")
	splitRef := strings.Split(r.Ref, "/")

	if len(splitRepo) > 0 {
		r.User = splitRepo[0]
		r.RepoName = splitRepo[len(splitRepo)-1]
	}
	if len(splitRef) > 0 {
		r.Branch = splitRef[len(splitRef)-1]
	}

	if strings.Contains(r.Ref, "refs/tags/") {
		r.ID = fmt.Sprintf("%s_%s_%s", r.User, r.RepoName, "tags")
		r.Tags = true
		return
	}

	r.ID = fmt.Sprintf("%s_%s_%s", r.User, r.RepoName, r.Branch)
}
