package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

/*
=============================================================================
 EXERCISES: Graceful Shutdown
=============================================================================

 These exercises teach the patterns for graceful shutdown of web services.
 They're designed to work without running an actual long-lived server —
 we simulate the lifecycle in tests.

 Run the tests with:

   make test 21

 Tip: Run a single test at a time while working:

   go test -v -run TestSignalHandler ./21-graceful-shutdown/

=============================================================================
*/

// =========================================================================
// Exercise 1: Signal Handler
// =========================================================================

// SignalHandler listens for OS signals and calls a shutdown function.
//
// It should:
//   - Listen for the given signals (typically SIGINT, SIGTERM)
//   - When a signal is received, call the onShutdown function
//   - Return a cancel function that stops listening for signals
//
// The returned cancel function should clean up the signal channel.
//
// Hint: Use signal.Notify and a goroutine that reads from the channel.
func SignalHandler(onShutdown func(os.Signal), signals ...os.Signal) (cancel func()) {
	// YOUR CODE HERE
	return func() {}
}

// =========================================================================
// Exercise 2: Graceful HTTP Server with Timeout
// =========================================================================

// ManagedServer wraps an http.Server with lifecycle management.
type ManagedServer struct {
	server          *http.Server
	shutdownTimeout time.Duration
	running         atomic.Bool
}

// NewManagedServer creates a new ManagedServer.
//
// Parameters:
//   - addr: the address to listen on (e.g., ":8080")
//   - handler: the HTTP handler
//   - shutdownTimeout: max time to wait for in-flight requests during shutdown
func NewManagedServer(addr string, handler http.Handler, shutdownTimeout time.Duration) *ManagedServer {
	// YOUR CODE HERE
	return nil
}

// Start begins listening for HTTP requests in a background goroutine.
// Returns a channel that will be closed when the server has fully stopped.
//
// The server goroutine should:
//  1. Set running to true
//  2. Call ListenAndServe
//  3. When ListenAndServe returns, set running to false
//  4. Close the done channel
//
// Returns the done channel and any immediate error (e.g., nil listener).
func (s *ManagedServer) Start() (<-chan struct{}, error) {
	// YOUR CODE HERE
	return nil, nil
}

// Stop gracefully shuts down the server.
// It should use the configured shutdownTimeout as a deadline.
// Returns any error from the shutdown process.
func (s *ManagedServer) Stop() error {
	// YOUR CODE HERE
	return nil
}

// IsRunning returns whether the server is currently running.
func (s *ManagedServer) IsRunning() bool {
	// YOUR CODE HERE
	return false
}

// =========================================================================
// Exercise 3: Health Check Handler
// =========================================================================

// HealthStatus tracks the health state of a service for both
// liveness and readiness probes.
type HealthStatus struct {
	alive atomic.Bool
	ready atomic.Bool
}

// NewHealthStatus creates a new HealthStatus.
// Initially: alive=true (the process is running), ready=false (still starting).
func NewHealthStatus() *HealthStatus {
	// YOUR CODE HERE
	return nil
}

// SetAlive sets the liveness state.
func (h *HealthStatus) SetAlive(alive bool) {
	// YOUR CODE HERE
}

// SetReady sets the readiness state.
func (h *HealthStatus) SetReady(ready bool) {
	// YOUR CODE HERE
}

// HealthHandler returns an http.HandlerFunc for the /health endpoint.
// Returns 200 "alive" if alive, 503 "not alive" otherwise.
func (h *HealthStatus) HealthHandler() http.HandlerFunc {
	// YOUR CODE HERE
	return nil
}

// ReadyHandler returns an http.HandlerFunc for the /ready endpoint.
// Returns 200 "ready" if ready, 503 "not ready" otherwise.
func (h *HealthStatus) ReadyHandler() http.HandlerFunc {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 4: Shutdown Coordinator
// =========================================================================

// ShutdownFunc is a named function that runs during shutdown.
type ShutdownFunc struct {
	Name string
	Fn   func(ctx context.Context) error
}

// ShutdownCoordinator manages an ordered shutdown of multiple components.
// Components are shut down in the reverse order they were registered
// (like a stack — last registered, first shut down).
type ShutdownCoordinator struct {
	funcs []ShutdownFunc
	mu    sync.Mutex
}

// NewShutdownCoordinator creates a new coordinator.
func NewShutdownCoordinator() *ShutdownCoordinator {
	// YOUR CODE HERE
	return nil
}

// Register adds a shutdown function. Functions are called in reverse
// order during shutdown (LIFO — last in, first out).
func (c *ShutdownCoordinator) Register(sf ShutdownFunc) {
	// YOUR CODE HERE
}

// Shutdown executes all registered shutdown functions in reverse order.
// Each function receives the given context (for timeout control).
//
// Returns a ShutdownReport with the result of each function.
// If a function returns an error, continue with the remaining functions
// (don't stop on first error — we want to clean up everything).
type ShutdownResult struct {
	Name     string
	Duration time.Duration
	Err      error
}

type ShutdownReport struct {
	Results  []ShutdownResult
	Duration time.Duration
}

func (c *ShutdownCoordinator) Shutdown(ctx context.Context) ShutdownReport {
	// YOUR CODE HERE
	return ShutdownReport{}
}

// =========================================================================
// Exercise 5: Server Group (Multiple Servers/Workers)
// =========================================================================

// Worker represents a background task that can be started and stopped.
type Worker interface {
	// Start begins the worker. It should block until ctx is cancelled
	// or an error occurs. Return nil on clean shutdown.
	Start(ctx context.Context) error

	// Name returns the worker's name for logging.
	Name() string
}

// ServerGroup manages multiple workers, starting them all and shutting
// them all down together.
//
// If any worker returns an error, all workers should be stopped.
type ServerGroup struct {
	workers []Worker
	mu      sync.Mutex
}

// NewServerGroup creates a new empty ServerGroup.
func NewServerGroup() *ServerGroup {
	// YOUR CODE HERE
	return nil
}

// Add adds a worker to the group.
func (g *ServerGroup) Add(w Worker) {
	// YOUR CODE HERE
}

// Run starts all workers and waits for them to complete.
// The ctx parameter controls the lifecycle — when it's cancelled,
// all workers should stop.
//
// Returns the first error from any worker, or nil if all shut down cleanly.
// All workers must be stopped before Run returns, even if one fails early.
//
// Implementation:
//  1. Create a child context with cancel
//  2. Start each worker in its own goroutine
//  3. If any worker returns a non-nil error, cancel the context
//  4. Wait for all workers to finish
//  5. Return the first error (if any)
func (g *ServerGroup) Run(ctx context.Context) error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 6: Request Drain
// =========================================================================

// RequestTracker tracks in-flight requests and supports waiting for
// them to complete during shutdown.
type RequestTracker struct {
	wg       sync.WaitGroup
	count    atomic.Int64
	draining atomic.Bool
}

// NewRequestTracker creates a new RequestTracker.
func NewRequestTracker() *RequestTracker {
	// YOUR CODE HERE
	return nil
}

// TrackRequest increments the in-flight request count.
// Returns false if the server is draining (not accepting new requests).
// Returns true if the request was accepted.
func (t *RequestTracker) TrackRequest() bool {
	// YOUR CODE HERE
	return false
}

// RequestDone decrements the in-flight request count.
// Call this when a request completes (use defer).
func (t *RequestTracker) RequestDone() {
	// YOUR CODE HERE
}

// InFlight returns the current number of in-flight requests.
func (t *RequestTracker) InFlight() int64 {
	// YOUR CODE HERE
	return 0
}

// Drain sets the draining flag (reject new requests) and waits for
// all in-flight requests to complete, or until the context is cancelled.
//
// Returns nil if all requests completed, or the context error if timed out.
func (t *RequestTracker) Drain(ctx context.Context) error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 7: Readiness Toggle
// =========================================================================

// ReadinessToggle manages the readiness state with an initialization
// and shutdown lifecycle.
//
// Lifecycle:
//  1. Created (not ready)
//  2. MarkReady() called after initialization (ready)
//  3. MarkShuttingDown() called when shutdown starts (not ready)
//
// The state only moves forward — once shutting down, it can't become
// ready again.
type ReadinessToggle struct {
	// state: 0=starting, 1=ready, 2=shutting_down
	state atomic.Int32
}

// NewReadinessToggle creates a new toggle in the "starting" state.
func NewReadinessToggle() *ReadinessToggle {
	// YOUR CODE HERE
	return nil
}

// MarkReady transitions from "starting" to "ready".
// Returns false if already shutting down (invalid transition).
func (r *ReadinessToggle) MarkReady() bool {
	// YOUR CODE HERE
	return false
}

// MarkShuttingDown transitions to "shutting_down" from any state.
// Always succeeds.
func (r *ReadinessToggle) MarkShuttingDown() {
	// YOUR CODE HERE
}

// IsReady returns true only if in the "ready" state.
func (r *ReadinessToggle) IsReady() bool {
	// YOUR CODE HERE
	return false
}

// State returns the current state as a string: "starting", "ready",
// or "shutting_down".
func (r *ReadinessToggle) State() string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 8: Complete Lifecycle Manager
// =========================================================================

// LifecycleManager wires together all the shutdown patterns:
// signal handling, health checks, request tracking, and ordered shutdown.
//
// Usage:
//  1. Create with NewLifecycleManager
//  2. Register shutdown functions
//  3. Call Run() — it starts health checks and waits for shutdown signal
//  4. On signal: mark not ready → drain requests → run shutdown funcs → exit
type LifecycleManager struct {
	readiness       *ReadinessToggle
	tracker         *RequestTracker
	coordinator     *ShutdownCoordinator
	shutdownTimeout time.Duration
	onShutdown      func(string) // callback for logging shutdown events
}

// NewLifecycleManager creates a new manager with the given shutdown timeout.
// The onLog callback is called with status messages during shutdown.
// If onLog is nil, messages are silently discarded.
func NewLifecycleManager(shutdownTimeout time.Duration, onLog func(string)) *LifecycleManager {
	// YOUR CODE HERE
	return nil
}

// Readiness returns the readiness toggle for external use.
func (m *LifecycleManager) Readiness() *ReadinessToggle {
	// YOUR CODE HERE
	return nil
}

// Tracker returns the request tracker for external use.
func (m *LifecycleManager) Tracker() *RequestTracker {
	// YOUR CODE HERE
	return nil
}

// RegisterShutdownFunc adds a function to run during shutdown.
func (m *LifecycleManager) RegisterShutdownFunc(sf ShutdownFunc) {
	// YOUR CODE HERE
}

// MarkReady marks the service as ready to handle traffic.
func (m *LifecycleManager) MarkReady() {
	// YOUR CODE HERE
}

// Shutdown performs the complete shutdown sequence:
//  1. Log "shutdown started"
//  2. Mark as not ready (readiness toggle → shutting_down)
//  3. Log "draining requests"
//  4. Drain in-flight requests (with timeout from shutdownTimeout)
//  5. Log "running shutdown functions"
//  6. Run all registered shutdown functions (with timeout)
//  7. Log "shutdown complete"
//  8. Return the ShutdownReport from the coordinator
func (m *LifecycleManager) Shutdown() ShutdownReport {
	// YOUR CODE HERE
	return ShutdownReport{}
}

// These suppress "imported and not used" errors for stubs.
var _ = os.Interrupt
var _ = syscall.SIGTERM
var _ = fmt.Sprintf
var _ = http.StatusOK
var _ = time.Second
