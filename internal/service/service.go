package service

import "fmt"

// Service holds list of data for all repos.
type Service []struct {
	Repo Model `yaml:"repo"`
}

// LookupRepo lookup for the repo in services that match given repo id.
func (s *Service) LookupRepo(id string) (Model, error) {
	for i, ss := range *s {
		if err := ss.Repo.Sanitization(); err != nil {
			return Model{}, fmt.Errorf("failed sanitizing #%d repo", i+1)
		}
		if ss.Repo.Id == id {
			return ss.Repo, nil
		}
	}

	return Model{}, fmt.Errorf("repo not found")
}

// Sanitization loop through all repos and run their each sanitization.
func (s *Service) Sanitization() error {
	for i, ss := range *s {
		if err := ss.Repo.Sanitization(); err != nil {
			return fmt.Errorf("failed sanitizing #%d repo in config: %s", i+1, err)
		}
	}

	return nil
}
