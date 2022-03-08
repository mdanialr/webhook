package worker

import (
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
		res, err := execCmd("sh", "-c", cmd).CombinedOutput()
		if err != nil {
			b.DockerWebhookChan.ErrC <- fmt.Sprintf("failed to execute docker pull commands: %v\n", err)
			b.DockerWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		b.DockerWebhookChan.InfC <- string(res)
		b.DockerWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
