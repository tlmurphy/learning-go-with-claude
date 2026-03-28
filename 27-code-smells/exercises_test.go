package codesmells

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// Exercise 1: KeyValueStore
// ---------------------------------------------------------------------------

func TestKeyValueStore(t *testing.T) {
	t.Run("new store is empty", func(t *testing.T) {
		store := NewKeyValueStore()
		if store == nil {
			t.Fatal("NewKeyValueStore() returned nil. Initialize the store with a map.")
		}
		_, err := store.Get("missing")
		if err == nil {
			t.Error("Get on empty store should return an error.")
		}
	})

	t.Run("set and get", func(t *testing.T) {
		store := NewKeyValueStore()
		store.Set("name", "Alice")
		val, err := store.Get("name")
		if err != nil {
			t.Errorf("Get after Set returned error: %v. Set should store the value in the map.", err)
		}
		if val != "Alice" {
			t.Errorf("Get(\"name\") = %q, want %q.", val, "Alice")
		}
	})

	t.Run("overwrite", func(t *testing.T) {
		store := NewKeyValueStore()
		store.Set("key", "v1")
		store.Set("key", "v2")
		val, _ := store.Get("key")
		if val != "v2" {
			t.Errorf("After overwrite, Get = %q, want %q.", val, "v2")
		}
	})

	t.Run("delete existing", func(t *testing.T) {
		store := NewKeyValueStore()
		store.Set("key", "value")
		err := store.Delete("key")
		if err != nil {
			t.Errorf("Delete existing key returned error: %v.", err)
		}
		_, err = store.Get("key")
		if err == nil {
			t.Error("Get after Delete should return error.")
		}
	})

	t.Run("delete missing", func(t *testing.T) {
		store := NewKeyValueStore()
		err := store.Delete("missing")
		if err == nil {
			t.Error("Delete on missing key should return an error.")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 2: SafeProducer (goroutine leak fix)
// ---------------------------------------------------------------------------

func TestSafeProducer(t *testing.T) {
	t.Run("reads all values", func(t *testing.T) {
		ctx := context.Background()
		data := []int{1, 2, 3, 4, 5}
		ch := SafeProducer(ctx, data)

		var got []int
		for v := range ch {
			got = append(got, v)
		}

		if len(got) != len(data) {
			t.Errorf("Got %d values, want %d. The goroutine should send all values and close the channel.", len(got), len(data))
		}
		for i, v := range got {
			if v != data[i] {
				t.Errorf("Value %d = %d, want %d.", i, v, data[i])
			}
		}
	})

	t.Run("stops on cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		data := make([]int, 1000)
		for i := range data {
			data[i] = i
		}

		ch := SafeProducer(ctx, data)

		// Read only a few values
		got := <-ch
		if got != 0 {
			t.Errorf("First value = %d, want 0.", got)
		}
		got = <-ch
		if got != 1 {
			t.Errorf("Second value = %d, want 1.", got)
		}

		// Cancel the context — goroutine should stop
		cancel()

		// Drain any remaining buffered values
		remaining := 0
		for range ch {
			remaining++
		}

		// The goroutine should have stopped well before sending all 1000 values
		if remaining > 998 {
			t.Error("SafeProducer sent all values despite context cancellation. Use select with ctx.Done() to stop early.")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 3: Config (fix stuttering)
// ---------------------------------------------------------------------------

func TestConfig(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		cfg := NewConfig("smtp.example.com", 587, "noreply@example.com", "secret")
		if cfg.Host != "smtp.example.com" {
			t.Errorf("Host = %q, want %q.", cfg.Host, "smtp.example.com")
		}
		if cfg.Port != 587 {
			t.Errorf("Port = %d, want 587.", cfg.Port)
		}
		if cfg.From != "noreply@example.com" {
			t.Errorf("From = %q, want %q.", cfg.From, "noreply@example.com")
		}
		if cfg.Password != "secret" {
			t.Errorf("Password = %q, want %q.", cfg.Password, "secret")
		}
	})

	t.Run("address", func(t *testing.T) {
		cfg := NewConfig("smtp.example.com", 587, "x", "x")
		addr := cfg.Address()
		want := "smtp.example.com:587"
		if addr != want {
			t.Errorf("Address() = %q, want %q. Use fmt.Sprintf to format host:port.", addr, want)
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 4: Cache (refactor global state)
// ---------------------------------------------------------------------------

func TestCache(t *testing.T) {
	t.Run("get missing", func(t *testing.T) {
		c := NewCache()
		if c == nil {
			t.Fatal("NewCache() returned nil. Initialize with a map.")
		}
		_, ok := c.Get("missing")
		if ok {
			t.Error("Get on empty cache returned ok=true.")
		}
	})

	t.Run("set and get", func(t *testing.T) {
		c := NewCache()
		c.Set("key", "value")
		val, ok := c.Get("key")
		if !ok {
			t.Error("Get after Set returned ok=false. Use the map to store values.")
		}
		if val != "value" {
			t.Errorf("Get = %q, want %q.", val, "value")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		c := NewCache()
		var wg sync.WaitGroup
		// Concurrent writes
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c.Set(fmt.Sprintf("key-%d", i), fmt.Sprintf("val-%d", i))
			}(i)
		}
		wg.Wait()

		// Verify all values
		for i := 0; i < 100; i++ {
			val, ok := c.Get(fmt.Sprintf("key-%d", i))
			if !ok {
				t.Errorf("Missing key-%d after concurrent writes. Ensure you're using sync.RWMutex.", i)
			}
			want := fmt.Sprintf("val-%d", i)
			if val != want {
				t.Errorf("key-%d = %q, want %q.", i, val, want)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 5: Error Handling
// ---------------------------------------------------------------------------

func TestErrNotFound(t *testing.T) {
	msg := ErrNotFound.Error()
	if msg != "not found" {
		t.Errorf("ErrNotFound.Error() = %q, want %q. Error strings should be lowercase with no punctuation.",
			msg, "not found")
	}
}

func TestFindItem(t *testing.T) {
	store := map[string]string{"apple": "fruit", "carrot": "vegetable"}

	t.Run("found", func(t *testing.T) {
		val, err := FindItem(store, "apple")
		if err != nil {
			t.Errorf("FindItem for existing key returned error: %v.", err)
		}
		if val != "fruit" {
			t.Errorf("FindItem = %q, want %q.", val, "fruit")
		}
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		_, err := FindItem(store, "pizza")
		if err == nil {
			t.Fatal("FindItem for missing key should return an error.")
		}
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("FindItem error = %v, want ErrNotFound. Return the ErrNotFound sentinel error.", err)
		}
	})
}

func TestProcessItem(t *testing.T) {
	store := map[string]string{"apple": "fruit"}

	t.Run("found item is uppercased", func(t *testing.T) {
		val, err := ProcessItem(store, "apple")
		if err != nil {
			t.Errorf("ProcessItem returned error: %v.", err)
		}
		if val != "FRUIT" {
			t.Errorf("ProcessItem = %q, want %q. Use strings.ToUpper on the found value.", val, "FRUIT")
		}
	})

	t.Run("not found returns DEFAULT", func(t *testing.T) {
		val, err := ProcessItem(store, "pizza")
		if err != nil {
			t.Errorf("ProcessItem should not return error for not-found, got: %v.", err)
		}
		if val != "DEFAULT" {
			t.Errorf("ProcessItem for missing key = %q, want %q. Use errors.Is(err, ErrNotFound) to check, then return \"DEFAULT\".",
				val, "DEFAULT")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 6: ServerConfig (simplify over-engineering)
// ---------------------------------------------------------------------------

func TestServerConfig(t *testing.T) {
	t.Run("zero value uses defaults", func(t *testing.T) {
		cfg := ServerConfig{}
		addr := cfg.Addr()
		want := "localhost:8080"
		if addr != want {
			t.Errorf("Zero-value ServerConfig.Addr() = %q, want %q. Use DefaultHost and DefaultPort when fields are zero.",
				addr, want)
		}
		timeout := cfg.EffectiveTimeout()
		if timeout != DefaultTimeout {
			t.Errorf("Zero-value EffectiveTimeout() = %d, want %d.", timeout, DefaultTimeout)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		cfg := ServerConfig{Host: "0.0.0.0", Port: 9090, Timeout: 60}
		addr := cfg.Addr()
		want := "0.0.0.0:9090"
		if addr != want {
			t.Errorf("Addr() = %q, want %q.", addr, want)
		}
		if cfg.EffectiveTimeout() != 60 {
			t.Errorf("EffectiveTimeout() = %d, want 60.", cfg.EffectiveTimeout())
		}
	})

	t.Run("partial custom", func(t *testing.T) {
		cfg := ServerConfig{Port: 3000}
		addr := cfg.Addr()
		want := "localhost:3000"
		if addr != want {
			t.Errorf("Addr() with only Port set = %q, want %q.", addr, want)
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 7: RequestHandler (fix context misuse)
// ---------------------------------------------------------------------------

type mockLogger struct {
	messages []string
}

func (l *mockLogger) Log(msg string) {
	l.messages = append(l.messages, msg)
}

type mockDB struct {
	users map[int]string
}

func (d *mockDB) GetUser(_ context.Context, id int) string {
	return d.users[id]
}

func TestRequestHandler(t *testing.T) {
	logger := &mockLogger{}
	db := &mockDB{users: map[int]string{1: "Alice", 2: "Bob"}}

	handler := NewRequestHandler(logger, db)
	if handler == nil {
		t.Fatal("NewRequestHandler returned nil. Store logger and db as struct fields.")
	}

	t.Run("handles request with context", func(t *testing.T) {
		ctx := WithRequestID(context.Background(), "req-123")
		result := handler.HandleRequest(ctx, 1)
		if result != "Alice" {
			t.Errorf("HandleRequest = %q, want %q. Use h.db.GetUser to fetch the user.", result, "Alice")
		}
		if len(logger.messages) == 0 {
			t.Error("No log messages recorded. Use h.logger.Log to log the action.")
		}
		// Check that the log includes the request ID
		lastMsg := logger.messages[len(logger.messages)-1]
		if !strings.Contains(lastMsg, "req-123") {
			t.Errorf("Log message %q should contain request ID %q. Use GetRequestID(ctx) to extract it.",
				lastMsg, "req-123")
		}
	})

	t.Run("different user", func(t *testing.T) {
		ctx := WithRequestID(context.Background(), "req-456")
		result := handler.HandleRequest(ctx, 2)
		if result != "Bob" {
			t.Errorf("HandleRequest = %q, want %q.", result, "Bob")
		}
	})
}

// ---------------------------------------------------------------------------
// Exercise 8: Decomposed Order Processing
// ---------------------------------------------------------------------------

func TestCalculateSubtotal(t *testing.T) {
	prices := map[string]float64{"apple": 1.50, "banana": 0.75, "coffee": 4.00}

	t.Run("valid items", func(t *testing.T) {
		total, err := CalculateSubtotal([]string{"apple", "banana"}, prices)
		if err != nil {
			t.Errorf("Unexpected error: %v.", err)
		}
		if total != 2.25 {
			t.Errorf("Subtotal = %.2f, want 2.25. Sum the prices from the map.", total)
		}
	})

	t.Run("unknown item", func(t *testing.T) {
		_, err := CalculateSubtotal([]string{"apple", "steak"}, prices)
		if err == nil {
			t.Error("Should return error for unknown item.")
		}
	})

	t.Run("empty", func(t *testing.T) {
		total, err := CalculateSubtotal([]string{}, prices)
		if err != nil {
			t.Errorf("Unexpected error: %v.", err)
		}
		if total != 0 {
			t.Errorf("Subtotal = %.2f, want 0.", total)
		}
	})
}

func TestApplyDiscount(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		discount float64
		want     float64
	}{
		{"10% off", 100.0, 0.1, 90.0},
		{"no discount", 100.0, 0.0, 100.0},
		{"50% off", 200.0, 0.5, 100.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyDiscount(tt.amount, tt.discount)
			if got != tt.want {
				t.Errorf("ApplyDiscount(%.2f, %.2f) = %.2f, want %.2f. Multiply amount by (1 - discount).",
					tt.amount, tt.discount, got, tt.want)
			}
		})
	}
}

func TestCalculateTax(t *testing.T) {
	got := CalculateTax(100.0, 0.08)
	if got != 8.0 {
		t.Errorf("CalculateTax(100, 0.08) = %.2f, want 8.00. Multiply amount by taxRate.", got)
	}
}

func TestFormatReceipt(t *testing.T) {
	got := FormatReceipt(3, 10.00, 0.80, 10.80)
	want := "Items: 3, Subtotal: $10.00, Tax: $0.80, Total: $10.80"
	if got != want {
		t.Errorf("FormatReceipt = %q, want %q.", got, want)
	}
}

func TestComposeReceipt(t *testing.T) {
	prices := map[string]float64{"apple": 1.00, "banana": 2.00, "coffee": 5.00}

	t.Run("full order", func(t *testing.T) {
		receipt, err := ComposeReceipt(
			[]string{"apple", "banana", "coffee"},
			prices, 0.10, 0.0,
		)
		if err != nil {
			t.Errorf("Unexpected error: %v.", err)
		}
		// Subtotal: 8.00, Tax: 0.80, Total: 8.80
		want := "Items: 3, Subtotal: $8.00, Tax: $0.80, Total: $8.80"
		if receipt != want {
			t.Errorf("ComposeReceipt = %q, want %q. Chain CalculateSubtotal -> ApplyDiscount -> CalculateTax -> FormatReceipt.",
				receipt, want)
		}
	})

	t.Run("with discount", func(t *testing.T) {
		receipt, err := ComposeReceipt(
			[]string{"coffee", "coffee"},
			prices, 0.10, 0.20, // 20% discount
		)
		if err != nil {
			t.Errorf("Unexpected error: %v.", err)
		}
		// Subtotal: 10.00, after 20% discount: 8.00, Tax: 0.80, Total: 8.80
		want := "Items: 2, Subtotal: $8.00, Tax: $0.80, Total: $8.80"
		if receipt != want {
			t.Errorf("ComposeReceipt with discount = %q, want %q.", receipt, want)
		}
	})

	t.Run("unknown item", func(t *testing.T) {
		_, err := ComposeReceipt([]string{"steak"}, prices, 0.10, 0.0)
		if err == nil {
			t.Error("Should return error for unknown item.")
		}
	})
}
