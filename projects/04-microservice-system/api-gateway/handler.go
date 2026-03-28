package main

import (
	"net/http"

	"learning-go-with-claude/projects/04-microservice-system/proto"
)

// GatewayHandler handles REST requests and delegates to the UserService.
//
// TODO: Implement each handler:
//   - Decode JSON request body
//   - Call the appropriate UserService method
//   - Translate ServiceError codes to HTTP status codes
//   - Encode JSON response

type GatewayHandler struct {
	userService proto.UserService
}

func NewGatewayHandler(us proto.UserService) *GatewayHandler {
	return &GatewayHandler{userService: us}
}

// Register handles POST /api/v1/register
func (h *GatewayHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode RegisterRequest, validate, call service, return 201
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// Login handles POST /api/v1/login
func (h *GatewayHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Decode LoginRequest, call service, return TokenPair
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// GetProfile handles GET /api/v1/profile (auth required)
func (h *GatewayHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract user ID from context (set by auth middleware), call service
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// UpdateProfile handles PUT /api/v1/profile (auth required)
func (h *GatewayHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract user ID from context, decode body, call service
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
