package configuration

import (
	"strings"
	"testing"
	"time"
)

// =========================================================================
// Test Exercise 1: Config from Environment Variables
// =========================================================================

func TestEnvConfig(t *testing.T) {
	t.Run("defaults when no env vars set", func(t *testing.T) {
		// Clear any env vars that might be set
		t.Setenv("APP_HOST", "")
		// Use a unique approach: unset by setting empty, but our loader
		// should use LookupEnv. We need to properly test defaults.
		// t.Setenv sets the var, so we need a different approach.

		// For this test, we rely on env vars not being set for these keys.
		// We'll test with explicit env vars in the next subtest.
		cfg := LoadAppConfig()

		// These tests may be affected by actual env vars, so we just
		// check that the function doesn't panic and returns valid values.
		if cfg.Port == 0 && cfg.Host == "" {
			// This means the function returned zero values — not implemented yet
			t.Skip("LoadAppConfig not yet implemented (returns zero values)")
		}
	})

	t.Run("loads from environment variables", func(t *testing.T) {
		t.Setenv("APP_HOST", "0.0.0.0")
		t.Setenv("APP_PORT", "3000")
		t.Setenv("DATABASE_URL", "postgres://db:5432/testdb")
		t.Setenv("LOG_LEVEL", "debug")
		t.Setenv("APP_DEBUG", "true")
		t.Setenv("READ_TIMEOUT", "30s")
		t.Setenv("WRITE_TIMEOUT", "45s")

		cfg := LoadAppConfig()

		if cfg.Host != "0.0.0.0" {
			t.Errorf("Expected Host='0.0.0.0', got %q", cfg.Host)
		}
		if cfg.Port != 3000 {
			t.Errorf("Expected Port=3000, got %d", cfg.Port)
		}
		if cfg.DatabaseURL != "postgres://db:5432/testdb" {
			t.Errorf("Expected DatabaseURL='postgres://db:5432/testdb', got %q", cfg.DatabaseURL)
		}
		if cfg.LogLevel != "debug" {
			t.Errorf("Expected LogLevel='debug', got %q", cfg.LogLevel)
		}
		if cfg.Debug != true {
			t.Errorf("Expected Debug=true, got %t", cfg.Debug)
		}
		if cfg.ReadTimeout != 30*time.Second {
			t.Errorf("Expected ReadTimeout=30s, got %v", cfg.ReadTimeout)
		}
		if cfg.WriteTimeout != 45*time.Second {
			t.Errorf("Expected WriteTimeout=45s, got %v", cfg.WriteTimeout)
		}
	})

	t.Run("uses defaults for missing env vars", func(t *testing.T) {
		// Unset all relevant env vars by not setting them
		// (t.Setenv cleanup restores previous values after each subtest)
		cfg := LoadAppConfig()

		if cfg.Host == "" {
			t.Error("Expected a default Host, got empty string")
		}
		if cfg.Port == 0 {
			t.Error("Expected a default Port, got 0")
		}
	})
}

// =========================================================================
// Test Exercise 2: Functional Options Pattern
// =========================================================================

func TestServiceConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := NewServiceConfig()

		if cfg.Name != "service" {
			t.Errorf("Expected Name='service', got %q", cfg.Name)
		}
		if cfg.Port != 8080 {
			t.Errorf("Expected Port=8080, got %d", cfg.Port)
		}
		if cfg.ReadTimeout != 15*time.Second {
			t.Errorf("Expected ReadTimeout=15s, got %v", cfg.ReadTimeout)
		}
		if cfg.MaxBodySize != 1048576 {
			t.Errorf("Expected MaxBodySize=1048576, got %d", cfg.MaxBodySize)
		}
		if cfg.EnableCORS != false {
			t.Errorf("Expected EnableCORS=false, got %t", cfg.EnableCORS)
		}
	})

	t.Run("with options", func(t *testing.T) {
		cfg := NewServiceConfig(
			WithName("api"),
			WithServicePort(3000),
			WithTimeouts(30*time.Second, 60*time.Second),
			WithMaxBodySize(5*1024*1024),
			WithCORS("https://example.com", "https://other.com"),
		)

		if cfg.Name != "api" {
			t.Errorf("Expected Name='api', got %q", cfg.Name)
		}
		if cfg.Port != 3000 {
			t.Errorf("Expected Port=3000, got %d", cfg.Port)
		}
		if cfg.ReadTimeout != 30*time.Second {
			t.Errorf("Expected ReadTimeout=30s, got %v", cfg.ReadTimeout)
		}
		if cfg.WriteTimeout != 60*time.Second {
			t.Errorf("Expected WriteTimeout=60s, got %v", cfg.WriteTimeout)
		}
		if cfg.MaxBodySize != 5*1024*1024 {
			t.Errorf("Expected MaxBodySize=5242880, got %d", cfg.MaxBodySize)
		}
		if !cfg.EnableCORS {
			t.Error("Expected EnableCORS=true")
		}
		if len(cfg.AllowedOrigins) != 2 {
			t.Errorf("Expected 2 allowed origins, got %d", len(cfg.AllowedOrigins))
		}
	})

	t.Run("with single option", func(t *testing.T) {
		cfg := NewServiceConfig(WithServicePort(9090))

		if cfg.Port != 9090 {
			t.Errorf("Expected Port=9090, got %d", cfg.Port)
		}
		// Other values should be defaults
		if cfg.Name != "service" {
			t.Errorf("Expected Name='service', got %q", cfg.Name)
		}
	})
}

// =========================================================================
// Test Exercise 3: Config Validator
// =========================================================================

func TestConfigValidator(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := AppConfig{
			Host:         "localhost",
			Port:         8080,
			DatabaseURL:  "postgres://localhost:5432/app",
			LogLevel:     "info",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err != nil {
			t.Errorf("Expected no error for valid config, got: %v", err)
		}
	})

	t.Run("empty host", func(t *testing.T) {
		cfg := AppConfig{
			Port:         8080,
			DatabaseURL:  "postgres://localhost:5432/app",
			LogLevel:     "info",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Error("Expected error for empty host")
		}
	})

	t.Run("invalid port", func(t *testing.T) {
		cfg := AppConfig{
			Host:         "localhost",
			Port:         0,
			DatabaseURL:  "postgres://localhost:5432/app",
			LogLevel:     "info",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Error("Expected error for port=0")
		}
	})

	t.Run("port too high", func(t *testing.T) {
		cfg := AppConfig{
			Host:         "localhost",
			Port:         70000,
			DatabaseURL:  "postgres://localhost:5432/app",
			LogLevel:     "info",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Error("Expected error for port=70000")
		}
	})

	t.Run("empty database URL", func(t *testing.T) {
		cfg := AppConfig{
			Host:         "localhost",
			Port:         8080,
			LogLevel:     "info",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Error("Expected error for empty DatabaseURL")
		}
	})

	t.Run("invalid log level", func(t *testing.T) {
		cfg := AppConfig{
			Host:         "localhost",
			Port:         8080,
			DatabaseURL:  "postgres://localhost:5432/app",
			LogLevel:     "verbose",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Error("Expected error for invalid log level 'verbose'")
		}
	})

	t.Run("negative timeout", func(t *testing.T) {
		cfg := AppConfig{
			Host:         "localhost",
			Port:         8080,
			DatabaseURL:  "postgres://localhost:5432/app",
			LogLevel:     "info",
			ReadTimeout:  -1 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Error("Expected error for negative read timeout")
		}
	})

	t.Run("multiple errors reported", func(t *testing.T) {
		cfg := AppConfig{
			// Host is empty, Port is 0, DatabaseURL is empty
			LogLevel:     "invalid",
			ReadTimeout:  -1,
			WriteTimeout: -1,
		}
		err := ValidateAppConfig(cfg)
		if err == nil {
			t.Fatal("Expected validation errors")
		}
		ce, ok := err.(*ConfigErrors)
		if !ok {
			t.Fatalf("Expected *ConfigErrors, got %T", err)
		}
		if len(ce.Messages()) < 3 {
			t.Errorf("Expected at least 3 errors, got %d: %v", len(ce.Messages()), ce.Messages())
		}
	})
}

// =========================================================================
// Test Exercise 4: Config Layering
// =========================================================================

func TestLayeredConfig(t *testing.T) {
	t.Run("defaults only", func(t *testing.T) {
		cfg := LoadLayeredConfig(nil, nil)

		if cfg.Host.Value != "localhost" {
			t.Errorf("Expected Host='localhost', got %q", cfg.Host.Value)
		}
		if cfg.Host.Source != "default" {
			t.Errorf("Expected Host source='default', got %q", cfg.Host.Source)
		}
		if cfg.Port.Value != 8080 {
			t.Errorf("Expected Port=8080, got %d", cfg.Port.Value)
		}
		if cfg.Debug.Value != false {
			t.Errorf("Expected Debug=false, got %t", cfg.Debug.Value)
		}
		if cfg.Timeout.Value != 30*time.Second {
			t.Errorf("Expected Timeout=30s, got %v", cfg.Timeout.Value)
		}
	})

	t.Run("env overrides defaults", func(t *testing.T) {
		envVars := map[string]string{
			"APP_HOST": "0.0.0.0",
			"APP_PORT": "3000",
		}
		cfg := LoadLayeredConfig(envVars, nil)

		if cfg.Host.Value != "0.0.0.0" {
			t.Errorf("Expected Host='0.0.0.0', got %q", cfg.Host.Value)
		}
		if cfg.Host.Source != "env" {
			t.Errorf("Expected Host source='env', got %q", cfg.Host.Source)
		}
		if cfg.Port.Value != 3000 {
			t.Errorf("Expected Port=3000, got %d", cfg.Port.Value)
		}
		if cfg.Port.Source != "env" {
			t.Errorf("Expected Port source='env', got %q", cfg.Port.Source)
		}
		// Debug should still be default
		if cfg.Debug.Source != "default" {
			t.Errorf("Expected Debug source='default', got %q", cfg.Debug.Source)
		}
	})

	t.Run("flags override env", func(t *testing.T) {
		envVars := map[string]string{
			"APP_HOST":  "0.0.0.0",
			"APP_PORT":  "3000",
			"APP_DEBUG": "true",
		}
		flags := map[string]string{
			"host": "custom.host",
			"port": "9090",
		}
		cfg := LoadLayeredConfig(envVars, flags)

		if cfg.Host.Value != "custom.host" {
			t.Errorf("Expected Host='custom.host' (flag), got %q", cfg.Host.Value)
		}
		if cfg.Host.Source != "flag" {
			t.Errorf("Expected Host source='flag', got %q", cfg.Host.Source)
		}
		if cfg.Port.Value != 9090 {
			t.Errorf("Expected Port=9090 (flag), got %d", cfg.Port.Value)
		}
		// Debug should come from env (no flag override)
		if cfg.Debug.Value != true {
			t.Errorf("Expected Debug=true (env), got %t", cfg.Debug.Value)
		}
		if cfg.Debug.Source != "env" {
			t.Errorf("Expected Debug source='env', got %q", cfg.Debug.Source)
		}
	})

	t.Run("invalid env port keeps default", func(t *testing.T) {
		envVars := map[string]string{
			"APP_PORT": "not-a-number",
		}
		cfg := LoadLayeredConfig(envVars, nil)

		if cfg.Port.Value != 8080 {
			t.Errorf("Expected Port=8080 (default, invalid env), got %d", cfg.Port.Value)
		}
		if cfg.Port.Source != "default" {
			t.Errorf("Expected Port source='default' (invalid env), got %q", cfg.Port.Source)
		}
	})

	t.Run("timeout parsing", func(t *testing.T) {
		envVars := map[string]string{
			"APP_TIMEOUT": "5s",
		}
		cfg := LoadLayeredConfig(envVars, nil)

		if cfg.Timeout.Value != 5*time.Second {
			t.Errorf("Expected Timeout=5s, got %v", cfg.Timeout.Value)
		}
		if cfg.Timeout.Source != "env" {
			t.Errorf("Expected Timeout source='env', got %q", cfg.Timeout.Source)
		}
	})
}

// =========================================================================
// Test Exercise 5: Feature Flags
// =========================================================================

func TestFeatureFlags(t *testing.T) {
	t.Run("get bool flag", func(t *testing.T) {
		ff := NewFeatureFlags(map[string]interface{}{
			"dark_mode": true,
			"new_ui":    false,
		})
		if ff == nil {
			t.Fatal("NewFeatureFlags returned nil")
		}

		if !ff.IsEnabled("dark_mode", false) {
			t.Error("Expected dark_mode to be enabled")
		}
		if ff.IsEnabled("new_ui", true) {
			t.Error("Expected new_ui to be disabled")
		}
		if ff.IsEnabled("missing", false) != false {
			t.Error("Expected missing flag to return default (false)")
		}
		if ff.IsEnabled("missing", true) != true {
			t.Error("Expected missing flag to return default (true)")
		}
	})

	t.Run("get string flag", func(t *testing.T) {
		ff := NewFeatureFlags(map[string]interface{}{
			"theme": "blue",
		})
		if ff == nil {
			t.Fatal("NewFeatureFlags returned nil")
		}

		if ff.GetString("theme", "default") != "blue" {
			t.Errorf("Expected theme='blue', got %q", ff.GetString("theme", "default"))
		}
		if ff.GetString("missing", "fallback") != "fallback" {
			t.Errorf("Expected 'fallback', got %q", ff.GetString("missing", "fallback"))
		}
	})

	t.Run("get int flag", func(t *testing.T) {
		ff := NewFeatureFlags(map[string]interface{}{
			"max_items": 50,
		})
		if ff == nil {
			t.Fatal("NewFeatureFlags returned nil")
		}

		if ff.GetInt("max_items", 10) != 50 {
			t.Errorf("Expected max_items=50, got %d", ff.GetInt("max_items", 10))
		}
		if ff.GetInt("missing", 10) != 10 {
			t.Errorf("Expected default 10, got %d", ff.GetInt("missing", 10))
		}
	})

	t.Run("type mismatch returns default", func(t *testing.T) {
		ff := NewFeatureFlags(map[string]interface{}{
			"count": "not-an-int",
		})
		if ff == nil {
			t.Fatal("NewFeatureFlags returned nil")
		}

		if ff.GetInt("count", 42) != 42 {
			t.Error("Expected default for type mismatch")
		}
	})

	t.Run("set flag", func(t *testing.T) {
		ff := NewFeatureFlags(nil)
		if ff == nil {
			t.Fatal("NewFeatureFlags returned nil")
		}

		ff.Set("new_feature", true)
		if !ff.IsEnabled("new_feature", false) {
			t.Error("Expected new_feature to be true after Set")
		}
	})

	t.Run("all returns copy", func(t *testing.T) {
		ff := NewFeatureFlags(map[string]interface{}{
			"a": true,
			"b": "hello",
		})
		if ff == nil {
			t.Fatal("NewFeatureFlags returned nil")
		}

		all := ff.All()
		if len(all) != 2 {
			t.Errorf("Expected 2 flags, got %d", len(all))
		}
		// Modifying the returned map shouldn't affect the original
		all["c"] = "new"
		if len(ff.All()) != 2 {
			t.Error("Modifying All() result should not affect original")
		}
	})
}

// =========================================================================
// Test Exercise 6: Config Watcher
// =========================================================================

func TestConfigWatcher(t *testing.T) {
	t.Run("get initial values", func(t *testing.T) {
		w := NewConfigWatcher(map[string]string{"host": "localhost"})
		if w == nil {
			t.Fatal("NewConfigWatcher returned nil")
		}

		if w.Get("host") != "localhost" {
			t.Errorf("Expected host='localhost', got %q", w.Get("host"))
		}
		if w.Get("missing") != "" {
			t.Errorf("Expected empty string for missing key, got %q", w.Get("missing"))
		}
	})

	t.Run("watch receives changes", func(t *testing.T) {
		w := NewConfigWatcher(map[string]string{"host": "localhost"})
		if w == nil {
			t.Fatal("NewConfigWatcher returned nil")
		}

		ch := w.Watch()

		w.Update("host", "0.0.0.0")

		select {
		case change := <-ch:
			if change.Key != "host" {
				t.Errorf("Expected key='host', got %q", change.Key)
			}
			if change.OldValue != "localhost" {
				t.Errorf("Expected OldValue='localhost', got %q", change.OldValue)
			}
			if change.NewValue != "0.0.0.0" {
				t.Errorf("Expected NewValue='0.0.0.0', got %q", change.NewValue)
			}
		default:
			t.Error("Expected to receive a change notification")
		}
	})

	t.Run("no notification for same value", func(t *testing.T) {
		w := NewConfigWatcher(map[string]string{"host": "localhost"})
		if w == nil {
			t.Fatal("NewConfigWatcher returned nil")
		}

		ch := w.Watch()
		w.Update("host", "localhost") // same value

		select {
		case change := <-ch:
			t.Errorf("Should not receive notification for unchanged value, got %+v", change)
		default:
			// Good — no notification
		}
	})

	t.Run("new key triggers notification", func(t *testing.T) {
		w := NewConfigWatcher(map[string]string{})
		if w == nil {
			t.Fatal("NewConfigWatcher returned nil")
		}

		ch := w.Watch()
		w.Update("port", "8080")

		select {
		case change := <-ch:
			if change.Key != "port" {
				t.Errorf("Expected key='port', got %q", change.Key)
			}
			if change.OldValue != "" {
				t.Errorf("Expected OldValue='', got %q", change.OldValue)
			}
			if change.NewValue != "8080" {
				t.Errorf("Expected NewValue='8080', got %q", change.NewValue)
			}
		default:
			t.Error("Expected notification for new key")
		}
	})

	t.Run("multiple watchers", func(t *testing.T) {
		w := NewConfigWatcher(map[string]string{"x": "1"})
		if w == nil {
			t.Fatal("NewConfigWatcher returned nil")
		}

		ch1 := w.Watch()
		ch2 := w.Watch()

		w.Update("x", "2")

		received1 := false
		received2 := false

		select {
		case <-ch1:
			received1 = true
		default:
		}
		select {
		case <-ch2:
			received2 = true
		default:
		}

		if !received1 || !received2 {
			t.Error("Expected both watchers to receive notification")
		}
	})
}

// =========================================================================
// Test Exercise 7: Database Config Parser
// =========================================================================

func TestDatabaseConfigParser(t *testing.T) {
	t.Run("full URL", func(t *testing.T) {
		cfg, err := ParseDatabaseURL("postgres://admin:secret@db.example.com:5433/mydb?sslmode=require")
		if err != nil {
			t.Fatalf("ParseDatabaseURL error: %v", err)
		}
		if cfg.Host != "db.example.com" {
			t.Errorf("Expected Host='db.example.com', got %q", cfg.Host)
		}
		if cfg.Port != 5433 {
			t.Errorf("Expected Port=5433, got %d", cfg.Port)
		}
		if cfg.Database != "mydb" {
			t.Errorf("Expected Database='mydb', got %q", cfg.Database)
		}
		if cfg.User != "admin" {
			t.Errorf("Expected User='admin', got %q", cfg.User)
		}
		if cfg.Password != "secret" {
			t.Errorf("Expected Password='secret', got %q", cfg.Password)
		}
		if cfg.SSLMode != "require" {
			t.Errorf("Expected SSLMode='require', got %q", cfg.SSLMode)
		}
	})

	t.Run("URL without port defaults to 5432", func(t *testing.T) {
		cfg, err := ParseDatabaseURL("postgres://user:pass@localhost/testdb")
		if err != nil {
			t.Fatalf("ParseDatabaseURL error: %v", err)
		}
		if cfg.Port != 5432 {
			t.Errorf("Expected default Port=5432, got %d", cfg.Port)
		}
	})

	t.Run("URL without sslmode defaults to disable", func(t *testing.T) {
		cfg, err := ParseDatabaseURL("postgres://user:pass@localhost/testdb")
		if err != nil {
			t.Fatalf("ParseDatabaseURL error: %v", err)
		}
		if cfg.SSLMode != "disable" {
			t.Errorf("Expected default SSLMode='disable', got %q", cfg.SSLMode)
		}
	})

	t.Run("URL with extra options", func(t *testing.T) {
		cfg, err := ParseDatabaseURL("postgres://user:pass@localhost/db?sslmode=require&connect_timeout=10")
		if err != nil {
			t.Fatalf("ParseDatabaseURL error: %v", err)
		}
		if cfg.Options == nil {
			t.Fatal("Expected Options to be non-nil")
		}
		if cfg.Options["connect_timeout"] != "10" {
			t.Errorf("Expected Options[connect_timeout]='10', got %q", cfg.Options["connect_timeout"])
		}
	})

	t.Run("build URL from config", func(t *testing.T) {
		cfg := DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "mydb",
			User:     "admin",
			Password: "secret",
			SSLMode:  "disable",
		}

		url := cfg.BuildDatabaseURL()
		if !strings.Contains(url, "postgres://") {
			t.Errorf("Expected URL to start with 'postgres://', got %q", url)
		}
		if !strings.Contains(url, "admin:secret@") {
			t.Errorf("Expected URL to contain 'admin:secret@', got %q", url)
		}
		if !strings.Contains(url, "localhost:5432") {
			t.Errorf("Expected URL to contain 'localhost:5432', got %q", url)
		}
		if !strings.Contains(url, "/mydb") {
			t.Errorf("Expected URL to contain '/mydb', got %q", url)
		}
		if !strings.Contains(url, "sslmode=disable") {
			t.Errorf("Expected URL to contain 'sslmode=disable', got %q", url)
		}
	})

	t.Run("invalid URL returns error", func(t *testing.T) {
		_, err := ParseDatabaseURL("not a valid url ://")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})
}

// =========================================================================
// Test Exercise 8: Complete Server Configuration
// =========================================================================

func TestFullServerConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := DefaultFullServerConfig()

		if cfg.Server.Port != 8080 {
			t.Errorf("Expected Port=8080, got %d", cfg.Server.Port)
		}
		if cfg.Server.ReadTimeout != 15*time.Second {
			t.Errorf("Expected ReadTimeout=15s, got %v", cfg.Server.ReadTimeout)
		}
		if cfg.Database.MaxOpenConns != 25 {
			t.Errorf("Expected MaxOpenConns=25, got %d", cfg.Database.MaxOpenConns)
		}
		if cfg.Logging.Level != "info" {
			t.Errorf("Expected LogLevel='info', got %q", cfg.Logging.Level)
		}
		if cfg.Logging.Format != "json" {
			t.Errorf("Expected LogFormat='json', got %q", cfg.Logging.Format)
		}
		if cfg.Features == nil {
			t.Error("Expected Features to be non-nil (empty map)")
		}
	})

	t.Run("valid config passes validation", func(t *testing.T) {
		cfg := DefaultFullServerConfig()
		err := ValidateFullServerConfig(cfg)
		if err != nil {
			t.Errorf("Expected no error for default config, got: %v", err)
		}
	})

	t.Run("invalid config fails validation", func(t *testing.T) {
		cfg := FullServerConfig{
			Server: ServerSection{
				Port:         0,
				ReadTimeout:  -1,
				WriteTimeout: -1,
			},
			Database: DatabaseSection{
				MaxOpenConns: 0,
			},
			Logging: LoggingSection{
				Level:  "verbose",
				Format: "xml",
			},
		}
		err := ValidateFullServerConfig(cfg)
		if err == nil {
			t.Fatal("Expected validation errors")
		}
		ce, ok := err.(*ConfigErrors)
		if !ok {
			t.Fatalf("Expected *ConfigErrors, got %T", err)
		}
		if len(ce.Messages()) < 4 {
			t.Errorf("Expected at least 4 errors, got %d: %v", len(ce.Messages()), ce.Messages())
		}
	})

	t.Run("string representation masks password", func(t *testing.T) {
		cfg := DefaultFullServerConfig()
		cfg.Database.URL = "postgres://admin:secret@localhost:5432/app"
		cfg.Features = map[string]bool{
			"dark_mode": true,
		}

		s := cfg.String()
		if s == "" {
			t.Fatal("String() returned empty string")
		}
		if strings.Contains(s, "secret") {
			t.Error("String() should not contain the database password")
		}
		if !strings.Contains(s, "****") {
			t.Error("String() should mask the password with ****")
		}
	})
}
