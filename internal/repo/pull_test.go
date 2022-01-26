package repo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_ParsePullCommand(t *testing.T) {
	testCases := []struct {
		name     string
		sample   Model
		expected string
	}{
		{
			name: "Success parsing",
			sample: Model{
				RootPath: "/path/to/this/test-repo/",
				Cmd:      "",
			},
			expected: "cd /path/to/this/test-repo/;git stash;git pull;git stash clear;",
		},
		{
			name: "Another success parsing",
			sample: Model{
				RootPath: "/path/to/this/another-repo/",
				Cmd:      "systemctl reload nginx",
			},
			expected: "cd /path/to/this/another-repo/;git stash;git pull;git stash clear;systemctl reload nginx",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.sample.ParsePullCommand()
			assert.Equal(t, tt.expected, out)
		})
	}
}
