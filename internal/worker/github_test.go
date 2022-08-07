package worker

import (
	"os/exec"
	"testing"

	"github.com/mdanialr/webhook/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGithubActionWebhookWorker(t *testing.T) {
	githubActionSamples := model.Service{
		Github: []*model.Repo{
			{Name: "repo-one", Path: "/path/to/repo-one/", CMD: []string{"pwd"}},
			{Name: "repo-two", Path: "/path/to/repo-two/", CMD: []string{"systemctl reload nginx"}},
			{Name: "repo-three", Path: "/path/to/repo-three/"},
		},
	}

	testCases := []struct {
		name     string
		job      string
		errCount uint8
		wantErr  bool
	}{
		{
			name: "Repo that exist should has 0 err count",
			job:  "user_repo-one_master",
		},
		{
			name:     "Using fake exec Command should error and has 1 err count",
			job:      "fake-command",
			errCount: 1,
			wantErr:  true,
		},
	}

	// prepare the channels
	bag := BagOfChannels{
		GithubActionChan: &Channel{JobC: make(chan string, 10), InfC: make(chan string, 10), ErrC: make(chan string, 10)},
	}
	// spawn worker
	go GithubActionWebhookWorker(bag, &githubActionSamples)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			execCmd = fakeExecCommand
			defer func() { execCmd = exec.Command }()

			var gotErr uint8
			bag.GithubActionChan.JobC <- tt.job

			if tt.wantErr {
				select {
				case <-bag.GithubActionChan.ErrC:
					gotErr++
				}
			}

			assert.Equal(t, tt.errCount, gotErr)
		})
	}
}
