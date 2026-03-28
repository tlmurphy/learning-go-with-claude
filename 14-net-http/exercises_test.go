package nethttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// TestHelloHandler tests Exercise 1: basic handler with content type and body.
func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantBody       string
		wantType       string
		wantStatusCode int
	}{
		{
			name:           "returns hello world",
			method:         http.MethodGet,
			wantBody:       "Hello, World!",
			wantType:       "text/plain; charset=utf-8",
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/hello", nil)
			rec := httptest.NewRecorder()

			HelloHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("status code = %d, want %d", res.StatusCode, tt.wantStatusCode)
			}

			body, _ := io.ReadAll(res.Body)
			if string(body) != tt.wantBody {
				t.Errorf("body = %q, want %q", string(body), tt.wantBody)
			}

			ct := res.Header.Get("Content-Type")
			if ct != tt.wantType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantType)
			}
		})
	}
}

// TestQueryParamHandler tests Exercise 2: query parameter extraction.
func TestQueryParamHandler(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		wantName string
		wantAge  string
	}{
		{
			name:     "both params provided",
			query:    "?name=Alice&age=30",
			wantName: "Alice",
			wantAge:  "30",
		},
		{
			name:     "name only",
			query:    "?name=Bob",
			wantName: "Bob",
			wantAge:  "0",
		},
		{
			name:     "age only",
			query:    "?age=25",
			wantName: "anonymous",
			wantAge:  "25",
		},
		{
			name:     "no params",
			query:    "",
			wantName: "anonymous",
			wantAge:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/info"+tt.query, nil)
			rec := httptest.NewRecorder()

			QueryParamHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			ct := res.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			var result map[string]string
			if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode JSON response: %v", err)
			}

			if result["name"] != tt.wantName {
				t.Errorf("name = %q, want %q", result["name"], tt.wantName)
			}
			if result["age"] != tt.wantAge {
				t.Errorf("age = %q, want %q", result["age"], tt.wantAge)
			}
		})
	}
}

// TestEchoBodyHandler tests Exercise 3: request body reading.
func TestEchoBodyHandler(t *testing.T) {
	t.Run("echoes POST body", func(t *testing.T) {
		body := "Hello, this is the request body!"
		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
		rec := httptest.NewRecorder()

		EchoBodyHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusOK)
		}

		ct := res.Header.Get("Content-Type")
		if ct != "application/octet-stream" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/octet-stream")
		}

		responseBody, _ := io.ReadAll(res.Body)
		if string(responseBody) != body {
			t.Errorf("body = %q, want %q", string(responseBody), body)
		}
	})

	t.Run("rejects non-POST methods", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/echo", nil)
		rec := httptest.NewRecorder()

		EchoBodyHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusMethodNotAllowed)
		}
	})

	t.Run("handles empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(""))
		rec := httptest.NewRecorder()

		EchoBodyHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusOK)
		}

		responseBody, _ := io.ReadAll(res.Body)
		if string(responseBody) != "" {
			t.Errorf("body = %q, want empty", string(responseBody))
		}
	})
}

// TestCustomHeaderHandler tests Exercise 4: custom headers and status codes.
func TestCustomHeaderHandler(t *testing.T) {
	t.Run("sets custom headers and status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/custom", nil)
		rec := httptest.NewRecorder()

		CustomHeaderHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusAccepted {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusAccepted)
		}

		expectedHeaders := map[string]string{
			"X-Request-Id": "12345",
			"X-Powered-By": "Go",
			"Cache-Control": "no-store",
			"Content-Type":  "application/json",
		}

		for header, want := range expectedHeaders {
			got := res.Header.Get(header)
			if got != want {
				t.Errorf("header %s = %q, want %q", header, got, want)
			}
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if result["status"] != "accepted" {
			t.Errorf("status = %q, want %q", result["status"], "accepted")
		}
	})
}

// TestHealthCheckHandler tests Exercise 5: health check endpoint.
func TestHealthCheckHandler(t *testing.T) {
	t.Run("returns healthy status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		HealthCheckHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusOK)
		}

		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if result["status"] != "healthy" {
			t.Errorf("status = %q, want %q", result["status"], "healthy")
		}
		if result["version"] != "1.0.0" {
			t.Errorf("version = %q, want %q", result["version"], "1.0.0")
		}
	})

	t.Run("rejects non-GET methods", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		rec := httptest.NewRecorder()

		HealthCheckHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusMethodNotAllowed)
		}
	})
}

// TestContentNegotiationHandler tests Exercise 6: Accept header handling.
func TestContentNegotiationHandler(t *testing.T) {
	tests := []struct {
		name       string
		accept     string
		wantType   string
		wantBody   string
		isJSON     bool
	}{
		{
			name:     "returns JSON for application/json",
			accept:   "application/json",
			wantType: "application/json",
			isJSON:   true,
		},
		{
			name:     "returns plain text for text/plain",
			accept:   "text/plain",
			wantType: "text/plain; charset=utf-8",
			wantBody: "hello",
			isJSON:   false,
		},
		{
			name:     "returns plain text for empty accept",
			accept:   "",
			wantType: "text/plain; charset=utf-8",
			wantBody: "hello",
			isJSON:   false,
		},
		{
			name:     "returns plain text for unknown accept",
			accept:   "text/html",
			wantType: "text/plain; charset=utf-8",
			wantBody: "hello",
			isJSON:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/content", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			rec := httptest.NewRecorder()

			ContentNegotiationHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			ct := res.Header.Get("Content-Type")
			if ct != tt.wantType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantType)
			}

			if tt.isJSON {
				var result map[string]string
				if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
					t.Fatalf("failed to decode JSON: %v", err)
				}
				if result["message"] != "hello" {
					t.Errorf("message = %q, want %q", result["message"], "hello")
				}
			} else {
				body, _ := io.ReadAll(res.Body)
				if strings.TrimSpace(string(body)) != tt.wantBody {
					t.Errorf("body = %q, want %q", strings.TrimSpace(string(body)), tt.wantBody)
				}
			}
		})
	}
}

// TestVisitCounter tests Exercise 7: struct-based handler with state.
func TestVisitCounter(t *testing.T) {
	t.Run("increments count on each request", func(t *testing.T) {
		counter := &VisitCounter{}

		for i := 1; i <= 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/visits", nil)
			rec := httptest.NewRecorder()

			counter.ServeHTTP(rec, req)

			res := rec.Result()

			ct := res.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("request %d: Content-Type = %q, want %q", i, ct, "application/json")
			}

			var result map[string]int
			if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
				t.Fatalf("request %d: failed to decode JSON: %v", i, err)
			}
			res.Body.Close()

			if result["visits"] != i {
				t.Errorf("request %d: visits = %d, want %d", i, result["visits"], i)
			}
		}

		// Verify final count through the accessor
		if counter.CurrentCount() != 5 {
			t.Errorf("CurrentCount() = %d, want 5", counter.CurrentCount())
		}
	})

	t.Run("handles concurrent requests safely", func(t *testing.T) {
		counter := &VisitCounter{}
		numRequests := 100

		var wg sync.WaitGroup
		wg.Add(numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodGet, "/visits", nil)
				rec := httptest.NewRecorder()
				counter.ServeHTTP(rec, req)
			}()
		}

		wg.Wait()

		if counter.CurrentCount() != numRequests {
			t.Errorf("after %d concurrent requests, count = %d", numRequests, counter.CurrentCount())
		}
	})
}

// TestFormValidationHandler tests Exercise 8: form data with validation.
func TestFormValidationHandler(t *testing.T) {
	t.Run("accepts valid form data", func(t *testing.T) {
		form := url.Values{}
		form.Set("username", "alice")
		form.Set("email", "alice@example.com")

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		FormValidationHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if result["username"] != "alice" {
			t.Errorf("username = %q, want %q", result["username"], "alice")
		}
		if result["email"] != "alice@example.com" {
			t.Errorf("email = %q, want %q", result["email"], "alice@example.com")
		}
	})

	t.Run("rejects missing username", func(t *testing.T) {
		form := url.Values{}
		form.Set("email", "alice@example.com")

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		FormValidationHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if result["error"] != "missing required field: username" {
			t.Errorf("error = %q, want %q", result["error"], "missing required field: username")
		}
	})

	t.Run("rejects missing email", func(t *testing.T) {
		form := url.Values{}
		form.Set("username", "alice")

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		FormValidationHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		if result["error"] != "missing required field: email" {
			t.Errorf("error = %q, want %q", result["error"], "missing required field: email")
		}
	})

	t.Run("rejects non-POST methods", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/register", nil)
		rec := httptest.NewRecorder()

		FormValidationHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusMethodNotAllowed)
		}
	})

	t.Run("rejects empty field values", func(t *testing.T) {
		form := url.Values{}
		form.Set("username", "")
		form.Set("email", "alice@example.com")

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		FormValidationHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("status code = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}
	})
}
