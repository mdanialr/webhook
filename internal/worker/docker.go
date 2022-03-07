package worker

import (
	"fmt"

	"github.com/mdanialr/webhook/internal/config"
)

// DockerChannel used by worker to exchange messages, either receive job
// or send any information.
type DockerChannel struct {
	JobC chan string // receive job
	InfC chan string // send any information
	ErrC chan string // send any error information
}

// DockerCDWorker worker that would always listen to job channel and do
// continuous delivery based on the docker's id.
func DockerCDWorker(ch *DockerChannel, conf *config.Model) {
	for job := range ch.JobC {
		ch.InfC <- fmt.Sprintf("START working on: %v\n", job)

		// make sure docker's id exist
		dock, err := conf.Dockers.LookupRepo(job)
		if err != nil {
			ch.ErrC <- err.Error()
			return
		}

		// setup and prepare command
		cmd := dock.ParsePullCommand()

		// execute the command
		res, err := execCmd("sh", "-c", cmd).CombinedOutput()
		if err != nil {
			ch.ErrC <- fmt.Sprintf("failed to execute docker pull commands: %v\n", err)
			return
		}

		ch.InfC <- string(res)
		ch.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
