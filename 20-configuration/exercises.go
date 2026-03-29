package configuration

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 EXERCISES: Configuration
=============================================================================

 Work through these exercises in order. Each one builds on concepts from
 the lesson. Run the tests with:

   make test 20

 Tip: Run a single test at a time while working:

   go test -v -run TestEnvConfig ./20-configuration/

=============================================================================
*/

// =========================================================================
// Exercise 1: Config from Environment Variables
// =========================================================================

// AppConfig represents the configuration for an application.
// Load each field from the corresponding environment variable with
// the given default.
type AppConfig struct {
	Host         string        // env: APP_HOST, default: "0.0.0.0"
	Port         int           // env: APP_PORT, default: 8080
	DatabaseURL  string        // env: DATABASE_URL, default: "postgres://localhost:5432/app"
	LogLevel     string        // env: LOG_LEVEL, default: "info"
	Debug        bool          // env: APP_DEBUG, default: false
	ReadTimeout  time.Duration // env: READ_TIMEOUT, default: 15s
	WriteTimeout time.Duration // env: WRITE_TIMEOUT, default: 15s
}

// LoadAppConfig loads configuration from environment variables with defaults.
// Use os.LookupEnv (or the helpers from the lesson) to read env vars.
//
// Environment variable mapping:
//
//	APP_HOST      → Host      (string, default "0.0.0.0")
//	APP_PORT      → Port      (int, default 8080)
//	DATABASE_URL  → DatabaseURL (string, default "postgres://localhost:5432/app")
//	LOG_LEVEL     → LogLevel  (string, default "info")
//	APP_DEBUG     → Debug     (bool, default false)
//	READ_TIMEOUT  → ReadTimeout (duration, default 15s)
//	WRITE_TIMEOUT → WriteTimeout (duration, default 15s)
func LoadAppConfig() AppConfig {
	// YOUR CODE HERE
	return AppConfig{}
}

// =========================================================================
// Exercise 2: Functional Options Pattern
// =========================================================================

// ServiceConfig holds configuration for an HTTP service.
type ServiceConfig struct {
	Name           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxBodySize    int64
	EnableCORS     bool
	AllowedOrigins []string
}

// ServiceOption is a function that modifies ServiceConfig.
type ServiceOption func(*ServiceConfig)

// NewServiceConfig creates a ServiceConfig with defaults and applies options.
//
// Defaults:
//   - Name: "service"
//   - Port: 8080
//   - ReadTimeout: 15s
//   - WriteTimeout: 15s
//   - MaxBodySize: 1MB (1048576 bytes)
//   - EnableCORS: false
//   - AllowedOrigins: empty slice
func NewServiceConfig(opts ...ServiceOption) ServiceConfig {
	// YOUR CODE HERE
	return ServiceConfig{}
}

// WithName returns a ServiceOption that sets the service name.
func WithName(name string) ServiceOption {
	// YOUR CODE HERE
	return nil
}

// WithServicePort returns a ServiceOption that sets the port.
func WithServicePort(port int) ServiceOption {
	// YOUR CODE HERE
	return nil
}

// WithTimeouts returns a ServiceOption that sets both read and write timeouts.
func WithTimeouts(read, write time.Duration) ServiceOption {
	// YOUR CODE HERE
	return nil
}

// WithMaxBodySize returns a ServiceOption that sets the max body size.
func WithMaxBodySize(size int64) ServiceOption {
	// YOUR CODE HERE
	return nil
}

// WithCORS returns a ServiceOption that enables CORS with the given origins.
func WithCORS(origins ...string) ServiceOption {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 3: Config Validator
// =========================================================================

// ConfigErrors holds a list of validation error messages.
type ConfigErrors struct {
	errors []string
}

// Add adds an error message.
func (e *ConfigErrors) Add(msg string) {
	e.errors = append(e.errors, msg)
}

// HasErrors returns true if there are any errors.
func (e *ConfigErrors) HasErrors() bool {
	return len(e.errors) > 0
}

// Error implements the error interface.
func (e *ConfigErrors) Error() string {
	return strings.Join(e.errors, "; ")
}

// Messages returns all error messages.
func (e *ConfigErrors) Messages() []string {
	return e.errors
}

// ValidateAppConfig validates an AppConfig and returns all errors found.
// Return nil if the config is valid.
//
// Validation rules:
//   - Host must not be empty
//   - Port must be between 1 and 65535
//   - DatabaseURL must not be empty
//   - DatabaseURL must be a valid URL (use url.Parse — invalid if err != nil or scheme is empty)
//   - LogLevel must be one of: "debug", "info", "warn", "error"
//   - ReadTimeout must be positive
//   - WriteTimeout must be positive
func ValidateAppConfig(cfg AppConfig) error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 4: Config Layering
// =========================================================================

// LayeredConfig represents config that can come from multiple sources.
// Each field has a value and a source indicating where it came from.
type LayeredConfig struct {
	Host    ConfigValue[string]
	Port    ConfigValue[int]
	Debug   ConfigValue[bool]
	Timeout ConfigValue[time.Duration]
}

// ConfigValue holds a value and its source.
type ConfigValue[T comparable] struct {
	Value  T
	Source string // "default", "env", "flag"
}

// LoadLayeredConfig builds a config by applying layers in priority order:
// 1. defaults (lowest priority)
// 2. env vars (medium priority)
// 3. flags (highest priority)
//
// Parameters:
//   - envVars: map of environment variable values (simulated, not real env vars)
//     Keys: "APP_HOST", "APP_PORT", "APP_DEBUG", "APP_TIMEOUT"
//   - flags: map of CLI flag values (simulated)
//     Keys: "host", "port", "debug", "timeout"
//
// Defaults:
//   - Host: "localhost"
//   - Port: 8080
//   - Debug: false
//   - Timeout: 30s
//
// Only override a value if the map contains the key (even if empty string).
// For port: parse from string to int, skip if invalid.
// For debug: "true"/"1"/"yes" = true, "false"/"0"/"no" = false.
// For timeout: parse as duration (e.g., "5s"), skip if invalid.
func LoadLayeredConfig(envVars map[string]string, flags map[string]string) LayeredConfig {
	// YOUR CODE HERE
	return LayeredConfig{}
}

// =========================================================================
// Exercise 5: Feature Flag System
// =========================================================================

// FeatureFlags manages a set of feature flags with typed getters.
type FeatureFlags struct {
	flags map[string]interface{}
	mu    sync.RWMutex
}

// NewFeatureFlags creates a new FeatureFlags instance with the given
// initial flags.
func NewFeatureFlags(initial map[string]interface{}) *FeatureFlags {
	// YOUR CODE HERE
	return nil
}

// IsEnabled returns whether a boolean feature flag is enabled.
// Returns defaultVal if the flag doesn't exist or isn't a bool.
func (f *FeatureFlags) IsEnabled(name string, defaultVal bool) bool {
	// YOUR CODE HERE
	return false
}

// GetString returns a string feature flag value.
// Returns defaultVal if the flag doesn't exist or isn't a string.
func (f *FeatureFlags) GetString(name string, defaultVal string) string {
	// YOUR CODE HERE
	return ""
}

// GetInt returns an int feature flag value.
// Returns defaultVal if the flag doesn't exist or isn't an int.
func (f *FeatureFlags) GetInt(name string, defaultVal int) int {
	// YOUR CODE HERE
	return 0
}

// Set sets a feature flag value. It's safe for concurrent use.
func (f *FeatureFlags) Set(name string, value interface{}) {
	// YOUR CODE HERE
}

// All returns a copy of all current flag values.
func (f *FeatureFlags) All() map[string]interface{} {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 6: Config Watcher
// =========================================================================

// ConfigWatcher watches for configuration changes and notifies listeners.
// It simulates a config watcher (in production, this would watch a file
// or poll a config service).
type ConfigWatcher struct {
	current   map[string]string
	listeners []chan ConfigChange
	mu        sync.RWMutex
}

// ConfigChange represents a change to a configuration value.
type ConfigChange struct {
	Key      string
	OldValue string
	NewValue string
}

// NewConfigWatcher creates a new watcher with initial configuration.
func NewConfigWatcher(initial map[string]string) *ConfigWatcher {
	// YOUR CODE HERE
	return nil
}

// Get returns the current value for a key, or empty string if not set.
func (w *ConfigWatcher) Get(key string) string {
	// YOUR CODE HERE
	return ""
}

// Watch returns a channel that receives ConfigChange notifications.
// The channel is buffered with size 10 to prevent blocking.
func (w *ConfigWatcher) Watch() <-chan ConfigChange {
	// YOUR CODE HERE
	return nil
}

// Update sets a new value for a key. If the value has changed, notify
// all watchers. If the value hasn't changed, do nothing.
func (w *ConfigWatcher) Update(key, value string) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 7: Database Config Parser
// =========================================================================

// DatabaseConfig represents database connection configuration.
type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
	Options  map[string]string // additional connection options
}

// ParseDatabaseURL parses a database connection string URL into a
// DatabaseConfig.
//
// Format: postgres://user:password@host:port/database?sslmode=disable&key=value
//
// If port is not specified, default to 5432.
// If sslmode is not specified, default to "disable".
// Any other query parameters go into Options.
//
// Return an error if the URL is not parseable.
func ParseDatabaseURL(rawURL string) (DatabaseConfig, error) {
	// YOUR CODE HERE
	_ = url.Parse // hint
	return DatabaseConfig{}, nil
}

// BuildDatabaseURL converts a DatabaseConfig back into a connection URL string.
// Format: postgres://user:password@host:port/database?sslmode=MODE
// Include any Options as additional query parameters (sorted alphabetically).
func (c DatabaseConfig) BuildDatabaseURL() string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 8: Complete Server Configuration
// =========================================================================

// FullServerConfig represents a complete server configuration combining
// all the patterns from this module.
type FullServerConfig struct {
	Server   ServerSection
	Database DatabaseSection
	Logging  LoggingSection
	Features map[string]bool
}

type ServerSection struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseSection struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type LoggingSection struct {
	Level  string
	Format string // "json" or "text"
}

// DefaultFullServerConfig returns a FullServerConfig with production-ready
// defaults.
//
// Server defaults:
//   - Host: "0.0.0.0"
//   - Port: 8080
//   - ReadTimeout: 15s
//   - WriteTimeout: 15s
//   - IdleTimeout: 60s
//
// Database defaults:
//   - URL: "postgres://localhost:5432/app"
//   - MaxOpenConns: 25
//   - MaxIdleConns: 5
//   - ConnMaxLifetime: 5 minutes
//
// Logging defaults:
//   - Level: "info"
//   - Format: "json"
//
// Features: empty map (not nil)
func DefaultFullServerConfig() FullServerConfig {
	// YOUR CODE HERE
	return FullServerConfig{}
}

// ValidateFullServerConfig validates the complete config.
// Return a *ConfigErrors with all validation errors, or nil if valid.
//
// Rules:
//   - Server.Port must be 1-65535
//   - Server.ReadTimeout must be positive
//   - Server.WriteTimeout must be positive
//   - Database.URL must not be empty
//   - Database.MaxOpenConns must be > 0
//   - Logging.Level must be one of: "debug", "info", "warn", "error"
//   - Logging.Format must be one of: "json", "text"
func ValidateFullServerConfig(cfg FullServerConfig) error {
	// YOUR CODE HERE
	return nil
}

// String returns a human-readable representation of the config.
// IMPORTANT: Do NOT include sensitive values (like database passwords).
// Mask the database URL password if present.
//
// Format example:
//
//	Server: 0.0.0.0:8080 (read=15s, write=15s, idle=60s)
//	Database: postgres://user:****@localhost:5432/app (pool: 25/5, lifetime=5m0s)
//	Logging: level=info, format=json
//	Features: feature1=true, feature2=false
func (c FullServerConfig) String() string {
	// YOUR CODE HERE
	return ""
}

// These suppress "unused import" errors in stubs.
var _ = fmt.Sprintf
var _ = os.LookupEnv
var _ = strconv.Atoi
var _ = strings.Join
var _ = time.Second
var _ = url.Parse
