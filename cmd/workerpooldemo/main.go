//You have a fixed number of worker goroutines reading from a shared channel of jobs.
//Each worker takes one job, processes it, and (optionally) sends a result into another channel.
//This is useful because spawning one goroutine per job can explode memory so you can reuse 
//a fixed number of workers to continously process jobs from a queue.
package main

import (
	"sync"
	"fmt"
	"math/rand"
	"time"
	"context"
)

type Job struct {
	ID      int
	Payload int   //data to process
	RetryCount int //how many retries until now
	MaxRetries int //max retry attempts
	Backoff time.Duration
}

type Result struct {
	JobID    int
	WorkerID int
	Output   int
	Duration time.Duration
	Status string // "success/failed"
}

//Gets jobs from the job queue and sends results to results channel
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, retries chan<- Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Worker %d] Context canceled, exiting.\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("[Worker %d] No more jobs, exiting. \n", id)
				return
			}
			start := time.Now()
			if rand.Float32() < 0.70 {
				fmt.Printf("[Worker %d] Failed job %d (attempt %d)\n", id, job.ID, job.RetryCount + 1)
				job.RetryCount++
				if job.RetryCount <= job.MaxRetries	{
					retries <- job 
				} else {
					results <- Result {JobID: job.ID, WorkerID: id, Output: -1, Duration: time.Since(start), Status:"failed"}
				}	
				continue;
			}
			fmt.Printf("[Worker %d] Processing job %d (payload %d)\n", id, job.ID, job.Payload)
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			results <- Result {JobID: job.ID, WorkerID: id, Output: job.Payload * 2, Duration: time.Since(start), Status:"success"}
		}
	}
	fmt.Printf("[Worker %d] Exiting\n", id)
}

func retryManager(ctx context.Context, retries <-chan Job, jobs chan<- Job, wg *sync.WaitGroup)	{
	defer wg.Done()
	for {
		select{
		case <- ctx.Done():
			fmt.Println("[RetryManager] Context canceled — stopping retries.")
			return

		case job, ok := <-retries:
			if !ok {
				fmt.Println("[RetryManager] Retry channel closed — done.")
				return
			}
			backoff := job.Backoff * time.Duration(1<<job.RetryCount)
			fmt.Printf("[RetryManager] Scheduling retry for job %d (attempt %d) in %v\n",
				job.ID, job.RetryCount, backoff)
			
			go func(j Job, d time.Duration)	{
				select {
				case <-ctx.Done():
					return
				case <-time.After(d):
					fmt.Printf("[Retry Manager] Requeueing job %d (attempt %d) \n", j.ID, j.RetryCount)
					jobs <- j
				}
			}(job, backoff)
		}
	}
}

func receiver(ctx context.Context, ch <-chan Result, wg *sync.WaitGroup)	{
	defer wg.Done()
	success, failed := 0, 0
	for {
		select { //concurrency version of a switch. Waits for whichever channel operations becomes ready
		case <- ctx.Done():
			fmt.Println("[Receiver] Context Canceled - stopping.")
			return
		case res, ok := <- ch:
			if !ok {
				fmt.Printf("[Receiver] Channel closed. Success: %d, Fail: %d\n", success, failed)
				return
			}
			if res.Status == "failed" {
				failed++
				fmt.Printf("❌ Job %d permanently failed (Worker %d)\n", res.JobID, res.WorkerID)
			}  else {
				success++
				fmt.Printf("✅ Job %d done by Worker %d → %d (%.2fs)\n",
					res.JobID, res.WorkerID, res.Output, res.Duration.Seconds())	
			}
		}
	}
	fmt.Printf("[Receiver] Summary: %d successes, %d failures\n", success, failed)
}

func main() {
	start := time.Now()
	const numWorkers = 3
	const numJobs = 20

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	retries := make(chan Job, numJobs)

	var wgWorkers sync.WaitGroup
	var wgReceiver sync.WaitGroup
	var wgRetry sync.WaitGroup

	//start retry manager
	wgRetry.Add(1)
	go retryManager(ctx, retries, jobs, &wgRetry)

	//Start workers
	for w:=1; w<=numWorkers; w++ {
		wgWorkers.Add(1)
		go worker(ctx, w, jobs, results, retries, &wgWorkers)
	}

	//Send jobs to the job queue
	for j:=1; j<=numJobs; j++ {		
		jobs <- Job {ID: j, Payload: rand.Intn(100), MaxRetries: 3, Backoff: 1 * time.Second}
	}
	// close(jobs) //all jobs are sent to the channel. No sending back to the channel so safe to close.

	wgReceiver.Add(1)
	go receiver(ctx, results, &wgReceiver)

	wgWorkers.Wait()
	close(retries)
	wgRetry.Wait()
	close(results)
	wgReceiver.Wait()

	fmt.Printf("All jobs processed. Main exiting. Total time: %.2f seconds\n", time.Since(start).Seconds())

}
