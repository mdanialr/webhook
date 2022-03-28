package worker

import (
	"bytes"
	"fmt"

	"github.com/mdanialr/webhook/internal/config"
)

// DockerCDWorker worker that would always listen to job channel and do
// continuous delivery based on the docker's id.
func DockerCDWorker(b BagOfChannels, conf *config.Model) {
	for job := range b.DockerWebhookChan.JobC {
		b.DockerWebhookChan.InfC <- fmt.Sprintf("START working on: %v\n", job)

		// make sure docker's id exist
		dock, err := conf.Dockers.LookupRepo(job)
		if err != nil {
			b.DockerWebhookChan.ErrC <- err.Error() + " id: " + job
			b.DockerWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		// setup and prepare command
		cmd := dock.ParsePullCommand()

		// execute the command
		var stdErr, stdOut bytes.Buffer
		execCommand := execCmd("sh", "-c", cmd)
		execCommand.Stdout = &stdOut
		execCommand.Stderr = &stdErr
		err = execCommand.Run()
		if err != nil {
			b.DockerWebhookChan.ErrC <- fmt.Sprintf("failed to execute git pull from remote repo: %v ~detail: %s\n", err, stdErr.String())
			b.DockerWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		b.DockerWebhookChan.InfC <- stdOut.String()
		b.DockerWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
