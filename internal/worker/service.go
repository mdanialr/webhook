package worker

import (
	"fmt"

	"github.com/mdanialr/webhook/internal/config"
)

// GithubCDWorker worker that would always listen to job channel and do
// continuous delivery based on the repo's name.
func GithubCDWorker(b BagOfChannels, m *config.Model) {
	for job := range b.GithubWebhookChan.JobC {
		b.GithubWebhookChan.InfC <- fmt.Sprintf("START working on: %v\n", job)

		// make sure repo exist
		r, err := m.Service.LookupRepo(job)
		if err != nil {
			b.GithubWebhookChan.ErrC <- err.Error() + " id: " + job
			b.GithubWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		// setup and prepare command
		cmd := r.ParsePullCommand()

		// execute the command
		res, err := execCmd("sh", "-c", cmd).CombinedOutput()
		if err != nil {
			b.GithubWebhookChan.ErrC <- fmt.Sprintf("failed to execute git pull from remote repo: %v\n", err)
			b.GithubWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		b.GithubWebhookChan.InfC <- string(res)
		b.GithubWebhookChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
