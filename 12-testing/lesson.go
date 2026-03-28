package testingmod

/*
Module 12: Testing
===================

Go has testing built into the language toolchain — not as a third-party
framework, but as a first-class citizen. "go test" is as fundamental as
"go build". This philosophy shows: Go's testing approach is simpler and
more pragmatic than most languages.

Test File Naming
-----------------
Test files end in _test.go. The Go toolchain automatically excludes them
from production builds. They live alongside the code they test:

  user.go           ← Production code
  user_test.go      ← Tests for user.go

This convention is enforced by the toolchain, not just a suggestion.

Test Function Signatures
------------------------
Test functions follow a strict naming pattern:

  func TestXxx(t *testing.T) { ... }

  - Must start with "Test" followed by an uppercase letter
  - Takes exactly one parameter: *testing.T
  - No return value — signal failure with t.Error/t.Fatal

The testing.T type provides:
  t.Error(args...)     — Log failure and continue
  t.Errorf(fmt, ...)   — Formatted failure, continue
  t.Fatal(args...)     — Log failure and STOP this test immediately
  t.Fatalf(fmt, ...)   — Formatted failure, stop
  t.Log(args...)       — Log info (only shown with -v flag)
  t.Skip(args...)      — Skip this test
  t.Helper()           — Mark as helper (fixes line numbers in output)

Use t.Error when you want to report multiple failures.
Use t.Fatal when continuing would panic or produce meaningless results.

Table-Driven Tests: THE Go Pattern
-----------------------------------
This is the dominant testing pattern in Go. Instead of writing separate test
functions for each case, you define a table of inputs and expected outputs:

  func TestAdd(t *testing.T) {
      tests := []struct {
          name string
          a, b int
          want int
      }{
          {"positives", 2, 3, 5},
          {"zero", 0, 5, 5},
          {"negatives", -1, -1, -2},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              got := Add(tt.a, tt.b)
              if got != tt.want {
                  t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
              }
          })
      }
  }

Why this works so well:
  - Adding a test case is just adding a line to the table
  - Each case has a name (visible in test output)
  - t.Run creates subtests that can be run individually
  - The error message tells you exactly what failed

t.Run Subtests
--------------
t.Run creates named subtests. This gives you:
  - Hierarchical test output
  - The ability to run specific subtests: go test -run TestAdd/negatives
  - Separate test contexts (each subtest gets its own t)

t.Helper()
-----------
When you write helper functions that call t.Error or t.Fatal, mark them
with t.Helper(). This makes test output report the CALLER's line number
instead of the helper's line number:

  func assertEq(t *testing.T, got, want int) {
      t.Helper()  // Without this, failures point to THIS line
      if got != want {
          t.Errorf("got %d, want %d", got, want)
      }
  }

t.Parallel()
-------------
Call t.Parallel() at the start of a test or subtest to run it concurrently
with other parallel tests:

  func TestSomething(t *testing.T) {
      t.Parallel() // This test runs in parallel with other parallel tests
      // ...
  }

Gotcha: When using t.Parallel() with table-driven tests, you MUST capture
the loop variable:

  for _, tt := range tests {
      tt := tt  // Capture! (Go 1.22+ fixes this, but be explicit)
      t.Run(tt.name, func(t *testing.T) {
          t.Parallel()
          // use tt safely
      })
  }

Benchmarks
-----------
Benchmark functions measure performance:

  func BenchmarkConcat(b *testing.B) {
      for i := 0; i < b.N; i++ {
          _ = concat("hello", "world")
      }
  }

Run with: go test -bench=. -benchmem
The framework automatically determines b.N for stable measurements.

Example Functions
------------------
Example functions serve as both documentation and tests:

  func ExampleReverse() {
      fmt.Println(Reverse("hello"))
      // Output: olleh
  }

If the output comment doesn't match, the test fails. These examples appear
in godoc.

testdata/ Directory
--------------------
Test fixtures go in a testdata/ directory. The Go toolchain ignores this
directory during builds. Use it for:
  - Golden files (expected output)
  - Sample input data
  - Configuration files for tests

Coverage
---------
  go test -cover ./...                    — Show coverage percentage
  go test -coverprofile=coverage.out      — Generate coverage data
  go tool cover -html=coverage.out        — View in browser

Aim for meaningful coverage, not 100%. Test the tricky parts, not getters.

Mocking with Interfaces
------------------------
Go doesn't need a mocking framework. Instead:
  1. Define dependencies as interfaces
  2. Pass real implementations in production
  3. Pass mock implementations in tests

  type UserStore interface {
      GetUser(id int) (*User, error)
  }

  // In production: realStore := &PostgresUserStore{db}
  // In tests:      mockStore := &MockUserStore{users: testData}

Golden File Testing
--------------------
Compare output against a saved "golden" file:
  1. Run the code, get actual output
  2. Compare against testdata/expected_output.golden
  3. To update golden files: go test -update (with a custom flag)

This is great for testing serialization, formatting, or code generation.

The -race Flag
---------------
ALWAYS run tests with the race detector in CI:

  go test -race ./...

It detects data races — concurrent access to shared memory where at least
one access is a write. A test passing without -race proves nothing about
thread safety.
*/

import "fmt"

// ==========================================
// Functions to be Tested (used by exercises)
// ==========================================

// Reverse reverses a string. Used as a simple example for testing.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string reads the same forwards and backwards.
func IsPalindrome(s string) bool {
	return s == Reverse(s)
}

// ConcatStrings concatenates strings using different approaches.
// Used for benchmarking exercises.
func ConcatWithPlus(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

// ConcatWithBuilder concatenates strings efficiently using a builder.
func ConcatWithBuilder(strs []string) string {
	var b []byte
	for _, s := range strs {
		b = append(b, s...)
	}
	return string(b)
}

// Abs returns the absolute value of an integer.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Clamp restricts a value to a range [min, max].
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ==========================================
// Types for Mocking Exercises
// ==========================================

// UserRecord represents a user in the system.
type UserRecord struct {
	ID    int
	Name  string
	Email string
}

// UserStore defines the interface for user persistence.
// In production, this might be backed by PostgreSQL.
// In tests, you'll create a mock implementation.
type UserStore interface {
	GetUser(id int) (*UserRecord, error)
	ListUsers() ([]*UserRecord, error)
	CreateUser(name, email string) (*UserRecord, error)
}

// UserService contains business logic that depends on UserStore.
type UserService struct {
	store UserStore
}

// NewUserService creates a UserService with the given store.
func NewUserService(store UserStore) *UserService {
	return &UserService{store: store}
}

// GetUserDisplayName returns a formatted display name for a user.
// Returns an error if the user is not found.
func (s *UserService) GetUserDisplayName(id int) (string, error) {
	user, err := s.store.GetUser(id)
	if err != nil {
		return "", fmt.Errorf("getting user %d: %w", id, err)
	}
	return fmt.Sprintf("%s <%s>", user.Name, user.Email), nil
}

// ListUserNames returns the names of all users.
func (s *UserService) ListUserNames() ([]string, error) {
	users, err := s.store.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	names := make([]string, len(users))
	for i, u := range users {
		names[i] = u.Name
	}
	return names, nil
}

// CreateAndGreet creates a user and returns a greeting message.
func (s *UserService) CreateAndGreet(name, email string) (string, error) {
	user, err := s.store.CreateUser(name, email)
	if err != nil {
		return "", fmt.Errorf("creating user: %w", err)
	}
	return fmt.Sprintf("Welcome, %s! Your ID is %d.", user.Name, user.ID), nil
}

// ==========================================
// Golden File Testing Support
// ==========================================

// FormatReport generates a formatted report from data.
// This is used for golden file testing exercises.
func FormatReport(title string, items []string) string {
	result := fmt.Sprintf("=== %s ===\n", title)
	for i, item := range items {
		result += fmt.Sprintf("  %d. %s\n", i+1, item)
	}
	result += fmt.Sprintf("Total: %d items\n", len(items))
	return result
}

// ==========================================
// Function for Parallel Test Exercises
// ==========================================

// SlowHash simulates a slow hashing operation (for parallel test exercises).
// In real life, this might be bcrypt or scrypt.
func SlowHash(input string) string {
	// Simulate computation
	hash := 0
	for _, ch := range input {
		hash = hash*31 + int(ch)
	}
	if hash < 0 {
		hash = -hash
	}
	return fmt.Sprintf("%x", hash)
}
