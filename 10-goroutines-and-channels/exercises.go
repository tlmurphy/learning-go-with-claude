package goroutines

/*
=============================================================================
 EXERCISES: Goroutines and Channels
=============================================================================

 Work through these exercises in order. Run tests with:

   make test 10

 Run a single test:

   go test -v -run TestSquareNumbers ./10-goroutines-and-channels/

=============================================================================
*/

import (
	"time"
)

// =============================================================================
// Exercise 1: Launch and Collect
// =============================================================================

// SquareNumbers launches a goroutine for each number in the input slice,
// computes its square concurrently, and returns all the squares.
//
// Requirements:
//   - Each square computation must happen in its own goroutine
//   - Use a channel to collect results
//   - Return a slice of all squared values
//   - The order of results does NOT need to match input order
//
// Hint: Use a buffered channel with capacity len(numbers).
func SquareNumbers(numbers []int) []int {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 2: Producer-Consumer
// =============================================================================

// ProduceConsume implements a producer-consumer pattern.
//
// The producer should:
//   - Generate integers from 1 to count (inclusive)
//   - Send each integer to a buffered channel of the given bufferSize
//   - Close the channel when done
//
// The consumer should:
//   - Read all values from the channel (use range)
//   - Double each value
//   - Collect the doubled values into a slice
//
// Returns the slice of doubled values.
//
// Both producer and consumer should be separate goroutines.
// The function waits for the consumer to finish and returns its results.
func ProduceConsume(count int, bufferSize int) []int {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 3: Pipeline with Channel Directions
// =============================================================================

// BuildPipeline creates a three-stage processing pipeline:
//
// Stage 1 (generate): Sends integers 1 through count to a channel, then closes it.
// Stage 2 (double):   Reads from stage 1's channel, doubles each value,
//
//	sends to the next channel, then closes it.
//
// Stage 3 (offset):   Reads from stage 2's channel, adds 'offset' to each value,
//
//	sends to the output channel, then closes it.
//
// Use proper channel directions (chan<- and <-chan) in your helper functions
// or inline goroutines to enforce correct usage.
//
// Returns a receive-only channel that produces the final output.
func BuildPipeline(count int, offset int) <-chan int {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 4: Ping Pong
// =============================================================================

// PingPong simulates a ping-pong game between two goroutines.
//
// Rules:
//   - Two goroutines alternate sending a counter back and forth
//   - The counter starts at 0 and increments by 1 each time it's sent
//   - After 'rounds' total exchanges, the game ends
//   - Return the final counter value
//
// Example: PingPong(4) → counter goes 0→1→2→3→4, returns 4
//
// Use two channels to coordinate the exchange.
func PingPong(rounds int) int {
	// YOUR CODE HERE
	return 0
}

// =============================================================================
// Exercise 5: Number Generator with Close Signaling
// =============================================================================

// Fibonacci returns a receive-only channel that produces the first n
// Fibonacci numbers, then closes.
//
// The Fibonacci sequence: 0, 1, 1, 2, 3, 5, 8, 13, 21, ...
//   - Fibonacci(0) → channel closes immediately (no values)
//   - Fibonacci(1) → sends 0, then closes
//   - Fibonacci(5) → sends 0, 1, 1, 2, 3, then closes
//
// The generation must happen in a goroutine. The returned channel
// should be closed when all values have been sent.
func Fibonacci(n int) <-chan int {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 6: Fix the Deadlock
// =============================================================================

// BrokenCollect has a deadlock bug. Find and fix it.
//
// It's supposed to launch goroutines that each send a value to a channel,
// then collect all the values into a slice.
//
// The fix should be minimal — change as little as possible.
//
// Hint: Think about whether the channel needs to be buffered, or whether
// sends and receives need to happen concurrently.
func BrokenCollect(values []string) []string {
	ch := make(chan string) // Bug is related to this line or the structure below

	// YOUR CODE HERE — fix the deadlock
	// The original (broken) logic:
	//   1. Send all values through the channel in goroutines
	//   2. Collect all results from the channel
	// But the current structure deadlocks. Fix it.

	for _, v := range values {
		_ = v
		// go func(s string) { ch <- s }(v)  // Uncomment and fix as needed
	}

	var result []string
	for range values {
		_ = <-ch
	}

	return result
}

// =============================================================================
// Exercise 7: Timeout Pattern
// =============================================================================

// SlowOperation simulates a slow operation that takes the given duration.
// It returns "completed" when done.
// Do NOT modify this function.
func SlowOperation(duration time.Duration) chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(duration)
		ch <- "completed"
	}()
	return ch
}

// WithTimeout runs SlowOperation with the given work duration, but gives up
// if it takes longer than the timeout.
//
// Returns:
//   - ("completed", true)  if the operation finishes within the timeout
//   - ("timeout", false)   if the timeout expires first
//
// Use select with time.After to implement the timeout.
// Do NOT use time.Sleep.
func WithTimeout(workDuration, timeout time.Duration) (string, bool) {
	// YOUR CODE HERE
	return "", false
}

// =============================================================================
// Exercise 8: Channel-Based Message Broker
// =============================================================================

// Broker is a simple publish-subscribe message broker.
// Subscribers register a channel to receive messages.
// When a message is published, ALL subscribers receive it.
type Broker struct {
	subscribers []chan string
}

// NewBroker creates a new message broker.
func NewBroker() *Broker {
	// YOUR CODE HERE
	return nil
}

// Subscribe adds a new subscriber and returns their message channel.
// The channel should be buffered (capacity 10) to prevent slow subscribers
// from blocking the publisher.
func (b *Broker) Subscribe() <-chan string {
	// YOUR CODE HERE
	return nil
}

// Publish sends a message to ALL subscribers.
// Use a non-blocking send (select with default) so a slow subscriber
// doesn't block delivery to others. If a subscriber's buffer is full,
// skip them for this message.
func (b *Broker) Publish(msg string) {
	// YOUR CODE HERE
}

// Close closes all subscriber channels, signaling no more messages.
func (b *Broker) Close() {
	// YOUR CODE HERE
}
