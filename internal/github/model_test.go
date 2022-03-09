package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModel_Sanitization(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Model
		expect  Model
		wantErr bool
	}{
		{
			name:    "Should error when `name` field is not provided",
			sample:  Model{},
			wantErr: true,
		},
		{
			name:    "Should error when `path` field is not provided",
			sample:  Model{Name: "repo"},
			wantErr: true,
		},
		{
			name:   "Default value for `user` field if not provided is 'user'",
			sample: Model{Name: "repo", Path: "/fake/path"},
			expect: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "master", Id: "user_repo_master"},
		},
		{
			name:   "Default value for `branch` field if not provided is 'master'",
			sample: Model{Name: "repo", User: "user", Path: "/fake/path"},
			expect: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "master", Id: "user_repo_master"},
		},
		{
			name:   "Should use provided value for `branch` field if provided instead the default one",
			sample: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "main"},
			expect: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "main", Id: "user_repo_main"},
		},
		{
			name:   "Field `path` should has prefix slash '/'",
			sample: Model{Name: "repo", User: "user", Path: "fake/path", Branch: "main"},
			expect: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "main", Id: "user_repo_main"},
		},
		{
			name:   "Field `id` value and structure should be as expected which combination of User_Name_Branch",
			sample: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "main"},
			expect: Model{Name: "repo", User: "user", Path: "/fake/path", Branch: "main", Id: "user_repo_main"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, tc.sample)
			case true:
				require.Error(t, err)
			}
		})
	}
}

func TestModel_ParsePullCommand(t *testing.T) {
	testCases := []struct {
		name     string
		sample   Model
		expected string
	}{
		{
			name: "1# Should result as expected even `cmd` field is empty",
			sample: Model{
				Path: "/path/to/this/test-repo/",
				Cmd:  "",
			},
			expected: "cd /path/to/this/test-repo/ && git stash && git pull && git stash clear && ",
		},
		{
			name: "2# Should result as expected with `cmd` field provided",
			sample: Model{
				Path: "/path/to/this/another-repo/",
				Cmd:  "systemctl reload nginx",
			},
			expected: "cd /path/to/this/another-repo/ && git stash && git pull && git stash clear && systemctl reload nginx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.sample.ParsePullCommand()
			assert.Equal(t, tc.expected, out)
		})
	}
}
