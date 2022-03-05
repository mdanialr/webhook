package worker

import (
	"os/exec"
	"testing"

	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/docker"
	"github.com/stretchr/testify/assert"
)

func TestDockerCDWorker(t *testing.T) {
	dockerSamples := docker.Docker{
		{docker.Model{
			Id:   "nzk_repo_latest",
			User: "nzk",
			Pass: "pass",
			Repo: "repo",
		}},
		{docker.Model{
			Id:   "user_repo-one_latest",
			User: "user",
			Pass: "sec",
			Repo: "repo-one",
		}},
		{docker.Model{
			Id:   "nzk_repo1_newest",
			User: "nzk",
			Pass: "pass-me",
			Repo: "repo1",
			Tag:  "newest",
		}},
	}

	conf := config.Model{Dockers: dockerSamples}

	testCases := []struct {
		name     string
		job      string
		errCount uint8
		wantErr  bool
	}{
		{
			name: "1# Docker that does exist should has 0 err count",
			job:  "user_repo-one_latest",
		},
		{
			name: "2# Docker that does exist should has 0 err count",
			job:  "nzk_repo1_newest",
		},
		{
			name:     "Repo that does not exist should has 1 err count",
			job:      "not-exist-repo",
			errCount: 1,
			wantErr:  true,
		},
	}

	// prepare the channels
	ch := &DockerChannel{
		JobC: make(chan string, 10),
		InfC: make(chan string, 10),
		ErrC: make(chan string, 10),
	}
	// spawn worker
	go DockerCDWorker(ch, &conf)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			execCmd = fakeExecCommand
			defer func() { execCmd = exec.Command }()

			var gotErr uint8
			ch.JobC <- tt.job

			if tt.wantErr {
				select {
				case <-ch.ErrC:
					gotErr++
				}
			}

			assert.Equal(t, tt.errCount, gotErr)
		})
	}
}
