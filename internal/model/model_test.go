package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo_ParseCommand(t *testing.T) {
	testCases := []struct {
		name   string
		sample *Repo
		expect string
	}{
		{
			name:   "Should return empty string if CMD is empty",
			sample: &Repo{},
		},
		{
			name:   "Should return the joined CMD fields using '&&'",
			sample: &Repo{Path: "/tmp", CMD: []string{"echo", "hello"}},
			expect: "cd /tmp && echo && hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, tc.sample.ParseCommand())
		})
	}
}

func TestService_LookupGithub(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Service
		id      string
		wantErr bool
	}{
		{
			name:    "Should error when the id is not found",
			sample:  Service{Github: []*Repo{{ID: "user_repo_tags"}}},
			id:      "user_repo_main",
			wantErr: true,
		},
		{
			name:   "Should return the repo when the id is found",
			sample: Service{Github: []*Repo{{ID: "user_repo_tags"}}},
			id:     "user_repo_tags",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.LookupGithub(tc.id)

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.id, out.ID)
			}
		})
	}
}

func TestService_ValidateAndParseID(t *testing.T) {
	okSample := Service{Github: []*Repo{{Name: "hello-world", User: "user", Branch: "main", Path: "/tmp", Tags: true}}}
	testCases := []struct {
		name      string
		sample    Service
		expectErr string
		wantErr   bool
	}{
		{
			name:      "Should error if name is not set and give error message as expected",
			sample:    Service{Github: []*Repo{{Event: "push"}}},
			expectErr: "`name` field is required",
			wantErr:   true,
		},
		{
			name:      "Should error if user is not set and give error message as expected",
			sample:    Service{Github: []*Repo{{Name: "hello-world"}}},
			expectErr: "`user` field is required",
			wantErr:   true,
		},
		{
			name:      "Should error if branch is not set and give error message as expected",
			sample:    Service{Github: []*Repo{{Name: "hello-world", User: "user"}}},
			expectErr: "`branch` field is required",
			wantErr:   true,
		},
		{
			name:      "Should error if path is not set and give error message as expected",
			sample:    Service{Github: []*Repo{{Name: "hello-world", User: "user", Branch: "main"}}},
			expectErr: "`path` field is required",
			wantErr:   true,
		},
		{
			name:   "Should pass and event is push if not set",
			sample: okSample,
		},
		{
			name:   "Should pass and id is using tags if its tags type",
			sample: okSample,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.ValidateAndParseID()

			switch tc.wantErr {
			case true:
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.expectErr)
			case false:
				require.NoError(t, err)
				assert.Equal(t, "push", tc.sample.Github[0].Event)
				assert.Equal(t, "user_hello-world_tags", tc.sample.Github[0].ID)
			}
		})
	}
}
