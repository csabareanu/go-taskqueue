//Goroutines : Running workers and receivers concurrently
//Channels: Safely sending and receiving typed data
//WaitGroups: Coordinating clean shutdowns
//Channel fan-out: Having multiple receivers read from one channel concurrently
//Graceful channel closing: Ending communication cleanly with no panics
//Multiple results per worker: Streaming data from producers to consumers


package main

import (
	"fmt"
	"time"
	"sync"
	"math/rand"
)

type Result struct {
    Message  string
    Value    int
	Duration time.Duration
}

func worker(id int, ch chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done() //Mark this worker as done then the function returns. defer-run at the end of the function.
	for i:=1; i<=3; i++ {
		nowStart := time.Now();
		fmt.Printf("[Worker %d] Starting work on task %d...\n", id, i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		duration := time.Since(nowStart)

		ch <- Result {
			Message: fmt.Sprintf("[Worker %d, Task %d] Done!", id, i),
			Value: id * i * 10,
			Duration: duration,
		}
	}
	fmt.Printf("[Worker %d] Finished all tasks.\n", id)
}

func receiver(id int, ch <-chan Result, wg *sync.WaitGroup) {
	defer wg.Done()
	counter :=0
	for res := range ch {
		counter++
		fmt.Print("[Receiver 1] ")
		fmt.Print(res.Message)
		fmt.Printf(" with value: %d and a duration of: %.2f seconds\n", res.Value, res.Duration.Seconds())
	}
	fmt.Printf("[Manager] All results received — channel closed. Receiver %d messages\n", counter)
}

func main()	{
	const numWorkers = 3
	var wgWorkers sync.WaitGroup //WaitGroup variable
	var wgReceivers sync.WaitGroup 

	//Create channel
	results := make(chan Result)

	//Start the receivers
	wgReceivers.Add(2)
	go receiver(1, results, &wgReceivers)
	go receiver(2, results, &wgReceivers)

	//Start workers
	for i := 1; i<= numWorkers; i++ {
		wgWorkers.Add(1) // Announce a new worker will start
		go worker(i, results, &wgWorkers)
	}

	wgWorkers.Wait() //Wait for all the workers then close
	close(results)
	wgReceivers.Wait() //Wait for receiver to finish printing.

	fmt.Println("Main goroutine exiting.") //kills all goroutines.
}