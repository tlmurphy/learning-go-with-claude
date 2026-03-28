package restapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// jsonBody is a test helper that creates a reader from a JSON string.
func jsonBody(s string) *strings.Reader {
	return strings.NewReader(s)
}

// decodeJSONResponse is a test helper that decodes a JSON response.
func decodeJSONResponse(t *testing.T, body io.Reader, v interface{}) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(v); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
}

// TestCreateHandler tests Exercise 1: creating a todo.
func TestCreateHandler(t *testing.T) {
	t.Run("creates a todo successfully", func(t *testing.T) {
		api := NewTodoAPI()
		req := httptest.NewRequest(http.MethodPost, "/todos", jsonBody(`{"title": "Buy groceries"}`))
		rec := httptest.NewRecorder()

		api.CreateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusCreated)
		}

		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		location := res.Header.Get("Location")
		if location == "" {
			t.Error("Location header should be set")
		}

		var todo Todo
		decodeJSONResponse(t, res.Body, &todo)

		if todo.ID == "" {
			t.Error("todo ID should not be empty")
		}
		if todo.Title != "Buy groceries" {
			t.Errorf("title = %q, want %q", todo.Title, "Buy groceries")
		}
		if todo.Completed {
			t.Error("new todo should not be completed")
		}
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		api := NewTodoAPI()
		req := httptest.NewRequest(http.MethodPost, "/todos", jsonBody(`{invalid`))
		rec := httptest.NewRecorder()

		api.CreateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}

		var errResp APIError
		decodeJSONResponse(t, res.Body, &errResp)
		if errResp.Error.Code != "invalid_request" {
			t.Errorf("error code = %q, want %q", errResp.Error.Code, "invalid_request")
		}
	})

	t.Run("rejects empty title", func(t *testing.T) {
		api := NewTodoAPI()
		req := httptest.NewRequest(http.MethodPost, "/todos", jsonBody(`{"title": ""}`))
		rec := httptest.NewRecorder()

		api.CreateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}

		var errResp APIError
		decodeJSONResponse(t, res.Body, &errResp)
		if errResp.Error.Code != "validation_error" {
			t.Errorf("error code = %q, want %q", errResp.Error.Code, "validation_error")
		}
		if errResp.Error.Message != "title is required" {
			t.Errorf("error message = %q, want %q", errResp.Error.Message, "title is required")
		}
	})
}

// TestListHandler tests Exercise 2: listing todos.
func TestListHandler(t *testing.T) {
	t.Run("returns empty array when no todos", func(t *testing.T) {
		api := NewTodoAPI()
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()

		api.ListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		body, _ := io.ReadAll(res.Body)
		// Should be an empty array, not null
		trimmed := strings.TrimSpace(string(body))
		if trimmed != "[]" {
			t.Errorf("body = %q, want %q", trimmed, "[]")
		}
	})

	t.Run("returns todos after creation", func(t *testing.T) {
		api := NewTodoAPI()
		api.Store.Create("Todo 1")
		api.Store.Create("Todo 2")

		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()

		api.ListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var todos []Todo
		decodeJSONResponse(t, res.Body, &todos)

		if len(todos) != 2 {
			t.Errorf("got %d todos, want 2", len(todos))
		}
	})
}

// TestGetByIDHandler tests Exercise 3: getting a specific todo.
func TestGetByIDHandler(t *testing.T) {
	t.Run("returns existing todo", func(t *testing.T) {
		api := NewTodoAPI()
		created := api.Store.Create("Test todo")

		// Use the mux approach to populate PathValue
		mux := http.NewServeMux()
		mux.HandleFunc("GET /todos/{id}", api.GetByIDHandler)

		req := httptest.NewRequest(http.MethodGet, "/todos/"+created.ID, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var todo Todo
		decodeJSONResponse(t, res.Body, &todo)

		if todo.ID != created.ID {
			t.Errorf("id = %q, want %q", todo.ID, created.ID)
		}
		if todo.Title != "Test todo" {
			t.Errorf("title = %q, want %q", todo.Title, "Test todo")
		}
	})

	t.Run("returns 404 for missing todo", func(t *testing.T) {
		api := NewTodoAPI()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /todos/{id}", api.GetByIDHandler)

		req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
		}

		var errResp APIError
		decodeJSONResponse(t, res.Body, &errResp)
		if errResp.Error.Code != "not_found" {
			t.Errorf("error code = %q, want %q", errResp.Error.Code, "not_found")
		}
	})
}

// TestUpdateHandler tests Exercise 4: updating a todo.
func TestUpdateHandler(t *testing.T) {
	t.Run("updates title", func(t *testing.T) {
		api := NewTodoAPI()
		created := api.Store.Create("Original title")

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /todos/{id}", api.UpdateHandler)

		req := httptest.NewRequest(http.MethodPut, "/todos/"+created.ID,
			jsonBody(`{"title": "Updated title"}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var todo Todo
		decodeJSONResponse(t, res.Body, &todo)
		if todo.Title != "Updated title" {
			t.Errorf("title = %q, want %q", todo.Title, "Updated title")
		}
	})

	t.Run("updates completed status", func(t *testing.T) {
		api := NewTodoAPI()
		created := api.Store.Create("Test todo")

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /todos/{id}", api.UpdateHandler)

		req := httptest.NewRequest(http.MethodPut, "/todos/"+created.ID,
			jsonBody(`{"completed": true}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var todo Todo
		decodeJSONResponse(t, res.Body, &todo)
		if !todo.Completed {
			t.Error("todo should be completed after update")
		}
		// Title should remain unchanged
		if todo.Title != "Test todo" {
			t.Errorf("title = %q, want %q (should be unchanged)", todo.Title, "Test todo")
		}
	})

	t.Run("returns 404 for missing todo", func(t *testing.T) {
		api := NewTodoAPI()

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /todos/{id}", api.UpdateHandler)

		req := httptest.NewRequest(http.MethodPut, "/todos/999",
			jsonBody(`{"title": "Updated"}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
		}
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		api := NewTodoAPI()
		api.Store.Create("Test todo")

		mux := http.NewServeMux()
		mux.HandleFunc("PUT /todos/{id}", api.UpdateHandler)

		req := httptest.NewRequest(http.MethodPut, "/todos/1", jsonBody(`{bad json`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusBadRequest)
		}
	})
}

// TestDeleteHandler tests Exercise 5: deleting a todo.
func TestDeleteHandler(t *testing.T) {
	t.Run("deletes existing todo", func(t *testing.T) {
		api := NewTodoAPI()
		created := api.Store.Create("To delete")

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /todos/{id}", api.DeleteHandler)

		req := httptest.NewRequest(http.MethodDelete, "/todos/"+created.ID, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNoContent)
		}

		// Verify it was actually deleted
		if api.Store.Count() != 0 {
			t.Errorf("store count = %d, want 0 after delete", api.Store.Count())
		}
	})

	t.Run("returns 404 for missing todo", func(t *testing.T) {
		api := NewTodoAPI()

		mux := http.NewServeMux()
		mux.HandleFunc("DELETE /todos/{id}", api.DeleteHandler)

		req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
		}
	})
}

// TestPaginatedListHandler tests Exercise 6: pagination.
func TestPaginatedListHandler(t *testing.T) {
	// Helper to populate the store with n todos
	setup := func(n int) *TodoAPI {
		api := NewTodoAPI()
		for i := 0; i < n; i++ {
			api.Store.Create("Todo " + intToString(i+1))
		}
		return api
	}

	t.Run("returns paginated results with defaults", func(t *testing.T) {
		api := setup(25)
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()

		api.PaginatedListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var result struct {
			Data   []Todo `json:"data"`
			Total  int    `json:"total"`
			Offset int    `json:"offset"`
			Limit  int    `json:"limit"`
		}
		decodeJSONResponse(t, res.Body, &result)

		if result.Total != 25 {
			t.Errorf("total = %d, want 25", result.Total)
		}
		if result.Offset != 0 {
			t.Errorf("offset = %d, want 0", result.Offset)
		}
		if result.Limit != 10 {
			t.Errorf("limit = %d, want 10", result.Limit)
		}
		if len(result.Data) != 10 {
			t.Errorf("data length = %d, want 10", len(result.Data))
		}
	})

	t.Run("respects offset and limit", func(t *testing.T) {
		api := setup(25)
		req := httptest.NewRequest(http.MethodGet, "/todos?offset=20&limit=10", nil)
		rec := httptest.NewRecorder()

		api.PaginatedListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var result struct {
			Data   []Todo `json:"data"`
			Total  int    `json:"total"`
			Offset int    `json:"offset"`
			Limit  int    `json:"limit"`
		}
		decodeJSONResponse(t, res.Body, &result)

		if result.Total != 25 {
			t.Errorf("total = %d, want 25", result.Total)
		}
		if result.Offset != 20 {
			t.Errorf("offset = %d, want 20", result.Offset)
		}
		if len(result.Data) != 5 {
			t.Errorf("data length = %d, want 5 (only 5 remaining)", len(result.Data))
		}
	})

	t.Run("returns empty data for out-of-range offset", func(t *testing.T) {
		api := setup(5)
		req := httptest.NewRequest(http.MethodGet, "/todos?offset=100", nil)
		rec := httptest.NewRecorder()

		api.PaginatedListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var result struct {
			Data  []Todo `json:"data"`
			Total int    `json:"total"`
		}
		decodeJSONResponse(t, res.Body, &result)

		if len(result.Data) != 0 {
			t.Errorf("data length = %d, want 0", len(result.Data))
		}
		if result.Total != 5 {
			t.Errorf("total = %d, want 5", result.Total)
		}
	})

	t.Run("falls back to defaults for invalid params", func(t *testing.T) {
		api := setup(5)
		req := httptest.NewRequest(http.MethodGet, "/todos?offset=abc&limit=-1", nil)
		rec := httptest.NewRecorder()

		api.PaginatedListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var result struct {
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
		}
		decodeJSONResponse(t, res.Body, &result)

		if result.Offset != 0 {
			t.Errorf("offset = %d, want 0 (default)", result.Offset)
		}
		if result.Limit != 10 {
			t.Errorf("limit = %d, want 10 (default)", result.Limit)
		}
	})

	t.Run("returns empty array not null when no data", func(t *testing.T) {
		api := NewTodoAPI()
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()

		api.PaginatedListHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		if strings.Contains(string(body), "null") {
			t.Error("data field should be empty array [], not null")
		}
	})
}

// TestRespondWithError tests Exercise 7: consistent error responses.
func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		code    string
		message string
	}{
		{
			name:    "bad request",
			status:  http.StatusBadRequest,
			code:    "invalid_request",
			message: "the request body is malformed",
		},
		{
			name:    "not found",
			status:  http.StatusNotFound,
			code:    "not_found",
			message: "resource not found",
		},
		{
			name:    "internal error",
			status:  http.StatusInternalServerError,
			code:    "internal_error",
			message: "an unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			RespondWithError(rec, tt.status, tt.code, tt.message)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.status {
				t.Errorf("status = %d, want %d", res.StatusCode, tt.status)
			}

			ct := res.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			var errResp APIError
			decodeJSONResponse(t, res.Body, &errResp)

			if errResp.Error.Code != tt.code {
				t.Errorf("error code = %q, want %q", errResp.Error.Code, tt.code)
			}
			if errResp.Error.Message != tt.message {
				t.Errorf("error message = %q, want %q", errResp.Error.Message, tt.message)
			}
		})
	}
}

// TestCRUDHandler tests Exercise 8: complete CRUD integration.
func TestCRUDHandler(t *testing.T) {
	api := NewTodoAPI()
	mux := api.CRUDHandler()

	// Create a todo
	t.Run("POST /todos creates a todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/todos", jsonBody(`{"title": "Integration test"}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusCreated)
		}

		var todo Todo
		decodeJSONResponse(t, res.Body, &todo)
		if todo.Title != "Integration test" {
			t.Errorf("title = %q, want %q", todo.Title, "Integration test")
		}
	})

	// List todos
	t.Run("GET /todos returns todos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var todos []Todo
		decodeJSONResponse(t, res.Body, &todos)
		if len(todos) != 1 {
			t.Errorf("got %d todos, want 1", len(todos))
		}
	})

	// Get specific todo
	t.Run("GET /todos/1 returns the todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
	})

	// Update todo
	t.Run("PUT /todos/1 updates the todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/todos/1",
			jsonBody(`{"title": "Updated integration test"}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}

		var todo Todo
		decodeJSONResponse(t, res.Body, &todo)
		if todo.Title != "Updated integration test" {
			t.Errorf("title = %q, want %q", todo.Title, "Updated integration test")
		}
	})

	// Delete todo
	t.Run("DELETE /todos/1 removes the todo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNoContent)
		}
	})

	// Verify deletion
	t.Run("GET /todos/1 returns 404 after deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
		}
	})
}
