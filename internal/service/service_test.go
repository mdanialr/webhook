package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_LookupRepo(t *testing.T) {
	serviceSample := Service{
		{Model{
			User: "user",
			Name: "repo-one",
			Path: "/path/to/repo-one/",
			Cmd:  "pwd",
		}},
		{Model{
			User: "user",
			Name: "repo-two",
			Path: "/path/to/repo-two/",
			Cmd:  "systemctl reload nginx",
		}},
		{Model{
			User: "user",
			Name: "repo-three",
			Path: "/path/to/repo-three/",
		}},
	}

	serviceSampleWithError := Service{
		{Model{User: "user", Name: "repo-one"}},
	}

	testCases := []struct {
		name     string
		sample   Service
		lookupId string
		expected Model
		wantErr  bool
	}{
		{
			name:     "1# Founded repo should be identical with expected result",
			sample:   serviceSample,
			lookupId: "user_repo-one_master",
			expected: Model{
				Id:     "user_repo-one_master",
				User:   "user",
				Name:   "repo-one",
				Branch: "master",
				Path:   "/path/to/repo-one/",
				Cmd:    "pwd",
			},
		},
		{
			name:     "2# Founded repo should be identical with expected result",
			sample:   serviceSample,
			lookupId: "user_repo-three_master",
			expected: Model{
				Id:     "user_repo-three_master",
				User:   "user",
				Name:   "repo-three",
				Branch: "master",
				Path:   "/path/to/repo-three/",
			},
		},
		{
			name:     "Should error when any required fields are not provided",
			sample:   serviceSampleWithError,
			lookupId: "user_repo_master",
			wantErr:  true,
		},
		{
			name:     "1# Should error when repo not found and the result should be empty Model",
			sample:   serviceSample,
			lookupId: "repo-x",
			wantErr:  true,
		},
		{
			name:     "2# Should error when repo not found and the result should be empty Model",
			sample:   serviceSample,
			lookupId: "repo-z",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.LookupRepo(tc.lookupId)

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expected, out)
			case true:
				require.Error(t, err)
			}

		})
	}
}

func TestService_Sanitization(t *testing.T) {
	serviceSample := []Service{
		{{Model{
			User: "user",
			Name: "repo",
			Path: "/path/to/repo",
		}}},
		{{Model{
			Name: "repo",
			Path: "/path/to/repo",
		}}},
		{{Model{
			User: "user",
			Path: "/path/to/repo",
		}}},
		{{Model{
			User: "user",
			Name: "repo",
		}}},
	}

	testCases := []struct {
		name    string
		sample  Service
		wantErr bool
	}{
		{
			name:   "Should pass when all required fields are provided",
			sample: serviceSample[0],
		},
		{
			name:    "Should fail when any required fields are not provided in this case is `user` field",
			sample:  serviceSample[1],
			wantErr: true,
		},
		{
			name:    "Should fail when any required fields are not provided in this case is `name` field",
			sample:  serviceSample[2],
			wantErr: true,
		},
		{
			name:    "Should fail when any required fields are not provided in this case is `path` field",
			sample:  serviceSample[3],
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
