package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// okHandler is a simple handler that returns 200 OK with a JSON body.
func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
}

// statusHandler returns a handler that responds with the given status code.
func statusHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		fmt.Fprintf(w, "status: %d", code)
	})
}

// panicHandler always panics.
func panicHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
}

// TestLoggingMW tests Exercise 1: logging middleware.
func TestLoggingMW(t *testing.T) {
	t.Run("logs request details in headers", func(t *testing.T) {
		handler := LoggingMW(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		method := res.Header.Get("X-Log-Method")
		if method != "GET" {
			t.Errorf("X-Log-Method = %q, want %q", method, "GET")
		}

		path := res.Header.Get("X-Log-Path")
		if path != "/api/users" {
			t.Errorf("X-Log-Path = %q, want %q", path, "/api/users")
		}

		status := res.Header.Get("X-Log-Status")
		if status != "200" {
			t.Errorf("X-Log-Status = %q, want %q", status, "200")
		}
	})

	t.Run("captures non-200 status codes", func(t *testing.T) {
		handler := LoggingMW(statusHandler(http.StatusNotFound))

		req := httptest.NewRequest(http.MethodPost, "/missing", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		method := res.Header.Get("X-Log-Method")
		if method != "POST" {
			t.Errorf("X-Log-Method = %q, want %q", method, "POST")
		}

		status := res.Header.Get("X-Log-Status")
		if status != "404" {
			t.Errorf("X-Log-Status = %q, want %q", status, "404")
		}
	})
}

// TestRecoveryMW tests Exercise 2: panic recovery middleware.
func TestRecoveryMW(t *testing.T) {
	t.Run("recovers from panic", func(t *testing.T) {
		handler := RecoveryMW(panicHandler())

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()

		// This should NOT panic — the middleware catches it.
		// We use a deferred recover here so that if the middleware isn't
		// implemented yet (stub just passes through), the test fails
		// gracefully instead of crashing the test runner.
		panicked := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			handler.ServeHTTP(rec, req)
		}()
		if panicked {
			t.Fatal("RecoveryMW did not catch the panic — implement defer/recover in the middleware")
		}

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusInternalServerError)
		}

		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if result["error"] != "internal server error" {
			t.Errorf("error = %q, want %q", result["error"], "internal server error")
		}
	})

	t.Run("passes through normal requests", func(t *testing.T) {
		handler := RecoveryMW(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
	})
}

// TestRequestIDMW tests Exercise 3: request ID middleware.
func TestRequestIDMW(t *testing.T) {
	counter := 0
	generator := func() string {
		counter++
		return fmt.Sprintf("test-id-%d", counter)
	}

	t.Run("generates request ID when not provided", func(t *testing.T) {
		counter = 0
		handler := RequestIDMW(generator)(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		reqID := res.Header.Get("X-Request-Id")
		if reqID != "test-id-1" {
			t.Errorf("X-Request-Id = %q, want %q", reqID, "test-id-1")
		}
	})

	t.Run("uses client-provided request ID", func(t *testing.T) {
		counter = 0
		handler := RequestIDMW(generator)(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-Id", "client-provided-id")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		reqID := res.Header.Get("X-Request-Id")
		if reqID != "client-provided-id" {
			t.Errorf("X-Request-Id = %q, want %q", reqID, "client-provided-id")
		}
	})

	t.Run("adds request ID to context", func(t *testing.T) {
		counter = 0
		// Handler that reads the request ID from context
		contextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := GetRequestID(r.Context())
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, id)
		})

		handler := RequestIDMW(generator)(contextHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		if string(body) != "test-id-1" {
			t.Errorf("context request ID = %q, want %q", string(body), "test-id-1")
		}
	})
}

// TestCORSMW tests Exercise 4: CORS middleware.
func TestCORSMW(t *testing.T) {
	origins := []string{"https://example.com", "https://app.example.com"}
	methods := []string{"GET", "POST", "PUT"}

	t.Run("adds CORS headers for allowed origin", func(t *testing.T) {
		handler := CORSMW(origins, methods)(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		allowOrigin := res.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin != "https://example.com" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", allowOrigin, "https://example.com")
		}

		allowMethods := res.Header.Get("Access-Control-Allow-Methods")
		if allowMethods != "GET, POST, PUT" {
			t.Errorf("Access-Control-Allow-Methods = %q, want %q", allowMethods, "GET, POST, PUT")
		}
	})

	t.Run("handles preflight OPTIONS request", func(t *testing.T) {
		handler := CORSMW(origins, methods)(okHandler())

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNoContent)
		}

		allowOrigin := res.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin != "https://example.com" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", allowOrigin, "https://example.com")
		}
	})

	t.Run("no CORS headers for disallowed origin", func(t *testing.T) {
		handler := CORSMW(origins, methods)(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "https://evil.com")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		allowOrigin := res.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin != "" {
			t.Errorf("Access-Control-Allow-Origin = %q, want empty for disallowed origin", allowOrigin)
		}
	})

	t.Run("no CORS headers when no Origin header", func(t *testing.T) {
		handler := CORSMW(origins, methods)(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		allowOrigin := res.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin != "" {
			t.Errorf("Access-Control-Allow-Origin = %q, want empty when no Origin header", allowOrigin)
		}
	})
}

// TestSimpleRateLimitMW tests Exercise 5: rate limiting middleware.
func TestSimpleRateLimitMW(t *testing.T) {
	t.Run("allows requests within limit", func(t *testing.T) {
		handler := SimpleRateLimitMW(5, time.Second)(okHandler())

		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("request %d: status = %d, want %d", i+1, res.StatusCode, http.StatusOK)
			}
		}
	})

	t.Run("rejects requests over limit", func(t *testing.T) {
		handler := SimpleRateLimitMW(3, time.Hour)(okHandler())

		// Use up all tokens
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			rec.Result().Body.Close()
		}

		// This request should be rejected
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusTooManyRequests {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusTooManyRequests)
		}

		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		var result map[string]string
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if result["error"] != "rate limit exceeded" {
			t.Errorf("error = %q, want %q", result["error"], "rate limit exceeded")
		}
	})
}

// TestStatusCapture tests Exercise 6: ResponseWriter wrapper.
func TestStatusCapture(t *testing.T) {
	t.Run("captures status code", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sc := NewStatusCapture(rec)

		sc.WriteHeader(http.StatusNotFound)

		if sc.Code != http.StatusNotFound {
			t.Errorf("Code = %d, want %d", sc.Code, http.StatusNotFound)
		}
	})

	t.Run("captures first status code only", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sc := NewStatusCapture(rec)

		sc.WriteHeader(http.StatusNotFound)
		sc.WriteHeader(http.StatusOK) // second call shouldn't change Code

		if sc.Code != http.StatusNotFound {
			t.Errorf("Code = %d, want %d (should capture first call)", sc.Code, http.StatusNotFound)
		}
	})

	t.Run("defaults to 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sc := NewStatusCapture(rec)

		if sc.Code != http.StatusOK {
			t.Errorf("Code = %d, want %d (default)", sc.Code, http.StatusOK)
		}
	})

	t.Run("captures bytes written", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sc := NewStatusCapture(rec)

		data := []byte("Hello, World!")
		n, err := sc.Write(data)

		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
		if n != len(data) {
			t.Errorf("Write returned %d, want %d", n, len(data))
		}
		if sc.Written != len(data) {
			t.Errorf("Written = %d, want %d", sc.Written, len(data))
		}

		// Verify the data was forwarded to the underlying writer
		body := rec.Body.String()
		if body != "Hello, World!" {
			t.Errorf("body = %q, want %q", body, "Hello, World!")
		}
	})

	t.Run("accumulates bytes across multiple writes", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sc := NewStatusCapture(rec)

		sc.Write([]byte("Hello, "))
		sc.Write([]byte("World!"))

		if sc.Written != 13 {
			t.Errorf("Written = %d, want 13", sc.Written)
		}
	})

	t.Run("implicit WriteHeader on first Write", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sc := NewStatusCapture(rec)

		sc.Write([]byte("data"))

		if sc.Code != http.StatusOK {
			t.Errorf("Code = %d, want %d after implicit WriteHeader", sc.Code, http.StatusOK)
		}
		if !sc.wroteHeader {
			t.Error("wroteHeader should be true after Write")
		}
	})
}

// TestTimeoutMW tests Exercise 7: timeout middleware.
func TestTimeoutMW(t *testing.T) {
	t.Run("allows fast requests through", func(t *testing.T) {
		handler := TimeoutMW(time.Second)(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("times out slow requests", func(t *testing.T) {
		slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-time.After(5 * time.Second):
				w.Write([]byte("done"))
			case <-r.Context().Done():
				return
			}
		})

		handler := TimeoutMW(50 * time.Millisecond)(slowHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusServiceUnavailable)
		}

		body, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(body), "request timeout") {
			t.Errorf("body = %q, should contain %q", string(body), "request timeout")
		}
	})
}

// TestChainMiddleware tests Exercise 8: middleware chaining.
func TestChainMiddleware(t *testing.T) {
	t.Run("applies middlewares in correct order", func(t *testing.T) {
		// Each middleware adds a header showing the order of execution
		var order []string

		mw1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw1-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw1-after")
			})
		}

		mw2 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw2-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw2-after")
			})
		}

		mw3 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw3-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw3-after")
			})
		}

		innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "handler")
			w.WriteHeader(http.StatusOK)
		})

		handler := ChainMiddleware(innerHandler, mw1, mw2, mw3)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		// mw1 should be outermost (first to run)
		expected := []string{
			"mw1-before",
			"mw2-before",
			"mw3-before",
			"handler",
			"mw3-after",
			"mw2-after",
			"mw1-after",
		}

		if len(order) != len(expected) {
			t.Fatalf("got %d calls, want %d: %v", len(order), len(expected), order)
		}

		for i, want := range expected {
			if order[i] != want {
				t.Errorf("order[%d] = %q, want %q (full order: %v)", i, order[i], want, order)
			}
		}
	})

	t.Run("works with no middlewares", func(t *testing.T) {
		handler := ChainMiddleware(okHandler())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("single middleware works", func(t *testing.T) {
		headerMW := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware", "applied")
				next.ServeHTTP(w, r)
			})
		}

		handler := ChainMiddleware(okHandler(), headerMW)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.Header.Get("X-Middleware") != "applied" {
			t.Errorf("X-Middleware = %q, want %q", res.Header.Get("X-Middleware"), "applied")
		}
	})
}
