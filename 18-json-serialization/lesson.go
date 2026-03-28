package jsonserialization

/*
=============================================================================
 Module 18: JSON and Serialization
=============================================================================

 JSON is the lingua franca of web APIs. Every REST endpoint you build will
 marshal Go structs into JSON responses and unmarshal JSON request bodies
 into Go structs. Getting this right is essential — bugs in serialization
 are insidious because they often manifest as subtle data corruption or
 silent data loss rather than loud errors.

 Go's encoding/json package uses reflection to map between JSON and Go
 structs. It's well-designed, but it has sharp edges you need to know about.

 WHY THIS MATTERS FOR WEB SERVICES:
 - Every API response is a JSON marshal operation
 - Every request body parse is a JSON unmarshal operation
 - Incorrect struct tags silently drop fields
 - Missing omitempty causes null-vs-absent confusion for API consumers
 - Custom marshaling lets you present clean APIs over messy internal types

 NOTE: Go 1.25 introduces an experimental encoding/json/v2 with better
 performance, stricter defaults, and more control. Keep an eye on it, but
 encoding/json (v1) will remain the standard for years to come.

=============================================================================
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// -------------------------------------------------------------------------
// Marshal and Unmarshal Basics
// -------------------------------------------------------------------------

/*
 json.Marshal converts a Go value to a JSON byte slice.
 json.Unmarshal converts a JSON byte slice back to a Go value.

 The key rule: only EXPORTED fields (uppercase) are visible to encoding/json.
 Unexported fields are completely invisible — they won't appear in JSON
 output and they'll be silently ignored during parsing. This is one of the
 most common gotchas for Go beginners.

 json.Marshal NEVER returns invalid JSON — if it succeeds, the output
 is valid. If it fails (e.g., you pass a channel or function), it returns
 an error.
*/

type Product struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

func DemoMarshalUnmarshal() {
	// Marshal: Go struct -> JSON bytes
	p := Product{ID: 1, Name: "Widget", Price: 9.99, InStock: true}
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	fmt.Println("Marshaled:", string(data))
	// Output: {"id":1,"name":"Widget","price":9.99,"in_stock":true}

	// Unmarshal: JSON bytes -> Go struct
	jsonStr := `{"id":2,"name":"Gadget","price":24.99,"in_stock":false}`
	var p2 Product
	err = json.Unmarshal([]byte(jsonStr), &p2)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}
	fmt.Printf("Unmarshaled: %+v\n", p2)

	// MarshalIndent for human-readable output (useful for debugging, logs)
	pretty, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println("Pretty:\n", string(pretty))
}

// -------------------------------------------------------------------------
// Struct Tags for JSON
// -------------------------------------------------------------------------

/*
 Struct tags control how fields map to JSON. The syntax is:
   `json:"name"`           — use "name" as the JSON key
   `json:"name,omitempty"` — omit if the field is the zero value
   `json:"-"`              — always omit (never serialize)
   `json:"-,"`             — use literal "-" as the key (rare)
   `json:",omitempty"`     — use Go field name, but omit if zero value

 The omitempty option is critical for APIs. Without it, zero values appear
 in your JSON, which can confuse clients:
   {"email": ""}     — is email empty or was it not provided?
   (vs. no field)    — clearly not provided

 Common pitfall: omitempty considers 0, false, "", nil, and empty
 slices/maps as "empty". This means `json:",omitempty"` on a boolean
 field will omit it when false, which might not be what you want!
*/

type UserProfile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"` // omit if empty string
	Bio      string `json:"bio,omitempty"`   // omit if empty string
	IsAdmin  bool   `json:"is_admin"`        // NOT omitempty — false is meaningful
	Password string `json:"-"`               // never include in JSON output!
}

func DemoStructTags() {
	// With omitempty, empty fields are excluded from output
	u := UserProfile{
		ID:       1,
		Username: "gopher",
		Password: "secret123", // this will never appear in JSON
		IsAdmin:  false,       // this WILL appear because no omitempty
	}
	data, _ := json.MarshalIndent(u, "", "  ")
	fmt.Println("With omitempty:\n", string(data))
	// Note: email and bio are absent, password is absent, is_admin is present (false)
}

// -------------------------------------------------------------------------
// Null vs Absent: Pointer Types for Optional Fields
// -------------------------------------------------------------------------

/*
 One of the trickiest aspects of JSON in Go: distinguishing between a
 field that's absent from the JSON and one that's explicitly null or
 has a zero value.

 Consider: {"name": "Alice"} vs {"name": "Alice", "age": 0}
 With a plain int field, both unmarshal to age=0. You can't tell them
 apart.

 The solution: use pointer types for optional fields.

   *string with omitempty:
     nil     → field omitted from JSON
     ""      → "field": ""  (present but empty)

   *int with omitempty:
     nil     → field omitted
     0       → "field": 0   (present and zero — meaningful!)

 This is essential for PATCH endpoints where you need to distinguish
 "don't change this field" (absent) from "set this to null/zero"
 (explicitly provided).
*/

type UpdateRequest struct {
	Name  *string `json:"name,omitempty"`  // nil = don't change, "" = set empty
	Age   *int    `json:"age,omitempty"`   // nil = don't change, 0 = set to zero
	Email *string `json:"email,omitempty"` // nil = don't change
}

func DemoNullVsAbsent() {
	// Absent fields: pointer stays nil
	input1 := `{"name": "Alice"}`
	var req1 UpdateRequest
	json.Unmarshal([]byte(input1), &req1)
	fmt.Printf("Name: %v, Age: %v\n", req1.Name, req1.Age)
	// Name: &"Alice", Age: <nil> — we know age was not provided

	// Explicit null: pointer is nil too (same as absent for most cases)
	input2 := `{"name": "Alice", "age": null}`
	var req2 UpdateRequest
	json.Unmarshal([]byte(input2), &req2)
	fmt.Printf("Name: %v, Age: %v\n", req2.Name, req2.Age)

	// Explicit zero: pointer is non-nil, points to 0
	input3 := `{"name": "Alice", "age": 0}`
	var req3 UpdateRequest
	json.Unmarshal([]byte(input3), &req3)
	if req3.Age != nil {
		fmt.Printf("Age explicitly set to: %d\n", *req3.Age)
	}
}

// -------------------------------------------------------------------------
// Handling Unknown Fields
// -------------------------------------------------------------------------

/*
 By default, json.Unmarshal silently ignores JSON fields that don't match
 any struct field. This is usually fine — it makes your API forward
 compatible (clients can send extra fields).

 But sometimes you WANT to reject unknown fields — for example, to catch
 typos in configuration files or to enforce strict API contracts.

 Use json.Decoder with DisallowUnknownFields for this.
*/

func DemoDisallowUnknownFields() {
	type StrictConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	input := `{"host": "localhost", "port": 8080, "typo_field": true}`

	// Default: silently ignores unknown fields
	var c1 StrictConfig
	json.Unmarshal([]byte(input), &c1) // no error!
	fmt.Printf("Default (no error): %+v\n", c1)

	// Strict: rejects unknown fields
	var c2 StrictConfig
	dec := json.NewDecoder(strings.NewReader(input))
	dec.DisallowUnknownFields()
	err := dec.Decode(&c2)
	if err != nil {
		fmt.Println("Strict mode error:", err)
		// json: unknown field "typo_field"
	}
}

// -------------------------------------------------------------------------
// Custom MarshalJSON and UnmarshalJSON
// -------------------------------------------------------------------------

/*
 Sometimes the default marshaling isn't enough. You might need to:
   - Serialize cents as "$12.34" for display
   - Validate fields during deserialization
   - Handle legacy formats from external APIs
   - Flatten nested structures

 Implement the json.Marshaler and json.Unmarshaler interfaces:
   type Marshaler interface {
       MarshalJSON() ([]byte, error)
   }
   type Unmarshaler interface {
       UnmarshalJSON([]byte) error
   }

 IMPORTANT: In MarshalJSON, you must return valid JSON. Use json.Marshal
 on a simple type (string, map, etc.) rather than hand-building JSON
 strings — it's easy to forget to quote strings or escape special chars.
*/

// Money stores amount in cents to avoid floating-point precision issues.
// Externally, it serializes as a string like "12.34".
type Money struct {
	Cents int64
}

func (m Money) MarshalJSON() ([]byte, error) {
	// Format cents as dollars.cents string
	dollars := m.Cents / 100
	cents := m.Cents % 100
	if cents < 0 {
		cents = -cents
	}
	s := fmt.Sprintf("%d.%02d", dollars, cents)
	// Use json.Marshal on the string to get proper JSON quoting
	return json.Marshal(s)
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var dollars, cents int64
	_, err := fmt.Sscanf(s, "%d.%d", &dollars, &cents)
	if err != nil {
		return fmt.Errorf("invalid money format: %s", s)
	}

	m.Cents = dollars*100 + cents
	return nil
}

func DemoCustomMarshal() {
	price := Money{Cents: 1234}
	data, _ := json.Marshal(price)
	fmt.Println("Money as JSON:", string(data)) // "12.34"

	var parsed Money
	json.Unmarshal([]byte(`"56.78"`), &parsed)
	fmt.Printf("Parsed money: %d cents\n", parsed.Cents) // 5678
}

// -------------------------------------------------------------------------
// json.RawMessage for Deferred Parsing
// -------------------------------------------------------------------------

/*
 json.RawMessage is a raw encoded JSON value. It lets you delay parsing
 part of a JSON document until you know what type it should be.

 This is incredibly useful for:
   - Heterogeneous arrays/objects (events with different payload types)
   - API responses where the "data" field varies by endpoint
   - Storing JSON blobs without parsing them (pass-through)
   - Two-pass parsing: first pass reads the type, second pass parses data
*/

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // parse later based on Type
}

type ClickPayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type KeyPayload struct {
	Key  string `json:"key"`
	Code int    `json:"code"`
}

func DemoRawMessage() {
	events := `[
		{"type": "click", "payload": {"x": 100, "y": 200}},
		{"type": "key", "payload": {"key": "Enter", "code": 13}}
	]`

	var parsed []Event
	json.Unmarshal([]byte(events), &parsed)

	for _, e := range parsed {
		switch e.Type {
		case "click":
			var p ClickPayload
			json.Unmarshal(e.Payload, &p)
			fmt.Printf("Click at (%d, %d)\n", p.X, p.Y)
		case "key":
			var p KeyPayload
			json.Unmarshal(e.Payload, &p)
			fmt.Printf("Key press: %s (code %d)\n", p.Key, p.Code)
		}
	}
}

// -------------------------------------------------------------------------
// Streaming JSON with Decoder/Encoder
// -------------------------------------------------------------------------

/*
 For large JSON payloads, loading the entire thing into memory (via
 json.Unmarshal) is wasteful. json.Decoder reads from an io.Reader and
 can process tokens incrementally. json.Encoder writes to an io.Writer.

 When to use streaming:
   - Large files or API responses (don't load 100MB of JSON into memory)
   - HTTP request/response bodies (io.Reader/io.Writer)
   - Line-delimited JSON (NDJSON) — one JSON object per line

 In HTTP handlers, always use json.NewDecoder(r.Body) rather than
 reading the entire body into a byte slice first. It's more efficient
 and lets the decoder report errors as they occur.
*/

func DemoStreamingJSON() {
	// Encoder: write JSON directly to an io.Writer
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ") // optional: pretty print

	products := []Product{
		{ID: 1, Name: "Widget", Price: 9.99, InStock: true},
		{ID: 2, Name: "Gadget", Price: 24.99, InStock: false},
	}

	for _, p := range products {
		enc.Encode(p) // writes one JSON object per line
	}
	fmt.Println("Encoded stream:\n", buf.String())

	// Decoder: read JSON from an io.Reader
	reader := strings.NewReader(`{"id":1,"name":"Widget","price":9.99,"in_stock":true}
{"id":2,"name":"Gadget","price":24.99,"in_stock":false}`)

	dec := json.NewDecoder(reader)
	for {
		var p Product
		err := dec.Decode(&p)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("decode error:", err)
			break
		}
		fmt.Printf("Decoded: %+v\n", p)
	}
}

// -------------------------------------------------------------------------
// Time Serialization
// -------------------------------------------------------------------------

/*
 time.Time marshals to RFC 3339 format by default: "2024-01-15T14:30:00Z"
 This is the standard for JSON APIs and is what most clients expect.

 If you need a different format (e.g., Unix timestamp, custom layout),
 you'll need a custom type with MarshalJSON/UnmarshalJSON.
*/

type AuditEntry struct {
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"` // RFC 3339 by default
}

func DemoTimeSerialization() {
	entry := AuditEntry{
		Action:    "user.login",
		Timestamp: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
	}
	data, _ := json.MarshalIndent(entry, "", "  ")
	fmt.Println("Time in JSON:\n", string(data))
	// timestamp will be "2024-01-15T14:30:00Z"
}

// -------------------------------------------------------------------------
// json.Number for Arbitrary Precision
// -------------------------------------------------------------------------

/*
 By default, JSON numbers unmarshal to float64 when decoding into
 interface{}. This can lose precision for large integers:

   {"id": 9007199254740993}  →  float64(9007199254740992)  (precision lost!)

 Use json.Decoder with UseNumber() to keep numbers as strings (json.Number)
 until you explicitly convert them. This is critical when handling IDs or
 financial data where precision matters.
*/

func DemoJSONNumber() {
	input := `{"id": 9007199254740993, "amount": "12.34"}`

	// Default: numbers become float64 (precision risk)
	var m1 map[string]any
	json.Unmarshal([]byte(input), &m1)
	fmt.Printf("Default float64: %v (type: %T)\n", m1["id"], m1["id"])

	// With UseNumber: numbers stay as json.Number strings
	dec := json.NewDecoder(strings.NewReader(input))
	dec.UseNumber()
	var m2 map[string]any
	dec.Decode(&m2)
	if num, ok := m2["id"].(json.Number); ok {
		// Convert to int64 with full precision
		id, _ := num.Int64()
		fmt.Printf("json.Number: %d (preserved!)\n", id)
	}
}

// -------------------------------------------------------------------------
// Maps and Generic JSON
// -------------------------------------------------------------------------

/*
 When you don't know the JSON structure at compile time, unmarshal into
 map[string]interface{} (or map[string]any in Go 1.18+).

 The type mappings are:
   JSON object  → map[string]interface{}
   JSON array   → []interface{}
   JSON string  → string
   JSON number  → float64 (or json.Number with UseNumber)
   JSON boolean → bool
   JSON null    → nil

 This is useful for:
   - Proxying JSON (receive and forward without parsing)
   - Dynamic configuration
   - Exploring unknown API responses

 But prefer typed structs when you know the structure — they're safer,
 faster, and self-documenting.
*/

func DemoMapJSON() {
	// Unmarshal into a map when structure is unknown
	input := `{"name": "Alice", "scores": [95, 87, 92], "active": true}`
	var m map[string]any
	json.Unmarshal([]byte(input), &m)

	name := m["name"].(string)
	scores := m["scores"].([]any)
	fmt.Printf("Name: %s, First score: %v\n", name, scores[0])

	// Marshal from a map (useful for building dynamic JSON)
	response := map[string]any{
		"status": "ok",
		"count":  42,
		"items":  []string{"a", "b", "c"},
	}
	data, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(data))
}

// -------------------------------------------------------------------------
// Performance Considerations
// -------------------------------------------------------------------------

/*
 encoding/json uses reflection, which has overhead. For most web services,
 this is fine — network I/O dwarfs JSON processing time. But for hot paths
 processing thousands of objects per second, consider:

 1. Use json.Decoder/Encoder for streaming (avoid intermediate []byte)
 2. Pre-allocate slices when you know the approximate size
 3. For extreme performance, consider third-party libraries:
    - github.com/json-iterator/go (drop-in replacement, 2-3x faster)
    - github.com/goccy/go-json (even faster, also drop-in)
    - github.com/bytedance/sonic (fastest, but uses code generation)
 4. Avoid marshaling in hot loops — marshal once, cache if possible
 5. The upcoming encoding/json/v2 in Go 1.25 aims to be significantly
    faster while remaining in the standard library

 For most applications: use standard encoding/json. Optimize only when
 profiling shows JSON processing is a bottleneck.
*/

// -------------------------------------------------------------------------
// Common Gotchas Summary
// -------------------------------------------------------------------------

/*
 1. UNEXPORTED FIELDS ARE INVISIBLE
    type user struct { name string }
    json.Marshal(user{name: "hi"}) → {}  // empty! name is lowercase

 2. MAPS SERIALIZE TO OBJECTS (keys become strings)
    map[int]string{1: "one"} → {"1": "one"}  // int keys become strings

 3. NIL SLICE vs EMPTY SLICE
    var s []int = nil     → null in JSON
    s := []int{}          → [] in JSON
    Always initialize slices you'll marshal: make([]T, 0)

 4. INTERFACE{} VALUES UNMARSHAL AS FLOAT64 FOR NUMBERS
    var v interface{}; json.Unmarshal([]byte("42"), &v)
    v is float64(42), not int(42)

 5. EMBEDDED STRUCTS CAN CAUSE FIELD CONFLICTS
    If an embedded struct and the outer struct both have a field with
    the same JSON name, behavior depends on nesting depth.

 6. MARSHAL ERRORS ARE RARE BUT REAL
    Channels, functions, and complex numbers can't be marshaled.
    Always check the error from json.Marshal in production code.
*/
