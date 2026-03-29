package jsonserialization

import (
	"encoding/json"
	"fmt"
	"io"
)

/*
=============================================================================
 EXERCISES: JSON and Serialization
=============================================================================

 Work through these exercises in order. Each one builds on concepts from
 the lesson. Run the tests with:

   make test 18

 Tip: Run a single test at a time while working:

   go test -v -run TestBasicMarshalUnmarshal ./18-json-serialization/

=============================================================================
*/

// =========================================================================
// Exercise 1: Basic Struct with JSON Tags
// =========================================================================

// Book represents a book in a library system. Define it with these fields
// and appropriate JSON tags:
//   - ID       int       → json key "id"
//   - Title    string    → json key "title"
//   - Author   string    → json key "author"
//   - Pages    int       → json key "pages", omit if zero
//   - ISBN     string    → json key "isbn", omit if empty
//   - internal string    → should NEVER appear in JSON (unexported is fine,
//     but also add json:"-" tag for documentation)
type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Pages    int    `json:"pages,omitempty"` // YOUR CODE HERE — add the right tag
	ISBN     string `json:"isbn,omitempty"`  // YOUR CODE HERE — add the right tag
	internal string `json:"-"`               // YOUR CODE HERE — add the right tag
}

// MarshalBook takes a Book and returns its JSON representation as a string.
// Return the JSON string and any error from marshaling.
func MarshalBook(b Book) (string, error) {
	// YOUR CODE HERE
	return "", nil
}

// UnmarshalBook takes a JSON string and returns a Book.
// Return the Book and any error from unmarshaling.
func UnmarshalBook(jsonStr string) (Book, error) {
	// YOUR CODE HERE
	return Book{}, nil
}

// =========================================================================
// Exercise 2: Optional Fields with Pointer Types
// =========================================================================

// UserUpdate represents a partial update request for a user profile.
// All fields are optional (pointer types) — nil means "don't change".
// A non-nil pointer means "set to this value" (even if that value is
// zero/empty).
//
// Fields:
//   - Name     *string  → json key "name", omitempty
//   - Age      *int     → json key "age", omitempty
//   - Email    *string  → json key "email", omitempty
//   - Active   *bool    → json key "active", omitempty
type UserUpdate struct {
	Name   *string `json:"name,omitempty"`   // YOUR CODE HERE
	Age    *int    `json:"age,omitempty"`    // YOUR CODE HERE
	Email  *string `json:"email,omitempty"`  // YOUR CODE HERE
	Active *bool   `json:"active,omitempty"` // YOUR CODE HERE
}

// ApplyUpdate takes an existing user (as a map) and a UserUpdate, and
// returns a new map with the updates applied. Only non-nil fields in
// the update should modify the map.
//
// The input map has string keys matching the JSON field names:
// "name", "age", "email", "active"
//
// Example:
//
//	existing: {"name": "Alice", "age": 30, "email": "a@b.com", "active": true}
//	update:   UserUpdate{Name: ptr("Bob")}   (only Name is non-nil)
//	result:   {"name": "Bob", "age": 30, "email": "a@b.com", "active": true}
func ApplyUpdate(existing map[string]interface{}, update UserUpdate) map[string]interface{} {
	// YOUR CODE HERE
	return existing
}

// =========================================================================
// Exercise 3: Custom MarshalJSON — Money Type
// =========================================================================

// Price stores a monetary amount in cents to avoid floating-point issues.
// When marshaled to JSON, it should appear as a string formatted like
// "12.34" (dollars and cents, always two decimal places).
// When unmarshaled, it should parse that string format back to cents.
//
// Examples:
//
//	Price{Cents: 1299}  → JSON: "12.99"
//	Price{Cents: 500}   → JSON: "5.00"
//	Price{Cents: 7}     → JSON: "0.07"
//	JSON "25.50"        → Price{Cents: 2550}
type Price struct {
	Cents int64
}

// MarshalJSON implements json.Marshaler for Price.
// Format the cents as "dollars.cents" with exactly two decimal places.
func (p Price) MarshalJSON() ([]byte, error) {
	// YOUR CODE HERE
	return nil, nil
}

// UnmarshalJSON implements json.Unmarshaler for Price.
// Parse a JSON string like "12.34" into cents (1234).
// Return an error if the format is invalid.
func (p *Price) UnmarshalJSON(data []byte) error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 4: Custom UnmarshalJSON with Validation
// =========================================================================

// Rating represents a product rating from 1 to 5.
// During unmarshaling, it should validate that the value is in range.
// If the value is out of range, return a descriptive error.
type Rating struct {
	Value int
}

// UnmarshalJSON implements json.Unmarshaler for Rating.
// It should:
//   - Parse the JSON number
//   - Validate that it's between 1 and 5 (inclusive)
//   - Return a descriptive error if out of range
func (r *Rating) UnmarshalJSON(data []byte) error {
	// YOUR CODE HERE
	return nil
}

// MarshalJSON implements json.Marshaler for Rating.
// Marshal the Value as a plain JSON number.
func (r Rating) MarshalJSON() ([]byte, error) {
	// YOUR CODE HERE
	return nil, nil
}

// =========================================================================
// Exercise 5: Partial Parsing with json.RawMessage
// =========================================================================

// Notification represents a notification with a type-dependent payload.
// The Type field determines how to parse the Payload.
//
// Types and their payload structures:
//
//	"email"   → EmailPayload{To: string, Subject: string}
//	"sms"     → SMSPayload{Phone: string, Message: string}
//	"push"    → PushPayload{DeviceID: string, Title: string, Body: string}
type Notification struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
}

type SMSPayload struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type PushPayload struct {
	DeviceID string `json:"device_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

// ParseNotification takes a JSON string containing a notification and
// returns the type and the parsed payload as an interface{}.
//
// Based on the "type" field, parse the payload into the correct struct:
//
//	"email" → *EmailPayload
//	"sms"   → *SMSPayload
//	"push"  → *PushPayload
//
// Return an error for unknown types or invalid JSON.
func ParseNotification(jsonStr string) (string, interface{}, error) {
	// YOUR CODE HERE
	return "", nil, nil
}

// =========================================================================
// Exercise 6: Streaming JSON Encoder
// =========================================================================

// LogEntry represents a structured log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// WriteJSONStream writes a slice of LogEntry values to the given io.Writer
// as newline-delimited JSON (one JSON object per line, no array brackets).
// This is the NDJSON format, commonly used for log shipping and streaming.
//
// Each entry should be written using json.Encoder (not json.Marshal).
// Return the number of entries written and any error encountered.
func WriteJSONStream(w io.Writer, entries []LogEntry) (int, error) {
	// YOUR CODE HERE
	return 0, nil
}

// ReadJSONStream reads newline-delimited JSON from the given io.Reader
// and returns a slice of LogEntry values.
// Stop reading when io.EOF is reached.
// Return all successfully parsed entries and any non-EOF error.
func ReadJSONStream(r io.Reader) ([]LogEntry, error) {
	// YOUR CODE HERE
	return nil, nil
}

// =========================================================================
// Exercise 7: API Response Envelope
// =========================================================================

// APIResponse is a generic response envelope for API endpoints.
// It wraps any data type with metadata.
//
// JSON structure:
//
//	{
//	  "data": <whatever T is>,
//	  "meta": {
//	    "request_id": "abc-123",
//	    "count": 42
//	  }
//	}
type Meta struct {
	RequestID string `json:"request_id"`
	Count     int    `json:"count,omitempty"`
}

type APIEnvelope struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// NewAPIResponse creates an APIEnvelope wrapping the given data with
// the provided request ID. If data is a slice, set Meta.Count to the
// length of the slice. For non-slice data, Count should be 0 (omitted).
//
// Hint: Use a type switch or reflect to check if data is a slice.
// For simplicity, you can accept that Count will only be set for
// []interface{} or known slice types passed as interface{}.
// A simpler approach: accept an explicit count parameter.
func NewAPIResponse(data interface{}, requestID string, count int) APIEnvelope {
	// YOUR CODE HERE
	return APIEnvelope{}
}

// MarshalAPIResponse marshals an APIEnvelope to a JSON string.
// Use MarshalIndent with two-space indentation for readability.
func MarshalAPIResponse(env APIEnvelope) (string, error) {
	// YOUR CODE HERE
	return "", nil
}

// =========================================================================
// Exercise 8: Flexible JSON Field (string or number)
// =========================================================================

// FlexibleID can be either a string or a number in JSON.
// Many external APIs are inconsistent about ID types — sometimes they
// send "123" and sometimes 123. This type handles both.
//
// Internally, always store as a string.
// When marshaling, always output as a string.
type FlexibleID struct {
	Value string
}

// UnmarshalJSON implements json.Unmarshaler for FlexibleID.
// It should handle both JSON strings ("123") and JSON numbers (123),
// converting numbers to their string representation.
func (f *FlexibleID) UnmarshalJSON(data []byte) error {
	// YOUR CODE HERE
	// Hint: Try unmarshaling as a string first. If that fails, try as a
	// json.Number or float64 and convert to string.
	_ = fmt.Sprintf // hint: may be useful for number-to-string conversion
	return nil
}

// MarshalJSON implements json.Marshaler for FlexibleID.
// Always marshal as a JSON string.
func (f FlexibleID) MarshalJSON() ([]byte, error) {
	// YOUR CODE HERE
	return nil, nil
}
