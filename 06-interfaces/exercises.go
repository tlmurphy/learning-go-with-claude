package interfaces

import (
	"fmt"
	"io"
	"strings"
)

/*
=============================================================================
 EXERCISES: Interfaces
=============================================================================

 Work through these exercises in order. Run tests with:

   make test 06

 Run a single test:

   go test -v -run TestColorString ./06-interfaces/

=============================================================================
*/

// =========================================================================
// Exercise 1: Implement fmt.Stringer
// =========================================================================

// Color represents an RGB color.
type Color struct {
	R, G, B uint8
}

// String implements fmt.Stringer for Color.
// Return the color as a hex string: "#rrggbb"
// Example: Color{255, 0, 128} -> "#ff0080"
// Use fmt.Sprintf with "%02x" format for each component.
func (c Color) String() string {
	// YOUR CODE HERE
	return ""
}

// Book represents a book with a title and author.
type Book struct {
	Title  string
	Author string
	Pages  int
}

// String implements fmt.Stringer for Book.
// Format: "Title by Author (Pages pages)"
// Example: "The Go Programming Language by Donovan & Kernighan (380 pages)"
func (b Book) String() string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 2: Shape Interface — Multiple Implementations
// =========================================================================

// Shape defines the behavior of a geometric shape.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rect represents a rectangle.
type Rect struct {
	Width, Height float64
}

// Area returns the area of the rectangle.
func (r Rect) Area() float64 {
	// YOUR CODE HERE
	return 0
}

// Perimeter returns the perimeter of the rectangle.
func (r Rect) Perimeter() float64 {
	// YOUR CODE HERE
	return 0
}

// Triangle represents a triangle with three sides and a height relative
// to the base (side A) for area calculation.
type Triangle struct {
	SideA, SideB, SideC float64
	Height              float64 // height relative to SideA (base)
}

// Area returns the area of the triangle (0.5 * base * height).
func (t Triangle) Area() float64 {
	// YOUR CODE HERE
	return 0
}

// Perimeter returns the perimeter of the triangle (sum of all sides).
func (t Triangle) Perimeter() float64 {
	// YOUR CODE HERE
	return 0
}

// CircleShape represents a circle with a given radius.
// (Named CircleShape to avoid collision with the lesson's Circle type.)
type CircleShape struct {
	Radius float64
}

// Area returns the area of the circle (pi * r^2).
// Use math.Pi from the "math" package.
func (c CircleShape) Area() float64 {
	// YOUR CODE HERE
	return 0
}

// Perimeter returns the circumference of the circle (2 * pi * r).
func (c CircleShape) Perimeter() float64 {
	// YOUR CODE HERE
	return 0
}

// TotalArea takes a slice of Shapes and returns the sum of all their areas.
// This demonstrates the power of interfaces — one function works for all shapes.
func TotalArea(shapes []Shape) float64 {
	// YOUR CODE HERE
	return 0
}

// =========================================================================
// Exercise 3: Custom io.Reader — ROT13 Reader
// =========================================================================

// ROT13Reader wraps another io.Reader and applies ROT13 transformation
// to all alphabetic characters as they are read.
//
// ROT13 shifts each letter by 13 positions:
//
//	'A' -> 'N', 'B' -> 'O', ..., 'N' -> 'A', etc.
//	'a' -> 'n', 'b' -> 'o', ..., 'n' -> 'a', etc.
//
// Non-alphabetic characters pass through unchanged.
//
// This exercise teaches you to implement io.Reader, one of Go's most
// important interfaces. HTTP request bodies, file contents, and compressed
// streams are all io.Readers.
type ROT13Reader struct {
	reader io.Reader
}

// NewROT13Reader creates a ROT13Reader wrapping the given reader.
func NewROT13Reader(r io.Reader) *ROT13Reader {
	// YOUR CODE HERE
	return nil
}

// Read implements io.Reader. Read from the underlying reader, then apply
// ROT13 to each byte before returning.
//
// Steps:
//  1. Read from r.reader into p
//  2. For each byte read, apply ROT13 transformation
//  3. Return the number of bytes read and any error
func (r *ROT13Reader) Read(p []byte) (int, error) {
	// YOUR CODE HERE
	return 0, nil
}

// rot13 transforms a single byte using ROT13.
// Letters a-z and A-Z are rotated; everything else is unchanged.
func rot13(b byte) byte {
	// YOUR CODE HERE
	return b
}

// =========================================================================
// Exercise 4: Type Switch — JSON-like Value Handler
// =========================================================================

// Describe takes any value and returns a human-readable description
// of its type and value. Handle these types:
//
//   - int:       "integer: 42"
//   - float64:   "float: 3.14"
//   - string:    "string: hello (length: 5)"
//   - bool:      "boolean: true" or "boolean: false"
//   - []int:     "int slice: [1 2 3] (length: 3)"
//   - nil:       "nil value"
//   - any other: "unknown type: <type>"
//
// Use a type switch — it's the idiomatic way to handle this.
func Describe(value any) string {
	// YOUR CODE HERE
	_ = fmt.Sprintf // hint
	return ""
}

// Summarize takes a slice of any values and returns a map counting how
// many of each type are present. The keys should be:
// "int", "float64", "string", "bool", "nil", "other"
func Summarize(values []any) map[string]int {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 5: Interface Composition
// =========================================================================

// Reader can read data
type Reader interface {
	Read() (string, error)
}

// Writer can write data
type Writer interface {
	Write(data string) error
}

// Closer can be closed
type Closer interface {
	Close() error
}

// ReadWriter composes Reader and Writer
type ReadWriter interface {
	Reader
	Writer
}

// ReadWriteCloser composes all three
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// Buffer is an in-memory implementation that satisfies ReadWriteCloser.
// It stores data as a string. Read returns all stored data and clears it.
// Write appends to the stored data. Close marks it as closed.
type Buffer struct {
	data   string
	closed bool
}

// NewBuffer creates a new empty Buffer.
func NewBuffer() *Buffer {
	// YOUR CODE HERE
	return nil
}

// Read returns all data in the buffer and clears it.
// Returns an error if the buffer is closed.
// If empty, return ("", nil).
func (b *Buffer) Read() (string, error) {
	// YOUR CODE HERE
	return "", nil
}

// Write appends data to the buffer.
// Returns an error if the buffer is closed.
func (b *Buffer) Write(data string) error {
	// YOUR CODE HERE
	return nil
}

// Close marks the buffer as closed. Further Read/Write calls should error.
// Closing an already-closed buffer returns an error.
func (b *Buffer) Close() error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 6: Middleware Chain (Handler Pattern)
// =========================================================================

// Handler processes a request string and returns a response string.
// This is a simplified version of http.Handler.
type Handler interface {
	Handle(request string) string
}

// HandlerFunc is a function type that implements Handler.
// This is the same pattern as http.HandlerFunc.
type HandlerFunc func(request string) string

// Handle calls the function itself — this makes any function a Handler.
func (f HandlerFunc) Handle(request string) string {
	return f(request)
}

// Middleware is a function that wraps a Handler and returns a new Handler.
type Middleware func(Handler) Handler

// LoggingMiddleware returns a Middleware that prefixes the response with
// "[LOG] " and the original request.
// Format: "[LOG] request=<request> response=<response>"
//
// Example: if the inner handler returns "OK" for request "GET /api",
// the logging middleware should return "[LOG] request=GET /api response=OK"
func LoggingMiddleware() Middleware {
	// YOUR CODE HERE
	// Return a function that takes a Handler and returns a Handler
	return nil
}

// AuthMiddleware returns a Middleware that checks if the request contains
// the word "authorized". If not, return "401 Unauthorized" without calling
// the next handler. If authorized, strip "authorized " from the front of
// the request and pass the rest to the next handler.
func AuthMiddleware() Middleware {
	// YOUR CODE HERE
	return nil
}

// Chain applies middlewares to a handler in order (first middleware is outermost).
// Chain(handler, mw1, mw2) means: mw1 wraps mw2 wraps handler
// So the request flows: mw1 -> mw2 -> handler -> mw2 -> mw1
func Chain(h Handler, middlewares ...Middleware) Handler {
	// YOUR CODE HERE
	return h
}

// =========================================================================
// Exercise 7: Dependency Injection — Storage Interface
// =========================================================================

// Item represents a stored item with an ID and data.
type Item struct {
	ID   string
	Data string
}

// Storage defines the interface for persisting items.
// In a real app, you'd have database and cache implementations.
// For testing, you'd use the in-memory version.
type Storage interface {
	Get(id string) (Item, bool)
	Put(item Item) error
	Delete(id string) error
	List() []Item
}

// MemoryStorage implements Storage using an in-memory map.
type MemoryStorage struct {
	items map[string]Item
}

// NewMemoryStorage creates a new MemoryStorage ready for use.
func NewMemoryStorage() *MemoryStorage {
	// YOUR CODE HERE
	return nil
}

// Get retrieves an item by ID. Returns the item and true if found,
// or zero Item and false if not found.
func (ms *MemoryStorage) Get(id string) (Item, bool) {
	// YOUR CODE HERE
	return Item{}, false
}

// Put stores an item. If an item with the same ID exists, it's replaced.
// Return an error if the ID is empty.
func (ms *MemoryStorage) Put(item Item) error {
	// YOUR CODE HERE
	return nil
}

// Delete removes an item by ID. Return an error if the ID is not found.
func (ms *MemoryStorage) Delete(id string) error {
	// YOUR CODE HERE
	return nil
}

// List returns all items in storage (order not guaranteed).
func (ms *MemoryStorage) List() []Item {
	// YOUR CODE HERE
	return nil
}

// ItemService uses a Storage interface — it doesn't know or care whether
// the storage is in-memory, a database, or anything else.
type ItemService struct {
	storage Storage
}

// NewItemService creates a service backed by the given storage.
// This is dependency injection — the storage is "injected" by the caller.
func NewItemService(storage Storage) *ItemService {
	// YOUR CODE HERE
	return nil
}

// GetItem retrieves an item, returning an error if not found.
func (s *ItemService) GetItem(id string) (Item, error) {
	// YOUR CODE HERE
	return Item{}, nil
}

// SaveItem stores an item.
func (s *ItemService) SaveItem(item Item) error {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 8: The Nil Interface Gotcha
// =========================================================================

// Validator is an interface for things that can validate themselves.
type Validator interface {
	Validate() error
}

// EmailValidator validates an email address.
type EmailValidator struct {
	Email string
}

// Validate checks that the email contains "@".
// Return an error (using fmt.Errorf) if it doesn't.
func (e *EmailValidator) Validate() error {
	// YOUR CODE HERE
	return nil
}

// GetValidator returns a Validator based on the validatorType string.
//
// IMPORTANT: This exercise demonstrates the nil interface gotcha.
//
// If validatorType is "email", return an EmailValidator with the given value.
// For any other type, return nil (the interface nil, NOT a typed nil pointer).
//
// WRONG (returns non-nil interface holding nil pointer):
//
//	var v *EmailValidator = nil
//	return v  // interface is NOT nil! It holds (*EmailValidator, nil)
//
// RIGHT (returns truly nil interface):
//
//	return nil  // interface is nil — both type and value are nil
func GetValidator(validatorType string, value string) Validator {
	// YOUR CODE HERE
	return nil
}

// IsNilInterface checks whether the given interface value is truly nil.
// This is a helper to demonstrate the concept — in practice you'd just
// use == nil, but this shows you understand the gotcha.
//
// Return true only when the interface itself is nil (no type, no value).
// An interface holding a nil pointer is NOT nil.
func IsNilInterface(v any) bool {
	// YOUR CODE HERE
	return false
}

// Ensure unused imports don't cause compile errors
var _ = strings.Contains
var _ = fmt.Sprintf
