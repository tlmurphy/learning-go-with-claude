package errorhandling

import (
	"errors"
	"fmt"
	"testing"
)

// =========================================================================
// Exercise 1 Tests: Division
// =========================================================================

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"normal division", 10, 2, 5, nil},
		{"decimal result", 7, 2, 3.5, nil},
		{"divide by zero", 10, 0, 0, ErrDivisionByZero},
		{"zero numerator", 0, 5, 0, nil},
		{"negative numbers", -10, 2, -5, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Divide(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
				return
			}
			if got != tt.want {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name       string
		a, b, def  float64
		want       float64
	}{
		{"normal", 10, 2, -1, 5},
		{"divide by zero returns default", 10, 0, -1, -1},
		{"divide by zero custom default", 10, 0, 999, 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeDivide(tt.a, tt.b, tt.def)
			if got != tt.want {
				t.Errorf("SafeDivide(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.def, got, tt.want)
			}
		})
	}
}

// =========================================================================
// Exercise 2 Tests: Sentinel Errors
// =========================================================================

func TestStoreFind(t *testing.T) {
	store := NewStore(map[string]StoreItem{
		"widget": {Name: "widget", Price: 9.99, Quantity: 10},
	})
	if store == nil {
		t.Fatal("NewStore returned nil")
	}

	t.Run("find existing", func(t *testing.T) {
		item, err := store.FindItem("widget")
		if err != nil {
			t.Fatalf("FindItem(\"widget\") error: %v", err)
		}
		if item.Name != "widget" {
			t.Errorf("item.Name = %q, want %q", item.Name, "widget")
		}
	})

	t.Run("find missing", func(t *testing.T) {
		_, err := store.FindItem("gadget")
		if err == nil {
			t.Fatal("FindItem(\"gadget\") should return error")
		}
		if !errors.Is(err, ErrItemNotFound) {
			t.Errorf("error should wrap ErrItemNotFound, got: %v", err)
		}
	})
}

func TestStorePurchase(t *testing.T) {
	store := NewStore(map[string]StoreItem{
		"widget": {Name: "widget", Price: 10.00, Quantity: 5},
	})
	if store == nil {
		t.Fatal("NewStore returned nil")
	}

	t.Run("valid purchase", func(t *testing.T) {
		total, err := store.Purchase("widget", 2)
		if err != nil {
			t.Fatalf("Purchase error: %v", err)
		}
		if total != 20.00 {
			t.Errorf("total = %f, want 20.00", total)
		}
		// Check quantity decreased
		item, _ := store.FindItem("widget")
		if item.Quantity != 3 {
			t.Errorf("remaining quantity = %d, want 3", item.Quantity)
		}
	})

	t.Run("item not found", func(t *testing.T) {
		_, err := store.Purchase("nonexistent", 1)
		if !errors.Is(err, ErrItemNotFound) {
			t.Errorf("error should wrap ErrItemNotFound, got: %v", err)
		}
	})

	t.Run("invalid quantity", func(t *testing.T) {
		_, err := store.Purchase("widget", 0)
		if !errors.Is(err, ErrInvalidQuantity) {
			t.Errorf("error should wrap ErrInvalidQuantity, got: %v", err)
		}
		_, err = store.Purchase("widget", -1)
		if !errors.Is(err, ErrInvalidQuantity) {
			t.Errorf("error should wrap ErrInvalidQuantity for negative qty, got: %v", err)
		}
	})

	t.Run("out of stock", func(t *testing.T) {
		_, err := store.Purchase("widget", 100)
		if !errors.Is(err, ErrItemOutOfStock) {
			t.Errorf("error should wrap ErrItemOutOfStock, got: %v", err)
		}
	})
}

// =========================================================================
// Exercise 3 Tests: Custom Error Types
// =========================================================================

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "email", Message: "must contain @", Code: 2001}
	want := `validation error on field "email": must contain @ (code: 2001)`
	got := err.Error()
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestIsValidationError(t *testing.T) {
	t.Run("direct ValidationError", func(t *testing.T) {
		err := &ValidationError{Field: "age", Message: "too young", Code: 1001}
		ve, ok := IsValidationError(err)
		if !ok {
			t.Fatal("IsValidationError should return true for *ValidationError")
		}
		if ve.Field != "age" {
			t.Errorf("Field = %q, want %q", ve.Field, "age")
		}
	})

	t.Run("wrapped ValidationError", func(t *testing.T) {
		inner := &ValidationError{Field: "name", Message: "too short", Code: 3001}
		wrapped := fmt.Errorf("processing form: %w", inner)
		ve, ok := IsValidationError(wrapped)
		if !ok {
			t.Fatal("IsValidationError should find wrapped *ValidationError")
		}
		if ve.Field != "name" {
			t.Errorf("Field = %q, want %q", ve.Field, "name")
		}
	})

	t.Run("not a ValidationError", func(t *testing.T) {
		err := errors.New("random error")
		_, ok := IsValidationError(err)
		if ok {
			t.Error("IsValidationError should return false for non-ValidationError")
		}
	})
}

func TestValidateAge(t *testing.T) {
	tests := []struct {
		age     int
		wantErr bool
		code    int
	}{
		{25, false, 0},
		{0, false, 0},
		{150, false, 0},
		{-1, true, 1001},
		{151, true, 1002},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("age_%d", tt.age), func(t *testing.T) {
			err := ValidateAge(tt.age)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidateAge(%d) should return error", tt.age)
				}
				ve, ok := IsValidationError(err)
				if !ok {
					t.Fatalf("error should be *ValidationError, got %T", err)
				}
				if ve.Code != tt.code {
					t.Errorf("Code = %d, want %d", ve.Code, tt.code)
				}
				if ve.Field != "age" {
					t.Errorf("Field = %q, want %q", ve.Field, "age")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAge(%d) unexpected error: %v", tt.age, err)
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"user@example.com", false},
		{"a@b.c", false},
		{"invalid", true},
		{"no-at-sign.com", true},
		{"no-dot@here", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr && err == nil {
				t.Errorf("ValidateEmail(%q) should return error", tt.email)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidateEmail(%q) unexpected error: %v", tt.email, err)
			}
			if err != nil {
				ve, ok := IsValidationError(err)
				if !ok {
					t.Errorf("error should be *ValidationError, got %T", err)
				}
				if ok && ve.Code != 2001 {
					t.Errorf("Code = %d, want 2001", ve.Code)
				}
			}
		})
	}
}

// =========================================================================
// Exercise 4 Tests: Error Wrapping Chain
// =========================================================================

func TestConnectToDatabase(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		err := ConnectToDatabase("timeout")
		if err == nil {
			t.Fatal("should return error for timeout host")
		}
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("error should wrap ErrTimeout, got: %v", err)
		}
	})

	t.Run("refused", func(t *testing.T) {
		err := ConnectToDatabase("refused")
		if err == nil {
			t.Fatal("should return error for refused host")
		}
		if !errors.Is(err, ErrConnection) {
			t.Errorf("error should wrap ErrConnection, got: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		err := ConnectToDatabase("localhost")
		if err != nil {
			t.Errorf("should succeed for normal host, got: %v", err)
		}
	})
}

func TestHandleUserRequest(t *testing.T) {
	t.Run("wraps through layers", func(t *testing.T) {
		err := HandleUserRequest("timeout")
		if err == nil {
			t.Fatal("should return error")
		}
		// Should still be able to find ErrTimeout through the chain
		if !errors.Is(err, ErrTimeout) {
			t.Error("error chain should contain ErrTimeout")
		}
		// Error message should contain context from all layers
		msg := err.Error()
		if msg == "" {
			t.Error("error message should not be empty")
		}
	})

	t.Run("success", func(t *testing.T) {
		err := HandleUserRequest("localhost")
		if err != nil {
			t.Errorf("should succeed, got: %v", err)
		}
	})
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"timeout", fmt.Errorf("wrapped: %w", ErrTimeout), "timeout"},
		{"connection", fmt.Errorf("wrapped: %w", ErrConnection), "connection"},
		{"validation", fmt.Errorf("wrapped: %w", &ValidationError{Field: "x", Message: "bad"}), "validation"},
		{"unknown", errors.New("something else"), "unknown"},
		{"deeply wrapped timeout", fmt.Errorf("a: %w", fmt.Errorf("b: %w", ErrTimeout)), "timeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyError(tt.err)
			if got != tt.want {
				t.Errorf("ClassifyError() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =========================================================================
// Exercise 5 Tests: Error Aggregator
// =========================================================================

func TestMultiError(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		me := &MultiError{}
		if me.HasErrors() {
			t.Error("empty MultiError should not have errors")
		}
		if me.ErrorOrNil() != nil {
			t.Error("ErrorOrNil() should return nil when empty")
		}
	})

	t.Run("add errors", func(t *testing.T) {
		me := &MultiError{}
		me.Add(errors.New("error 1"))
		me.Add(nil) // should be ignored
		me.Add(errors.New("error 2"))

		if !me.HasErrors() {
			t.Error("MultiError with errors should have errors")
		}
		if len(me.Errors) != 2 {
			t.Errorf("len(Errors) = %d, want 2 (nil should be ignored)", len(me.Errors))
		}
	})

	t.Run("error message", func(t *testing.T) {
		me := &MultiError{}
		me.Add(errors.New("first"))
		me.Add(errors.New("second"))
		me.Add(errors.New("third"))

		want := "first; second; third"
		got := me.Error()
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("error or nil returns error", func(t *testing.T) {
		me := &MultiError{}
		me.Add(errors.New("oops"))

		err := me.ErrorOrNil()
		if err == nil {
			t.Error("ErrorOrNil() should return non-nil when errors exist")
		}
	})
}

// =========================================================================
// Exercise 6 Tests: HTTP Error
// =========================================================================

func TestHTTPError(t *testing.T) {
	t.Run("without underlying error", func(t *testing.T) {
		err := NewHTTPError(404, "not found")
		if err == nil {
			t.Fatal("NewHTTPError returned nil")
		}
		want := "HTTP 404: not found"
		got := err.Error()
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("with underlying error", func(t *testing.T) {
		inner := errors.New("connection refused")
		err := WrapHTTPError(502, "bad gateway", inner)
		if err == nil {
			t.Fatal("WrapHTTPError returned nil")
		}
		want := "HTTP 502: bad gateway: connection refused"
		got := err.Error()
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("unwrap", func(t *testing.T) {
		inner := ErrConnection
		err := WrapHTTPError(502, "bad gateway", inner)
		if err == nil {
			t.Fatal("WrapHTTPError returned nil")
		}
		if !errors.Is(err, ErrConnection) {
			t.Error("Unwrap should allow errors.Is to find the inner error")
		}
	})
}

func TestIsHTTPError(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		err := NewHTTPError(400, "bad request")
		he, ok := IsHTTPError(err)
		if !ok {
			t.Fatal("IsHTTPError should return true")
		}
		if he.Status != 400 {
			t.Errorf("Status = %d, want 400", he.Status)
		}
	})

	t.Run("wrapped", func(t *testing.T) {
		inner := NewHTTPError(404, "not found")
		wrapped := fmt.Errorf("API error: %w", inner)
		he, ok := IsHTTPError(wrapped)
		if !ok {
			t.Fatal("IsHTTPError should find wrapped *HTTPError")
		}
		if he.Status != 404 {
			t.Errorf("Status = %d, want 404", he.Status)
		}
	})

	t.Run("not HTTP error", func(t *testing.T) {
		_, ok := IsHTTPError(errors.New("generic"))
		if ok {
			t.Error("IsHTTPError should return false for generic error")
		}
	})
}

func TestStatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"HTTP error", NewHTTPError(404, "not found"), 404},
		{"wrapped HTTP error", fmt.Errorf("context: %w", NewHTTPError(400, "bad")), 400},
		{"non-HTTP error", errors.New("generic"), 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StatusCode(tt.err)
			if got != tt.want {
				t.Errorf("StatusCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

// =========================================================================
// Exercise 7 Tests: Recover from Panic
// =========================================================================

func TestSafeRun(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		err := SafeRun(func() {
			// normal execution
		})
		if err != nil {
			t.Errorf("SafeRun should return nil for normal function, got: %v", err)
		}
	})

	t.Run("with panic string", func(t *testing.T) {
		err := SafeRun(func() {
			panic("something went wrong")
		})
		if err == nil {
			t.Fatal("SafeRun should return error when function panics")
		}
	})

	t.Run("with panic error", func(t *testing.T) {
		err := SafeRun(func() {
			panic(errors.New("an error"))
		})
		if err == nil {
			t.Fatal("SafeRun should return error when function panics")
		}
	})
}

func TestSafeRunWithResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result, err := SafeRunWithResult(func() string {
			return "hello"
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "hello" {
			t.Errorf("result = %q, want %q", result, "hello")
		}
	})

	t.Run("panic", func(t *testing.T) {
		result, err := SafeRunWithResult(func() string {
			panic("boom")
		})
		if err == nil {
			t.Fatal("should return error on panic")
		}
		if result != "" {
			t.Errorf("result should be empty on panic, got %q", result)
		}
	})
}

func TestMustPositive(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		got := MustPositive(42)
		if got != 42 {
			t.Errorf("MustPositive(42) = %d, want 42", got)
		}
	})

	t.Run("zero panics", func(t *testing.T) {
		err := SafeRun(func() {
			MustPositive(0)
		})
		if err == nil {
			t.Error("MustPositive(0) should panic")
		}
	})

	t.Run("negative panics", func(t *testing.T) {
		err := SafeRun(func() {
			MustPositive(-5)
		})
		if err == nil {
			t.Error("MustPositive(-5) should panic")
		}
	})
}

// =========================================================================
// Exercise 8 Tests: Registration Form Validation
// =========================================================================

func TestRegistrationFormValidate(t *testing.T) {
	t.Run("all valid", func(t *testing.T) {
		f := RegistrationForm{
			Username: "alice",
			Email:    "alice@example.com",
			Password: "securepassword",
			Age:      25,
		}
		err := f.Validate()
		if err != nil {
			t.Errorf("valid form should not error, got: %v", err)
		}
	})

	t.Run("all invalid", func(t *testing.T) {
		f := RegistrationForm{
			Username: "",
			Email:    "invalid",
			Password: "short",
			Age:      5,
		}
		err := f.Validate()
		if err == nil {
			t.Fatal("invalid form should return error")
		}

		// Should have multiple errors
		var me *MultiError
		if !errors.As(err, &me) {
			t.Fatal("error should be *MultiError")
		}
		if len(me.Errors) < 3 {
			t.Errorf("expected at least 3 errors, got %d: %v", len(me.Errors), me.Errors)
		}
	})

	t.Run("username too short", func(t *testing.T) {
		f := RegistrationForm{
			Username: "ab",
			Email:    "ok@ok.com",
			Password: "longpassword",
			Age:      25,
		}
		err := f.Validate()
		if err == nil {
			t.Error("username 'ab' (2 chars) should fail validation")
		}
	})

	t.Run("age too young", func(t *testing.T) {
		f := RegistrationForm{
			Username: "alice",
			Email:    "alice@example.com",
			Password: "securepassword",
			Age:      12,
		}
		err := f.Validate()
		if err == nil {
			t.Error("age 12 should fail validation (minimum 13)")
		}
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		errs := ValidationErrors(nil)
		if len(errs) != 0 {
			t.Errorf("ValidationErrors(nil) should return empty slice, got %v", errs)
		}
	})

	t.Run("multi error", func(t *testing.T) {
		me := &MultiError{}
		me.Add(errors.New("first"))
		me.Add(errors.New("second"))

		errs := ValidationErrors(me)
		if len(errs) != 2 {
			t.Errorf("expected 2 errors, got %d", len(errs))
		}
	})

	t.Run("single error", func(t *testing.T) {
		errs := ValidationErrors(errors.New("just one"))
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errs))
		}
		if errs[0] != "just one" {
			t.Errorf("errs[0] = %q, want %q", errs[0], "just one")
		}
	})
}
