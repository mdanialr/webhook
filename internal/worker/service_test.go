package worker

import (
	"os/exec"
	"testing"

	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestJobCD(t *testing.T) {
	serviceSample := service.Service{
		{service.Model{
			Name: "repo-one",
			Path: "/path/to/repo-one/",
			Cmd:  "pwd",
		}},
		{service.Model{
			Name: "repo-two",
			Path: "/path/to/repo-two/",
			Cmd:  "systemctl reload nginx",
		}},
		{service.Model{
			Name: "repo-three",
			Path: "/path/to/repo-three/",
			Cmd:  "",
		}},
	}

	m := config.Model{
		Env:     "dev",
		PortNum: 5005,
		Secret:  "secret",
		LogDir:  "/home/nzk/test-app/webhook/log",
		Service: serviceSample,
	}

	testCases := []struct {
		name     string
		job      string
		errCount uint8
		isErr    bool
	}{
		{
			name: "Repo that exist should has 0 err count",
			job:  "repo-one",
		},
		{
			name:     "Repo that does not exist should has 1 err count",
			job:      "not-exist-repo",
			errCount: 1,
			isErr:    true,
		},
	}

	// prepare the channels
	bag := BagOfChannels{
		GithubWebhookChan: &Channel{JobC: make(chan string, 10), InfC: make(chan string, 10), ErrC: make(chan string, 10)},
	}
	// spawn worker
	go GithubCDWorker(bag, &m)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			execCmd = fakeExecCommand
			defer func() { execCmd = exec.Command }()

			var gotErr uint8
			bag.GithubWebhookChan.JobC <- tt.job

			if tt.isErr {
				select {
				case <-bag.GithubWebhookChan.ErrC:
					gotErr++
				}
			}

			assert.Equal(t, tt.errCount, gotErr)
		})
	}
}
