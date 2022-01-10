package helpers

import "log"

var WorkerChan chan string

// CDWorker Continuous Delivery worker will always listen to jobChan to receive
// the job and trigger PullRepo.
func CDWorker(jobChan <-chan string) {
	for job := range jobChan {
		log.Println("Working on:", job)
		pullRepo(job)
	}
}
