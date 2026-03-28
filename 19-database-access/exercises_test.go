package database

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

// =========================================================================
// Test Exercise 1 & 2: Product Repository
// =========================================================================

func TestProductRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("create and get product", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		p := &Product{
			ID:          "prod-1",
			Name:        "Widget",
			Description: "A fine widget",
			PriceCents:  999,
			Category:    "gadgets",
			InStock:     true,
		}

		err := repo.Create(ctx, p)
		if err != nil {
			t.Fatalf("Create error: %v", err)
		}

		got, err := repo.GetByID(ctx, "prod-1")
		if err != nil {
			t.Fatalf("GetByID error: %v", err)
		}

		if got.Name != "Widget" {
			t.Errorf("Expected Name='Widget', got %q", got.Name)
		}
		if got.PriceCents != 999 {
			t.Errorf("Expected PriceCents=999, got %d", got.PriceCents)
		}
		if got.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set, got zero time")
		}
	})

	t.Run("create duplicate returns ErrAlreadyExists", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		p := &Product{ID: "prod-1", Name: "Widget"}
		repo.Create(ctx, p)

		err := repo.Create(ctx, &Product{ID: "prod-1", Name: "Other"})
		if !errors.Is(err, ErrAlreadyExists) {
			t.Errorf("Expected ErrAlreadyExists, got: %v", err)
		}
	})

	t.Run("create with empty ID returns ErrInvalidInput", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		err := repo.Create(ctx, &Product{Name: "No ID"})
		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("Expected ErrInvalidInput, got: %v", err)
		}
	})

	t.Run("get nonexistent returns ErrNotFound", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		_, err := repo.GetByID(ctx, "nonexistent")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got: %v", err)
		}
	})

	t.Run("list with pagination", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		for i := 0; i < 5; i++ {
			repo.Create(ctx, &Product{
				ID:   fmt.Sprintf("prod-%d", i),
				Name: fmt.Sprintf("Product %d", i),
			})
		}

		// Get all
		all, err := repo.List(ctx, 0, 0)
		if err != nil {
			t.Fatalf("List error: %v", err)
		}
		if len(all) != 5 {
			t.Errorf("Expected 5 products, got %d", len(all))
		}

		// Get with limit
		limited, err := repo.List(ctx, 2, 0)
		if err != nil {
			t.Fatalf("List error: %v", err)
		}
		if len(limited) != 2 {
			t.Errorf("Expected 2 products with limit=2, got %d", len(limited))
		}

		// Get with offset beyond range
		empty, err := repo.List(ctx, 10, 100)
		if err != nil {
			t.Fatalf("List error: %v", err)
		}
		if len(empty) != 0 {
			t.Errorf("Expected 0 products with large offset, got %d", len(empty))
		}
		if empty == nil {
			t.Error("Expected empty slice (not nil) when no results")
		}
	})

	t.Run("update product", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		repo.Create(ctx, &Product{
			ID:   "prod-1",
			Name: "Widget",
		})

		time.Sleep(time.Millisecond) // ensure UpdatedAt differs

		err := repo.Update(ctx, &Product{
			ID:   "prod-1",
			Name: "Super Widget",
		})
		if err != nil {
			t.Fatalf("Update error: %v", err)
		}

		got, _ := repo.GetByID(ctx, "prod-1")
		if got.Name != "Super Widget" {
			t.Errorf("Expected Name='Super Widget', got %q", got.Name)
		}
	})

	t.Run("update nonexistent returns ErrNotFound", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		err := repo.Update(ctx, &Product{ID: "nonexistent"})
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got: %v", err)
		}
	})

	t.Run("delete product", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		repo.Create(ctx, &Product{ID: "prod-1", Name: "Widget"})

		err := repo.Delete(ctx, "prod-1")
		if err != nil {
			t.Fatalf("Delete error: %v", err)
		}

		_, err = repo.GetByID(ctx, "prod-1")
		if !errors.Is(err, ErrNotFound) {
			t.Error("Expected product to be deleted")
		}
	})

	t.Run("delete nonexistent returns ErrNotFound", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		err := repo.Delete(ctx, "nonexistent")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got: %v", err)
		}
	})

	t.Run("list by category", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		repo.Create(ctx, &Product{ID: "1", Name: "Widget", Category: "gadgets"})
		repo.Create(ctx, &Product{ID: "2", Name: "Gizmo", Category: "gadgets"})
		repo.Create(ctx, &Product{ID: "3", Name: "Book", Category: "books"})

		gadgets, err := repo.ListByCategory(ctx, "gadgets")
		if err != nil {
			t.Fatalf("ListByCategory error: %v", err)
		}
		if len(gadgets) != 2 {
			t.Errorf("Expected 2 gadgets, got %d", len(gadgets))
		}

		toys, err := repo.ListByCategory(ctx, "toys")
		if err != nil {
			t.Fatalf("ListByCategory error: %v", err)
		}
		if len(toys) != 0 {
			t.Errorf("Expected 0 toys, got %d", len(toys))
		}
		if toys == nil {
			t.Error("Expected empty slice (not nil) for no results")
		}
	})

	t.Run("stored products are copies (no external mutation)", func(t *testing.T) {
		repo := NewInMemoryProductRepo()
		if repo == nil {
			t.Fatal("NewInMemoryProductRepo returned nil")
		}

		original := &Product{ID: "prod-1", Name: "Widget"}
		repo.Create(ctx, original)

		// Mutate the original — should NOT affect stored copy
		original.Name = "Mutated"

		got, _ := repo.GetByID(ctx, "prod-1")
		if got.Name != "Widget" {
			t.Errorf("Expected stored name='Widget' (not mutated), got %q", got.Name)
		}

		// Mutate the returned value — should NOT affect stored copy
		got.Name = "Also Mutated"
		got2, _ := repo.GetByID(ctx, "prod-1")
		if got2.Name != "Widget" {
			t.Errorf("Expected stored name='Widget' (not mutated), got %q", got2.Name)
		}
	})
}

// =========================================================================
// Test Exercise 3: Parameterized Query Builder
// =========================================================================

func TestBuildSelectQuery(t *testing.T) {
	t.Run("select all", func(t *testing.T) {
		result := BuildSelectQuery("users", nil)
		if result.Query != "SELECT * FROM users" {
			t.Errorf("Expected 'SELECT * FROM users', got %q", result.Query)
		}
		if result.Params != nil {
			t.Errorf("Expected nil params, got %v", result.Params)
		}
	})

	t.Run("select with one condition", func(t *testing.T) {
		result := BuildSelectQuery("users", map[string]interface{}{
			"name": "Alice",
		})
		if result.Query != "SELECT * FROM users WHERE name = $1" {
			t.Errorf("Expected 'SELECT * FROM users WHERE name = $1', got %q", result.Query)
		}
		if len(result.Params) != 1 || result.Params[0] != "Alice" {
			t.Errorf("Expected params=[Alice], got %v", result.Params)
		}
	})

	t.Run("select with multiple conditions (sorted)", func(t *testing.T) {
		result := BuildSelectQuery("products", map[string]interface{}{
			"category": "gadgets",
			"active":   true,
		})
		expected := "SELECT * FROM products WHERE active = $1 AND category = $2"
		if result.Query != expected {
			t.Errorf("Expected %q, got %q", expected, result.Query)
		}
		if len(result.Params) != 2 {
			t.Fatalf("Expected 2 params, got %d", len(result.Params))
		}
		if result.Params[0] != true {
			t.Errorf("Expected first param=true, got %v", result.Params[0])
		}
		if result.Params[1] != "gadgets" {
			t.Errorf("Expected second param='gadgets', got %v", result.Params[1])
		}
	})
}

func TestBuildInsertQuery(t *testing.T) {
	t.Run("insert with values", func(t *testing.T) {
		result := BuildInsertQuery("users", map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
		})
		expected := "INSERT INTO users (email, name) VALUES ($1, $2)"
		if result.Query != expected {
			t.Errorf("Expected %q, got %q", expected, result.Query)
		}
		if len(result.Params) != 2 {
			t.Fatalf("Expected 2 params, got %d", len(result.Params))
		}
		if result.Params[0] != "alice@example.com" {
			t.Errorf("Expected first param='alice@example.com', got %v", result.Params[0])
		}
		if result.Params[1] != "Alice" {
			t.Errorf("Expected second param='Alice', got %v", result.Params[1])
		}
	})

	t.Run("insert single value", func(t *testing.T) {
		result := BuildInsertQuery("logs", map[string]interface{}{
			"message": "hello",
		})
		expected := "INSERT INTO logs (message) VALUES ($1)"
		if result.Query != expected {
			t.Errorf("Expected %q, got %q", expected, result.Query)
		}
	})
}

// =========================================================================
// Test Exercise 4: Transaction Pattern
// =========================================================================

func TestTransactionalRepo(t *testing.T) {
	ctx := context.Background()

	t.Run("successful transaction commits", func(t *testing.T) {
		txRepo := NewTransactionalProductRepo()
		if txRepo == nil {
			t.Fatal("NewTransactionalProductRepo returned nil")
		}

		err := txRepo.WithTransaction(func(repo ProductRepository) error {
			return repo.Create(ctx, &Product{ID: "prod-1", Name: "Widget"})
		})
		if err != nil {
			t.Fatalf("WithTransaction error: %v", err)
		}

		// Verify product exists after commit
		got, err := txRepo.Repo().GetByID(ctx, "prod-1")
		if err != nil {
			t.Fatalf("GetByID after commit: %v", err)
		}
		if got.Name != "Widget" {
			t.Errorf("Expected Name='Widget', got %q", got.Name)
		}
	})

	t.Run("failed transaction rolls back", func(t *testing.T) {
		txRepo := NewTransactionalProductRepo()
		if txRepo == nil {
			t.Fatal("NewTransactionalProductRepo returned nil")
		}

		// Pre-create a product
		txRepo.Repo().Create(ctx, &Product{ID: "existing", Name: "Existing"})

		// Transaction that creates then fails
		err := txRepo.WithTransaction(func(repo ProductRepository) error {
			repo.Create(ctx, &Product{ID: "new-prod", Name: "New"})
			return fmt.Errorf("something went wrong")
		})
		if err == nil {
			t.Fatal("Expected error from failed transaction")
		}

		// Verify the new product was rolled back
		_, err = txRepo.Repo().GetByID(ctx, "new-prod")
		if !errors.Is(err, ErrNotFound) {
			t.Error("Expected new product to be rolled back, but it still exists")
		}

		// Verify pre-existing product still exists
		got, err := txRepo.Repo().GetByID(ctx, "existing")
		if err != nil {
			t.Fatalf("Existing product should still exist: %v", err)
		}
		if got.Name != "Existing" {
			t.Errorf("Expected Name='Existing', got %q", got.Name)
		}
	})

	t.Run("transaction with multiple operations", func(t *testing.T) {
		txRepo := NewTransactionalProductRepo()
		if txRepo == nil {
			t.Fatal("NewTransactionalProductRepo returned nil")
		}

		err := txRepo.WithTransaction(func(repo ProductRepository) error {
			if err := repo.Create(ctx, &Product{ID: "p1", Name: "Product 1"}); err != nil {
				return err
			}
			if err := repo.Create(ctx, &Product{ID: "p2", Name: "Product 2"}); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			t.Fatalf("WithTransaction error: %v", err)
		}

		// Both products should exist
		if _, err := txRepo.Repo().GetByID(ctx, "p1"); err != nil {
			t.Error("Expected p1 to exist after successful transaction")
		}
		if _, err := txRepo.Repo().GetByID(ctx, "p2"); err != nil {
			t.Error("Expected p2 to exist after successful transaction")
		}
	})
}

// =========================================================================
// Test Exercise 5: Query Builder with Filtering
// =========================================================================

func TestFilterProducts(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	boolPtr := func(b bool) *bool { return &b }
	int64Ptr := func(i int64) *int64 { return &i }

	products := []*Product{
		{ID: "1", Name: "Blue Widget", PriceCents: 999, Category: "gadgets", InStock: true},
		{ID: "2", Name: "Red Widget", PriceCents: 1499, Category: "gadgets", InStock: false},
		{ID: "3", Name: "Go Book", PriceCents: 3999, Category: "books", InStock: true},
		{ID: "4", Name: "Cheap Gadget", PriceCents: 299, Category: "gadgets", InStock: true},
		{ID: "5", Name: "Expensive Book", PriceCents: 7999, Category: "books", InStock: false},
	}

	t.Run("no filter returns all", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{})
		if len(result) != 5 {
			t.Errorf("Expected 5 products, got %d", len(result))
		}
	})

	t.Run("filter by category", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{Category: strPtr("gadgets")})
		if len(result) != 3 {
			t.Errorf("Expected 3 gadgets, got %d", len(result))
		}
	})

	t.Run("filter by in stock", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{InStock: boolPtr(true)})
		if len(result) != 3 {
			t.Errorf("Expected 3 in-stock products, got %d", len(result))
		}
	})

	t.Run("filter by price range", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{
			MinPrice: int64Ptr(1000),
			MaxPrice: int64Ptr(5000),
		})
		if len(result) != 2 {
			t.Errorf("Expected 2 products in price range, got %d", len(result))
		}
	})

	t.Run("filter by name contains", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{NameContains: "Widget"})
		if len(result) != 2 {
			t.Errorf("Expected 2 products with 'Widget' in name, got %d", len(result))
		}
	})

	t.Run("combined filters", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{
			Category: strPtr("gadgets"),
			InStock:  boolPtr(true),
		})
		if len(result) != 2 {
			t.Errorf("Expected 2 in-stock gadgets, got %d", len(result))
		}
	})

	t.Run("no matches returns empty slice", func(t *testing.T) {
		result := FilterProducts(products, ProductFilter{Category: strPtr("nonexistent")})
		if len(result) != 0 {
			t.Errorf("Expected 0 products, got %d", len(result))
		}
		if result == nil {
			t.Error("Expected empty slice (not nil)")
		}
	})
}

// =========================================================================
// Test Exercise 6: Nullable Field Handling
// =========================================================================

func TestNullableProfile(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }

	t.Run("profile with all fields to map", func(t *testing.T) {
		p := NullableProfile{
			ID:      "user-1",
			Name:    "Alice",
			Bio:     strPtr("Developer"),
			Website: strPtr("https://alice.dev"),
			Age:     intPtr(30),
		}

		m := ProfileToMap(p)
		if m["id"] != "user-1" {
			t.Errorf("Expected id='user-1', got %v", m["id"])
		}
		if m["bio"] != "Developer" {
			t.Errorf("Expected bio='Developer', got %v", m["bio"])
		}
		if m["age"] != 30 {
			t.Errorf("Expected age=30, got %v", m["age"])
		}
	})

	t.Run("profile with nil fields to map", func(t *testing.T) {
		p := NullableProfile{
			ID:   "user-2",
			Name: "Bob",
			// Bio, Website, Age are nil
		}

		m := ProfileToMap(p)
		if m["bio"] != nil {
			t.Errorf("Expected bio=nil, got %v", m["bio"])
		}
		if m["website"] != nil {
			t.Errorf("Expected website=nil, got %v", m["website"])
		}
		if m["age"] != nil {
			t.Errorf("Expected age=nil, got %v", m["age"])
		}
	})

	t.Run("map with all fields to profile", func(t *testing.T) {
		m := map[string]interface{}{
			"id":      "user-1",
			"name":    "Alice",
			"bio":     "Developer",
			"website": "https://alice.dev",
			"age":     30,
		}

		p := MapToProfile(m)
		if p.ID != "user-1" {
			t.Errorf("Expected ID='user-1', got %q", p.ID)
		}
		if p.Bio == nil || *p.Bio != "Developer" {
			t.Errorf("Expected Bio='Developer', got %v", p.Bio)
		}
		if p.Age == nil || *p.Age != 30 {
			t.Errorf("Expected Age=30, got %v", p.Age)
		}
	})

	t.Run("map with missing fields to profile", func(t *testing.T) {
		m := map[string]interface{}{
			"id":   "user-2",
			"name": "Bob",
		}

		p := MapToProfile(m)
		if p.Bio != nil {
			t.Errorf("Expected Bio=nil, got %v", p.Bio)
		}
		if p.Website != nil {
			t.Errorf("Expected Website=nil, got %v", p.Website)
		}
		if p.Age != nil {
			t.Errorf("Expected Age=nil, got %v", p.Age)
		}
	})

	t.Run("map with nil values to profile", func(t *testing.T) {
		m := map[string]interface{}{
			"id":      "user-3",
			"name":    "Charlie",
			"bio":     nil,
			"website": nil,
			"age":     nil,
		}

		p := MapToProfile(m)
		if p.Bio != nil {
			t.Errorf("Expected Bio=nil for nil map value, got %v", p.Bio)
		}
	})
}

// =========================================================================
// Test Exercise 7: Connection Pool Options
// =========================================================================

func TestDBConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := NewDBConfig()
		if cfg.Host != "localhost" {
			t.Errorf("Expected Host='localhost', got %q", cfg.Host)
		}
		if cfg.Port != 5432 {
			t.Errorf("Expected Port=5432, got %d", cfg.Port)
		}
		if cfg.MaxOpenConns != 25 {
			t.Errorf("Expected MaxOpenConns=25, got %d", cfg.MaxOpenConns)
		}
		if cfg.SSLMode != "disable" {
			t.Errorf("Expected SSLMode='disable', got %q", cfg.SSLMode)
		}
	})

	t.Run("with options", func(t *testing.T) {
		cfg := NewDBConfig(
			WithHost("db.example.com"),
			WithPort(5433),
			WithDatabase("myapp"),
			WithCredentials("admin", "secret"),
			WithSSLMode("require"),
		)
		if cfg.Host != "db.example.com" {
			t.Errorf("Expected Host='db.example.com', got %q", cfg.Host)
		}
		if cfg.Port != 5433 {
			t.Errorf("Expected Port=5433, got %d", cfg.Port)
		}
		if cfg.Database != "myapp" {
			t.Errorf("Expected Database='myapp', got %q", cfg.Database)
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

	t.Run("pool config option", func(t *testing.T) {
		cfg := NewDBConfig(
			WithPoolConfig(50, 10, 10*time.Minute, 5*time.Minute),
		)
		if cfg.MaxOpenConns != 50 {
			t.Errorf("Expected MaxOpenConns=50, got %d", cfg.MaxOpenConns)
		}
		if cfg.MaxIdleConns != 10 {
			t.Errorf("Expected MaxIdleConns=10, got %d", cfg.MaxIdleConns)
		}
		if cfg.ConnMaxLifetime != 10*time.Minute {
			t.Errorf("Expected ConnMaxLifetime=10m, got %v", cfg.ConnMaxLifetime)
		}
	})

	t.Run("connection string", func(t *testing.T) {
		cfg := NewDBConfig(
			WithHost("db.example.com"),
			WithPort(5432),
			WithDatabase("myapp"),
			WithCredentials("admin", "secret"),
			WithSSLMode("require"),
		)
		connStr := cfg.ConnectionString()

		expected := "host=db.example.com port=5432 dbname=myapp user=admin password=secret sslmode=require"
		if connStr != expected {
			t.Errorf("Expected connection string:\n  %q\ngot:\n  %q", expected, connStr)
		}
	})
}

// =========================================================================
// Test Exercise 8: Migration Registry
// =========================================================================

func TestMigrationRegistry(t *testing.T) {
	t.Run("register and run migrations", func(t *testing.T) {
		reg := NewMigrationRegistry()
		if reg == nil {
			t.Fatal("NewMigrationRegistry returned nil")
		}

		var log []string

		reg.Register(MigrationEntry{
			ID:   "001_create_users",
			Up:   func() error { log = append(log, "up:001"); return nil },
			Down: func() error { log = append(log, "down:001"); return nil },
		})
		reg.Register(MigrationEntry{
			ID:   "002_create_products",
			Up:   func() error { log = append(log, "up:002"); return nil },
			Down: func() error { log = append(log, "down:002"); return nil },
		})

		n, err := reg.MigrateUp()
		if err != nil {
			t.Fatalf("MigrateUp error: %v", err)
		}
		if n != 2 {
			t.Errorf("Expected 2 migrations applied, got %d", n)
		}
		if len(log) != 2 || log[0] != "up:001" || log[1] != "up:002" {
			t.Errorf("Expected log=[up:001, up:002], got %v", log)
		}
	})

	t.Run("idempotent migrate up", func(t *testing.T) {
		reg := NewMigrationRegistry()
		if reg == nil {
			t.Fatal("NewMigrationRegistry returned nil")
		}

		count := 0
		reg.Register(MigrationEntry{
			ID:   "001",
			Up:   func() error { count++; return nil },
			Down: func() error { return nil },
		})

		reg.MigrateUp()
		n, _ := reg.MigrateUp() // second call should be no-op
		if n != 0 {
			t.Errorf("Expected 0 migrations on second run, got %d", n)
		}
		if count != 1 {
			t.Errorf("Expected migration to run only once, ran %d times", count)
		}
	})

	t.Run("migrate down rolls back in reverse", func(t *testing.T) {
		reg := NewMigrationRegistry()
		if reg == nil {
			t.Fatal("NewMigrationRegistry returned nil")
		}

		var log []string

		reg.Register(MigrationEntry{
			ID:   "001",
			Up:   func() error { return nil },
			Down: func() error { log = append(log, "down:001"); return nil },
		})
		reg.Register(MigrationEntry{
			ID:   "002",
			Up:   func() error { return nil },
			Down: func() error { log = append(log, "down:002"); return nil },
		})
		reg.Register(MigrationEntry{
			ID:   "003",
			Up:   func() error { return nil },
			Down: func() error { log = append(log, "down:003"); return nil },
		})

		reg.MigrateUp()

		n, err := reg.MigrateDown(2)
		if err != nil {
			t.Fatalf("MigrateDown error: %v", err)
		}
		if n != 2 {
			t.Errorf("Expected 2 rolled back, got %d", n)
		}
		if len(log) != 2 || log[0] != "down:003" || log[1] != "down:002" {
			t.Errorf("Expected [down:003, down:002], got %v", log)
		}

		// 001 should still be applied
		applied := reg.Applied()
		if len(applied) != 1 || applied[0] != "001" {
			t.Errorf("Expected [001] still applied, got %v", applied)
		}
	})

	t.Run("pending returns unapplied migrations", func(t *testing.T) {
		reg := NewMigrationRegistry()
		if reg == nil {
			t.Fatal("NewMigrationRegistry returned nil")
		}

		reg.Register(MigrationEntry{ID: "001", Up: func() error { return nil }, Down: func() error { return nil }})
		reg.Register(MigrationEntry{ID: "002", Up: func() error { return nil }, Down: func() error { return nil }})
		reg.Register(MigrationEntry{ID: "003", Up: func() error { return nil }, Down: func() error { return nil }})

		pending := reg.Pending()
		if len(pending) != 3 {
			t.Errorf("Expected 3 pending, got %d", len(pending))
		}

		reg.MigrateUp()
		pending = reg.Pending()
		if len(pending) != 0 {
			t.Errorf("Expected 0 pending after migrate up, got %d", len(pending))
		}
	})

	t.Run("migration failure stops execution", func(t *testing.T) {
		reg := NewMigrationRegistry()
		if reg == nil {
			t.Fatal("NewMigrationRegistry returned nil")
		}

		reg.Register(MigrationEntry{
			ID:   "001",
			Up:   func() error { return nil },
			Down: func() error { return nil },
		})
		reg.Register(MigrationEntry{
			ID:   "002",
			Up:   func() error { return fmt.Errorf("migration failed") },
			Down: func() error { return nil },
		})
		reg.Register(MigrationEntry{
			ID:   "003",
			Up:   func() error { return nil },
			Down: func() error { return nil },
		})

		n, err := reg.MigrateUp()
		if err == nil {
			t.Error("Expected error from failed migration")
		}
		if n != 1 {
			t.Errorf("Expected 1 successful migration before failure, got %d", n)
		}

		applied := reg.Applied()
		if len(applied) != 1 || applied[0] != "001" {
			t.Errorf("Expected only [001] applied, got %v", applied)
		}
	})

	t.Run("duplicate registration returns error", func(t *testing.T) {
		reg := NewMigrationRegistry()
		if reg == nil {
			t.Fatal("NewMigrationRegistry returned nil")
		}

		err := reg.Register(MigrationEntry{ID: "001", Up: func() error { return nil }, Down: func() error { return nil }})
		if err != nil {
			t.Fatalf("First register should succeed: %v", err)
		}

		err = reg.Register(MigrationEntry{ID: "001", Up: func() error { return nil }, Down: func() error { return nil }})
		if err == nil {
			t.Error("Expected error for duplicate migration ID")
		}
	})
}
