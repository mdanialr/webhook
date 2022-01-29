package service

import (
	"fmt"

	"github.com/mdanialr/webhook/internal/repo"
)

// Model holds list of data for all repos.
type Model []struct {
	Repos repo.Model `yaml:"repo"`
}

// LookupRepo lookup for the repo in services that match given
// repo name. then return it if found and error would be nil otherwise
// error would be ErrRepoNotFound.
func (m *Model) LookupRepo(name string) (repo.Model, error) {
	for _, s := range *m {
		if s.Repos.Name == name {
			return s.Repos, nil
		}
	}

	return repo.Model{}, fmt.Errorf("repo not found")
}
