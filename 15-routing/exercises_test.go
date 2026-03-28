package routing

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper decodes a JSON response body into a map.
func decodeJSON(t *testing.T, body io.Reader) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	return result
}

// TestMethodRouter tests Exercise 1: method-based routing.
func TestMethodRouter(t *testing.T) {
	mux := MethodRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantMethod string
		wantAction string
	}{
		{"GET items", http.MethodGet, "/items", 200, "GET", "list"},
		{"POST items", http.MethodPost, "/items", 201, "POST", "create"},
		{"PUT items", http.MethodPut, "/items", 200, "PUT", "update"},
		{"DELETE items", http.MethodDelete, "/items", 200, "DELETE", "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", res.StatusCode, tt.wantStatus)
			}

			ct := res.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			result := decodeJSON(t, res.Body)
			if result["method"] != tt.wantMethod {
				t.Errorf("method = %q, want %q", result["method"], tt.wantMethod)
			}
			if result["action"] != tt.wantAction {
				t.Errorf("action = %q, want %q", result["action"], tt.wantAction)
			}
		})
	}
}

// TestPathParamExtractor tests Exercise 2: path parameter extraction.
func TestPathParamExtractor(t *testing.T) {
	mux := PathParamExtractor()

	tests := []struct {
		name         string
		path         string
		wantResource string
		wantParamKey string
		wantParamVal string
	}{
		{"user by id", "/users/42", "user", "id", "42"},
		{"user by string id", "/users/abc", "user", "id", "abc"},
		{"post by slug", "/posts/my-first-post", "post", "slug", "my-first-post"},
		{"post by numeric slug", "/posts/123", "post", "slug", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
			}

			result := decodeJSON(t, res.Body)
			if result["resource"] != tt.wantResource {
				t.Errorf("resource = %q, want %q", result["resource"], tt.wantResource)
			}
			if result[tt.wantParamKey] != tt.wantParamVal {
				t.Errorf("%s = %q, want %q", tt.wantParamKey, result[tt.wantParamKey], tt.wantParamVal)
			}
		})
	}
}

// TestResourceRouter tests Exercise 3: complete CRUD routing.
func TestResourceRouter(t *testing.T) {
	mux := ResourceRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantAction string
		wantID     string
	}{
		{"list books", http.MethodGet, "/books", 200, "list", ""},
		{"create book", http.MethodPost, "/books", 201, "create", ""},
		{"get book", http.MethodGet, "/books/1", 200, "get", "1"},
		{"update book", http.MethodPut, "/books/1", 200, "update", "1"},
		{"delete book", http.MethodDelete, "/books/1", 204, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", res.StatusCode, tt.wantStatus)
			}

			// DELETE returns no body
			if tt.wantStatus == 204 {
				return
			}

			result := decodeJSON(t, res.Body)
			if result["action"] != tt.wantAction {
				t.Errorf("action = %q, want %q", result["action"], tt.wantAction)
			}
			if tt.wantID != "" {
				if result["id"] != tt.wantID {
					t.Errorf("id = %q, want %q", result["id"], tt.wantID)
				}
			}
		})
	}
}

// TestWildcardRouter tests Exercise 4: catch-all wildcard routing.
func TestWildcardRouter(t *testing.T) {
	mux := WildcardRouter()

	tests := []struct {
		name     string
		path     string
		wantFile string
	}{
		{"single file", "/static/style.css", "style.css"},
		{"nested path", "/static/images/logo.png", "images/logo.png"},
		{"deep path", "/static/a/b/c/d.txt", "a/b/c/d.txt"},
		{"root path", "/static/", "index.html"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
			}

			result := decodeJSON(t, res.Body)
			if result["file"] != tt.wantFile {
				t.Errorf("file = %q, want %q", result["file"], tt.wantFile)
			}
		})
	}
}

// TestVersionedAPI tests Exercise 5: versioned API routes.
func TestVersionedAPI(t *testing.T) {
	mux := VersionedAPI()

	t.Run("v1 status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["version"] != "v1" {
			t.Errorf("version = %q, want %q", result["version"], "v1")
		}
		if result["status"] != "ok" {
			t.Errorf("status = %q, want %q", result["status"], "ok")
		}
	})

	t.Run("v2 status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/status", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["version"] != "v2" {
			t.Errorf("version = %q, want %q", result["version"], "v2")
		}
	})

	t.Run("v1 users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["version"] != "v1" {
			t.Errorf("version = %q, want %q", result["version"], "v1")
		}
	})

	t.Run("v2 users with meta", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v2/users", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["version"] != "v2" {
			t.Errorf("version = %q, want %q", result["version"], "v2")
		}
		if result["meta"] == nil {
			t.Error("v2 users response should include 'meta' field")
		}
	})
}

// TestStripPrefixRouter tests Exercise 6: subrouting with StripPrefix.
func TestStripPrefixRouter(t *testing.T) {
	mux := StripPrefixRouter()

	tests := []struct {
		name     string
		path     string
		wantPage string
	}{
		{"admin dashboard", "/admin/dashboard", "dashboard"},
		{"admin settings", "/admin/settings", "settings"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
			}

			result := decodeJSON(t, res.Body)
			if result["page"] != tt.wantPage {
				t.Errorf("page = %q, want %q", result["page"], tt.wantPage)
			}
		})
	}
}

// TestCustomErrorRouter tests Exercise 7: custom 404 responses.
func TestCustomErrorRouter(t *testing.T) {
	mux := CustomErrorRouter()

	t.Run("health endpoint works", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["status"] != "ok" {
			t.Errorf("status = %q, want %q", result["status"], "ok")
		}
	})

	t.Run("data endpoint works", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
	})

	t.Run("unknown path returns custom 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
		}

		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		result := decodeJSON(t, res.Body)
		if result["error"] != "not found" {
			t.Errorf("error = %q, want %q", result["error"], "not found")
		}
		if result["path"] != "/api/unknown" {
			t.Errorf("path = %q, want %q", result["path"], "/api/unknown")
		}
	})
}

// TestBlogRoutes tests Exercise 8: complete blog API route table.
func TestBlogRoutes(t *testing.T) {
	mux := BlogRoutes()

	t.Run("list posts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "list_posts" {
			t.Errorf("route = %q, want %q", result["route"], "list_posts")
		}
	})

	t.Run("create post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/posts", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusCreated)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "create_post" {
			t.Errorf("route = %q, want %q", result["route"], "create_post")
		}
	})

	t.Run("get post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/42", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "get_post" {
			t.Errorf("route = %q, want %q", result["route"], "get_post")
		}
		if result["id"] != "42" {
			t.Errorf("id = %q, want %q", result["id"], "42")
		}
	})

	t.Run("update post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/posts/42", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "update_post" {
			t.Errorf("route = %q, want %q", result["route"], "update_post")
		}
		if result["id"] != "42" {
			t.Errorf("id = %q, want %q", result["id"], "42")
		}
	})

	t.Run("delete post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/posts/42", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNoContent)
		}
	})

	t.Run("list comments for post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/42/comments", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "list_comments" {
			t.Errorf("route = %q, want %q", result["route"], "list_comments")
		}
		if result["post_id"] != "42" {
			t.Errorf("post_id = %q, want %q", result["post_id"], "42")
		}
	})

	t.Run("create comment for post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/posts/42/comments", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusCreated)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "create_comment" {
			t.Errorf("route = %q, want %q", result["route"], "create_comment")
		}
		if result["post_id"] != "42" {
			t.Errorf("post_id = %q, want %q", result["post_id"], "42")
		}
	})

	t.Run("list users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "list_users" {
			t.Errorf("route = %q, want %q", result["route"], "list_users")
		}
	})

	t.Run("get user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/7", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		result := decodeJSON(t, res.Body)
		if result["route"] != "get_user" {
			t.Errorf("route = %q, want %q", result["route"], "get_user")
		}
		if result["id"] != "7" {
			t.Errorf("id = %q, want %q", result["id"], "7")
		}
	})
}
