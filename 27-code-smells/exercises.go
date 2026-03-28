package codesmells

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

/*
=============================================================================
 EXERCISES: Code Smells
=============================================================================

 Each exercise presents SMELLY code (described in comments) and asks you to
 write the CLEAN version. The tests verify your refactored version produces
 correct results.

 Run the tests with:

   go test -v ./27-code-smells/

 Tip: Read the "smelly" version carefully, understand what it does, then
 write a cleaner version that produces the same results.

=============================================================================
*/

// =========================================================================
// Exercise 1: Refactor Over-Abstracted Interfaces
// =========================================================================
//
// SMELLY CODE (do NOT implement this — it's what you're fixing):
//
//   // Three interfaces for one implementation...
//   type DataReader interface {
//       ReadData(id string) (string, error)
//   }
//   type DataWriter interface {
//       WriteData(id string, value string) error
//   }
//   type DataDeleter interface {
//       DeleteData(id string) error
//   }
//   type DataProcessor interface {
//       DataReader
//       DataWriter
//       DataDeleter
//   }
//   // Only one implementation exists, and every consumer uses all methods.
//
// YOUR TASK: Create a simple KeyValueStore struct with these methods:
//   - Get(key string) (string, error) — returns error if key not found
//   - Set(key string, value string)   — stores the key-value pair
//   - Delete(key string) error        — returns error if key not found
//
// Use a map[string]string as the underlying storage. No interfaces needed
// since there's only one implementation.

// KeyValueStore is a simple in-memory key-value store.
// Replace the over-abstracted interface mess with a concrete struct.
type KeyValueStore struct {
	// YOUR CODE HERE
	data map[string]string
}

// NewKeyValueStore creates a ready-to-use KeyValueStore.
func NewKeyValueStore() *KeyValueStore {
	// YOUR CODE HERE
	return &KeyValueStore{}
}

// Get retrieves a value by key. Returns an error if the key doesn't exist.
func (s *KeyValueStore) Get(key string) (string, error) {
	// YOUR CODE HERE
	return "", fmt.Errorf("key %q not found", key)
}

// Set stores a key-value pair.
func (s *KeyValueStore) Set(key string, value string) {
	// YOUR CODE HERE
}

// Delete removes a key. Returns an error if the key doesn't exist.
func (s *KeyValueStore) Delete(key string) error {
	// YOUR CODE HERE
	return fmt.Errorf("key %q not found", key)
}

// =========================================================================
// Exercise 2: Fix Goroutine Leaks
// =========================================================================
//
// SMELLY CODE:
//
//   func LeakyProducer(data []int) <-chan int {
//       ch := make(chan int)
//       go func() {
//           for _, v := range data {
//               ch <- v    // blocks forever if consumer stops reading
//           }
//           close(ch)
//       }()
//       return ch
//   }
//
// The goroutine leaks if the consumer doesn't read all values (e.g., it
// only wants the first 3 items from a slice of 1000).
//
// YOUR TASK: Implement SafeProducer that accepts a context.Context.
// The goroutine should stop when the context is cancelled.

// SafeProducer sends values from data onto a channel, but stops if
// ctx is cancelled. The channel is closed when the goroutine exits.
func SafeProducer(ctx context.Context, data []int) <-chan int {
	// YOUR CODE HERE
	ch := make(chan int)
	go func() {
		defer close(ch)
		_ = ctx
		_ = data
		// Send each value, but check ctx.Done() to avoid leaking
	}()
	return ch
}

// =========================================================================
// Exercise 3: Fix Stuttering Names
// =========================================================================
//
// SMELLY CODE:
//
//   type EmailConfig struct {
//       EmailHost     string
//       EmailPort     int
//       EmailFrom     string
//       EmailPassword string
//   }
//
//   func NewEmailConfig(host string, port int, from, pass string) EmailConfig {
//       return EmailConfig{
//           EmailHost: host, EmailPort: port,
//           EmailFrom: from, EmailPassword: pass,
//       }
//   }
//
//   func (c EmailConfig) EmailAddress() string {
//       return fmt.Sprintf("%s:%d", c.EmailHost, c.EmailPort)
//   }
//
// Everything stutters: EmailConfig.EmailHost is read as "email config email host".
//
// YOUR TASK: Create a Config struct (not EmailConfig) with non-stuttering
// field names, plus a constructor and an Address() method.

// Config holds email configuration with clean, non-stuttering field names.
// Since this is in a hypothetical "email" package, the caller would write
// email.Config, email.Config.Host, etc.
type Config struct {
	// YOUR CODE HERE — use Host, Port, From, Password (no "Email" prefix)
	Host     string
	Port     int
	From     string
	Password string
}

// NewConfig creates an email Config.
func NewConfig(host string, port int, from, password string) Config {
	// YOUR CODE HERE
	return Config{}
}

// Address returns the host:port string for connecting to the mail server.
func (c Config) Address() string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 4: Refactor Global Mutable State to Dependency Injection
// =========================================================================
//
// SMELLY CODE:
//
//   var globalCache = map[string]string{}
//   var globalMu sync.RWMutex
//
//   func CacheGet(key string) (string, bool) {
//       globalMu.RLock()
//       defer globalMu.RUnlock()
//       v, ok := globalCache[key]
//       return v, ok
//   }
//
//   func CacheSet(key string, value string) {
//       globalMu.Lock()
//       defer globalMu.Unlock()
//       globalCache[key] = value
//   }
//
// YOUR TASK: Create a Cache struct with Get and Set methods.
// Use sync.RWMutex for thread safety. No global state.

// Cache is a thread-safe in-memory cache.
type Cache struct {
	// YOUR CODE HERE — need a mutex and a map
	mu   sync.RWMutex
	data map[string]string
}

// NewCache creates a ready-to-use Cache.
func NewCache() *Cache {
	// YOUR CODE HERE
	return &Cache{}
}

// Get retrieves a value from the cache. The bool indicates whether
// the key was found.
func (c *Cache) Get(key string) (string, bool) {
	// YOUR CODE HERE — use RLock for reads
	return "", false
}

// Set stores a key-value pair in the cache.
func (c *Cache) Set(key string, value string) {
	// YOUR CODE HERE — use Lock for writes
}

// =========================================================================
// Exercise 5: Fix Error Handling
// =========================================================================
//
// SMELLY CODE:
//
//   var errNotFound = errors.New("Not Found!")  // Bad formatting
//
//   func findItem(store map[string]string, key string) (string, error) {
//       val, ok := store[key]
//       if !ok {
//           return "", errNotFound
//       }
//       return val, nil
//   }
//
//   func processItem(store map[string]string, key string) (string, error) {
//       val, err := findItem(store, key)
//       if err != nil {
//           // Bad: string matching
//           if err.Error() == "Not Found!" {
//               return "default", nil
//           }
//           return "", err
//       }
//       return strings.ToUpper(val), nil
//   }
//
// YOUR TASK:
// 1. Define ErrNotFound as a proper sentinel error (lowercase, no punctuation)
// 2. Implement FindItem that returns ErrNotFound when key is missing
// 3. Implement ProcessItem that uses errors.Is to check for ErrNotFound

// ErrNotFound is a sentinel error for missing items.
// Follow Go conventions: lowercase, no punctuation.
var ErrNotFound = errors.New("placeholder")

// FindItem looks up a key in the store. Returns ErrNotFound if missing.
func FindItem(store map[string]string, key string) (string, error) {
	// YOUR CODE HERE
	return "", nil
}

// ProcessItem retrieves an item, returning "DEFAULT" if not found.
// Use errors.Is to check for ErrNotFound — never match on error strings.
func ProcessItem(store map[string]string, key string) (string, error) {
	// YOUR CODE HERE
	_ = errors.Is
	return "", nil
}

// =========================================================================
// Exercise 6: Simplify Over-Engineered Config
// =========================================================================
//
// SMELLY CODE (unnecessary factory pattern):
//
//   type ServerConfigBuilder struct {
//       host    string
//       port    int
//       timeout int
//   }
//   func NewServerConfigBuilder() *ServerConfigBuilder {
//       return &ServerConfigBuilder{host: "localhost", port: 8080, timeout: 30}
//   }
//   func (b *ServerConfigBuilder) WithHost(h string) *ServerConfigBuilder {
//       b.host = h; return b
//   }
//   func (b *ServerConfigBuilder) WithPort(p int) *ServerConfigBuilder {
//       b.port = p; return b
//   }
//   func (b *ServerConfigBuilder) Build() ServerConfig {
//       return ServerConfig{Host: b.host, Port: b.port, Timeout: b.timeout}
//   }
//
// A builder pattern for 3 fields is massive overkill. Go's zero values
// and struct literals handle this perfectly.
//
// YOUR TASK: Define ServerConfig as a simple struct with useful zero values.
// If Host is empty, DefaultHost should be used. If Port is 0, DefaultPort.
// If Timeout is 0, DefaultTimeout. Implement an Addr() method.

const (
	DefaultHost    = "localhost"
	DefaultPort    = 8080
	DefaultTimeout = 30
)

// ServerConfig holds server settings. Zero values should be meaningful —
// use the defaults above when fields aren't set.
type ServerConfig struct {
	Host    string
	Port    int
	Timeout int
}

// Addr returns "host:port", using defaults for zero values.
func (c ServerConfig) Addr() string {
	// YOUR CODE HERE
	// If Host is "", use DefaultHost. If Port is 0, use DefaultPort.
	return ""
}

// EffectiveTimeout returns the timeout, using DefaultTimeout for zero value.
func (c ServerConfig) EffectiveTimeout() int {
	// YOUR CODE HERE
	return 0
}

// =========================================================================
// Exercise 7: Fix Context Misuse
// =========================================================================
//
// SMELLY CODE:
//
//   type ctxKey string
//   func HandleRequest(ctx context.Context) string {
//       logger := ctx.Value(ctxKey("logger")).(Logger)
//       db := ctx.Value(ctxKey("db")).(Database)
//       user := db.GetUser(ctx, 1)
//       logger.Log("fetched user: " + user)
//       return user
//   }
//
// Storing services (logger, db) in context is an anti-pattern. Context
// should only carry request-scoped values (request ID, auth claims).
//
// YOUR TASK: Create a RequestHandler struct that takes dependencies via
// constructor, and uses context only for request-scoped data (request ID).

// Logger is a simple logging interface for the exercise.
type Logger interface {
	Log(msg string)
}

// Database is a simple database interface for the exercise.
type Database interface {
	GetUser(ctx context.Context, id int) string
}

// RequestHandler handles requests with properly injected dependencies.
type RequestHandler struct {
	// YOUR CODE HERE — store logger and db as struct fields
	logger Logger
	db     Database
}

// NewRequestHandler creates a handler with its dependencies.
func NewRequestHandler(logger Logger, db Database) *RequestHandler {
	// YOUR CODE HERE
	return nil
}

// requestIDKey is a typed context key for request IDs.
type requestIDKey struct{}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// GetRequestID extracts the request ID from context.
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// HandleRequest processes a request using injected dependencies and
// request-scoped data from context (request ID).
func (h *RequestHandler) HandleRequest(ctx context.Context, userID int) string {
	// YOUR CODE HERE
	// 1. Get the request ID from context (using GetRequestID)
	// 2. Use h.db to fetch the user
	// 3. Use h.logger to log the action (include request ID)
	// 4. Return the user string
	return ""
}

// =========================================================================
// Exercise 8: Fix a Function That Does Too Many Things
// =========================================================================
//
// SMELLY CODE:
//
//   func ProcessOrder(items []string, prices map[string]float64,
//       taxRate float64, discount float64) (string, error) {
//       // Calculate subtotal
//       subtotal := 0.0
//       for _, item := range items {
//           p, ok := prices[item]
//           if !ok {
//               return "", fmt.Errorf("unknown item: %s", item)
//           }
//           subtotal += p
//       }
//       // Apply discount
//       if discount > 0 {
//           subtotal = subtotal * (1 - discount)
//       }
//       // Calculate tax
//       tax := subtotal * taxRate
//       total := subtotal + tax
//       // Format receipt
//       receipt := fmt.Sprintf("Items: %d, Subtotal: $%.2f, Tax: $%.2f, Total: $%.2f",
//           len(items), subtotal, tax, total)
//       return receipt, nil
//   }
//
// This function calculates subtotal, applies discount, calculates tax,
// AND formats a receipt. Each of these should be its own function.
//
// YOUR TASK: Break this into focused functions. The test calls each one
// independently and also calls a ComposeReceipt that ties them together.

// CalculateSubtotal sums the prices of all items.
// Returns an error if any item isn't in the price list.
func CalculateSubtotal(items []string, prices map[string]float64) (float64, error) {
	// YOUR CODE HERE
	return 0, nil
}

// ApplyDiscount applies a discount rate (0.0 to 1.0) to an amount.
// A discount of 0.1 means 10% off. If discount is 0, returns amount unchanged.
func ApplyDiscount(amount float64, discount float64) float64 {
	// YOUR CODE HERE
	return 0
}

// CalculateTax computes the tax on an amount at the given rate.
func CalculateTax(amount float64, taxRate float64) float64 {
	// YOUR CODE HERE
	return 0
}

// FormatReceipt creates a formatted receipt string.
func FormatReceipt(itemCount int, subtotal, tax, total float64) string {
	// YOUR CODE HERE
	// Format: "Items: %d, Subtotal: $%.2f, Tax: $%.2f, Total: $%.2f"
	_ = fmt.Sprintf
	return ""
}

// ComposeReceipt ties together the focused functions above to process
// a complete order. This replaces the monolithic ProcessOrder function.
func ComposeReceipt(items []string, prices map[string]float64,
	taxRate float64, discount float64) (string, error) {
	// YOUR CODE HERE
	// 1. CalculateSubtotal
	// 2. ApplyDiscount
	// 3. CalculateTax
	// 4. FormatReceipt
	return "", nil
}
