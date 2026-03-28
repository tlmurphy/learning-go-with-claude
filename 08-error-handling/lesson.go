package errorhandling

/*
=============================================================================
 Module 08: Error Handling
=============================================================================

 Error handling in Go is deliberately simple: errors are values. There's
 no try/catch, no exceptions, no stack unwinding. A function that can fail
 returns an error as its last return value, and the caller checks it.

 This might feel tedious at first — you'll write "if err != nil" a LOT.
 But there's a reason for it: explicit error handling forces you to think
 about failure at every step. Nothing is silently swallowed. You always
 know where an error can occur and what happens when it does.

 WHY THIS MATTERS FOR WEB SERVICES:
 Every HTTP handler, every database query, every external API call can
 fail. Go's error handling makes failure paths explicit and visible.
 In production, the difference between "handled errors gracefully" and
 "crashed because we forgot to check an error" is the difference between
 a good service and a 2 AM pager alert.

=============================================================================
*/

import (
	"errors"
	"fmt"
	"strings"
)

// -------------------------------------------------------------------------
// The error Interface
// -------------------------------------------------------------------------

/*
 The error type is just an interface:

   type error interface {
       Error() string
   }

 Any type that has an Error() string method is an error. That's it.
 This simplicity is powerful — you can create rich error types that carry
 additional context while still being compatible with everything that
 expects an error.
*/

// -------------------------------------------------------------------------
// Creating Errors: errors.New and fmt.Errorf
// -------------------------------------------------------------------------

/*
 For simple error messages, use errors.New:
   return errors.New("something went wrong")

 For formatted errors with context, use fmt.Errorf:
   return fmt.Errorf("failed to load user %d: %v", id, err)

 For error wrapping (preserving the original error), use %w:
   return fmt.Errorf("failed to load user %d: %w", id, err)

 The difference between %v and %w is critical:
  - %v formats the error as a string (loses the original error)
  - %w wraps the original error (preserves it for errors.Is/errors.As)
*/

func DemoCreatingErrors() {
	// Simple error
	err1 := errors.New("connection refused")
	fmt.Println(err1) // "connection refused"

	// Formatted error with context
	userID := 42
	err2 := fmt.Errorf("failed to load user %d: %s", userID, err1)
	fmt.Println(err2) // "failed to load user 42: connection refused"

	// Wrapped error — preserves the original for unwrapping
	err3 := fmt.Errorf("service unavailable: %w", err1)
	fmt.Println(err3)                      // "service unavailable: connection refused"
	fmt.Println(errors.Is(err3, err1))     // true — the original error is preserved
	fmt.Println(errors.Is(err2, err1))     // false — %v doesn't preserve it
}

// -------------------------------------------------------------------------
// The if err != nil Pattern
// -------------------------------------------------------------------------

/*
 You'll write this hundreds of times. Embrace it.

   result, err := doSomething()
   if err != nil {
       return fmt.Errorf("context about what we were doing: %w", err)
   }

 The pattern:
 1. Call a function that returns an error
 2. Check if err != nil immediately
 3. If it's an error, add context and return (or handle it)
 4. If not, proceed with the result

 The key discipline: ALWAYS add context when wrapping errors. Don't just
 return err — the caller needs to know WHAT failed, not just that
 something failed.

 BAD:  return err
 GOOD: return fmt.Errorf("loading config from %s: %w", path, err)
*/

func ReadConfig(path string) (string, error) {
	// Simulating a file read that might fail
	if path == "" {
		return "", errors.New("path cannot be empty")
	}
	if strings.HasSuffix(path, ".yaml") {
		return "host: localhost", nil
	}
	return "", fmt.Errorf("unsupported config format: %s", path)
}

func InitializeApp(configPath string) error {
	config, err := ReadConfig(configPath)
	if err != nil {
		// Add context: WHAT were we doing when this error occurred?
		return fmt.Errorf("initializing app with config %q: %w", configPath, err)
	}

	// Use config...
	_ = config
	return nil
}

// -------------------------------------------------------------------------
// Sentinel Errors
// -------------------------------------------------------------------------

/*
 Sentinel errors are package-level variables that represent specific error
 conditions. They let callers check for specific errors by identity:

   if errors.Is(err, ErrNotFound) { ... }

 Convention:
 - Name starts with "Err" (ErrNotFound, ErrUnauthorized, ErrTimeout)
 - Defined at package level using errors.New()
 - Don't create sentinels for every error — only for errors that callers
   need to distinguish and handle differently

 Standard library examples:
   io.EOF            — end of input (not really an "error")
   sql.ErrNoRows     — query returned no results
   os.ErrNotExist    — file doesn't exist
   context.Canceled  — operation was canceled
*/

// Sentinel errors for a user service
var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("invalid user id %d: %w", id, ErrInvalidInput)
	}
	if id == 999 {
		return "", fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	return "Alice", nil
}

func DemoSentinelErrors() {
	_, err := FindUser(999)
	if err != nil {
		// errors.Is checks the entire error chain
		if errors.Is(err, ErrNotFound) {
			fmt.Println("User not found — return 404")
		} else if errors.Is(err, ErrInvalidInput) {
			fmt.Println("Bad input — return 400")
		} else {
			fmt.Println("Unexpected error — return 500")
		}
	}
}

// -------------------------------------------------------------------------
// Custom Error Types
// -------------------------------------------------------------------------

/*
 When you need errors to carry structured data (not just a message),
 create a custom error type. This is a struct that implements the error
 interface.

 Use custom error types when callers need to:
 - Extract specific fields (status code, field name, etc.)
 - Make decisions based on error properties
 - Log structured error information

 Use errors.As to extract a custom error type from an error chain.
*/

// FieldError carries details about what validation failed.
// (Named FieldError here to avoid conflict with the exercises' ValidationError.)
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("validation error: field %q %s", e.Field, e.Message)
}

// RequestError represents an error with an HTTP status code.
// (Named RequestError here to avoid conflict with the exercises' HTTPError.)
type RequestError struct {
	StatusCode int
	Message    string
	Err        error // the underlying error
}

func (e *RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

// Unwrap lets errors.Is and errors.As see the underlying error.
func (e *RequestError) Unwrap() error {
	return e.Err
}

func DemoCustomErrors() {
	// Create a RequestError wrapping a FieldError
	valErr := &FieldError{Field: "email", Message: "must contain @"}
	apiErr := &RequestError{
		StatusCode: 400,
		Message:    "bad request",
		Err:        valErr,
	}

	fmt.Println(apiErr) // "[400] bad request: validation error: field "email" must contain @"

	// errors.As extracts the specific error type
	var ve *FieldError
	if errors.As(apiErr, &ve) {
		fmt.Printf("Validation failed on field: %s\n", ve.Field)
	}

	// Go 1.26 alternative — cleaner, no pointer variable needed:
	// if ve, ok := errors.AsType[*FieldError](apiErr); ok {
	//     fmt.Printf("Validation failed on field: %s\n", ve.Field)
	// }

	// errors.Is still works through the chain if needed
	// (In this case, valErr doesn't wrap anything, but the chain is RequestError -> FieldError)
}

// -------------------------------------------------------------------------
// Error Wrapping with %w and Unwrapping
// -------------------------------------------------------------------------

/*
 Error wrapping creates a chain of errors, each adding context:

   original error
       └── wrapped by: "loading config: <original>"
              └── wrapped by: "initializing app: loading config: <original>"

 You can inspect this chain with:
 - errors.Is(err, target) — is target anywhere in the chain?
 - errors.As(err, &target) — extract a specific type from the chain
 - errors.Unwrap(err) — get the next error in the chain (rarely needed)

 This is how you build informative error messages in layered applications:
 each layer adds its own context while preserving the original error.
*/

func connectDB(host string) error {
	return fmt.Errorf("connection refused to %s", host)
}

func loadUsers(dbHost string) error {
	err := connectDB(dbHost)
	if err != nil {
		return fmt.Errorf("loading users: %w", err)
	}
	return nil
}

func handleRequest() error {
	err := loadUsers("db.example.com")
	if err != nil {
		return fmt.Errorf("handling /api/users request: %w", err)
	}
	return nil
}

func DemoErrorWrapping() {
	err := handleRequest()
	if err != nil {
		// Full error message with complete context chain:
		fmt.Println(err)
		// "handling /api/users request: loading users: connection refused to db.example.com"
	}
}

// -------------------------------------------------------------------------
// Don't Panic — When panic IS Appropriate
// -------------------------------------------------------------------------

/*
 panic() causes the program to crash (unless recovered). Use it ONLY for:

 1. Truly unrecoverable situations (out of memory, corrupted state)
 2. Programmer bugs that should never happen in correct code
 3. Failed invariants during initialization (can't start the server)

 DO NOT panic for:
 - Bad user input (return an error)
 - Network failures (return an error)
 - File not found (return an error)
 - Any expected failure condition (return an error)

 The rule: if a caller could reasonably handle the failure, return an error.
 If the program is in an inconsistent state and can't continue, panic.

 Standard library examples of panic:
 - Index out of bounds (programmer bug)
 - nil pointer dereference (programmer bug)
 - regexp.MustCompile (fails only with invalid regex — a programmer bug)
*/

// MustParsePort panics if the port is invalid. The "Must" prefix is a
// convention that signals "this panics on failure."
func MustParsePort(port int) int {
	if port < 1 || port > 65535 {
		panic(fmt.Sprintf("invalid port number: %d", port))
	}
	return port
}

// -------------------------------------------------------------------------
// recover() — Catching Panics
// -------------------------------------------------------------------------

/*
 recover() catches a panic and prevents the program from crashing. It only
 works inside a deferred function. This is commonly used in:

 1. HTTP middleware (prevent one bad request from crashing the server)
 2. Worker goroutines (prevent one task from killing all workers)
 3. Plugin systems (isolate plugin failures)

 Pattern:
   defer func() {
       if r := recover(); r != nil {
           log.Printf("recovered from panic: %v", r)
       }
   }()
*/

// SafeExecute runs a function and recovers from any panic, returning it
// as an error instead.
func SafeExecute(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	fn()
	return nil
}

func DemoRecover() {
	err := SafeExecute(func() {
		panic("something terrible happened")
	})
	fmt.Println(err) // "recovered from panic: something terrible happened"

	// Program continues normally — the panic was caught
	err = SafeExecute(func() {
		fmt.Println("this runs fine")
	})
	fmt.Println(err) // nil
}

// -------------------------------------------------------------------------
// Error Handling in Web Services
// -------------------------------------------------------------------------

/*
 In web services, errors typically need to be translated into HTTP responses.
 A common pattern is a centralized error handler that maps error types to
 status codes:

   func handleError(w http.ResponseWriter, err error) {
       var apiErr *APIError
       if errors.As(err, &apiErr) {
           w.WriteHeader(apiErr.StatusCode)
           json.NewEncoder(w).Encode(apiErr)
           return
       }
       if errors.Is(err, ErrNotFound) {
           w.WriteHeader(404)
           return
       }
       // Default: internal server error
       w.WriteHeader(500)
   }

 COMMON MISTAKE: Not checking errors from Close(), Write(), etc.

   f, err := os.Open("file.txt")
   if err != nil { return err }
   defer f.Close()  // BUG: ignoring the error from Close()!

 While you can't return from a defer, at minimum log the error:
   defer func() {
       if err := f.Close(); err != nil {
           log.Printf("error closing file: %v", err)
       }
   }()

 GO 1.26: errors.AsType[T]()
 Go 1.26 adds errors.AsType[T](), a generic helper that makes type
 assertions on errors cleaner:

   if apiErr, ok := errors.AsType[*APIError](err); ok {
       // use apiErr directly — no separate variable declaration needed
   }

 This replaces the awkward two-step of declaring a typed variable then
 passing its address to errors.As. Prefer AsType in new code.
*/
