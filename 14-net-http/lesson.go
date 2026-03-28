// Package nethttp covers the fundamentals of Go's net/http package — the
// foundation for building web services. Go's standard library HTTP support
// is remarkably capable, and understanding it deeply will serve you well
// whether you use it directly or through a framework.
package nethttp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 NET/HTTP FUNDAMENTALS
=============================================================================

Go's net/http package is unusual among programming languages: it provides a
production-quality HTTP server and client in the standard library. Many
companies run Go's stdlib HTTP server in production handling millions of
requests per day without any third-party server framework.

The key insight is that the entire HTTP serving model in Go rests on a
single interface:

  type Handler interface {
      ServeHTTP(ResponseWriter, *Request)
  }

That's it. Everything else — routers, middleware, frameworks — is built
on top of this one interface. If you understand Handler, you understand
Go's web model.

=============================================================================
 THE http.Handler INTERFACE
=============================================================================

The Handler interface is the cornerstone. Any type that implements
ServeHTTP(http.ResponseWriter, *http.Request) can handle HTTP requests.

Why an interface and not a function? Because sometimes your handler needs
state — a database connection, a logger, configuration. By using an
interface, your handler can be a struct with fields:

  type UserHandler struct {
      DB     *sql.DB
      Logger *log.Logger
  }

  func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
      // has access to h.DB and h.Logger
  }

This is the dependency injection pattern in Go — no framework magic, just
passing dependencies through struct fields.

=============================================================================
*/

// GreetingHandler is a struct that implements http.Handler. This pattern is
// how you create handlers with state (dependencies like databases, config, etc.).
type GreetingHandler struct {
	DefaultName string
}

// ServeHTTP implements the http.Handler interface. This is the method that
// gets called for every incoming request routed to this handler.
func (h *GreetingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = h.DefaultName
	}

	// Always set Content-Type BEFORE writing the body.
	// Once you call Write() or WriteHeader(), headers are sent and can't be changed.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello, %s!", name)
}

/*
=============================================================================
 http.HandlerFunc — THE ADAPTER
=============================================================================

Most of the time you don't need a full struct — you just want a function.
Go provides http.HandlerFunc, which is a type adapter that converts any
function with the right signature into an http.Handler:

  type HandlerFunc func(ResponseWriter, *Request)

  func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
      f(w, r)
  }

This is a beautiful example of Go's adapter pattern. HandlerFunc is both
a type and an implementation of Handler. When you see:

  http.HandleFunc("/path", myFunc)

...it's wrapping your function in HandlerFunc behind the scenes.

This duality — struct handlers for stateful things, function handlers for
simple things — gives you flexibility without complexity.

=============================================================================
*/

// DemoHandlerFunc shows how plain functions serve as HTTP handlers.
// This is the most common pattern for simple endpoints.
func DemoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "This is a plain handler function.")
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
}

/*
=============================================================================
 http.ServeMux — THE ROUTER
=============================================================================

ServeMux is Go's built-in HTTP request router (multiplexer). It matches
incoming requests to registered patterns and dispatches to the right handler.

As of Go 1.22, ServeMux got a major upgrade with method-based routing and
path parameters. Before 1.22 you needed third-party routers for these
features. Now the stdlib is sufficient for most applications.

Key behaviors:
- Patterns are matched most-specific-first
- A trailing slash "/path/" matches all paths under that prefix
- A pattern without trailing slash "/path" matches exactly that path
- Go 1.22+ supports "GET /path" for method-specific routes
- Go 1.22+ supports "/path/{param}" for path parameters

We'll cover routing in depth in Module 15. For now, just know that ServeMux
is the thing that decides which handler gets each request.

=============================================================================
*/

// DemoServeMux shows how to set up a basic router with multiple handlers.
func DemoServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register a handler function (most common pattern)
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Hello from GET /hello")
	})

	// Register a struct handler (for handlers that need state)
	greeting := &GreetingHandler{DefaultName: "World"}
	mux.Handle("GET /greet", greeting)

	return mux
}

/*
=============================================================================
 http.Request — WHAT THE CLIENT SENT
=============================================================================

The *http.Request struct contains everything about the incoming request.
Here are the fields you'll use constantly:

  r.Method          string              "GET", "POST", "PUT", "DELETE", etc.
  r.URL             *url.URL            Parsed URL with Path, Query, Fragment
  r.URL.Path        string              Just the path: "/users/42"
  r.URL.Query()     url.Values          Parsed query params: ?key=value
  r.Header          http.Header         Request headers (case-insensitive keys)
  r.Body            io.ReadCloser       Request body (must close when done!)
  r.Context()       context.Context     Request context (for cancellation, values)
  r.Form            url.Values          Parsed form data (call ParseForm first)
  r.Host            string              Host header value
  r.RemoteAddr      string              Client IP:port

Important: r.Body is an io.ReadCloser. You can only read it once! If you
need the body in multiple places, read it into a []byte first:

  body, err := io.ReadAll(r.Body)
  defer r.Body.Close()

Also important: r.Body is NEVER nil on the server side. For requests with
no body (like GET), reading it returns EOF immediately.

=============================================================================
*/

// DemoRequestInspection shows how to extract information from an HTTP request.
func DemoRequestInspection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract various request properties
	info := map[string]any{
		"method":      r.Method,
		"path":        r.URL.Path,
		"host":        r.Host,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.Header.Get("User-Agent"),
	}

	// Extract query parameters
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		queryParams[key] = strings.Join(values, ", ")
	}
	info["query_params"] = queryParams

	json.NewEncoder(w).Encode(info)
}

/*
=============================================================================
 http.ResponseWriter — WHAT YOU SEND BACK
=============================================================================

ResponseWriter is an interface with three methods:

  type ResponseWriter interface {
      Header() http.Header           // Access response headers
      Write([]byte) (int, error)     // Write the response body
      WriteHeader(statusCode int)    // Set the status code
  }

THE ORDER MATTERS. This is the #1 source of bugs for Go HTTP beginners:

  1. Set headers with w.Header().Set(...)
  2. Call w.WriteHeader(statusCode) to set the status code
  3. Call w.Write(body) to write the response body

If you call Write() before WriteHeader(), Go automatically sends a 200 OK
status. If you try to set headers after Write(), they're silently ignored.

Another critical gotcha: ResponseWriter is a write-only interface. There's
no way to read back what you've already written. If you need to inspect
the response (for logging, for example), you need a wrapper — which is
exactly what the httptest.ResponseRecorder is, and what we'll build in
the middleware module.

Also: after you write an error response, YOU MUST RETURN. Go doesn't stop
executing your handler just because you wrote a 400 or 500 response:

  if err != nil {
      http.Error(w, "bad request", http.StatusBadRequest)
      return  // <-- Without this, your handler keeps running!
  }

=============================================================================
*/

// DemoResponseWriter shows the correct order of operations for HTTP responses.
func DemoResponseWriter(w http.ResponseWriter, r *http.Request) {
	// Step 1: Set headers FIRST (before any Write call)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Custom-Header", "custom-value")
	w.Header().Set("Cache-Control", "no-cache")

	// Step 2: Set status code (optional — defaults to 200)
	w.WriteHeader(http.StatusCreated) // 201

	// Step 3: Write the body
	// After this point, headers and status code cannot be changed.
	response := map[string]string{
		"status":  "created",
		"message": "Resource was created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

/*
=============================================================================
 http.Server — THE PRODUCTION SERVER
=============================================================================

http.ListenAndServe is fine for quick prototypes, but for production you
want http.Server, which gives you control over timeouts, TLS, graceful
shutdown, and more:

  server := &http.Server{
      Addr:         ":8080",
      Handler:      mux,
      ReadTimeout:  5 * time.Second,
      WriteTimeout: 10 * time.Second,
      IdleTimeout:  120 * time.Second,
  }

Why timeouts matter: without them, a slow client can hold your server's
goroutine hostage indefinitely. Go spawns a goroutine per connection, and
without timeouts you're vulnerable to Slowloris-style attacks where
attackers open connections and trickle data to exhaust your resources.

  ReadTimeout:  Max time to read the entire request (headers + body)
  WriteTimeout: Max time to write the response
  IdleTimeout:  Max time to wait for the next request on a keep-alive connection

A good rule of thumb: set ReadTimeout < WriteTimeout, and set IdleTimeout
to something reasonable (60-120 seconds). These should be tuned based on
your specific use case.

=============================================================================
*/

// DemoProductionServer shows how to configure a production-ready HTTP server.
// This returns the server without starting it — in production you'd call
// server.ListenAndServe() or server.ListenAndServeTLS().
func DemoProductionServer(mux http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,

		// Timeouts prevent resource exhaustion from slow or malicious clients.
		ReadTimeout:  5 * time.Second,  // max time to read request
		WriteTimeout: 10 * time.Second, // max time to write response
		IdleTimeout:  120 * time.Second, // max time between requests on keep-alive

		// MaxHeaderBytes limits the size of request headers.
		// Default is 1MB. For most APIs, 1MB is more than enough.
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

/*
=============================================================================
 REQUEST LIFECYCLE
=============================================================================

Understanding the full lifecycle of an HTTP request in Go helps you debug
issues and write better handlers:

  1. Client sends HTTP request
  2. Go's HTTP server accepts the TCP connection
  3. A new goroutine is spawned for this connection
  4. The request is parsed into an *http.Request
  5. ServeMux matches the request to a handler
  6. The handler's ServeHTTP is called
  7. The handler writes to ResponseWriter
  8. The response is sent to the client
  9. The goroutine waits for the next request (keep-alive) or exits

Key implications:
- Each request runs in its own goroutine — handlers MUST be safe for
  concurrent access. If your handler reads/writes shared state, you
  need synchronization (mutexes, channels, etc.).
- The goroutine persists for the connection lifetime, not just the
  request. HTTP keep-alive means multiple requests share a goroutine.
- Context cancellation: if the client disconnects, r.Context() is
  canceled. Check this in long-running handlers to avoid wasted work.

=============================================================================
*/

// RequestCounter is a handler that demonstrates the need for concurrency
// safety. Since each request runs in its own goroutine, the counter
// must be protected by a mutex.
type RequestCounter struct {
	mu    sync.Mutex
	count int
}

// ServeHTTP increments the counter and returns the current count.
// The mutex ensures correct behavior even under concurrent requests.
func (rc *RequestCounter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rc.mu.Lock()
	rc.count++
	current := rc.count
	rc.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": current})
}

// GetCount returns the current count safely.
func (rc *RequestCounter) GetCount() int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.count
}

/*
=============================================================================
 COMMON MISTAKES AND PRODUCTION TIPS
=============================================================================

1. Writing after WriteHeader / Writing headers after Write:
   Once you call Write() or WriteHeader(), the HTTP status line and headers
   are flushed to the network. Any subsequent Header().Set() calls are
   silently ignored. This is the single most common HTTP bug in Go.

2. Not returning after error responses:
   http.Error() writes an error response but doesn't stop your handler.
   Always follow it with return.

3. Forgetting to close the request body:
   For handlers that read r.Body, always defer r.Body.Close(). The server
   handles this for you in most cases, but being explicit is good practice.

4. Not setting Content-Type:
   Go's http.DetectContentType sniffs the first 512 bytes if you don't set
   Content-Type explicitly. This can lead to surprising results — always
   set it yourself.

5. Goroutine leaks with background work:
   If your handler spawns goroutines, they outlive the request. Use context
   to bound their lifetime:

     go func(ctx context.Context) {
         select {
         case <-ctx.Done():
             return
         case result := <-doWork():
             // process result
         }
     }(r.Context())

6. Using the default ServeMux in production:
   http.HandleFunc registers on the default global mux. This is fine for
   examples but bad for production — use http.NewServeMux() instead. The
   global mux is a shared mutable state that any imported package can
   register handlers on (some packages do this as a side effect of import!).

=============================================================================
*/

// DemoCommonMistakes shows patterns you'll see (and should avoid) in the wild.
func DemoCommonMistakes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// WRONG: Not returning after error
		// if err := validateRequest(r); err != nil {
		//     http.Error(w, "bad request", 400)
		//     // handler continues executing! Bugs ensue.
		// }

		// RIGHT: Always return after error responses
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return // <-- Critical!
		}

		// WRONG: Setting headers after Write
		// w.Write([]byte("hello"))
		// w.Header().Set("X-Custom", "too-late") // silently ignored

		// RIGHT: Headers first, then write
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Custom", "this-works")
		fmt.Fprintln(w, "Headers were set correctly before writing.")
	}
}

/*
=============================================================================
 READING REQUEST BODIES
=============================================================================

For POST, PUT, and PATCH requests, the client sends data in the request
body. In Go, r.Body is an io.ReadCloser — it's a stream that you can only
read once.

For JSON APIs (which is most modern APIs), you'll use json.NewDecoder to
decode the body directly into a struct. This is more efficient than reading
the entire body into memory first because it streams the JSON parsing.

For forms (HTML form submissions), you'll call r.ParseForm() first, then
read values from r.Form or r.PostForm.

Important: always validate and limit the size of request bodies. An
attacker could send a multi-gigabyte body to exhaust your server's memory.
Use http.MaxBytesReader to limit body size:

  r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1 MB limit

=============================================================================
*/

// DemoBodyReading shows how to safely read and parse a JSON request body.
func DemoBodyReading(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit body size to prevent abuse (1 MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	// Read the body
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// Echo it back
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

// DemoJSONDecoding shows the preferred way to decode JSON request bodies.
func DemoJSONDecoding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit body size
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	// Decode JSON directly from the body stream — no need to read into []byte first.
	// This is more memory-efficient for large payloads.
	var payload map[string]any
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // strict mode: reject unknown JSON keys
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with the decoded data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"received": payload,
		"message":  "JSON decoded successfully",
	})
}

/*
=============================================================================
 FORM HANDLING
=============================================================================

While JSON APIs dominate modern development, you'll still encounter HTML
forms. Go handles them through r.ParseForm() and r.ParseMultipartForm().

After calling ParseForm():
  r.Form         — contains BOTH URL query params and POST body params
  r.PostForm     — contains ONLY POST body params

This distinction matters when a form POSTs to a URL with query params.

=============================================================================
*/

// DemoFormHandling shows how to process HTML form submissions.
func DemoFormHandling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ParseForm must be called before accessing r.Form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// Now you can access form values
	name := r.FormValue("name")         // shortcut that calls ParseForm if needed
	email := r.PostFormValue("email")    // only from POST body, not query params

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"name":  name,
		"email": email,
	})
}

/*
=============================================================================
 PUTTING IT ALL TOGETHER
=============================================================================

Here's a pattern you'll see in real Go applications: a function that
creates and configures the entire server, ready to start.

Notice how we:
- Use http.NewServeMux() instead of the default mux
- Create handler structs with their dependencies
- Register routes on the mux
- Return a configured *http.Server with proper timeouts

This pattern makes the server testable — you can create the server in
tests and use httptest without actually listening on a port.

=============================================================================
*/

// NewDemoServer creates a fully configured HTTP server demonstrating all
// the concepts covered in this module. It does NOT start the server —
// call server.ListenAndServe() to start it.
func NewDemoServer() *http.Server {
	mux := http.NewServeMux()

	// Simple function handler
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Welcome to the Go HTTP demo server!")
	})

	// Struct handler with state
	counter := &RequestCounter{}
	mux.Handle("GET /count", counter)

	// JSON handler
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	})

	// Request inspector
	mux.HandleFunc("GET /inspect", DemoRequestInspection)

	// Body echo
	mux.HandleFunc("POST /echo", DemoBodyReading)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown would be set up here in production
	// using signal.NotifyContext and server.Shutdown(ctx)

	log.Printf("Server configured on %s", server.Addr)
	return server
}
