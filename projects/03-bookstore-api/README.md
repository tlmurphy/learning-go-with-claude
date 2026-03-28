# Project 03: Bookstore API

**Prerequisite modules:** 14-21 (Web Services)

## Overview

Build a complete REST API for a bookstore. This project brings together HTTP
handling, routing, middleware, JSON serialization, data access patterns,
configuration, and graceful shutdown into a well-architected service.

The emphasis is on **clean architecture** — even though storage is in-memory,
the code should be structured so swapping in a real database later is trivial.

## Resources

| Resource | Endpoints |
|----------|-----------|
| Books    | Full CRUD + search |
| Authors  | Full CRUD |
| Reviews  | CRUD, nested under books (`/books/:id/reviews`) |

## Requirements

### API Design
- RESTful routes with correct HTTP methods and status codes
- JSON request/response bodies with input validation
- Consistent error response format across all endpoints:
  ```json
  {
    "error": {
      "code": "NOT_FOUND",
      "message": "Book with ID 42 not found"
    }
  }
  ```
- Pagination for list endpoints (`?page=1&per_page=20`)
- Filtering (`?author=Tolkien&year=1954`) and sorting (`?sort=title&order=asc`)

### Storage
- In-memory storage using the **repository pattern**
- Define repository interfaces so storage can be swapped without changing handlers
- Each resource gets its own repository interface

### Middleware Stack
- **Logging** — log method, path, status code, and duration for every request
- **Recovery** — catch panics and return 500 instead of crashing
- **Request ID** — generate a unique ID for each request, include in response headers and logs
- **CORS** — allow configurable origins

### Operational Requirements
- Health check endpoint (`GET /health`)
- Graceful shutdown on SIGINT/SIGTERM
- Configuration loaded from environment variables (port, allowed origins, log level, etc.)

## Architecture

```
projects/03-bookstore-api/
  cmd/server/main.go          — entry point, wiring
  internal/
    model/
      book.go                 — Book, Author, Review types
    repository/
      repository.go           — interfaces
      memory.go               — in-memory implementations
    service/
      book_service.go         — business logic
    handler/
      book.go                 — HTTP handlers for books
      author.go               — HTTP handlers for authors
      review.go               — HTTP handlers for reviews
      routes.go               — route registration
    middleware/
      logging.go
      recovery.go
      request_id.go
      cors.go
    config/
      config.go               — configuration struct + loader
```

## Hints

<details>
<summary>Build order</summary>

1. Define model types first (`Book`, `Author`, `Review`)
2. Define repository interfaces
3. Implement in-memory repositories
4. Build handlers for one resource (start with Books)
5. Wire up routes in main.go and test with curl
6. Add middleware one at a time
7. Add pagination, filtering, sorting
8. Repeat handlers for Authors and Reviews
9. Add graceful shutdown

</details>

<details>
<summary>Repository interface pattern</summary>

```go
type BookRepository interface {
    FindAll(ctx context.Context, filter BookFilter) ([]Book, int, error)
    FindByID(ctx context.Context, id string) (Book, error)
    Create(ctx context.Context, book Book) (Book, error)
    Update(ctx context.Context, id string, book Book) (Book, error)
    Delete(ctx context.Context, id string) error
}
```

The second return value from FindAll is the total count (for pagination).

</details>

<details>
<summary>Handler pattern</summary>

Group handlers in a struct that holds their dependencies:

```go
type BookHandler struct {
    service BookService
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) { ... }
func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) { ... }
```

</details>

<details>
<summary>Validation hint</summary>

Write a simple validation helper rather than pulling in a framework.
Return a map of field names to error messages:

```go
func (b *CreateBookRequest) Validate() map[string]string { ... }
```

</details>

## Stretch Goals

- **SQLite persistence** — add a second repository implementation backed by SQLite
- **OpenAPI documentation** — write an OpenAPI 3.0 spec for your API
- **Rate limiting** — per-client rate limiting using a token bucket
- **ETag caching** — support `If-None-Match` / `ETag` for conditional GETs
- **WebSocket notifications** — push real-time events when new reviews are posted
