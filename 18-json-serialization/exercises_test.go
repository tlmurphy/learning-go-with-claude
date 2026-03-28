package jsonserialization

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// =========================================================================
// Test Exercise 1: Basic Struct with JSON Tags
// =========================================================================

func TestBasicMarshalUnmarshal(t *testing.T) {
	t.Run("marshal book with all fields", func(t *testing.T) {
		b := Book{
			ID:     1,
			Title:  "The Go Programming Language",
			Author: "Donovan & Kernighan",
			Pages:  380,
			ISBN:   "978-0134190440",
		}
		result, err := MarshalBook(b)
		if err != nil {
			t.Fatalf("MarshalBook returned error: %v", err)
		}

		// Unmarshal back to verify round-trip
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			t.Fatalf("Result is not valid JSON: %v", err)
		}

		if parsed["id"] != float64(1) {
			t.Errorf("Expected id=1, got %v", parsed["id"])
		}
		if parsed["title"] != "The Go Programming Language" {
			t.Errorf("Expected title='The Go Programming Language', got %v", parsed["title"])
		}
		if parsed["author"] != "Donovan & Kernighan" {
			t.Errorf("Expected author='Donovan & Kernighan', got %v", parsed["author"])
		}
		if parsed["pages"] != float64(380) {
			t.Errorf("Expected pages=380, got %v", parsed["pages"])
		}
		if parsed["isbn"] != "978-0134190440" {
			t.Errorf("Expected isbn='978-0134190440', got %v", parsed["isbn"])
		}
	})

	t.Run("marshal book omits zero/empty fields", func(t *testing.T) {
		b := Book{
			ID:     2,
			Title:  "Learning Go",
			Author: "Jon Bodner",
			// Pages and ISBN intentionally left as zero values
		}
		result, err := MarshalBook(b)
		if err != nil {
			t.Fatalf("MarshalBook returned error: %v", err)
		}

		if strings.Contains(result, "pages") {
			t.Error("Expected 'pages' to be omitted when zero, but it was present in JSON")
		}
		if strings.Contains(result, "isbn") {
			t.Error("Expected 'isbn' to be omitted when empty, but it was present in JSON")
		}
	})

	t.Run("unmarshal book from JSON", func(t *testing.T) {
		input := `{"id":3,"title":"Concurrency in Go","author":"Katherine Cox-Buday","pages":238}`
		b, err := UnmarshalBook(input)
		if err != nil {
			t.Fatalf("UnmarshalBook returned error: %v", err)
		}

		if b.ID != 3 {
			t.Errorf("Expected ID=3, got %d", b.ID)
		}
		if b.Title != "Concurrency in Go" {
			t.Errorf("Expected Title='Concurrency in Go', got %q", b.Title)
		}
		if b.Author != "Katherine Cox-Buday" {
			t.Errorf("Expected Author='Katherine Cox-Buday', got %q", b.Author)
		}
		if b.Pages != 238 {
			t.Errorf("Expected Pages=238, got %d", b.Pages)
		}
	})

	t.Run("unmarshal ignores unknown fields", func(t *testing.T) {
		input := `{"id":4,"title":"Test","author":"Author","extra_field":"ignored"}`
		_, err := UnmarshalBook(input)
		if err != nil {
			t.Errorf("UnmarshalBook should ignore unknown fields, but got error: %v", err)
		}
	})

	t.Run("round trip preserves data", func(t *testing.T) {
		original := Book{
			ID:     5,
			Title:  "Go in Action",
			Author: "William Kennedy",
			Pages:  264,
			ISBN:   "978-1617291784",
		}
		jsonStr, err := MarshalBook(original)
		if err != nil {
			t.Fatalf("MarshalBook error: %v", err)
		}

		restored, err := UnmarshalBook(jsonStr)
		if err != nil {
			t.Fatalf("UnmarshalBook error: %v", err)
		}

		if restored.ID != original.ID || restored.Title != original.Title ||
			restored.Author != original.Author || restored.Pages != original.Pages ||
			restored.ISBN != original.ISBN {
			t.Errorf("Round trip failed.\n  Original: %+v\n  Restored: %+v", original, restored)
		}
	})
}

// =========================================================================
// Test Exercise 2: Optional Fields with Pointer Types
// =========================================================================

func TestOptionalFields(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }
	boolPtr := func(b bool) *bool { return &b }

	t.Run("apply name update only", func(t *testing.T) {
		existing := map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"email":  "alice@example.com",
			"active": true,
		}
		update := UserUpdate{Name: strPtr("Bob")}
		result := ApplyUpdate(existing, update)

		if result["name"] != "Bob" {
			t.Errorf("Expected name='Bob', got %v", result["name"])
		}
		if result["age"] != 30 {
			t.Errorf("Expected age=30 (unchanged), got %v", result["age"])
		}
	})

	t.Run("apply multiple updates", func(t *testing.T) {
		existing := map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"email":  "alice@example.com",
			"active": true,
		}
		update := UserUpdate{
			Name:   strPtr("Charlie"),
			Age:    intPtr(25),
			Active: boolPtr(false),
		}
		result := ApplyUpdate(existing, update)

		if result["name"] != "Charlie" {
			t.Errorf("Expected name='Charlie', got %v", result["name"])
		}
		if result["age"] != 25 {
			t.Errorf("Expected age=25, got %v", result["age"])
		}
		if result["active"] != false {
			t.Errorf("Expected active=false, got %v", result["active"])
		}
		if result["email"] != "alice@example.com" {
			t.Errorf("Expected email unchanged, got %v", result["email"])
		}
	})

	t.Run("nil update changes nothing", func(t *testing.T) {
		existing := map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"email":  "alice@example.com",
			"active": true,
		}
		update := UserUpdate{} // all nil
		result := ApplyUpdate(existing, update)

		if result["name"] != "Alice" {
			t.Errorf("Expected name='Alice' (unchanged), got %v", result["name"])
		}
		if result["age"] != 30 {
			t.Errorf("Expected age=30 (unchanged), got %v", result["age"])
		}
	})

	t.Run("update to zero values", func(t *testing.T) {
		existing := map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"email":  "alice@example.com",
			"active": true,
		}
		// Setting age to 0 and active to false — these are valid updates
		update := UserUpdate{
			Age:    intPtr(0),
			Active: boolPtr(false),
		}
		result := ApplyUpdate(existing, update)

		if result["age"] != 0 {
			t.Errorf("Expected age=0 (explicitly set), got %v", result["age"])
		}
		if result["active"] != false {
			t.Errorf("Expected active=false (explicitly set), got %v", result["active"])
		}
	})

	t.Run("json unmarshal with optional fields", func(t *testing.T) {
		input := `{"name": "Updated Name"}`
		var update UserUpdate
		err := json.Unmarshal([]byte(input), &update)
		if err != nil {
			t.Fatalf("Failed to unmarshal UserUpdate: %v", err)
		}

		if update.Name == nil {
			t.Error("Expected Name to be non-nil after unmarshaling")
		} else if *update.Name != "Updated Name" {
			t.Errorf("Expected Name='Updated Name', got %q", *update.Name)
		}

		if update.Age != nil {
			t.Error("Expected Age to be nil (not in JSON), but it was non-nil")
		}
		if update.Email != nil {
			t.Error("Expected Email to be nil (not in JSON), but it was non-nil")
		}
	})
}

// =========================================================================
// Test Exercise 3: Custom MarshalJSON — Money Type
// =========================================================================

func TestPriceMarshal(t *testing.T) {
	tests := []struct {
		name     string
		price    Price
		expected string
	}{
		{"standard amount", Price{Cents: 1299}, `"12.99"`},
		{"whole dollars", Price{Cents: 500}, `"5.00"`},
		{"small cents", Price{Cents: 7}, `"0.07"`},
		{"zero", Price{Cents: 0}, `"0.00"`},
		{"large amount", Price{Cents: 99999}, `"999.99"`},
		{"one cent", Price{Cents: 1}, `"0.01"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.price)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestPriceUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"standard amount", `"12.99"`, 1299},
		{"whole dollars", `"5.00"`, 500},
		{"small cents", `"0.07"`, 7},
		{"zero", `"0.00"`, 0},
		{"large amount", `"999.99"`, 99999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Price
			err := json.Unmarshal([]byte(tt.input), &p)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if p.Cents != tt.expected {
				t.Errorf("Expected %d cents, got %d", tt.expected, p.Cents)
			}
		})
	}
}

func TestPriceRoundTrip(t *testing.T) {
	original := Price{Cents: 4299}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored Price
	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if restored.Cents != original.Cents {
		t.Errorf("Round trip failed: original=%d, restored=%d", original.Cents, restored.Cents)
	}
}

// =========================================================================
// Test Exercise 4: Custom UnmarshalJSON with Validation
// =========================================================================

func TestRatingValidation(t *testing.T) {
	t.Run("valid ratings", func(t *testing.T) {
		for _, val := range []int{1, 2, 3, 4, 5} {
			input, _ := json.Marshal(val)
			var r Rating
			err := json.Unmarshal(input, &r)
			if err != nil {
				t.Errorf("Rating %d should be valid, got error: %v", val, err)
			}
			if r.Value != val {
				t.Errorf("Expected Value=%d, got %d", val, r.Value)
			}
		}
	})

	t.Run("invalid rating too low", func(t *testing.T) {
		var r Rating
		err := json.Unmarshal([]byte("0"), &r)
		if err == nil {
			t.Error("Expected error for rating 0, got nil")
		}
	})

	t.Run("invalid rating too high", func(t *testing.T) {
		var r Rating
		err := json.Unmarshal([]byte("6"), &r)
		if err == nil {
			t.Error("Expected error for rating 6, got nil")
		}
	})

	t.Run("negative rating", func(t *testing.T) {
		var r Rating
		err := json.Unmarshal([]byte("-1"), &r)
		if err == nil {
			t.Error("Expected error for rating -1, got nil")
		}
	})

	t.Run("marshal rating", func(t *testing.T) {
		r := Rating{Value: 4}
		data, err := json.Marshal(r)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if string(data) != "4" {
			t.Errorf("Expected '4', got %s", string(data))
		}
	})

	t.Run("rating in struct", func(t *testing.T) {
		type Review struct {
			Rating Rating `json:"rating"`
			Text   string `json:"text"`
		}

		input := `{"rating": 5, "text": "Excellent!"}`
		var review Review
		err := json.Unmarshal([]byte(input), &review)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if review.Rating.Value != 5 {
			t.Errorf("Expected rating=5, got %d", review.Rating.Value)
		}
	})
}

// =========================================================================
// Test Exercise 5: Partial Parsing with json.RawMessage
// =========================================================================

func TestParseNotification(t *testing.T) {
	t.Run("email notification", func(t *testing.T) {
		input := `{"type":"email","payload":{"to":"user@example.com","subject":"Hello"}}`
		typ, payload, err := ParseNotification(input)
		if err != nil {
			t.Fatalf("ParseNotification error: %v", err)
		}
		if typ != "email" {
			t.Errorf("Expected type='email', got %q", typ)
		}
		ep, ok := payload.(*EmailPayload)
		if !ok {
			t.Fatalf("Expected *EmailPayload, got %T", payload)
		}
		if ep.To != "user@example.com" {
			t.Errorf("Expected To='user@example.com', got %q", ep.To)
		}
		if ep.Subject != "Hello" {
			t.Errorf("Expected Subject='Hello', got %q", ep.Subject)
		}
	})

	t.Run("sms notification", func(t *testing.T) {
		input := `{"type":"sms","payload":{"phone":"+1234567890","message":"Hi there"}}`
		typ, payload, err := ParseNotification(input)
		if err != nil {
			t.Fatalf("ParseNotification error: %v", err)
		}
		if typ != "sms" {
			t.Errorf("Expected type='sms', got %q", typ)
		}
		sp, ok := payload.(*SMSPayload)
		if !ok {
			t.Fatalf("Expected *SMSPayload, got %T", payload)
		}
		if sp.Phone != "+1234567890" {
			t.Errorf("Expected Phone='+1234567890', got %q", sp.Phone)
		}
		if sp.Message != "Hi there" {
			t.Errorf("Expected Message='Hi there', got %q", sp.Message)
		}
	})

	t.Run("push notification", func(t *testing.T) {
		input := `{"type":"push","payload":{"device_id":"dev-001","title":"Alert","body":"Something happened"}}`
		typ, payload, err := ParseNotification(input)
		if err != nil {
			t.Fatalf("ParseNotification error: %v", err)
		}
		if typ != "push" {
			t.Errorf("Expected type='push', got %q", typ)
		}
		pp, ok := payload.(*PushPayload)
		if !ok {
			t.Fatalf("Expected *PushPayload, got %T", payload)
		}
		if pp.DeviceID != "dev-001" {
			t.Errorf("Expected DeviceID='dev-001', got %q", pp.DeviceID)
		}
		if pp.Title != "Alert" {
			t.Errorf("Expected Title='Alert', got %q", pp.Title)
		}
	})

	t.Run("unknown type returns error", func(t *testing.T) {
		input := `{"type":"webhook","payload":{}}`
		_, _, err := ParseNotification(input)
		if err == nil {
			t.Error("Expected error for unknown notification type, got nil")
		}
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		_, _, err := ParseNotification(`{invalid json}`)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})
}

// =========================================================================
// Test Exercise 6: Streaming JSON Encoder
// =========================================================================

func TestWriteJSONStream(t *testing.T) {
	t.Run("write multiple entries", func(t *testing.T) {
		entries := []LogEntry{
			{Timestamp: "2024-01-15T10:00:00Z", Level: "INFO", Message: "Server started"},
			{Timestamp: "2024-01-15T10:00:01Z", Level: "DEBUG", Message: "Processing request"},
			{Timestamp: "2024-01-15T10:00:02Z", Level: "ERROR", Message: "Connection failed"},
		}

		var buf bytes.Buffer
		n, err := WriteJSONStream(&buf, entries)
		if err != nil {
			t.Fatalf("WriteJSONStream error: %v", err)
		}
		if n != 3 {
			t.Errorf("Expected 3 entries written, got %d", n)
		}

		// Verify each line is valid JSON
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		if len(lines) != 3 {
			t.Fatalf("Expected 3 lines, got %d", len(lines))
		}

		for i, line := range lines {
			var entry LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				t.Errorf("Line %d is not valid JSON: %v", i, err)
			}
		}
	})

	t.Run("write empty slice", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := WriteJSONStream(&buf, []LogEntry{})
		if err != nil {
			t.Fatalf("WriteJSONStream error: %v", err)
		}
		if n != 0 {
			t.Errorf("Expected 0 entries written, got %d", n)
		}
		if buf.Len() != 0 {
			t.Errorf("Expected empty output, got %q", buf.String())
		}
	})
}

func TestReadJSONStream(t *testing.T) {
	t.Run("read multiple entries", func(t *testing.T) {
		input := `{"timestamp":"2024-01-15T10:00:00Z","level":"INFO","message":"Started"}
{"timestamp":"2024-01-15T10:00:01Z","level":"ERROR","message":"Failed"}`

		entries, err := ReadJSONStream(strings.NewReader(input))
		if err != nil {
			t.Fatalf("ReadJSONStream error: %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("Expected 2 entries, got %d", len(entries))
		}

		if entries[0].Level != "INFO" {
			t.Errorf("Expected first entry level='INFO', got %q", entries[0].Level)
		}
		if entries[1].Message != "Failed" {
			t.Errorf("Expected second entry message='Failed', got %q", entries[1].Message)
		}
	})

	t.Run("read empty input", func(t *testing.T) {
		entries, err := ReadJSONStream(strings.NewReader(""))
		if err != nil {
			t.Fatalf("ReadJSONStream error: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("Expected 0 entries, got %d", len(entries))
		}
	})

	t.Run("round trip write then read", func(t *testing.T) {
		original := []LogEntry{
			{Timestamp: "2024-01-15T10:00:00Z", Level: "INFO", Message: "Test"},
			{Timestamp: "2024-01-15T10:00:01Z", Level: "WARN", Message: "Warning"},
		}

		var buf bytes.Buffer
		WriteJSONStream(&buf, original)

		restored, err := ReadJSONStream(&buf)
		if err != nil {
			t.Fatalf("ReadJSONStream error: %v", err)
		}

		if len(restored) != len(original) {
			t.Fatalf("Expected %d entries, got %d", len(original), len(restored))
		}

		for i, entry := range restored {
			if entry.Message != original[i].Message {
				t.Errorf("Entry %d: expected message=%q, got %q", i, original[i].Message, entry.Message)
			}
		}
	})
}

// =========================================================================
// Test Exercise 7: API Response Envelope
// =========================================================================

func TestAPIResponseEnvelope(t *testing.T) {
	t.Run("single item response", func(t *testing.T) {
		data := map[string]interface{}{"id": 1, "name": "Alice"}
		env := NewAPIResponse(data, "req-123", 0)

		if env.Meta.RequestID != "req-123" {
			t.Errorf("Expected request_id='req-123', got %q", env.Meta.RequestID)
		}
		if env.Meta.Count != 0 {
			t.Errorf("Expected count=0, got %d", env.Meta.Count)
		}
	})

	t.Run("list response with count", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		env := NewAPIResponse(items, "req-456", 3)

		if env.Meta.Count != 3 {
			t.Errorf("Expected count=3, got %d", env.Meta.Count)
		}
	})

	t.Run("marshal response", func(t *testing.T) {
		data := map[string]string{"message": "hello"}
		env := NewAPIResponse(data, "req-789", 0)

		jsonStr, err := MarshalAPIResponse(env)
		if err != nil {
			t.Fatalf("MarshalAPIResponse error: %v", err)
		}

		// Verify it's valid JSON with expected structure
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
			t.Fatalf("Result is not valid JSON: %v", err)
		}

		if _, ok := parsed["data"]; !ok {
			t.Error("Expected 'data' field in JSON output")
		}
		if _, ok := parsed["meta"]; !ok {
			t.Error("Expected 'meta' field in JSON output")
		}

		meta, _ := parsed["meta"].(map[string]interface{})
		if meta["request_id"] != "req-789" {
			t.Errorf("Expected request_id='req-789', got %v", meta["request_id"])
		}
	})
}

// =========================================================================
// Test Exercise 8: Flexible JSON Field (string or number)
// =========================================================================

func TestFlexibleID(t *testing.T) {
	t.Run("unmarshal string ID", func(t *testing.T) {
		var f FlexibleID
		err := json.Unmarshal([]byte(`"abc-123"`), &f)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if f.Value != "abc-123" {
			t.Errorf("Expected Value='abc-123', got %q", f.Value)
		}
	})

	t.Run("unmarshal numeric ID", func(t *testing.T) {
		var f FlexibleID
		err := json.Unmarshal([]byte(`456`), &f)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if f.Value != "456" {
			t.Errorf("Expected Value='456', got %q", f.Value)
		}
	})

	t.Run("unmarshal float ID", func(t *testing.T) {
		var f FlexibleID
		err := json.Unmarshal([]byte(`789.0`), &f)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		// Accept either "789" or "789.0" — both are reasonable
		if f.Value != "789" && f.Value != "789.0" {
			t.Errorf("Expected Value='789' or '789.0', got %q", f.Value)
		}
	})

	t.Run("marshal always produces string", func(t *testing.T) {
		f := FlexibleID{Value: "test-id"}
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if string(data) != `"test-id"` {
			t.Errorf("Expected '\"test-id\"', got %s", string(data))
		}
	})

	t.Run("flexible ID in struct", func(t *testing.T) {
		type Record struct {
			ID   FlexibleID `json:"id"`
			Name string     `json:"name"`
		}

		// Test with string ID
		input1 := `{"id": "str-1", "name": "Alice"}`
		var r1 Record
		if err := json.Unmarshal([]byte(input1), &r1); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if r1.ID.Value != "str-1" {
			t.Errorf("Expected ID='str-1', got %q", r1.ID.Value)
		}

		// Test with numeric ID
		input2 := `{"id": 42, "name": "Bob"}`
		var r2 Record
		if err := json.Unmarshal([]byte(input2), &r2); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if r2.ID.Value != "42" {
			t.Errorf("Expected ID='42', got %q", r2.ID.Value)
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original := FlexibleID{Value: "round-trip-123"}
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		var restored FlexibleID
		err = json.Unmarshal(data, &restored)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if restored.Value != original.Value {
			t.Errorf("Round trip failed: original=%q, restored=%q", original.Value, restored.Value)
		}
	})
}
