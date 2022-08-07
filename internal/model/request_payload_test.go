package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestPayload_CreateId(t *testing.T) {
	testCases := []struct {
		name   string
		sample RequestPayload
		expect RequestPayload
	}{
		{
			name:   "1# Result should be as expected",
			sample: RequestPayload{Repo: "user/repo", Ref: "refs/heads/main"},
			expect: RequestPayload{Repo: "user/repo", Ref: "refs/heads/main", User: "user", RepoName: "repo", Branch: "main", ID: "user_repo_main"},
		},
		{
			name:   "2# Result should be as expected",
			sample: RequestPayload{Repo: "us/repo-y", Ref: "refs/heads/stable"},
			expect: RequestPayload{Repo: "us/repo-y", Ref: "refs/heads/stable", User: "us", RepoName: "repo-y", Branch: "stable", ID: "us_repo-y_stable"},
		},
		{
			name:   "1# Should use branch `tags` instead if the payload is a tag",
			sample: RequestPayload{Repo: "user/repo", Ref: "refs/tags/v1.0"},
			expect: RequestPayload{Repo: "user/repo", Ref: "refs/tags/v1.0", User: "user", RepoName: "repo", Branch: "v1.0", ID: "user_repo_tags", Tags: true},
		},
		{
			name:   "2# Should use the branch `tags` instead if the payload is a tag",
			sample: RequestPayload{Repo: "user/repo", Ref: "refs/tags/v2.4"},
			expect: RequestPayload{Repo: "user/repo", Ref: "refs/tags/v2.4", User: "user", RepoName: "repo", Branch: "v2.4", ID: "user_repo_tags", Tags: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.CreateID()
			assert.Equal(t, tc.expect, tc.sample)
		})
	}
}
