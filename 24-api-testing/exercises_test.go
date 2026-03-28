package apitesting

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// =============================================================================
// Exercise 1: Table-Driven API Test Cases
// =============================================================================

func TestTableDrivenCRUD(t *testing.T) {
	store := NewUserStore()

	// Build a simple mux for testing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", HandleGetUser(store))
	mux.HandleFunc("GET /users", HandleListUsers(store))
	mux.HandleFunc("POST /users", HandleCreateUser(store))
	mux.HandleFunc("DELETE /users/{id}", HandleDeleteUser(store))

	t.Run("RunAPITest single case", func(t *testing.T) {
		tc := APITestCase{
			Name:                 "get existing user",
			Method:               "GET",
			Path:                 "/users/user-1",
			ExpectedStatus:       http.StatusOK,
			ExpectedBodyContains: "Alice",
		}

		// RunAPITest should not cause the test to fail for this valid case
		RunAPITest(t, mux, tc)
	})

	t.Run("RunAPITests multiple cases", func(t *testing.T) {
		tests := []APITestCase{
			{
				Name:                 "get existing user",
				Method:               "GET",
				Path:                 "/users/user-1",
				ExpectedStatus:       http.StatusOK,
				ExpectedBodyContains: "Alice",
			},
			{
				Name:           "get missing user",
				Method:         "GET",
				Path:           "/users/nonexistent",
				ExpectedStatus: http.StatusNotFound,
			},
			{
				Name:   "create user",
				Method: "POST",
				Path:   "/users",
				Body:   `{"id":"user-new","name":"New User","email":"new@test.com"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				ExpectedStatus:       http.StatusCreated,
				ExpectedBodyContains: "New User",
			},
		}

		RunAPITests(t, mux, tests)
	})

	t.Run("RunAPITest sets headers", func(t *testing.T) {
		// The test case should set Content-Type header
		tc := APITestCase{
			Name:   "create with json content type",
			Method: "POST",
			Path:   "/users",
			Body:   `{"id":"user-h","name":"Header Test","email":"h@test.com"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			ExpectedStatus: http.StatusCreated,
		}
		RunAPITest(t, mux, tc)
	})
}

// =============================================================================
// Exercise 2: Test Middleware in Isolation
// =============================================================================

func TestTestMiddleware(t *testing.T) {
	t.Run("middleware that passes through", func(t *testing.T) {
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Custom", "added")
				next.ServeHTTP(w, r)
			})
		}

		req := httptest.NewRequest("GET", "/test", nil)
		result := TestMiddleware(middleware, req)

		if !result.NextCalled {
			t.Error("NextCalled should be true when middleware calls next. " +
				"Track whether your recording handler was invoked.")
		}
		if result.Headers.Get("X-Custom") != "added" {
			t.Error("Headers should include headers set by middleware. " +
				"Capture the response headers from the recorder.")
		}
	})

	t.Run("middleware that blocks", func(t *testing.T) {
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"forbidden"}`))
				// Note: does NOT call next
			})
		}

		req := httptest.NewRequest("GET", "/test", nil)
		result := TestMiddleware(middleware, req)

		if result.NextCalled {
			t.Error("NextCalled should be false when middleware blocks the request.")
		}
		if result.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d. " +
				"Capture the status code from the recorder.", result.StatusCode)
		}
		if !strings.Contains(result.Body, "forbidden") {
			t.Error("Body should contain the response written by middleware.")
		}
	})

	t.Run("middleware modifies request headers", func(t *testing.T) {
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Header.Set("X-Injected", "by-middleware")
				next.ServeHTTP(w, r)
			})
		}

		req := httptest.NewRequest("GET", "/test", nil)
		result := TestMiddleware(middleware, req)

		if !result.NextCalled {
			t.Error("NextCalled should be true.")
		}
		if result.RequestHeaders.Get("X-Injected") != "by-middleware" {
			t.Error("RequestHeaders should capture headers as seen by the next handler. " +
				"Record r.Header inside your next handler.")
		}
	})
}

// =============================================================================
// Exercise 3: Test Helper Functions
// =============================================================================

func TestAssertHelpers(t *testing.T) {
	t.Run("AssertStatus passes on match", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.WriteHeader(http.StatusOK)

		// This should not fail
		mt := &mockT{}
		AssertStatus(mt, rr, http.StatusOK)
		if mt.failed {
			t.Error("AssertStatus should not fail when status matches.")
		}
	})

	t.Run("AssertStatus fails on mismatch", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.WriteHeader(http.StatusNotFound)

		mt := &mockT{}
		AssertStatus(mt, rr, http.StatusOK)
		if !mt.failed {
			t.Error("AssertStatus should fail when status doesn't match. " +
				"Compare rr.Code with the expected status and call t.Errorf on mismatch.")
		}
	})

	t.Run("AssertJSON passes for JSON content type", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.Header().Set("Content-Type", "application/json")

		mt := &mockT{}
		AssertJSON(mt, rr)
		if mt.failed {
			t.Error("AssertJSON should pass when Content-Type is application/json.")
		}
	})

	t.Run("AssertJSON fails for non-JSON", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.Header().Set("Content-Type", "text/plain")

		mt := &mockT{}
		AssertJSON(mt, rr)
		if !mt.failed {
			t.Error("AssertJSON should fail when Content-Type is not application/json. " +
				"Check that Content-Type header starts with 'application/json'.")
		}
	})

	t.Run("AssertBodyContains passes when present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.Write([]byte(`{"name":"Alice","email":"alice@test.com"}`))

		mt := &mockT{}
		AssertBodyContains(mt, rr, "Alice")
		if mt.failed {
			t.Error("AssertBodyContains should pass when the body contains the expected string.")
		}
	})

	t.Run("AssertBodyContains fails when absent", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.Write([]byte(`{"name":"Bob"}`))

		mt := &mockT{}
		AssertBodyContains(mt, rr, "Alice")
		if !mt.failed {
			t.Error("AssertBodyContains should fail when the body doesn't contain the expected string. " +
				"Use strings.Contains to check.")
		}
	})

	t.Run("AssertHeader passes on match", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.Header().Set("X-Custom", "value")

		mt := &mockT{}
		AssertHeader(mt, rr, "X-Custom", "value")
		if mt.failed {
			t.Error("AssertHeader should pass when header matches.")
		}
	})

	t.Run("AssertHeader fails on mismatch", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rr.Header().Set("X-Custom", "wrong")

		mt := &mockT{}
		AssertHeader(mt, rr, "X-Custom", "expected")
		if !mt.failed {
			t.Error("AssertHeader should fail when header doesn't match.")
		}
	})
}

func TestMakeRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"method": r.Method, "path": r.URL.Path})
	})

	t.Run("makes GET request", func(t *testing.T) {
		rr := MakeRequest(t, handler, "GET", "/test", "")
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d. "+
				"MakeRequest should create a request, recorder, call the handler, and return the recorder.",
				rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "GET") {
			t.Error("Request method should be GET.")
		}
	})

	t.Run("makes POST request with body", func(t *testing.T) {
		bodyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			w.Write(body)
		})

		rr := MakeRequest(t, bodyHandler, "POST", "/test", `{"key":"value"}`)
		if !strings.Contains(rr.Body.String(), "value") {
			t.Error("POST body should be sent to the handler. " +
				"Use strings.NewReader(body) for the request body.")
		}
	})
}

func TestMakeAuthRequest(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	}))

	t.Run("sends auth header", func(t *testing.T) {
		rr := MakeAuthRequest(t, handler, "GET", "/protected", "", "my-token")
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d. "+
				"MakeAuthRequest should set Authorization: Bearer <token>.", rr.Code)
		}
	})

	t.Run("missing token gets 401", func(t *testing.T) {
		rr := MakeRequest(t, handler, "GET", "/protected", "")
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 without token, got %d", rr.Code)
		}
	})
}

// =============================================================================
// Exercise 4: Integration Test with httptest.Server
// =============================================================================

func TestSetupTestServer(t *testing.T) {
	t.Run("basic server setup", func(t *testing.T) {
		ts := SetupTestServer()
		if ts == nil {
			t.Fatal("SetupTestServer should return a non-nil TestServer.")
		}
		if ts.Server == nil {
			t.Fatal("TestServer.Server should not be nil. "+
				"Use httptest.NewServer to create it.")
		}
		defer ts.Server.Close()

		if ts.Store == nil {
			t.Fatal("TestServer.Store should not be nil. " +
				"Create a UserStore and use it in your handlers.")
		}
	})

	t.Run("GET existing user", func(t *testing.T) {
		ts := SetupTestServer()
		if ts.Server == nil {
			t.Fatal("SetupTestServer must return a TestServer with a non-nil Server.")
		}
		defer ts.Server.Close()

		resp, err := http.Get(ts.Server.URL + "/users/user-1")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d. "+
				"Register HandleGetUser on a path like GET /users/{id}", resp.StatusCode)
		}

		var user User
		json.NewDecoder(resp.Body).Decode(&user)
		if user.Name != "Alice" {
			t.Errorf("Expected Alice, got %q", user.Name)
		}
	})

	t.Run("GET missing user returns 404", func(t *testing.T) {
		ts := SetupTestServer()
		if ts.Server == nil {
			t.Fatal("SetupTestServer must return a TestServer with a non-nil Server.")
		}
		defer ts.Server.Close()

		resp, err := http.Get(ts.Server.URL + "/users/nonexistent")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("POST creates user", func(t *testing.T) {
		ts := SetupTestServer()
		if ts.Server == nil {
			t.Fatal("SetupTestServer must return a TestServer with a non-nil Server.")
		}
		defer ts.Server.Close()

		body := strings.NewReader(`{"id":"user-new","name":"Charlie","email":"charlie@test.com"}`)
		resp, err := http.Post(ts.Server.URL+"/users", "application/json", body)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d. "+
				"Register HandleCreateUser on POST /users.", resp.StatusCode)
		}

		// Verify user was actually stored
		if _, ok := ts.Store.Get("user-new"); !ok {
			t.Error("Created user should be in the store.")
		}
	})

	t.Run("DELETE removes user", func(t *testing.T) {
		ts := SetupTestServer()
		if ts.Server == nil {
			t.Fatal("SetupTestServer must return a TestServer with a non-nil Server.")
		}
		defer ts.Server.Close()

		req, _ := http.NewRequest("DELETE", ts.Server.URL+"/users/user-1", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected 204, got %d. "+
				"Register HandleDeleteUser on DELETE /users/{id}.", resp.StatusCode)
		}
	})

	t.Run("server with middleware", func(t *testing.T) {
		ts := SetupTestServer(LoggingMiddleware)
		if ts.Server == nil {
			t.Fatal("SetupTestServer must return a TestServer with a non-nil Server.")
		}
		defer ts.Server.Close()

		resp, err := http.Get(ts.Server.URL + "/users/user-1")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.Header.Get("X-Middleware") != "logging" {
			t.Error("Middleware should be applied. "+
				"Wrap your mux with the provided middlewares before creating the server. "+
				"The LoggingMiddleware sets X-Middleware: logging.",
			)
		}
	})
}

// =============================================================================
// Exercise 5: Test Fixture Factory
// =============================================================================

func TestUserFactory(t *testing.T) {
	t.Run("creates user with defaults", func(t *testing.T) {
		f := NewUserFactory()
		user := f.MakeUser()

		if user.ID == "" {
			t.Error("MakeUser should generate an ID like 'user-1'. " +
				"Use the counter field to auto-increment.")
		}
		if user.Name == "" {
			t.Error("MakeUser should set a default Name like 'Test User 1'.")
		}
		if user.Email == "" {
			t.Error("MakeUser should set a default Email like 'user-1@test.com'.")
		}
	})

	t.Run("auto-increments IDs", func(t *testing.T) {
		f := NewUserFactory()
		u1 := f.MakeUser()
		u2 := f.MakeUser()
		u3 := f.MakeUser()

		if u1.ID == u2.ID || u2.ID == u3.ID {
			t.Errorf("Each user should have a unique ID. Got: %q, %q, %q. "+
				"Increment the counter in MakeUser.", u1.ID, u2.ID, u3.ID)
		}
	})

	t.Run("WithName overrides name", func(t *testing.T) {
		f := NewUserFactory()
		user := f.MakeUser(WithName("Custom Name"))

		if user.Name != "Custom Name" {
			t.Errorf("Expected name 'Custom Name', got %q. "+
				"Apply options after setting defaults.", user.Name)
		}
	})

	t.Run("WithEmail overrides email", func(t *testing.T) {
		f := NewUserFactory()
		user := f.MakeUser(WithEmail("custom@example.com"))

		if user.Email != "custom@example.com" {
			t.Errorf("Expected email 'custom@example.com', got %q", user.Email)
		}
	})

	t.Run("WithID overrides ID", func(t *testing.T) {
		f := NewUserFactory()
		user := f.MakeUser(WithID("custom-id"))

		if user.ID != "custom-id" {
			t.Errorf("Expected ID 'custom-id', got %q", user.ID)
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		f := NewUserFactory()
		user := f.MakeUser(WithName("Alice"), WithEmail("alice@test.com"))

		if user.Name != "Alice" {
			t.Errorf("Expected name 'Alice', got %q", user.Name)
		}
		if user.Email != "alice@test.com" {
			t.Errorf("Expected email 'alice@test.com', got %q", user.Email)
		}
	})

	t.Run("MakeUsers creates multiple", func(t *testing.T) {
		f := NewUserFactory()
		users := f.MakeUsers(5)

		if len(users) != 5 {
			t.Errorf("Expected 5 users, got %d. "+
				"Call MakeUser in a loop.", len(users))
		}

		// All should have unique IDs
		ids := make(map[string]bool)
		for _, u := range users {
			if ids[u.ID] {
				t.Errorf("Duplicate ID: %q", u.ID)
			}
			ids[u.ID] = true
		}
	})
}

// =============================================================================
// Exercise 6: Golden File Testing
// =============================================================================

func TestGoldenFile(t *testing.T) {
	store := NewUserStore()

	t.Run("get user matches golden file", func(t *testing.T) {
		handler := HandleGetUser(store)
		req := httptest.NewRequest("GET", "/users/user-1", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		CompareResponseWithGolden(t, "get_user_response.golden.json", rr)
	})

	t.Run("not found matches golden file", func(t *testing.T) {
		handler := HandleGetUser(store)
		req := httptest.NewRequest("GET", "/users/nonexistent", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		CompareResponseWithGolden(t, "error_not_found.golden.json", rr)
	})

	t.Run("create user matches golden file", func(t *testing.T) {
		testStore := NewUserStore()
		handler := HandleCreateUser(testStore)
		body := `{"id":"user-new","name":"Charlie","email":"charlie@example.com"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		CompareResponseWithGolden(t, "create_user_response.golden.json", rr)
	})

	t.Run("validation error matches golden file", func(t *testing.T) {
		handler := HandleCreateUser(store)
		body := `{"id":"bad","name":"","email":"not-an-email"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		CompareResponseWithGolden(t, "error_validation.golden.json", rr)
	})
}

// =============================================================================
// Exercise 7: Error Scenario Testing
// =============================================================================

func TestRunErrorTests(t *testing.T) {
	store := NewUserStore()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", HandleGetUser(store))
	mux.HandleFunc("POST /users", HandleCreateUser(store))

	// Wrap with auth middleware for testing unauthorized errors
	protectedMux := AuthMiddleware(mux)

	tests := []ErrorTestCase{
		{
			Name:           "missing auth header",
			Method:         "GET",
			Path:           "/users/user-1",
			ExpectedStatus: http.StatusUnauthorized,
			ExpectedError:  "unauthorized",
		},
		{
			Name:   "not found",
			Method: "GET",
			Path:   "/users/nonexistent",
			Headers: map[string]string{
				"Authorization": "Bearer valid-token",
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "user not found",
		},
		{
			Name:   "invalid JSON body",
			Method: "POST",
			Path:   "/users",
			Body:   "not json",
			Headers: map[string]string{
				"Authorization": "Bearer valid-token",
				"Content-Type":  "application/json",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "invalid request body",
		},
	}

	RunErrorTests(t, protectedMux, tests)
}

// =============================================================================
// Exercise 8: Complete API Test Suite
// =============================================================================

func TestAPITestSuite(t *testing.T) {
	suite := NewAPITestSuite()
	if suite == nil {
		t.Fatal("NewAPITestSuite should return a non-nil suite.")
	}

	suite.Setup(t)
	if suite.server == nil {
		t.Fatal("Setup should initialize the server. " +
			"Use SetupTestServer to create it.")
	}
	if suite.factory == nil {
		t.Fatal("Setup should initialize the factory. " +
			"Use NewUserFactory to create it.")
	}
	suite.Teardown(t)

	// RunAll should execute all sub-tests without panicking
	t.Run("RunAll", func(t *testing.T) {
		suite.RunAll(t)
	})
}

// =============================================================================
// mockT: A mock testing.T for testing assertion helpers
// =============================================================================

// mockT satisfies the TB interface used by our assertion helpers.
// It captures whether the test "failed" without actually failing the real test.
type mockT struct {
	failed   bool
	messages []string
}

func (m *mockT) Helper() {}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}

func (m *mockT) Fatalf(format string, args ...interface{}) {
	m.failed = true
	m.messages = append(m.messages, fmt.Sprintf(format, args...))
}
