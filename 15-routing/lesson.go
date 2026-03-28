// Package routing covers Go 1.22+'s enhanced ServeMux routing capabilities,
// including method-based routing, path parameters, wildcards, and patterns
// that previously required third-party routers.
package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
=============================================================================
 ROUTING AND URL PATTERNS
=============================================================================

Before Go 1.22, the standard library's ServeMux was embarrassingly basic.
It could only match exact paths and path prefixes — no method-based
routing, no path parameters, no wildcards. Every serious Go web application
needed a third-party router like gorilla/mux, chi, or echo.

Go 1.22 changed everything. The enhanced ServeMux now supports:
  - Method-based routing:  "GET /users"
  - Path parameters:       "/users/{id}"
  - Wildcard catch-alls:   "/files/{path...}"
  - Precedence rules:      most specific pattern wins

This means the standard library is now sufficient for the vast majority of
web applications. You no longer need a third-party router unless you need
very specific features like regex patterns or middleware chaining built
into the router itself.

=============================================================================
 METHOD-BASED ROUTING
=============================================================================

The simplest enhancement: prefix your pattern with an HTTP method.

  mux.HandleFunc("GET /users", listUsers)     // only matches GET
  mux.HandleFunc("POST /users", createUser)   // only matches POST

Without a method prefix, the pattern matches ALL methods:

  mux.HandleFunc("/users", anyMethodHandler)  // matches GET, POST, PUT, etc.

This is a significant improvement. Before 1.22, you had to do method
checking inside every handler:

  func usersHandler(w http.ResponseWriter, r *http.Request) {
      switch r.Method {
      case "GET":
          listUsers(w, r)
      case "POST":
          createUser(w, r)
      default:
          http.Error(w, "method not allowed", 405)
      }
  }

Now the mux does this for you, and automatically returns 405 Method Not
Allowed for unregistered methods on a matched path.

=============================================================================
*/

// DemoMethodRouting creates a mux demonstrating method-based routing.
// The same path ("/items") handles different HTTP methods.
func DemoMethodRouting() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"action": "list items"})
	})

	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"action": "create item"})
	})

	mux.HandleFunc("DELETE /items", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

/*
=============================================================================
 PATH PARAMETERS
=============================================================================

Path parameters let you capture dynamic segments of the URL path:

  mux.HandleFunc("GET /users/{id}", getUser)
  mux.HandleFunc("GET /posts/{slug}", getPost)

Inside the handler, extract the parameter using r.PathValue():

  func getUser(w http.ResponseWriter, r *http.Request) {
      id := r.PathValue("id")
      // id is a string — convert to int if needed
  }

Key behaviors:
- Path parameters match any non-empty segment (everything between slashes)
- The parameter name must be unique within a pattern
- r.PathValue() returns "" if the parameter doesn't exist

Before Go 1.22, extracting path parameters required either:
  1. Third-party router (gorilla/mux, chi)
  2. Manual string splitting: parts := strings.Split(r.URL.Path, "/")
  3. Regex matching on the path

The new approach is cleaner and type-safe (at the routing level).

=============================================================================
*/

// DemoPathParams creates a mux showing path parameter extraction.
func DemoPathParams() *http.ServeMux {
	mux := http.NewServeMux()

	// Single path parameter
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"resource": "user",
			"id":       id,
		})
	})

	// Multiple path parameters
	mux.HandleFunc("GET /users/{userID}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		postID := r.PathValue("postID")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"user_id": userID,
			"post_id": postID,
		})
	})

	return mux
}

/*
=============================================================================
 WILDCARD (CATCH-ALL) ROUTES
=============================================================================

The {name...} syntax creates a wildcard that matches the rest of the path,
including slashes:

  mux.HandleFunc("GET /files/{path...}", serveFile)

For a request to /files/images/logo.png:
  r.PathValue("path") returns "images/logo.png"

This is useful for:
- Static file servers
- Proxy routes (forward everything under /api/v1/ to another service)
- SPA catch-all routes (serve index.html for any unmatched path)

The wildcard must be the last segment in the pattern. You can't have
anything after {path...}.

=============================================================================
*/

// DemoWildcardRoutes creates a mux with catch-all wildcard routes.
func DemoWildcardRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Catch-all route — matches /files/anything/here/including/slashes
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.PathValue("path")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"requested_path": filePath,
		})
	})

	// Specific route takes precedence over wildcard
	mux.HandleFunc("GET /files/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"note": "specific route wins over wildcard",
		})
	})

	return mux
}

/*
=============================================================================
 PATTERN PRECEDENCE
=============================================================================

When multiple patterns could match a request, Go uses these rules:

  1. Most specific pattern wins
  2. Patterns with methods are more specific than those without
  3. Longer paths are more specific than shorter paths
  4. Fixed segments are more specific than wildcards
  5. A {param} is more specific than {param...}

Examples:
  "GET /users/{id}"      wins over  "/users/{id}"       (has method)
  "/users/{id}"          wins over  "/users/{path...}"   (exact vs wildcard)
  "/users/me"            wins over  "/users/{id}"        (literal vs param)
  "/api/v1/users/{id}"   wins over  "/api/{path...}"     (more specific)

This is intuitive — think of it as "the pattern that describes the request
most precisely wins." If two patterns are equally specific (which shouldn't
happen in a well-designed API), registration panics at startup rather than
causing subtle runtime bugs.

=============================================================================
*/

// DemoPrecedence creates a mux showing how pattern specificity works.
func DemoPrecedence() *http.ServeMux {
	mux := http.NewServeMux()

	// Literal path — most specific
	mux.HandleFunc("GET /users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"matched": "GET /users/me (literal)"})
	})

	// Path parameter — less specific than literal
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"matched": "GET /users/{id} (param)",
			"id":      id,
		})
	})

	// Any method — less specific than method-specific
	mux.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"matched": "/users/{id} (any method)",
			"id":      id,
			"method":  r.Method,
		})
	})

	return mux
}

/*
=============================================================================
 TRAILING SLASH BEHAVIOR
=============================================================================

Trailing slashes in patterns have specific meaning in ServeMux:

  "/users"   — matches EXACTLY /users (no trailing slash)
  "/users/"  — matches /users/ AND anything under /users/...

This subtlety trips up many developers. A pattern "/api/" acts as a
prefix match — it matches /api/, /api/foo, /api/foo/bar, etc.

With Go 1.22+ enhanced routing, if you register "/users/" but a client
requests "/users" (no trailing slash), the mux redirects to "/users/"
with a 301 Moved Permanently. This is usually what you want for
browsable URLs but can be surprising for APIs.

For APIs, be explicit about your patterns:
  "GET /users"      — list users (no trailing slash)
  "GET /users/{id}" — get specific user

=============================================================================
*/

// DemoTrailingSlash creates a mux demonstrating trailing slash behavior.
func DemoTrailingSlash() *http.ServeMux {
	mux := http.NewServeMux()

	// Exact match — only "/about" matches
	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "About page (exact match)")
	})

	// Prefix match — "/docs/" and anything under it matches
	mux.HandleFunc("GET /docs/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Docs page, path: %s\n", r.URL.Path)
	})

	return mux
}

/*
=============================================================================
 SUBROUTING WITH StripPrefix
=============================================================================

http.StripPrefix is a middleware that removes a prefix from the URL path
before passing the request to the inner handler. This enables modular
route groups:

  apiMux := http.NewServeMux()
  apiMux.HandleFunc("GET /users", listUsers)
  apiMux.HandleFunc("GET /posts", listPosts)

  mainMux := http.NewServeMux()
  mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

Now requests to /api/v1/users are seen as /users by apiMux. This is
powerful for:
- Versioning APIs
- Mounting sub-applications
- Organizing routes into logical groups

The key detail: StripPrefix modifies r.URL.Path before forwarding. This
means the inner handler sees a "clean" path without the prefix.

=============================================================================
*/

// DemoSubrouting creates a mux with subrouted API versions.
func DemoSubrouting() *http.ServeMux {
	// V1 routes
	v1 := http.NewServeMux()
	v1.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"version": "v1",
			"action":  "list users",
		})
	})

	// V2 routes (maybe with different response format)
	v2 := http.NewServeMux()
	v2.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"version": "v2",
			"action":  "list users",
			"meta":    map[string]int{"total": 0, "page": 1},
		})
	})

	// Mount on the main mux
	main := http.NewServeMux()
	main.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
	main.Handle("/api/v2/", http.StripPrefix("/api/v2", v2))

	return main
}

/*
=============================================================================
 WHEN TO USE THIRD-PARTY ROUTERS
=============================================================================

Since Go 1.22, the standard ServeMux handles the vast majority of routing
needs. But there are cases where a third-party router adds value:

 chi (github.com/go-chi/chi):
  - Built-in middleware chaining
  - Route groups with shared middleware
  - Elegant API: r.Route("/users", func(r chi.Router) { ... })
  - Very lightweight, follows stdlib conventions

 gorilla/mux (now community-maintained):
  - Regex path matching: "/users/{id:[0-9]+}"
  - Host-based routing
  - Historical significance (was THE Go router for years)

 echo / gin:
  - Full frameworks, not just routers
  - Built-in middleware, validation, binding
  - Higher-level abstractions

The Go community has increasingly moved toward "use the stdlib" since 1.22.
The advice is: start with net/http. If you find yourself fighting it, THEN
consider chi (which wraps net/http cleanly). Only reach for echo/gin if
you want a full framework experience.

=============================================================================
*/

// DemoCustomNotFound shows how to handle unmatched routes with a custom 404.
// The default ServeMux returns a plain "404 page not found" response.
// You can override this by registering a catch-all handler.
func DemoCustomNotFound() *http.ServeMux {
	mux := http.NewServeMux()

	// Register your actual routes
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Custom 404 — this catches anything the above routes don't match.
	// The "/" pattern matches everything because every path starts with "/".
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "not found",
			"message": fmt.Sprintf("no route matches %s %s", r.Method, r.URL.Path),
		})
	})

	return mux
}

/*
=============================================================================
 BUILDING A COMPLETE ROUTE TABLE
=============================================================================

In real applications, you organize routes into logical groups. Here's a
pattern that scales well:

  1. Define route groups (users, posts, admin, etc.)
  2. Each group is a function that takes a mux and registers its routes
  3. A top-level function assembles all groups

This keeps your routing code organized even as your API grows. Each
route group can live in its own file if needed.

=============================================================================
*/

// RegisterRoutes demonstrates organizing routes into logical groups.
// This is the pattern you'll use in real applications.
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Each function registers its own routes
	registerUserRoutes(mux)
	registerPostRoutes(mux)
	registerHealthRoutes(mux)

	return mux
}

func registerUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"id": "1", "name": "Alice"},
			{"id": "2", "name": "Bob"},
		})
	})

	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "name": "User " + id})
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "created"})
	})
}

func registerPostRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"id": "1", "title": "First Post"},
		})
	})

	mux.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id, "title": "Post " + id})
	})
}

func registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
}

// Ensure strings is used (referenced in exercises too)
var _ = strings.TrimSpace
