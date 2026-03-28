package packages

/*
Exercises: Packages and Modules
================================

These exercises focus on the concepts that make Go's package system powerful:
visibility rules, initialization, API design, and proper encapsulation.

Since we can't create multiple packages in a single exercise file, some
exercises simulate multi-package scenarios within one package.
*/

// =============================================================================
// Exercise 1: Exported vs Unexported
// =============================================================================

// BankAccount represents a bank account with proper encapsulation.
// The balance should NOT be directly accessible from outside this package —
// all modifications must go through methods that enforce business rules.
//
// Implement:
//   - An exported field: AccountNumber (string)
//   - An exported field: Owner (string)
//   - An unexported field: balance (float64)
//   - An unexported field: isActive (bool)
type BankAccount struct {
	// YOUR CODE HERE
	AccountNumber string
	Owner         string
	balance       float64
	isActive      bool
}

// NewBankAccount creates a new active bank account with the given details
// and an initial balance. Returns an error if initialBalance is negative.
func NewBankAccount(accountNumber, owner string, initialBalance float64) (BankAccount, error) {
	// YOUR CODE HERE
	return BankAccount{}, nil
}

// Deposit adds money to the account. Returns an error if:
//   - The account is not active
//   - The amount is negative or zero
func (a *BankAccount) Deposit(amount float64) error {
	// YOUR CODE HERE
	return nil
}

// Withdraw removes money from the account. Returns an error if:
//   - The account is not active
//   - The amount is negative or zero
//   - The amount exceeds the current balance (no overdrafts)
func (a *BankAccount) Withdraw(amount float64) error {
	// YOUR CODE HERE
	return nil
}

// Balance returns the current balance. This is the only way to read the
// balance from outside the package — a "getter" that provides read-only access.
func (a *BankAccount) Balance() float64 {
	// YOUR CODE HERE
	return 0
}

// Deactivate marks the account as inactive. Once deactivated,
// deposits and withdrawals should fail.
func (a *BankAccount) Deactivate() {
	// YOUR CODE HERE
}

// IsActive returns whether the account is currently active.
func (a *BankAccount) IsActive() bool {
	// YOUR CODE HERE
	return false
}

// =============================================================================
// Exercise 2: Package-Level Variables and Init Ordering
// =============================================================================

// registeredPlugins simulates a plugin registry that gets populated
// during package initialization.
var registeredPlugins []string

// RegisterPlugin adds a plugin name to the registry.
// This simulates what might happen in init() functions across multiple files.
func RegisterPlugin(name string) {
	// YOUR CODE HERE
}

// ListPlugins returns a copy of the registered plugin names.
// Important: return a COPY, not the original slice. Why? Because if you
// return the original, callers could modify your internal state.
func ListPlugins() []string {
	// YOUR CODE HERE
	return nil
}

// ClearPlugins resets the plugin registry. Useful for testing.
func ClearPlugins() {
	// YOUR CODE HERE
}

// =============================================================================
// Exercise 3: Interface-Based API Design
// =============================================================================

// Logger defines a logging interface. In a real application, you'd have
// different implementations: console logger, file logger, structured logger, etc.
// The key insight: consumers depend on this interface, not on any concrete type.
type Logger interface {
	Info(msg string)
	Error(msg string)
	Messages() []string
}

// logEntry is unexported — an implementation detail.
type logEntry struct {
	level   string
	message string
}

// memoryLogger is an unexported implementation of Logger.
// It stores log messages in memory — useful for testing.
type memoryLogger struct {
	entries []logEntry
}

// NewMemoryLogger creates a Logger that stores messages in memory.
// Notice: it returns the Logger interface, not *memoryLogger.
// This is proper API design — hide the implementation.
func NewMemoryLogger() Logger {
	// YOUR CODE HERE
	return nil
}

// Implement the Logger interface methods on memoryLogger:

// Info logs a message at INFO level.
// Store it as a logEntry with level "INFO".
func (l *memoryLogger) Info(msg string) {
	// YOUR CODE HERE
}

// Error logs a message at ERROR level.
// Store it as a logEntry with level "ERROR".
func (l *memoryLogger) Error(msg string) {
	// YOUR CODE HERE
}

// Messages returns all logged messages formatted as "[LEVEL] message".
// For example: "[INFO] server started", "[ERROR] connection failed"
func (l *memoryLogger) Messages() []string {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 4: Functional Options Pattern
// =============================================================================

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host           string
	Port           int
	Database       string
	MaxConnections int
	TimeoutSeconds int
	SSLEnabled     bool
}

// DBOption is a functional option for configuring DatabaseConfig.
type DBOption func(*DatabaseConfig)

// WithDBHost returns a DBOption that sets the database host.
func WithDBHost(host string) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithDBPort returns a DBOption that sets the database port.
func WithDBPort(port int) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithDBName returns a DBOption that sets the database name.
func WithDBName(name string) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithMaxConnections returns a DBOption that sets the max connection pool size.
func WithMaxConnections(n int) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithDBTimeout returns a DBOption that sets the connection timeout in seconds.
func WithDBTimeout(seconds int) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithSSL returns a DBOption that enables or disables SSL.
func WithSSL(enabled bool) DBOption {
	// YOUR CODE HERE
	return nil
}

// NewDatabaseConfig creates a DatabaseConfig with sensible defaults,
// then applies the provided options.
//
// Defaults:
//   - Host: "localhost"
//   - Port: 5432
//   - Database: "app"
//   - MaxConnections: 10
//   - TimeoutSeconds: 30
//   - SSLEnabled: false
func NewDatabaseConfig(opts ...DBOption) DatabaseConfig {
	// YOUR CODE HERE
	return DatabaseConfig{}
}

// =============================================================================
// Exercise 5: Preventing Circular Dependencies with Interfaces
// =============================================================================

// This exercise simulates a common scenario: an OrderService needs to
// check user permissions, and a UserService might need to look up order history.
// Without careful design, this creates a circular dependency.
//
// Solution: Define interfaces for what each service needs from the other.

// PermissionChecker is what OrderService needs — it doesn't need
// the full UserService, just the ability to check permissions.
type PermissionChecker interface {
	HasPermission(userID int, permission string) bool
}

// OrderLookup is what UserService needs — it doesn't need
// the full OrderService, just the ability to look up orders.
type OrderLookup interface {
	OrderCountForUser(userID int) int
}

// Order represents a simple order.
type Order struct {
	ID     int
	UserID int
	Amount float64
}

// OrderService manages orders and depends on PermissionChecker (not UserService directly).
type OrderService struct {
	orders      []Order
	permissions PermissionChecker
}

// NewOrderService creates an OrderService with the given permission checker.
func NewOrderService(pc PermissionChecker) *OrderService {
	// YOUR CODE HERE
	return nil
}

// PlaceOrder adds an order if the user has the "place_order" permission.
// Returns an error if the user lacks permission.
func (s *OrderService) PlaceOrder(userID int, amount float64) error {
	// YOUR CODE HERE
	return nil
}

// OrderCountForUser returns the number of orders for a given user.
// This satisfies the OrderLookup interface.
func (s *OrderService) OrderCountForUser(userID int) int {
	// YOUR CODE HERE
	return 0
}

// UserService manages users and depends on OrderLookup (not OrderService directly).
type UserService struct {
	permissions map[int]map[string]bool // userID -> permission -> granted
	orders      OrderLookup
}

// NewUserService creates a UserService with the given order lookup.
func NewUserService(ol OrderLookup) *UserService {
	// YOUR CODE HERE
	return nil
}

// GrantPermission gives a user a specific permission.
func (s *UserService) GrantPermission(userID int, permission string) {
	// YOUR CODE HERE
}

// HasPermission checks if a user has a specific permission.
// This satisfies the PermissionChecker interface.
func (s *UserService) HasPermission(userID int, permission string) bool {
	// YOUR CODE HERE
	return false
}

// UserOrderCount returns how many orders a user has, using the OrderLookup.
func (s *UserService) UserOrderCount(userID int) int {
	// YOUR CODE HERE
	return 0
}

// =============================================================================
// Exercise 6: Godoc-Style Package API
// =============================================================================

// EmailValidator validates email addresses. It demonstrates proper godoc
// commenting style and clean API design.
//
// Usage:
//
//	v := NewEmailValidator(WithBlockedDomains("example.com", "test.com"))
//	err := v.Validate("user@example.com")
//	if err != nil {
//	    // handle invalid email
//	}
type EmailValidator struct {
	blockedDomains map[string]bool
	maxLength      int
}

// EmailOption is a functional option for configuring EmailValidator.
type EmailOption func(*EmailValidator)

// WithBlockedDomains returns an EmailOption that blocks emails from
// the specified domains.
func WithBlockedDomains(domains ...string) EmailOption {
	// YOUR CODE HERE
	return nil
}

// WithMaxLength returns an EmailOption that sets the maximum allowed
// email length. Emails exceeding this length will fail validation.
func WithMaxLength(length int) EmailOption {
	// YOUR CODE HERE
	return nil
}

// NewEmailValidator creates an EmailValidator with the given options.
//
// Defaults:
//   - No blocked domains
//   - MaxLength: 254 (per RFC 5321)
func NewEmailValidator(opts ...EmailOption) *EmailValidator {
	// YOUR CODE HERE
	return nil
}

// Validate checks if the given email address is valid.
// An email is valid if:
//   - It is not empty
//   - It contains exactly one "@" symbol
//   - It has a non-empty local part (before @) and domain part (after @)
//   - The domain contains at least one "." (dot)
//   - It does not exceed the maximum length
//   - Its domain is not in the blocked domains list
//
// Returns nil if valid, or an error describing why it's invalid.
func (v *EmailValidator) Validate(email string) error {
	// YOUR CODE HERE
	return nil
}
