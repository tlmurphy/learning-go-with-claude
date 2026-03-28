// Package dependencyinjection covers Go's approach to dependency injection:
// constructor injection, interface-based DI, functional options, and how
// to structure a service for testability and flexibility.
package dependencyinjection

import (
	"context"
	"fmt"
	"log"
	"time"
)

/*
=============================================================================
 DEPENDENCY INJECTION IN GO
=============================================================================

Dependency injection (DI) is a fancy term for a simple idea: instead of a
function or struct creating its own dependencies, you pass them in from
the outside.

  Bad:  func NewService() *Service { return &Service{db: sql.Open(...)} }
  Good: func NewService(db *sql.DB) *Service { return &Service{db: db} }

Why does this matter?
  1. TESTABILITY: You can pass in a mock database for tests
  2. FLEXIBILITY: You can swap implementations without changing the service
  3. CLARITY: Dependencies are visible in the constructor signature
  4. NO GLOBALS: Each instance has its own dependencies

Go doesn't need a DI framework. The language has everything built in:
  - Interfaces for abstraction
  - Struct fields for storage
  - Constructor functions for wiring
  - Functional options for configuration

=============================================================================
 CONSTRUCTOR INJECTION — THE GO WAY
=============================================================================

The most common DI pattern in Go: pass dependencies to a New function.

  type UserService struct {
      db     *sql.DB
      cache  *redis.Client
      logger *slog.Logger
  }

  func NewUserService(db *sql.DB, cache *redis.Client, logger *slog.Logger) *UserService {
      return &UserService{
          db:     db,
          cache:  cache,
          logger: logger,
      }
  }

The constructor makes dependencies explicit. Anyone reading this code
immediately knows what UserService needs to work.

=============================================================================
 INTERFACE-BASED DI — "ACCEPT INTERFACES, RETURN STRUCTS"
=============================================================================

This is Go's most powerful DI pattern. Define interfaces at the CONSUMER
site, not the provider site.

  // In the service package — define what WE need
  type UserRepository interface {
      GetByID(ctx context.Context, id int) (*User, error)
      Save(ctx context.Context, u *User) error
  }

  type UserService struct {
      repo UserRepository  // accepts the interface
  }

  func NewUserService(repo UserRepository) *UserService {
      return &UserService{repo: repo}
  }

Now you can pass in:
  - A PostgreSQL implementation for production
  - An in-memory implementation for tests
  - A caching wrapper that decorates the real implementation

The interface belongs to the CONSUMER (UserService), not the PRODUCER
(PostgresRepo). This is the opposite of Java, where interfaces are
defined next to their implementations.

Why return structs instead of interfaces?
  - Callers get the full API of the concrete type
  - No unnecessary abstraction
  - Types are documented more clearly
  - You can always assign a struct to an interface later

=============================================================================
 FUNCTIONAL OPTIONS PATTERN
=============================================================================

When you have optional configuration, the functional options pattern is
elegant:

  type Server struct {
      host    string
      port    int
      timeout time.Duration
      logger  *log.Logger
  }

  type Option func(*Server)

  func WithHost(host string) Option {
      return func(s *Server) { s.host = host }
  }

  func WithPort(port int) Option {
      return func(s *Server) { s.port = port }
  }

  func WithTimeout(d time.Duration) Option {
      return func(s *Server) { s.timeout = d }
  }

  func NewServer(opts ...Option) *Server {
      s := &Server{
          host:    "localhost",   // sensible defaults
          port:    8080,
          timeout: 30 * time.Second,
      }
      for _, opt := range opts {
          opt(s)
      }
      return s
  }

Usage:
  srv := NewServer(WithHost("0.0.0.0"), WithPort(9090))

Benefits:
  - Sensible defaults (zero config works)
  - Self-documenting option names
  - Backwards compatible (adding options doesn't break existing callers)
  - Options can do validation or complex setup

=============================================================================
 THE SERVICE LAYER PATTERN
=============================================================================

A common architecture for Go web services:

  Handler → Service → Repository
     ↓          ↓          ↓
  HTTP      Business    Database
  concern   logic       access

Each layer depends on interfaces from the layer below:

  // Repository layer (data access)
  type UserRepo interface {
      FindByID(ctx context.Context, id int) (*User, error)
      Save(ctx context.Context, u *User) error
  }

  // Service layer (business logic)
  type UserService struct { repo UserRepo }
  func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
      // Business logic: validation, caching, transformation
      return s.repo.FindByID(ctx, id)
  }

  // Handler layer (HTTP concerns)
  type UserHandler struct { svc *UserService }
  func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
      // Parse request, call service, write response
  }

=============================================================================
 PORTS AND ADAPTERS (HEXAGONAL ARCHITECTURE)
=============================================================================

A more formal version of the service layer:

  "Ports" = interfaces defining how the core communicates with the outside
  "Adapters" = implementations that connect to real infrastructure

  Core domain logic has NO infrastructure dependencies. It defines ports
  (interfaces) that adapters implement.

  // Port (defined by the core)
  type NotificationSender interface {
      Send(ctx context.Context, to, subject, body string) error
  }

  // Adapter (implements the port)
  type SMTPSender struct { ... }
  func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error { ... }

  // Another adapter (for tests)
  type FakeSender struct { Sent []Message }
  func (f *FakeSender) Send(ctx context.Context, to, subject, body string) error { ... }

=============================================================================
 DON'T USE DI FRAMEWORKS
=============================================================================

In Java/C#, DI frameworks (Spring, Guice) are common because wiring
is verbose. In Go, explicit wiring is the norm and is considered a feature.

  func main() {
      // Create dependencies bottom-up
      db := connectDB()
      cache := connectCache()
      logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

      // Wire the dependency graph
      userRepo := postgres.NewUserRepo(db)
      userCache := redis.NewUserCache(cache)
      userSvc := service.NewUserService(userRepo, userCache, logger)
      userHandler := handler.NewUserHandler(userSvc)

      // Set up routes
      mux := http.NewServeMux()
      mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
  }

This is explicit, readable, and debuggable. No magic, no reflection,
no runtime surprises.

Google's Wire tool (golang.org/x/tools/cmd/wire) generates this wiring
code at compile time from dependency declarations. It's useful for very
large projects but overkill for most Go services.

=============================================================================
 TESTING WITH DI
=============================================================================

DI makes testing easy because you can swap real implementations for fakes:

  func TestUserService_GetUser(t *testing.T) {
      // Fake repository for testing
      repo := &FakeUserRepo{
          users: map[int]*User{
              1: {ID: 1, Name: "Alice"},
          },
      }

      svc := NewUserService(repo)
      user, err := svc.GetUser(context.Background(), 1)
      if err != nil {
          t.Fatal(err)
      }
      if user.Name != "Alice" {
          t.Errorf("got %q, want %q", user.Name, "Alice")
      }
  }

No mocking library needed. Just implement the interface with a struct
that has the behavior you want. Go interfaces are implicitly satisfied,
so any struct with the right methods works.

=============================================================================
 WHERE DEPENDENCIES GET CREATED
=============================================================================

The "composition root" is where all dependencies are created and wired.
In Go, this is main() or a setup function called from main().

  func main() {
      // Parse config
      cfg := loadConfig()

      // Create leaf dependencies (no deps of their own)
      db := mustConnectDB(cfg.DatabaseURL)
      logger := newLogger(cfg.LogLevel)

      // Create mid-level dependencies
      userRepo := postgres.NewUserRepo(db)
      emailSvc := smtp.NewEmailService(cfg.SMTPHost)

      // Create top-level dependencies
      userSvc := service.NewUserService(userRepo, emailSvc, logger)

      // Create HTTP handlers
      handler := api.NewHandler(userSvc)

      // Start server
      log.Fatal(http.ListenAndServe(":8080", handler))
  }

Dependencies flow DOWN: main creates everything, passes deps to constructors.
No package reaches up to grab what it needs from globals.

=============================================================================
*/

// --- Demo types and functions ---

// Notifier is an interface for sending notifications — a "port" in
// hexagonal architecture terms.
type Notifier interface {
	Notify(ctx context.Context, recipient, message string) error
}

// EmailNotifier is a concrete "adapter" that sends email notifications.
type EmailNotifier struct {
	SMTPHost string
	From     string
}

// Notify sends an email notification.
func (e *EmailNotifier) Notify(ctx context.Context, recipient, message string) error {
	// In production, this would actually send an email
	fmt.Printf("Sending email from %s via %s to %s: %s\n",
		e.From, e.SMTPHost, recipient, message)
	return nil
}

// SlackNotifier is another adapter — sends to Slack instead of email.
type SlackNotifier struct {
	WebhookURL string
}

// Notify sends a Slack notification.
func (s *SlackNotifier) Notify(ctx context.Context, recipient, message string) error {
	fmt.Printf("Sending Slack message to %s: %s\n", recipient, message)
	return nil
}

// AlertService uses dependency injection — it doesn't know or care
// whether it's sending email, Slack messages, or something else.
type AlertService struct {
	notifier Notifier
}

// NewAlertService demonstrates constructor injection.
func NewAlertService(n Notifier) *AlertService {
	return &AlertService{notifier: n}
}

// SendAlert sends an alert using whatever notifier was injected.
func (a *AlertService) SendAlert(ctx context.Context, recipient, msg string) error {
	return a.notifier.Notify(ctx, recipient, fmt.Sprintf("ALERT: %s", msg))
}

// DemoConstructorInjection shows how different implementations can be
// swapped by injecting different dependencies.
func DemoConstructorInjection() {
	fmt.Println("=== Constructor Injection ===")

	// Production: use email
	emailNotifier := &EmailNotifier{
		SMTPHost: "smtp.company.com",
		From:     "alerts@company.com",
	}
	alertSvc := NewAlertService(emailNotifier)
	_ = alertSvc.SendAlert(context.Background(), "oncall@company.com", "Server CPU > 90%")

	// Alternative: use Slack
	slackNotifier := &SlackNotifier{WebhookURL: "https://hooks.slack.com/..."}
	alertSvc2 := NewAlertService(slackNotifier)
	_ = alertSvc2.SendAlert(context.Background(), "#alerts", "Server CPU > 90%")
}

// --- Functional Options Demo ---

// HTTPServer demonstrates the functional options pattern.
type HTTPServer struct {
	host    string
	port    int
	timeout time.Duration
	logger  *log.Logger
}

// ServerOption configures an HTTPServer.
type ServerOption func(*HTTPServer)

// WithHost sets the server's listen host.
func WithHost(host string) ServerOption {
	return func(s *HTTPServer) {
		s.host = host
	}
}

// WithPort sets the server's listen port.
func WithPort(port int) ServerOption {
	return func(s *HTTPServer) {
		s.port = port
	}
}

// WithTimeout sets the server's request timeout.
func WithTimeout(d time.Duration) ServerOption {
	return func(s *HTTPServer) {
		s.timeout = d
	}
}

// WithLogger sets the server's logger.
func WithLogger(l *log.Logger) ServerOption {
	return func(s *HTTPServer) {
		s.logger = l
	}
}

// NewHTTPServer creates a server with sensible defaults, overridden by options.
func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	s := &HTTPServer{
		host:    "localhost",
		port:    8080,
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Addr returns the listen address.
func (s *HTTPServer) Addr() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// DemoFunctionalOptions shows the functional options pattern in action.
func DemoFunctionalOptions() {
	fmt.Println("=== Functional Options ===")

	// Zero config — all defaults
	s1 := NewHTTPServer()
	fmt.Println("Default:", s1.Addr())

	// Custom config
	s2 := NewHTTPServer(
		WithHost("0.0.0.0"),
		WithPort(9090),
		WithTimeout(60*time.Second),
	)
	fmt.Println("Custom:", s2.Addr())
}

// --- Service Layer Demo ---

// User represents a user in the system.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserRepository defines the data access interface (port).
type UserRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	Save(ctx context.Context, u *User) error
}

// UserService contains business logic and depends on the repository interface.
type UserService struct {
	repo   UserRepository
	notify Notifier
}

// NewUserService demonstrates a service with multiple injected dependencies.
func NewUserService(repo UserRepository, notify Notifier) *UserService {
	return &UserService{
		repo:   repo,
		notify: notify,
	}
}

// GetUser retrieves a user by ID with business logic.
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}
	return user, nil
}

// UpdateUser updates a user and sends a notification.
func (s *UserService) UpdateUser(ctx context.Context, u *User) error {
	if err := s.repo.Save(ctx, u); err != nil {
		return fmt.Errorf("saving user %d: %w", u.ID, err)
	}
	if err := s.notify.Notify(ctx, u.Email, "Your profile was updated"); err != nil {
		// Log but don't fail — notification is best-effort
		fmt.Printf("notification failed for user %d: %v\n", u.ID, err)
	}
	return nil
}
