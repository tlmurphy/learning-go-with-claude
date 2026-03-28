package production

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Exercise 1: Structured Logging Setup
// ---------------------------------------------------------------------------

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"unknown", slog.LevelInfo},
		{"", slog.LevelInfo},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseLogLevel(tt.input)
			if got != tt.want {
				t.Errorf("ParseLogLevel(%q) = %v, want %v. Use strings.ToLower and a switch statement.",
					tt.input, got, tt.want)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	logger := NewLogger(handler, "info", "service", "test-app")
	if logger == nil {
		t.Fatal("NewLogger returned nil.")
	}

	logger.Info("test message", "key", "value")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Error("Logger output should contain the message. Create the logger from the provided handler.")
	}
	if !strings.Contains(output, "test-app") {
		t.Error("Logger output should contain the persistent attrs. Use logger.With(attrs...) to add them.")
	}
}

// ---------------------------------------------------------------------------
// Exercise 2: Request-Scoped Logging
// ---------------------------------------------------------------------------

func TestContextLogger(t *testing.T) {
	t.Run("store and retrieve", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, nil)
		logger := slog.New(handler)

		ctx := ContextWithLogger(context.Background(), logger)
		retrieved := LoggerFromContext(ctx)
		if retrieved == nil {
			t.Fatal("LoggerFromContext returned nil.")
		}
		retrieved.Info("from context")
		if !strings.Contains(buf.String(), "from context") {
			t.Error("Retrieved logger should be the one we stored. Use context.WithValue and type assert on retrieval.")
		}
	})

	t.Run("default when missing", func(t *testing.T) {
		logger := LoggerFromContext(context.Background())
		if logger == nil {
			t.Error("LoggerFromContext should return slog.Default() when no logger in context.")
		}
	})
}

func TestRequestContext(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	baseLogger := slog.New(handler)

	ctx := RequestContext(context.Background(), baseLogger, "req-xyz-789")
	logger := LoggerFromContext(ctx)
	logger.Info("processing")

	output := buf.String()
	if !strings.Contains(output, "req-xyz-789") {
		t.Error("Request logger should include request_id. Use baseLogger.With(\"request_id\", requestID).")
	}
	if !strings.Contains(output, "processing") {
		t.Error("Request logger should produce output.")
	}
}

// ---------------------------------------------------------------------------
// Exercise 3: Circuit Breaker
// ---------------------------------------------------------------------------

func TestCircuitBreaker(t *testing.T) {
	t.Run("starts closed", func(t *testing.T) {
		cb := NewCircuitBreaker(3, time.Second)
		if cb == nil {
			t.Fatal("NewCircuitBreaker returned nil. Initialize all fields.")
		}
		if cb.State() != StateClosed {
			t.Errorf("Initial state = %v, want CLOSED.", cb.State())
		}
	})

	t.Run("success keeps closed", func(t *testing.T) {
		cb := NewCircuitBreaker(3, time.Second)
		err := cb.Execute(func() error { return nil })
		if err != nil {
			t.Errorf("Execute returned error: %v. Successful calls should pass through.", err)
		}
		if cb.State() != StateClosed {
			t.Errorf("State = %v, want CLOSED after success.", cb.State())
		}
	})

	t.Run("opens after threshold failures", func(t *testing.T) {
		cb := NewCircuitBreaker(3, time.Second)
		fail := func() error { return errors.New("fail") }

		for i := 0; i < 3; i++ {
			_ = cb.Execute(fail)
		}

		if cb.State() != StateOpen {
			t.Errorf("State = %v, want OPEN after %d failures. Transition to Open when failures >= threshold.",
				cb.State(), 3)
		}
	})

	t.Run("open returns ErrCircuitOpen", func(t *testing.T) {
		cb := NewCircuitBreaker(1, time.Hour) // 1 failure opens, long timeout
		_ = cb.Execute(func() error { return errors.New("fail") })

		err := cb.Execute(func() error { return nil })
		if !errors.Is(err, ErrCircuitOpen) {
			t.Errorf("Execute when open returned %v, want ErrCircuitOpen.", err)
		}
	})

	t.Run("half-open after timeout", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 50*time.Millisecond)
		_ = cb.Execute(func() error { return errors.New("fail") })

		// Wait for reset timeout
		time.Sleep(100 * time.Millisecond)

		// The next call should be allowed (half-open)
		called := false
		_ = cb.Execute(func() error {
			called = true
			return nil
		})

		if !called {
			t.Error("After reset timeout, circuit should allow a test call (half-open). Check if resetTimeout has elapsed.")
		}
		if cb.State() != StateClosed {
			t.Errorf("State = %v, want CLOSED after successful half-open test.", cb.State())
		}
	})

	t.Run("half-open failure returns to open", func(t *testing.T) {
		cb := NewCircuitBreaker(1, 50*time.Millisecond)
		_ = cb.Execute(func() error { return errors.New("fail") })

		time.Sleep(100 * time.Millisecond)

		_ = cb.Execute(func() error { return errors.New("still failing") })

		if cb.State() != StateOpen {
			t.Errorf("State = %v, want OPEN after half-open failure.", cb.State())
		}
	})

	t.Run("success resets failure count", func(t *testing.T) {
		cb := NewCircuitBreaker(3, time.Second)
		fail := func() error { return errors.New("fail") }

		// 2 failures (not enough to open)
		_ = cb.Execute(fail)
		_ = cb.Execute(fail)
		// 1 success (should reset count)
		_ = cb.Execute(func() error { return nil })
		// 2 more failures (should NOT open, because count was reset)
		_ = cb.Execute(fail)
		_ = cb.Execute(fail)

		if cb.State() != StateClosed {
			t.Error("State should be CLOSED. A success should reset the failure count to 0.")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 4: Retry with Backoff
// ---------------------------------------------------------------------------

func TestRetryWithBackoff(t *testing.T) {
	t.Run("succeeds first try", func(t *testing.T) {
		config := DefaultRetryConfig()
		calls := 0
		err := RetryWithBackoff(context.Background(), config, func(time.Duration) {}, func() error {
			calls++
			return nil
		})
		if err != nil {
			t.Errorf("Expected nil error, got: %v.", err)
		}
		if calls != 1 {
			t.Errorf("Function called %d times, want 1. Don't retry on success.", calls)
		}
	})

	t.Run("succeeds on retry", func(t *testing.T) {
		config := DefaultRetryConfig()
		calls := 0
		err := RetryWithBackoff(context.Background(), config, func(time.Duration) {}, func() error {
			calls++
			if calls < 3 {
				return errors.New("transient error")
			}
			return nil
		})
		if err != nil {
			t.Errorf("Expected nil error after retry, got: %v.", err)
		}
		if calls != 3 {
			t.Errorf("Function called %d times, want 3.", calls)
		}
	})

	t.Run("exhausts retries", func(t *testing.T) {
		config := RetryConfig{MaxRetries: 2, InitialWait: time.Millisecond, MaxWait: time.Second, Multiplier: 2}
		calls := 0
		err := RetryWithBackoff(context.Background(), config, func(time.Duration) {}, func() error {
			calls++
			return errors.New("permanent error")
		})
		if err == nil {
			t.Error("Expected error after exhausting retries.")
		}
		// 1 initial + 2 retries = 3 calls
		if calls != 3 {
			t.Errorf("Function called %d times, want 3 (1 initial + 2 retries).", calls)
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		config := RetryConfig{MaxRetries: 10, InitialWait: time.Hour, MaxWait: time.Hour, Multiplier: 2}
		ctx, cancel := context.WithCancel(context.Background())

		calls := 0
		sleepFn := func(d time.Duration) {
			cancel() // cancel after first sleep
		}

		err := RetryWithBackoff(ctx, config, sleepFn, func() error {
			calls++
			return errors.New("fail")
		})

		if err == nil {
			t.Error("Expected error when context is cancelled.")
		}
		if calls > 2 {
			t.Errorf("Should stop retrying when context is cancelled. Called %d times.", calls)
		}
	})

	t.Run("backoff increases", func(t *testing.T) {
		config := RetryConfig{MaxRetries: 3, InitialWait: 100 * time.Millisecond, MaxWait: 10 * time.Second, Multiplier: 2}
		var sleepDurations []time.Duration
		sleepFn := func(d time.Duration) {
			sleepDurations = append(sleepDurations, d)
		}

		_ = RetryWithBackoff(context.Background(), config, sleepFn, func() error {
			return errors.New("fail")
		})

		if len(sleepDurations) < 2 {
			t.Fatalf("Expected at least 2 sleep calls, got %d.", len(sleepDurations))
		}
		// Second sleep should be longer than first (exponential backoff)
		// Account for jitter by checking it's at least somewhat longer
		if sleepDurations[1] <= sleepDurations[0]/2 {
			t.Errorf("Backoff should increase. Got durations: %v.", sleepDurations)
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 5: Health Check Aggregator
// ---------------------------------------------------------------------------

func TestHealthAggregator(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		agg := NewHealthAggregator()
		if agg == nil {
			t.Fatal("NewHealthAggregator returned nil. Initialize the checks map.")
		}
		agg.AddCheck("db", func(_ context.Context) CheckResult {
			return CheckResult{Name: "db", Status: "ok"}
		})
		agg.AddCheck("cache", func(_ context.Context) CheckResult {
			return CheckResult{Name: "cache", Status: "ok"}
		})

		status, results := agg.Check(context.Background())
		if status != "ok" {
			t.Errorf("Overall status = %q, want %q. All checks passed.", status, "ok")
		}
		if len(results) != 2 {
			t.Errorf("Got %d results, want 2.", len(results))
		}
	})

	t.Run("one degraded", func(t *testing.T) {
		agg := NewHealthAggregator()
		agg.AddCheck("db", func(_ context.Context) CheckResult {
			return CheckResult{Name: "db", Status: "ok"}
		})
		agg.AddCheck("cache", func(_ context.Context) CheckResult {
			return CheckResult{Name: "cache", Status: "degraded", Message: "high latency"}
		})

		status, _ := agg.Check(context.Background())
		if status != "degraded" {
			t.Errorf("Overall status = %q, want %q. One check is degraded.", status, "degraded")
		}
	})

	t.Run("one down", func(t *testing.T) {
		agg := NewHealthAggregator()
		agg.AddCheck("db", func(_ context.Context) CheckResult {
			return CheckResult{Name: "db", Status: "down", Message: "connection refused"}
		})
		agg.AddCheck("cache", func(_ context.Context) CheckResult {
			return CheckResult{Name: "cache", Status: "ok"}
		})

		status, _ := agg.Check(context.Background())
		if status != "down" {
			t.Errorf("Overall status = %q, want %q. A critical dependency is down.", status, "down")
		}
	})

	t.Run("no checks", func(t *testing.T) {
		agg := NewHealthAggregator()
		status, results := agg.Check(context.Background())
		if status != "ok" {
			t.Errorf("Status with no checks = %q, want %q.", status, "ok")
		}
		if len(results) != 0 {
			t.Errorf("Results with no checks = %d, want 0.", len(results))
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 6: Metrics Collector
// ---------------------------------------------------------------------------

func TestCounter(t *testing.T) {
	c := &Counter{}

	t.Run("starts at zero", func(t *testing.T) {
		if c.Value() != 0 {
			t.Errorf("Initial value = %d, want 0.", c.Value())
		}
	})

	t.Run("inc", func(t *testing.T) {
		c.Inc()
		if c.Value() != 1 {
			t.Errorf("After Inc, value = %d, want 1.", c.Value())
		}
	})

	t.Run("add", func(t *testing.T) {
		c.Add(5)
		if c.Value() != 6 {
			t.Errorf("After Add(5), value = %d, want 6.", c.Value())
		}
	})

	t.Run("concurrent safety", func(t *testing.T) {
		counter := &Counter{}
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				counter.Inc()
			}()
		}
		wg.Wait()
		if counter.Value() != 1000 {
			t.Errorf("After 1000 concurrent increments, value = %d, want 1000. Use a mutex.", counter.Value())
		}
	})
}

func TestGauge(t *testing.T) {
	g := &Gauge{}

	t.Run("set", func(t *testing.T) {
		g.Set(42.5)
		if g.Value() != 42.5 {
			t.Errorf("After Set(42.5), value = %f, want 42.5.", g.Value())
		}
	})

	t.Run("inc and dec", func(t *testing.T) {
		g.Set(10)
		g.Inc()
		if g.Value() != 11 {
			t.Errorf("After Set(10) + Inc, value = %f, want 11.", g.Value())
		}
		g.Dec()
		g.Dec()
		if g.Value() != 9 {
			t.Errorf("After two Dec, value = %f, want 9.", g.Value())
		}
	})
}

func TestMetricsRegistry(t *testing.T) {
	reg := NewMetricsRegistry()
	if reg == nil {
		t.Fatal("NewMetricsRegistry returned nil. Initialize the maps.")
	}

	t.Run("counter", func(t *testing.T) {
		c1 := reg.Counter("requests_total")
		c1.Inc()
		c2 := reg.Counter("requests_total") // same name → same counter
		if c2.Value() != 1 {
			t.Errorf("Same-name counter should return the same instance. Got value %d, want 1.", c2.Value())
		}
	})

	t.Run("gauge", func(t *testing.T) {
		g1 := reg.Gauge("active_connections")
		g1.Set(5)
		g2 := reg.Gauge("active_connections") // same name → same gauge
		if g2.Value() != 5 {
			t.Errorf("Same-name gauge should return the same instance. Got value %f, want 5.", g2.Value())
		}
	})

	t.Run("different names", func(t *testing.T) {
		c1 := reg.Counter("a")
		c2 := reg.Counter("b")
		c1.Inc()
		if c2.Value() != 0 {
			t.Error("Different-name counters should be independent.")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 7: Graceful Degradation
// ---------------------------------------------------------------------------

func TestWithFallback(t *testing.T) {
	t.Run("primary succeeds", func(t *testing.T) {
		result, fallbackErr, used := WithFallback(
			func() (string, error) { return "primary result", nil },
			func() (string, error) { return "fallback result", nil },
		)
		if result != "primary result" {
			t.Errorf("Result = %q, want %q. Return primary result on success.", result, "primary result")
		}
		if fallbackErr != nil {
			t.Errorf("FallbackErr should be nil when primary succeeds, got: %v.", fallbackErr)
		}
		if used {
			t.Error("usedFallback should be false when primary succeeds.")
		}
	})

	t.Run("primary fails, fallback succeeds", func(t *testing.T) {
		primaryErr := errors.New("primary failed")
		result, fallbackErr, used := WithFallback(
			func() (string, error) { return "", primaryErr },
			func() (string, error) { return "fallback result", nil },
		)
		if result != "fallback result" {
			t.Errorf("Result = %q, want %q. Return fallback result when primary fails.", result, "fallback result")
		}
		if fallbackErr != primaryErr {
			t.Errorf("FallbackErr should be the primary error, got: %v.", fallbackErr)
		}
		if !used {
			t.Error("usedFallback should be true when fallback is used.")
		}
	})

	t.Run("both fail", func(t *testing.T) {
		fallbackError := errors.New("fallback also failed")
		_, fallbackErr, used := WithFallback(
			func() (string, error) { return "", errors.New("primary failed") },
			func() (string, error) { return "", fallbackError },
		)
		if fallbackErr != fallbackError {
			t.Errorf("When both fail, should return fallback's error. Got: %v.", fallbackErr)
		}
		if !used {
			t.Error("usedFallback should be true when fallback is attempted.")
		}
	})
}

func TestCachedFallback(t *testing.T) {
	fallback := CachedFallback("cached value")
	result, err := fallback()
	if err != nil {
		t.Errorf("CachedFallback should return nil error, got: %v.", err)
	}
	if result != "cached value" {
		t.Errorf("CachedFallback result = %q, want %q.", result, "cached value")
	}
}

func TestWithFallbackInt(t *testing.T) {
	// Test with a different type to verify generics work
	result, _, _ := WithFallback(
		func() (int, error) { return 0, errors.New("fail") },
		func() (int, error) { return 42, nil },
	)
	if result != 42 {
		t.Errorf("WithFallback[int] result = %d, want 42.", result)
	}
}

// ---------------------------------------------------------------------------
// Exercise 8: Production Handler
// ---------------------------------------------------------------------------

func TestProductionHandler(t *testing.T) {
	makeLogger := func() (*slog.Logger, *bytes.Buffer) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		return slog.New(handler), &buf
	}

	t.Run("success path", func(t *testing.T) {
		logger, _ := makeLogger()
		metrics := NewMetricsRegistry()
		cb := NewCircuitBreaker(5, time.Second)

		deps := HandlerDeps{Logger: logger, Metrics: metrics, Breaker: cb}
		result, err := ProductionHandler(
			context.Background(), deps,
			func() (string, error) { return "ok response", nil },
			func() (string, error) { return "fallback", nil },
		)
		if err != nil {
			t.Fatalf("Unexpected error: %v.", err)
		}
		if result != "ok response" {
			t.Errorf("Result = %q, want %q.", result, "ok response")
		}

		reqCount := metrics.Counter("requests_total").Value()
		if reqCount != 1 {
			t.Errorf("requests_total = %d, want 1. Increment the counter on each call.", reqCount)
		}
	})

	t.Run("service error falls back", func(t *testing.T) {
		logger, _ := makeLogger()
		metrics := NewMetricsRegistry()
		cb := NewCircuitBreaker(5, time.Second)

		deps := HandlerDeps{Logger: logger, Metrics: metrics, Breaker: cb}
		result, err := ProductionHandler(
			context.Background(), deps,
			func() (string, error) { return "", errors.New("service down") },
			func() (string, error) { return "fallback result", nil },
		)
		if err != nil {
			t.Fatalf("Should return fallback result, not error: %v.", err)
		}
		if result != "fallback result" {
			t.Errorf("Result = %q, want %q. Use fallback when service fails.", result, "fallback result")
		}

		errCount := metrics.Counter("errors_total").Value()
		if errCount != 1 {
			t.Errorf("errors_total = %d, want 1. Increment on service error.", errCount)
		}
	})

	t.Run("circuit open uses fallback", func(t *testing.T) {
		logger, _ := makeLogger()
		metrics := NewMetricsRegistry()
		cb := NewCircuitBreaker(1, time.Hour) // opens after 1 failure

		deps := HandlerDeps{Logger: logger, Metrics: metrics, Breaker: cb}

		// Trip the circuit
		_, _ = ProductionHandler(
			context.Background(), deps,
			func() (string, error) { return "", errors.New("fail") },
			func() (string, error) { return "fallback", nil },
		)

		// Now the circuit should be open
		result, err := ProductionHandler(
			context.Background(), deps,
			func() (string, error) { return "should not be called", nil },
			func() (string, error) { return "circuit fallback", nil },
		)
		if err != nil {
			t.Fatalf("Should return fallback when circuit is open: %v.", err)
		}
		if result != "circuit fallback" {
			t.Errorf("Result = %q, want %q.", result, "circuit fallback")
		}

		circuitCount := metrics.Counter("circuit_open_total").Value()
		if circuitCount < 1 {
			t.Error("circuit_open_total should be incremented when circuit is open.")
		}
	})
}

// Ensure imports are used.
var _ = strings.Contains
