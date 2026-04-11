package nethttp

import (
	"net/http"
	"sync"
)

/*
=============================================================================
 EXERCISES: net/http Fundamentals
=============================================================================

 Work through these exercises in order. Run tests with:

   make test 14

 Run a single test:

   go test -v -run TestHelloHandler ./14-net-http/

=============================================================================
*/

// Exercise 1: HelloHandler
//
// Write a handler function that returns "Hello, World!" as plain text.
//
// Requirements:
// - Set Content-Type header to "text/plain; charset=utf-8"
// - Set status code to 200 (OK)
// - Write "Hello, World!" as the response body (no trailing newline)
//
// This is the simplest possible HTTP handler, but pay attention to the
// details: setting content type before writing, and explicit status codes.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 2: QueryParamHandler
//
// Write a handler that reads query parameters and returns them as JSON.
//
// Requirements:
//   - Read the "name" query parameter from the request URL
//   - Read the "age" query parameter from the request URL
//   - If "name" is empty, use "anonymous" as the default
//   - If "age" is empty, use "0" as the default
//   - Return a JSON response: {"name": "...", "age": "..."}
//     (age remains a string in the response — no need to convert)
//   - Set Content-Type to "application/json"
//
// Example: GET /info?name=Alice&age=30
// Response: {"age":"30","name":"Alice"}
func QueryParamHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 3: EchoBodyHandler
//
// Write a handler that reads a POST request body and echoes it back.
//
// Requirements:
//   - Only accept POST requests; return 405 Method Not Allowed for others
//     (use http.Error with the message "method not allowed")
//   - Read the entire request body
//   - Set Content-Type to "application/octet-stream"
//   - Write the body back as the response
//   - If reading the body fails, return 400 Bad Request
//
// This exercises the fundamental skill of reading request bodies. Remember:
// the body is a stream — you can only read it once.
func EchoBodyHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 4: CustomHeaderHandler
//
// Write a handler that sets custom response headers and a specific status code.
//
// Requirements:
//   - Set these response headers:
//     X-Request-Id: "12345"
//     X-Powered-By: "Go"
//     Cache-Control: "no-store"
//   - Set status code to 202 (Accepted)
//   - Set Content-Type to "application/json"
//   - Write the JSON body: {"status": "accepted"}
//
// Remember: headers must be set BEFORE WriteHeader() and Write() calls.
func CustomHeaderHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 5: HealthCheckHandler
//
// Write a health check endpoint suitable for load balancers and orchestrators.
//
// Requirements:
//   - Only accept GET requests; return 405 for others
//     (use http.Error with the message "method not allowed")
//   - Set Content-Type to "application/json"
//   - Return status 200
//   - Return JSON body: {"status": "healthy", "version": "1.0.0"}
//
// Health check endpoints are critical in production. Load balancers (like
// AWS ALB, Kubernetes, etc.) poll these to determine if your service is
// ready to receive traffic.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 6: ContentNegotiationHandler
//
// Write a handler that serves different content based on the Accept header.
//
// Requirements:
// - Check the Accept header of the request
// - If Accept contains "application/json":
//   - Set Content-Type to "application/json"
//   - Return: {"message": "hello"}
//
// - If Accept contains "text/plain" (or anything else, as the default):
//   - Set Content-Type to "text/plain; charset=utf-8"
//   - Return: hello
//
// - Use strings.Contains to check the Accept header
//
// Content negotiation lets a single endpoint serve multiple formats.
// It's less common in pure APIs but important in services that serve
// both browsers and API clients.
func ContentNegotiationHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// VisitCounter is a struct-based handler that counts requests.
// You need to implement the ServeHTTP method to make it satisfy http.Handler.
//
// Fields:
// - mu: protects concurrent access to count
// - count: the number of requests handled
type VisitCounter struct {
	mu    sync.Mutex
	count int
}

// Exercise 7: VisitCounter.ServeHTTP
//
// Implement ServeHTTP on VisitCounter to make it an http.Handler.
//
// Requirements:
// - Increment the count safely using the mutex
// - Return JSON: {"visits": <count>} where count is the new total
// - Set Content-Type to "application/json"
// - Status code 200
//
// Why a struct handler? Because sometimes handlers need mutable state.
// Using a struct with a mutex is the Go-idiomatic way to handle this.
// This pattern comes up constantly: rate limiters, caches, connection
// pools, and metrics collectors are all struct handlers.
func (vc *VisitCounter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// CurrentCount returns the current visit count safely.
func (vc *VisitCounter) CurrentCount() int {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return vc.count
}

// Exercise 8: FormValidationHandler
//
// Write a handler that reads form data and validates required fields.
//
// Requirements:
//   - Only accept POST requests; return 405 for others
//     (use http.Error with the message "method not allowed")
//   - Parse the form data from the request body
//   - Require these fields: "username", "email"
//   - If any required field is empty or missing, return 400 Bad Request
//     with JSON: {"error": "missing required field: <fieldname>"}
//     Check username first, then email.
//   - If all fields are present, return 200 with JSON:
//     {"username": "...", "email": "..."}
//   - Set Content-Type to "application/json" for all responses
//
// Form validation is a fundamental server-side concern. Never trust
// client-side validation alone — always validate on the server.
func FormValidationHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}
