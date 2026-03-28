// Package production covers patterns essential for running Go services in
// production: structured logging, circuit breakers, retries, health checks,
// metrics, graceful degradation, and observability.
package production

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"
)

/*
=============================================================================
 PRODUCTION PATTERNS IN GO
=============================================================================

Writing Go code that works on your laptop is step one. Making it survive
in production — with real traffic, flaky networks, overloaded databases,
and 3 AM incidents — is where these patterns come in.

This module covers the "operational" side of Go: the patterns that keep
your service running when things go wrong (and things WILL go wrong).

=============================================================================
 STRUCTURED LOGGING WITH log/slog
=============================================================================

Since Go 1.21, the standard library includes log/slog — a structured
logging package that produces machine-parseable log output (JSON, logfmt).

Why structured logging?
  - Grep works on plain text. But when you have millions of logs, you need
    to QUERY them: "show me all errors for user 42 in the last hour."
  - Structured logs are key-value pairs that log aggregation tools (Datadog,
    Grafana Loki, CloudWatch) can index and search.

Basic usage:
  slog.Info("user logged in",
      "user_id", 42,
      "ip", "192.168.1.1",
      "method", "oauth",
  )
  // Output (JSON handler):
  // {"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"user logged in",
  //  "user_id":42,"ip":"192.168.1.1","method":"oauth"}

Log levels:
  - slog.Debug: verbose info for development
  - slog.Info:  normal operations (user logged in, request served)
  - slog.Warn:  something unexpected but recoverable (retrying connection)
  - slog.Error: something broke (failed to save order, DB connection lost)

Production tip: log at Info level by default. Use Debug for development.
Never log at Error unless something is actually broken — error fatigue
is real, and if everything is an "error," nothing is.

=============================================================================
 REQUEST-SCOPED LOGGING
=============================================================================

In a web service, you want every log message for a request to include the
request ID. This makes tracing a single request through your logs trivial.

Pattern: create a logger with the request ID, store it in context.

  func middleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          reqID := uuid.New().String()
          logger := slog.With("request_id", reqID)
          ctx := ContextWithLogger(r.Context(), logger)
          next.ServeHTTP(w, r.WithContext(ctx))
      })
  }

  func handler(w http.ResponseWriter, r *http.Request) {
      logger := LoggerFromContext(r.Context())
      logger.Info("processing request", "path", r.URL.Path)
      // All logs from this request include request_id automatically
  }

=============================================================================
 METRICS PATTERNS
=============================================================================

Three types of metrics you need:

  Counters: things that only go up
    - requests_total, errors_total, orders_placed_total
    - "How many X happened?"

  Gauges: things that go up and down
    - active_connections, queue_depth, goroutine_count
    - "What is the current value of X?"

  Histograms: distributions of values
    - request_duration_seconds, response_size_bytes
    - "What does the distribution of X look like?"

You don't need a specific library to understand these concepts. Prometheus
is the most common choice in Go (github.com/prometheus/client_golang),
but the concepts are universal.

=============================================================================
 HEALTH CHECKS: LIVENESS VS READINESS
=============================================================================

In Kubernetes (and similar platforms), health checks tell the orchestrator
whether your service is healthy:

  Liveness: "Is the process alive and not deadlocked?"
    - Returns 200 if the process is running
    - If this fails, Kubernetes RESTARTS the container
    - Keep it simple: if the process can respond, it's alive

  Readiness: "Can this instance serve traffic?"
    - Returns 200 if ALL dependencies are healthy (DB, cache, etc.)
    - If this fails, Kubernetes stops sending traffic to this instance
    - Check each dependency: can you ping the DB? Is the cache reachable?

Common mistake: making liveness checks too complex. If your liveness check
queries the database and the DB is slow, Kubernetes restarts your service,
which makes the DB even more overloaded. Cascade failure.

=============================================================================
 CIRCUIT BREAKER PATTERN
=============================================================================

When a downstream service is failing, you don't want to keep hammering it
with requests. The circuit breaker "trips" after too many failures:

  Closed (normal): requests flow through.
    → After N consecutive failures, transition to Open.

  Open (broken): requests fail immediately without calling downstream.
    → After a timeout, transition to Half-Open.

  Half-Open (testing): one request is allowed through.
    → If it succeeds, transition to Closed.
    → If it fails, transition back to Open.

This protects:
  1. The downstream service (doesn't get overwhelmed)
  2. Your service (doesn't waste time on requests that will fail)
  3. The user (gets a fast error instead of a timeout)

=============================================================================
 RETRY WITH EXPONENTIAL BACKOFF AND JITTER
=============================================================================

When a transient error occurs (network blip, temporary overload), retry.
But retry SMART:

  - Exponential backoff: wait 1s, 2s, 4s, 8s, ... between retries
  - Jitter: add randomness to prevent the "thundering herd" problem
    (if 1000 clients all retry at exactly the same time, you've just
     created the same overload that caused the failure)

  func retryWithBackoff(ctx context.Context, fn func() error) error {
      backoff := time.Second
      for attempt := 0; attempt < maxRetries; attempt++ {
          err := fn()
          if err == nil {
              return nil
          }
          jitter := time.Duration(rand.Int63n(int64(backoff)))
          select {
          case <-time.After(backoff + jitter):
          case <-ctx.Done():
              return ctx.Err()
          }
          backoff *= 2
      }
      return errors.New("max retries exceeded")
  }

=============================================================================
 GRACEFUL DEGRADATION
=============================================================================

When a dependency is down, don't crash. Degrade gracefully:

  - Cache is down? Serve from the database (slower but works).
  - Recommendation service is down? Show popular items instead.
  - Payment processor is down? Queue the payment for later.

The key insight: partial functionality is better than no functionality.

  func GetRecommendations(ctx context.Context, userID int) ([]Product, error) {
      products, err := recommender.GetPersonalized(ctx, userID)
      if err != nil {
          logger.Warn("recommender down, falling back to popular items",
              "error", err)
          return GetPopularProducts(ctx) // degraded but functional
      }
      return products, nil
  }

=============================================================================
 OBSERVABILITY: THE THREE PILLARS
=============================================================================

  1. Logs: What happened? (discrete events)
  2. Metrics: How is it performing? (aggregated measurements)
  3. Traces: How does a request flow? (distributed call chains)

Logs tell you WHAT happened. Metrics tell you HOW MUCH. Traces tell you
WHERE (across services).

For Go services:
  - Logs: log/slog (stdlib)
  - Metrics: Prometheus client or OpenTelemetry
  - Traces: OpenTelemetry (Go 1.25 adds runtime/trace.FlightRecorder for
    continuous trace recording — keep a rolling buffer and dump on demand)

=============================================================================
 PANIC RECOVERY
=============================================================================

In production, a panic in one goroutine crashes the ENTIRE process. Use
a recovery middleware at the top level:

  func recoveryMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          defer func() {
              if err := recover(); err != nil {
                  slog.Error("panic recovered",
                      "error", err,
                      "stack", string(debug.Stack()),
                  )
                  http.Error(w, "Internal Server Error", 500)
              }
          }()
          next.ServeHTTP(w, r)
      })
  }

Never silently swallow panics. Log the stack trace — you need it to fix
the bug.

=============================================================================
 CONFIGURATION HOT-RELOADING
=============================================================================

Some configuration changes shouldn't require a restart:
  - Feature flags
  - Log level
  - Rate limits
  - A/B test percentages

Pattern: use atomic values or RWMutex-protected config:

  type Config struct {
      mu       sync.RWMutex
      logLevel slog.Level
      features map[string]bool
  }

  func (c *Config) LogLevel() slog.Level {
      c.mu.RLock()
      defer c.mu.RUnlock()
      return c.logLevel
  }

  func (c *Config) Reload(newLevel slog.Level, features map[string]bool) {
      c.mu.Lock()
      defer c.mu.Unlock()
      c.logLevel = newLevel
      c.features = features
  }

=============================================================================
*/

// --- Demo: Structured Logging ---

// DemoStructuredLogging shows slog basics.
func DemoStructuredLogging() {
	// JSON handler for production
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(jsonHandler)

	logger.Info("server starting",
		"host", "0.0.0.0",
		"port", 8080,
		"version", "1.2.3",
	)

	logger.Warn("high memory usage",
		"current_mb", 450,
		"limit_mb", 512,
	)

	// Logger with persistent fields (great for request-scoped logging)
	reqLogger := logger.With("request_id", "req-abc-123", "user_id", 42)
	reqLogger.Info("processing request")
	reqLogger.Info("request complete", "duration_ms", 150)
}

// --- Demo: Circuit Breaker Concept ---

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	StateClosed   CircuitState = iota // Normal operation
	StateOpen                         // Failing fast
	StateHalfOpen                     // Testing recovery
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

// DemoCircuitBreakerConcept explains the state machine.
func DemoCircuitBreakerConcept() {
	fmt.Println("=== Circuit Breaker States ===")
	fmt.Println("CLOSED → requests flow through normally")
	fmt.Println("  → after N failures → OPEN")
	fmt.Println("OPEN → requests fail immediately (fast failure)")
	fmt.Println("  → after timeout → HALF-OPEN")
	fmt.Println("HALF-OPEN → one test request allowed")
	fmt.Println("  → if success → CLOSED")
	fmt.Println("  → if failure → OPEN")
}

// --- Demo: Retry with Backoff ---

// DemoRetryBackoff shows exponential backoff with jitter.
func DemoRetryBackoff() {
	fmt.Println("=== Retry Backoff Schedule ===")
	backoff := 100 * time.Millisecond
	for attempt := 1; attempt <= 5; attempt++ {
		jitter := time.Duration(rand.Int63n(int64(backoff)))
		wait := backoff + jitter
		fmt.Printf("  Attempt %d: wait %v (base %v + jitter %v)\n",
			attempt, wait, backoff, jitter)
		backoff *= 2
	}
}

// --- Demo: Health Check ---

// HealthStatus represents the result of a health check.
type HealthStatus struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// DemoHealthCheck shows a basic health check aggregator.
func DemoHealthCheck() {
	status := HealthStatus{
		Status: "ok",
		Checks: map[string]string{
			"database": "ok",
			"cache":    "ok",
			"queue":    "degraded",
		},
	}

	// If any check is not "ok", overall status is degraded
	for name, s := range status.Checks {
		if s != "ok" {
			status.Status = "degraded"
			fmt.Printf("  Warning: %s is %s\n", name, s)
		}
	}
	fmt.Printf("Overall status: %s\n", status.Status)
}

// --- Demo: Panic Recovery ---

// DemoSafeCall shows how to call a function that might panic, safely.
func DemoSafeCall() {
	fmt.Println("=== Panic Recovery ===")

	err := safeCall(func() error {
		panic("something went terribly wrong")
	})
	fmt.Println("Recovered from panic:", err)

	err = safeCall(func() error {
		return nil // no panic
	})
	fmt.Println("Normal function:", err)
}

func safeCall(fn func() error) (retErr error) {
	defer func() {
		if r := recover(); r != nil {
			retErr = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	return fn()
}

// --- Demo: Hot-Reloadable Config ---

// HotConfig demonstrates thread-safe configuration hot-reloading.
type HotConfig struct {
	mu       sync.RWMutex
	logLevel string
	features map[string]bool
}

// NewHotConfig creates a config with defaults.
func NewHotConfig() *HotConfig {
	return &HotConfig{
		logLevel: "info",
		features: map[string]bool{},
	}
}

// LogLevel returns the current log level (thread-safe).
func (c *HotConfig) LogLevel() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logLevel
}

// IsFeatureEnabled checks if a feature flag is on (thread-safe).
func (c *HotConfig) IsFeatureEnabled(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.features[name]
}

// Reload updates config atomically (thread-safe).
func (c *HotConfig) Reload(logLevel string, features map[string]bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logLevel = logLevel
	c.features = features
}

// DemoHotConfig shows config reload in action.
func DemoHotConfig() {
	fmt.Println("=== Hot Config Reload ===")
	cfg := NewHotConfig()
	fmt.Println("Initial log level:", cfg.LogLevel())
	fmt.Println("Dark mode enabled:", cfg.IsFeatureEnabled("dark_mode"))

	// Simulate config reload (e.g., from file watcher or API call)
	cfg.Reload("debug", map[string]bool{"dark_mode": true})
	fmt.Println("After reload log level:", cfg.LogLevel())
	fmt.Println("Dark mode enabled:", cfg.IsFeatureEnabled("dark_mode"))
}

// Ensure imports are used.
var _ = errors.New
var _ = context.Background
