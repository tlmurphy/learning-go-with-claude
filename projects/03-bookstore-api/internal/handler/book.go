package handler

import "net/http"

// BookHandler handles HTTP requests for the /books resource.
//
// TODO: Add a service dependency (BookService interface or struct)
// and implement each handler method.

type BookHandler struct {
	// TODO: Add service dependency
}

// List handles GET /books
// Should support query params: page, per_page, author_id, title, sort, order
func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse query params into a BookFilter, call service, return JSON
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Get handles GET /books/{id}
func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract ID from path, call service, return JSON or 404
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Create handles POST /books
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode JSON body, validate, call service, return 201
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Update handles PUT /books/{id}
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract ID, decode body, validate, call service, return JSON
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Delete handles DELETE /books/{id}
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract ID, call service, return 204
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
