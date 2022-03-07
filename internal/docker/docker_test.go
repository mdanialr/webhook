package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocker_LookupRepoUsingModel(t *testing.T) {
	dockerSamples := Docker{
		{Model{
			Id:    "nzk_repo_latest",
			User:  "nzk",
			Pass:  "pass",
			Repo:  "repo",
			Tag:   "latest",
			Image: "nzk/repo:latest",
			Name:  "repo_latest",
		}},
		{Model{
			Id:    "user_repo-one_latest",
			User:  "user",
			Pass:  "sec",
			Repo:  "repo-one",
			Tag:   "latest",
			Image: "user/repo-one:latest",
			Name:  "repo-one_latest",
		}},
		{Model{
			Id:    "nzk_repo-me_latest",
			User:  "nzk",
			Pass:  "pass-me",
			Repo:  "repo-me",
			Tag:   "latest",
			Image: "nzk/repo-me:latest",
			Name:  "repo-me_latest",
		}},
	}

	testCases := []struct {
		name     string
		sample   Docker
		lookupId string
		expect   Model
		wantErr  bool
	}{
		{
			name:     "1# Founded docker should be identical with the expected result",
			sample:   dockerSamples,
			lookupId: "user_repo-one_latest",
			expect:   dockerSamples[1].Docker,
		},
		{
			name:     "2# Founded docker should be identical with the expected result",
			sample:   dockerSamples,
			lookupId: "nzk_repo_latest",
			expect:   dockerSamples[0].Docker,
		}, {
			name:     "1# Not found docker should return error and empty Model",
			sample:   dockerSamples,
			lookupId: "user_repo-latest_first",
			wantErr:  true,
		},
		{
			name:     "2# Not found docker should return error and empty Model",
			sample:   dockerSamples,
			lookupId: "user_repo-not-found_latest",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.LookupRepo(tc.lookupId)

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, out)
			case true:
				require.Error(t, err)
				assert.Equal(t, tc.expect, out)
			}
		})
	}
}

func TestDocker_SanitizationUsingModel(t *testing.T) {
	dockerSamples := []Docker{
		{{Model{Id: "nzk_repo_latest",
			User: "nzk",
			Pass: "pass",
			Repo: "repo"}},
		},
		{{Model{
			User: "",
			Pass: "sec",
			Repo: "repo-one",
		}}},
		{{Model{
			User: "nzk",
			Pass: "",
			Repo: "repo-me",
		}}},
	}

	testCases := []struct {
		name    string
		sample  Docker
		wantErr bool
	}{
		{
			name:   "Should pass when all required fields are provided",
			sample: dockerSamples[0],
		},
		{
			name:    "Should fail when any required fields are not provided in this case is `user` field",
			sample:  dockerSamples[1],
			wantErr: true,
		},
		{
			name:    "Should fail when any required fields are not provided in this case is `pass` field",
			sample:  dockerSamples[2],
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}
