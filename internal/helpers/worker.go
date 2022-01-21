package helpers

import "github.com/mdanialr/webhook/internal/config"

var WorkerChan chan string

// CDWorker Continuous Delivery worker will always listen to jobChan to receive
// the job and trigger PullRepo.
func CDWorker(jobChan <-chan string) {
	for job := range jobChan {
		NzLogInf.Println("Working on:", job)
		pRepo := pullRepoT{
			Name:    job,
			Service: config.Conf.Service,
		}
		out, err := pullRepo(pRepo)
		if err != nil {
			NzLogErr.Println(err)
			return
		}

		NzLogInf.Println(out)
	}
}
