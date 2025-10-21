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

	nowStart := time.Now();
	fmt.Printf("[Worker %d] Starting work...\n", id)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	duration := time.Since(nowStart)

	res := Result {
		Message: fmt.Sprintf("[Worker %d] Done!", id),
		Value: id * 2,
		Duration: duration,
	}
	ch <- res
}


func main()	{
	const numWorkers = 5
	var wgWorkers sync.WaitGroup //WaitGroup variable
	var wgReceivers sync.WaitGroup 

	//Create channel
	results := make(chan Result)

	//Start workers
	for i := 1; i<= numWorkers; i++ {
		wgWorkers.Add(1) // Announce a new worker will start
		go worker(i, results, &wgWorkers)
	}

	//receiver goroutine
	wgReceivers.Add(1)
	go func()	{
		defer wgReceivers.Done()
		for r := range results {
			fmt.Print(r.Message)
			fmt.Printf(" with value: %d and a duration of: %.2f seconds\n", r.Value, r.Duration.Seconds())
		}
		fmt.Println("[Manager] All results received — channel closed.")
	}()
	
	wgWorkers.Wait() //Wait for all the workers then close
	close(results)
	wgReceivers.Wait() //Wait for receiver to finish printing.

	fmt.Println("Main goroutine exiting.") //kills all goroutines.
}