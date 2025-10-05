package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Shared resource: A simple counter that multiple goroutines will try to update.
var sharedCounter int

// 1. --- GLOBAL MUTEX ---
// This mutex protects access to sharedCounter.
var mu sync.Mutex

// This wait group will ensure the main function waits for all goroutines to finish.
var wg sync.WaitGroup

// The total number of increments we expect to perform.
const totalIncrements = 1000

// increment simulates a concurrent task that updates the shared resource.
func increment() {
	// The WaitGroup counter must be decremented when this goroutine finishes.
	defer wg.Done()

	// 2. --- ACQUIRE LOCK (MUTEX) ---
	// This blocks until the lock is available. Only one goroutine can pass this point at a time.
	mu.Lock()

	// 3. --- CRITICAL SECTION: INCREMENT THE SHARED COUNTER SAFELY ---
	sharedCounter++

	// 4. --- RELEASE LOCK (MUTEX) ---
	// This unblocks the lock, allowing the next waiting goroutine to proceed.
	mu.Unlock()

	// Simulate some work being done before exiting
	time.Sleep(time.Millisecond)
}

func main() {
	fmt.Println("--- Project 7: Concurrency and Synchronization ---")

	// We will launch 1000 goroutines, each trying to increment the counter once.
	for i := 0; i < totalIncrements; i++ {
		// Before launching a goroutine, tell the WaitGroup we expect one more job.
		wg.Add(1)
		// Launch the concurrent task
		go increment()
	}

	// Wait for all 1000 goroutines to call wg.Done()
	wg.Wait()

	// The expected result should now reliably be 1000.
	log.Printf("Final Counter Value: %d (Expected: %d)", sharedCounter, totalIncrements)
}
