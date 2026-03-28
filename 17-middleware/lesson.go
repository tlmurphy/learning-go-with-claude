// Package middleware covers the middleware pattern in Go's net/http — the
// standard way to add cross-cutting concerns like logging, authentication,
// rate limiting, and panic recovery to HTTP handlers.
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 MIDDLEWARE IN GO
=============================================================================

Middleware is a function that wraps an HTTP handler to add behavior before
and/or after the handler executes. It's how you add cross-cutting concerns
to your web application without repeating code in every handler.

The canonical middleware signature in Go is:

  func(http.Handler) http.Handler

A middleware takes a handler, wraps it with additional logic, and returns
a new handler. This composability is the beauty of the pattern — you can
chain any number of middlewares together.

Why this pattern? Because http.Handler is an interface. Middleware doesn't
change the interface — it wraps it. Your router doesn't know (or care)
that a handler has been wrapped by three middlewares. Everything is still
just http.Handler all the way down.

Typical middleware stack for a production API:
  1. Recovery (outermost — catches panics from everything below)
  2. Request ID (adds a unique ID to every request)
  3. Logging (logs every request with timing)
  4. CORS (handles cross-origin requests)
  5. Rate limiting (prevents abuse)
  6. Authentication (verifies identity)
  7. Your handler (innermost)

Order matters! Recovery should be outermost so it catches panics from
all other middleware. Logging should be before auth so you log both
authenticated and rejected requests. Think of it as layers of an onion.

=============================================================================
 THE BASIC PATTERN
=============================================================================

Here's the simplest possible middleware to understand the pattern:

  func SimpleMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          // Code here runs BEFORE the handler
          fmt.Println("before")

          next.ServeHTTP(w, r)  // Call the actual handler

          // Code here runs AFTER the handler
          fmt.Println("after")
      })
  }

The key insight: next.ServeHTTP(w, r) is where the actual handler (or the
next middleware in the chain) executes. Everything before it is "pre-
processing" and everything after it is "post-processing."

=============================================================================
*/

// TimingMiddleware measures how long each request takes and adds the
// duration as a response header. This is the simplest useful middleware.
func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// After the handler completes, record the duration.
		// Note: this header might not be sent if the handler already
		// flushed the response. In production, you'd use a ResponseWriter
		// wrapper to capture this information instead.
		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())
	})
}

/*
=============================================================================
 THE ResponseWriter WRAPPER PATTERN
=============================================================================

Here's a fundamental problem with middleware: http.ResponseWriter is a
write-only interface. Once the handler calls WriteHeader() or Write(),
you can't read back what was written. But middleware often needs to know:
- What status code was sent?
- How many bytes were written?
- Did the handler write anything at all?

The solution is a wrapper type that implements http.ResponseWriter but
intercepts the calls to capture metadata:

  type responseWriter struct {
      http.ResponseWriter          // embed the real writer
      statusCode    int            // captured status code
      bytesWritten  int            // captured byte count
      wroteHeader   bool           // whether WriteHeader was called
  }

This wrapper is passed to the handler instead of the real ResponseWriter.
The handler doesn't know the difference (it still implements the interface),
but the middleware can now inspect what happened.

This pattern is so common and so important that it's worth mastering. Nearly
every production middleware that logs or records metrics uses it.

=============================================================================
*/

// ResponseCapture wraps http.ResponseWriter to capture the status code
// and bytes written. This is the fundamental building block for logging
// and metrics middleware.
type ResponseCapture struct {
	http.ResponseWriter
	StatusCode   int
	BytesWritten int
	wroteHeader  bool
}

// NewResponseCapture creates a new ResponseCapture wrapping the given writer.
func NewResponseCapture(w http.ResponseWriter) *ResponseCapture {
	return &ResponseCapture{
		ResponseWriter: w,
		StatusCode:     http.StatusOK, // default, like the real ResponseWriter
	}
}

// WriteHeader captures the status code and forwards to the real writer.
func (rc *ResponseCapture) WriteHeader(code int) {
	if !rc.wroteHeader {
		rc.StatusCode = code
		rc.wroteHeader = true
	}
	rc.ResponseWriter.WriteHeader(code)
}

// Write captures the byte count and forwards to the real writer.
// It also calls WriteHeader(200) implicitly if not already called,
// matching the behavior of the real ResponseWriter.
func (rc *ResponseCapture) Write(b []byte) (int, error) {
	if !rc.wroteHeader {
		rc.WriteHeader(http.StatusOK)
	}
	n, err := rc.ResponseWriter.Write(b)
	rc.BytesWritten += n
	return n, err
}

/*
=============================================================================
 LOGGING MIDDLEWARE
=============================================================================

Every production API should log every request. At minimum you want:
- HTTP method and path
- Status code
- Response time
- Client IP (for abuse detection)

Using our ResponseCapture wrapper, we can capture the status code and
log a complete request summary after the handler finishes.

=============================================================================
*/

// LoggingMiddleware logs each request with method, path, status code, and
// duration. It uses the ResponseCapture wrapper to capture the status code.
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the ResponseWriter to capture the status code
			capture := NewResponseCapture(w)

			// Call the next handler with our wrapper
			next.ServeHTTP(capture, r)

			// Log the request after it completes
			duration := time.Since(start)
			logger.Printf("%s %s -> %d (%s, %d bytes)",
				r.Method,
				r.URL.Path,
				capture.StatusCode,
				duration,
				capture.BytesWritten,
			)
		})
	}
}

/*
=============================================================================
 RECOVERY MIDDLEWARE
=============================================================================

Go handlers run in goroutines, and a panic in a goroutine kills the
entire process unless recovered. Recovery middleware catches panics,
logs the stack trace, and returns a 500 Internal Server Error instead
of crashing the server.

This should be the OUTERMOST middleware in your chain — it needs to
wrap everything else to catch panics from any layer.

Why not just use recover() in every handler? Because:
1. You'd have to add it to every handler (tedious, error-prone)
2. You'd miss panics from middleware
3. Centralizing it ensures consistent error responses for panics

=============================================================================
*/

// RecoveryMiddleware catches panics and returns a 500 JSON error response.
// It should be the outermost middleware in your chain.
func RecoveryMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					stack := debug.Stack()
					logger.Printf("PANIC: %v\n%s", err, stack)

					// Return a generic error — never expose internal details
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "internal server error",
					})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

/*
=============================================================================
 REQUEST ID MIDDLEWARE
=============================================================================

Assigning a unique ID to every request is invaluable for debugging. When
a user reports a problem, they can give you the request ID and you can
trace it through all your logs.

The pattern:
1. Check if the client sent an X-Request-Id header (for tracing across services)
2. If not, generate a new one
3. Add it to the response headers (so the client can reference it)
4. Add it to the request context (so handlers and other middleware can access it)

=============================================================================
*/

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

// RequestIDKey is the context key for the request ID.
const RequestIDKey contextKey = "request_id"

// RequestIDMiddleware adds a unique request ID to each request.
// If the client sends X-Request-Id, it's reused. Otherwise a new one
// is generated.
func RequestIDMiddleware(next http.Handler) http.Handler {
	var counter uint64
	var mu sync.Mutex

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for client-provided request ID
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			// Generate a simple sequential ID.
			// In production, use a UUID library (e.g., google/uuid).
			mu.Lock()
			counter++
			reqID = fmt.Sprintf("req-%d", counter)
			mu.Unlock()
		}

		// Add to response headers
		w.Header().Set("X-Request-Id", reqID)

		// Add to request context so handlers can access it
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID extracts the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

/*
=============================================================================
 CORS MIDDLEWARE
=============================================================================

CORS (Cross-Origin Resource Sharing) is necessary when your API is called
from a browser running on a different domain. Without CORS headers, the
browser blocks the response.

The CORS dance:
1. Browser sends a "preflight" OPTIONS request with Origin header
2. Server responds with Access-Control-Allow-* headers
3. If allowed, browser sends the actual request
4. Server includes CORS headers on the actual response too

CORS configuration is security-sensitive. Never use "*" for allowed
origins in production — it means any website can call your API.

=============================================================================
*/

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int // preflight cache duration in seconds
}

// CORSMiddleware adds CORS headers based on the provided configuration.
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	allowedOrigins := make(map[string]bool)
	for _, origin := range config.AllowedOrigins {
		allowedOrigins[origin] = true
	}

	methods := strings.Join(config.AllowedMethods, ", ")
	headers := strings.Join(config.AllowedHeaders, ", ")
	maxAge := fmt.Sprintf("%d", config.MaxAge)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if the origin is allowed
			if allowedOrigins[origin] || allowedOrigins["*"] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", methods)
				w.Header().Set("Access-Control-Allow-Headers", headers)
				w.Header().Set("Access-Control-Max-Age", maxAge)
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
=============================================================================
 RATE LIMITING MIDDLEWARE
=============================================================================

Rate limiting prevents abuse by limiting how many requests a client can
make in a given time window. There are several algorithms:

  Token Bucket: Start with N tokens. Each request costs 1 token. Tokens
  refill at a fixed rate. If no tokens, reject the request.

  Sliding Window: Count requests in the last N seconds. If over the
  limit, reject.

  Fixed Window: Count requests in the current time window (e.g., the
  current minute). Simpler but can have edge effects at window boundaries.

For production rate limiting, consider:
- Per-IP limiting (for unauthenticated endpoints)
- Per-user limiting (for authenticated endpoints)
- Using a distributed rate limiter (Redis-based) for multi-server deployments
- Returning Retry-After headers so clients know when to retry
- Rate limiting by endpoint (login pages need stricter limits)

We'll implement a simple token bucket per-IP rate limiter here.

=============================================================================
*/

// RateLimiter tracks request rates per client IP.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // tokens per interval
	interval time.Duration // refill interval
}

type visitor struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a rate limiter that allows 'rate' requests per
// 'interval' per client IP.
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		interval: interval,
	}
}

// Allow checks if a request from the given IP should be allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[ip]
	if !ok {
		rl.visitors[ip] = &visitor{
			tokens:   rl.rate - 1, // -1 for this request
			lastSeen: time.Now(),
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := time.Since(v.lastSeen)
	v.lastSeen = time.Now()

	tokensToAdd := int(elapsed/rl.interval) * rl.rate
	v.tokens += tokensToAdd
	if v.tokens > rl.rate {
		v.tokens = rl.rate
	}

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

// RateLimitMiddleware limits requests per client IP using a token bucket.
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract client IP (simplified — in production, check X-Forwarded-For)
			ip := r.RemoteAddr

			if !limiter.Allow(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "rate limit exceeded",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
=============================================================================
 MIDDLEWARE CHAINING
=============================================================================

Applying multiple middlewares can get verbose:

  handler := RecoveryMiddleware(logger)(
      LoggingMiddleware(logger)(
          RequestIDMiddleware(
              myHandler,
          ),
      ),
  )

A chain helper makes this much cleaner:

  handler := Chain(myHandler, RecoveryMiddleware(logger), LoggingMiddleware(logger), RequestIDMiddleware)

The chain applies middlewares in order from LEFT to RIGHT, meaning the
leftmost middleware is the outermost wrapper (first to run, last to
complete). This matches how you'd read a middleware stack specification.

=============================================================================
*/

// Chain applies middlewares to a handler in the order given.
// The first middleware in the list is the outermost wrapper.
//
// Usage: Chain(handler, recovery, logging, requestID)
// Equivalent to: recovery(logging(requestID(handler)))
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Apply in reverse order so the first middleware is outermost
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

/*
=============================================================================
 AUTHENTICATION MIDDLEWARE (PREVIEW)
=============================================================================

Authentication middleware validates that the caller has proper credentials.
Full auth is covered in Module 23 — this is a simplified preview to show
the middleware pattern.

The pattern:
1. Extract credentials from the request (header, cookie, etc.)
2. Validate the credentials
3. If invalid, return 401 Unauthorized
4. If valid, add user info to the request context
5. Call the next handler

=============================================================================
*/

// AuthMiddleware is a simplified auth middleware that checks for an
// API key in the Authorization header. In production, you'd validate
// JWTs, check sessions, etc.
func AuthMiddleware(validKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Authorization")
			if key == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "authorization header required",
				})
				return
			}

			// Strip "Bearer " prefix if present
			key = strings.TrimPrefix(key, "Bearer ")

			if key != validKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid API key",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
=============================================================================
 TIMEOUT MIDDLEWARE
=============================================================================

Timeout middleware wraps each request's context with a deadline. If the
handler takes too long, the context is canceled and the client gets a
timeout error.

This is important for:
- Preventing slow handlers from holding goroutines indefinitely
- Cascading timeouts (if your handler calls another service, the context
  cancellation propagates)
- Resource protection (a slow database query doesn't block forever)

Note: http.TimeoutHandler is built into the standard library and does
this, but understanding the pattern helps you build more sophisticated
timeout behavior.

=============================================================================
*/

// TimeoutMiddleware adds a deadline to each request's context.
// If the handler doesn't complete within the duration, the context
// is canceled (but the handler must check ctx.Done() to stop work).
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, `{"error": "request timeout"}`)
	}
}

/*
=============================================================================
 PUTTING IT ALL TOGETHER
=============================================================================

Here's how you'd assemble a production middleware stack:

=============================================================================
*/

// NewProductionStack demonstrates a production middleware configuration.
func NewProductionStack(handler http.Handler, logger *log.Logger) http.Handler {
	limiter := NewRateLimiter(100, time.Minute) // 100 requests per minute

	cors := CORSConfig{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         86400,
	}

	// Apply middlewares — outermost first
	return Chain(handler,
		RecoveryMiddleware(logger),   // 1. Catch panics
		RequestIDMiddleware,          // 2. Add request IDs
		LoggingMiddleware(logger),    // 3. Log requests
		CORSMiddleware(cors),         // 4. Handle CORS
		RateLimitMiddleware(limiter), // 5. Rate limit
	)
}
