package database

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 EXERCISES: Database Access
=============================================================================

 These exercises teach database patterns using interfaces and in-memory
 implementations — no actual database required!

 Run the tests with:

   make test 19

 Tip: Run a single test at a time while working:

   go test -v -run TestProductRepository ./19-database-access/

=============================================================================
*/

// =========================================================================
// Exercise 1: Define a Repository Interface
// =========================================================================

// Product represents a product in an e-commerce system.
type Product struct {
	ID          string
	Name        string
	Description string
	PriceCents  int64
	Category    string
	InStock     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProductRepository defines the interface for product data access.
// Define methods for:
//   - Create(ctx context.Context, product *Product) error
//   - GetByID(ctx context.Context, id string) (*Product, error)
//   - List(ctx context.Context, limit, offset int) ([]*Product, error)
//   - Update(ctx context.Context, product *Product) error
//   - Delete(ctx context.Context, id string) error
//   - ListByCategory(ctx context.Context, category string) ([]*Product, error)
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, limit, offset int) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	ListByCategory(ctx context.Context, category string) ([]*Product, error)
}

// =========================================================================
// Exercise 2: Implement an In-Memory Repository
// =========================================================================

// InMemoryProductRepo implements ProductRepository with an in-memory map.
// It should be safe for concurrent use (use sync.RWMutex).
//
// Requirements:
//   - Create: fail with ErrAlreadyExists if ID exists, fail with ErrInvalidInput if ID is empty
//   - GetByID: fail with ErrNotFound if not found
//   - List: return paginated results (apply offset first, then limit). Return empty slice (not nil) if no results.
//   - Update: fail with ErrNotFound if not found, update the UpdatedAt timestamp
//   - Delete: fail with ErrNotFound if not found
//   - ListByCategory: return all products matching the category. Return empty slice if none.
//   - Store copies of products to prevent external mutation
//   - Set CreatedAt and UpdatedAt on Create
type InMemoryProductRepo struct {
	mu       sync.RWMutex
	products map[string]*Product
}

// NewInMemoryProductRepo creates a new empty in-memory product repository.
func NewInMemoryProductRepo() *InMemoryProductRepo {
	// YOUR CODE HERE
	return nil
}

// Implement all ProductRepository methods on InMemoryProductRepo.

// Create adds a new product to the repository.
// YOUR CODE HERE — implement validation and storage logic
func (r *InMemoryProductRepo) Create(_ context.Context, product *Product) error {
	// YOUR CODE HERE
	return nil
}

// GetByID retrieves a product by its ID.
// YOUR CODE HERE — return ErrNotFound if not found
func (r *InMemoryProductRepo) GetByID(_ context.Context, id string) (*Product, error) {
	// YOUR CODE HERE
	return nil, fmt.Errorf("%w: product with ID %s", ErrNotFound, id)
}

// List returns a paginated list of products.
// YOUR CODE HERE — apply offset then limit, return empty slice (not nil) if no results
func (r *InMemoryProductRepo) List(_ context.Context, limit, offset int) ([]*Product, error) {
	// YOUR CODE HERE
	return []*Product{}, nil
}

// Update modifies an existing product.
// YOUR CODE HERE — return ErrNotFound if not found, set UpdatedAt
func (r *InMemoryProductRepo) Update(_ context.Context, product *Product) error {
	// YOUR CODE HERE
	return fmt.Errorf("%w: product with ID %s", ErrNotFound, product.ID)
}

// Delete removes a product by its ID.
// YOUR CODE HERE — return ErrNotFound if not found
func (r *InMemoryProductRepo) Delete(_ context.Context, id string) error {
	// YOUR CODE HERE
	return fmt.Errorf("%w: product with ID %s", ErrNotFound, id)
}

// ListByCategory returns all products in the given category.
// YOUR CODE HERE — return empty slice (not nil) if no matches
func (r *InMemoryProductRepo) ListByCategory(_ context.Context, category string) ([]*Product, error) {
	// YOUR CODE HERE
	_ = category
	return []*Product{}, nil
}

// =========================================================================
// Exercise 3: Parameterized Query Builder
// =========================================================================

// QueryResult holds a SQL query string and its parameters.
// This exercise teaches you to think about parameterized queries
// without needing a real database.
type QueryResult struct {
	Query  string
	Params []interface{}
}

// BuildSelectQuery builds a SELECT query for the given table with optional
// WHERE conditions. Each condition is a column name and value pair.
//
// Use positional placeholders ($1, $2, etc.) for PostgreSQL style.
//
// Examples:
//
//	BuildSelectQuery("users", nil)
//	  → QueryResult{Query: "SELECT * FROM users", Params: nil}
//
//	BuildSelectQuery("users", map[string]interface{}{"name": "Alice", "active": true})
//	  → QueryResult{Query: "SELECT * FROM users WHERE active = $1 AND name = $2", Params: [true, "Alice"]}
//
// Note: Sort condition keys alphabetically for deterministic output.
func BuildSelectQuery(table string, conditions map[string]interface{}) QueryResult {
	// YOUR CODE HERE
	return QueryResult{}
}

// BuildInsertQuery builds an INSERT query for the given table and column values.
//
// Example:
//
//	BuildInsertQuery("users", map[string]interface{}{"name": "Alice", "email": "a@b.com"})
//	  → QueryResult{Query: "INSERT INTO users (email, name) VALUES ($1, $2)", Params: ["a@b.com", "Alice"]}
//
// Note: Sort column names alphabetically for deterministic output.
func BuildInsertQuery(table string, values map[string]interface{}) QueryResult {
	// YOUR CODE HERE
	return QueryResult{}
}

// =========================================================================
// Exercise 4: Transaction Pattern
// =========================================================================

// TxFunc is a function that runs within a transaction context.
// If it returns an error, the transaction should be rolled back.
// If it returns nil, the transaction should be committed.
type TxFunc func(repo ProductRepository) error

// TransactionalProductRepo wraps a ProductRepository and adds
// transaction support using a snapshot/restore pattern.
//
// WithTransaction should:
//  1. Take a snapshot of the current state
//  2. Execute the function
//  3. If the function returns an error, restore the snapshot (rollback)
//  4. If the function returns nil, keep the changes (commit)
type TransactionalProductRepo struct {
	repo *InMemoryProductRepo
}

// NewTransactionalProductRepo creates a new transactional wrapper.
func NewTransactionalProductRepo() *TransactionalProductRepo {
	// YOUR CODE HERE
	return nil
}

// Repo returns the underlying repository for direct operations.
func (t *TransactionalProductRepo) Repo() *InMemoryProductRepo {
	// YOUR CODE HERE
	return nil
}

// WithTransaction executes fn within a transaction. If fn returns an error,
// all changes made during fn are rolled back. If fn returns nil, changes
// are kept.
//
// Implementation hint: Before calling fn, copy the products map (snapshot).
// If fn returns an error, restore the map from the snapshot.
func (t *TransactionalProductRepo) WithTransaction(fn TxFunc) error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 5: Query Builder with Filtering
// =========================================================================

// ProductFilter defines filter criteria for querying products.
type ProductFilter struct {
	Category     *string // nil means no filter
	InStock      *bool   // nil means no filter
	MinPrice     *int64  // nil means no minimum
	MaxPrice     *int64  // nil means no maximum
	NameContains string  // empty means no filter
}

// FilterProducts takes a slice of products and a filter, returning only
// products that match ALL specified criteria.
//
// A nil filter field means "don't filter on this criterion."
// An empty NameContains means "don't filter on name."
//
// This simulates what a WHERE clause would do in SQL, but in Go.
func FilterProducts(products []*Product, filter ProductFilter) []*Product {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 6: Nullable Field Handling
// =========================================================================

// NullableProfile represents a user profile where some fields may be null.
// This mimics how you'd handle nullable database columns in Go.
type NullableProfile struct {
	ID      string
	Name    string
	Bio     *string // nullable
	Website *string // nullable
	Age     *int    // nullable
}

// ProfileToMap converts a NullableProfile to a map[string]interface{},
// representing how you might build a JSON response or database row.
//
// Rules:
//   - Always include "id" and "name"
//   - For pointer fields: include the key with the dereferenced value
//     if non-nil, or include the key with nil if the pointer is nil
func ProfileToMap(p NullableProfile) map[string]interface{} {
	// YOUR CODE HERE
	return nil
}

// MapToProfile converts a map[string]interface{} back to a NullableProfile.
// This simulates reading from a database row where some columns may be NULL.
//
// The map uses these keys: "id", "name", "bio", "website", "age"
// Missing or nil values for pointer fields should result in nil pointers.
// The "id" and "name" fields default to "" if missing.
// The "age" field in the map will be an int (not *int) if present.
func MapToProfile(m map[string]interface{}) NullableProfile {
	// YOUR CODE HERE
	return NullableProfile{}
}

// =========================================================================
// Exercise 7: Connection Pool Options (Functional Options Pattern)
// =========================================================================

// DBConfig represents database connection configuration.
type DBConfig struct {
	Host            string
	Port            int
	Database        string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DBOption is a function that modifies a DBConfig.
type DBOption func(*DBConfig)

// NewDBConfig creates a DBConfig with sensible defaults and applies
// any provided options.
//
// Defaults:
//   - Host: "localhost"
//   - Port: 5432
//   - Database: "app"
//   - User: "postgres"
//   - Password: ""
//   - SSLMode: "disable"
//   - MaxOpenConns: 25
//   - MaxIdleConns: 5
//   - ConnMaxLifetime: 5 minutes
//   - ConnMaxIdleTime: 2 minutes
//
// Apply each option function to the config before returning.
func NewDBConfig(opts ...DBOption) DBConfig {
	// YOUR CODE HERE
	return DBConfig{}
}

// Implement these option functions:

// WithHost returns a DBOption that sets the host.
func WithHost(host string) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithPort returns a DBOption that sets the port.
func WithPort(port int) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithDatabase returns a DBOption that sets the database name.
func WithDatabase(db string) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithCredentials returns a DBOption that sets user and password.
func WithCredentials(user, password string) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithSSLMode returns a DBOption that sets the SSL mode.
func WithSSLMode(mode string) DBOption {
	// YOUR CODE HERE
	return nil
}

// WithPoolConfig returns a DBOption that sets all pool parameters.
func WithPoolConfig(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) DBOption {
	// YOUR CODE HERE
	return nil
}

// ConnectionString returns a PostgreSQL-style connection string.
// Format: "host=HOST port=PORT dbname=DB user=USER password=PASS sslmode=MODE"
func (c DBConfig) ConnectionString() string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 8: Migration Registry
// =========================================================================

// MigrationFunc represents a migration function (up or down).
type MigrationFunc func() error

// MigrationEntry represents a single migration with up and down functions.
type MigrationEntry struct {
	ID   string
	Up   MigrationFunc
	Down MigrationFunc
}

// MigrationRegistry manages database migrations.
// It tracks registered migrations and which ones have been applied.
type MigrationRegistry struct {
	migrations []MigrationEntry
	applied    map[string]bool
	mu         sync.Mutex
}

// NewMigrationRegistry creates a new empty migration registry.
func NewMigrationRegistry() *MigrationRegistry {
	// YOUR CODE HERE
	return nil
}

// Register adds a migration to the registry.
// Migrations should be registered in order (by ID).
// Return an error if a migration with the same ID already exists.
func (r *MigrationRegistry) Register(entry MigrationEntry) error {
	// YOUR CODE HERE
	return nil
}

// MigrateUp runs all unapplied migrations in order.
// If any migration fails, stop and return the error.
// Successfully applied migrations before the failure should remain applied.
// Return the number of migrations that were applied.
func (r *MigrationRegistry) MigrateUp() (int, error) {
	// YOUR CODE HERE
	return 0, nil
}

// MigrateDown rolls back the last N applied migrations in reverse order.
// If n is 0 or negative, roll back all applied migrations.
// Return the number of migrations that were rolled back.
func (r *MigrationRegistry) MigrateDown(n int) (int, error) {
	// YOUR CODE HERE
	return 0, nil
}

// Applied returns the IDs of all applied migrations in order.
func (r *MigrationRegistry) Applied() []string {
	// YOUR CODE HERE
	return nil
}

// Pending returns the IDs of all unapplied migrations in order.
func (r *MigrationRegistry) Pending() []string {
	// YOUR CODE HERE
	return nil
}

// These are used to suppress "unused import" errors in stubs.
var _ = context.Background
var _ = fmt.Sprintf
var _ = strings.Join
