package routing

import (
	"net/http"
)

/*
=============================================================================
 EXERCISES: Routing and URL Patterns
=============================================================================

These exercises build your skills with Go 1.22+ routing patterns. You'll
work with method routing, path parameters, wildcards, subrouting, and
custom error handling.

All exercises return *http.ServeMux so they can be tested with httptest.

=============================================================================
*/

// Exercise 1: MethodRouter
//
// Create a ServeMux that routes different HTTP methods to different handlers
// for the same path "/items".
//
// Requirements:
// - GET /items    -> 200, JSON body: {"method": "GET", "action": "list"}
// - POST /items   -> 201, JSON body: {"method": "POST", "action": "create"}
// - PUT /items    -> 200, JSON body: {"method": "PUT", "action": "update"}
// - DELETE /items -> 200, JSON body: {"method": "DELETE", "action": "delete"}
// - All responses should have Content-Type: application/json
//
// This demonstrates the core improvement of Go 1.22 routing: method dispatch
// at the mux level instead of inside every handler.
func MethodRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 2: PathParamExtractor
//
// Create a ServeMux that extracts path parameters and returns them in the
// response.
//
// Requirements:
// - GET /users/{id} -> 200, JSON body: {"resource": "user", "id": "<id>"}
// - GET /posts/{slug} -> 200, JSON body: {"resource": "post", "slug": "<slug>"}
// - All responses should have Content-Type: application/json
//
// The path parameter values should come directly from r.PathValue().
func PathParamExtractor() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 3: ResourceRouter
//
// Create a complete CRUD route set for a "books" resource.
//
// Requirements:
// - GET /books         -> 200, JSON body: {"action": "list"}
// - POST /books        -> 201, JSON body: {"action": "create"}
// - GET /books/{id}    -> 200, JSON body: {"action": "get", "id": "<id>"}
// - PUT /books/{id}    -> 200, JSON body: {"action": "update", "id": "<id>"}
// - DELETE /books/{id} -> 204, no body
// - All responses (except DELETE) should have Content-Type: application/json
//
// This is the standard REST pattern: collection routes (/books) and
// individual resource routes (/books/{id}).
func ResourceRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 4: WildcardRouter
//
// Create a ServeMux that uses the {path...} wildcard for catch-all routing.
//
// Requirements:
//   - GET /static/{path...} -> 200, JSON body: {"file": "<path>"}
//     where <path> is the wildcard value from r.PathValue("path")
//   - GET /static/ (no path) -> 200, JSON body: {"file": "index.html"}
//     (default to "index.html" when path is empty)
//   - All responses should have Content-Type: application/json
//
// Wildcards are useful for file servers, SPA routing, and proxy handlers.
func WildcardRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 5: VersionedAPI
//
// Create a ServeMux with versioned API routes using subrouting.
//
// Requirements:
//   - GET /api/v1/status -> 200, JSON body: {"version": "v1", "status": "ok"}
//   - GET /api/v2/status -> 200, JSON body: {"version": "v2", "status": "ok"}
//   - GET /api/v1/users  -> 200, JSON body: {"version": "v1", "users": []}
//   - GET /api/v2/users  -> 200, JSON body: {"version": "v2", "users": [], "meta": {}}
//     where "meta" is an empty map (map[string]interface{}{})
//   - All responses should have Content-Type: application/json
//
// Use http.StripPrefix and separate muxes for each version. This pattern
// lets different API versions evolve independently.
func VersionedAPI() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 6: StripPrefixRouter
//
// Create a ServeMux that uses StripPrefix to mount a sub-router.
//
// Requirements:
//   - Create an "admin" sub-router with these routes:
//     GET /dashboard -> 200, JSON body: {"page": "dashboard"}
//     GET /settings  -> 200, JSON body: {"page": "settings"}
//   - Mount it at /admin/ on the main mux using StripPrefix
//   - So GET /admin/dashboard -> 200, JSON body: {"page": "dashboard"}
//   - And GET /admin/settings -> 200, JSON body: {"page": "settings"}
//   - All responses should have Content-Type: application/json
//
// StripPrefix removes the prefix so the sub-router sees clean paths.
func StripPrefixRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 7: CustomErrorRouter
//
// Create a ServeMux with custom 404 and 405 error responses.
//
// Requirements:
//   - Register: GET /api/health -> 200, JSON body: {"status": "ok"}
//   - Register: POST /api/data -> 200, JSON body: {"received": true}
//   - For unmatched routes (404), return JSON:
//     {"error": "not found", "path": "<request path>"}
//     with status 404 and Content-Type: application/json
//   - The catch-all should match "/" to handle all unmatched paths
//
// Note: 405 (Method Not Allowed) is handled automatically by the mux
// when a path matches but the method doesn't. We focus on 404 here.
func CustomErrorRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}

// Exercise 8: BlogRoutes
//
// Build a complete API route table for a blog with posts, comments, and users.
//
// Requirements:
// All routes return JSON with Content-Type: application/json.
// Each response includes a "route" field identifying which handler matched.
//
// Posts:
// - GET /posts              -> 200, {"route": "list_posts"}
// - POST /posts             -> 201, {"route": "create_post"}
// - GET /posts/{id}         -> 200, {"route": "get_post", "id": "<id>"}
// - PUT /posts/{id}         -> 200, {"route": "update_post", "id": "<id>"}
// - DELETE /posts/{id}      -> 204, no body
//
// Comments (nested under posts):
// - GET /posts/{postID}/comments          -> 200, {"route": "list_comments", "post_id": "<postID>"}
// - POST /posts/{postID}/comments         -> 201, {"route": "create_comment", "post_id": "<postID>"}
//
// Users:
// - GET /users              -> 200, {"route": "list_users"}
// - GET /users/{id}         -> 200, {"route": "get_user", "id": "<id>"}
//
// This exercises nested resources and the complete CRUD pattern.
func BlogRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	// YOUR CODE HERE
	return mux
}
