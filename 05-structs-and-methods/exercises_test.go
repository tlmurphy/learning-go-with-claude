package structs

import (
	"fmt"
	"testing"
)

// =========================================================================
// Exercise 1 Tests: User Struct
// =========================================================================

func TestNewUser(t *testing.T) {
	u := NewUser("Alice", "Smith", "alice@example.com", 30)
	if u == nil {
		t.Fatal("NewUser returned nil — make sure to return a pointer to User")
	}
	if u.FirstName != "Alice" {
		t.Errorf("FirstName = %q, want %q", u.FirstName, "Alice")
	}
	if u.LastName != "Smith" {
		t.Errorf("LastName = %q, want %q", u.LastName, "Smith")
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", u.Email, "alice@example.com")
	}
	if u.Age != 30 {
		t.Errorf("Age = %d, want %d", u.Age, 30)
	}
}

func TestUserFullName(t *testing.T) {
	tests := []struct {
		first, last string
		want        string
	}{
		{"Alice", "Smith", "Alice Smith"},
		{"Bob", "Jones", "Bob Jones"},
		{"", "Solo", " Solo"},
	}

	for _, tt := range tests {
		t.Run(tt.first+" "+tt.last, func(t *testing.T) {
			u := User{FirstName: tt.first, LastName: tt.last}
			got := u.FullName()
			if got != tt.want {
				t.Errorf("FullName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUserIsAdult(t *testing.T) {
	tests := []struct {
		age  int
		want bool
	}{
		{17, false},
		{18, true},
		{21, true},
		{0, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("age_%d", tt.age), func(t *testing.T) {
			u := User{Age: tt.age}
			got := u.IsAdult()
			if got != tt.want {
				t.Errorf("IsAdult() with age %d = %v, want %v", tt.age, got, tt.want)
			}
		})
	}
}

func TestUserUpdateEmail(t *testing.T) {
	u := NewUser("Alice", "Smith", "old@example.com", 30)
	if u == nil {
		t.Fatal("NewUser returned nil")
	}
	u.UpdateEmail("new@example.com")
	if u.Email != "new@example.com" {
		t.Errorf("After UpdateEmail, Email = %q, want %q — did you use a pointer receiver?",
			u.Email, "new@example.com")
	}
}

// =========================================================================
// Exercise 2 Tests: Rectangle
// =========================================================================

func TestNewRectangle(t *testing.T) {
	tests := []struct {
		name          string
		width, height float64
		wantW, wantH  float64
	}{
		{"positive", 5, 3, 5, 3},
		{"negative width", -5, 3, 5, 3},
		{"negative height", 5, -3, 5, 3},
		{"both negative", -5, -3, 5, 3},
		{"zero", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRectangle(tt.width, tt.height)
			if r.Width != tt.wantW {
				t.Errorf("Width = %f, want %f", r.Width, tt.wantW)
			}
			if r.Height != tt.wantH {
				t.Errorf("Height = %f, want %f", r.Height, tt.wantH)
			}
		})
	}
}

func TestRectangleArea(t *testing.T) {
	tests := []struct {
		width, height float64
		want          float64
	}{
		{5, 3, 15},
		{10, 10, 100},
		{0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.0fx%.0f", tt.width, tt.height), func(t *testing.T) {
			r := Rectangle{Width: tt.width, Height: tt.height}
			got := r.Area()
			if got != tt.want {
				t.Errorf("Area() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestRectanglePerimeter(t *testing.T) {
	r := Rectangle{Width: 5, Height: 3}
	got := r.Perimeter()
	want := 16.0
	if got != want {
		t.Errorf("Perimeter() = %f, want %f", got, want)
	}
}

func TestRectangleIsSquare(t *testing.T) {
	tests := []struct {
		name          string
		width, height float64
		want          bool
	}{
		{"square", 5, 5, true},
		{"not square", 5, 3, false},
		{"zero square", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rectangle{Width: tt.width, Height: tt.height}
			got := r.IsSquare()
			if got != tt.want {
				t.Errorf("IsSquare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRectangleScale(t *testing.T) {
	r := Rectangle{Width: 5, Height: 3}
	scaled := r.Scale(2)

	if scaled.Width != 10 || scaled.Height != 6 {
		t.Errorf("Scale(2) = {%f, %f}, want {10, 6}", scaled.Width, scaled.Height)
	}

	// Original should be unchanged (value receiver)
	if r.Width != 5 || r.Height != 3 {
		t.Error("Original rectangle was modified — Scale should use a value receiver")
	}
}

// =========================================================================
// Exercise 3 Tests: BankAccount
// =========================================================================

func TestNewBankAccount(t *testing.T) {
	t.Run("positive balance", func(t *testing.T) {
		a := NewBankAccount("Alice", 100)
		if a == nil {
			t.Fatal("NewBankAccount returned nil")
		}
		if a.Owner != "Alice" {
			t.Errorf("Owner = %q, want %q", a.Owner, "Alice")
		}
		if a.Balance() != 100 {
			t.Errorf("Balance() = %f, want %f", a.Balance(), 100.0)
		}
	})

	t.Run("negative balance becomes zero", func(t *testing.T) {
		a := NewBankAccount("Bob", -50)
		if a == nil {
			t.Fatal("NewBankAccount returned nil")
		}
		if a.Balance() != 0 {
			t.Errorf("Balance() = %f, want 0 (negative initial balance should become 0)", a.Balance())
		}
	})
}

func TestBankAccountDeposit(t *testing.T) {
	a := NewBankAccount("Alice", 100)
	if a == nil {
		t.Fatal("NewBankAccount returned nil")
	}

	t.Run("valid deposit", func(t *testing.T) {
		msg := a.Deposit(50)
		if msg != "" {
			t.Errorf("Deposit(50) returned error %q, want empty string", msg)
		}
		if a.Balance() != 150 {
			t.Errorf("Balance after deposit = %f, want 150", a.Balance())
		}
	})

	t.Run("zero deposit", func(t *testing.T) {
		msg := a.Deposit(0)
		if msg == "" {
			t.Error("Deposit(0) should return an error message")
		}
	})

	t.Run("negative deposit", func(t *testing.T) {
		msg := a.Deposit(-10)
		if msg == "" {
			t.Error("Deposit(-10) should return an error message")
		}
	})
}

func TestBankAccountWithdraw(t *testing.T) {
	a := NewBankAccount("Alice", 100)
	if a == nil {
		t.Fatal("NewBankAccount returned nil")
	}

	t.Run("valid withdraw", func(t *testing.T) {
		msg := a.Withdraw(30)
		if msg != "" {
			t.Errorf("Withdraw(30) returned error %q, want empty string", msg)
		}
		if a.Balance() != 70 {
			t.Errorf("Balance after withdraw = %f, want 70", a.Balance())
		}
	})

	t.Run("overdraft", func(t *testing.T) {
		msg := a.Withdraw(1000)
		if msg == "" {
			t.Error("Withdraw(1000) should return error — insufficient funds")
		}
		if a.Balance() != 70 {
			t.Errorf("Balance should be unchanged after failed withdraw, got %f", a.Balance())
		}
	})

	t.Run("negative withdraw", func(t *testing.T) {
		msg := a.Withdraw(-10)
		if msg == "" {
			t.Error("Withdraw(-10) should return an error message")
		}
	})
}

func TestBankAccountTransfer(t *testing.T) {
	from := NewBankAccount("Alice", 100)
	to := NewBankAccount("Bob", 50)
	if from == nil || to == nil {
		t.Fatal("NewBankAccount returned nil")
	}

	t.Run("valid transfer", func(t *testing.T) {
		msg := from.Transfer(30, to)
		if msg != "" {
			t.Errorf("Transfer returned error %q, want empty string", msg)
		}
		if from.Balance() != 70 {
			t.Errorf("Sender balance = %f, want 70", from.Balance())
		}
		if to.Balance() != 80 {
			t.Errorf("Receiver balance = %f, want 80", to.Balance())
		}
	})

	t.Run("insufficient funds", func(t *testing.T) {
		msg := from.Transfer(1000, to)
		if msg == "" {
			t.Error("Transfer with insufficient funds should return error")
		}
	})
}

// =========================================================================
// Exercise 4 Tests: Admin Embedding
// =========================================================================

func TestNewAdmin(t *testing.T) {
	a := NewAdmin("Ada", "Lovelace", "ada@example.com", 30, "superadmin", []string{"read", "write", "delete"})
	if a == nil {
		t.Fatal("NewAdmin returned nil")
	}

	// Test promoted fields from embedded User
	if a.FirstName != "Ada" {
		t.Errorf("FirstName = %q, want %q (promoted from User)", a.FirstName, "Ada")
	}
	if a.LastName != "Lovelace" {
		t.Errorf("LastName = %q, want %q", a.LastName, "Lovelace")
	}
	if a.Email != "ada@example.com" {
		t.Errorf("Email = %q, want %q", a.Email, "ada@example.com")
	}

	// Test promoted method
	if a.FullName() != "Ada Lovelace" {
		t.Errorf("FullName() = %q, want %q (promoted from User)", a.FullName(), "Ada Lovelace")
	}

	// Test Admin-specific fields
	if a.Role != "superadmin" {
		t.Errorf("Role = %q, want %q", a.Role, "superadmin")
	}
}

func TestAdminHasPermission(t *testing.T) {
	a := Admin{
		Permissions: []string{"read", "write", "delete"},
	}

	tests := []struct {
		perm string
		want bool
	}{
		{"read", true},
		{"write", true},
		{"delete", true},
		{"admin", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.perm, func(t *testing.T) {
			got := a.HasPermission(tt.perm)
			if got != tt.want {
				t.Errorf("HasPermission(%q) = %v, want %v", tt.perm, got, tt.want)
			}
		})
	}
}

func TestAdminPromote(t *testing.T) {
	a := &Admin{
		Permissions: []string{"read"},
	}

	t.Run("add new permission", func(t *testing.T) {
		a.Promote("write")
		if !a.HasPermission("write") {
			t.Error("After Promote(\"write\"), HasPermission(\"write\") should be true")
		}
	})

	t.Run("no duplicates", func(t *testing.T) {
		a.Promote("read") // already exists
		count := 0
		for _, p := range a.Permissions {
			if p == "read" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("Found %d copies of \"read\" — Promote should not add duplicates", count)
		}
	})
}

// =========================================================================
// Exercise 5 Tests: Linked List
// =========================================================================

func TestNewLinkedList(t *testing.T) {
	t.Run("from slice", func(t *testing.T) {
		head := NewLinkedList([]int{1, 2, 3})
		if head == nil {
			t.Fatal("NewLinkedList returned nil for non-empty slice")
		}
		got := head.ToSlice()
		want := []int{1, 2, 3}
		if len(got) != len(want) {
			t.Fatalf("ToSlice() length = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("ToSlice()[%d] = %d, want %d", i, got[i], want[i])
			}
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		head := NewLinkedList([]int{})
		if head != nil {
			t.Error("NewLinkedList([]) should return nil")
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		head := NewLinkedList(nil)
		if head != nil {
			t.Error("NewLinkedList(nil) should return nil")
		}
	})
}

func TestLinkedListToSlice(t *testing.T) {
	t.Run("nil node returns empty slice", func(t *testing.T) {
		var n *ListNode
		got := n.ToSlice()
		if got == nil {
			t.Error("ToSlice() on nil node should return empty (non-nil) slice, got nil")
		}
		if len(got) != 0 {
			t.Errorf("ToSlice() on nil node should return empty slice, got length %d", len(got))
		}
	})
}

func TestLinkedListLen(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"three elements", []int{1, 2, 3}, 3},
		{"single element", []int{42}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := NewLinkedList(tt.values)
			if head == nil {
				t.Fatal("NewLinkedList returned nil")
			}
			got := head.Len()
			if got != tt.want {
				t.Errorf("Len() = %d, want %d", got, tt.want)
			}
		})
	}

	t.Run("nil node", func(t *testing.T) {
		var n *ListNode
		if n.Len() != 0 {
			t.Errorf("Len() on nil node = %d, want 0", n.Len())
		}
	})
}

func TestLinkedListAppend(t *testing.T) {
	head := NewLinkedList([]int{1, 2})
	if head == nil {
		t.Fatal("NewLinkedList returned nil")
	}
	head.Append(3)
	head.Append(4)

	got := head.ToSlice()
	want := []int{1, 2, 3, 4}
	if len(got) != len(want) {
		t.Fatalf("After Append, length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("After Append, [%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

// =========================================================================
// Exercise 6 Tests: Request
// =========================================================================

func TestNewRequest(t *testing.T) {
	r := NewRequest("GET", "/api/users")
	if r == nil {
		t.Fatal("NewRequest returned nil")
	}
	if r.Method != "GET" {
		t.Errorf("Method = %q, want %q", r.Method, "GET")
	}
	if r.Path != "/api/users" {
		t.Errorf("Path = %q, want %q", r.Path, "/api/users")
	}
	if r.Headers == nil {
		t.Error("Headers should be initialized (not nil) — did you make the map?")
	}
}

func TestRequestHeaders(t *testing.T) {
	r := NewRequest("GET", "/")
	if r == nil {
		t.Fatal("NewRequest returned nil")
	}

	t.Run("add and get header", func(t *testing.T) {
		r.AddHeader("Content-Type", "application/json")
		got := r.GetHeader("Content-Type")
		if got != "application/json" {
			t.Errorf("GetHeader(\"Content-Type\") = %q, want %q", got, "application/json")
		}
	})

	t.Run("multiple values", func(t *testing.T) {
		r.AddHeader("Accept", "text/html")
		r.AddHeader("Accept", "application/json")
		all := r.GetAllHeaders("Accept")
		if len(all) != 2 {
			t.Fatalf("GetAllHeaders(\"Accept\") length = %d, want 2", len(all))
		}
		if all[0] != "text/html" || all[1] != "application/json" {
			t.Errorf("GetAllHeaders = %v, want [text/html application/json]", all)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		got := r.GetHeader("X-Missing")
		if got != "" {
			t.Errorf("GetHeader for missing key = %q, want empty string", got)
		}
	})

	t.Run("missing header all", func(t *testing.T) {
		got := r.GetAllHeaders("X-Missing")
		if got != nil {
			t.Errorf("GetAllHeaders for missing key = %v, want nil", got)
		}
	})
}

func TestRequestIsSecure(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"https://example.com/api", true},
		{"http://example.com/api", false},
		{"/api/users", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			r := &Request{Path: tt.path}
			got := r.IsSecure()
			if got != tt.want {
				t.Errorf("IsSecure() for path %q = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// =========================================================================
// Exercise 7 Tests: Temperature / Stringer
// =========================================================================

func TestTemperatureString(t *testing.T) {
	tests := []struct {
		temp Temperature
		want string
	}{
		{Temperature{72.0, 'F'}, "72.0°F"},
		{Temperature{22.5, 'C'}, "22.5°C"},
		{Temperature{0.0, 'C'}, "0.0°C"},
		{Temperature{100.0, 'F'}, "100.0°F"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.temp.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}

			// Also verify it works with fmt.Sprintf %v (Stringer interface)
			gotFmt := fmt.Sprintf("%v", tt.temp)
			if gotFmt != tt.want {
				t.Errorf("fmt.Sprintf(\"%%v\") = %q, want %q — make sure String() is implemented", gotFmt, tt.want)
			}
		})
	}
}

func TestTemperatureToFahrenheit(t *testing.T) {
	tests := []struct {
		name string
		temp Temperature
		want float64
	}{
		{"0C to F", Temperature{0, 'C'}, 32.0},
		{"100C to F", Temperature{100, 'C'}, 212.0},
		{"already F", Temperature{72, 'F'}, 72.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.temp.ToFahrenheit()
			if got.Unit != 'F' {
				t.Errorf("Unit = %c, want F", got.Unit)
			}
			if got.Degrees < tt.want-0.1 || got.Degrees > tt.want+0.1 {
				t.Errorf("Degrees = %f, want %f", got.Degrees, tt.want)
			}
		})
	}
}

func TestTemperatureToCelsius(t *testing.T) {
	tests := []struct {
		name string
		temp Temperature
		want float64
	}{
		{"32F to C", Temperature{32, 'F'}, 0.0},
		{"212F to C", Temperature{212, 'F'}, 100.0},
		{"already C", Temperature{22.5, 'C'}, 22.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.temp.ToCelsius()
			if got.Unit != 'C' {
				t.Errorf("Unit = %c, want C", got.Unit)
			}
			if got.Degrees < tt.want-0.1 || got.Degrees > tt.want+0.1 {
				t.Errorf("Degrees = %f, want %f", got.Degrees, tt.want)
			}
		})
	}
}

// =========================================================================
// Exercise 8 Tests: Builder Pattern
// =========================================================================

func TestServerBuilderDefaults(t *testing.T) {
	b := NewServerBuilder()
	if b == nil {
		t.Fatal("NewServerBuilder returned nil")
	}

	cfg := b.Build()

	if cfg.Host != "localhost" {
		t.Errorf("Default Host = %q, want %q", cfg.Host, "localhost")
	}
	if cfg.Port != 8080 {
		t.Errorf("Default Port = %d, want %d", cfg.Port, 8080)
	}
	if cfg.TLSEnabled != false {
		t.Error("Default TLSEnabled should be false")
	}
	if cfg.ReadTimeout != 30 {
		t.Errorf("Default ReadTimeout = %d, want 30", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 30 {
		t.Errorf("Default WriteTimeout = %d, want 30", cfg.WriteTimeout)
	}
	if cfg.MaxConns != 100 {
		t.Errorf("Default MaxConns = %d, want 100", cfg.MaxConns)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Default LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
}

func TestServerBuilderChaining(t *testing.T) {
	cfg := NewServerBuilder().
		WithHost("api.example.com").
		WithPort(443).
		WithTLS(true).
		WithTimeouts(60, 60).
		WithMaxConns(500).
		WithLogLevel("debug").
		Build()

	if cfg.Host != "api.example.com" {
		t.Errorf("Host = %q, want %q", cfg.Host, "api.example.com")
	}
	if cfg.Port != 443 {
		t.Errorf("Port = %d, want 443", cfg.Port)
	}
	if cfg.TLSEnabled != true {
		t.Error("TLSEnabled should be true")
	}
	if cfg.ReadTimeout != 60 {
		t.Errorf("ReadTimeout = %d, want 60", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 60 {
		t.Errorf("WriteTimeout = %d, want 60", cfg.WriteTimeout)
	}
	if cfg.MaxConns != 500 {
		t.Errorf("MaxConns = %d, want 500", cfg.MaxConns)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
}

func TestServerConfigString(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want string
	}{
		{
			"with TLS",
			ServerConfig{Host: "example.com", Port: 443, TLSEnabled: true},
			"example.com:443 (TLS: enabled)",
		},
		{
			"without TLS",
			ServerConfig{Host: "localhost", Port: 8080, TLSEnabled: false},
			"localhost:8080 (TLS: disabled)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
