package packages

/*
Module 09: Packages and Modules
================================

Go's package system is one of its greatest strengths for building maintainable
software. Unlike languages where you import individual files, Go organizes code
into packages — directories of related .go files that compile together as a unit.

Why Packages Matter
-------------------
Packages solve several fundamental problems:
  - Namespacing: Avoid name collisions between different parts of your code
  - Encapsulation: Control what's visible outside a package (the API surface)
  - Reusability: Share code between projects via modules
  - Compilation speed: Go compiles packages independently and caches results

Package Naming Conventions
--------------------------
Go has strong opinions about package names, and following them makes your code
feel idiomatic:

  - Use short, lowercase, single-word names: "http", "json", "bytes"
  - Avoid underscores, hyphens, or mixedCaps in package names
  - The package name should describe what the package provides, not what it contains
  - Avoid generic names like "util", "common", "misc" — they're a code smell
  - The package name is part of every call site: http.Get(), not httppackage.Get()

Bad:  package string_utils   →  Good: package strings
Bad:  package myHelpers      →  Good: package auth
Bad:  package common         →  Good: package validation

Exported vs Unexported (Capitalization Rule)
--------------------------------------------
This is Go's visibility system — simple but powerful:

  - Exported (public):   Starts with an uppercase letter — User, ParseJSON, ErrNotFound
  - Unexported (private): Starts with a lowercase letter — validate, dbConn, errInternal

This applies to EVERYTHING: types, functions, methods, variables, constants, fields.

Why not keywords like "public" and "private"? Because Go's approach is:
  1. Visible at a glance — you can tell from the name whether it's exported
  2. Enforced by the compiler — not just a convention you can accidentally break
  3. Package-level, not class-level — everything in the same package can see
     everything else in that package, regardless of which file it's in

The go.mod File and Module System
---------------------------------
A module is a collection of packages released together. The go.mod file at the
root of your project defines your module:

  module github.com/yourname/yourproject    // Module path (also the import path)
  go 1.24                                    // Minimum Go version
  require github.com/some/dependency v1.2.3  // Dependencies

Key commands:
  go mod init <module-path>  — Create a new module
  go mod tidy                — Add missing deps, remove unused ones
  go get <package>@<version> — Add or update a dependency
  go mod vendor              — Copy dependencies into vendor/ directory
  go mod verify              — Verify dependencies haven't been tampered with

Tip: "go mod tidy" is your best friend. Run it often. It keeps your go.mod and
go.sum clean and accurate.

Internal Packages
-----------------
Go has a special directory name: "internal". Code inside an internal/ directory
can only be imported by code in the parent of that internal/ directory.

  myproject/
    cmd/
      server/          ← can import internal/auth
    internal/
      auth/            ← only importable within myproject/
    pkg/
      validation/      ← importable by anyone

This is enforced by the Go toolchain — not a convention, a rule. It's how you
share code within your project without exposing it to the world.

Package Initialization (init Functions)
---------------------------------------
Each package can have init() functions that run automatically when the package
is loaded:

  func init() {
      // Runs before main(), used for setup
  }

A single file can have multiple init() functions. They run in the order they
appear. Across files in a package, the order is alphabetical by filename (but
don't rely on this — if order matters, you're probably doing it wrong).

init() is useful for:
  - Registering database drivers
  - Setting up package-level state
  - Validating configuration

init() should NOT be used for:
  - Complex logic that could fail (it can't return errors)
  - Anything with side effects that make testing hard
  - Heavy computation (it blocks program startup)

A common pattern is the "blank import" that triggers init() for side effects:
  import _ "github.com/lib/pq"  // Registers the PostgreSQL driver

Project Layout for Web Services
-------------------------------
The Go community has converged on a standard project layout:

  myservice/
    cmd/
      myservice/         ← main.go lives here (the entry point)
      migrate/           ← another binary in the same project
    internal/
      server/            ← HTTP handlers, middleware
      database/          ← Database access layer
      auth/              ← Authentication logic
    pkg/                 ← Code you're OK with others importing (use sparingly)
    api/                 ← OpenAPI specs, proto files
    configs/             ← Configuration files
    go.mod
    go.sum

Why cmd/? Because a single module can produce multiple binaries. Each
subdirectory of cmd/ is a separate main package.

Why internal/? Because most of your code shouldn't be importable by others.
If you're building a web service (not a library), almost everything goes in
internal/.

Build Tags and Conditional Compilation
---------------------------------------
Build tags let you include or exclude files from compilation based on
conditions like OS, architecture, or custom tags.

Modern syntax (Go 1.17+):
  //go:build linux
  //go:build !windows
  //go:build integration

Legacy syntax (still works):
  // +build linux

You'll see build tags used for:
  - Platform-specific code (syscalls, file paths)
  - Integration tests that need external services
  - Debug/release builds with different behavior

Circular Dependencies
---------------------
Go does NOT allow circular imports. If package A imports B, then B cannot
import A (directly or indirectly). This is a feature, not a limitation:
  - It forces you to think about dependency direction
  - It keeps compilation fast (no need to resolve cycles)
  - It often reveals design problems early

If you hit a circular dependency, common solutions are:
  1. Extract the shared types into a third package
  2. Use interfaces to invert the dependency
  3. Merge the packages if they're too tightly coupled
*/

import "fmt"

// ==========================================
// Exported vs Unexported Identifiers
// ==========================================

// User is an exported type — visible outside this package.
// In a real web service, this might be your domain model.
type User struct {
	// ID is exported — other packages can read and write it.
	ID   int
	Name string

	// passwordHash is unexported — only code within the 'packages' package
	// can access this field. This is encapsulation at work.
	passwordHash string
}

// NewUser is an exported constructor function.
// This pattern is common in Go: since struct fields can be unexported,
// you need a constructor to set them.
func NewUser(id int, name, password string) User {
	return User{
		ID:           id,
		Name:         name,
		passwordHash: hashPassword(password),
	}
}

// hashPassword is unexported — it's an implementation detail.
// External packages don't need to know HOW we hash passwords.
func hashPassword(password string) string {
	// In real code, use bcrypt. This is just for demonstration.
	return "hashed:" + password
}

// CheckPassword is exported — it's part of the public API.
// Notice how it accesses the unexported passwordHash field.
// Within the same package, everything is visible.
func (u User) CheckPassword(password string) bool {
	return u.passwordHash == hashPassword(password)
}

// ==========================================
// Package-Level Variables and Constants
// ==========================================

// MaxUsers is an exported constant — part of the package's public API.
const MaxUsers = 1000

// defaultRegion is unexported — an internal default.
const defaultRegion = "us-east-1"

// Version is an exported package-level variable.
// In web services, this is often set at build time using ldflags:
//
//	go build -ldflags "-X packages.Version=1.2.3"
var Version = "dev"

// registry is an unexported package-level variable.
// Only code within this package can access it.
var registry = make(map[string]User)

// ==========================================
// Demonstrating Package Organization
// ==========================================

// Repository defines an exported interface with unexported implementation.
// This is a powerful pattern: consumers depend on the interface, not the
// concrete type. You can swap implementations without changing callers.
type Repository interface {
	FindByID(id int) (User, error)
	Save(user User) error
}

// memoryRepository is unexported — you can't create one directly from
// outside this package. You MUST use NewMemoryRepository().
type memoryRepository struct {
	users map[int]User
}

// NewMemoryRepository is the exported constructor.
// It returns the Repository interface, not *memoryRepository.
// This hides the implementation details completely.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		users: make(map[int]User),
	}
}

func (r *memoryRepository) FindByID(id int) (User, error) {
	user, ok := r.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %d not found", id)
	}
	return user, nil
}

func (r *memoryRepository) Save(user User) error {
	if user.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", user.ID)
	}
	r.users[user.ID] = user
	return nil
}

// ==========================================
// Init Functions and Initialization Order
// ==========================================

// initOrder tracks the order in which init functions run.
// This is for demonstration only — don't rely on init ordering in real code.
var initOrder []string

// InitOrder returns a copy of the initialization order log.
// We return a copy to prevent external code from modifying our internal state.
func InitOrder() []string {
	result := make([]string, len(initOrder))
	copy(result, initOrder)
	return result
}

func init() {
	initOrder = append(initOrder, "first init in lesson.go")
}

func init() {
	initOrder = append(initOrder, "second init in lesson.go")
}

// ==========================================
// Demonstrating API Design Patterns
// ==========================================

// Option is a functional option type. This pattern is common in well-designed
// Go packages — it lets you provide a clean API with sensible defaults while
// still allowing customization.
type Option func(*config)

// config is unexported — the configuration details are hidden.
type config struct {
	port    int
	host    string
	timeout int // seconds
}

// WithPort sets the server port. Returns an Option for use with NewServer.
func WithPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

// WithHost sets the server host. Returns an Option for use with NewServer.
func WithHost(host string) Option {
	return func(c *config) {
		c.host = host
	}
}

// WithTimeout sets the server timeout in seconds.
func WithTimeout(seconds int) Option {
	return func(c *config) {
		c.timeout = seconds
	}
}

// Server represents a configured server. Exported type, but the config
// details are encapsulated.
type Server struct {
	cfg config
}

// NewServer creates a server with functional options.
// Usage: NewServer(WithPort(8080), WithHost("localhost"))
func NewServer(opts ...Option) *Server {
	// Start with sensible defaults
	cfg := config{
		port:    8080,
		host:    "0.0.0.0",
		timeout: 30,
	}

	// Apply each option
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Server{cfg: cfg}
}

// Addr returns the server's address string.
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.cfg.host, s.cfg.port)
}

// Timeout returns the server's timeout in seconds.
func (s *Server) Timeout() int {
	return s.cfg.timeout
}

// ==========================================
// Build Tags (Conceptual Demonstration)
// ==========================================

// BuildInfo holds build metadata. In a real project, you might have
// platform-specific files that set this differently:
//
//	build_linux.go   — //go:build linux
//	build_darwin.go  — //go:build darwin
//	build_windows.go — //go:build windows
//
// Each file would set Platform to the appropriate value.
type BuildInfo struct {
	Platform string
	Version  string
	Debug    bool
}

// DefaultBuildInfo returns build info for the current compilation.
func DefaultBuildInfo() BuildInfo {
	return BuildInfo{
		Platform: "portable", // Would be OS-specific with build tags
		Version:  Version,
		Debug:    false,
	}
}

// ==========================================
// Preventing Circular Dependencies
// ==========================================

// When you have two packages that need to know about each other,
// use interfaces to break the cycle.
//
// BAD (circular):
//   package auth imports package user (to get User type)
//   package user imports package auth (to check permissions)
//
// GOOD (interface breaks the cycle):
//   package auth defines Authenticatable interface
//   package user has User that satisfies Authenticatable
//   package auth never imports package user
//
// Or extract shared types:
//   package models defines User (shared types, no business logic)
//   package auth imports models
//   package user imports models
//   No cycle!

// Authenticatable demonstrates using an interface to avoid circular deps.
// Instead of importing a specific User type from another package,
// we define what we need right here.
type Authenticatable interface {
	GetID() int
	GetPasswordHash() string
}

// Authenticate works with any type that satisfies Authenticatable.
// It doesn't need to import the package where the concrete type lives.
func Authenticate(a Authenticatable, password string) bool {
	return a.GetPasswordHash() == hashPassword(password)
}

// Make User satisfy Authenticatable.
func (u User) GetID() int              { return u.ID }
func (u User) GetPasswordHash() string { return u.passwordHash }
