// Package codesmells catalogs common Go anti-patterns and demonstrates how to
// recognize and refactor them into idiomatic, maintainable code.
package codesmells

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

/*
=============================================================================
 CODE SMELLS IN GO
=============================================================================

"Code smell" is a term from Martin Fowler's refactoring work. It means code
that technically works but hints at deeper problems — poor design, future
maintenance nightmares, or misunderstanding of the language's idioms.

Go has its own set of smells. The language is opinionated, and fighting its
conventions makes code harder to read for everyone. The Go community has
strong norms, and code that ignores them stands out like a Java class
hierarchy in a Go codebase.

This lesson covers the most common Go-specific code smells. Each exercise
gives you SMELLY code to refactor into something clean.

=============================================================================
 SMELL 1: INTERFACE POLLUTION
=============================================================================

The #1 Go smell, especially from developers coming from Java or C#.

Bad: Defining interfaces before you need them.
  type UserRepository interface {
      GetUser(id int) (*User, error)
      CreateUser(u *User) error
      UpdateUser(u *User) error
      DeleteUser(id int) error
  }
  // ...with exactly one implementation.

Good: Define interfaces where they're CONSUMED, not where they're produced.
  // In the handler package:
  type UserGetter interface {
      GetUser(id int) (*User, error)
  }
  // Only declares what THIS package needs.

The Go proverb: "Accept interfaces, return structs."

Interfaces should emerge from usage, not be designed upfront. If you have
an interface with one implementation, you probably don't need the interface
(unless you're writing tests, which is a legitimate reason).

=============================================================================
 SMELL 2: PACKAGE-LEVEL MUTABLE STATE
=============================================================================

Global variables are the enemy of testability and concurrency safety.

Bad:
  var db *sql.DB  // package-level, set during init
  func GetUser(id int) (*User, error) { return db.Query(...) }

Good:
  type UserStore struct { db *sql.DB }
  func (s *UserStore) GetUser(id int) (*User, error) { return s.db.Query(...) }

Package-level mutable state makes it impossible to:
  - Run tests in parallel (they share the global)
  - Test with different configurations
  - Reason about what a function depends on

Constants are fine at package level. Mutable state is not.

=============================================================================
 SMELL 3: STUTTERING NAMES
=============================================================================

Go uses the package name as a qualifier, so names shouldn't repeat it.

Bad:
  package user
  type UserService struct{} // caller writes: user.UserService
  type UserConfig struct{}  // caller writes: user.UserConfig

Good:
  package user
  type Service struct{}     // caller writes: user.Service
  type Config struct{}      // caller writes: user.Config

This extends to functions too:
  Bad:  user.NewUserService()
  Good: user.NewService()
  Good: user.New()  // if there's only one main type

=============================================================================
 SMELL 4: NOT USING ZERO VALUES
=============================================================================

Go types have meaningful zero values. Don't write constructors when the
zero value works fine.

Bad:
  func NewBuffer() *Buffer {
      return &Buffer{
          data:  make([]byte, 0),
          count: 0,
          ready: false,
      }
  }

Good:
  // Just use Buffer{}. All those fields have useful zero values already.
  var b Buffer
  b.Write(data)

The standard library is full of this: bytes.Buffer, strings.Builder,
sync.Mutex, sync.WaitGroup — none need constructors.

=============================================================================
 SMELL 5: IGNORING ERRORS
=============================================================================

This one gets people in production. The most common ignored errors:

  f.Close()              // Bad! Close can flush buffers and fail
  fmt.Fprintf(w, ...)    // Bad in HTTP handlers! Write errors mean the
                         //   client disconnected
  json.NewEncoder(w).Encode(v)  // Bad! Encode can fail

At minimum, log the error. For Close, a common pattern:

  defer func() {
      if err := f.Close(); err != nil {
          log.Printf("failed to close file: %v", err)
      }
  }()

For writes in HTTP handlers:
  if _, err := w.Write(data); err != nil {
      // Client disconnected. Log and return — don't try to write an error
      // response to a broken connection.
      return
  }

=============================================================================
 SMELL 6: CONTEXT MISUSE
=============================================================================

context.Context is for request-scoped values, deadlines, and cancellation.
It is NOT a dependency injection container.

Bad:
  ctx = context.WithValue(ctx, "db", database)
  ctx = context.WithValue(ctx, "logger", logger)
  // Now pull them out everywhere...
  db := ctx.Value("db").(*sql.DB)

Good:
  type Handler struct {
      db     *sql.DB
      logger *slog.Logger
  }

Context values are appropriate for:
  - Request ID / trace ID
  - Authentication claims
  - Deadline/cancellation
That's about it. If you're putting services in context, use struct fields.

=============================================================================
 SMELL 7: GOROUTINE LEAKS
=============================================================================

Starting goroutines without cleanup paths is one of the most common
production bugs in Go.

Bad:
  func Process(ch <-chan int) {
      go func() {
          for v := range ch {
              // process v
          }
      }()
      // Who closes ch? When does this goroutine stop?
  }

Good:
  func Process(ctx context.Context, ch <-chan int) {
      go func() {
          for {
              select {
              case v, ok := <-ch:
                  if !ok { return }
                  // process v
              case <-ctx.Done():
                  return
              }
          }
      }()
  }

Every goroutine needs an exit strategy:
  1. The channel it reads from is closed
  2. A context is cancelled
  3. A done channel is signaled
  4. The work naturally completes

=============================================================================
 SMELL 8: NAKED RETURNS IN LONG FUNCTIONS
=============================================================================

Naked returns (returning without specifying values) are fine in short
functions. In long functions, they make it impossible to see what's returned.

Bad:
  func process(data []byte) (result string, err error) {
      // 50 lines of code...
      if something {
          result = "processed"
          return  // What is err here? Did we set it somewhere above?
      }
      // 30 more lines...
      return  // What gets returned? I have to trace through all 80 lines.
  }

Good: Use explicit returns in functions longer than ~10 lines.

=============================================================================
 SMELL 9: GIANT INTERFACES
=============================================================================

Bad:
  type DataStore interface {
      GetUser(id int) (*User, error)
      CreateUser(u *User) error
      UpdateUser(u *User) error
      DeleteUser(id int) error
      GetProduct(id int) (*Product, error)
      CreateProduct(p *Product) error
      ListOrders(userID int) ([]Order, error)
      // ...15 more methods
  }

Good: Small, focused interfaces.
  type UserGetter interface { GetUser(id int) (*User, error) }
  type UserCreator interface { CreateUser(u *User) error }

The standard library's best interfaces have 1-2 methods:
  io.Reader, io.Writer, io.Closer, fmt.Stringer, sort.Interface (3 methods)

"The bigger the interface, the weaker the abstraction." — Rob Pike

=============================================================================
 SMELL 10: INIT FUNCTION ABUSE
=============================================================================

init() functions run at package load time. They can't be tested, can't
accept parameters, and their order within a package is based on source
file name (fragile).

Bad:
  func init() {
      db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
      if err != nil {
          log.Fatal(err)  // crashes at startup, can't test
      }
      globalDB = db
  }

Good: Explicit initialization in main().
  func main() {
      db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
      if err != nil {
          log.Fatal(err)
      }
      svc := NewService(db)
      // ...
  }

Legitimate uses of init(): registering drivers/codecs (like database/sql
drivers, image format decoders). That's about it.

=============================================================================
 SMELL 11: ERROR STRING FORMATTING
=============================================================================

Go convention: error strings should NOT be capitalized and should NOT
end with punctuation.

Bad:
  return fmt.Errorf("Failed to connect to database.")
  return errors.New("User not found!")

Good:
  return fmt.Errorf("failed to connect to database")
  return errors.New("user not found")

Why? Because errors are often wrapped:
  fmt.Errorf("processing request: %w", err)
  // "processing request: Failed to connect to database." looks wrong
  // "processing request: failed to connect to database" reads naturally

=============================================================================
 SMELL 12: STRING MATCHING INSTEAD OF errors.Is/errors.As
=============================================================================

Bad:
  if err.Error() == "not found" { ... }
  if strings.Contains(err.Error(), "timeout") { ... }

Good:
  if errors.Is(err, ErrNotFound) { ... }
  var netErr *net.OpError
  if errors.As(err, &netErr) { ... }

String matching breaks when error messages change and doesn't work with
wrapped errors. Sentinel errors and error types are robust.

=============================================================================
 SMELL 13: MAP OF BOOLEANS AS SETS
=============================================================================

Bad:
  seen := map[string]bool{}
  seen["alice"] = true
  if seen["bob"] { ... }

Better:
  seen := map[string]struct{}{}
  seen["alice"] = struct{}{}
  if _, ok := seen["bob"]; ok { ... }

The struct{} version uses zero bytes per value. The bool version uses 1 byte.
For small maps it doesn't matter, but for millions of entries it adds up.
More importantly, struct{} signals intent: "this is a set."

That said, the bool version IS more readable and simpler. In most code,
readability wins. Use struct{} when memory actually matters.

=============================================================================
 SMELL 14: CHANNEL OVERUSE
=============================================================================

Channels are great for communicating between goroutines. But sometimes a
mutex is simpler and faster.

Bad (overengineered):
  type Counter struct {
      ch chan func()
  }
  func (c *Counter) Increment() {
      c.ch <- func() { doWork() }
  }

Good (simple mutex):
  type Counter struct {
      mu    sync.Mutex
      count int
  }
  func (c *Counter) Increment() {
      c.mu.Lock()
      defer c.mu.Unlock()
      c.count++
  }

Use channels for: pipelines, fan-out/fan-in, signaling (done channels).
Use mutexes for: protecting shared state with simple access patterns.

=============================================================================
*/

// --- Demo types and functions for the lesson ---

// DemoInterfacePollution shows the bad vs good pattern for interfaces.
func DemoInterfacePollution() {
	fmt.Println("=== Interface Pollution ===")
	fmt.Println("Bad: Define huge interface at the implementation site")
	fmt.Println("Good: Define small interfaces at the consumer site")
	fmt.Println()

	// Good: consumer defines only what it needs
	var getter UserGetter = &InMemoryUserStore{
		users: map[int]string{1: "Alice", 2: "Bob"},
	}
	name, err := getter.GetUser(1)
	if err == nil {
		fmt.Println("Got user:", name)
	}
}

// UserGetter is a small, focused interface — defined where it's consumed.
type UserGetter interface {
	GetUser(id int) (string, error)
}

// InMemoryUserStore is a concrete implementation.
type InMemoryUserStore struct {
	users map[int]string
}

func (s *InMemoryUserStore) GetUser(id int) (string, error) {
	name, ok := s.users[id]
	if !ok {
		return "", fmt.Errorf("user %d not found", id)
	}
	return name, nil
}

// DemoGoroutineLeak shows a goroutine that leaks vs one with proper cleanup.
func DemoGoroutineLeak() {
	fmt.Println("=== Goroutine Leak Demo ===")

	// BAD: This goroutine will leak if nobody reads from ch
	ch := make(chan int)
	go func() {
		ch <- 42 // blocks forever if nobody reads
	}()
	// We read it here so it doesn't actually leak in this demo
	v := <-ch
	fmt.Println("Got value:", v)

	// GOOD: Use context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan int, 1)
	go func() {
		select {
		case resultCh <- 42:
		case <-ctx.Done():
			return // clean exit
		}
	}()
	fmt.Println("Got value:", <-resultCh)
	cancel() // clean up
}

// DemoErrorFormatting shows proper Go error formatting conventions.
func DemoErrorFormatting() {
	fmt.Println("=== Error Formatting ===")

	// Bad: capitalized, punctuated
	bad := errors.New("Failed to connect to database.")
	fmt.Println("Bad error:", bad)

	// Good: lowercase, no punctuation
	good := errors.New("failed to connect to database")
	fmt.Println("Good error:", good)

	// Wrapping reads naturally when lowercase
	wrapped := fmt.Errorf("processing user 42: %w", good)
	fmt.Println("Wrapped:", wrapped)
}

// DemoContextMisuse shows the bad and good patterns for context values.
func DemoContextMisuse() {
	fmt.Println("=== Context Values ===")
	fmt.Println("Good: request ID in context")

	type contextKey string
	const requestIDKey contextKey = "requestID"

	ctx := context.WithValue(context.Background(), requestIDKey, "req-abc-123")
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Println("Request ID:", id)
	}

	fmt.Println("Bad: putting services in context. Use struct fields instead.")
}

// DemoSets shows map[string]bool vs map[string]struct{} for sets.
func DemoSets() {
	fmt.Println("=== Sets in Go ===")

	// Simple but slightly wasteful
	seen1 := map[string]bool{
		"alice": true,
		"bob":   true,
	}
	if seen1["alice"] {
		fmt.Println("bool map: alice is in the set")
	}

	// Memory efficient, signals "this is a set"
	seen2 := map[string]struct{}{
		"alice": {},
		"bob":   {},
	}
	if _, ok := seen2["alice"]; ok {
		fmt.Println("struct{} map: alice is in the set")
	}
}

// DemoMutexVsChannel shows when a mutex is simpler than a channel.
func DemoMutexVsChannel() {
	fmt.Println("=== Mutex vs Channel ===")

	// Simple shared state → use mutex
	var mu sync.Mutex
	count := 0

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("Mutex counter:", count)
}

// Ensure imports are used.
var _ io.Reader
var _ = strings.Builder{}
