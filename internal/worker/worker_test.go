package worker

import (
	"os"
	"os/exec"
	"testing"

	"github.com/mdanialr/webhook/internal/config"
	"github.com/mdanialr/webhook/internal/repo"
	"github.com/mdanialr/webhook/internal/service"
	"github.com/stretchr/testify/assert"
)

type fakeWriter struct{}

func (_ fakeWriter) Write(_ []byte) (_ int, _ error) { return }

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Stdout = fakeWriter{}
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestJobCD(t *testing.T) {
	serviceSample := service.Model{
		{repo.Model{
			Name:     "repo-one",
			RootPath: "/path/to/repo-one/",
			Cmd:      "pwd",
		}},
		{repo.Model{
			Name:     "repo-two",
			RootPath: "/path/to/repo-two/",
			Cmd:      "systemctl reload nginx",
		}},
		{repo.Model{
			Name:     "repo-three",
			RootPath: "/path/to/repo-three/",
			Cmd:      "",
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
	ch := &Channel{
		JobC: make(chan string, 10),
		InfC: make(chan string, 10),
		ErrC: make(chan string, 10),
	}
	// spawn worker
	go JobCD(ch, &m)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			execCmd = fakeExecCommand
			defer func() { execCmd = exec.Command }()

			var gotErr uint8
			ch.JobC <- tt.job

			if tt.isErr {
				select {
				case <-ch.ErrC:
					gotErr++
				}
			}

			assert.Equal(t, tt.errCount, gotErr)
		})
	}
}
