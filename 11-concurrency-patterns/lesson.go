package concurrency

/*
Module 11: Concurrency Patterns
=================================

Module 10 covered the raw primitives — goroutines and channels. This module
covers the PATTERNS built on top of them. These patterns solve recurring
concurrency problems that come up constantly in production Go code, especially
in web services.

sync.WaitGroup
--------------
WaitGroup waits for a collection of goroutines to finish:

  var wg sync.WaitGroup
  for i := 0; i < 10; i++ {
      wg.Add(1)
      go func() {
          defer wg.Done()
          // do work
      }()
  }
  wg.Wait() // blocks until all 10 goroutines call Done()

Critical rules:
  - Call Add() BEFORE launching the goroutine (not inside it)
  - Call Done() with defer as the first statement in the goroutine
  - Never copy a WaitGroup after first use (pass by pointer)

Note: As of Go 1.25, WaitGroup has a .Go() method for cleaner usage:
  wg.Go(func() { ... })  // Handles Add and Done automatically

sync.Mutex and sync.RWMutex
----------------------------
When you DO need shared state (and sometimes you do), mutexes protect it:

  var mu sync.Mutex
  var count int

  mu.Lock()
  count++        // Only one goroutine can be here at a time
  mu.Unlock()

RWMutex allows multiple concurrent readers OR one writer:
  var mu sync.RWMutex
  mu.RLock()   // Multiple goroutines can hold read locks
  mu.RUnlock()
  mu.Lock()    // Exclusive — waits for all readers to finish
  mu.Unlock()

Use RWMutex when reads vastly outnumber writes (e.g., an in-memory cache
in a web service).

Gotcha: Always use defer mu.Unlock() to avoid forgetting to unlock,
especially if a function can return early or panic.

sync.Once
---------
Once ensures a function runs exactly once, no matter how many goroutines
call it:

  var once sync.Once
  once.Do(func() {
      // This runs exactly once, even if called from 100 goroutines
  })

Perfect for lazy initialization of expensive resources (database connections,
configuration loading).

The select Statement
---------------------
select is like switch but for channels. It waits on multiple channel
operations and proceeds with whichever is ready first:

  select {
  case msg := <-ch1:
      // ch1 sent a value
  case ch2 <- value:
      // Sent value to ch2
  case <-time.After(5 * time.Second):
      // Timeout after 5 seconds
  default:
      // No channel ready (non-blocking)
  }

Key behaviors:
  - If multiple cases are ready, one is chosen at random (fair)
  - Without default, select blocks until a case is ready
  - With default, select is non-blocking

context.Context
----------------
Context is the backbone of cancellation and deadline propagation in Go.
Every web service handler gets a context from the HTTP request:

  func handler(w http.ResponseWriter, r *http.Request) {
      ctx := r.Context() // Cancelled if the client disconnects
      result, err := doWork(ctx)
  }

Creating contexts:
  ctx, cancel := context.WithCancel(parentCtx)    // Manual cancellation
  ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
  ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5*time.Second))

ALWAYS call cancel() when done, even if the context expired:
  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()  // Always! Prevents resource leaks.

Checking for cancellation:
  select {
  case <-ctx.Done():
      return ctx.Err() // context.Canceled or context.DeadlineExceeded
  case result := <-ch:
      return result
  }

Worker Pool Pattern
--------------------
A fixed number of goroutines processing jobs from a shared channel:

  jobs := make(chan Job, 100)
  for i := 0; i < numWorkers; i++ {
      go worker(jobs)
  }

Why not one goroutine per job? Because:
  - Limits concurrent resource usage (DB connections, file handles)
  - Prevents overwhelming downstream services
  - Provides backpressure (jobs channel fills up → producers slow down)

Fan-Out / Fan-In
-----------------
Fan-out: Distribute work across multiple goroutines (each reads from the
same channel).
Fan-in: Merge results from multiple goroutines into one channel.

This is how you parallelize CPU-bound work or I/O operations.

Pipeline Pattern
-----------------
Chain of processing stages connected by channels:
  input → stage1 → stage2 → stage3 → output

Each stage is a goroutine that reads from one channel and writes to another.
This naturally limits memory usage (backpressure through channels) and
composes well.

Rate Limiting
--------------
Use time.Ticker or time.After to limit the rate of operations:

  limiter := time.NewTicker(100 * time.Millisecond) // 10 ops/second
  defer limiter.Stop()
  for request := range requests {
      <-limiter.C  // Wait for the next tick
      process(request)
  }

The Race Detector
------------------
Go has a built-in race detector. USE IT:

  go test -race ./...
  go run -race main.go

It detects data races at runtime — when two goroutines access the same
memory concurrently and at least one is a write. It slows down your program
(~2-10x) so don't run it in production, but ALWAYS use it in CI.

A race condition passing without -race doesn't mean it's correct. It means
the race didn't happen to trigger during that run.
*/

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ==========================================
// WaitGroup Pattern
// ==========================================

// DemoWaitGroup shows the standard WaitGroup pattern for launching
// multiple goroutines and waiting for all to complete.
func DemoWaitGroup() []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []string

	tasks := []string{"fetch users", "fetch orders", "fetch products"}

	for _, task := range tasks {
		wg.Add(1) // Add BEFORE launching goroutine
		go func(t string) {
			defer wg.Done() // Done when goroutine exits

			// Simulate work
			result := fmt.Sprintf("completed: %s", t)

			// Protect shared slice with mutex
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(task)
	}

	wg.Wait() // Block until all goroutines are done
	return results
}

// Go 1.25+ Modern Alternative:
// sync.WaitGroup now has a .Go() method that handles Add/Done automatically:
//
//   var wg sync.WaitGroup
//   for _, task := range tasks {
//       t := task
//       wg.Go(func() {
//           result := fmt.Sprintf("completed: %s", t)
//           mu.Lock()
//           results = append(results, result)
//           mu.Unlock()
//       })
//   }
//   wg.Wait()
//
// This eliminates the most common WaitGroup bug: mismatched Add/Done calls.
// The classic pattern above is still valid and appears in most existing codebases.

// ==========================================
// Mutex Pattern: Thread-Safe Counter
// ==========================================

// SafeCounter is a thread-safe counter protected by a mutex.
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// Increment safely increments the counter.
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // defer ensures unlock even if code below panics
	c.value++
}

// Value safely reads the counter.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// DemoMutex shows concurrent access to a shared counter.
func DemoMutex(goroutines, incrementsEach int) int {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsEach; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()
	return counter.Value()
}

// ==========================================
// RWMutex Pattern: Concurrent Cache
// ==========================================

// Cache is a thread-safe key-value cache using RWMutex.
// Multiple goroutines can read concurrently, but writes are exclusive.
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewCache creates a new thread-safe cache.
func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

// Get reads a value (concurrent-safe with other reads).
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock() // Read lock — multiple readers allowed
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

// Set writes a value (exclusive access).
func (c *Cache) Set(key, value string) {
	c.mu.Lock() // Write lock — exclusive access
	defer c.mu.Unlock()
	c.data[key] = value
}

// ==========================================
// sync.Once Pattern
// ==========================================

// DBConnection simulates a database connection that's expensive to create.
type DBConnection struct {
	Host string
}

var (
	dbInstance *DBConnection
	dbOnce     sync.Once
)

// GetDB returns the singleton database connection. It's created exactly
// once, even if called from multiple goroutines simultaneously.
func GetDB() *DBConnection {
	dbOnce.Do(func() {
		// This runs exactly once, regardless of how many goroutines call GetDB
		dbInstance = &DBConnection{Host: "localhost:5432"}
	})
	return dbInstance
}

// ==========================================
// Select Statement
// ==========================================

// DemoSelect shows the select statement multiplexing multiple channels.
func DemoSelect() string {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func() {
		time.Sleep(1 * time.Millisecond)
		ch1 <- "from channel 1"
	}()
	go func() {
		time.Sleep(2 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// Wait for whichever channel is ready first
	select {
	case msg := <-ch1:
		return msg
	case msg := <-ch2:
		return msg
	case <-time.After(1 * time.Second):
		return "timeout"
	}
}

// ==========================================
// Context: Cancellation
// ==========================================

// DemoContextCancellation shows how to propagate cancellation.
func DemoContextCancellation() (string, error) {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		// Simulate a long-running operation that respects context
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		case <-time.After(10 * time.Second):
			resultCh <- "completed"
		}
	}()

	// Cancel after a short time
	time.Sleep(5 * time.Millisecond)
	cancel() // Signal all goroutines using this context to stop

	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return "", err
	}
}

// ==========================================
// Context: Timeout
// ==========================================

// DemoContextTimeout shows automatic cancellation after a deadline.
func DemoContextTimeout() error {
	// Context auto-cancels after 50ms
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel() // ALWAYS defer cancel, even with timeout

	// Simulate work that checks for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err() // context.DeadlineExceeded
	case <-time.After(1 * time.Second):
		return nil // Won't reach this — timeout fires first
	}
}

// ==========================================
// Worker Pool Pattern
// ==========================================

// Job represents a unit of work.
type Job struct {
	ID    int
	Input int
}

// Result represents the output of processing a Job.
type Result struct {
	JobID  int
	Output int
}

// DemoWorkerPool demonstrates the worker pool pattern.
func DemoWorkerPool(numWorkers int, jobs []Job) []Result {
	jobCh := make(chan Job, len(jobs))
	resultCh := make(chan Result, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				// Process the job (square the input)
				resultCh <- Result{
					JobID:  job.ID,
					Output: job.Input * job.Input,
				}
			}
		}()
	}

	// Send all jobs
	for _, job := range jobs {
		jobCh <- job
	}
	close(jobCh) // Signal workers that no more jobs are coming

	// Wait for all workers to finish, then close results
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	var results []Result
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}

// ==========================================
// Fan-Out / Fan-In Pattern
// ==========================================

// DemoFanOutFanIn distributes work across multiple goroutines and merges results.
func DemoFanOutFanIn(input []int, numWorkers int) []int {
	inputCh := make(chan int, len(input))

	// Fan-out: each worker reads from the same input channel
	var workerChannels []<-chan int
	for i := 0; i < numWorkers; i++ {
		ch := make(chan int)
		workerChannels = append(workerChannels, ch)
		go func(out chan<- int) {
			defer close(out)
			for n := range inputCh {
				out <- n * 2 // Each worker doubles the value
			}
		}(ch)
	}

	// Send input
	for _, v := range input {
		inputCh <- v
	}
	close(inputCh)

	// Fan-in: merge all worker output channels into one
	merged := fanIn(workerChannels...)

	var results []int
	for v := range merged {
		results = append(results, v)
	}
	return results
}

// fanIn merges multiple channels into one.
func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				merged <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// ==========================================
// Rate Limiter Pattern
// ==========================================

// RateLimiter limits the rate of operations using a ticker.
type RateLimiter struct {
	ticker *time.Ticker
}

// NewRateLimiter creates a rate limiter that allows one operation per interval.
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(interval),
	}
}

// Wait blocks until the next tick, enforcing the rate limit.
func (r *RateLimiter) Wait() {
	<-r.ticker.C
}

// Stop stops the rate limiter. Always call this when done.
func (r *RateLimiter) Stop() {
	r.ticker.Stop()
}

// ==========================================
// Pipeline Pattern
// ==========================================

// DemoPipeline shows a multi-stage pipeline with proper shutdown.
func DemoPipeline(input []int) []string {
	// Stage 1: Generate values
	stage1 := make(chan int)
	go func() {
		defer close(stage1)
		for _, v := range input {
			stage1 <- v
		}
	}()

	// Stage 2: Filter (keep even numbers)
	stage2 := make(chan int)
	go func() {
		defer close(stage2)
		for v := range stage1 {
			if v%2 == 0 {
				stage2 <- v
			}
		}
	}()

	// Stage 3: Transform (format as strings)
	stage3 := make(chan string)
	go func() {
		defer close(stage3)
		for v := range stage2 {
			stage3 <- fmt.Sprintf("even:%d", v)
		}
	}()

	// Collect results
	var results []string
	for s := range stage3 {
		results = append(results, s)
	}
	return results
}
