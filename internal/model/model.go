package model

import (
	"fmt"
	"strings"
)

type Repo struct {
	// Name the name of the repo.
	Name string `mapstructure:"name"`
	// User owner of this repository.
	User string `mapstructure:"user"`
	// Branch the Branch which would be listened & pulled.
	Branch string ` mapstructure:"branch"`
	// Path a path where the local repo is located.
	Path string `mapstructure:"path"`
	// CMD additional bash commands that would be executed last.
	CMD []string `mapstructure:"commands"`
	// Tags whether you want to listen to tags or normal branch.
	Tags bool `mapstructure:"tags"`
	// Event what event you want to listen to.
	Event string `mapstructure:"event"`
	// ID additional key to make it easier when lookup for a repo which structured from User, Name and Branch (User_Name_Branch).
	ID string
}

// ParseCommand join CMD fields using '&&' and return it. Return empty string if CMD is empty.
func (m *Repo) ParseCommand() string {
	if len(m.CMD) > 0 {

		return fmt.Sprintf("cd %s && %s", m.Path, strings.Join(m.CMD, " && "))
	}

	return ""
}

type Service struct {
	Github []*Repo `mapstructure:"github"`
}

// LookupGithub find a GitHub repo by the given custom ID.
func (s *Service) LookupGithub(id string) (*Repo, error) {
	if len(s.Github) > 0 {
		for _, r := range s.Github {
			if r.ID == id {
				return r, nil
			}
		}
	}

	return nil, fmt.Errorf("repo not found with the given id: %s", id)
}

// ValidateAndParseID make sure required fields is set then parse unique ID which is the combination of
// User_Name_Branch for every Repo. Also set default value of 'event' field is `push`.
func (s *Service) ValidateAndParseID() error {
	if len(s.Github) > 0 {
		for _, r := range s.Github {
			if err := s.validate(r); err != nil {
				return err
			}

			if r.Event == "" {
				r.Event = "push"
			}

			if r.Tags { // if a tags type
				r.ID = fmt.Sprintf("%s_%s_%s", r.User, r.Name, "tags")
				continue
			}
			r.ID = fmt.Sprintf("%s_%s_%s", r.User, r.Name, r.Branch)
		}
	}

	return nil
}

// validate make sure required fields is set for the given Repo. Currently, required fields are: name, user, branch
// & path.
func (s *Service) validate(r *Repo) error {
	if r.Name == "" {
		return fmt.Errorf("`name` field is required")
	}
	if r.User == "" {
		return fmt.Errorf("`user` field is required")
	}
	if r.Branch == "" {
		return fmt.Errorf("`branch` field is required")
	}
	if r.Path == "" {
		return fmt.Errorf("`path` field is required")
	}

	return nil
}
