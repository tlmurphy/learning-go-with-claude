package production

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 EXERCISES: Production Patterns
=============================================================================

 These exercises build real production-quality components. Run tests with:

   go test -v ./29-production-patterns/

 Each exercise is a building block that could be used in a real service.

=============================================================================
*/

// =========================================================================
// Exercise 1: Structured Logging Setup
// =========================================================================
//
// Create a function that sets up a structured logger with a custom handler.
// The logger should support different log levels based on a string parameter.

// NewLogger creates an slog.Logger that writes to the given handler.
// The level parameter should be one of: "debug", "info", "warn", "error".
// If the level string is unrecognized, default to "info".
//
// The attrs parameter provides key-value pairs that should be included
// in EVERY log message (e.g., "service", "my-app", "version", "1.0").
//
// Return the configured logger.
func NewLogger(handler slog.Handler, level string, attrs ...any) *slog.Logger {
	// YOUR CODE HERE
	// 1. Parse the level string into an slog.Level
	// 2. Create a logger from the handler
	// 3. Add the attrs using logger.With(attrs...)
	// 4. Return the logger
	_ = slog.LevelDebug
	_ = slog.LevelInfo
	_ = slog.LevelWarn
	_ = slog.LevelError
	_ = handler
	_ = level
	_ = attrs
	return slog.Default()
}

// ParseLogLevel converts a level string to slog.Level.
// Supports: "debug", "info", "warn", "error" (case-insensitive).
// Returns slog.LevelInfo for unrecognized values.
func ParseLogLevel(level string) slog.Level {
	// YOUR CODE HERE
	_ = strings.ToLower
	return slog.LevelInfo
}

// =========================================================================
// Exercise 2: Request-Scoped Logging Middleware
// =========================================================================
//
// Implement context-based logger passing for request-scoped logging.

type loggerKey struct{}

// ContextWithLogger stores a logger in the context.
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	// YOUR CODE HERE
	return ctx
}

// LoggerFromContext retrieves a logger from the context.
// If no logger is found, return slog.Default().
func LoggerFromContext(ctx context.Context) *slog.Logger {
	// YOUR CODE HERE
	return slog.Default()
}

// RequestContext creates a new context with a request-scoped logger.
// The logger should include the given requestID as a persistent field.
// This simulates what HTTP middleware would do.
func RequestContext(ctx context.Context, baseLogger *slog.Logger, requestID string) context.Context {
	// YOUR CODE HERE
	// 1. Create a new logger with the requestID: baseLogger.With("request_id", requestID)
	// 2. Store it in context using ContextWithLogger
	_ = baseLogger
	_ = requestID
	return ctx
}

// =========================================================================
// Exercise 3: Circuit Breaker
// =========================================================================
//
// Implement a circuit breaker that protects against cascading failures.

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu sync.Mutex

	// Configuration
	failureThreshold int           // failures before opening
	resetTimeout     time.Duration // how long to stay open before half-open

	// State
	state        CircuitState
	failures     int       // consecutive failure count
	lastFailTime time.Time // when the last failure occurred

	// For testing: allow injecting a clock
	now func() time.Time
}

// NewCircuitBreaker creates a circuit breaker.
//
//	failureThreshold: number of consecutive failures before opening
//	resetTimeout: how long to wait in Open state before trying Half-Open
func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	// YOUR CODE HERE
	return &CircuitBreaker{
		now: time.Now,
	}
}

// Execute runs the given function through the circuit breaker.
//
// Behavior:
//   - CLOSED: execute the function. If it fails, increment failure count.
//     If failures >= threshold, transition to OPEN.
//   - OPEN: check if resetTimeout has elapsed. If not, return ErrCircuitOpen
//     immediately. If yes, transition to HALF-OPEN and allow the call.
//   - HALF-OPEN: execute the function. If it succeeds, transition to CLOSED
//     and reset failures. If it fails, transition back to OPEN.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// YOUR CODE HERE
	_ = fn
	return nil
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Failures returns the current consecutive failure count.
func (cb *CircuitBreaker) Failures() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.failures
}

// =========================================================================
// Exercise 4: Retry with Exponential Backoff and Jitter
// =========================================================================
//
// Implement a retry function with configurable behavior.

// RetryConfig configures the retry behavior.
type RetryConfig struct {
	MaxRetries  int           // maximum number of retry attempts
	InitialWait time.Duration // initial backoff duration
	MaxWait     time.Duration // cap on backoff duration
	Multiplier  float64       // backoff multiplier (typically 2.0)
}

// DefaultRetryConfig returns sensible defaults for retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:  3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     5 * time.Second,
		Multiplier:  2.0,
	}
}

// RetryWithBackoff executes fn with retry logic. It retries on error
// with exponential backoff and jitter.
//
// - Calls fn up to config.MaxRetries+1 times (1 initial + MaxRetries retries)
// - Wait between retries: min(initialWait * multiplier^attempt + jitter, maxWait)
// - Jitter: random duration between 0 and current backoff
// - Respects context cancellation (stop retrying if ctx is done)
// - Returns the first nil error, or the last error if all attempts fail
//
// The sleepFn parameter replaces time.Sleep for testing. In production,
// pass time.Sleep. In tests, pass a no-op or a fake.
func RetryWithBackoff(ctx context.Context, config RetryConfig, sleepFn func(time.Duration), fn func() error) error {
	// YOUR CODE HERE
	_ = rand.Int63n
	_ = ctx
	_ = config
	_ = sleepFn
	_ = fn
	return nil
}

// =========================================================================
// Exercise 5: Health Check Aggregator
// =========================================================================
//
// Build a health check system that checks multiple dependencies.

// CheckResult represents the result of a single health check.
type CheckResult struct {
	Name    string
	Status  string // "ok", "degraded", "down"
	Message string // optional detail
}

// HealthChecker is a function that checks a dependency's health.
type HealthChecker func(ctx context.Context) CheckResult

// HealthAggregator checks multiple dependencies and reports overall status.
type HealthAggregator struct {
	checks map[string]HealthChecker
}

// NewHealthAggregator creates a new health aggregator.
func NewHealthAggregator() *HealthAggregator {
	// YOUR CODE HERE
	return &HealthAggregator{}
}

// AddCheck registers a named health check.
func (h *HealthAggregator) AddCheck(name string, checker HealthChecker) {
	// YOUR CODE HERE
}

// Check runs all registered health checks and returns the results.
//
// Overall status rules:
//   - If all checks are "ok", overall status is "ok"
//   - If any check is "degraded" (but none "down"), overall status is "degraded"
//   - If any check is "down", overall status is "down"
//
// Returns: overall status string, and slice of individual results.
func (h *HealthAggregator) Check(ctx context.Context) (string, []CheckResult) {
	// YOUR CODE HERE
	return "ok", nil
}

// =========================================================================
// Exercise 6: Simple Metrics Collector
// =========================================================================
//
// Implement basic counter and gauge metric types with a registry.

// Counter is a metric that only goes up.
type Counter struct {
	mu    sync.Mutex
	value int64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	// YOUR CODE HERE
}

// Add adds the given value to the counter. Value must be non-negative.
func (c *Counter) Add(v int64) {
	// YOUR CODE HERE
}

// Value returns the current counter value.
func (c *Counter) Value() int64 {
	// YOUR CODE HERE
	return 0
}

// Gauge is a metric that can go up and down.
type Gauge struct {
	mu    sync.Mutex
	value float64
}

// Set sets the gauge to the given value.
func (g *Gauge) Set(v float64) {
	// YOUR CODE HERE
}

// Inc increments the gauge by 1.
func (g *Gauge) Inc() {
	// YOUR CODE HERE
}

// Dec decrements the gauge by 1.
func (g *Gauge) Dec() {
	// YOUR CODE HERE
}

// Value returns the current gauge value.
func (g *Gauge) Value() float64 {
	// YOUR CODE HERE
	return 0
}

// MetricsRegistry holds named metrics.
type MetricsRegistry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
	gauges   map[string]*Gauge
}

// NewMetricsRegistry creates a new registry.
func NewMetricsRegistry() *MetricsRegistry {
	// YOUR CODE HERE
	return &MetricsRegistry{}
}

// Counter returns a named counter, creating it if needed.
// Subsequent calls with the same name return the same counter.
func (r *MetricsRegistry) Counter(name string) *Counter {
	// YOUR CODE HERE
	return &Counter{}
}

// Gauge returns a named gauge, creating it if needed.
func (r *MetricsRegistry) Gauge(name string) *Gauge {
	// YOUR CODE HERE
	return &Gauge{}
}

// =========================================================================
// Exercise 7: Graceful Degradation Wrapper
// =========================================================================
//
// Build a wrapper that tries a primary function and falls back to a
// secondary function if the primary fails.

// WithFallback calls primary. If it fails, calls fallback instead.
// Returns the result from whichever succeeds, or the fallback's error
// if both fail.
//
// The fallbackErr return value indicates whether the fallback was used.
// If primary succeeds, fallbackErr is nil.
// If primary fails and fallback succeeds, fallbackErr is the primary's error.
// If both fail, returns fallback's result and fallback's error.
func WithFallback[T any](
	primary func() (T, error),
	fallback func() (T, error),
) (result T, fallbackErr error, usedFallback bool) {
	// YOUR CODE HERE
	_ = primary
	_ = fallback
	var zero T
	return zero, nil, false
}

// CachedFallback creates a fallback function that returns a cached value.
// This is useful for graceful degradation: try the real service, fall back
// to a cached/default value.
func CachedFallback[T any](cachedValue T) func() (T, error) {
	// YOUR CODE HERE
	// Return a function that returns the cachedValue with nil error
	return func() (T, error) {
		var zero T
		return zero, nil
	}
}

// =========================================================================
// Exercise 8: Production-Ready Handler
// =========================================================================
//
// Wire together logging, metrics, circuit breaking, and error handling
// into a "production-ready" handler function.

// HandlerDeps holds all the production dependencies for a handler.
type HandlerDeps struct {
	Logger  *slog.Logger
	Metrics *MetricsRegistry
	Breaker *CircuitBreaker
}

// ProductionHandler demonstrates a handler with all production patterns.
//
// It should:
//  1. Increment a "requests_total" counter
//  2. Use the circuit breaker to call the serviceFn
//  3. If circuit breaker returns ErrCircuitOpen, increment "circuit_open_total"
//     counter and return the fallback result
//  4. If serviceFn succeeds, return the result
//  5. If serviceFn fails (but not circuit open), increment "errors_total"
//     counter and return the fallback result
//  6. Log each outcome (success, circuit open, error) using deps.Logger
//
// The fallbackFn provides the degraded result.
func ProductionHandler(
	ctx context.Context,
	deps HandlerDeps,
	serviceFn func() (string, error),
	fallbackFn func() (string, error),
) (string, error) {
	// YOUR CODE HERE
	_ = deps
	_ = serviceFn
	_ = fallbackFn
	_ = fmt.Sprintf
	_ = ErrCircuitOpen
	return "", nil
}
