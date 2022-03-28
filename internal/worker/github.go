package worker

import (
	"bytes"
	"fmt"

	"github.com/mdanialr/webhook/internal/config"
)

// GithubActionWebhookWorker worker that receive job and execute CD job but using lookup by `id`
// instead by `name`. Mainly used for webhook that use GitHub actions that trigger webhook instead
// of the webhook from GitHub in repo's setting.
func GithubActionWebhookWorker(b BagOfChannels, m *config.Model) {
	for job := range b.GithubActionChan.JobC {
		b.GithubActionChan.InfC <- fmt.Sprintf("START working on: %v\n", job)

		// make sure repo exist
		r, err := m.Service.LookupRepoById(job)
		if err != nil {
			b.GithubActionChan.ErrC <- err.Error() + " id: " + job
			b.GithubActionChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		// setup and prepare command
		cmd := r.ParsePullCommand()

		// execute the command
		var stdErr, stdOut bytes.Buffer
		execCommand := execCmd("sh", "-c", cmd)
		execCommand.Stdout = &stdOut
		execCommand.Stderr = &stdErr
		err = execCommand.Run()
		if err != nil {
			b.GithubActionChan.ErrC <- fmt.Sprintf("failed to execute git pull from remote repo: %v ~detail: %s\n", err, stdErr.String())
			b.GithubActionChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		b.GithubActionChan.InfC <- stdOut.String()
		b.GithubActionChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
