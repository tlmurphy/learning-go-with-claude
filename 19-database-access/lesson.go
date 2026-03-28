package database

/*
=============================================================================
 Module 19: Database Access
=============================================================================

 Every non-trivial web service talks to a database. Go's database/sql
 package provides a universal interface for SQL databases — it's
 driver-agnostic, connection-pooled, and context-aware out of the box.

 However, database/sql has more sharp edges than most Go packages. It's
 designed for correctness and generality, which means it's easy to
 accidentally leak connections, miss error checks, or misconfigure the
 pool.

 This module teaches the patterns without requiring an actual database.
 We'll define interfaces and in-memory implementations that mirror what
 real database code looks like. When you add a real database later, you'll
 swap in a real implementation without changing your business logic.

 WHY THIS MATTERS FOR WEB SERVICES:
 - Connection pool misconfiguration is a top cause of production outages
 - SQL injection is still a top-10 vulnerability (use parameterized queries!)
 - Forgetting to close rows leaks connections (your app slowly dies)
 - The repository pattern decouples business logic from storage details

=============================================================================
*/

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// -------------------------------------------------------------------------
// database/sql Overview
// -------------------------------------------------------------------------

/*
 The database/sql package has three key concepts:

 1. sql.DB — NOT a single connection! It's a CONNECTION POOL.
    - sql.Open() creates the pool but doesn't actually connect
    - Always call db.Ping() or db.PingContext() to verify connectivity
    - It's safe to share across goroutines (it manages concurrency)
    - Create ONE sql.DB per database for your entire application

 2. sql.Tx — A transaction. Wraps multiple operations atomically.
    - Always defer tx.Rollback() — it's a no-op after Commit
    - Context cancellation rolls back automatically

 3. sql.Rows / sql.Row — Query results.
    - ALWAYS close Rows (defer rows.Close())
    - Always check rows.Err() after iteration
    - sql.Row (singular) is for QueryRow — doesn't need explicit Close

 Here's what typical database code looks like (pseudocode since we
 don't have a real DB):

   db, err := sql.Open("postgres", connString)
   if err != nil { return err }
   defer db.Close()

   // IMPORTANT: Open doesn't connect! Ping does.
   if err := db.PingContext(ctx); err != nil { return err }

   // Configure the pool
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
*/

// -------------------------------------------------------------------------
// Connection Pool Configuration
// -------------------------------------------------------------------------

/*
 Connection pool settings are critical in production. Wrong settings
 cause either:
   - Too few connections → requests queue up, latency spikes
   - Too many connections → overwhelm the database server
   - Stale connections → random "connection reset" errors

 Key settings:

 MaxOpenConns (default: unlimited!)
   Maximum number of open connections to the database.
   ALWAYS set this. The default of 0 (unlimited) will exhaust your
   database's connection limit under load.
   Rule of thumb: start with 25, tune based on load testing.

 MaxIdleConns (default: 2)
   Maximum number of idle connections in the pool.
   Set this to something reasonable (5-10). Too low means constant
   reconnection overhead. Too high wastes database resources.

 ConnMaxLifetime (default: unlimited)
   Maximum time a connection can be reused.
   Set this to ~5 minutes. It ensures connections are periodically
   refreshed, which helps with database failovers and DNS changes.

 ConnMaxIdleTime (default: unlimited)
   Maximum time a connection can sit idle before being closed.
   Set this to prevent idle connections from going stale.
*/

// PoolConfig holds connection pool configuration.
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultPoolConfig returns sensible defaults for a connection pool.
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}
}

// -------------------------------------------------------------------------
// Query vs QueryRow vs Exec
// -------------------------------------------------------------------------

/*
 database/sql has three main query methods. Using the wrong one is
 a common source of connection leaks:

 db.Query(query, args...) → *sql.Rows
   Use for SELECT that returns multiple rows.
   MUST close the rows when done: defer rows.Close()
   Forgetting to close leaks a connection!

 db.QueryRow(query, args...) → *sql.Row
   Use for SELECT that returns at most one row.
   No need to close — Row is consumed by Scan().
   If no rows match, Scan returns sql.ErrNoRows.

 db.Exec(query, args...) → sql.Result
   Use for INSERT, UPDATE, DELETE (no result rows).
   Returns RowsAffected() and LastInsertId().
   No rows to close — no leak risk.

 All three have Context variants (QueryContext, QueryRowContext, ExecContext)
 which you should ALWAYS use in web services. The context carries the
 request deadline — if the client disconnects, the query cancels.
*/

// -------------------------------------------------------------------------
// Scanning Results and Handling NULL
// -------------------------------------------------------------------------

/*
 sql.Rows.Scan() reads columns into Go variables. It does type conversion
 but has important rules:

   rows.Scan(&id, &name, &email)  // scan into addressable variables

 For NULL values, you have two options:

 1. sql.NullString, sql.NullInt64, sql.NullBool, etc.
    These are structs with a Valid field:
      var name sql.NullString
      rows.Scan(&name)
      if name.Valid { fmt.Println(name.String) }

 2. Pointer types (*string, *int, etc.)
    NULL becomes nil:
      var name *string
      rows.Scan(&name)
      if name != nil { fmt.Println(*name) }

 Pointer types are simpler and more idiomatic. The sql.Null* types
 exist for historical reasons and because they implement Scanner/Valuer
 interfaces for custom types.

 ALWAYS check the error from Scan. A common mistake is to ignore it and
 use partially-scanned data.
*/

// -------------------------------------------------------------------------
// Transactions
// -------------------------------------------------------------------------

/*
 Transactions group multiple operations into an atomic unit.
 The pattern in Go is:

   tx, err := db.BeginTx(ctx, nil)
   if err != nil { return err }
   defer tx.Rollback()  // no-op after commit, safety net on error

   // Do work within the transaction
   _, err = tx.ExecContext(ctx, "INSERT ...", args...)
   if err != nil { return err }

   _, err = tx.ExecContext(ctx, "UPDATE ...", args...)
   if err != nil { return err }

   return tx.Commit()  // commit if everything succeeded

 The key insight: defer tx.Rollback() is always safe. If Commit() has
 already been called, Rollback() does nothing. If any error occurs before
 Commit(), the deferred Rollback() cleans up.

 NEVER start a transaction and forget to either Commit or Rollback.
 That holds a connection open and eventually exhausts your pool.
*/

// -------------------------------------------------------------------------
// SQL Injection Prevention
// -------------------------------------------------------------------------

/*
 NEVER build SQL queries with string concatenation or fmt.Sprintf:

   // DANGEROUS — SQL injection vulnerability!
   query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)

   // If name is: '; DROP TABLE users; --
   // The query becomes: SELECT * FROM users WHERE name = ''; DROP TABLE users; --'

 ALWAYS use parameterized queries:

   // SAFE — the database driver escapes the parameter
   rows, err := db.Query("SELECT * FROM users WHERE name = $1", name)

 The placeholder syntax varies by database:
   PostgreSQL: $1, $2, $3
   MySQL:      ?, ?, ?
   SQLite:     ?, ?, ? or $1, $2, $3

 Parameterized queries are not just safer — they're also faster because
 the database can cache the query plan.
*/

// -------------------------------------------------------------------------
// The Repository Pattern
// -------------------------------------------------------------------------

/*
 The repository pattern is the most common way to organize database code
 in Go. It provides a clean interface between your business logic and
 your data storage:

   Business Logic → Repository Interface → Database Implementation

 Benefits:
   1. Business logic doesn't know about SQL — it calls repo.GetUser(id)
   2. Easy to test — swap in a mock/in-memory implementation
   3. Easy to switch databases — implement a new repository
   4. Consistent error handling — repository translates DB errors

 The interface defines WHAT operations are available.
 The implementation defines HOW they're performed.
*/

// User represents a user in our system.
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository defines the interface for user data access.
// This is the contract that any storage backend must fulfill.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// Common errors that repositories should return.
// Using sentinel errors lets callers check with errors.Is().
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

// -------------------------------------------------------------------------
// In-Memory Repository Implementation
// -------------------------------------------------------------------------

/*
 An in-memory implementation is valuable for:
   1. Testing — fast, no external dependencies
   2. Prototyping — get your API working before adding a database
   3. Reference — shows the expected behavior of each method

 Note how it implements the same interface as a database-backed repo would.
 Your HTTP handlers shouldn't know or care which implementation they're using.
*/

// InMemoryUserRepo implements UserRepository using an in-memory map.
type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserRepo creates an empty in-memory user repository.
func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepo) Create(_ context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == "" {
		return fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("%w: user with ID %s", ErrAlreadyExists, user.ID)
	}

	// Store a copy to prevent external mutation
	stored := *user
	stored.CreatedAt = time.Now()
	stored.UpdatedAt = stored.CreatedAt
	r.users[user.ID] = &stored
	return nil
}

func (r *InMemoryUserRepo) GetByID(_ context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("%w: user with ID %s", ErrNotFound, id)
	}

	// Return a copy to prevent external mutation
	result := *user
	return &result, nil
}

func (r *InMemoryUserRepo) GetByEmail(_ context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			result := *user
			return &result, nil
		}
	}
	return nil, fmt.Errorf("%w: user with email %s", ErrNotFound, email)
}

func (r *InMemoryUserRepo) List(_ context.Context, limit, offset int) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all users (map iteration order is random, but that's fine for demo)
	all := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		u := *user
		all = append(all, &u)
	}

	// Apply offset
	if offset >= len(all) {
		return []*User{}, nil
	}
	all = all[offset:]

	// Apply limit
	if limit > 0 && limit < len(all) {
		all = all[:limit]
	}

	return all, nil
}

func (r *InMemoryUserRepo) Update(_ context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("%w: user with ID %s", ErrNotFound, user.ID)
	}

	stored := *user
	stored.UpdatedAt = time.Now()
	r.users[user.ID] = &stored
	return nil
}

func (r *InMemoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("%w: user with ID %s", ErrNotFound, id)
	}

	delete(r.users, id)
	return nil
}

// -------------------------------------------------------------------------
// Transaction Pattern
// -------------------------------------------------------------------------

/*
 Even with an in-memory store, the transaction pattern is worth learning.
 The key idea: wrap a set of operations in a function, and let the
 transaction manager handle commit/rollback.

 In production, this function would call tx.Commit() on success and
 tx.Rollback() on error. Here we simulate it with a "snapshot and
 restore" approach.

 The pattern:

   err := store.WithTransaction(ctx, func(tx TxContext) error {
       if err := tx.DoThing1(); err != nil { return err }
       if err := tx.DoThing2(); err != nil { return err }
       return nil  // success → commit
   })
   // if err != nil, the transaction was rolled back
*/

// -------------------------------------------------------------------------
// Prepared Statements
// -------------------------------------------------------------------------

/*
 Prepared statements are pre-compiled SQL queries that can be executed
 multiple times with different parameters. They offer:

 1. Performance — the database parses and plans the query once
 2. Safety — parameters are always properly escaped
 3. Clarity — separates SQL structure from data

 In Go:
   stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
   if err != nil { return err }
   defer stmt.Close()

   // Use many times
   row1 := stmt.QueryRowContext(ctx, "user-1")
   row2 := stmt.QueryRowContext(ctx, "user-2")

 CAVEAT: Prepared statements have connection affinity. The statement is
 bound to a specific database connection. If that connection goes away,
 database/sql transparently re-prepares it, but this adds latency.
 For simple queries, the overhead of preparing may not be worth it.
*/

// -------------------------------------------------------------------------
// Migration Strategies
// -------------------------------------------------------------------------

/*
 Database migrations track schema changes over time. The typical approach:

 1. Each migration has an ID (usually a timestamp) and up/down functions
 2. A migrations table tracks which have been applied
 3. On startup, run any unapplied migrations in order

 Common tools:
   - golang-migrate/migrate (most popular)
   - pressly/goose
   - atlas (by Ariga)

 Key principles:
   - Migrations should be idempotent when possible
   - Always write both up AND down migrations
   - Never modify a migration that's been applied to production
   - Test migrations against a copy of production data
   - Consider backward compatibility (can old code run with new schema?)

 We'll implement a simple migration registry to understand the pattern.
*/

// Migration represents a database schema migration.
type Migration struct {
	ID   string
	Up   func() error
	Down func() error
}

// -------------------------------------------------------------------------
// Common Gotchas Summary
// -------------------------------------------------------------------------

/*
 1. FORGETTING TO CLOSE ROWS
    rows, err := db.Query(...)
    // If you forget defer rows.Close(), the connection leaks.
    // Your pool eventually exhausts and all queries block.

 2. NOT CHECKING rows.Err()
    for rows.Next() { rows.Scan(...) }
    // Always check: if err := rows.Err(); err != nil { ... }
    // rows.Next() can stop early due to an error, not just EOF.

 3. USING db.Query FOR NON-SELECT QUERIES
    db.Query("DELETE FROM users WHERE id = $1", id)
    // This LEAKS a connection because you never close the rows!
    // Use db.Exec for INSERT, UPDATE, DELETE.

 4. sql.Open DOESN'T CONNECT
    db, err := sql.Open("postgres", connString)
    // err is only non-nil if the driver name is wrong.
    // Call db.Ping() to actually test the connection.

 5. UNLIMITED MaxOpenConns
    // The default is 0 (unlimited). Under load, this overwhelms the DB.
    // ALWAYS set db.SetMaxOpenConns(25) or similar.

 6. SCANNING INTO WRONG TYPES
    // columns must be scanned in order, with matching types.
    // Mismatches cause runtime errors, not compile-time errors.
*/
