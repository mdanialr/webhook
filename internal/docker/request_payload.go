package docker

import (
	"fmt"
	"strings"
)

// RequestPayload most outer layer of requests payload from docker hub's webhook.
type RequestPayload struct {
	Id          string     // To identify and lookup in config which structured from RepoName and Tag (User_Repo_Tag).
	CallbackUrl string     `json:"callback_url"` // To send POST request back to docker hub to validate the webhook.
	PushData    PushData   `json:"push_data"`
	Repository  Repository `json:"repository"`
}

// PushData contain data regarding the push process that trigger this webhook.
type PushData struct {
	Pusher string `json:"pusher"` // The user that push to docker hub and trigger this webhook. Most likely the user that own this repo.
	Tag    string `json:"tag"`    // The tag of this repo.
}

// Repository contain data regarding this detail of the repo where this webhook triggered.
type Repository struct {
	RepoName string `json:"repo_name"` // Contain repository name with the owner in format USER/REPO.
}

// StdResponse standard structure to send back a response to docker hub's webhook.
type StdResponse struct {
	State   string `json:"state"`       // Required. Accepted values are success, failure, and error. If the state isn’t success, the Webhook chain is interrupted.
	Context string `json:"context"`     // A string containing the context of the operation. Can be retrieved from the Docker Hub. Maximum 100 characters.
	Desc    string `json:"description"` // A string containing miscellaneous information that is available on Docker Hub. Maximum 255 characters.
}

// CreateId create id that would be lookup in config. So the structure of the id should be
// the same as structuring id in `docker.Model`.
func (r *RequestPayload) CreateId() {
	r.Id = fmt.Sprintf(
		"%s_%s",
		strings.Replace(r.Repository.RepoName, "/", "_", 1),
		r.PushData.Tag,
	)
}
