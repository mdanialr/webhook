package worker

import (
	"fmt"
	"os/exec"

	"github.com/mdanialr/webhook/internal/config"
)

var execCmd = exec.Command

// Channel used by worker to exchange messages, either receive job
// or send any information.
type Channel struct {
	JobC chan string // receive job
	InfC chan string // send any information
	ErrC chan string // send any error information
}

// JobCD worker that would always listen to job channel and do
// continuous delivery based on the repo's name.
func JobCD(ch *Channel, m *config.Model) {
	for job := range ch.JobC {
		ch.InfC <- fmt.Sprintf("START working on: %v\n", job)

		// make sure repo exist
		r, err := m.Service.LookupRepo(job)
		if err != nil {
			ch.ErrC <- err.Error()
		}

		// setup and prepare command
		cmd := r.ParsePullCommand()

		// execute the command
		res, err := execCmd("sh", "-c", cmd).CombinedOutput()
		if err != nil {
			ch.ErrC <- fmt.Sprintf("failed to execute git pull from remote repo: %v\n", err)
		}

		ch.InfC <- string(res)
		ch.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
