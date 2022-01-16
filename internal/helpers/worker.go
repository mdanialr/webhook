package helpers

var WorkerChan chan string

// CDWorker Continuous Delivery worker will always listen to jobChan to receive
// the job and trigger PullRepo.
func CDWorker(jobChan <-chan string) {
	for job := range jobChan {
		NzLogInf.Println("Working on:", job)
		go pullRepo(job)
	}
}
