// Package apitesting teaches how to test Go web services effectively —
// from unit testing individual handlers to integration testing complete
// API flows, using the powerful net/http/httptest standard library package.
package apitesting

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
=============================================================================
 API TESTING IN GO
=============================================================================

Go has one of the best standard libraries for testing HTTP services.
The net/http/httptest package gives you everything you need to test
handlers without starting a real server or making real network calls.

Why Go's approach is great:
  1. httptest.NewRecorder() — captures response without a network
  2. httptest.NewRequest() — creates request without a network
  3. httptest.NewServer() — spins up a real test server on localhost
  4. Everything works with the standard http.Handler interface

This means you can test:
  - Individual handlers (unit tests) with Recorder
  - Full middleware chains (integration tests) with Recorder or Server
  - Client code (end-to-end tests) with Server

The testing pyramid for APIs:
  Many:   Handler unit tests (fast, isolated, table-driven)
  Some:   Integration tests (middleware + handler together)
  Few:    End-to-end tests (real HTTP, real database)

=============================================================================
 httptest.NewRecorder: THE WORKHORSE
=============================================================================

httptest.NewRecorder returns an httptest.ResponseRecorder that implements
http.ResponseWriter. You pass it to your handler, and it captures everything
the handler writes — status code, headers, body.

  func TestMyHandler(t *testing.T) {
      // Create a request
      req := httptest.NewRequest("GET", "/users/123", nil)

      // Create a recorder (fake ResponseWriter)
      rr := httptest.NewRecorder()

      // Call the handler directly — no server needed!
      MyHandler(rr, req)

      // Inspect what the handler wrote
      if rr.Code != http.StatusOK {
          t.Errorf("expected 200, got %d", rr.Code)
      }

      body := rr.Body.String()
      // ... check the body
  }

This is the pattern you'll use for 90% of your API tests. It's fast
(no network I/O), deterministic, and easy to debug.

=============================================================================
*/

// --- Types used throughout the lesson and exercises ---

// User is a simple user model for our test examples.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    int               `json:"code"`
	Details []ValidationError `json:"details,omitempty"`
}

// ValidationError represents a field-level validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ListUsersResponse represents a paginated list of users.
type ListUsersResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// UserStore is a simple in-memory user store for testing examples.
type UserStore struct {
	users map[string]User
}

// NewUserStore creates a store pre-populated with test data.
func NewUserStore() *UserStore {
	return &UserStore{
		users: map[string]User{
			"user-1": {ID: "user-1", Name: "Alice", Email: "alice@example.com"},
			"user-2": {ID: "user-2", Name: "Bob", Email: "bob@example.com"},
		},
	}
}

// Get returns a user by ID.
func (s *UserStore) Get(id string) (User, bool) {
	u, ok := s.users[id]
	return u, ok
}

// List returns all users.
func (s *UserStore) List() []User {
	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

// Create adds a new user.
func (s *UserStore) Create(user User) {
	s.users[user.ID] = user
}

// Delete removes a user.
func (s *UserStore) Delete(id string) bool {
	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

/*
=============================================================================
 EXAMPLE HANDLERS (used in tests)
=============================================================================

These are simple handlers that demonstrate common patterns. The exercises
will have you test handlers like these.

=============================================================================
*/

// HandleGetUser is a handler that returns a user by ID.
// The ID is expected in the URL path as the last segment.
func HandleGetUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path (last segment)
		parts := strings.Split(r.URL.Path, "/")
		id := parts[len(parts)-1]

		user, ok := store.Get(id)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "user not found",
				Code:  http.StatusNotFound,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// HandleListUsers returns all users.
func HandleListUsers(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := store.List()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListUsersResponse{
			Users: users,
			Total: len(users),
		})
	}
}

// HandleCreateUser creates a new user from the request body.
func HandleCreateUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "invalid request body",
				Code:  http.StatusBadRequest,
			})
			return
		}

		// Validate
		var errors []ValidationError
		if user.Name == "" {
			errors = append(errors, ValidationError{Field: "name", Message: "name is required"})
		}
		if user.Email == "" || !strings.Contains(user.Email, "@") {
			errors = append(errors, ValidationError{Field: "email", Message: "invalid email format"})
		}
		if len(errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "validation failed",
				Code:    http.StatusBadRequest,
				Details: errors,
			})
			return
		}

		store.Create(user)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// HandleDeleteUser deletes a user by ID.
func HandleDeleteUser(store *UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		id := parts[len(parts)-1]

		if !store.Delete(id) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "user not found",
				Code:  http.StatusNotFound,
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

/*
=============================================================================
 EXAMPLE MIDDLEWARE
=============================================================================
*/

// LoggingMiddleware is a simple middleware that adds a header to track
// that it ran. Useful for testing that middleware chains work correctly.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Middleware", "logging")
		next.ServeHTTP(w, r)
	})
}

// RequireJSONMiddleware rejects requests without Content-Type: application/json.
func RequireJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				json.NewEncoder(w).Encode(ErrorResponse{
					Error: "Content-Type must be application/json",
					Code:  http.StatusUnsupportedMediaType,
				})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware checks for a Bearer token in the Authorization header.
// For testing purposes, it accepts any non-empty token and sets X-User
// header to "test-user".
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "unauthorized",
				Code:  http.StatusUnauthorized,
			})
			return
		}
		w.Header().Set("X-User", "test-user")
		next.ServeHTTP(w, r)
	})
}

/*
=============================================================================
 DemoRecorder shows how httptest.NewRecorder works.
=============================================================================
*/

// DemoRecorder demonstrates testing a handler with httptest.ResponseRecorder.
func DemoRecorder() {
	fmt.Println("=== httptest.NewRecorder ===")
	fmt.Println()

	store := NewUserStore()
	handler := HandleGetUser(store)

	// Test a successful request
	req := httptest.NewRequest("GET", "/users/user-1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	fmt.Printf("Status: %d\n", rr.Code)
	fmt.Printf("Content-Type: %s\n", rr.Header().Get("Content-Type"))
	fmt.Printf("Body: %s\n", rr.Body.String())

	// Test a not-found request
	req = httptest.NewRequest("GET", "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	fmt.Printf("\nNot Found Status: %d\n", rr.Code)
	fmt.Printf("Not Found Body: %s\n", rr.Body.String())
}

/*
=============================================================================
 httptest.NewServer: REAL HTTP FOR INTEGRATION TESTS
=============================================================================

When you need to test actual HTTP behavior (redirects, cookies, TLS,
connection handling), use httptest.NewServer. It starts a real HTTP
server on a random port on localhost.

  func TestWithServer(t *testing.T) {
      store := NewUserStore()
      mux := http.NewServeMux()
      mux.HandleFunc("/users", HandleListUsers(store))

      server := httptest.NewServer(mux)
      defer server.Close() // Always close!

      // Make real HTTP requests
      resp, err := http.Get(server.URL + "/users")
      if err != nil {
          t.Fatalf("request failed: %v", err)
      }
      defer resp.Body.Close()

      if resp.StatusCode != 200 {
          t.Errorf("expected 200, got %d", resp.StatusCode)
      }
  }

When to use NewServer vs NewRecorder:
  - NewRecorder: Testing handler logic (most tests)
  - NewServer: Testing HTTP behavior, client code, middleware chains

=============================================================================
*/

// DemoServer shows how httptest.NewServer works.
func DemoServer() {
	fmt.Println("=== httptest.NewServer ===")
	fmt.Println()

	store := NewUserStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/users/", HandleGetUser(store))

	server := httptest.NewServer(mux)
	defer server.Close()

	fmt.Printf("Test server running at: %s\n", server.URL)

	// Make a real HTTP request
	resp, err := http.Get(server.URL + "/users/user-1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", string(body))
}

/*
=============================================================================
 TABLE-DRIVEN API TESTS
=============================================================================

Table-driven tests are the Go way to test multiple scenarios. For API
tests, each row in the table specifies the request and expected response:

  tests := []struct {
      name           string
      method         string
      path           string
      body           string
      expectedStatus int
      expectedBody   string  // or a check function
  }{
      {"get existing user", "GET", "/users/1", "", 200, `{"id":"1"...}`},
      {"get missing user", "GET", "/users/99", "", 404, `{"error":...}`},
      {"create user", "POST", "/users", `{"name":"X"}`, 201, ""},
      {"bad request", "POST", "/users", "not json", 400, ""},
  }

  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
          rr := httptest.NewRecorder()
          handler.ServeHTTP(rr, req)
          if rr.Code != tt.expectedStatus {
              t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
          }
      })
  }

Benefits of table-driven API tests:
  - Easy to add new test cases (just add a row)
  - Easy to see all scenarios at a glance
  - Consistent testing approach
  - Each case runs as a subtest (t.Run) for clear output

=============================================================================
 TESTING MIDDLEWARE
=============================================================================

Middleware should be tested in isolation. Create a simple "next" handler
that records whether it was called and what it received:

  func TestAuthMiddleware(t *testing.T) {
      // Create a simple handler that records it was called
      called := false
      next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          called = true
          w.WriteHeader(200)
      })

      // Wrap with middleware
      handler := AuthMiddleware(next)

      // Test: without auth header → should NOT call next
      req := httptest.NewRequest("GET", "/", nil)
      rr := httptest.NewRecorder()
      handler.ServeHTTP(rr, req)
      if called {
          t.Error("middleware should block unauthenticated requests")
      }

      // Test: with auth header → should call next
      called = false
      req.Header.Set("Authorization", "Bearer valid-token")
      rr = httptest.NewRecorder()
      handler.ServeHTTP(rr, req)
      if !called {
          t.Error("middleware should pass authenticated requests")
      }
  }

=============================================================================
 GOLDEN FILE TESTING
=============================================================================

Golden file testing compares API responses against saved "golden" files.
The first time you run the test, you save the output as the golden file.
On subsequent runs, you compare against it.

This is great for:
  - Complex JSON responses (easier than writing assertions for every field)
  - Detecting unexpected changes in API output
  - Documenting what your API responses look like

The pattern:

  func TestGolden(t *testing.T) {
      // Get the actual response
      rr := httptest.NewRecorder()
      handler.ServeHTTP(rr, req)
      actual := rr.Body.Bytes()

      // Load the golden file
      golden := filepath.Join("testdata", "response.golden.json")
      expected, err := os.ReadFile(golden)

      // Compare (usually with normalization)
      if !bytes.Equal(normalize(actual), normalize(expected)) {
          t.Errorf("response doesn't match golden file")
      }
  }

To update golden files when the API intentionally changes:

  go test -run TestGolden -update

(You'd check for an -update flag and write the actual output to the file.)

=============================================================================
 TEST HELPERS
=============================================================================

Good test helpers reduce boilerplate and make tests more readable.
Common helpers for API testing:

  assertStatus(t, rr, http.StatusOK)
  assertJSON(t, rr)
  assertBodyContains(t, rr, "Alice")
  makeAuthRequest("GET", "/users", "token123")

Important: Test helpers should call t.Helper() so that test failures
report the caller's line number, not the helper's.

=============================================================================
 DEPENDENCY INJECTION FOR TESTING
=============================================================================

The key to testable handlers is dependency injection. Instead of handlers
reaching into global state, they receive their dependencies:

  // BAD: Handler uses global variable
  var db *sql.DB
  func GetUser(w http.ResponseWriter, r *http.Request) {
      db.Query(...) // Can't test without a real database!
  }

  // GOOD: Handler receives dependency via closure
  func GetUser(store UserStore) http.HandlerFunc {
      return func(w http.ResponseWriter, r *http.Request) {
          store.Get(...) // Can pass a fake store in tests!
      }
  }

  // GOOD: Handler receives dependency via interface
  type UserGetter interface {
      GetUser(id string) (User, error)
  }
  func GetUser(getter UserGetter) http.HandlerFunc { ... }

In tests, you pass a fake/mock/stub implementation:

  store := NewInMemoryStore() // Test double
  handler := GetUser(store)
  // ... test the handler

=============================================================================
 TEST COVERAGE: WHAT TO ACTUALLY TEST
=============================================================================

Don't aim for 100% coverage. Aim for confidence. Focus on:

  1. Happy paths (the thing works correctly)
  2. Error paths (bad input, missing resources, auth failures)
  3. Edge cases (empty lists, nil values, concurrent access)
  4. Security boundaries (auth middleware, authorization checks)

Things NOT worth testing:
  - Simple getters/setters
  - Third-party library behavior
  - Exact error message wording (test the status code instead)

Run coverage:
  go test -cover ./24-api-testing/
  go test -coverprofile=coverage.out ./24-api-testing/
  go tool cover -html=coverage.out  # Visual coverage report

=============================================================================
*/

// LoadGoldenFile reads a golden file from the testdata directory.
// This is a helper for golden file testing.
func LoadGoldenFile(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v", path, err)
	}
	return data
}

// NormalizeJSON re-marshals JSON to normalize whitespace and key order.
// This makes golden file comparisons reliable regardless of formatting.
func NormalizeJSON(t *testing.T, data []byte) []byte {
	t.Helper()
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nInput: %s", err, string(data))
	}
	normalized, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to re-marshal JSON: %v", err)
	}
	return normalized
}
