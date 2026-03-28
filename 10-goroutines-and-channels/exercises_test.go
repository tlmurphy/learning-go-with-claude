package goroutines

import (
	"sort"
	"testing"
	"time"
)

// =============================================================================
// Exercise 1: Launch and Collect
// =============================================================================

func TestSquareNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int // sorted for comparison since order is non-deterministic
	}{
		{
			name:  "basic squares",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{1, 4, 9, 16, 25},
		},
		{
			name:  "includes zero",
			input: []int{0, 3, 7},
			want:  []int{0, 9, 49},
		},
		{
			name:  "single element",
			input: []int{10},
			want:  []int{100},
		},
		{
			name:  "empty input",
			input: []int{},
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SquareNumbers(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("SquareNumbers(%v) returned %d values, want %d", tt.input, len(got), len(tt.want))
			}

			// Sort both for comparison (order is non-deterministic)
			sort.Ints(got)
			sort.Ints(tt.want)

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("sorted result[%d] = %d, want %d (full result: %v)", i, got[i], tt.want[i], got)
				}
			}
		})
	}
}

// =============================================================================
// Exercise 2: Producer-Consumer
// =============================================================================

func TestProduceConsume(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		bufferSize int
		want       []int
	}{
		{
			name:       "five items buffered",
			count:      5,
			bufferSize: 3,
			want:       []int{2, 4, 6, 8, 10},
		},
		{
			name:       "unbuffered channel",
			count:      3,
			bufferSize: 0,
			want:       []int{2, 4, 6},
		},
		{
			name:       "single item",
			count:      1,
			bufferSize: 1,
			want:       []int{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProduceConsume(tt.count, tt.bufferSize)

			if len(got) != len(tt.want) {
				t.Fatalf("ProduceConsume(%d, %d) returned %d values, want %d",
					tt.count, tt.bufferSize, len(got), len(tt.want))
			}

			// Producer sends 1..count in order, consumer reads in order
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("result[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// =============================================================================
// Exercise 3: Pipeline with Channel Directions
// =============================================================================

func TestBuildPipeline(t *testing.T) {
	tests := []struct {
		name   string
		count  int
		offset int
		want   []int
	}{
		{
			name:   "double then add 10",
			count:  5,
			offset: 10,
			want:   []int{12, 14, 16, 18, 20}, // (1*2+10, 2*2+10, 3*2+10, 4*2+10, 5*2+10)
		},
		{
			name:   "double then add 0",
			count:  3,
			offset: 0,
			want:   []int{2, 4, 6},
		},
		{
			name:   "single value",
			count:  1,
			offset: 100,
			want:   []int{102},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := BuildPipeline(tt.count, tt.offset)
			if ch == nil {
				t.Fatal("BuildPipeline returned nil channel")
			}

			var got []int
			for v := range ch {
				got = append(got, v)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("pipeline produced %d values, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("pipeline[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// =============================================================================
// Exercise 4: Ping Pong
// =============================================================================

func TestPingPong(t *testing.T) {
	tests := []struct {
		name   string
		rounds int
		want   int
	}{
		{"zero rounds", 0, 0},
		{"one round", 1, 1},
		{"four rounds", 4, 4},
		{"ten rounds", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PingPong(tt.rounds)
			if got != tt.want {
				t.Errorf("PingPong(%d) = %d, want %d", tt.rounds, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Exercise 5: Fibonacci Generator
// =============================================================================

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"zero values", 0, nil},
		{"one value", 1, []int{0}},
		{"two values", 2, []int{0, 1}},
		{"five values", 5, []int{0, 1, 1, 2, 3}},
		{"eight values", 8, []int{0, 1, 1, 2, 3, 5, 8, 13}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := Fibonacci(tt.n)
			if ch == nil {
				t.Fatal("Fibonacci returned nil channel")
			}

			var got []int
			for v := range ch {
				got = append(got, v)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("Fibonacci(%d) produced %d values, want %d: got %v",
					tt.n, len(got), len(tt.want), got)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Fibonacci(%d)[%d] = %d, want %d", tt.n, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// =============================================================================
// Exercise 6: Fix the Deadlock
// =============================================================================

func TestBrokenCollect(t *testing.T) {
	t.Run("collects all values", func(t *testing.T) {
		input := []string{"alpha", "beta", "gamma"}

		// Use a timeout to detect deadlocks
		done := make(chan []string, 1)
		go func() {
			done <- BrokenCollect(input)
		}()

		select {
		case got := <-done:
			if len(got) != len(input) {
				t.Fatalf("BrokenCollect returned %d values, want %d", len(got), len(input))
			}
			// Values may be in any order, so sort both
			sort.Strings(got)
			sorted := make([]string, len(input))
			copy(sorted, input)
			sort.Strings(sorted)
			for i := range got {
				if got[i] != sorted[i] {
					t.Errorf("sorted result[%d] = %q, want %q", i, got[i], sorted[i])
				}
			}
		case <-time.After(2 * time.Second):
			t.Fatal("BrokenCollect deadlocked (timed out after 2 seconds)")
		}
	})

	t.Run("handles empty input", func(t *testing.T) {
		done := make(chan []string, 1)
		go func() {
			done <- BrokenCollect([]string{})
		}()

		select {
		case got := <-done:
			if len(got) != 0 {
				t.Errorf("expected empty result for empty input, got %v", got)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("BrokenCollect deadlocked on empty input")
		}
	})
}

// =============================================================================
// Exercise 7: Timeout Pattern
// =============================================================================

func TestWithTimeout(t *testing.T) {
	t.Run("completes before timeout", func(t *testing.T) {
		result, ok := WithTimeout(10*time.Millisecond, 500*time.Millisecond)
		if !ok {
			t.Error("expected operation to complete, but it timed out")
		}
		if result != "completed" {
			t.Errorf("result = %q, want %q", result, "completed")
		}
	})

	t.Run("times out", func(t *testing.T) {
		result, ok := WithTimeout(500*time.Millisecond, 10*time.Millisecond)
		if ok {
			t.Error("expected timeout, but operation completed")
		}
		if result != "timeout" {
			t.Errorf("result = %q, want %q", result, "timeout")
		}
	})
}

// =============================================================================
// Exercise 8: Message Broker
// =============================================================================

func TestBroker(t *testing.T) {
	t.Run("single subscriber receives messages", func(t *testing.T) {
		broker := NewBroker()
		if broker == nil {
			t.Fatal("NewBroker() returned nil")
		}

		sub := broker.Subscribe()
		if sub == nil {
			t.Fatal("Subscribe() returned nil")
		}

		broker.Publish("hello")
		broker.Publish("world")

		msg1 := <-sub
		if msg1 != "hello" {
			t.Errorf("first message = %q, want %q", msg1, "hello")
		}

		msg2 := <-sub
		if msg2 != "world" {
			t.Errorf("second message = %q, want %q", msg2, "world")
		}
	})

	t.Run("multiple subscribers all receive", func(t *testing.T) {
		broker := NewBroker()
		if broker == nil {
			t.Fatal("NewBroker() returned nil")
		}

		sub1 := broker.Subscribe()
		sub2 := broker.Subscribe()
		sub3 := broker.Subscribe()

		broker.Publish("broadcast")

		for i, sub := range []<-chan string{sub1, sub2, sub3} {
			select {
			case msg := <-sub:
				if msg != "broadcast" {
					t.Errorf("subscriber %d received %q, want %q", i, msg, "broadcast")
				}
			case <-time.After(1 * time.Second):
				t.Errorf("subscriber %d did not receive message within timeout", i)
			}
		}
	})

	t.Run("close signals end to subscribers", func(t *testing.T) {
		broker := NewBroker()
		if broker == nil {
			t.Fatal("NewBroker() returned nil")
		}

		sub := broker.Subscribe()
		broker.Publish("before close")
		<-sub // consume the message

		broker.Close()

		// After close, receiving from the channel should return zero value
		select {
		case _, ok := <-sub:
			if ok {
				// It's possible to receive buffered messages, but eventually ok should be false
				// Try again to get the close signal
				_, ok = <-sub
				if ok {
					t.Error("expected channel to be closed after broker.Close()")
				}
			}
			// ok is false, channel is closed — this is correct
		case <-time.After(1 * time.Second):
			t.Error("subscriber channel not closed within timeout")
		}
	})
}

// =============================================================================
// Lesson Function Tests
// =============================================================================

func TestLessonGoroutines(t *testing.T) {
	t.Run("DemoBasicGoroutine", func(t *testing.T) {
		result := DemoBasicGoroutine()
		if result != "hello from goroutine" {
			t.Errorf("got %q, want %q", result, "hello from goroutine")
		}
	})

	t.Run("DemoMultipleGoroutines", func(t *testing.T) {
		results := DemoMultipleGoroutines(5)
		if len(results) != 5 {
			t.Fatalf("expected 5 results, got %d", len(results))
		}
		// Check that all expected squares are present
		sort.Ints(results)
		expected := []int{0, 1, 4, 9, 16}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("sorted results[%d] = %d, want %d", i, results[i], v)
			}
		}
	})

	t.Run("DemoBufferedChannel", func(t *testing.T) {
		results := DemoBufferedChannel()
		expected := []int{10, 20, 30}
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("results[%d] = %d, want %d", i, results[i], v)
			}
		}
	})

	t.Run("DemoPipeline", func(t *testing.T) {
		results := DemoPipeline(4)
		expected := []int{0, 1, 4, 9}
		if len(results) != 4 {
			t.Fatalf("expected 4 results, got %d", len(results))
		}
		for i, v := range expected {
			if results[i] != v {
				t.Errorf("results[%d] = %d, want %d", i, results[i], v)
			}
		}
	})

	t.Run("DemoRangeOverChannel", func(t *testing.T) {
		result := DemoRangeOverChannel()
		expected := []string{"Go", "is", "concurrent"}
		if len(result) != 3 {
			t.Fatalf("expected 3 words, got %d", len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("result[%d] = %q, want %q", i, result[i], v)
			}
		}
	})

	t.Run("GenerateSequence", func(t *testing.T) {
		result := FormatSequence(1, 3)
		expected := []string{"#1", "#2", "#3"}
		if len(result) != 3 {
			t.Fatalf("expected 3 items, got %d", len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("result[%d] = %q, want %q", i, result[i], v)
			}
		}
	})
}
