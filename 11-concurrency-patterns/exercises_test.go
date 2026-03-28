package concurrency

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Exercise 1: WaitGroup — ParallelSquare
// =============================================================================

func TestParallelSquare(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "basic squares",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{1, 4, 9, 16, 25},
		},
		{
			name:  "includes zero and negatives",
			input: []int{-2, 0, 3},
			want:  []int{0, 4, 9},
		},
		{
			name:  "empty input",
			input: []int{},
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParallelSquare(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("ParallelSquare returned %d values, want %d", len(got), len(tt.want))
			}
			sort.Ints(got)
			sort.Ints(tt.want)
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("sorted result[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// =============================================================================
// Exercise 2: SafeMap
// =============================================================================

func TestSafeMap(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		m := NewSafeMap()
		if m == nil {
			t.Fatal("NewSafeMap() returned nil")
		}

		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		if v, ok := m.Get("a"); !ok || v != 1 {
			t.Errorf("Get(a) = (%d, %v), want (1, true)", v, ok)
		}
		if v, ok := m.Get("b"); !ok || v != 2 {
			t.Errorf("Get(b) = (%d, %v), want (2, true)", v, ok)
		}
		if _, ok := m.Get("missing"); ok {
			t.Error("Get(missing) should return false")
		}
		if m.Len() != 3 {
			t.Errorf("Len() = %d, want 3", m.Len())
		}
	})

	t.Run("delete", func(t *testing.T) {
		m := NewSafeMap()
		if m == nil {
			t.Fatal("NewSafeMap() returned nil")
		}

		m.Set("x", 10)
		m.Delete("x")
		if _, ok := m.Get("x"); ok {
			t.Error("Get(x) should return false after Delete")
		}
		if m.Len() != 0 {
			t.Errorf("Len() = %d, want 0 after Delete", m.Len())
		}
	})

	t.Run("keys", func(t *testing.T) {
		m := NewSafeMap()
		if m == nil {
			t.Fatal("NewSafeMap() returned nil")
		}

		m.Set("z", 26)
		m.Set("a", 1)
		m.Set("m", 13)

		keys := m.Keys()
		sort.Strings(keys)
		expected := []string{"a", "m", "z"}
		if len(keys) != len(expected) {
			t.Fatalf("Keys() returned %d keys, want %d", len(keys), len(expected))
		}
		for i := range keys {
			if keys[i] != expected[i] {
				t.Errorf("keys[%d] = %q, want %q", i, keys[i], expected[i])
			}
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		m := NewSafeMap()
		if m == nil {
			t.Fatal("NewSafeMap() returned nil")
		}

		var wg sync.WaitGroup
		// Concurrent writes
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				m.Set(fmt.Sprintf("key%d", n), n)
			}(i)
		}
		wg.Wait()

		if m.Len() != 100 {
			t.Errorf("after 100 concurrent Sets, Len() = %d, want 100", m.Len())
		}

		// Concurrent reads
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				v, ok := m.Get(fmt.Sprintf("key%d", n))
				if !ok || v != n {
					t.Errorf("Get(key%d) = (%d, %v), want (%d, true)", n, v, ok, n)
				}
			}(i)
		}
		wg.Wait()
	})
}

// =============================================================================
// Exercise 3: Multiplex
// =============================================================================

func TestMultiplex(t *testing.T) {
	t.Run("merges two channels", func(t *testing.T) {
		ch1 := make(chan int, 3)
		ch2 := make(chan int, 3)

		ch1 <- 1
		ch1 <- 2
		ch1 <- 3
		close(ch1)

		ch2 <- 10
		ch2 <- 20
		close(ch2)

		out := Multiplex(ch1, ch2)
		if out == nil {
			t.Fatal("Multiplex returned nil channel")
		}

		var results []int
		for v := range out {
			results = append(results, v)
		}

		sort.Ints(results)
		expected := []int{1, 2, 3, 10, 20}
		if len(results) != len(expected) {
			t.Fatalf("got %d values, want %d: %v", len(results), len(expected), results)
		}
		for i := range results {
			if results[i] != expected[i] {
				t.Errorf("sorted results[%d] = %d, want %d", i, results[i], expected[i])
			}
		}
	})

	t.Run("one empty channel", func(t *testing.T) {
		ch1 := make(chan int, 2)
		ch2 := make(chan int)

		ch1 <- 5
		ch1 <- 6
		close(ch1)
		close(ch2)

		out := Multiplex(ch1, ch2)
		if out == nil {
			t.Fatal("Multiplex returned nil channel")
		}

		var results []int
		for v := range out {
			results = append(results, v)
		}

		sort.Ints(results)
		expected := []int{5, 6}
		if len(results) != len(expected) {
			t.Fatalf("got %d values, want %d", len(results), len(expected))
		}
		for i := range results {
			if results[i] != expected[i] {
				t.Errorf("results[%d] = %d, want %d", i, results[i], expected[i])
			}
		}
	})

	t.Run("both empty channels", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)
		close(ch1)
		close(ch2)

		out := Multiplex(ch1, ch2)
		if out == nil {
			t.Fatal("Multiplex returned nil channel")
		}

		var results []int
		for v := range out {
			results = append(results, v)
		}
		if len(results) != 0 {
			t.Errorf("expected 0 values from empty channels, got %d", len(results))
		}
	})
}

// =============================================================================
// Exercise 4: Context Cancellation Chain
// =============================================================================

func TestSlowComputation(t *testing.T) {
	t.Run("completes before cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		result, err := SlowComputation(ctx, 10*time.Millisecond, "done")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "done" {
			t.Errorf("result = %q, want %q", result, "done")
		}
	})

	t.Run("cancelled before completion", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := SlowComputation(ctx, 1*time.Second, "done")
		if err == nil {
			t.Error("expected error from cancelled context")
		}
	})
}

func TestChainedComputation(t *testing.T) {
	t.Run("all steps complete", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		results, err := ChainedComputation(ctx, 10*time.Millisecond, 10*time.Millisecond, 10*time.Millisecond)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		expected := []string{"step1", "step2", "step3"}
		for i, want := range expected {
			if results[i] != want {
				t.Errorf("results[%d] = %q, want %q", i, results[i], want)
			}
		}
	})

	t.Run("cancelled mid-chain", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()

		results, err := ChainedComputation(ctx, 10*time.Millisecond, 10*time.Millisecond, 500*time.Millisecond)
		if err == nil {
			t.Error("expected error from cancelled context")
		}
		// Should have completed at least the first step
		if len(results) < 1 {
			t.Error("expected at least 1 completed step before cancellation")
		}
	})
}

// =============================================================================
// Exercise 5: Worker Pool
// =============================================================================

func TestWorkerPool(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		tasks      []Task
		want       map[int]int // taskID -> expected result
	}{
		{
			name:       "three workers five tasks",
			numWorkers: 3,
			tasks: []Task{
				{ID: 1, Value: 2},
				{ID: 2, Value: 3},
				{ID: 3, Value: 4},
				{ID: 4, Value: 5},
				{ID: 5, Value: 6},
			},
			want: map[int]int{1: 8, 2: 27, 3: 64, 4: 125, 5: 216},
		},
		{
			name:       "single worker",
			numWorkers: 1,
			tasks: []Task{
				{ID: 1, Value: 3},
				{ID: 2, Value: 4},
			},
			want: map[int]int{1: 27, 2: 64},
		},
		{
			name:       "more workers than tasks",
			numWorkers: 10,
			tasks: []Task{
				{ID: 1, Value: 2},
			},
			want: map[int]int{1: 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := WorkerPool(tt.numWorkers, tt.tasks)

			if len(results) != len(tt.want) {
				t.Fatalf("got %d results, want %d", len(results), len(tt.want))
			}

			for _, r := range results {
				expected, ok := tt.want[r.TaskID]
				if !ok {
					t.Errorf("unexpected task ID %d in results", r.TaskID)
					continue
				}
				if r.Result != expected {
					t.Errorf("task %d: result = %d, want %d (value^3)", r.TaskID, r.Result, expected)
				}
			}
		})
	}
}

// =============================================================================
// Exercise 6: Fan-Out/Fan-In
// =============================================================================

func TestFanOutFanIn(t *testing.T) {
	t.Run("double all values", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8}
		double := func(n int) int { return n * 2 }

		got := FanOutFanIn(input, 3, double)
		sort.Ints(got)

		want := []int{2, 4, 6, 8, 10, 12, 14, 16}
		if len(got) != len(want) {
			t.Fatalf("got %d values, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("sorted result[%d] = %d, want %d", i, got[i], want[i])
			}
		}
	})

	t.Run("single worker", func(t *testing.T) {
		input := []int{5, 10}
		square := func(n int) int { return n * n }

		got := FanOutFanIn(input, 1, square)
		sort.Ints(got)

		want := []int{25, 100}
		if len(got) != len(want) {
			t.Fatalf("got %d values, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("sorted result[%d] = %d, want %d", i, got[i], want[i])
			}
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got := FanOutFanIn([]int{}, 3, func(n int) int { return n })
		if len(got) != 0 {
			t.Errorf("expected empty result for empty input, got %v", got)
		}
	})
}

// =============================================================================
// Exercise 7: Token Bucket Rate Limiter
// =============================================================================

func TestTokenBucket(t *testing.T) {
	t.Run("starts full", func(t *testing.T) {
		tb := NewTokenBucket(3, 1*time.Hour) // Long refill so it doesn't interfere
		if tb == nil {
			t.Fatal("NewTokenBucket returned nil")
		}
		defer tb.Stop()

		// Should allow 'capacity' operations immediately
		for i := 0; i < 3; i++ {
			if !tb.Allow() {
				t.Errorf("Allow() call %d should succeed (bucket starts full)", i+1)
			}
		}

		// Next should fail (bucket empty, refill interval is very long)
		if tb.Allow() {
			t.Error("Allow() should fail when bucket is empty")
		}
	})

	t.Run("refills over time", func(t *testing.T) {
		tb := NewTokenBucket(2, 20*time.Millisecond)
		if tb == nil {
			t.Fatal("NewTokenBucket returned nil")
		}
		defer tb.Stop()

		// Drain the bucket
		for tb.Allow() {
			// drain
		}

		// Wait for a refill
		time.Sleep(50 * time.Millisecond)

		// Should have at least 1 token now
		if !tb.Allow() {
			t.Error("Allow() should succeed after refill interval")
		}
	})
}

// =============================================================================
// Exercise 8: Processing Pipeline
// =============================================================================

func TestProcessingPipeline(t *testing.T) {
	t.Run("filter-transform-format", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		isEven := func(n int) bool { return n%2 == 0 }
		double := func(n int) int { return n * 2 }
		format := func(n int) string { return fmt.Sprintf("val=%d", n) }

		results := ProcessingPipeline(numbers, isEven, double, format)

		expected := []PipelineResult{
			{Original: 2, Processed: "val=4"},
			{Original: 4, Processed: "val=8"},
			{Original: 6, Processed: "val=12"},
			{Original: 8, Processed: "val=16"},
			{Original: 10, Processed: "val=20"},
		}

		if len(results) != len(expected) {
			t.Fatalf("got %d results, want %d", len(results), len(expected))
		}

		for i := range results {
			if results[i].Original != expected[i].Original {
				t.Errorf("results[%d].Original = %d, want %d", i, results[i].Original, expected[i].Original)
			}
			if results[i].Processed != expected[i].Processed {
				t.Errorf("results[%d].Processed = %q, want %q", i, results[i].Processed, expected[i].Processed)
			}
		}
	})

	t.Run("no values pass filter", func(t *testing.T) {
		numbers := []int{1, 3, 5}

		isEven := func(n int) bool { return n%2 == 0 }
		identity := func(n int) int { return n }
		format := func(n int) string { return fmt.Sprintf("%d", n) }

		results := ProcessingPipeline(numbers, isEven, identity, format)
		if len(results) != 0 {
			t.Errorf("expected 0 results when nothing passes filter, got %d", len(results))
		}
	})

	t.Run("all values pass filter", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		passAll := func(n int) bool { return true }
		square := func(n int) int { return n * n }
		format := func(n int) string { return fmt.Sprintf("sq=%d", n) }

		results := ProcessingPipeline(numbers, passAll, square, format)

		if len(results) != 3 {
			t.Fatalf("got %d results, want 3", len(results))
		}
		if results[0].Processed != "sq=1" {
			t.Errorf("results[0].Processed = %q, want %q", results[0].Processed, "sq=1")
		}
		if results[1].Processed != "sq=4" {
			t.Errorf("results[1].Processed = %q, want %q", results[1].Processed, "sq=4")
		}
		if results[2].Processed != "sq=9" {
			t.Errorf("results[2].Processed = %q, want %q", results[2].Processed, "sq=9")
		}
	})
}

// =============================================================================
// Lesson Function Tests
// =============================================================================

func TestLessonPatterns(t *testing.T) {
	t.Run("DemoMutex", func(t *testing.T) {
		result := DemoMutex(10, 100)
		if result != 1000 {
			t.Errorf("DemoMutex(10, 100) = %d, want 1000", result)
		}
	})

	t.Run("Cache", func(t *testing.T) {
		c := NewCache()
		c.Set("key", "value")
		v, ok := c.Get("key")
		if !ok || v != "value" {
			t.Errorf("Cache.Get(key) = (%q, %v), want (value, true)", v, ok)
		}
	})

	t.Run("GetDB singleton", func(t *testing.T) {
		db1 := GetDB()
		db2 := GetDB()
		if db1 != db2 {
			t.Error("GetDB() should return the same instance")
		}
	})

	t.Run("DemoWorkerPool", func(t *testing.T) {
		jobs := []Job{{ID: 1, Input: 3}, {ID: 2, Input: 4}}
		results := DemoWorkerPool(2, jobs)
		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("DemoPipeline", func(t *testing.T) {
		results := DemoPipeline([]int{1, 2, 3, 4, 5})
		expected := []string{"even:2", "even:4"}
		if len(results) != len(expected) {
			t.Fatalf("expected %d results, got %d: %v", len(expected), len(results), results)
		}
		for i := range expected {
			if results[i] != expected[i] {
				t.Errorf("results[%d] = %q, want %q", i, results[i], expected[i])
			}
		}
	})
}
