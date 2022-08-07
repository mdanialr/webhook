package worker

import (
	"bytes"
	"fmt"

	"github.com/mdanialr/webhook/internal/model"
)

// GithubActionWebhookWorker worker that receive job and execute CD job but using lookup by `id`
// instead by `name`. Mainly used for webhook that use GitHub actions that trigger webhook instead
// of the webhook from GitHub in repo's setting.
func GithubActionWebhookWorker(b BagOfChannels, svc *model.Service) {
	for job := range b.GithubActionChan.JobC {
		b.GithubActionChan.InfC <- fmt.Sprintf("START working on: %v\n", job)

		// make sure repo exist
		repo, err := svc.LookupGithub(job)
		if err != nil {
			b.GithubActionChan.ErrC <- err.Error()
			b.GithubActionChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		// execute the command
		var stdErr, stdOut bytes.Buffer
		execCommand := execCmd("sh", "-c", repo.ParseCommand())
		execCommand.Stdout = &stdOut
		execCommand.Stderr = &stdErr
		err = execCommand.Run()
		if err != nil {
			b.GithubActionChan.ErrC <- fmt.Sprintf("failed to execute commands: %v ~detail: %s\n", err, stdErr.String())
			b.GithubActionChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
			return
		}

		b.GithubActionChan.InfC <- stdOut.String()
		b.GithubActionChan.InfC <- fmt.Sprintf("DONE working on: %v\n", job)
	}
}
