// Package grpcmod teaches gRPC service patterns in Go. Since we can't run
// the protobuf compiler in this tutorial, we define Go interfaces and types
// that mirror what protoc would generate, letting you focus on the service
// implementation patterns that matter most in production.
package grpcmod

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 gRPC SERVICES IN GO
=============================================================================

gRPC is a high-performance RPC (Remote Procedure Call) framework originally
developed at Google. If REST is the lingua franca of public APIs, gRPC is
the backbone of internal microservice communication at most large tech
companies.

Why does gRPC exist when we already have REST?

1. PERFORMANCE: gRPC uses Protocol Buffers (protobuf) for serialization,
   which is binary and dramatically smaller/faster than JSON. A typical
   protobuf message is 3-10x smaller and 20-100x faster to parse than
   its JSON equivalent.

2. TYPE SAFETY: You define your API in a .proto file, and the compiler
   generates strongly-typed client and server code. No more hoping that
   the JSON field "user_id" is actually an integer.

3. STREAMING: gRPC natively supports streaming in both directions. REST
   can do server-sent events or WebSockets, but they're bolted on.
   gRPC streaming is a first-class citizen.

4. HTTP/2: gRPC runs on HTTP/2, giving you multiplexing (many requests
   over one connection), header compression, and binary framing.

5. CODE GENERATION: Write one .proto file, generate clients in any
   language. Your Go service and Python ML pipeline speak the same
   protocol without hand-written serialization.

=============================================================================
 PROTOCOL BUFFERS BASICS
=============================================================================

Protocol Buffers (protobuf) is the serialization format and interface
definition language (IDL) used by gRPC. Here's what a .proto file looks
like conceptually:

  syntax = "proto3";
  package userservice;

  // A message is like a Go struct
  message User {
    string id = 1;           // Field number 1 (used in binary encoding)
    string name = 2;
    string email = 3;
    int64 created_at = 4;    // Unix timestamp
  }

  message GetUserRequest {
    string id = 1;
  }

  message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
  }

  // A service is like a Go interface
  service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
    rpc CreateUser(User) returns (User);
  }

The protobuf compiler (protoc) would generate Go code from this, including:
  - Struct types for each message (with serialization methods)
  - An interface for the server to implement
  - A client that calls the server over the network

Since we can't run protoc here, we'll define equivalent Go types manually.

Key protobuf concepts:
  - Field numbers (1, 2, 3...) are used in the binary encoding, not names
  - proto3 uses zero values as defaults (sound familiar, Go developers?)
  - Repeated fields map to slices in Go
  - Oneof fields map to interface types in Go
  - Enums map to int32 with named constants

=============================================================================
 gRPC vs REST: WHEN TO USE WHICH
=============================================================================

Use gRPC when:
  - Service-to-service communication (internal microservices)
  - Performance matters (high throughput, low latency)
  - You need streaming (real-time updates, large data transfers)
  - You want strong contracts between teams/services
  - You're building in a polyglot environment (many languages)

Use REST when:
  - Public-facing APIs (browsers can't easily use gRPC directly)
  - Simple CRUD operations where JSON is fine
  - You want curl-friendly debugging
  - Your team is small and the overhead of proto files isn't worth it
  - You're building a web application (REST is more natural with HTML)

In practice, many companies use BOTH: gRPC between microservices,
REST (often via gRPC-Gateway) for public APIs.

=============================================================================
 THE FOUR TYPES OF RPC
=============================================================================

gRPC supports four communication patterns:

1. UNARY: Client sends one request, server sends one response.
   Like a regular function call. This is the most common pattern.
   Example: GetUser(id) -> User

2. SERVER STREAMING: Client sends one request, server sends a stream
   of responses. Like an iterator.
   Example: ListUsers(filter) -> stream of Users

3. CLIENT STREAMING: Client sends a stream of requests, server sends
   one response. Like uploading chunks.
   Example: UploadFile(stream of chunks) -> UploadResult

4. BIDIRECTIONAL STREAMING: Both client and server send streams.
   Like a chat conversation.
   Example: Chat(stream of messages) <-> stream of messages

=============================================================================
*/

// ---- Types that mirror what protoc would generate ----

// User represents a user message, similar to what protobuf would generate.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt int64 // Unix timestamp
}

// GetUserRequest mirrors a protobuf request message.
type GetUserRequest struct {
	ID string
}

// ListUsersRequest mirrors a protobuf request with pagination.
type ListUsersRequest struct {
	PageSize  int32
	PageToken string
}

// CreateUserRequest mirrors a protobuf request for creating a user.
type CreateUserRequest struct {
	Name  string
	Email string
}

// UpdateUserRequest mirrors a protobuf request for updating a user.
type UpdateUserRequest struct {
	ID    string
	Name  string
	Email string
}

// DeleteUserRequest mirrors a protobuf request for deleting a user.
type DeleteUserRequest struct {
	ID string
}

// Empty mirrors google.protobuf.Empty — used when no data is needed.
type Empty struct{}

/*
=============================================================================
 SERVICE INTERFACES
=============================================================================

In real gRPC, the protoc compiler generates a server interface that you
implement. Here we define it manually. Notice how every method takes a
context.Context as its first parameter — this is mandatory in gRPC and
is how deadlines, cancellation, and metadata flow through the system.
=============================================================================
*/

// UserServiceServer defines the interface that a gRPC server would implement.
// In real gRPC, this interface is generated by protoc-gen-go-grpc.
//
// Every method takes context.Context — this carries deadlines, cancellation
// signals, and metadata (like HTTP headers). ALWAYS respect context in your
// implementations.
type UserServiceServer interface {
	GetUser(ctx context.Context, req *GetUserRequest) (*User, error)
	ListUsers(req *ListUsersRequest, stream UserStream) error
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest) (*Empty, error)
}

// UserStream represents a server-side stream for sending users.
// In real gRPC, this would be generated as UserService_ListUsersServer.
type UserStream interface {
	Send(user *User) error
	Context() context.Context
}

// DemoServiceInterface shows how gRPC service interfaces work.
func DemoServiceInterface() {
	fmt.Println("=== gRPC Service Interface Pattern ===")
	fmt.Println()
	fmt.Println("In real gRPC, you'd define a .proto file and run protoc to generate:")
	fmt.Println("  1. Message types (Go structs with serialization)")
	fmt.Println("  2. A server interface to implement")
	fmt.Println("  3. A client that calls the server over the network")
	fmt.Println()
	fmt.Println("The generated server interface always looks like:")
	fmt.Println("  type UserServiceServer interface {")
	fmt.Println("      GetUser(ctx context.Context, req *GetUserRequest) (*User, error)")
	fmt.Println("      ListUsers(req *ListUsersRequest, stream UserService_ListUsersServer) error")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("Your job is to implement this interface with your business logic.")
}

/*
=============================================================================
 gRPC STATUS CODES AND ERROR HANDLING
=============================================================================

gRPC has its own set of status codes, similar to HTTP status codes but
designed for RPC semantics. They're defined in the google.golang.org/grpc/codes
package. Here's the mapping you need to internalize:

  gRPC Code          | HTTP Equiv | When to Use
  -------------------|------------|--------------------------------------------
  OK                 | 200        | Success
  Cancelled          | 499        | Client cancelled the request
  Unknown            | 500        | Catch-all for unknown errors
  InvalidArgument    | 400        | Bad input from the client
  DeadlineExceeded   | 504        | Operation took too long
  NotFound           | 404        | Requested resource doesn't exist
  AlreadyExists      | 409        | Resource already exists (e.g., duplicate)
  PermissionDenied   | 403        | Caller doesn't have permission
  Unauthenticated    | 401        | No valid auth credentials
  ResourceExhausted  | 429        | Quota exceeded or rate limited
  FailedPrecondition | 400        | System not in required state
  Aborted            | 409        | Operation aborted (e.g., concurrency)
  Unimplemented      | 501        | Method not implemented
  Internal           | 500        | Internal server error
  Unavailable        | 503        | Service temporarily unavailable
  DataLoss           | 500        | Unrecoverable data loss

The critical insight: NEVER return raw Go errors from gRPC methods. Always
wrap them in a status error so the client gets a meaningful code. An
unwrapped error becomes "Unknown" on the client side, which is useless
for debugging.

=============================================================================
*/

// Code represents a gRPC status code, mirroring google.golang.org/grpc/codes.Code.
type Code int

const (
	OK                 Code = 0
	Cancelled          Code = 1
	Unknown            Code = 2
	InvalidArgument    Code = 3
	DeadlineExceeded   Code = 4
	NotFound           Code = 5
	AlreadyExists      Code = 6
	PermissionDenied   Code = 7
	ResourceExhausted  Code = 8
	FailedPrecondition Code = 9
	Aborted            Code = 10
	Unimplemented      Code = 12
	Internal           Code = 13
	Unavailable        Code = 14
	DataLoss           Code = 15
	Unauthenticated    Code = 16
)

// codeNames maps codes to their string names for readable output.
var codeNames = map[Code]string{
	OK: "OK", Cancelled: "Cancelled", Unknown: "Unknown",
	InvalidArgument: "InvalidArgument", DeadlineExceeded: "DeadlineExceeded",
	NotFound: "NotFound", AlreadyExists: "AlreadyExists",
	PermissionDenied: "PermissionDenied", ResourceExhausted: "ResourceExhausted",
	FailedPrecondition: "FailedPrecondition", Aborted: "Aborted",
	Unimplemented: "Unimplemented", Internal: "Internal",
	Unavailable: "Unavailable", DataLoss: "DataLoss",
	Unauthenticated: "Unauthenticated",
}

// String returns the name of the status code.
func (c Code) String() string {
	if name, ok := codeNames[c]; ok {
		return name
	}
	return fmt.Sprintf("Code(%d)", int(c))
}

// Status represents a gRPC status, mirroring google.golang.org/grpc/status.Status.
// In real gRPC, you'd use status.New(codes.NotFound, "user not found").
type Status struct {
	code    Code
	message string
}

// StatusError returns a new Status as an error.
func StatusError(code Code, msg string) error {
	return &Status{code: code, message: msg}
}

// StatusErrorf returns a new Status error with a formatted message.
func StatusErrorf(code Code, format string, args ...any) error {
	return &Status{code: code, message: fmt.Sprintf(format, args...)}
}

// Error implements the error interface.
func (s *Status) Error() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", s.code, s.message)
}

// Code returns the status code.
func (s *Status) Code() Code {
	return s.code
}

// Message returns the status message.
func (s *Status) Message() string {
	return s.message
}

// FromError extracts a Status from an error. If the error is not a Status,
// it returns a Status with code Unknown.
func FromError(err error) *Status {
	if err == nil {
		return &Status{code: OK, message: ""}
	}
	if s, ok := err.(*Status); ok {
		return s
	}
	return &Status{code: Unknown, message: err.Error()}
}

// DemoErrorHandling shows proper gRPC error handling patterns.
func DemoErrorHandling() {
	fmt.Println("=== gRPC Error Handling ===")
	fmt.Println()

	// Create a status error — this is how you return errors from gRPC methods
	err := StatusError(NotFound, "user 'abc123' not found")
	fmt.Printf("Error: %v\n", err)

	// Extract the status from the error — this is how clients inspect errors
	st := FromError(err)
	fmt.Printf("Code: %s, Message: %s\n", st.Code(), st.Message())

	// A raw Go error becomes Unknown — this is why you should always use Status
	rawErr := fmt.Errorf("database connection failed")
	st2 := FromError(rawErr)
	fmt.Printf("Raw error code: %s (not helpful for the client!)\n", st2.Code())
	fmt.Println()
	fmt.Println("Key takeaway: Always return Status errors from gRPC methods.")
	fmt.Println("The client needs the code to decide what to do (retry? show error? etc.)")
}

/*
=============================================================================
 INTERCEPTORS (gRPC MIDDLEWARE)
=============================================================================

Interceptors are gRPC's equivalent of HTTP middleware. They sit between
the client/server and the handler, letting you add cross-cutting concerns
like logging, authentication, metrics, and tracing.

There are four types of interceptors:
  1. Unary Server Interceptor — wraps unary (single request/response) handlers
  2. Stream Server Interceptor — wraps streaming handlers
  3. Unary Client Interceptor — wraps unary client calls
  4. Stream Client Interceptor — wraps streaming client calls

The pattern is similar to HTTP middleware: you receive the handler (or
"invoker" in client terms), do something before/after, and call through.

Real signature from grpc package:

  type UnaryServerInterceptor func(
      ctx context.Context,
      req interface{},
      info *UnaryServerInfo,
      handler UnaryHandler,
  ) (interface{}, error)

=============================================================================
*/

// UnaryServerInfo contains information about the unary RPC being handled.
type UnaryServerInfo struct {
	FullMethod string // e.g., "/userservice.UserService/GetUser"
}

// UnaryHandler is the handler function for a unary RPC.
type UnaryHandler func(ctx context.Context, req any) (any, error)

// UnaryServerInterceptor is a function that intercepts unary RPC calls.
type UnaryServerInterceptor func(
	ctx context.Context,
	req any,
	info *UnaryServerInfo,
	handler UnaryHandler,
) (any, error)

// DemoLoggingInterceptor shows a simple logging interceptor.
func DemoLoggingInterceptor() UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *UnaryServerInfo,
		handler UnaryHandler,
	) (any, error) {
		start := time.Now()
		fmt.Printf("[gRPC] %s started\n", info.FullMethod)

		// Call the actual handler
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			st := FromError(err)
			fmt.Printf("[gRPC] %s failed: code=%s duration=%v\n",
				info.FullMethod, st.Code(), duration)
		} else {
			fmt.Printf("[gRPC] %s completed: duration=%v\n",
				info.FullMethod, duration)
		}

		return resp, err
	}
}

// ChainUnaryInterceptors chains multiple interceptors into one.
// The first interceptor in the slice is the outermost (runs first).
// This mirrors grpc.ChainUnaryInterceptor from the real library.
func ChainUnaryInterceptors(interceptors ...UnaryServerInterceptor) UnaryServerInterceptor {
	if len(interceptors) == 0 {
		// No-op interceptor — just call the handler directly
		return func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error) {
			return handler(ctx, req)
		}
	}
	if len(interceptors) == 1 {
		return interceptors[0]
	}

	return func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error) {
		// Build the chain from inside out: each interceptor wraps the next
		currentHandler := handler
		for i := len(interceptors) - 1; i > 0; i-- {
			// Capture loop variable
			interceptor := interceptors[i]
			next := currentHandler
			currentHandler = func(ctx context.Context, req any) (any, error) {
				return interceptor(ctx, req, info, next)
			}
		}
		return interceptors[0](ctx, req, info, currentHandler)
	}
}

/*
=============================================================================
 METADATA (gRPC HEADERS)
=============================================================================

Metadata is gRPC's equivalent of HTTP headers. It's a map of string keys
to string (or binary) values, sent alongside requests and responses.

Common uses:
  - Authentication tokens ("authorization": "Bearer <token>")
  - Request IDs for tracing ("x-request-id": "abc-123")
  - Custom routing hints
  - Client version information

In real gRPC, you use the metadata package:

  import "google.golang.org/grpc/metadata"

  // Client side: attach metadata to outgoing context
  md := metadata.Pairs("authorization", "Bearer token123")
  ctx := metadata.NewOutgoingContext(ctx, md)

  // Server side: extract metadata from incoming context
  md, ok := metadata.FromIncomingContext(ctx)
  tokens := md.Get("authorization")

We'll simulate this with a simple map-based approach.

=============================================================================
*/

// Metadata represents gRPC metadata (key-value pairs like HTTP headers).
// Keys are always lowercase. Values can have multiple entries per key.
type Metadata map[string][]string

// NewMetadata creates metadata from key-value pairs. Keys are lowercased.
// Panics if an odd number of strings is provided.
func NewMetadata(kv ...string) Metadata {
	if len(kv)%2 != 0 {
		panic("metadata: odd number of key-value pairs")
	}
	md := make(Metadata)
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		md[key] = append(md[key], kv[i+1])
	}
	return md
}

// Get returns the first value for a key, or empty string if not present.
func (md Metadata) Get(key string) string {
	vals := md[strings.ToLower(key)]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

// Set sets a key to a single value, replacing any existing values.
func (md Metadata) Set(key, value string) {
	md[strings.ToLower(key)] = []string{value}
}

// Context keys for storing metadata
type mdIncomingKey struct{}
type mdOutgoingKey struct{}

// NewIncomingContext attaches incoming metadata to a context (server side).
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

// FromIncomingContext extracts incoming metadata from a context (server side).
func FromIncomingContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(mdIncomingKey{}).(Metadata)
	return md, ok
}

// NewOutgoingContext attaches outgoing metadata to a context (client side).
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, md)
}

// DemoMetadata shows how gRPC metadata works.
func DemoMetadata() {
	fmt.Println("=== gRPC Metadata ===")
	fmt.Println()

	// Create metadata (like HTTP headers)
	md := NewMetadata(
		"authorization", "Bearer token123",
		"x-request-id", "req-456",
	)

	// Attach to context (server side receives this)
	ctx := NewIncomingContext(context.Background(), md)

	// Extract on the server side
	inMD, ok := FromIncomingContext(ctx)
	if ok {
		fmt.Printf("Auth: %s\n", inMD.Get("authorization"))
		fmt.Printf("Request ID: %s\n", inMD.Get("x-request-id"))
	}
}

/*
=============================================================================
 DEADLINES AND TIMEOUTS
=============================================================================

This is perhaps the most important production gRPC concept:
ALWAYS SET DEADLINES.

In gRPC, a deadline is an absolute point in time by which the RPC must
complete. If the deadline passes, the call fails with DeadlineExceeded.
Unlike HTTP timeouts, gRPC deadlines propagate through the entire call
chain — if Service A calls Service B calls Service C, and A's deadline
expires, all three stop working on the request.

Why this matters:
  - Without deadlines, a slow downstream service can cause requests to
    pile up, eventually crashing your entire system (cascading failure).
  - With deadlines, slow requests fail fast, freeing resources.
  - Google's internal rule: every RPC must have a deadline. No exceptions.

In Go, deadlines are implemented via context.Context:

  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()
  resp, err := client.GetUser(ctx, req)

The server should check ctx.Err() or use ctx.Done() to stop work early
if the deadline has passed.

=============================================================================
*/

// DemoDeadlines shows how deadlines work in gRPC.
func DemoDeadlines() {
	fmt.Println("=== gRPC Deadlines ===")
	fmt.Println()

	// Simulate a handler that respects deadlines
	handler := func(ctx context.Context) (string, error) {
		// Simulate work that takes 100ms
		select {
		case <-time.After(100 * time.Millisecond):
			return "result", nil
		case <-ctx.Done():
			// Deadline exceeded — stop working and return error
			return "", StatusError(DeadlineExceeded, "operation timed out")
		}
	}

	// Case 1: Plenty of time — succeeds
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()
	result, err := handler(ctx1)
	fmt.Printf("With 1s deadline: result=%q, err=%v\n", result, err)

	// Case 2: Too little time — fails with DeadlineExceeded
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel2()
	result, err = handler(ctx2)
	fmt.Printf("With 10ms deadline: result=%q, err=%v\n", result, err)

	fmt.Println()
	fmt.Println("Rule: ALWAYS set deadlines on gRPC calls.")
	fmt.Println("A missing deadline is a production incident waiting to happen.")
}

/*
=============================================================================
 HEALTH CHECKING
=============================================================================

The gRPC Health Checking Protocol is a standardized way for services to
report their health status. Load balancers and orchestrators (Kubernetes)
use this to decide whether to send traffic to an instance.

The protocol defines a simple service:

  service Health {
    rpc Check(HealthCheckRequest) returns (HealthCheckResponse);
    rpc Watch(HealthCheckRequest) returns (stream HealthCheckResponse);
  }

  message HealthCheckResponse {
    enum ServingStatus {
      UNKNOWN = 0;
      SERVING = 1;
      NOT_SERVING = 2;
      SERVICE_UNKNOWN = 3;
    }
    ServingStatus status = 1;
  }

Key points:
  - Empty service name means "overall server health"
  - Named services let you check individual components
  - Watch provides streaming health updates (for client-side load balancing)
  - Kubernetes gRPC health probes use this protocol directly

=============================================================================
*/

// ServingStatus represents the health of a service.
type ServingStatus int

const (
	StatusUnknown        ServingStatus = 0
	StatusServing        ServingStatus = 1
	StatusNotServing     ServingStatus = 2
	StatusServiceUnknown ServingStatus = 3
)

// String returns the name of the serving status.
func (s ServingStatus) String() string {
	switch s {
	case StatusUnknown:
		return "UNKNOWN"
	case StatusServing:
		return "SERVING"
	case StatusNotServing:
		return "NOT_SERVING"
	case StatusServiceUnknown:
		return "SERVICE_UNKNOWN"
	default:
		return fmt.Sprintf("ServingStatus(%d)", int(s))
	}
}

// HealthCheckRequest mirrors the gRPC health check request.
type HealthCheckRequest struct {
	Service string // Empty string means overall server health
}

// HealthCheckResponse mirrors the gRPC health check response.
type HealthCheckResponse struct {
	Status ServingStatus
}

// HealthServer implements the gRPC Health Checking Protocol.
// It's safe for concurrent use.
type HealthServer struct {
	mu       sync.RWMutex
	statuses map[string]ServingStatus
}

// NewHealthServer creates a new health server with all services unknown.
func NewHealthServer() *HealthServer {
	return &HealthServer{
		statuses: make(map[string]ServingStatus),
	}
}

// SetServingStatus sets the health status for a service.
// Use empty string for overall server health.
func (h *HealthServer) SetServingStatus(service string, status ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.statuses[service] = status
}

// Check returns the health status for a service.
func (h *HealthServer) Check(_ context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	status, ok := h.statuses[req.Service]
	if !ok {
		return nil, StatusError(NotFound, fmt.Sprintf("service %q not registered", req.Service))
	}
	return &HealthCheckResponse{Status: status}, nil
}

// DemoHealthCheck shows the health checking protocol in action.
func DemoHealthCheck() {
	fmt.Println("=== gRPC Health Checking ===")
	fmt.Println()

	hs := NewHealthServer()

	// Register services
	hs.SetServingStatus("", StatusServing)              // Overall: healthy
	hs.SetServingStatus("user-service", StatusServing)   // User service: healthy
	hs.SetServingStatus("cache", StatusNotServing)       // Cache: unhealthy

	// Check health
	ctx := context.Background()

	resp, err := hs.Check(ctx, &HealthCheckRequest{Service: ""})
	if err == nil {
		fmt.Printf("Overall: %s\n", resp.Status)
	}

	resp, err = hs.Check(ctx, &HealthCheckRequest{Service: "cache"})
	if err == nil {
		fmt.Printf("Cache: %s\n", resp.Status)
	}

	_, err = hs.Check(ctx, &HealthCheckRequest{Service: "unknown"})
	if err != nil {
		fmt.Printf("Unknown service: %v\n", err)
	}
}

/*
=============================================================================
 gRPC-GATEWAY: REST FACADE FOR gRPC
=============================================================================

gRPC-Gateway is a protoc plugin that generates a reverse proxy server
which translates RESTful JSON API calls into gRPC calls. This lets you
write your service once in gRPC and expose it as both gRPC and REST.

How it works:
  1. You annotate your .proto file with HTTP mappings
  2. The plugin generates a reverse proxy
  3. REST clients talk to the proxy, which translates to gRPC

  service UserService {
    rpc GetUser(GetUserRequest) returns (User) {
      option (google.api.http) = {
        get: "/v1/users/{id}"
      };
    }
  }

This means: GET /v1/users/abc123 → GetUser({id: "abc123"})

Benefits:
  - Single source of truth (.proto file)
  - Internal services use gRPC (fast)
  - External clients use REST (familiar)
  - Swagger/OpenAPI docs generated automatically

=============================================================================
 CONNECTION MANAGEMENT AND LOAD BALANCING
=============================================================================

gRPC connections are long-lived (unlike HTTP/1.1 which often creates new
connections). This changes how load balancing works:

Client-Side Load Balancing:
  - gRPC supports built-in client-side load balancing
  - The client maintains connections to multiple backends
  - Each RPC is routed to an appropriate backend
  - Policies: pick_first, round_robin, or custom

Server-Side Load Balancing:
  - Traditional L4/L7 load balancers work, but need HTTP/2 support
  - Envoy, Nginx (with gRPC module), and cloud LBs support this
  - Be careful: L4 balancing with long-lived connections means
    one connection = one backend (defeats the purpose)

Connection Keepalive:
  - gRPC has built-in keepalive pings
  - Important for detecting dead connections quickly
  - Configure on both client and server sides

Best Practices:
  - Use client-side load balancing for service mesh architectures
  - Use server-side L7 load balancing for simpler deployments
  - Always configure keepalive parameters
  - Monitor connection counts and states

=============================================================================
 RETRY LOGIC AND BACKOFF
=============================================================================

Network calls fail. gRPC calls are network calls. Therefore, gRPC calls
fail. You need retry logic.

But not all errors are retryable:
  - Unavailable: YES, retry (server might come back)
  - DeadlineExceeded: MAYBE (depends on whether the operation is idempotent)
  - Internal: MAYBE (might be transient)
  - InvalidArgument: NO (sending the same bad request won't help)
  - NotFound: NO (retrying won't create the resource)
  - Unauthenticated: NO (need new credentials first)

Exponential backoff is the standard retry strategy:
  1st retry: wait 100ms
  2nd retry: wait 200ms
  3rd retry: wait 400ms
  4th retry: wait 800ms
  ...with jitter (random variation) to prevent thundering herd

=============================================================================
*/

// RetryConfig controls retry behavior for gRPC calls.
type RetryConfig struct {
	MaxRetries  int           // Maximum number of retry attempts
	InitialWait time.Duration // Wait time before first retry
	MaxWait     time.Duration // Maximum wait time between retries
	Multiplier  float64       // Backoff multiplier (typically 2.0)
}

// DefaultRetryConfig returns sensible defaults for retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:  3,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     5 * time.Second,
		Multiplier:  2.0,
	}
}

// IsRetryable returns whether an error's status code is worth retrying.
func IsRetryable(err error) bool {
	st := FromError(err)
	switch st.Code() {
	case Unavailable, ResourceExhausted, Aborted:
		return true
	default:
		return false
	}
}

// CalculateBackoff returns the wait duration for a given attempt number.
// It uses exponential backoff capped at MaxWait.
func (rc RetryConfig) CalculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return rc.InitialWait
	}
	backoff := float64(rc.InitialWait) * math.Pow(rc.Multiplier, float64(attempt))
	if backoff > float64(rc.MaxWait) {
		backoff = float64(rc.MaxWait)
	}
	return time.Duration(backoff)
}

// DemoRetryLogic shows retry with exponential backoff.
func DemoRetryLogic() {
	fmt.Println("=== gRPC Retry Logic ===")
	fmt.Println()

	config := DefaultRetryConfig()

	fmt.Println("Backoff schedule:")
	for i := 0; i <= config.MaxRetries; i++ {
		wait := config.CalculateBackoff(i)
		fmt.Printf("  Attempt %d: wait %v\n", i, wait)
	}

	fmt.Println()
	fmt.Println("Only retry on transient errors (Unavailable, ResourceExhausted, Aborted).")
	fmt.Println("Never retry InvalidArgument, NotFound, or Unauthenticated.")
}

/*
=============================================================================
 PUTTING IT ALL TOGETHER
=============================================================================

A production gRPC service typically has:

  1. Proto definitions (.proto files)
  2. Generated code (from protoc)
  3. Server implementation (your business logic)
  4. Interceptors (logging, auth, metrics, tracing)
  5. Health checking
  6. Client wrappers (with retry logic)
  7. Configuration (timeouts, connection pools, TLS)
  8. Graceful shutdown

The server setup looks something like:

  server := grpc.NewServer(
      grpc.ChainUnaryInterceptor(
          loggingInterceptor,
          authInterceptor,
          metricsInterceptor,
      ),
  )

  pb.RegisterUserServiceServer(server, &userServer{})
  healthpb.RegisterHealthServer(server, healthServer)

  lis, _ := net.Listen("tcp", ":50051")
  server.Serve(lis)

And the client:

  conn, _ := grpc.Dial("localhost:50051",
      grpc.WithTransportCredentials(insecure.NewCredentials()),
      grpc.WithUnaryInterceptor(retryInterceptor),
  )
  defer conn.Close()

  client := pb.NewUserServiceClient(conn)
  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()

  user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "abc"})

=============================================================================
*/
