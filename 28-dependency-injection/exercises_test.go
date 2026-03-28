package dependencyinjection

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Shared test helpers / fakes
// ---------------------------------------------------------------------------

// fakeOrderStore is a test double for OrderStore.
type fakeOrderStore struct {
	orders    map[string]fakeOrder
	nextID    int
	insertErr error
}

type fakeOrder struct {
	item string
	qty  int
}

func newFakeOrderStore() *fakeOrderStore {
	return &fakeOrderStore{orders: make(map[string]fakeOrder)}
}

func (f *fakeOrderStore) InsertOrder(_ context.Context, item string, qty int) (string, error) {
	if f.insertErr != nil {
		return "", f.insertErr
	}
	f.nextID++
	id := fmt.Sprintf("ORD-%d", f.nextID)
	f.orders[id] = fakeOrder{item: item, qty: qty}
	return id, nil
}

func (f *fakeOrderStore) GetOrder(_ context.Context, orderID string) (string, int, error) {
	o, ok := f.orders[orderID]
	if !ok {
		return "", 0, fmt.Errorf("order %q not found", orderID)
	}
	return o.item, o.qty, nil
}

// fakeNotifier records notifications.
type fakeNotifier struct {
	notifications []string
	err           error
}

func (f *fakeNotifier) NotifyOrderPlaced(_ context.Context, orderID string) error {
	if f.err != nil {
		return f.err
	}
	f.notifications = append(f.notifications, orderID)
	return nil
}

// ---------------------------------------------------------------------------
// Exercise 1: OrderService with constructor injection
// ---------------------------------------------------------------------------

func TestOrderService_PlaceOrder(t *testing.T) {
	store := newFakeOrderStore()
	notifier := &fakeNotifier{}

	svc := NewOrderService(store, notifier)
	if svc == nil {
		t.Fatal("NewOrderService returned nil. Store dependencies as struct fields and return the struct.")
	}

	t.Run("successful order", func(t *testing.T) {
		orderID, err := svc.PlaceOrder(context.Background(), "widget", 5)
		if err != nil {
			t.Fatalf("PlaceOrder failed: %v. Call store.InsertOrder and notifier.NotifyOrderPlaced.", err)
		}
		if orderID == "" {
			t.Error("PlaceOrder returned empty orderID. Return the ID from store.InsertOrder.")
		}
		if len(notifier.notifications) == 0 {
			t.Error("No notification sent. Call notifier.NotifyOrderPlaced after inserting.")
		}
	})

	t.Run("retrieve placed order", func(t *testing.T) {
		orderID, _ := svc.PlaceOrder(context.Background(), "gadget", 3)
		item, qty, err := svc.GetOrder(context.Background(), orderID)
		if err != nil {
			t.Fatalf("GetOrder failed: %v. Delegate to store.GetOrder.", err)
		}
		if item != "gadget" || qty != 3 {
			t.Errorf("GetOrder = (%q, %d), want (%q, %d).", item, qty, "gadget", 3)
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 2: Interface definitions (compile-time check)
// ---------------------------------------------------------------------------

func TestInterfaceDefinitions(t *testing.T) {
	// These are compile-time checks that the interfaces are defined correctly.
	// If this compiles, the interfaces have the right method signatures.

	t.Run("Cache interface", func(t *testing.T) {
		var c Cache = &testCacheImpl{}
		ctx := context.Background()
		_, _, _ = c.Get(ctx, "key")
		_ = c.Set(ctx, "key", "value", time.Second)
	})

	t.Run("MessageQueue interface", func(t *testing.T) {
		var q MessageQueue = &testMQImpl{}
		ctx := context.Background()
		_ = q.Publish(ctx, "topic", "msg")
		_, _ = q.Subscribe(ctx, "topic")
	})
}

// Minimal implementations to verify interface signatures
type testCacheImpl struct{}

func (c *testCacheImpl) Get(_ context.Context, _ string) (string, bool, error) { return "", false, nil }
func (c *testCacheImpl) Set(_ context.Context, _, _ string, _ time.Duration) error { return nil }

type testMQImpl struct{}

func (q *testMQImpl) Publish(_ context.Context, _, _ string) error { return nil }
func (q *testMQImpl) Subscribe(_ context.Context, _ string) (<-chan string, error) {
	return make(chan string), nil
}

// ---------------------------------------------------------------------------
// Exercise 3: Functional Options for DatabaseClient
// ---------------------------------------------------------------------------

func TestDatabaseClientDefaults(t *testing.T) {
	c := NewDatabaseClient()
	if c == nil {
		t.Fatal("NewDatabaseClient() returned nil. Create a DatabaseClient with default values.")
	}
	if c.Host != "localhost" {
		t.Errorf("Default Host = %q, want %q.", c.Host, "localhost")
	}
	if c.Port != 5432 {
		t.Errorf("Default Port = %d, want 5432.", c.Port)
	}
	if c.Database != "app" {
		t.Errorf("Default Database = %q, want %q.", c.Database, "app")
	}
	if c.MaxConnections != 10 {
		t.Errorf("Default MaxConnections = %d, want 10.", c.MaxConnections)
	}
	if c.ConnectTimeout != 5*time.Second {
		t.Errorf("Default ConnectTimeout = %v, want 5s.", c.ConnectTimeout)
	}
	if c.ReadOnly {
		t.Error("Default ReadOnly = true, want false.")
	}
}

func TestDatabaseClientOptions(t *testing.T) {
	c := NewDatabaseClient(
		WithDBHost("db.production.com"),
		WithDBPort(5433),
		WithDBName("mydb"),
		WithMaxConnections(50),
		WithConnectTimeout(10*time.Second),
		WithReadOnly(),
	)

	if c.Host != "db.production.com" {
		t.Errorf("Host = %q, want %q. WithDBHost should set the host.", c.Host, "db.production.com")
	}
	if c.Port != 5433 {
		t.Errorf("Port = %d, want 5433. WithDBPort should set the port.", c.Port)
	}
	if c.Database != "mydb" {
		t.Errorf("Database = %q, want %q. WithDBName should set the database name.", c.Database, "mydb")
	}
	if c.MaxConnections != 50 {
		t.Errorf("MaxConnections = %d, want 50.", c.MaxConnections)
	}
	if c.ConnectTimeout != 10*time.Second {
		t.Errorf("ConnectTimeout = %v, want 10s.", c.ConnectTimeout)
	}
	if !c.ReadOnly {
		t.Error("ReadOnly = false, want true. WithReadOnly should set ReadOnly to true.")
	}
}

func TestDatabaseClientDSN(t *testing.T) {
	c := NewDatabaseClient(WithDBHost("myhost"), WithDBPort(5433), WithDBName("mydb"))
	dsn := c.DSN()
	want := "myhost:5433/mydb"
	if dsn != want {
		t.Errorf("DSN() = %q, want %q. Format as host:port/database.", dsn, want)
	}
}

// ---------------------------------------------------------------------------
// Exercise 4: ProductService with DI
// ---------------------------------------------------------------------------

func TestProductService_GetProduct(t *testing.T) {
	repo := NewMockProductRepo(&Product{ID: "p1", Name: "Widget", Price: 9.99})
	logger := &MockLogger{}
	svc := NewProductService(repo, logger)
	if svc == nil {
		t.Fatal("NewProductService returned nil. Store repo and logger as struct fields.")
	}

	t.Run("found", func(t *testing.T) {
		p, err := svc.GetProduct(context.Background(), "p1")
		if err != nil {
			t.Fatalf("GetProduct failed: %v. Delegate to repo.FindByID.", err)
		}
		if p == nil || p.Name != "Widget" {
			t.Errorf("GetProduct = %v, want Widget.", p)
		}
		if len(logger.InfoMessages) == 0 {
			t.Error("No info messages logged. Log the lookup with logger.Info.")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetProduct(context.Background(), "missing")
		if err == nil {
			t.Error("GetProduct for missing ID should return error.")
		}
		if len(logger.ErrorMessages) == 0 {
			t.Error("No error messages logged for missing product. Log errors with logger.Error.")
		}
	})
}

func TestProductService_ListProducts(t *testing.T) {
	repo := NewMockProductRepo(
		&Product{ID: "p1", Name: "Widget", Price: 9.99},
		&Product{ID: "p2", Name: "Gadget", Price: 19.99},
	)
	logger := &MockLogger{}
	svc := NewProductService(repo, logger)

	products, err := svc.ListProducts(context.Background())
	if err != nil {
		t.Fatalf("ListProducts failed: %v.", err)
	}
	if len(products) != 2 {
		t.Errorf("ListProducts returned %d products, want 2.", len(products))
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	repo := NewMockProductRepo()
	logger := &MockLogger{}
	svc := NewProductService(repo, logger)

	p := &Product{ID: "p1", Name: "New Thing", Price: 5.00}
	err := svc.CreateProduct(context.Background(), p)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v.", err)
	}

	// Verify it was saved
	found, err := repo.FindByID(context.Background(), "p1")
	if err != nil {
		t.Fatalf("Product not found in repo after save: %v.", err)
	}
	if found.Name != "New Thing" {
		t.Errorf("Saved product name = %q, want %q.", found.Name, "New Thing")
	}
}

// ---------------------------------------------------------------------------
// Exercise 5: MockProductRepo
// ---------------------------------------------------------------------------

func TestMockProductRepo(t *testing.T) {
	t.Run("initialized with products", func(t *testing.T) {
		repo := NewMockProductRepo(
			&Product{ID: "1", Name: "A", Price: 1.0},
			&Product{ID: "2", Name: "B", Price: 2.0},
		)
		if repo == nil {
			t.Fatal("NewMockProductRepo returned nil. Initialize the Products map.")
		}
		p, err := repo.FindByID(context.Background(), "1")
		if err != nil {
			t.Fatalf("FindByID failed: %v. Check the Products map.", err)
		}
		if p.Name != "A" {
			t.Errorf("FindByID = %q, want %q.", p.Name, "A")
		}
	})

	t.Run("save and find", func(t *testing.T) {
		repo := NewMockProductRepo()
		_ = repo.Save(context.Background(), &Product{ID: "x", Name: "X", Price: 10.0})
		p, err := repo.FindByID(context.Background(), "x")
		if err != nil {
			t.Fatalf("FindByID after Save failed: %v.", err)
		}
		if p.Name != "X" {
			t.Errorf("Got %q, want %q.", p.Name, "X")
		}
	})

	t.Run("find all", func(t *testing.T) {
		repo := NewMockProductRepo(
			&Product{ID: "1", Name: "A", Price: 1.0},
			&Product{ID: "2", Name: "B", Price: 2.0},
		)
		all, err := repo.FindAll(context.Background())
		if err != nil {
			t.Fatalf("FindAll failed: %v.", err)
		}
		if len(all) != 2 {
			t.Errorf("FindAll returned %d, want 2.", len(all))
		}
	})

	t.Run("save error", func(t *testing.T) {
		repo := NewMockProductRepo()
		repo.SaveErr = errors.New("disk full")
		err := repo.Save(context.Background(), &Product{ID: "x", Name: "X"})
		if err == nil {
			t.Error("Save should return SaveErr when set.")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 6: BuildApplication factory
// ---------------------------------------------------------------------------

func TestBuildApplication(t *testing.T) {
	productRepo := NewMockProductRepo(&Product{ID: "p1", Name: "Widget", Price: 9.99})
	orderStore := newFakeOrderStore()
	notifier := &fakeNotifier{}
	logger := &MockLogger{}

	app := BuildApplication(productRepo, orderStore, notifier, logger)
	if app == nil {
		t.Fatal("BuildApplication returned nil. Create an Application with wired services.")
	}
	if app.ProductService == nil {
		t.Fatal("Application.ProductService is nil. Create it with NewProductService.")
	}
	if app.OrderService == nil {
		t.Fatal("Application.OrderService is nil. Create it with NewOrderService.")
	}

	// Test that services work through the application
	t.Run("product lookup through app", func(t *testing.T) {
		p, err := app.ProductService.GetProduct(context.Background(), "p1")
		if err != nil {
			t.Fatalf("Failed: %v.", err)
		}
		if p.Name != "Widget" {
			t.Errorf("Got %q, want %q.", p.Name, "Widget")
		}
	})

	t.Run("order placement through app", func(t *testing.T) {
		id, err := app.OrderService.PlaceOrder(context.Background(), "widget", 2)
		if err != nil {
			t.Fatalf("Failed: %v.", err)
		}
		if id == "" {
			t.Error("Empty order ID.")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 7: BuildRoutes
// ---------------------------------------------------------------------------

func TestBuildRoutes(t *testing.T) {
	repo := NewMockProductRepo(
		&Product{ID: "p1", Name: "Widget", Price: 9.99},
		&Product{ID: "p2", Name: "Gadget", Price: 19.99},
	)
	logger := &MockLogger{}
	svc := NewProductService(repo, logger)

	routes := BuildRoutes(svc)
	if routes == nil {
		t.Fatal("BuildRoutes returned nil. Create a map with route handlers.")
	}

	t.Run("GET /products", func(t *testing.T) {
		handler, ok := routes["GET /products"]
		if !ok {
			t.Fatal("Missing route 'GET /products'. Add it to the map.")
		}
		result, err := handler(context.Background(), nil)
		if err != nil {
			t.Fatalf("Handler failed: %v.", err)
		}
		// Should contain both product names
		if !strings.Contains(result, "Widget") || !strings.Contains(result, "Gadget") {
			t.Errorf("ListProducts result = %q, should contain Widget and Gadget.", result)
		}
	})

	t.Run("GET /products/:id", func(t *testing.T) {
		handler, ok := routes["GET /products/:id"]
		if !ok {
			t.Fatal("Missing route 'GET /products/:id'. Add it to the map.")
		}
		result, err := handler(context.Background(), map[string]string{"id": "p1"})
		if err != nil {
			t.Fatalf("Handler failed: %v.", err)
		}
		if result != "Widget" {
			t.Errorf("GetProduct result = %q, want %q.", result, "Widget")
		}
	})

	t.Run("POST /products", func(t *testing.T) {
		handler, ok := routes["POST /products"]
		if !ok {
			t.Fatal("Missing route 'POST /products'. Add it to the map.")
		}
		result, err := handler(context.Background(), map[string]string{
			"id": "p3", "name": "Doohickey", "price": "29.99",
		})
		if err != nil {
			t.Fatalf("Handler failed: %v.", err)
		}
		if result != "created" {
			t.Errorf("CreateProduct result = %q, want %q.", result, "created")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 8: Ports and Adapters
// ---------------------------------------------------------------------------

func TestConsoleAdapter(t *testing.T) {
	adapter := &ConsoleAdapter{}
	err := adapter.Send(context.Background(), "alice@example.com", "Hello", "World")
	if err != nil {
		t.Fatalf("Send failed: %v.", err)
	}
	if len(adapter.SentMessages) != 1 {
		t.Fatalf("Expected 1 sent message, got %d. Append to SentMessages in Send().", len(adapter.SentMessages))
	}
	msg := adapter.SentMessages[0]
	if msg.To != "alice@example.com" {
		t.Errorf("To = %q, want %q.", msg.To, "alice@example.com")
	}
	if msg.Subject != "Hello" {
		t.Errorf("Subject = %q, want %q.", msg.Subject, "Hello")
	}
	if msg.Body != "World" {
		t.Errorf("Body = %q, want %q.", msg.Body, "World")
	}
}

func TestNotificationService(t *testing.T) {
	adapter := &ConsoleAdapter{}
	svc := NewNotificationService(adapter)
	if svc == nil {
		t.Fatal("NewNotificationService returned nil. Store the adapter as a struct field.")
	}

	err := svc.NotifyUser(context.Background(), "alice@example.com", "Alice")
	if err != nil {
		t.Fatalf("NotifyUser failed: %v.", err)
	}

	if len(adapter.SentMessages) != 1 {
		t.Fatalf("Expected 1 message, got %d.", len(adapter.SentMessages))
	}

	msg := adapter.SentMessages[0]
	if msg.To != "alice@example.com" {
		t.Errorf("To = %q, want %q.", msg.To, "alice@example.com")
	}
	if msg.Subject != "Welcome, Alice!" {
		t.Errorf("Subject = %q, want %q. Format: \"Welcome, <name>!\"", msg.Subject, "Welcome, Alice!")
	}
	if msg.Body != "Hello Alice, welcome to our service." {
		t.Errorf("Body = %q, want %q.", msg.Body, "Hello Alice, welcome to our service.")
	}
}
