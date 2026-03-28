package config

import "os"

// Config holds all configuration for the bookstore API server.
type Config struct {
	Port           string
	AllowedOrigins []string
	LogLevel       string
}

// Load reads configuration from environment variables, falling back to
// sensible defaults.
//
// TODO: Expand this to read all config fields. Consider supporting a
// .env file for local development.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return Config{
		Port:           port,
		AllowedOrigins: []string{"*"},
		LogLevel:       logLevel,
	}
}
