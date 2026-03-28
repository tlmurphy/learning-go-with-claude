package proto

import (
	"context"
	"time"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest holds the data needed to register a new user.
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginRequest holds login credentials.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileRequest holds the fields that can be updated.
type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// TokenPair holds the authentication tokens returned on login.
type TokenPair struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

// UserService defines the contract for user operations.
// Both the service implementation and the API gateway depend on this interface.
type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (User, error)
	Login(ctx context.Context, req LoginRequest) (TokenPair, error)
	GetProfile(ctx context.Context, userID string) (User, error)
	UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (User, error)
}

// ErrorCode represents a gRPC-style status code for service errors.
type ErrorCode int

const (
	CodeOK ErrorCode = iota
	CodeNotFound
	CodeAlreadyExists
	CodeUnauthenticated
	CodePermissionDenied
	CodeInvalidArgument
	CodeInternal
)

// ServiceError is a typed error that carries a status code.
// Use errors.As in the gateway to extract the code and map it to HTTP status.
type ServiceError struct {
	Code    ErrorCode
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// Helper constructors for common errors.

func ErrNotFound(msg string) *ServiceError {
	return &ServiceError{Code: CodeNotFound, Message: msg}
}

func ErrAlreadyExists(msg string) *ServiceError {
	return &ServiceError{Code: CodeAlreadyExists, Message: msg}
}

func ErrUnauthenticated(msg string) *ServiceError {
	return &ServiceError{Code: CodeUnauthenticated, Message: msg}
}

func ErrPermissionDenied(msg string) *ServiceError {
	return &ServiceError{Code: CodePermissionDenied, Message: msg}
}

func ErrInvalidArgument(msg string) *ServiceError {
	return &ServiceError{Code: CodeInvalidArgument, Message: msg}
}
