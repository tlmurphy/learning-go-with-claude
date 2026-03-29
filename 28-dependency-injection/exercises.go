package dependencyinjection

import (
	"context"
	"fmt"
	"time"
)

/*
=============================================================================
 EXERCISES: Dependency Injection
=============================================================================

 Work through these exercises to practice Go's DI patterns. Run tests with:

   make test 28

 Each exercise builds toward a complete, properly-wired service.

=============================================================================
*/

// =========================================================================
// Exercise 1: Refactor a Tightly Coupled Service
// =========================================================================
//
// TIGHTLY COUPLED CODE (the "before"):
//
//   type OrderService struct{}
//
//   func (s *OrderService) PlaceOrder(item string, qty int) (string, error) {
//       // Directly calls database
//       db := connectToDatabase()  // hard-coded!
//       orderID := db.Insert(item, qty)
//       // Directly sends email
//       smtp := connectToSMTP()    // hard-coded!
//       smtp.Send("order " + orderID + " placed")
//       return orderID, nil
//   }
//
// YOUR TASK: Create an OrderService that accepts its dependencies through
// a constructor. Define two interfaces: OrderStore and OrderNotifier.

// OrderStore is the interface for order persistence.
type OrderStore interface {
	InsertOrder(ctx context.Context, item string, qty int) (string, error)
	GetOrder(ctx context.Context, orderID string) (string, int, error)
}

// OrderNotifier is the interface for sending order notifications.
type OrderNotifier interface {
	NotifyOrderPlaced(ctx context.Context, orderID string) error
}

// OrderService orchestrates order operations using injected dependencies.
type OrderService struct {
	// YOUR CODE HERE — store the dependencies
	store    OrderStore
	notifier OrderNotifier
}

// NewOrderService creates an OrderService with its dependencies.
func NewOrderService(store OrderStore, notifier OrderNotifier) *OrderService {
	// YOUR CODE HERE
	return nil
}

// PlaceOrder creates an order and notifies. Returns the order ID.
func (s *OrderService) PlaceOrder(ctx context.Context, item string, qty int) (string, error) {
	// YOUR CODE HERE
	// 1. Use s.store.InsertOrder to create the order
	// 2. Use s.notifier.NotifyOrderPlaced to send notification
	// 3. Return the order ID
	return "", nil
}

// GetOrder retrieves an order by ID.
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (string, int, error) {
	// YOUR CODE HERE
	// Delegate to s.store.GetOrder
	return "", 0, nil
}

// =========================================================================
// Exercise 2: Define Interfaces for External Dependencies
// =========================================================================
//
// Define interfaces that abstract away external services. These interfaces
// should be small and focused (1-2 methods each).

// Cache abstracts a caching layer (like Redis).
// Methods:
//   - Get(ctx, key) returns (value string, found bool, err error)
//   - Set(ctx, key, value, ttl) returns error
type Cache interface {
	// YOUR CODE HERE
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

// MessageQueue abstracts a message broker (like RabbitMQ or SQS).
// Methods:
//   - Publish(ctx, topic, message) returns error
//   - Subscribe(ctx, topic) returns (<-chan string, error)
type MessageQueue interface {
	// YOUR CODE HERE
	Publish(ctx context.Context, topic string, message string) error
	Subscribe(ctx context.Context, topic string) (<-chan string, error)
}

// =========================================================================
// Exercise 3: Functional Options Pattern
// =========================================================================
//
// Implement the functional options pattern for a DatabaseClient.

// DatabaseClient represents a database connection with configurable options.
type DatabaseClient struct {
	Host           string
	Port           int
	Database       string
	MaxConnections int
	ConnectTimeout time.Duration
	ReadOnly       bool
}

// DBOption is a functional option for configuring DatabaseClient.
type DBOption func(*DatabaseClient)

// WithDBHost sets the database host.
func WithDBHost(host string) DBOption {
	// YOUR CODE HERE
	return func(c *DatabaseClient) {}
}

// WithDBPort sets the database port.
func WithDBPort(port int) DBOption {
	// YOUR CODE HERE
	return func(c *DatabaseClient) {}
}

// WithDBName sets the database name.
func WithDBName(name string) DBOption {
	// YOUR CODE HERE
	return func(c *DatabaseClient) {}
}

// WithMaxConnections sets the max connection pool size.
func WithMaxConnections(n int) DBOption {
	// YOUR CODE HERE
	return func(c *DatabaseClient) {}
}

// WithConnectTimeout sets the connection timeout.
func WithConnectTimeout(d time.Duration) DBOption {
	// YOUR CODE HERE
	return func(c *DatabaseClient) {}
}

// WithReadOnly sets the client to read-only mode.
func WithReadOnly() DBOption {
	// YOUR CODE HERE
	return func(c *DatabaseClient) {}
}

// NewDatabaseClient creates a DatabaseClient with sensible defaults:
//
//	Host: "localhost", Port: 5432, Database: "app",
//	MaxConnections: 10, ConnectTimeout: 5s, ReadOnly: false
//
// Options override the defaults.
func NewDatabaseClient(opts ...DBOption) *DatabaseClient {
	// YOUR CODE HERE
	return &DatabaseClient{}
}

// DSN returns the connection string: "host:port/database"
func (c *DatabaseClient) DSN() string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 4: Service Layer with DI
// =========================================================================
//
// Build a simple service layer: ProductService depends on ProductRepo
// and a Logger. The handler layer is simulated through a function.

// Product represents a product in a catalog.
type Product struct {
	ID    string
	Name  string
	Price float64
}

// ProductRepo is the interface for product data access.
type ProductRepo interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Save(ctx context.Context, p *Product) error
}

// ServiceLogger is a simple logging interface.
type ServiceLogger interface {
	Info(msg string)
	Error(msg string)
}

// ProductService handles product business logic.
type ProductService struct {
	// YOUR CODE HERE — store repo and logger
	repo   ProductRepo
	logger ServiceLogger
}

// NewProductService creates a ProductService with its dependencies.
func NewProductService(repo ProductRepo, logger ServiceLogger) *ProductService {
	// YOUR CODE HERE
	return nil
}

// GetProduct retrieves a product by ID.
// It should log the lookup and return the product or an error.
func (s *ProductService) GetProduct(ctx context.Context, id string) (*Product, error) {
	// YOUR CODE HERE
	// 1. Log: "looking up product: <id>"
	// 2. Call s.repo.FindByID
	// 3. If error, log: "product not found: <id>" and return error
	// 4. Return product
	return nil, nil
}

// ListProducts retrieves all products.
func (s *ProductService) ListProducts(ctx context.Context) ([]*Product, error) {
	// YOUR CODE HERE
	// 1. Log: "listing all products"
	// 2. Call s.repo.FindAll
	// 3. Return results
	return nil, nil
}

// CreateProduct saves a new product.
func (s *ProductService) CreateProduct(ctx context.Context, p *Product) error {
	// YOUR CODE HERE
	// 1. Log: "creating product: <name>"
	// 2. Call s.repo.Save
	// 3. Return error if any
	return nil
}

// =========================================================================
// Exercise 5: Test with Mock/Stub Repository
// =========================================================================
//
// Implement a mock ProductRepo that the tests can use to verify
// ProductService behavior without a real database.

// MockProductRepo is an in-memory implementation of ProductRepo for testing.
type MockProductRepo struct {
	Products map[string]*Product
	SaveErr  error // Set this to simulate save errors
}

// NewMockProductRepo creates a MockProductRepo with optional initial products.
func NewMockProductRepo(products ...*Product) *MockProductRepo {
	// YOUR CODE HERE
	return &MockProductRepo{}
}

// FindByID looks up a product in the in-memory map.
func (r *MockProductRepo) FindByID(_ context.Context, id string) (*Product, error) {
	// YOUR CODE HERE
	return nil, fmt.Errorf("product %q not found", id)
}

// FindAll returns all products in the mock store.
func (r *MockProductRepo) FindAll(_ context.Context) ([]*Product, error) {
	// YOUR CODE HERE
	return nil, nil
}

// Save stores a product in the mock store. Returns SaveErr if set.
func (r *MockProductRepo) Save(_ context.Context, p *Product) error {
	// YOUR CODE HERE
	return nil
}

// MockLogger captures log messages for testing.
type MockLogger struct {
	InfoMessages  []string
	ErrorMessages []string
}

func (l *MockLogger) Info(msg string) {
	l.InfoMessages = append(l.InfoMessages, msg)
}

func (l *MockLogger) Error(msg string) {
	l.ErrorMessages = append(l.ErrorMessages, msg)
}

// =========================================================================
// Exercise 6: Factory Function for Dependency Graph
// =========================================================================
//
// In a real app, main() creates all dependencies. This exercise simulates
// that with a factory function that builds the entire dependency graph.

// AppConfig holds configuration for the application.
type AppConfig struct {
	DBHost   string
	DBPort   int
	DBName   string
	LogLevel string
}

// Application holds all the top-level services, fully wired.
type Application struct {
	ProductService *ProductService
	OrderService   *OrderService
}

// BuildApplication creates a fully wired Application from config and
// implementations. This is your "composition root."
//
// It takes a ProductRepo and OrderStore/OrderNotifier because in a real
// app, main() would create those from the config.
func BuildApplication(
	productRepo ProductRepo,
	orderStore OrderStore,
	orderNotifier OrderNotifier,
	logger ServiceLogger,
) *Application {
	// YOUR CODE HERE
	// 1. Create ProductService with productRepo and logger
	// 2. Create OrderService with orderStore and orderNotifier
	// 3. Return Application with both services
	return nil
}

// =========================================================================
// Exercise 7: Wire-Up Function for HTTP Server
// =========================================================================
//
// Create a function that returns a map of routes to handler functions,
// demonstrating how DI wiring connects to HTTP routing.
// (We use a simple map instead of http.ServeMux to avoid import complexity.)

// RouteHandler is a simplified HTTP handler function.
type RouteHandler func(ctx context.Context, params map[string]string) (string, error)

// BuildRoutes creates a route map wired to the given ProductService.
// Routes:
//
//	"GET /products"     — calls svc.ListProducts, returns product names joined by comma
//	"GET /products/:id" — calls svc.GetProduct with params["id"], returns product name
//	"POST /products"    — calls svc.CreateProduct with product from params, returns "created"
func BuildRoutes(svc *ProductService) map[string]RouteHandler {
	// YOUR CODE HERE
	// Create a map with 3 entries, each calling the appropriate service method
	return nil
}

// =========================================================================
// Exercise 8: Ports and Adapters Pattern
// =========================================================================
//
// Implement a notification system using the ports and adapters pattern.
// The "port" (interface) is defined by the core domain.
// "Adapters" are concrete implementations.

// NotificationPort is the port defined by the core domain.
// Any notification adapter must satisfy this interface.
type NotificationPort interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

// ConsoleAdapter is an adapter that "sends" notifications by storing them
// in memory (useful for testing and demos).
type ConsoleAdapter struct {
	SentMessages []SentMessage
}

// SentMessage records a message that was "sent".
type SentMessage struct {
	To      string
	Subject string
	Body    string
}

// Send records the message (doesn't actually send anything).
func (a *ConsoleAdapter) Send(_ context.Context, to, subject, body string) error {
	// YOUR CODE HERE
	// Append a SentMessage to a.SentMessages
	return nil
}

// NotificationService is the core domain service that uses the port.
type NotificationService struct {
	// YOUR CODE HERE — store the NotificationPort
	adapter NotificationPort
}

// NewNotificationService creates a NotificationService with its port.
func NewNotificationService(adapter NotificationPort) *NotificationService {
	// YOUR CODE HERE
	return nil
}

// NotifyUser sends a welcome notification to a user.
func (s *NotificationService) NotifyUser(ctx context.Context, email, name string) error {
	// YOUR CODE HERE
	// Send a notification with:
	//   to: email
	//   subject: "Welcome, <name>!"
	//   body: "Hello <name>, welcome to our service."
	_ = fmt.Sprintf
	return nil
}
