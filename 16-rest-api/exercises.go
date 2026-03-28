package restapi

import (
	"net/http"
	"sync"
	"time"
)

/*
=============================================================================
 EXERCISES: Building a REST API
=============================================================================

These exercises walk you through building a complete REST API for a Todo
application. You'll implement each CRUD operation, add pagination, create
consistent error responses, and build a full handler struct.

The Todo model and TodoStore are provided below. Your job is to implement
the HTTP handlers.

=============================================================================
*/

// Todo is the domain model for exercises.
type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoRequest is the expected JSON body for creating a todo.
type CreateTodoRequest struct {
	Title string `json:"title"`
}

// UpdateTodoRequest is the expected JSON body for updating a todo.
type UpdateTodoRequest struct {
	Title     *string `json:"title,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

// TodoStore is a thread-safe in-memory store for todos.
type TodoStore struct {
	mu     sync.RWMutex
	todos  map[string]Todo
	nextID int
}

// NewTodoStore creates an initialized TodoStore.
func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make(map[string]Todo),
		nextID: 1,
	}
}

// All returns all todos, sorted is not required.
func (s *TodoStore) All() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Todo, 0, len(s.todos))
	for _, t := range s.todos {
		result = append(result, t)
	}
	return result
}

// Get returns a todo by ID.
func (s *TodoStore) Get(id string) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	return t, ok
}

// Create adds a new todo and returns it.
func (s *TodoStore) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	todo := Todo{
		ID:        intToString(s.nextID),
		Title:     title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.nextID++
	s.todos[todo.ID] = todo
	return todo
}

// Update modifies an existing todo. Returns the updated todo and whether found.
func (s *TodoStore) Update(id string, req UpdateTodoRequest) (Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo, ok := s.todos[id]
	if !ok {
		return Todo{}, false
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = time.Now().UTC()
	s.todos[id] = todo
	return todo, true
}

// Delete removes a todo by ID. Returns whether it existed.
func (s *TodoStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.todos[id]
	if ok {
		delete(s.todos, id)
	}
	return ok
}

// Count returns the number of todos.
func (s *TodoStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.todos)
}

// intToString converts an int to string without importing strconv in exercises.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// TodoAPI holds the handlers for the todo API.
// Use this struct for all exercises.
type TodoAPI struct {
	Store *TodoStore
}

// NewTodoAPI creates a TodoAPI with a fresh store.
func NewTodoAPI() *TodoAPI {
	return &TodoAPI{Store: NewTodoStore()}
}

// Exercise 1: CreateHandler
//
// Implement a handler for POST /todos that creates a new todo.
//
// Requirements:
// - Decode the JSON body into a CreateTodoRequest
//   Use json.NewDecoder(r.Body).Decode(&req)
// - If JSON decoding fails, return 400 with:
//   {"error": {"code": "invalid_request", "message": "<error message>"}}
// - If Title is empty, return 400 with:
//   {"error": {"code": "validation_error", "message": "title is required"}}
// - On success, create the todo using api.Store.Create(req.Title)
// - Set the Location header to "/todos/" + todo.ID
// - Return 201 with the created todo as JSON
// - Set Content-Type to "application/json" for all responses
//
// Tip: Use the writeJSON and writeError helpers from lesson.go, or write
// your own. The test expects the exact error format shown above.
func (api *TodoAPI) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 2: ListHandler
//
// Implement a handler for GET /todos that returns all todos.
//
// Requirements:
// - Return all todos from the store as a JSON array
// - If there are no todos, return an empty array [] (not null)
// - Set Content-Type to "application/json"
// - Return status 200
//
// Remember: json.Marshal(nil slice) produces "null", but
// json.Marshal([]Todo{}) produces "[]". Always initialize your slice.
func (api *TodoAPI) ListHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 3: GetByIDHandler
//
// Implement a handler for GET /todos/{id} that returns a specific todo.
//
// Requirements:
// - Extract the "id" path parameter using r.PathValue("id")
// - Look up the todo in the store
// - If not found, return 404 with:
//   {"error": {"code": "not_found", "message": "todo not found"}}
// - If found, return 200 with the todo as JSON
// - Set Content-Type to "application/json" for all responses
func (api *TodoAPI) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 4: UpdateHandler
//
// Implement a handler for PUT /todos/{id} that updates an existing todo.
//
// Requirements:
// - Extract the "id" path parameter
// - Decode the JSON body into an UpdateTodoRequest
// - If JSON decoding fails, return 400 with:
//   {"error": {"code": "invalid_request", "message": "<error message>"}}
// - Call api.Store.Update(id, req)
// - If the todo is not found, return 404 with:
//   {"error": {"code": "not_found", "message": "todo not found"}}
// - On success, return 200 with the updated todo as JSON
// - Set Content-Type to "application/json" for all responses
func (api *TodoAPI) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 5: DeleteHandler
//
// Implement a handler for DELETE /todos/{id}.
//
// Requirements:
// - Extract the "id" path parameter
// - Call api.Store.Delete(id)
// - If not found, return 404 with:
//   {"error": {"code": "not_found", "message": "todo not found"}}
// - On success, return 204 No Content (no response body)
// - Set Content-Type to "application/json" for error responses
func (api *TodoAPI) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// Exercise 6: PaginatedListHandler
//
// Implement a handler for GET /todos with pagination support.
//
// Requirements:
// - Parse "offset" query parameter (default 0, must be >= 0)
// - Parse "limit" query parameter (default 10, must be 1-100)
// - Invalid values should fall back to defaults (not return errors)
// - Return a JSON object with:
//   {
//     "data": [...todos...],
//     "total": <total count>,
//     "offset": <offset used>,
//     "limit": <limit used>
//   }
// - The "data" field should contain the paginated slice (never null, use [])
// - Set Content-Type to "application/json"
// - Return status 200
//
// Hint: Use strconv.Atoi to parse query params. If parsing fails, keep default.
func (api *TodoAPI) PaginatedListHandler(w http.ResponseWriter, r *http.Request) {
	// YOUR CODE HERE
}

// APIError is the error response format for Exercise 7.
type APIError struct {
	Error APIErrorDetail `json:"error"`
}

// APIErrorDetail contains the error details.
type APIErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Exercise 7: RespondWithError
//
// Implement a helper function that writes a consistent JSON error response.
//
// Requirements:
// - Set Content-Type to "application/json"
// - Set the HTTP status code to the given status
// - Write a JSON body with the format: {"error": {"code": "...", "message": "..."}}
// - Use the APIError and APIErrorDetail structs defined above
//
// This is a utility function, not a handler. It's used by other handlers
// to ensure all errors have the same format.
func RespondWithError(w http.ResponseWriter, status int, code, message string) {
	// YOUR CODE HERE
}

// Exercise 8: CRUDHandler
//
// Build a complete CRUD handler by registering all todo routes on a mux.
//
// Requirements:
// - Register these routes:
//   GET    /todos       -> api.ListHandler
//   POST   /todos       -> api.CreateHandler
//   GET    /todos/{id}  -> api.GetByIDHandler
//   PUT    /todos/{id}  -> api.UpdateHandler
//   DELETE /todos/{id}  -> api.DeleteHandler
// - Return the configured mux
//
// This ties everything together into a working API.
func (api *TodoAPI) CRUDHandler() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}
