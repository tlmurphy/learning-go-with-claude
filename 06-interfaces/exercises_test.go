package interfaces

import (
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
)

// =========================================================================
// Exercise 1 Tests: fmt.Stringer
// =========================================================================

func TestColorString(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{Color{255, 0, 0}, "#ff0000"},
		{Color{0, 255, 0}, "#00ff00"},
		{Color{0, 0, 255}, "#0000ff"},
		{Color{255, 128, 0}, "#ff8000"},
		{Color{0, 0, 0}, "#000000"},
		{Color{255, 255, 255}, "#ffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.color.String()
			if got != tt.want {
				t.Errorf("Color{%d, %d, %d}.String() = %q, want %q",
					tt.color.R, tt.color.G, tt.color.B, got, tt.want)
			}

			// Verify it works through the Stringer interface
			gotFmt := fmt.Sprint(tt.color)
			if gotFmt != tt.want {
				t.Errorf("fmt.Sprint(color) = %q, want %q", gotFmt, tt.want)
			}
		})
	}
}

func TestBookString(t *testing.T) {
	b := Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan", Pages: 380}
	want := "The Go Programming Language by Donovan & Kernighan (380 pages)"
	got := b.String()
	if got != want {
		t.Errorf("Book.String() = %q, want %q", got, want)
	}
}

// =========================================================================
// Exercise 2 Tests: Shape Interface
// =========================================================================

func TestRectShape(t *testing.T) {
	r := Rect{Width: 5, Height: 3}

	t.Run("area", func(t *testing.T) {
		got := r.Area()
		if got != 15 {
			t.Errorf("Rect{5,3}.Area() = %f, want 15", got)
		}
	})

	t.Run("perimeter", func(t *testing.T) {
		got := r.Perimeter()
		if got != 16 {
			t.Errorf("Rect{5,3}.Perimeter() = %f, want 16", got)
		}
	})

	// Verify it satisfies Shape
	var _ Shape = r
}

func TestTriangleShape(t *testing.T) {
	tr := Triangle{SideA: 6, SideB: 5, SideC: 5, Height: 4}

	t.Run("area", func(t *testing.T) {
		got := tr.Area()
		if got != 12 {
			t.Errorf("Triangle.Area() = %f, want 12", got)
		}
	})

	t.Run("perimeter", func(t *testing.T) {
		got := tr.Perimeter()
		if got != 16 {
			t.Errorf("Triangle.Perimeter() = %f, want 16", got)
		}
	})

	var _ Shape = tr
}

func TestCircleShapeShape(t *testing.T) {
	c := CircleShape{Radius: 5}

	t.Run("area", func(t *testing.T) {
		got := c.Area()
		want := math.Pi * 25
		if math.Abs(got-want) > 0.001 {
			t.Errorf("CircleShape{5}.Area() = %f, want %f", got, want)
		}
	})

	t.Run("perimeter", func(t *testing.T) {
		got := c.Perimeter()
		want := 2 * math.Pi * 5
		if math.Abs(got-want) > 0.001 {
			t.Errorf("CircleShape{5}.Perimeter() = %f, want %f", got, want)
		}
	})

	var _ Shape = c
}

func TestTotalArea(t *testing.T) {
	shapes := []Shape{
		Rect{Width: 10, Height: 5},    // area: 50
		CircleShape{Radius: 1},        // area: pi
		Triangle{SideA: 6, Height: 4}, // area: 12
	}

	got := TotalArea(shapes)
	want := 50 + math.Pi + 12
	if math.Abs(got-want) > 0.001 {
		t.Errorf("TotalArea() = %f, want %f", got, want)
	}
}

func TestTotalAreaEmpty(t *testing.T) {
	got := TotalArea(nil)
	if got != 0 {
		t.Errorf("TotalArea(nil) = %f, want 0", got)
	}
}

// =========================================================================
// Exercise 3 Tests: ROT13 Reader
// =========================================================================

func TestROT13Reader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"hello", "hello", "uryyb"},
		{"HELLO", "HELLO", "URYYB"},
		{"mixed case", "Hello, World!", "Uryyb, Jbeyq!"},
		{"numbers unchanged", "abc123", "nop123"},
		{"double ROT13", "uryyb", "hello"}, // ROT13 is its own inverse
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewROT13Reader(strings.NewReader(tt.input))
			if r == nil {
				t.Fatal("NewROT13Reader returned nil")
			}

			result, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("ReadAll error: %v", err)
			}
			got := string(result)
			if got != tt.want {
				t.Errorf("ROT13(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestROT13ReaderImplementsReader(t *testing.T) {
	r := NewROT13Reader(strings.NewReader("test"))
	if r == nil {
		t.Fatal("NewROT13Reader returned nil")
	}
	// Verify it satisfies io.Reader
	var _ io.Reader = r
}

// =========================================================================
// Exercise 4 Tests: Type Switch
// =========================================================================

func TestDescribe(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"int", 42, "integer: 42"},
		{"float", 3.14, "float: 3.14"},
		{"string", "hello", "string: hello (length: 5)"},
		{"bool true", true, "boolean: true"},
		{"bool false", false, "boolean: false"},
		{"int slice", []int{1, 2, 3}, "int slice: [1 2 3] (length: 3)"},
		{"nil", nil, "nil value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Describe(tt.value)
			if got != tt.want {
				t.Errorf("Describe(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}

	t.Run("unknown type", func(t *testing.T) {
		got := Describe([]string{"a"})
		if !strings.HasPrefix(got, "unknown type:") {
			t.Errorf("Describe([]string) = %q, should start with \"unknown type:\"", got)
		}
	})
}

func TestSummarize(t *testing.T) {
	values := []any{1, 2, "hello", true, 3.14, nil, "world", false, []byte{}}
	got := Summarize(values)

	if got == nil {
		t.Fatal("Summarize returned nil")
	}

	checks := map[string]int{
		"int":     2,
		"string":  2,
		"bool":    2,
		"float64": 1,
		"nil":     1,
		"other":   1,
	}

	for key, want := range checks {
		if got[key] != want {
			t.Errorf("Summarize[%q] = %d, want %d", key, got[key], want)
		}
	}
}

// =========================================================================
// Exercise 5 Tests: Interface Composition / Buffer
// =========================================================================

func TestBufferReadWrite(t *testing.T) {
	b := NewBuffer()
	if b == nil {
		t.Fatal("NewBuffer returned nil")
	}

	// Verify it satisfies ReadWriteCloser
	var _ ReadWriteCloser = b

	t.Run("write and read", func(t *testing.T) {
		err := b.Write("hello")
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
		err = b.Write(" world")
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}

		data, err := b.Read()
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
		if data != "hello world" {
			t.Errorf("Read() = %q, want %q", data, "hello world")
		}
	})

	t.Run("read clears buffer", func(t *testing.T) {
		b2 := NewBuffer()
		if b2 == nil {
			t.Fatal("NewBuffer returned nil")
		}
		b2.Write("data")
		b2.Read()
		data, err := b2.Read()
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
		if data != "" {
			t.Errorf("Second Read() = %q, want empty string (buffer should be cleared)", data)
		}
	})
}

func TestBufferClose(t *testing.T) {
	b := NewBuffer()
	if b == nil {
		t.Fatal("NewBuffer returned nil")
	}

	t.Run("close then write", func(t *testing.T) {
		b.Close()
		err := b.Write("data")
		if err == nil {
			t.Error("Write after Close should return an error")
		}
	})

	t.Run("close then read", func(t *testing.T) {
		b2 := NewBuffer()
		if b2 == nil {
			t.Fatal("NewBuffer returned nil")
		}
		b2.Write("data")
		b2.Close()
		_, err := b2.Read()
		if err == nil {
			t.Error("Read after Close should return an error")
		}
	})

	t.Run("double close", func(t *testing.T) {
		b3 := NewBuffer()
		if b3 == nil {
			t.Fatal("NewBuffer returned nil")
		}
		b3.Close()
		err := b3.Close()
		if err == nil {
			t.Error("Closing an already-closed buffer should return an error")
		}
	})
}

// =========================================================================
// Exercise 6 Tests: Middleware Chain
// =========================================================================

func TestHandlerFunc(t *testing.T) {
	h := HandlerFunc(func(req string) string {
		return "response: " + req
	})

	got := h.Handle("test")
	if got != "response: test" {
		t.Errorf("HandlerFunc.Handle() = %q, want %q", got, "response: test")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	mw := LoggingMiddleware()
	if mw == nil {
		t.Fatal("LoggingMiddleware returned nil")
	}

	inner := HandlerFunc(func(req string) string {
		return "OK"
	})

	handler := mw(inner)
	got := handler.Handle("GET /api")
	want := "[LOG] request=GET /api response=OK"
	if got != want {
		t.Errorf("LoggingMiddleware result = %q, want %q", got, want)
	}
}

func TestAuthMiddleware(t *testing.T) {
	mw := AuthMiddleware()
	if mw == nil {
		t.Fatal("AuthMiddleware returned nil")
	}

	inner := HandlerFunc(func(req string) string {
		return "data: " + req
	})

	handler := mw(inner)

	t.Run("authorized", func(t *testing.T) {
		got := handler.Handle("authorized GET /api")
		want := "data: GET /api"
		if got != want {
			t.Errorf("AuthMiddleware(authorized) = %q, want %q", got, want)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		got := handler.Handle("GET /api")
		if got != "401 Unauthorized" {
			t.Errorf("AuthMiddleware(unauthorized) = %q, want %q", got, "401 Unauthorized")
		}
	})
}

func TestChain(t *testing.T) {
	inner := HandlerFunc(func(req string) string {
		return "OK"
	})

	logging := LoggingMiddleware()
	auth := AuthMiddleware()
	if logging == nil || auth == nil {
		t.Fatal("Middleware returned nil")
	}

	chained := Chain(inner, logging, auth)

	t.Run("authorized request flows through all", func(t *testing.T) {
		got := chained.Handle("authorized GET /api")
		want := "[LOG] request=authorized GET /api response=OK"
		if got != want {
			t.Errorf("Chained = %q, want %q", got, want)
		}
	})

	t.Run("unauthorized blocked by auth", func(t *testing.T) {
		got := chained.Handle("GET /api")
		want := "[LOG] request=GET /api response=401 Unauthorized"
		if got != want {
			t.Errorf("Chained (unauth) = %q, want %q", got, want)
		}
	})
}

// =========================================================================
// Exercise 7 Tests: Storage / Dependency Injection
// =========================================================================

func TestMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()
	if ms == nil {
		t.Fatal("NewMemoryStorage returned nil")
	}

	// Verify it satisfies Storage
	var _ Storage = ms

	t.Run("put and get", func(t *testing.T) {
		err := ms.Put(Item{ID: "1", Data: "hello"})
		if err != nil {
			t.Fatalf("Put error: %v", err)
		}

		item, ok := ms.Get("1")
		if !ok {
			t.Fatal("Get(\"1\") returned ok=false, want true")
		}
		if item.Data != "hello" {
			t.Errorf("item.Data = %q, want %q", item.Data, "hello")
		}
	})

	t.Run("get missing", func(t *testing.T) {
		_, ok := ms.Get("nonexistent")
		if ok {
			t.Error("Get(\"nonexistent\") should return ok=false")
		}
	})

	t.Run("put empty ID", func(t *testing.T) {
		err := ms.Put(Item{ID: "", Data: "no id"})
		if err == nil {
			t.Error("Put with empty ID should return an error")
		}
	})

	t.Run("delete", func(t *testing.T) {
		ms.Put(Item{ID: "del", Data: "delete me"})
		err := ms.Delete("del")
		if err != nil {
			t.Fatalf("Delete error: %v", err)
		}
		_, ok := ms.Get("del")
		if ok {
			t.Error("Item should be deleted")
		}
	})

	t.Run("delete missing", func(t *testing.T) {
		err := ms.Delete("nonexistent")
		if err == nil {
			t.Error("Delete(\"nonexistent\") should return an error")
		}
	})

	t.Run("list", func(t *testing.T) {
		ms2 := NewMemoryStorage()
		if ms2 == nil {
			t.Fatal("NewMemoryStorage returned nil")
		}
		ms2.Put(Item{ID: "a", Data: "alpha"})
		ms2.Put(Item{ID: "b", Data: "beta"})

		items := ms2.List()
		if len(items) != 2 {
			t.Errorf("List() length = %d, want 2", len(items))
		}
	})
}

func TestItemService(t *testing.T) {
	ms := NewMemoryStorage()
	if ms == nil {
		t.Fatal("NewMemoryStorage returned nil")
	}
	svc := NewItemService(ms)
	if svc == nil {
		t.Fatal("NewItemService returned nil")
	}

	t.Run("save and get", func(t *testing.T) {
		err := svc.SaveItem(Item{ID: "1", Data: "test"})
		if err != nil {
			t.Fatalf("SaveItem error: %v", err)
		}

		item, err := svc.GetItem("1")
		if err != nil {
			t.Fatalf("GetItem error: %v", err)
		}
		if item.Data != "test" {
			t.Errorf("item.Data = %q, want %q", item.Data, "test")
		}
	})

	t.Run("get missing returns error", func(t *testing.T) {
		_, err := svc.GetItem("nonexistent")
		if err == nil {
			t.Error("GetItem(\"nonexistent\") should return an error")
		}
	})
}

// =========================================================================
// Exercise 8 Tests: Nil Interface Gotcha
// =========================================================================

func TestEmailValidator(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"user@example.com", false},
		{"invalid-email", true},
		{"", true},
		{"@", false}, // contains @, so it passes our simple check
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			v := &EmailValidator{Email: tt.email}
			err := v.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("Validate(%q) should return error", tt.email)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate(%q) unexpected error: %v", tt.email, err)
			}
		})
	}
}

func TestGetValidator(t *testing.T) {
	t.Run("email validator", func(t *testing.T) {
		v := GetValidator("email", "user@example.com")
		if v == nil {
			t.Fatal("GetValidator(\"email\", ...) returned nil interface — should return an EmailValidator")
		}
		err := v.Validate()
		if err != nil {
			t.Errorf("Valid email should not error: %v", err)
		}
	})

	t.Run("unknown type returns truly nil", func(t *testing.T) {
		v := GetValidator("phone", "555-1234")
		if v != nil {
			t.Error("GetValidator(\"phone\", ...) should return nil (truly nil interface, not a typed nil pointer)")
		}
	})
}

func TestIsNilInterface(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		var v Validator
		if !IsNilInterface(v) {
			t.Error("IsNilInterface(nil Validator) should be true")
		}
	})

	t.Run("typed nil pointer", func(t *testing.T) {
		var e *EmailValidator = nil
		// When assigned to any, this becomes (type=*EmailValidator, value=nil)
		// which is NOT a nil interface
		if IsNilInterface(e) {
			t.Error("IsNilInterface(typed nil *EmailValidator) should be false — " +
				"it holds a type (*EmailValidator) even though the value is nil")
		}
	})

	t.Run("non-nil value", func(t *testing.T) {
		e := &EmailValidator{Email: "test@test.com"}
		if IsNilInterface(e) {
			t.Error("IsNilInterface(non-nil EmailValidator) should be false")
		}
	})

	t.Run("plain nil", func(t *testing.T) {
		if !IsNilInterface(nil) {
			t.Error("IsNilInterface(nil) should be true")
		}
	})
}
