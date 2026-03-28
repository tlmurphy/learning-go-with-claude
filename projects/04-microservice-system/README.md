# Project 04: Microservice System

**Prerequisite modules:** 22-24 (Advanced Web)

## Overview

Build a two-service system: a **User Service** (gRPC-style, using Go interfaces)
and an **API Gateway** (REST) that translates HTTP requests into calls to the
User Service.

This project teaches you how services communicate, how authentication flows
across boundaries, and how to translate between protocols. Instead of actual
gRPC/protobuf (which adds tooling complexity), you will define the service
contract as Go interfaces and types — the same mental model, without the code
generation.

## Services

### User Service
An internal service that owns user data and authentication logic.

- **Register** — create a new user with email/password
- **Login** — verify credentials, return a token
- **GetProfile** — fetch user profile by ID
- **UpdateProfile** — update user details

### API Gateway
A public-facing REST API that delegates to the User Service.

- Exposes REST endpoints to external clients
- Translates HTTP requests into User Service calls
- Translates User Service errors into proper HTTP status codes
- Handles JWT validation for protected endpoints
- Adds request logging and validation

## Requirements

### Shared Contract (proto/)
- Define `User`, `Credentials`, `Token` types
- Define a `UserService` interface with all operations
- Define error types that map to gRPC-style status codes

### User Service
- Implement the `UserService` interface
- Hash passwords (use `golang.org/x/crypto/bcrypt` or a simple hash for learning)
- Generate JWT tokens on login
- Validate authorization: users can only view/edit their own profile
- Structured logging for all operations

### API Gateway
- REST endpoints:
  ```
  POST   /api/v1/register      → UserService.Register
  POST   /api/v1/login         → UserService.Login
  GET    /api/v1/profile       → UserService.GetProfile (auth required)
  PUT    /api/v1/profile       → UserService.UpdateProfile (auth required)
  ```
- JWT middleware for protected routes
- Request validation (check required fields, email format, password strength)
- Translate service errors to HTTP status codes:
  - NotFound → 404
  - AlreadyExists → 409
  - Unauthenticated → 401
  - PermissionDenied → 403
  - InvalidArgument → 400
- Health check endpoint

### Cross-Cutting
- Structured logging in both services
- Health checks for both services
- Clean separation: the gateway depends only on the `UserService` interface, never on the implementation

## Architecture

```
projects/04-microservice-system/
  proto/
    user.go         — shared types and UserService interface
    errors.go       — service error types
  user-service/
    main.go         — service entry point
    service.go      — UserService implementation
    store.go        — in-memory user storage
  api-gateway/
    main.go         — gateway entry point
    handler.go      — REST handlers
    middleware.go   — JWT auth middleware
  pkg/auth/
    jwt.go          — JWT creation and validation
```

## Hints

<details>
<summary>Service interface design</summary>

```go
// In proto/user.go
type UserService interface {
    Register(ctx context.Context, req RegisterRequest) (User, error)
    Login(ctx context.Context, req LoginRequest) (TokenPair, error)
    GetProfile(ctx context.Context, userID string) (User, error)
    UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (User, error)
}
```

Both the implementation and the gateway depend on this interface.

</details>

<details>
<summary>Error translation</summary>

Define typed errors in proto/ that carry a code:

```go
type ServiceError struct {
    Code    ErrorCode
    Message string
}

func (e *ServiceError) Error() string { return e.Message }
```

In the gateway, use `errors.As` to check for `*ServiceError` and map the
code to an HTTP status.

</details>

<details>
<summary>JWT flow</summary>

1. User calls `POST /api/v1/login` with credentials
2. Gateway calls `UserService.Login`, which returns a JWT token
3. For protected routes, gateway middleware extracts the token from the
   `Authorization: Bearer <token>` header, validates it, and puts the
   user ID into the request context
4. Handlers read the user ID from context and pass it to service calls

</details>

<details>
<summary>Wiring it together</summary>

In a real system, the gateway would make network calls (gRPC) to the user
service. For this project, you can wire them in-process: create the
UserService implementation and pass it directly to the gateway handlers.
This keeps things simple while maintaining the same architecture.

Later, you could swap the direct call for a gRPC client that implements
the same interface — that is the power of coding to an interface.

</details>

## Stretch Goals

- **Second service** — add a Product Service with its own interface, and expose it through the gateway
- **Service discovery** — implement simple config-based service discovery (a registry of service addresses)
- **Request tracing** — propagate a trace ID from the gateway through to the user service via context
- **Circuit breaker** — if the user service is slow or failing, have the gateway fail fast instead of hanging
