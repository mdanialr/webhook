package docker

import (
	"fmt"
	"strings"
)

// Model holds data for each docker repo and their tags.
type Model struct {
	User  string `yaml:"user"` // User that would be authenticated to docker hub.
	Pass  string `yaml:"pass"` // Pass password or token for authenticating User to docker hub.
	Repo  string `yaml:"repo"` // Repo name for this User.
	Tag   string `yaml:"tag"`  // Tag that determine which tag that would be listened & pulled.
	Args  string `yaml:"args"` // Args is additional arguments that would be passed when running docker command's `docker run ...`.
	Id    string // Additional key to make it easier when lookup for a docker which structured from User, Repo and Tag (User_Repo_Tag).
	Image string // Image is needed when stopping and running docker after update the images from remote registry. Combined from User, Repo & Tag (User/Repo:Tag).
	Name  string // Name is needed when stopping and running docker after update the images from remote registry. Combined from Repo & Tag (Repo_Tag).
}

// Sanitization check and sanitize config Model's instance.
func (m *Model) Sanitization() error {
	if m.User == "" {
		return fmt.Errorf("`user` field is required")
	}

	if m.Pass == "" {
		return fmt.Errorf("`pass` field is required")
	}

	if m.Repo == "" {
		return fmt.Errorf("`repo` field is required")
	}

	if m.Tag == "" {
		m.Tag = "latest"
	}

	m.Id = fmt.Sprintf("%s_%s_%s", m.User, m.Repo, m.Tag)
	m.Image = fmt.Sprintf("%s/%s:%s", m.User, m.Repo, m.Tag)
	m.Name = fmt.Sprintf("%s_%s", m.Repo, m.Tag)

	return nil
}

// ParsePullCommand parse all necessary docker commands and
// append optional arguments then return it.
func (m *Model) ParsePullCommand() string {
	str := []string{
		fmt.Sprintf("docker login -u %s -p %s", m.User, m.Pass),
		fmt.Sprintf("docker pull -q %s", m.Image),
		fmt.Sprintf("docker stop %s", m.Name),
		"docker container prune -f",
		fmt.Sprintf("docker run --name %s -d %s %s", m.Name, m.Args, m.Image),
		`docker image prune --filter "dangling=true" -f`,
		"docker logout",
	}

	return strings.Join(str, " && ")
}
