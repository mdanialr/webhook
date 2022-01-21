package helpers

import (
	"testing"

	"github.com/mdanialr/webhook/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tConfSrv = config.Service{
	{Repos: config.Repos{Name: "TestRepo", Cmd: "pwd"}},
	{Repos: config.Repos{Name: "ATestRepo", Cmd: "systemctl"}},
	{Repos: config.Repos{Name: "BTestRepo", Cmd: ""}},
	{Repos: config.Repos{Name: "XTestRepo", Cmd: ""}},
}

func TestPullRepo(t *testing.T) {
	tests := []struct {
		name  string
		repo  pullRepoT
		want  string
		isErr bool
	}{
		{
			name: "Success PullRepo", repo: pullRepoT{
				Name:    "TestRepo",
				Service: tConfSrv,
			}, want: "", isErr: false,
		},
		{
			name: "Failed PullRepo", repo: pullRepoT{
				Name:    "NoExistTestRepo",
				Service: tConfSrv,
			}, want: "", isErr: true,
		},
		{
			name: "Also Failed PullRepo", repo: pullRepoT{
				Name:    "BTestRepo",
				Service: tConfSrv,
			}, want: "", isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pRepo := pullRepoT{
				Name:    tt.repo.Name,
				Service: tt.repo.Service,
			}

			_, err := pullRepo(pRepo)
			switch tt.isErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}

func TestParsePullCommand(t *testing.T) {
	tests := []struct {
		name  string
		repos config.Repos
		want  string
		isErr bool
	}{
		{
			name:  "Parse TestRepo",
			want:  "cd /path/to/this/repo/;git stash;git pull;git stash clear;pwd",
			repos: config.Repos{RootPath: "/path/to/this/repo/", Cmd: "pwd"},
			isErr: false,
		},
		{
			name:  "Parse AnotherTestRepo",
			want:  "cd /path/to/another/repo/;git stash;git pull;git stash clear;",
			repos: config.Repos{RootPath: "/path/to/another/repo/", Cmd: ""},
			isErr: false,
		},
		{
			name:  "Parse BadTestRepo",
			want:  "cd /path/to/bad/repo/;git stash;git pull;git stash clear;pwd",
			repos: config.Repos{RootPath: "/path/to/should/be/good/repo/", Cmd: ""},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := parsePullCommand(tt.repos)
			switch tt.isErr {
			case true:
				// This case should be error and result should not equal
				assert.NotEqual(t, tt.want, res)
			case false:
				// Should be no error and result should as expected
				assert.Equal(t, tt.want, res)
			}
		})
	}
}

func TestLookupRepo(t *testing.T) {
	tests := []struct {
		name  string
		repo  string
		srv   config.Service
		want  config.Repos
		isErr bool
	}{
		{
			name:  "Found TestRepo",
			repo:  "TestRepo",
			srv:   tConfSrv,
			want:  config.Repos{Name: "TestRepo", Cmd: "pwd"},
			isErr: false,
		},
		{
			name:  "Found BTestRepo",
			repo:  "BTestRepo",
			srv:   tConfSrv,
			want:  config.Repos{Name: "BTestRepo", Cmd: ""},
			isErr: false,
		},
		{
			name:  "NotFound ZTestRepo",
			repo:  "ZTestRepo",
			srv:   tConfSrv,
			want:  config.Repos{Name: "ZTestRepo", Cmd: "pwd", RootPath: "/path/to/Z/TestRepo/"},
			isErr: true,
		},
	}

	errMsg := "repo name not found in config file"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := lookupRepo(tt.repo, tt.srv)

			switch tt.isErr {
			case true:
				// Should be error not found and error message as expected also every value should not as expected
				require.Error(t, err)
				require.Equal(t, errMsg, err.Error())
				assert.NotEqual(t, tt.want.Name, res.Name)
				assert.NotEqual(t, tt.want.RootPath, res.RootPath)
				assert.NotEqual(t, tt.want.Cmd, res.Cmd)
			case false:
				// Should be no error and founded repos got returned
				require.NoError(t, err)
				assert.Equal(t, tt.want.Name, res.Name)
				assert.Equal(t, tt.want.RootPath, res.RootPath)
				assert.Equal(t, tt.want.Cmd, res.Cmd)
			}
		})
	}
}
