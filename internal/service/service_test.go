package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/mdanialr/webhook/internal/repo"
)

func TestModel_LookupRepo(t *testing.T) {
	serviceSample := Model{
		{repo.Model{
			Name:     "repo-one",
			RootPath: "/path/to/repo-one/",
			Cmd:      "pwd",
		}},
		{repo.Model{
			Name:     "repo-two",
			RootPath: "/path/to/repo-two/",
			Cmd:      "systemctl reload nginx",
		}},
		{repo.Model{
			Name:     "repo-three",
			RootPath: "/path/to/repo-three/",
			Cmd:      "",
		}},
	}

	testCases := []struct {
		name       string
		sample     Model
		lookupName string
		expected   repo.Model
		isErr      bool
	}{
		{
			name:       "1# Founded repo should be identical with expected result",
			sample:     serviceSample,
			lookupName: "repo-one",
			expected: repo.Model{
				Name:     "repo-one",
				RootPath: "/path/to/repo-one/",
				Cmd:      "pwd",
			},
			isErr: false,
		},
		{
			name:       "2# Founded repo should be identical with expected result",
			sample:     serviceSample,
			lookupName: "repo-two",
			expected: repo.Model{
				Name:     "repo-two",
				RootPath: "/path/to/repo-two/",
				Cmd:      "systemctl reload nginx",
			},
			isErr: false,
		},
		{
			name:       "1# Not found repo should return err and the err should be identical w ErrRepoNotFound",
			sample:     serviceSample,
			lookupName: "repo-x",
			expected: repo.Model{
				Name:     "repo-x",
				RootPath: "/path/to/repo-x/",
				Cmd:      "",
			},
			isErr: true,
		},
		{
			name:       "2# Not found repo should return err and the expected model w return model should no be identical",
			sample:     serviceSample,
			lookupName: "repo-z",
			expected:   repo.Model{Name: "repo-z"},
			isErr:      true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m, err := tt.sample.LookupRepo(tt.lookupName)
			switch tt.isErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tt.expected, m)
			case true:
				require.Error(t, err)
				assert.Equal(t, ErrRepoNotFound, err)
				assert.NotEqual(t, tt.expected, m)
			}
		})
	}
}
