package goroutines

/*
Module 10: Goroutines and Channels
====================================

Concurrency is Go's superpower. While most languages bolt on threading as an
afterthought, Go was designed from the ground up with concurrency as a
first-class citizen. The two core primitives are goroutines and channels.

What Are Goroutines?
--------------------
A goroutine is a lightweight thread of execution managed by the Go runtime.
They are NOT OS threads — they're much cheaper:

  - A goroutine starts with ~2KB of stack (an OS thread: ~1MB)
  - You can run millions of goroutines on a single machine
  - The Go scheduler multiplexes goroutines onto OS threads (M:N scheduling)

Start a goroutine with the 'go' keyword:

  go doSomething()       // Function call
  go func() { ... }()   // Anonymous function (note the () at the end!)

When main() returns, ALL goroutines are killed immediately — the program
doesn't wait for them. This is the #1 goroutine gotcha for beginners.

Channels: Communication Between Goroutines
-------------------------------------------
"Don't communicate by sharing memory; share memory by communicating."
  — Go Proverb

Channels are typed conduits for sending values between goroutines:

  ch := make(chan int)      // Unbuffered channel of ints
  ch := make(chan string, 5) // Buffered channel with capacity 5

Unbuffered vs Buffered
~~~~~~~~~~~~~~~~~~~~~~
Unbuffered channels synchronize sender and receiver:
  - Send blocks until someone receives
  - Receive blocks until someone sends
  - They force a "handshake" — both goroutines must be ready

Buffered channels are like a mailbox with a capacity:
  - Send blocks only when the buffer is full
  - Receive blocks only when the buffer is empty
  - They decouple sender and receiver timing

When to use which?
  - Unbuffered: When you need synchronization (handshake between goroutines)
  - Buffered: When you need to decouple producer/consumer speed (with a limit)
  - Buffered with cap 1: A simple semaphore

Channel Directions
------------------
You can restrict channel parameters to send-only or receive-only:

  func producer(out chan<- int) { ... }  // Can only SEND to out
  func consumer(in  <-chan int) { ... }  // Can only RECEIVE from in

Why bother? Because it makes the code self-documenting and prevents bugs.
If a function only needs to send, declaring it chan<- guarantees it won't
accidentally read from the channel.

Closing Channels
----------------
The sender closes a channel to signal "no more values":

  close(ch)

After closing:
  - Sending to a closed channel PANICS (this is a common bug!)
  - Receiving from a closed channel returns the zero value immediately
  - You can check: value, ok := <-ch  (ok is false if closed and empty)

The 'range' keyword works beautifully with channels:

  for value := range ch {
      // Receives values until ch is closed
  }

This is the idiomatic way to consume all values from a channel.

Channel Axioms (Memorize These!)
---------------------------------
  1. Send to a nil channel → blocks forever
  2. Receive from a nil channel → blocks forever
  3. Send to a closed channel → PANICS
  4. Receive from a closed channel → returns zero value, ok=false
  5. Close a nil channel → PANICS
  6. Close an already-closed channel → PANICS

These axioms explain almost every deadlock and panic you'll encounter.

Deadlocks
---------
A deadlock occurs when ALL goroutines are blocked. Go's runtime detects this
and crashes with "fatal error: all goroutines are asleep - deadlock!"

Common deadlock scenarios:
  - Single goroutine sending to an unbuffered channel (nobody to receive)
  - Two goroutines each waiting for the other to send
  - Forgetting to close a channel that a goroutine is ranging over

The Happens-Before Relationship
-------------------------------
Go's memory model defines "happens-before" relationships that guarantee
ordering. For channels:
  - A send on a channel happens-before the corresponding receive completes
  - The close of a channel happens-before a receive of the zero value

This means data sent through a channel is safely visible to the receiver.
You don't need additional synchronization.

Goroutine Leaks: A Real Production Problem
-------------------------------------------
A goroutine leak happens when you start a goroutine that never finishes.
Since goroutines are cheap, it's easy to create thousands without noticing —
until your service slowly consumes all available memory.

Common causes:
  - Goroutine blocked sending to a channel nobody reads
  - Goroutine blocked receiving from a channel nobody sends to
  - Goroutine in an infinite loop with no exit condition
  - Forgetting to handle context cancellation

In web services, this often happens when a request is cancelled (client
disconnects) but the goroutine handling the request keeps working. That's
why context.Context is so important (we'll cover it in the next module).

Prevention: Every goroutine should have a clear exit strategy.
*/

import (
	"fmt"
	"sync"
	"time"
)

// ==========================================
// Starting Goroutines
// ==========================================

// DemoBasicGoroutine shows how to launch a goroutine and wait for its result.
// In real code, you'd use sync.WaitGroup or channels — this is a minimal example.
func DemoBasicGoroutine() string {
	result := make(chan string)

	// Launch a goroutine. The 'go' keyword is all it takes.
	go func() {
		// This runs concurrently with the caller.
		// We send the result back through a channel.
		result <- "hello from goroutine"
	}()

	// Block until we receive the result.
	// Without this receive, main might exit before the goroutine finishes.
	return <-result
}

// DemoMultipleGoroutines launches several goroutines and collects results.
// Note: the order of results is NON-DETERMINISTIC. This is fundamental to
// concurrent programming.
func DemoMultipleGoroutines(n int) []int {
	results := make(chan int, n) // Buffered: goroutines won't block on send

	for i := 0; i < n; i++ {
		go func(id int) {
			// Each goroutine sends its ID through the channel
			results <- id * id
		}(i) // Pass i as argument — don't capture the loop variable!
	}

	// Collect all results. We know there will be exactly n.
	collected := make([]int, 0, n)
	for i := 0; i < n; i++ {
		collected = append(collected, <-results)
	}
	return collected
}

// ==========================================
// Unbuffered Channels (Synchronization)
// ==========================================

// DemoPingPong shows synchronization between two goroutines using
// unbuffered channels. Each send/receive is a synchronization point.
func DemoPingPong(rounds int) []string {
	ping := make(chan string) // Unbuffered — forces synchronization
	pong := make(chan string) // Unbuffered — forces synchronization
	done := make(chan []string)

	var log []string

	// Ping player
	go func() {
		for i := 0; i < rounds; i++ {
			ping <- "ping" // Send ping, block until pong player receives
			<-pong         // Wait for pong player to respond
		}
	}()

	// Pong player (also collects the log)
	go func() {
		for i := 0; i < rounds; i++ {
			msg := <-ping      // Wait for ping
			log = append(log, msg)
			pong <- "pong"     // Send pong back
			log = append(log, "pong")
		}
		done <- log
	}()

	return <-done
}

// ==========================================
// Buffered Channels
// ==========================================

// DemoBufferedChannel shows how buffered channels decouple producer and consumer.
func DemoBufferedChannel() []int {
	// Buffer of 3: producer can send 3 values without blocking
	ch := make(chan int, 3)

	// Producer fills the buffer without needing a receiver ready
	ch <- 10
	ch <- 20
	ch <- 30
	// ch <- 40 would block here because the buffer is full!

	// Consumer reads all values
	results := make([]int, 0, 3)
	results = append(results, <-ch) // 10
	results = append(results, <-ch) // 20
	results = append(results, <-ch) // 30

	return results
}

// ==========================================
// Channel Directions
// ==========================================

// producer only sends to the channel. The chan<- direction makes this explicit.
// Trying to receive from 'out' would be a compile error.
func producer(out chan<- int, count int) {
	for i := 0; i < count; i++ {
		out <- i
	}
	close(out) // The sender closes the channel
}

// transformer reads from 'in' and writes to 'out'. It squares each value.
func transformer(in <-chan int, out chan<- int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

// DemoPipeline shows a processing pipeline using directed channels.
// This is a fundamental pattern: producer → transform → consume.
func DemoPipeline(count int) []int {
	// Create the channels connecting pipeline stages
	raw := make(chan int)
	squared := make(chan int)

	// Start pipeline stages as goroutines
	go producer(raw, count)
	go transformer(raw, squared)

	// Consume the final output
	results := make([]int, 0, count)
	for v := range squared {
		results = append(results, v)
	}
	return results
}

// ==========================================
// Closing Channels and Range
// ==========================================

// DemoRangeOverChannel shows the idiomatic way to consume all values.
func DemoRangeOverChannel() []string {
	ch := make(chan string)

	go func() {
		words := []string{"Go", "is", "concurrent"}
		for _, w := range words {
			ch <- w
		}
		close(ch) // MUST close, or the range below will deadlock
	}()

	// 'range' over a channel receives until it's closed
	var result []string
	for word := range ch {
		result = append(result, word)
	}
	return result
}

// DemoCheckChannelClosed shows how to detect a closed channel.
func DemoCheckChannelClosed() (int, bool) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)

	// First receive: gets the value
	val, ok := <-ch
	if ok {
		// ok is true — we got a real value
		_ = val
	}

	// Second receive: channel is closed and empty
	val, ok = <-ch
	// ok is false, val is 0 (zero value for int)
	return val, ok
}

// ==========================================
// Goroutine Leak Demonstration
// ==========================================

// LeakyFunction demonstrates a goroutine leak.
// The goroutine tries to send to a channel, but if nobody reads from it,
// the goroutine blocks forever — it's leaked.
//
// DO NOT DO THIS IN PRODUCTION CODE. This is an anti-pattern.
func LeakyFunction() string {
	ch := make(chan string)

	go func() {
		// This goroutine will block here forever if nobody reads from ch!
		ch <- "I might leak"
	}()

	// Simulate: sometimes we read, sometimes we don't
	// In this case we do read, so no leak.
	return <-ch
}

// NonLeakyFunction shows the fix: use a buffered channel so the goroutine
// can always send and then exit, even if nobody receives.
func NonLeakyFunction() string {
	ch := make(chan string, 1) // Buffer of 1: send never blocks

	go func() {
		// Even if nobody reads, the goroutine can send to the buffer and exit.
		ch <- "I won't leak"
	}()

	return <-ch
}

// ==========================================
// Practical Pattern: Fan-Out Results
// ==========================================

// FetchResult represents the result of a concurrent fetch operation.
type FetchResult struct {
	URL   string
	Size  int
	Error error
}

// DemoConcurrentFetch simulates fetching multiple URLs concurrently.
// This is a common pattern in web services: fan out requests, collect results.
func DemoConcurrentFetch(urls []string) []FetchResult {
	results := make(chan FetchResult, len(urls))

	// Fan out: one goroutine per URL
	for _, url := range urls {
		go func(u string) {
			// Simulate fetching (in real code: http.Get(u))
			time.Sleep(1 * time.Millisecond) // Simulated latency
			results <- FetchResult{
				URL:  u,
				Size: len(u) * 10, // Fake response size
			}
		}(url)
	}

	// Collect all results
	var collected []FetchResult
	for i := 0; i < len(urls); i++ {
		collected = append(collected, <-results)
	}
	return collected
}

// ==========================================
// WaitGroup Preview (Detailed in Module 11)
// ==========================================

// DemoWaitGroup shows a quick preview of sync.WaitGroup.
// We'll cover this in depth in the concurrency patterns module.
func DemoWaitGroup() []int {
	var mu sync.Mutex
	var results []int
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			mu.Lock()
			results = append(results, n*n)
			mu.Unlock()
		}(i)
	}

	wg.Wait() // Block until all goroutines call Done()
	return results
}

// Go 1.25+ Modern Alternative:
// sync.WaitGroup now has a .Go() method that handles Add/Done automatically:
//
//   var wg sync.WaitGroup
//   for i := range 5 {
//       n := i
//       wg.Go(func() {
//           mu.Lock()
//           results = append(results, n*n)
//           mu.Unlock()
//       })
//   }
//   wg.Wait()
//
// This eliminates the most common WaitGroup bug: mismatched Add/Done calls.
// The classic pattern above is still valid and appears in most existing codebases.

// ==========================================
// Channel as Semaphore
// ==========================================

// DemoSemaphore uses a buffered channel as a semaphore to limit concurrency.
// This is a common pattern in web services: limit the number of concurrent
// database connections, HTTP requests, etc.
func DemoSemaphore(tasks int, maxConcurrent int) int {
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan int, tasks)

	for i := 0; i < tasks; i++ {
		go func(id int) {
			sem <- struct{}{} // Acquire semaphore (blocks if at capacity)
			defer func() { <-sem }() // Release semaphore

			// Simulate work
			results <- id
		}(i)
	}

	total := 0
	for i := 0; i < tasks; i++ {
		total += <-results
	}
	return total
}

// ==========================================
// Generator Pattern
// ==========================================

// GenerateSequence returns a channel that produces values from start to end.
// The caller owns the channel — read from it until it closes.
func GenerateSequence(start, end int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := start; i <= end; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

// FormatSequence demonstrates composing generators.
func FormatSequence(start, end int) []string {
	var result []string
	for n := range GenerateSequence(start, end) {
		result = append(result, fmt.Sprintf("#%d", n))
	}
	return result
}
