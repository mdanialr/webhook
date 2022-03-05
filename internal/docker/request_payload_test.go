package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestPayload_CreateId(t *testing.T) {
	testCases := []struct {
		name   string
		sample RequestPayload
		expect string
	}{
		{
			name: "1# `Id` field with the given data should be as expected",
			sample: RequestPayload{
				PushData:   PushData{Tag: "new"},
				Repository: Repository{RepoName: "user1/repo3"},
			},
			expect: "user1_repo3_new",
		},
		{
			name: "2# `Id` field with the given data should be as expected",
			sample: RequestPayload{
				PushData:   PushData{Tag: "latest"},
				Repository: Repository{RepoName: "user1/repo7"},
			},
			expect: "user1_repo7_latest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.CreateId()
			assert.Equal(t, tc.expect, tc.sample.Id)
		})
	}
}
