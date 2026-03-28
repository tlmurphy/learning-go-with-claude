package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

// =========================================================================
// Test Exercise 1: Signal Handler
// =========================================================================

func TestSignalHandler(t *testing.T) {
	t.Run("receives signal and calls handler", func(t *testing.T) {
		received := make(chan os.Signal, 1)
		cancel := SignalHandler(func(sig os.Signal) {
			received <- sig
		}, syscall.SIGUSR1) // Use SIGUSR1 to avoid interfering with test runner
		defer cancel()

		// Send the signal to ourselves
		proc, _ := os.FindProcess(os.Getpid())
		proc.Signal(syscall.SIGUSR1)

		select {
		case sig := <-received:
			if sig != syscall.SIGUSR1 {
				t.Errorf("Expected SIGUSR1, got %v", sig)
			}
		case <-time.After(2 * time.Second):
			t.Error("Timed out waiting for signal handler to be called")
		}
	})

	t.Run("cancel stops listening", func(t *testing.T) {
		callCount := 0
		cancel := SignalHandler(func(sig os.Signal) {
			callCount++
		}, syscall.SIGUSR1)

		cancel() // stop listening immediately

		// Send signal — should NOT trigger the handler
		proc, _ := os.FindProcess(os.Getpid())
		proc.Signal(syscall.SIGUSR1)

		time.Sleep(100 * time.Millisecond)
		if callCount > 0 {
			t.Errorf("Handler should not be called after cancel, but was called %d times", callCount)
		}
	})
}

// =========================================================================
// Test Exercise 2: Graceful HTTP Server with Timeout
// =========================================================================

func TestManagedServer(t *testing.T) {
	t.Run("start and stop", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})

		srv := NewManagedServer(":0", handler, 5*time.Second)
		if srv == nil {
			t.Fatal("NewManagedServer returned nil")
		}

		done, err := srv.Start()
		if err != nil {
			t.Fatalf("Start error: %v", err)
		}
		if done == nil {
			t.Fatal("Start returned nil done channel")
		}

		// Give the server a moment to start
		time.Sleep(50 * time.Millisecond)

		if !srv.IsRunning() {
			t.Error("Expected server to be running")
		}

		err = srv.Stop()
		if err != nil {
			t.Fatalf("Stop error: %v", err)
		}

		// Wait for the done channel
		select {
		case <-done:
			// Good
		case <-time.After(2 * time.Second):
			t.Error("Server did not stop within timeout")
		}

		if srv.IsRunning() {
			t.Error("Expected server to not be running after stop")
		}
	})
}

// =========================================================================
// Test Exercise 3: Health Check Handler
// =========================================================================

func TestHealthStatus(t *testing.T) {
	t.Run("initial state", func(t *testing.T) {
		h := NewHealthStatus()
		if h == nil {
			t.Fatal("NewHealthStatus returned nil")
		}

		// Initially alive but not ready
		handler := h.HealthHandler()
		if handler == nil {
			t.Fatal("HealthHandler returned nil")
		}

		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected liveness 200, got %d", rec.Code)
		}

		readyHandler := h.ReadyHandler()
		if readyHandler == nil {
			t.Fatal("ReadyHandler returned nil")
		}

		req = httptest.NewRequest("GET", "/ready", nil)
		rec = httptest.NewRecorder()
		readyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected readiness 503 initially, got %d", rec.Code)
		}
	})

	t.Run("ready state", func(t *testing.T) {
		h := NewHealthStatus()
		if h == nil {
			t.Fatal("NewHealthStatus returned nil")
		}

		h.SetReady(true)

		req := httptest.NewRequest("GET", "/ready", nil)
		rec := httptest.NewRecorder()
		h.ReadyHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected readiness 200 when ready, got %d", rec.Code)
		}
	})

	t.Run("not alive state", func(t *testing.T) {
		h := NewHealthStatus()
		if h == nil {
			t.Fatal("NewHealthStatus returned nil")
		}

		h.SetAlive(false)

		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()
		h.HealthHandler().ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected liveness 503 when not alive, got %d", rec.Code)
		}
	})

	t.Run("shutdown lifecycle", func(t *testing.T) {
		h := NewHealthStatus()
		if h == nil {
			t.Fatal("NewHealthStatus returned nil")
		}

		// Start: alive, not ready
		checkHealth(t, h.HealthHandler(), http.StatusOK, "alive initially")
		checkReady(t, h.ReadyHandler(), http.StatusServiceUnavailable, "not ready initially")

		// After init: alive, ready
		h.SetReady(true)
		checkHealth(t, h.HealthHandler(), http.StatusOK, "alive after init")
		checkReady(t, h.ReadyHandler(), http.StatusOK, "ready after init")

		// Shutdown: alive, not ready
		h.SetReady(false)
		checkHealth(t, h.HealthHandler(), http.StatusOK, "alive during shutdown")
		checkReady(t, h.ReadyHandler(), http.StatusServiceUnavailable, "not ready during shutdown")
	})
}

func checkHealth(t *testing.T, handler http.HandlerFunc, expectedCode int, msg string) {
	t.Helper()
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != expectedCode {
		t.Errorf("%s: expected %d, got %d", msg, expectedCode, rec.Code)
	}
}

func checkReady(t *testing.T, handler http.HandlerFunc, expectedCode int, msg string) {
	t.Helper()
	req := httptest.NewRequest("GET", "/ready", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != expectedCode {
		t.Errorf("%s: expected %d, got %d", msg, expectedCode, rec.Code)
	}
}

// =========================================================================
// Test Exercise 4: Shutdown Coordinator
// =========================================================================

func TestShutdownCoordinator(t *testing.T) {
	t.Run("executes in reverse order", func(t *testing.T) {
		coord := NewShutdownCoordinator()
		if coord == nil {
			t.Fatal("NewShutdownCoordinator returned nil")
		}

		var order []string
		var mu sync.Mutex

		coord.Register(ShutdownFunc{
			Name: "database",
			Fn: func(ctx context.Context) error {
				mu.Lock()
				order = append(order, "database")
				mu.Unlock()
				return nil
			},
		})
		coord.Register(ShutdownFunc{
			Name: "cache",
			Fn: func(ctx context.Context) error {
				mu.Lock()
				order = append(order, "cache")
				mu.Unlock()
				return nil
			},
		})
		coord.Register(ShutdownFunc{
			Name: "server",
			Fn: func(ctx context.Context) error {
				mu.Lock()
				order = append(order, "server")
				mu.Unlock()
				return nil
			},
		})

		ctx := context.Background()
		report := coord.Shutdown(ctx)

		if len(report.Results) != 3 {
			t.Fatalf("Expected 3 results, got %d", len(report.Results))
		}

		// Reverse order: server, cache, database
		if len(order) != 3 {
			t.Fatalf("Expected 3 items in order, got %d", len(order))
		}
		if order[0] != "server" {
			t.Errorf("Expected first shutdown='server', got %q", order[0])
		}
		if order[1] != "cache" {
			t.Errorf("Expected second shutdown='cache', got %q", order[1])
		}
		if order[2] != "database" {
			t.Errorf("Expected third shutdown='database', got %q", order[2])
		}
	})

	t.Run("continues on error", func(t *testing.T) {
		coord := NewShutdownCoordinator()
		if coord == nil {
			t.Fatal("NewShutdownCoordinator returned nil")
		}

		coord.Register(ShutdownFunc{
			Name: "first",
			Fn:   func(ctx context.Context) error { return nil },
		})
		coord.Register(ShutdownFunc{
			Name: "failing",
			Fn:   func(ctx context.Context) error { return fmt.Errorf("shutdown failed") },
		})
		coord.Register(ShutdownFunc{
			Name: "last",
			Fn:   func(ctx context.Context) error { return nil },
		})

		report := coord.Shutdown(context.Background())

		if len(report.Results) != 3 {
			t.Fatalf("Expected 3 results (even with error), got %d", len(report.Results))
		}

		// Find the failing result
		var foundError bool
		for _, r := range report.Results {
			if r.Name == "failing" && r.Err != nil {
				foundError = true
			}
		}
		if !foundError {
			t.Error("Expected to find error for 'failing' shutdown func")
		}
	})

	t.Run("reports duration", func(t *testing.T) {
		coord := NewShutdownCoordinator()
		if coord == nil {
			t.Fatal("NewShutdownCoordinator returned nil")
		}

		coord.Register(ShutdownFunc{
			Name: "slow",
			Fn: func(ctx context.Context) error {
				time.Sleep(50 * time.Millisecond)
				return nil
			},
		})

		report := coord.Shutdown(context.Background())

		if report.Duration < 50*time.Millisecond {
			t.Errorf("Expected total duration >= 50ms, got %v", report.Duration)
		}
		if len(report.Results) > 0 && report.Results[0].Duration < 50*time.Millisecond {
			t.Errorf("Expected 'slow' duration >= 50ms, got %v", report.Results[0].Duration)
		}
	})
}

// =========================================================================
// Test Exercise 5: Server Group
// =========================================================================

func TestServerGroup(t *testing.T) {
	t.Run("all workers run and stop", func(t *testing.T) {
		group := NewServerGroup()
		if group == nil {
			t.Fatal("NewServerGroup returned nil")
		}

		var started sync.Map

		group.Add(&testWorker{
			name: "worker-1",
			startFn: func(ctx context.Context) error {
				started.Store("worker-1", true)
				<-ctx.Done()
				return nil
			},
		})
		group.Add(&testWorker{
			name: "worker-2",
			startFn: func(ctx context.Context) error {
				started.Store("worker-2", true)
				<-ctx.Done()
				return nil
			},
		})

		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)
		go func() {
			errCh <- group.Run(ctx)
		}()

		// Wait a bit for workers to start
		time.Sleep(100 * time.Millisecond)

		_, ok1 := started.Load("worker-1")
		_, ok2 := started.Load("worker-2")
		if !ok1 || !ok2 {
			t.Error("Expected both workers to have started")
		}

		// Cancel and wait for completion
		cancel()
		select {
		case err := <-errCh:
			if err != nil {
				t.Errorf("Expected nil error on clean shutdown, got: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Error("ServerGroup.Run did not return within timeout")
		}
	})

	t.Run("error in one worker stops all", func(t *testing.T) {
		group := NewServerGroup()
		if group == nil {
			t.Fatal("NewServerGroup returned nil")
		}

		group.Add(&testWorker{
			name: "healthy",
			startFn: func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			},
		})
		group.Add(&testWorker{
			name: "failing",
			startFn: func(ctx context.Context) error {
				time.Sleep(50 * time.Millisecond)
				return fmt.Errorf("worker failed")
			},
		})

		ctx := context.Background()
		errCh := make(chan error, 1)
		go func() {
			errCh <- group.Run(ctx)
		}()

		select {
		case err := <-errCh:
			if err == nil {
				t.Error("Expected an error when worker fails")
			}
		case <-time.After(2 * time.Second):
			t.Error("ServerGroup.Run did not return within timeout")
		}
	})
}

// testWorker implements the Worker interface for testing.
type testWorker struct {
	name    string
	startFn func(ctx context.Context) error
}

func (w *testWorker) Start(ctx context.Context) error {
	return w.startFn(ctx)
}

func (w *testWorker) Name() string {
	return w.name
}

// =========================================================================
// Test Exercise 6: Request Drain
// =========================================================================

func TestRequestTracker(t *testing.T) {
	t.Run("track and complete requests", func(t *testing.T) {
		tracker := NewRequestTracker()
		if tracker == nil {
			t.Fatal("NewRequestTracker returned nil")
		}

		if !tracker.TrackRequest() {
			t.Error("Expected TrackRequest to return true")
		}
		if !tracker.TrackRequest() {
			t.Error("Expected TrackRequest to return true")
		}

		if tracker.InFlight() != 2 {
			t.Errorf("Expected 2 in-flight, got %d", tracker.InFlight())
		}

		tracker.RequestDone()
		if tracker.InFlight() != 1 {
			t.Errorf("Expected 1 in-flight after done, got %d", tracker.InFlight())
		}
	})

	t.Run("drain waits for requests", func(t *testing.T) {
		tracker := NewRequestTracker()
		if tracker == nil {
			t.Fatal("NewRequestTracker returned nil")
		}

		tracker.TrackRequest()
		tracker.TrackRequest()

		// Complete requests after a delay
		go func() {
			time.Sleep(50 * time.Millisecond)
			tracker.RequestDone()
			tracker.RequestDone()
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := tracker.Drain(ctx)
		if err != nil {
			t.Errorf("Drain returned error: %v", err)
		}

		if tracker.InFlight() != 0 {
			t.Errorf("Expected 0 in-flight after drain, got %d", tracker.InFlight())
		}
	})

	t.Run("drain rejects new requests", func(t *testing.T) {
		tracker := NewRequestTracker()
		if tracker == nil {
			t.Fatal("NewRequestTracker returned nil")
		}

		// Start draining in background
		go func() {
			ctx := context.Background()
			tracker.Drain(ctx)
		}()

		time.Sleep(50 * time.Millisecond) // let drain start

		if tracker.TrackRequest() {
			t.Error("Expected TrackRequest to return false during drain")
		}
	})

	t.Run("drain timeout returns error", func(t *testing.T) {
		tracker := NewRequestTracker()
		if tracker == nil {
			t.Fatal("NewRequestTracker returned nil")
		}

		tracker.TrackRequest() // never completed

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := tracker.Drain(ctx)
		if err == nil {
			t.Error("Expected timeout error from drain")
		}
	})
}

// =========================================================================
// Test Exercise 7: Readiness Toggle
// =========================================================================

func TestReadinessToggle(t *testing.T) {
	t.Run("lifecycle transitions", func(t *testing.T) {
		r := NewReadinessToggle()
		if r == nil {
			t.Fatal("NewReadinessToggle returned nil")
		}

		// Initial state: starting (not ready)
		if r.IsReady() {
			t.Error("Expected not ready in starting state")
		}
		if r.State() != "starting" {
			t.Errorf("Expected state='starting', got %q", r.State())
		}

		// Transition to ready
		ok := r.MarkReady()
		if !ok {
			t.Error("Expected MarkReady to return true")
		}
		if !r.IsReady() {
			t.Error("Expected ready after MarkReady")
		}
		if r.State() != "ready" {
			t.Errorf("Expected state='ready', got %q", r.State())
		}

		// Transition to shutting down
		r.MarkShuttingDown()
		if r.IsReady() {
			t.Error("Expected not ready during shutdown")
		}
		if r.State() != "shutting_down" {
			t.Errorf("Expected state='shutting_down', got %q", r.State())
		}
	})

	t.Run("cannot become ready after shutdown", func(t *testing.T) {
		r := NewReadinessToggle()
		if r == nil {
			t.Fatal("NewReadinessToggle returned nil")
		}

		r.MarkShuttingDown()
		ok := r.MarkReady()
		if ok {
			t.Error("Expected MarkReady to return false when shutting down")
		}
		if r.IsReady() {
			t.Error("Expected not ready when shutting down, even after MarkReady")
		}
	})

	t.Run("can skip ready and go to shutting down", func(t *testing.T) {
		r := NewReadinessToggle()
		if r == nil {
			t.Fatal("NewReadinessToggle returned nil")
		}

		r.MarkShuttingDown()
		if r.State() != "shutting_down" {
			t.Errorf("Expected state='shutting_down', got %q", r.State())
		}
	})
}

// =========================================================================
// Test Exercise 8: Complete Lifecycle Manager
// =========================================================================

func TestLifecycleManager(t *testing.T) {
	t.Run("full lifecycle", func(t *testing.T) {
		var logs []string
		var mu sync.Mutex

		logFn := func(msg string) {
			mu.Lock()
			logs = append(logs, msg)
			mu.Unlock()
		}

		mgr := NewLifecycleManager(5*time.Second, logFn)
		if mgr == nil {
			t.Fatal("NewLifecycleManager returned nil")
		}

		// Register shutdown functions
		var shutdownOrder []string
		mgr.RegisterShutdownFunc(ShutdownFunc{
			Name: "database",
			Fn: func(ctx context.Context) error {
				mu.Lock()
				shutdownOrder = append(shutdownOrder, "database")
				mu.Unlock()
				return nil
			},
		})
		mgr.RegisterShutdownFunc(ShutdownFunc{
			Name: "cache",
			Fn: func(ctx context.Context) error {
				mu.Lock()
				shutdownOrder = append(shutdownOrder, "cache")
				mu.Unlock()
				return nil
			},
		})

		// Mark ready
		mgr.MarkReady()
		if !mgr.Readiness().IsReady() {
			t.Error("Expected ready after MarkReady")
		}

		// Track some requests
		if mgr.Tracker() == nil {
			t.Fatal("Tracker returned nil")
		}
		mgr.Tracker().TrackRequest()

		// Complete the request
		mgr.Tracker().RequestDone()

		// Trigger shutdown
		report := mgr.Shutdown()

		// Verify readiness changed
		if mgr.Readiness().IsReady() {
			t.Error("Expected not ready after shutdown")
		}

		// Verify shutdown functions ran in reverse order
		if len(shutdownOrder) != 2 {
			t.Fatalf("Expected 2 shutdown functions, got %d", len(shutdownOrder))
		}
		if shutdownOrder[0] != "cache" {
			t.Errorf("Expected first shutdown='cache', got %q", shutdownOrder[0])
		}
		if shutdownOrder[1] != "database" {
			t.Errorf("Expected second shutdown='database', got %q", shutdownOrder[1])
		}

		// Verify report
		if len(report.Results) != 2 {
			t.Errorf("Expected 2 results in report, got %d", len(report.Results))
		}

		// Verify logs were written
		mu.Lock()
		logCount := len(logs)
		mu.Unlock()
		if logCount < 3 {
			t.Errorf("Expected at least 3 log messages, got %d", logCount)
		}
	})

	t.Run("shutdown drains in-flight requests", func(t *testing.T) {
		mgr := NewLifecycleManager(5*time.Second, nil)
		if mgr == nil {
			t.Fatal("NewLifecycleManager returned nil")
		}

		mgr.MarkReady()

		// Start a request that will complete during drain
		mgr.Tracker().TrackRequest()
		go func() {
			time.Sleep(50 * time.Millisecond)
			mgr.Tracker().RequestDone()
		}()

		// Shutdown should wait for the request
		report := mgr.Shutdown()
		_ = report

		if mgr.Tracker().InFlight() != 0 {
			t.Errorf("Expected 0 in-flight after shutdown, got %d", mgr.Tracker().InFlight())
		}
	})
}
