package repo

import "strings"

// Model holds data for each single repo.
type Model struct {
	Name     string `yaml:"name"`
	RootPath string `yaml:"root"`
	Cmd      string `yaml:"opt_cmd"`
}

// ParsePullCommand parse all necessary git command and
// append optional command then return it.
func (m *Model) ParsePullCommand() string {
	str := []string{
		"cd " + m.RootPath,
		"git stash",
		"git pull",
		"git stash clear",
		m.Cmd,
	}

	return strings.Join(str, ";")
}
