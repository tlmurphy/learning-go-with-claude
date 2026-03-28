package middleware

import (
	"net/http"
	"time"
)

/*
=============================================================================
 EXERCISES: Middleware
=============================================================================

These exercises build your middleware skills from basic wrappers to a
complete middleware chain. Middleware is where Go's HTTP model really
shines — the composability of the http.Handler interface makes it easy
to add cross-cutting concerns cleanly.

Every middleware follows the signature: func(http.Handler) http.Handler
(or a function that returns one, for configurable middleware).

=============================================================================
*/

// Exercise 1: LoggingMW
//
// Write a logging middleware that records request details.
//
// Requirements:
// - Use a StatusCapture (Exercise 6) or the ResponseCapture from lesson.go
//   to capture the status code. Since you might not have done Exercise 6 yet,
//   you can use ResponseCapture from the lesson directly.
// - After the handler runs, add these response headers:
//   X-Log-Method: the request method (e.g., "GET")
//   X-Log-Path:   the request URL path (e.g., "/hello")
//   X-Log-Status: the status code as a string (e.g., "200")
//
// Note: In a real application you'd write to a logger, not headers. But
// headers are testable without mocking a logger, which keeps the exercise
// focused on the middleware pattern.
//
// Hint: Use fmt.Sprintf("%d", code) to convert the status code to a string.
func LoggingMW(next http.Handler) http.Handler {
	// YOUR CODE HERE
	return next
}

// Exercise 2: RecoveryMW
//
// Write a recovery middleware that catches panics and returns a 500 error.
//
// Requirements:
// - Use defer/recover to catch panics from the next handler
// - If a panic occurs:
//   - Set Content-Type to "application/json"
//   - Return status 500
//   - Write JSON body: {"error": "internal server error"}
// - If no panic, the request should proceed normally
//
// This prevents a single bad handler from crashing your entire server.
func RecoveryMW(next http.Handler) http.Handler {
	// YOUR CODE HERE
	return next
}

// Exercise 3: RequestIDMW
//
// Write a request ID middleware.
//
// Requirements:
// - Check if the request has an "X-Request-Id" header
// - If present, use that value as the request ID
// - If not present, generate one using the provided generator function
// - Set the "X-Request-Id" response header to the request ID
// - Add the request ID to the request context using RequestIDKey
//   (defined in lesson.go) so downstream handlers can access it
//   via GetRequestID(r.Context())
// - Call the next handler with the updated request
//
// The generator parameter allows tests to provide deterministic IDs.
func RequestIDMW(generator func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// YOUR CODE HERE
		return next
	}
}

// Exercise 4: CORSMW
//
// Write a CORS middleware with configurable origins and methods.
//
// Requirements:
// - Accept a list of allowed origins and allowed methods
// - For every request with an "Origin" header that matches an allowed origin:
//   - Set Access-Control-Allow-Origin to the matched origin
//   - Set Access-Control-Allow-Methods to the allowed methods (joined by ", ")
// - For OPTIONS requests (preflight):
//   - Set the CORS headers as above
//   - Return 204 No Content immediately (don't call the next handler)
// - For non-OPTIONS requests, set CORS headers and call the next handler
// - If the origin doesn't match, just call the next handler without CORS headers
//
// Hint: strings.Join is useful for joining the methods list.
func CORSMW(allowedOrigins []string, allowedMethods []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// YOUR CODE HERE
		return next
	}
}

// Exercise 5: SimpleRateLimitMW
//
// Write a simple rate limiting middleware using a token bucket.
//
// Requirements:
// - Accept a maxRequests parameter (the bucket size / max burst)
// - Accept a refillInterval parameter (how often one token is added)
// - Track tokens globally (not per-IP — keep it simple)
// - Start with maxRequests tokens
// - Each request consumes one token
// - Tokens are refilled based on time elapsed since last request
// - If no tokens available, return:
//   - Status: 429 Too Many Requests
//   - Content-Type: application/json
//   - Body: {"error": "rate limit exceeded"}
// - If tokens available, call the next handler
//
// Note: This is a simplified single-bucket limiter for learning purposes.
// Production rate limiters are per-IP or per-user and often Redis-backed.
func SimpleRateLimitMW(maxRequests int, refillInterval time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// YOUR CODE HERE
		return next
	}
}

// StatusCapture wraps http.ResponseWriter to capture the response status code.
// Exercise 6: Implement the http.ResponseWriter interface methods.
//
// Requirements:
// - Embed http.ResponseWriter for the Header() method
// - Override WriteHeader to capture the status code in the Code field
//   (only capture the FIRST call to WriteHeader — subsequent calls should
//   still forward to the underlying writer but not update Code)
// - Override Write to:
//   a) If WriteHeader hasn't been called yet, call WriteHeader(200) first
//   b) Track the total number of bytes written in the Written field
//   c) Forward the write to the underlying ResponseWriter
// - Initialize Code to 200 (the HTTP default)
//
// This is the most important pattern in Go middleware. Once you understand
// it, you can build any observability middleware.
type StatusCapture struct {
	http.ResponseWriter
	Code        int
	Written     int
	wroteHeader bool
}

// NewStatusCapture creates a StatusCapture wrapping the given ResponseWriter.
func NewStatusCapture(w http.ResponseWriter) *StatusCapture {
	return &StatusCapture{
		ResponseWriter: w,
		Code:           http.StatusOK,
	}
}

// WriteHeader captures the status code and forwards to the underlying writer.
func (sc *StatusCapture) WriteHeader(code int) {
	// YOUR CODE HERE
}

// Write captures bytes written and forwards to the underlying writer.
func (sc *StatusCapture) Write(b []byte) (int, error) {
	// YOUR CODE HERE
	return 0, nil
}

// Exercise 7: TimeoutMW
//
// Write a middleware that adds a timeout to all requests using the standard
// library's http.TimeoutHandler.
//
// Requirements:
// - Accept a timeout duration parameter
// - Use http.TimeoutHandler to wrap the handler
// - Use the message "request timeout" as the timeout message
// - Return the wrapped handler
//
// http.TimeoutHandler is built into Go's stdlib. It runs the handler in
// a separate goroutine, and if it doesn't complete within the timeout,
// sends a 503 Service Unavailable response with the given message.
func TimeoutMW(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// YOUR CODE HERE
		return next
	}
}

// Exercise 8: ChainMiddleware
//
// Build a middleware chain helper that applies multiple middlewares to a handler.
//
// Requirements:
// - Accept an http.Handler and a variadic list of middleware functions
// - Apply the middlewares so that the FIRST middleware in the list is the
//   OUTERMOST wrapper (first to run on the request, last to complete)
// - Return the fully wrapped handler
//
// Example:
//   ChainMiddleware(handler, mw1, mw2, mw3)
//   Should produce: mw1(mw2(mw3(handler)))
//
// This means mw1 runs first, then mw2, then mw3, then the handler.
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// YOUR CODE HERE
	return handler
}
