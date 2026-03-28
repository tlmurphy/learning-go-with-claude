package packages

import (
	"testing"
)

// =============================================================================
// Exercise 1: BankAccount — Exported vs Unexported
// =============================================================================

func TestNewBankAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountNum     string
		owner          string
		initialBalance float64
		wantErr        bool
		wantBalance    float64
	}{
		{
			name:           "valid account",
			accountNum:     "ACC001",
			owner:          "Alice",
			initialBalance: 100.0,
			wantErr:        false,
			wantBalance:    100.0,
		},
		{
			name:           "zero balance is valid",
			accountNum:     "ACC002",
			owner:          "Bob",
			initialBalance: 0.0,
			wantErr:        false,
			wantBalance:    0.0,
		},
		{
			name:           "negative balance is invalid",
			accountNum:     "ACC003",
			owner:          "Charlie",
			initialBalance: -50.0,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct, err := NewBankAccount(tt.accountNum, tt.owner, tt.initialBalance)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error for negative initial balance, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acct.AccountNumber != tt.accountNum {
				t.Errorf("AccountNumber = %q, want %q", acct.AccountNumber, tt.accountNum)
			}
			if acct.Owner != tt.owner {
				t.Errorf("Owner = %q, want %q", acct.Owner, tt.owner)
			}
			if acct.Balance() != tt.wantBalance {
				t.Errorf("Balance() = %f, want %f", acct.Balance(), tt.wantBalance)
			}
			if !acct.IsActive() {
				t.Error("new account should be active")
			}
		})
	}
}

func TestBankAccountDeposit(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *BankAccount
		amount      float64
		wantErr     bool
		wantBalance float64
	}{
		{
			name: "valid deposit",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:      50.0,
			wantErr:     false,
			wantBalance: 150.0,
		},
		{
			name: "deposit zero amount",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:  0,
			wantErr: true,
		},
		{
			name: "deposit negative amount",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:  -10.0,
			wantErr: true,
		},
		{
			name: "deposit to inactive account",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				a.Deactivate()
				return &a
			},
			amount:  50.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct := tt.setup()
			err := acct.Deposit(tt.amount)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acct.Balance() != tt.wantBalance {
				t.Errorf("Balance() = %f, want %f", acct.Balance(), tt.wantBalance)
			}
		})
	}
}

func TestBankAccountWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *BankAccount
		amount      float64
		wantErr     bool
		wantBalance float64
	}{
		{
			name: "valid withdrawal",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:      30.0,
			wantErr:     false,
			wantBalance: 70.0,
		},
		{
			name: "withdraw entire balance",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:      100.0,
			wantErr:     false,
			wantBalance: 0.0,
		},
		{
			name: "overdraft not allowed",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:  150.0,
			wantErr: true,
		},
		{
			name: "withdraw negative amount",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				return &a
			},
			amount:  -10.0,
			wantErr: true,
		},
		{
			name: "withdraw from inactive account",
			setup: func() *BankAccount {
				a, _ := NewBankAccount("ACC001", "Alice", 100.0)
				a.Deactivate()
				return &a
			},
			amount:  50.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct := tt.setup()
			err := acct.Withdraw(tt.amount)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acct.Balance() != tt.wantBalance {
				t.Errorf("Balance() = %f, want %f", acct.Balance(), tt.wantBalance)
			}
		})
	}
}

// =============================================================================
// Exercise 2: Plugin Registry
// =============================================================================

func TestPluginRegistry(t *testing.T) {
	// Clean state for this test
	ClearPlugins()

	t.Run("register and list plugins", func(t *testing.T) {
		RegisterPlugin("auth")
		RegisterPlugin("logging")
		RegisterPlugin("metrics")

		plugins := ListPlugins()
		if len(plugins) != 3 {
			t.Fatalf("expected 3 plugins, got %d", len(plugins))
		}

		expected := []string{"auth", "logging", "metrics"}
		for i, want := range expected {
			if plugins[i] != want {
				t.Errorf("plugin[%d] = %q, want %q", i, plugins[i], want)
			}
		}
	})

	t.Run("ListPlugins returns a copy", func(t *testing.T) {
		ClearPlugins()
		RegisterPlugin("original")

		plugins := ListPlugins()
		if len(plugins) == 0 {
			t.Fatal("ListPlugins returned empty after RegisterPlugin — implement RegisterPlugin and ListPlugins first")
		}
		plugins[0] = "modified" // Modify the returned slice

		// The original should be unchanged
		original := ListPlugins()
		if len(original) == 0 {
			t.Fatal("ListPlugins returned empty on second call")
		}
		if original[0] != "original" {
			t.Error("ListPlugins should return a copy; modifying the returned slice changed the internal state")
		}
	})

	t.Run("ClearPlugins resets registry", func(t *testing.T) {
		ClearPlugins()
		RegisterPlugin("plugin1")
		ClearPlugins()

		plugins := ListPlugins()
		if len(plugins) != 0 {
			t.Errorf("expected 0 plugins after clear, got %d", len(plugins))
		}
	})
}

// =============================================================================
// Exercise 3: Interface-Based Logger
// =============================================================================

func TestMemoryLogger(t *testing.T) {
	t.Run("creates non-nil logger", func(t *testing.T) {
		logger := NewMemoryLogger()
		if logger == nil {
			t.Fatal("NewMemoryLogger() returned nil")
		}
	})

	t.Run("logs info messages", func(t *testing.T) {
		logger := NewMemoryLogger()
		if logger == nil {
			t.Fatal("NewMemoryLogger() returned nil — implement NewMemoryLogger first")
		}
		logger.Info("server started")
		logger.Info("listening on :8080")

		msgs := logger.Messages()
		if len(msgs) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(msgs))
		}
		if msgs[0] != "[INFO] server started" {
			t.Errorf("message[0] = %q, want %q", msgs[0], "[INFO] server started")
		}
		if msgs[1] != "[INFO] listening on :8080" {
			t.Errorf("message[1] = %q, want %q", msgs[1], "[INFO] listening on :8080")
		}
	})

	t.Run("logs error messages", func(t *testing.T) {
		logger := NewMemoryLogger()
		if logger == nil {
			t.Fatal("NewMemoryLogger() returned nil — implement NewMemoryLogger first")
		}
		logger.Error("connection failed")

		msgs := logger.Messages()
		if len(msgs) != 1 {
			t.Fatalf("expected 1 message, got %d", len(msgs))
		}
		if msgs[0] != "[ERROR] connection failed" {
			t.Errorf("message = %q, want %q", msgs[0], "[ERROR] connection failed")
		}
	})

	t.Run("mixed message types preserve order", func(t *testing.T) {
		logger := NewMemoryLogger()
		if logger == nil {
			t.Fatal("NewMemoryLogger() returned nil — implement NewMemoryLogger first")
		}
		logger.Info("starting")
		logger.Error("something broke")
		logger.Info("recovered")

		msgs := logger.Messages()
		if len(msgs) != 3 {
			t.Fatalf("expected 3 messages, got %d", len(msgs))
		}
		expected := []string{
			"[INFO] starting",
			"[ERROR] something broke",
			"[INFO] recovered",
		}
		for i, want := range expected {
			if msgs[i] != want {
				t.Errorf("message[%d] = %q, want %q", i, msgs[i], want)
			}
		}
	})
}

// =============================================================================
// Exercise 4: Functional Options for DatabaseConfig
// =============================================================================

func TestNewDatabaseConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		cfg := NewDatabaseConfig()
		if cfg.Host != "localhost" {
			t.Errorf("default Host = %q, want %q", cfg.Host, "localhost")
		}
		if cfg.Port != 5432 {
			t.Errorf("default Port = %d, want %d", cfg.Port, 5432)
		}
		if cfg.Database != "app" {
			t.Errorf("default Database = %q, want %q", cfg.Database, "app")
		}
		if cfg.MaxConnections != 10 {
			t.Errorf("default MaxConnections = %d, want %d", cfg.MaxConnections, 10)
		}
		if cfg.TimeoutSeconds != 30 {
			t.Errorf("default TimeoutSeconds = %d, want %d", cfg.TimeoutSeconds, 30)
		}
		if cfg.SSLEnabled != false {
			t.Errorf("default SSLEnabled = %v, want %v", cfg.SSLEnabled, false)
		}
	})

	t.Run("with custom options", func(t *testing.T) {
		cfg := NewDatabaseConfig(
			WithDBHost("db.example.com"),
			WithDBPort(3306),
			WithDBName("myapp"),
			WithMaxConnections(50),
			WithDBTimeout(60),
			WithSSL(true),
		)
		if cfg.Host != "db.example.com" {
			t.Errorf("Host = %q, want %q", cfg.Host, "db.example.com")
		}
		if cfg.Port != 3306 {
			t.Errorf("Port = %d, want %d", cfg.Port, 3306)
		}
		if cfg.Database != "myapp" {
			t.Errorf("Database = %q, want %q", cfg.Database, "myapp")
		}
		if cfg.MaxConnections != 50 {
			t.Errorf("MaxConnections = %d, want %d", cfg.MaxConnections, 50)
		}
		if cfg.TimeoutSeconds != 60 {
			t.Errorf("TimeoutSeconds = %d, want %d", cfg.TimeoutSeconds, 60)
		}
		if cfg.SSLEnabled != true {
			t.Errorf("SSLEnabled = %v, want %v", cfg.SSLEnabled, true)
		}
	})

	t.Run("partial options use defaults for the rest", func(t *testing.T) {
		cfg := NewDatabaseConfig(
			WithDBHost("custom-host"),
			WithSSL(true),
		)
		if cfg.Host != "custom-host" {
			t.Errorf("Host = %q, want %q", cfg.Host, "custom-host")
		}
		if cfg.Port != 5432 {
			t.Errorf("Port should use default 5432, got %d", cfg.Port)
		}
		if cfg.SSLEnabled != true {
			t.Errorf("SSLEnabled = %v, want %v", cfg.SSLEnabled, true)
		}
	})
}

// =============================================================================
// Exercise 5: Breaking Circular Dependencies
// =============================================================================

func TestCircularDependencyBreaking(t *testing.T) {
	// This test demonstrates how interfaces break circular dependencies.
	// OrderService depends on PermissionChecker (an interface).
	// UserService depends on OrderLookup (an interface).
	// Neither depends on the other's concrete type.

	t.Run("order service checks permissions", func(t *testing.T) {
		// We need to set up both services with cross-references.
		// Step 1: Create UserService with a nil OrderLookup (we'll set it after).
		userSvc := NewUserService(nil)
		if userSvc == nil {
			t.Fatal("NewUserService returned nil")
		}

		// Step 2: Create OrderService with UserService as the PermissionChecker.
		orderSvc := NewOrderService(userSvc)
		if orderSvc == nil {
			t.Fatal("NewOrderService returned nil")
		}

		// Grant permission and place order
		userSvc.GrantPermission(1, "place_order")
		err := orderSvc.PlaceOrder(1, 99.99)
		if err != nil {
			t.Fatalf("expected successful order, got error: %v", err)
		}

		// Verify order count
		count := orderSvc.OrderCountForUser(1)
		if count != 1 {
			t.Errorf("OrderCountForUser(1) = %d, want 1", count)
		}
	})

	t.Run("order rejected without permission", func(t *testing.T) {
		userSvc := NewUserService(nil)
		orderSvc := NewOrderService(userSvc)

		// Don't grant permission
		err := orderSvc.PlaceOrder(1, 99.99)
		if err == nil {
			t.Error("expected error when placing order without permission, got nil")
		}
	})

	t.Run("user service can look up order count", func(t *testing.T) {
		// Create order service first (with a permissive checker for simplicity)
		permissive := &alwaysAllowed{}
		orderSvc := NewOrderService(permissive)

		// Create user service with order lookup
		userSvc := NewUserService(orderSvc)

		// Place some orders
		_ = orderSvc.PlaceOrder(1, 10.0)
		_ = orderSvc.PlaceOrder(1, 20.0)
		_ = orderSvc.PlaceOrder(2, 30.0)

		// User service should be able to look up order counts
		if count := userSvc.UserOrderCount(1); count != 2 {
			t.Errorf("UserOrderCount(1) = %d, want 2", count)
		}
		if count := userSvc.UserOrderCount(2); count != 1 {
			t.Errorf("UserOrderCount(2) = %d, want 1", count)
		}
	})
}

// alwaysAllowed is a test helper that grants all permissions.
type alwaysAllowed struct{}

func (a *alwaysAllowed) HasPermission(userID int, permission string) bool {
	return true
}

// =============================================================================
// Exercise 6: Email Validator
// =============================================================================

func TestEmailValidator(t *testing.T) {
	t.Run("valid emails", func(t *testing.T) {
		v := NewEmailValidator()
		if v == nil {
			t.Fatal("NewEmailValidator() returned nil")
		}

		validEmails := []string{
			"user@example.com",
			"first.last@domain.org",
			"user+tag@mail.co.uk",
			"a@b.co",
		}
		for _, email := range validEmails {
			if err := v.Validate(email); err != nil {
				t.Errorf("Validate(%q) = %v, want nil (should be valid)", email, err)
			}
		}
	})

	t.Run("invalid emails", func(t *testing.T) {
		v := NewEmailValidator()
		if v == nil {
			t.Fatal("NewEmailValidator() returned nil")
		}

		invalidEmails := []string{
			"",               // empty
			"noatsign",       // no @
			"@domain.com",    // no local part
			"user@",          // no domain
			"user@@mail.com", // double @
			"user@nodot",     // domain has no dot
		}
		for _, email := range invalidEmails {
			if err := v.Validate(email); err == nil {
				t.Errorf("Validate(%q) = nil, want error (should be invalid)", email)
			}
		}
	})

	t.Run("blocked domains", func(t *testing.T) {
		v := NewEmailValidator(WithBlockedDomains("spam.com", "trash.net"))
		if v == nil {
			t.Fatal("NewEmailValidator() returned nil")
		}

		if err := v.Validate("user@spam.com"); err == nil {
			t.Error("expected error for blocked domain spam.com")
		}
		if err := v.Validate("user@trash.net"); err == nil {
			t.Error("expected error for blocked domain trash.net")
		}
		if err := v.Validate("user@legit.com"); err != nil {
			t.Errorf("legit.com should not be blocked: %v", err)
		}
	})

	t.Run("max length", func(t *testing.T) {
		v := NewEmailValidator(WithMaxLength(20))
		if v == nil {
			t.Fatal("NewEmailValidator() returned nil")
		}

		if err := v.Validate("a@b.co"); err != nil {
			t.Errorf("short email should be valid: %v", err)
		}
		if err := v.Validate("verylongemail@verylongdomain.com"); err == nil {
			t.Error("expected error for email exceeding max length")
		}
	})
}

// =============================================================================
// Verify Init Functions Ran
// =============================================================================

func TestInitFunctionsRan(t *testing.T) {
	order := InitOrder()
	if len(order) < 2 {
		t.Errorf("expected at least 2 init entries, got %d", len(order))
	}
}

// =============================================================================
// Verify Lesson Code Works
// =============================================================================

func TestLessonExamples(t *testing.T) {
	t.Run("NewUser and CheckPassword", func(t *testing.T) {
		user := NewUser(1, "Alice", "secret123")
		if !user.CheckPassword("secret123") {
			t.Error("CheckPassword should return true for correct password")
		}
		if user.CheckPassword("wrongpassword") {
			t.Error("CheckPassword should return false for incorrect password")
		}
	})

	t.Run("Repository interface", func(t *testing.T) {
		repo := NewMemoryRepository()
		user := NewUser(1, "Bob", "pass")
		if err := repo.Save(user); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
		found, err := repo.FindByID(1)
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}
		if found.Name != "Bob" {
			t.Errorf("found.Name = %q, want %q", found.Name, "Bob")
		}
	})

	t.Run("Functional options server", func(t *testing.T) {
		s := NewServer(WithPort(9090), WithHost("localhost"))
		if s.Addr() != "localhost:9090" {
			t.Errorf("Addr() = %q, want %q", s.Addr(), "localhost:9090")
		}
		if s.Timeout() != 30 {
			t.Errorf("Timeout() = %d, want default 30", s.Timeout())
		}
	})
}
