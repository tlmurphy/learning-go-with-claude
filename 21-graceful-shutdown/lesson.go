package shutdown

/*
=============================================================================
 Module 21: Graceful Shutdown
=============================================================================

 A web service that stops abruptly when you press Ctrl+C or when
 Kubernetes sends SIGTERM is a liability. In-flight requests get dropped,
 database transactions abort, and websocket connections break. Users see
 errors. Data can be lost.

 Graceful shutdown means: stop accepting new work, finish what's in
 progress, clean up resources, then exit. It sounds simple, but getting
 the ordering and timeouts right requires understanding several Go
 patterns working together.

 WHY THIS MATTERS FOR WEB SERVICES:
 - Kubernetes sends SIGTERM before killing your pod (default 30s grace)
 - Load balancers need time to drain connections
 - Database transactions need to complete or roll back cleanly
 - Background workers need to finish their current task
 - Health checks must reflect the shutdown state (ready → not ready)

 The good news: Go's standard library has excellent support for all of
 this. http.Server.Shutdown, os/signal, and context work together
 beautifully.

=============================================================================
*/

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// -------------------------------------------------------------------------
// Signal Handling: os/signal
// -------------------------------------------------------------------------

/*
 Unix signals are how the operating system (or container orchestrator)
 tells your process to do something. The two you care about:

 SIGINT (signal 2): Sent when you press Ctrl+C. Means "interrupt."
 SIGTERM (signal 15): Sent by Kubernetes, systemd, docker stop. Means
   "please terminate gracefully."

 By default, Go terminates immediately on both signals. To handle them
 gracefully, you register a notification channel:

   ctx, stop := signal.NotifyContext(context.Background(),
       syscall.SIGINT, syscall.SIGTERM)
   defer stop()

   <-ctx.Done() // blocks until signal received
   // ctx is now cancelled — propagate to all operations

 signal.NotifyContext (Go 1.16+) is the modern way. It creates a context
 that cancels when a signal arrives, which integrates perfectly with
 Go's context-based cancellation.

 The older approach uses signal.Notify with a channel:
   c := make(chan os.Signal, 1)
   signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
   <-c // blocks until signal received
*/

// DemoSignalContext shows how to use signal.NotifyContext for shutdown.
// This is illustrative — in real code, you'd use this in main().
func DemoSignalContext() {
	// Create a context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Running... press Ctrl+C to stop")

	// This blocks until a signal arrives
	<-ctx.Done()
	fmt.Println("Received signal, shutting down...")
}

// -------------------------------------------------------------------------
// HTTP Server Graceful Shutdown
// -------------------------------------------------------------------------

/*
 http.Server.Shutdown is the key method. It:
   1. Closes the listener (stops accepting new connections)
   2. Waits for all active connections to become idle
   3. Returns after all connections are closed

 It does NOT interrupt in-flight requests — they continue processing.
 The context you pass controls the maximum wait time.

 The pattern:

   srv := &http.Server{Addr: ":8080", Handler: mux}

   // Start the server in a goroutine
   go func() {
       if err := srv.ListenAndServe(); err != http.ErrServerClosed {
           log.Fatal(err)
       }
   }()

   // Wait for shutdown signal
   <-ctx.Done()

   // Give in-flight requests 10 seconds to finish
   shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()

   if err := srv.Shutdown(shutdownCtx); err != nil {
       log.Printf("Forced shutdown: %v", err)
   }

 IMPORTANT: Always set a timeout on the shutdown context. Without one,
 a single slow connection can prevent your service from ever stopping.
 Kubernetes will eventually SIGKILL your process, which is the ungraceful
 exit you're trying to avoid.
*/

// GracefulServer wraps an HTTP server with graceful shutdown support.
type GracefulServer struct {
	server *http.Server
	done   chan struct{}
}

// NewGracefulServer creates a new server with the given address and handler.
func NewGracefulServer(addr string, handler http.Handler) *GracefulServer {
	return &GracefulServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		done: make(chan struct{}),
	}
}

// Start begins listening for HTTP requests. This method returns
// immediately. The server runs in a background goroutine.
func (s *GracefulServer) Start() error {
	var startErr error
	started := make(chan struct{})

	go func() {
		close(started)
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			startErr = err
		}
		close(s.done)
	}()

	<-started
	return startErr
}

// Shutdown gracefully stops the server, waiting up to the given timeout
// for in-flight requests to complete.
func (s *GracefulServer) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := s.server.Shutdown(ctx)

	// Wait for the server goroutine to finish
	<-s.done
	return err
}

// -------------------------------------------------------------------------
// Health Check Endpoints
// -------------------------------------------------------------------------

/*
 Kubernetes (and most load balancers) use two types of health checks:

 LIVENESS PROBE (/health or /healthz):
   "Is the process alive and not deadlocked?"
   - Returns 200 if the process is running
   - If this fails, Kubernetes RESTARTS the pod
   - Keep this simple — just return 200
   - Don't check databases or dependencies (that's readiness)

 READINESS PROBE (/ready or /readyz):
   "Can this instance handle traffic?"
   - Returns 200 when the service can serve requests
   - Returns 503 during startup (before initialization completes)
   - Returns 503 during shutdown (while draining)
   - If this fails, Kubernetes removes from the load balancer
   - Check critical dependencies here (DB connection, etc.)

 Why two endpoints? Because "alive but not ready" is different from
 "dead." During a deployment, your new pod is alive but still loading
 data. You don't want Kubernetes to restart it — you want it to wait.
 During shutdown, it's alive but should stop getting new traffic.

 The readiness probe is your primary tool for zero-downtime deployments.
*/

// HealthChecker manages liveness and readiness state.
type HealthChecker struct {
	ready atomic.Bool
}

// NewHealthChecker creates a new HealthChecker. Initially not ready.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

// SetReady marks the service as ready to handle traffic.
func (h *HealthChecker) SetReady(ready bool) {
	h.ready.Store(ready)
}

// IsReady returns whether the service is ready.
func (h *HealthChecker) IsReady() bool {
	return h.ready.Load()
}

// LivenessHandler returns an HTTP handler for liveness checks.
// Always returns 200 OK with body "ok" if the process is running.
func (h *HealthChecker) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// ReadinessHandler returns an HTTP handler for readiness checks.
// Returns 200 OK if ready, 503 Service Unavailable if not.
func (h *HealthChecker) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.IsReady() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("not ready"))
		}
	}
}

// -------------------------------------------------------------------------
// Shutdown Ordering
// -------------------------------------------------------------------------

/*
 The order in which you shut things down matters. Get it wrong and you'll
 have requests hitting a closed database, or new connections arriving at
 a half-stopped server.

 Correct shutdown order:

   1. Receive signal (SIGINT/SIGTERM)
   2. Mark as not ready (readiness probe returns 503)
      → Load balancer stops sending new traffic
   3. Wait a brief period (2-5 seconds)
      → Let the load balancer propagate the change
   4. Stop accepting new connections (server.Shutdown)
   5. Wait for in-flight requests to complete (with timeout)
   6. Close background workers (wait for current task to finish)
   7. Close database connections
   8. Close other resources (message queues, caches, etc.)
   9. Exit

 Step 3 is subtle but important in Kubernetes. When a pod becomes
 not-ready, it takes a moment for the Endpoints controller to update
 and for kube-proxy to remove the pod from the Service. During this
 window, the pod may still receive traffic. The brief wait ensures
 the load balancer catches up.
*/

// -------------------------------------------------------------------------
// errgroup for Managing Multiple Servers/Workers
// -------------------------------------------------------------------------

/*
 Many production services run multiple things concurrently:
   - HTTP server for API endpoints
   - gRPC server for internal communication
   - Background worker for async tasks
   - Metrics server on a different port

 The errgroup package (golang.org/x/sync/errgroup) manages these
 beautifully. It:
   - Starts multiple goroutines
   - Waits for all to complete
   - Cancels the group context when any goroutine returns an error

 Here's the simplified version with just goroutines and a WaitGroup,
 since we're avoiding external dependencies:

   var wg sync.WaitGroup
   ctx, cancel := context.WithCancel(context.Background())

   // Start servers
   wg.Add(1)
   go func() {
       defer wg.Done()
       runHTTPServer(ctx)
   }()

   wg.Add(1)
   go func() {
       defer wg.Done()
       runWorker(ctx)
   }()

   // Wait for signal
   <-signalCtx.Done()
   cancel()  // tell everything to stop
   wg.Wait() // wait for everything to finish
*/

// -------------------------------------------------------------------------
// Container Lifecycle and Kubernetes
// -------------------------------------------------------------------------

/*
 When Kubernetes terminates a pod, here's what happens:

 1. Pod status changes to Terminating
 2. preStop hook runs (if configured)
 3. SIGTERM is sent to PID 1 in the container
 4. Kubernetes waits (terminationGracePeriodSeconds, default 30s)
 5. SIGKILL is sent if still running

 Your Go service sees step 3 (SIGTERM). You have until step 5 to
 clean up. If your shutdown takes longer than the grace period,
 Kubernetes kills the process forcefully.

 Best practices:
   - Set terminationGracePeriodSeconds to match your expected drain time
   - Use preStop to add a delay (sleep 5) to handle the LB propagation
   - Set your Go shutdown timeout to a few seconds less than the grace period
   - Log the shutdown process so you can debug timeout issues

 Example Kubernetes config:
   spec:
     terminationGracePeriodSeconds: 60
     containers:
     - name: app
       lifecycle:
         preStop:
           exec:
             command: ["/bin/sh", "-c", "sleep 5"]
*/

// -------------------------------------------------------------------------
// Putting It All Together
// -------------------------------------------------------------------------

/*
 Here's the complete pattern for a production Go web service's main():

   func main() {
       // 1. Load config
       cfg := loadConfig()

       // 2. Set up signal handling
       ctx, stop := signal.NotifyContext(context.Background(),
           syscall.SIGINT, syscall.SIGTERM)
       defer stop()

       // 3. Set up health checker
       health := NewHealthChecker()

       // 4. Set up dependencies (DB, etc.)
       db := connectDB(cfg)
       defer db.Close()

       // 5. Set up HTTP server
       mux := http.NewServeMux()
       mux.HandleFunc("/health", health.LivenessHandler())
       mux.HandleFunc("/ready", health.ReadinessHandler())
       // ... register other handlers ...

       srv := &http.Server{Addr: cfg.Addr, Handler: mux}

       // 6. Start server
       go func() {
           if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
               log.Fatal(err)
           }
       }()
       log.Printf("Server listening on %s", cfg.Addr)

       // 7. Mark as ready
       health.SetReady(true)

       // 8. Wait for shutdown signal
       <-ctx.Done()
       log.Println("Shutting down...")

       // 9. Mark as not ready
       health.SetReady(false)

       // 10. Wait for LB to catch up
       time.Sleep(2 * time.Second)

       // 11. Shutdown server (with timeout)
       shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
       defer cancel()
       if err := srv.Shutdown(shutdownCtx); err != nil {
           log.Printf("Server shutdown error: %v", err)
       }

       // 12. Close DB and other resources
       db.Close()

       log.Println("Shutdown complete")
   }
*/

// Suppress "imported and not used" errors for lesson demo code.
var _ = os.Interrupt
var _ = syscall.SIGTERM
var _ = context.Background
var _ = errors.New
var _ = sync.WaitGroup{}
var _ = time.Second
var _ = fmt.Sprintf
