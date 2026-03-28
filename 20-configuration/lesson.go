package configuration

/*
=============================================================================
 Module 20: Configuration
=============================================================================

 Configuration is one of those things that seems simple until it isn't.
 Every production service needs to handle config, and getting it wrong
 leads to outages, security incidents, or deploy nightmares.

 The 12-Factor App methodology (https://12factor.net) recommends storing
 config in environment variables. This has become the standard approach
 for cloud-native applications because:
   - Environment variables work everywhere (Docker, Kubernetes, CI/CD)
   - They naturally separate config from code
   - They don't get accidentally committed to version control
   - They can be changed without rebuilding the binary

 But environment variables alone aren't enough for a real application.
 You typically need:
   1. Sensible defaults (the app works out of the box)
   2. Config file support (for complex local development setups)
   3. Environment variable overrides (for deployment)
   4. CLI flag overrides (for debugging and one-off runs)
   5. Validation (fail fast on bad config)

 This module covers all of these patterns.

 WHY THIS MATTERS FOR WEB SERVICES:
 - Wrong port? Your service doesn't start.
 - Wrong DB connection string? Silent data loss or outage.
 - Secret in source code? Security incident.
 - No validation? A typo in a config value causes a 3 AM page.

=============================================================================
*/

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// -------------------------------------------------------------------------
// Environment Variables: os.Getenv and os.LookupEnv
// -------------------------------------------------------------------------

/*
 os.Getenv(key) returns the value or "" if not set.
 os.LookupEnv(key) returns (value, true) if set, or ("", false) if not.

 The difference matters: Getenv can't distinguish between "not set" and
 "set to empty string". LookupEnv can.

   os.Getenv("MISSING")     → ""
   os.Getenv("EMPTY")       → ""    (if EMPTY="")
   os.LookupEnv("MISSING")  → "", false
   os.LookupEnv("EMPTY")    → "", true
*/

// GetEnvOrDefault returns the environment variable value, or the default
// if the variable is not set.
func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// GetEnvAsInt returns the environment variable as an int, or the default
// if not set or not a valid integer.
func GetEnvAsInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvAsBool returns the environment variable as a bool, or the default.
// Truthy values: "true", "1", "yes" (case-insensitive)
// Everything else (including not set) returns the default.
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			return true
		case "false", "0", "no":
			return false
		}
	}
	return defaultValue
}

// GetEnvAsDuration returns the environment variable as a time.Duration,
// or the default if not set or not parseable.
// Accepts formats like "5s", "2m30s", "1h".
func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

// DemoEnvHelpers shows how to use the environment variable helpers.
func DemoEnvHelpers() {
	host := GetEnvOrDefault("APP_HOST", "localhost")
	port := GetEnvAsInt("APP_PORT", 8080)
	debug := GetEnvAsBool("APP_DEBUG", false)
	timeout := GetEnvAsDuration("APP_TIMEOUT", 30*time.Second)

	fmt.Printf("Host: %s, Port: %d, Debug: %t, Timeout: %v\n",
		host, port, debug, timeout)
}

// -------------------------------------------------------------------------
// Struct-Based Configuration
// -------------------------------------------------------------------------

/*
 The standard pattern in Go: define a struct for your config, then
 populate it from various sources. This gives you:

   1. Type safety — port is an int, not a string
   2. Documentation — struct fields and comments describe the config
   3. Defaults — set in a constructor or Load function
   4. Validation — check the struct after loading
   5. Testing — pass config structs to functions, easy to mock

 The struct becomes the single source of truth for all configuration.
*/

// ServerConfig represents the configuration for an HTTP server.
type ServerConfig struct {
	Host         string        `env:"SERVER_HOST" default:"localhost"`
	Port         int           `env:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" default:"60s"`
	Debug        bool          `env:"SERVER_DEBUG" default:"false"`
}

// LoadServerConfig loads server configuration from environment variables
// with sensible defaults.
func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Host:         GetEnvOrDefault("SERVER_HOST", "localhost"),
		Port:         GetEnvAsInt("SERVER_PORT", 8080),
		ReadTimeout:  GetEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: GetEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  GetEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		Debug:        GetEnvAsBool("SERVER_DEBUG", false),
	}
}

// Addr returns the host:port address string.
func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// -------------------------------------------------------------------------
// The flag Package for CLI Arguments
// -------------------------------------------------------------------------

/*
 The standard library's flag package provides simple CLI argument parsing.
 It's adequate for most services but lacks subcommands (use cobra or
 urfave/cli for complex CLIs).

 Typical pattern:
   port := flag.Int("port", 8080, "server port")
   flag.Parse()  // must call before using flag values
   fmt.Println(*port)

 Flags integrate well with config layering:
   default → config file → env var → CLI flag

 Each layer overrides the previous, with CLI flags having highest priority.
 This lets you set defaults in code, override in config files for
 deployment, override with env vars in containers, and override
 individual values with flags for debugging.
*/

// -------------------------------------------------------------------------
// Configuration Validation
// -------------------------------------------------------------------------

/*
 Fail fast on invalid configuration. It's far better to crash at startup
 with a clear error message than to discover a misconfiguration at 3 AM
 when a request hits a bad code path.

 Common validations:
   - Required fields are not empty
   - Ports are in valid range (1-65535)
   - URLs are parseable
   - Timeouts are positive
   - File paths exist
   - Enum values are in the allowed set

 The pattern: a Validate() method that returns a list of all errors.
 Don't stop at the first error — report them all so the operator can
 fix everything in one pass.
*/

// ValidationError holds multiple validation errors.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(v.Errors, "; "))
}

// Add adds an error message to the validation error list.
func (v *ValidationError) Add(msg string) {
	v.Errors = append(v.Errors, msg)
}

// HasErrors returns true if there are any validation errors.
func (v *ValidationError) HasErrors() bool {
	return len(v.Errors) > 0
}

// Validate checks the server configuration for common errors.
func (c ServerConfig) Validate() error {
	v := &ValidationError{}

	if c.Host == "" {
		v.Add("host is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		v.Add(fmt.Sprintf("port must be 1-65535, got %d", c.Port))
	}
	if c.ReadTimeout <= 0 {
		v.Add("read_timeout must be positive")
	}
	if c.WriteTimeout <= 0 {
		v.Add("write_timeout must be positive")
	}

	if v.HasErrors() {
		return v
	}
	return nil
}

// -------------------------------------------------------------------------
// Configuration Layering
// -------------------------------------------------------------------------

/*
 In production, config comes from multiple sources with a priority order:

   1. Defaults (hardcoded in the application)
   2. Config file (for environment-specific settings)
   3. Environment variables (for container deployments)
   4. CLI flags (for debugging and overrides)

 Higher priority sources override lower ones. This means:
   - Developers get sane defaults out of the box
   - Ops can configure via files or environment variables
   - Anyone can override individual settings via CLI flags

 The implementation is straightforward: start with defaults, then
 apply each layer in order, only overriding non-zero values.
*/

// -------------------------------------------------------------------------
// Secret Management
// -------------------------------------------------------------------------

/*
 NEVER put secrets in:
   - Source code (even for "development only" — they leak via git history)
   - Config files committed to version control
   - Docker image layers (even deleted files exist in earlier layers)
   - Logs or error messages

 In production, use:
   - Environment variables (set by the deployment system)
   - Secret management services (AWS Secrets Manager, HashiCorp Vault)
   - Kubernetes Secrets (mounted as files or env vars)

 In development, use:
   - .env files (NEVER commit these — add to .gitignore)
   - Local secret manager or environment variables
   - Separate config for development vs production

 The key principle: secrets are injected at deploy time, not baked
 into the application.
*/

// -------------------------------------------------------------------------
// Feature Flags
// -------------------------------------------------------------------------

/*
 Feature flags let you enable/disable features at runtime without
 deploying new code. They're essential for:

   - Gradual rollouts (enable for 1% of users, then 10%, then all)
   - Kill switches (instantly disable a misbehaving feature)
   - A/B testing
   - Trunk-based development (merge incomplete features behind flags)

 A simple implementation uses a map of flag names to values. In
 production, you'd use a service like LaunchDarkly, Unleash, or
 a simple database-backed system.
*/

// -------------------------------------------------------------------------
// Functional Options for Services
// -------------------------------------------------------------------------

/*
 The functional options pattern is particularly powerful for service
 configuration. Instead of a massive config struct constructor, you
 provide option functions that each modify one aspect:

   server := NewServer(
       WithPort(8080),
       WithTimeout(30 * time.Second),
       WithLogger(myLogger),
   )

 Benefits:
   - Self-documenting — each option is a named function
   - Extensible — add new options without breaking existing callers
   - Optional — only specify what you want to override
   - Testable — compose different option sets for different tests
*/

// -------------------------------------------------------------------------
// Testing with Different Configurations
// -------------------------------------------------------------------------

/*
 Config-as-struct makes testing easy:

   func TestHandler(t *testing.T) {
       cfg := Config{
           Port:    0,           // random available port
           Timeout: time.Second, // short timeout for tests
           Debug:   true,        // verbose logging in tests
       }
       handler := NewHandler(cfg)
       // test handler...
   }

 No need for environment variable manipulation in tests. Just construct
 the config struct directly. If you DO need to test env var loading,
 use t.Setenv() which automatically cleans up after the test.
*/

// -------------------------------------------------------------------------
// Common Gotchas
// -------------------------------------------------------------------------

/*
 1. FORGETTING TO CALL flag.Parse()
    All flag values will be their defaults if you don't call Parse().
    No error — just silently wrong values.

 2. ENVIRONMENT VARIABLES ARE STRINGS
    "0" is truthy in most languages but means false in config.
    Always use explicit parsing functions.

 3. MISSING REQUIRED CONFIG
    Don't silently use a zero value for a database connection string.
    Validate and fail fast.

 4. SECRETS IN LOG OUTPUT
    Be careful with fmt.Printf("%+v", config) — it will print passwords.
    Implement a String() method that redacts sensitive fields.

 5. CONFIG DRIFT
    Multiple services reading the same config can diverge over time.
    Use shared config packages and version your config format.

 6. PORT CONFLICTS
    Using a hardcoded port in tests causes parallel test failures.
    Use port 0 (the OS assigns a random available port).
*/

// These are used to suppress "unused import" warnings for the lesson demo code.
var _ = url.Parse
var _ = os.LookupEnv
var _ = strconv.Atoi
