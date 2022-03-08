package worker

import "os/exec"

// Channel used by worker to exchange messages, either receive job or send any information.
type Channel struct {
	JobC chan string // receive job
	InfC chan string // send any information
	ErrC chan string // send any error information
}

// BagOfChannels contain all channels that used by worker.
type BagOfChannels struct {
	GithubActionChan  *Channel // Channel for GitHub action worker.
	GithubWebhookChan *Channel // Channel for old school GitHub webhook CD worker.
	DockerWebhookChan *Channel // Channel for docker hub webhook CD worker.
}

// execCmd to make it possible to test exec.Command.
var execCmd = exec.Command
