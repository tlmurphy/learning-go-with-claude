// Package restapi covers building RESTful APIs in Go — from request parsing
// and response formatting to validation, error handling, pagination, and
// the patterns that make real APIs maintainable.
package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/*
=============================================================================
 BUILDING A REST API
=============================================================================

REST (Representational State Transfer) is the dominant architectural style
for web APIs. The core idea is simple: you model your domain as resources,
and clients interact with those resources using standard HTTP methods.

  Resource: /books
  GET    /books       -> List all books
  POST   /books       -> Create a new book
  GET    /books/{id}  -> Get a specific book
  PUT    /books/{id}  -> Update a specific book
  DELETE /books/{id}  -> Delete a specific book

The power of REST isn't just the URL structure — it's that each operation
maps to well-understood HTTP semantics:

  Method   | Idempotent? | Safe? | Typical Status Codes
  ---------|-------------|-------|---------------------
  GET      | Yes         | Yes   | 200, 404
  POST     | No          | No    | 201, 400, 409
  PUT      | Yes         | No    | 200, 204, 404
  DELETE   | Yes         | No    | 204, 404
  PATCH    | No          | No    | 200, 400, 404

"Idempotent" means calling it multiple times has the same effect as calling
it once. This matters for retries — if a network error occurs, the client
can safely retry a PUT or DELETE without worrying about duplicates. POST
is not idempotent, which is why creating resources needs more care.

=============================================================================
 THE MODEL
=============================================================================

Every API starts with a domain model. In Go, this is typically a struct
with JSON tags for serialization. The JSON tags control how the struct
fields map to JSON keys.

Important conventions:
- Use lowercase snake_case for JSON keys (not Go's PascalCase)
- Use omitempty for optional fields
- Use pointer types for fields that can be null/absent in updates
- Separate your API model from your database model when they diverge

=============================================================================
*/

// Book is the domain model for our example API. This struct serves double
// duty as both the internal representation and the JSON serialization format.
// In larger applications, you'd separate these concerns.
type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year"`
	ISBN      string    `json:"isbn,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateBookRequest represents the expected JSON body for creating a book.
// Separate from Book because the client shouldn't set ID, CreatedAt, etc.
type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	ISBN   string `json:"isbn,omitempty"`
}

// UpdateBookRequest represents the expected JSON body for updating a book.
// All fields are pointers so we can distinguish "not provided" from "set to zero value".
type UpdateBookRequest struct {
	Title  *string `json:"title,omitempty"`
	Author *string `json:"author,omitempty"`
	Year   *int    `json:"year,omitempty"`
	ISBN   *string `json:"isbn,omitempty"`
}

/*
=============================================================================
 ERROR HANDLING
=============================================================================

Consistent error responses are crucial for API usability. Every error from
your API should have the same structure so clients can parse errors
predictably.

A common format:

  {
    "error": {
      "code": "not_found",
      "message": "Book with ID '42' not found"
    }
  }

Some APIs include additional fields like "details" (for validation errors)
or "request_id" (for debugging). The key is consistency — pick a format
and use it everywhere.

Status code guidelines:
  400 Bad Request      -> Client sent invalid data (malformed JSON, validation failure)
  404 Not Found        -> Resource doesn't exist
  405 Method Not Allowed -> Wrong HTTP method for this endpoint
  409 Conflict         -> Resource already exists (duplicate key, etc.)
  422 Unprocessable    -> Semantically invalid (valid JSON, but invalid values)
  500 Internal Error   -> Something broke on our side (always log these!)

=============================================================================
*/

// ErrorResponse is the standard error format for our API.
// Using a consistent structure lets clients handle all errors uniformly.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error code and human-readable message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON is a helper that writes a JSON response with the given status code.
// Centralizing this avoids repeating Content-Type and encoding logic everywhere.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError is a helper that writes a consistent error response.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

/*
=============================================================================
 IN-MEMORY STORAGE
=============================================================================

For learning purposes, we use an in-memory map as our "database." This
keeps things simple while we focus on the HTTP handling patterns. In a
real application, you'd replace this with a database connection.

The storage must be safe for concurrent access because HTTP handlers run
in separate goroutines. We use a sync.RWMutex for this:
- RLock for reads (multiple readers allowed simultaneously)
- Lock for writes (exclusive access)

This is the same pattern you'd use with a connection pool — the pool
handles concurrency internally, and your handler just calls pool.Get().

=============================================================================
*/

// BookStore is a thread-safe in-memory store for books.
type BookStore struct {
	mu     sync.RWMutex
	books  map[string]Book
	nextID int
}

// NewBookStore creates an initialized BookStore.
func NewBookStore() *BookStore {
	return &BookStore{
		books:  make(map[string]Book),
		nextID: 1,
	}
}

// All returns all books in the store.
func (s *BookStore) All() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		books = append(books, b)
	}
	return books
}

// Get returns a book by ID and whether it was found.
func (s *BookStore) Get(id string) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.books[id]
	return b, ok
}

// Create adds a new book and returns it with the assigned ID.
func (s *BookStore) Create(req CreateBookRequest) Book {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	book := Book{
		ID:        strconv.Itoa(s.nextID),
		Title:     req.Title,
		Author:    req.Author,
		Year:      req.Year,
		ISBN:      req.ISBN,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.nextID++
	s.books[book.ID] = book
	return book
}

// Update modifies an existing book. Returns the updated book and whether
// the book was found.
func (s *BookStore) Update(id string, req UpdateBookRequest) (Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return Book{}, false
	}

	// Only update fields that were provided (non-nil pointers)
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Year != nil {
		book.Year = *req.Year
	}
	if req.ISBN != nil {
		book.ISBN = *req.ISBN
	}
	book.UpdatedAt = time.Now().UTC()

	s.books[id] = book
	return book, true
}

// Delete removes a book by ID. Returns whether it was found.
func (s *BookStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.books[id]
	if ok {
		delete(s.books, id)
	}
	return ok
}

// Count returns the number of books in the store.
func (s *BookStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.books)
}

/*
=============================================================================
 REQUEST PARSING
=============================================================================

Parsing request bodies safely requires several steps:

1. Limit the body size (prevent memory exhaustion)
2. Decode the JSON
3. Handle decode errors gracefully (report what went wrong)
4. Validate the data (required fields, value ranges, etc.)

Each step can fail, and each failure should produce a clear, actionable
error message. "Bad Request" tells the client nothing useful. "Title is
required and must be between 1 and 200 characters" tells them exactly
what to fix.

=============================================================================
*/

// decodeJSON is a helper that decodes a JSON request body with proper
// error handling. It limits the body size and returns a descriptive error.
func decodeJSON(r *http.Request, v any) error {
	// Limit request body to 1 MB
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// validateCreateBook checks that a CreateBookRequest has all required fields.
func validateCreateBook(req CreateBookRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Author == "" {
		return fmt.Errorf("author is required")
	}
	if req.Year < 0 || req.Year > time.Now().Year()+1 {
		return fmt.Errorf("year must be between 0 and %d", time.Now().Year()+1)
	}
	return nil
}

/*
=============================================================================
 THE HANDLERS
=============================================================================

Each CRUD operation gets its own handler function. This is cleaner than
a single handler with a switch on the method. With Go 1.22+ routing,
the mux dispatches to the right handler based on method and path.

Notice the consistent pattern in each handler:
1. Parse input (path params, query params, body)
2. Validate input
3. Execute business logic (call the store)
4. Format and return the response

This separation makes handlers easy to test and reason about.

=============================================================================
*/

// BookAPI holds the handlers for the books API. Having a struct lets us
// inject dependencies (the store) without global variables.
type BookAPI struct {
	Store *BookStore
}

// NewBookAPI creates a BookAPI with a fresh store.
func NewBookAPI() *BookAPI {
	return &BookAPI{Store: NewBookStore()}
}

// List handles GET /books — returns all books.
func (api *BookAPI) List(w http.ResponseWriter, r *http.Request) {
	books := api.Store.All()

	// Never return null for a list — return an empty array instead.
	// This makes client code simpler: they can always iterate.
	if books == nil {
		books = []Book{}
	}

	writeJSON(w, http.StatusOK, books)
}

// GetByID handles GET /books/{id} — returns a specific book.
func (api *BookAPI) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	book, ok := api.Store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found",
			fmt.Sprintf("book with ID %q not found", id))
		return
	}

	writeJSON(w, http.StatusOK, book)
}

// Create handles POST /books — creates a new book.
func (api *BookAPI) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if err := validateCreateBook(req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	book := api.Store.Create(req)

	// Set Location header to the URL of the new resource.
	// This is a REST convention that helps clients find what they just created.
	w.Header().Set("Location", "/books/"+book.ID)
	writeJSON(w, http.StatusCreated, book)
}

// Update handles PUT /books/{id} — updates an existing book.
func (api *BookAPI) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req UpdateBookRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	book, ok := api.Store.Update(id, req)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found",
			fmt.Sprintf("book with ID %q not found", id))
		return
	}

	writeJSON(w, http.StatusOK, book)
}

// Delete handles DELETE /books/{id} — removes a book.
func (api *BookAPI) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if !api.Store.Delete(id) {
		writeError(w, http.StatusNotFound, "not_found",
			fmt.Sprintf("book with ID %q not found", id))
		return
	}

	// 204 No Content — the standard response for successful deletes.
	// No body is sent, which is why we don't call writeJSON.
	w.WriteHeader(http.StatusNoContent)
}

/*
=============================================================================
 PAGINATION
=============================================================================

Any list endpoint that could return many items needs pagination. Without
it, a list of 100,000 books would be returned in a single response,
crushing both your server and the client.

There are two common approaches:

Offset/Limit (simpler, what we use here):
  GET /books?offset=20&limit=10
  Returns items 20-29. Simple, but slow for large offsets because the
  database has to skip past all the preceding records.

Cursor-based (better for large datasets):
  GET /books?cursor=abc123&limit=10
  Returns the next 10 items after the cursor. The cursor is usually an
  opaque token encoding the last item's sort key. More efficient but
  more complex to implement.

For most APIs, offset/limit is fine until you have performance issues.

=============================================================================
*/

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse struct {
	Data   []Book `json:"data"`
	Total  int    `json:"total"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// ListPaginated handles GET /books with pagination support.
func (api *BookAPI) ListPaginated(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters with defaults
	offset := 0
	limit := 10

	if v := r.URL.Query().Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	allBooks := api.Store.All()
	total := len(allBooks)

	// Apply pagination
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	page := allBooks[start:end]
	if page == nil {
		page = []Book{}
	}

	writeJSON(w, http.StatusOK, PaginatedResponse{
		Data:   page,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	})
}

/*
=============================================================================
 WIRING IT ALL TOGETHER
=============================================================================

The final step is registering all handlers on a mux. This is where the
route table comes together. Notice how the API struct's methods map
cleanly to HTTP methods and paths.

=============================================================================
*/

// RegisterBookRoutes wires up all book API routes on the given mux.
func (api *BookAPI) RegisterBookRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /books", api.List)
	mux.HandleFunc("POST /books", api.Create)
	mux.HandleFunc("GET /books/{id}", api.GetByID)
	mux.HandleFunc("PUT /books/{id}", api.Update)
	mux.HandleFunc("DELETE /books/{id}", api.Delete)
}

// NewBookServer creates a complete HTTP server with all book routes registered.
func NewBookServer() *http.ServeMux {
	mux := http.NewServeMux()
	api := NewBookAPI()
	api.RegisterBookRoutes(mux)
	return mux
}
