package service

import (
	"fmt"
	"strings"
)

// Model holds data for each single repo.
type Model struct {
	Name   string `yaml:"name"`    // The name of the repository.
	User   string `yaml:"user"`    // The user that own this repository.
	Branch string `yaml:"branch"`  // The Branch which would be listened & pulled.
	Path   string `yaml:"path"`    // A path where the local repo is located.
	Cmd    string `yaml:"opt_cmd"` // Additional bash commands that would be executed last.
	Id     string // Additional key to make it easier when lookup for a repo which structured from User, Name and Branch (User_Name_Branch).
}

// Sanitization check and sanitize config Model's instance.
func (m *Model) Sanitization() error {
	if m.User == "" {
		return fmt.Errorf("`user` field is required")
	}
	if m.Name == "" {
		return fmt.Errorf("`name` is required")
	}
	if m.Path == "" {
		return fmt.Errorf("`path` field is required")
	}

	if !strings.HasPrefix(m.Path, "/") {
		m.Path = "/" + m.Path
	}
	if m.Branch == "" {
		m.Branch = "master"
	}

	m.Id = fmt.Sprintf(
		"%s_%s_%s",
		m.User,
		m.Name,
		m.Branch,
	)

	return nil
}

// ParsePullCommand parse all necessary git command and append optional command
// then return it.
func (m *Model) ParsePullCommand() string {
	str := []string{
		"cd " + m.Path,
		"git stash",
		"git pull",
		"git stash clear",
		m.Cmd,
	}

	return strings.Join(str, " && ")
}
