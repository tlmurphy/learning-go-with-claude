package structs

import "fmt"

/*
=============================================================================
 EXERCISES: Structs and Methods
=============================================================================

 Work through these exercises in order. Run tests with:

   make test 05

 Run a single test:

   go test -v -run TestNewUser ./05-structs-and-methods/

=============================================================================
*/

// =========================================================================
// Exercise 1: User Struct with Constructor and Methods
// =========================================================================

// User represents a person in our system. Define it with the following fields:
//   - FirstName string
//   - LastName  string
//   - Email     string
//   - Age       int
type User struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
}

// NewUser creates a new User with the given details.
// Return a pointer to the User.
func NewUser(firstName, lastName, email string, age int) *User {
	// YOUR CODE HERE
	return nil
}

// FullName returns the user's first and last name separated by a space.
// This should use a value receiver since it doesn't modify the User.
func (u User) FullName() string {
	// YOUR CODE HERE
	return ""
}

// IsAdult returns true if the user is 18 or older.
func (u User) IsAdult() bool {
	// YOUR CODE HERE
	return false
}

// UpdateEmail changes the user's email address.
// Think carefully: should this be a value or pointer receiver? Why?
func (u *User) UpdateEmail(newEmail string) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 2: Rectangle with Value Receiver Methods
// =========================================================================

// Rectangle represents a rectangle with a width and height.
type Rectangle struct {
	Width  float64
	Height float64
}

// NewRectangle creates a new Rectangle. If width or height is negative,
// use the absolute value (multiply by -1).
func NewRectangle(width, height float64) Rectangle {
	// YOUR CODE HERE
	return Rectangle{}
}

// Area returns the area of the rectangle (width * height).
// Use a value receiver — this is a pure computation, no mutation needed.
func (r Rectangle) Area() float64 {
	// YOUR CODE HERE
	return 0
}

// Perimeter returns the perimeter of the rectangle (2*width + 2*height).
func (r Rectangle) Perimeter() float64 {
	// YOUR CODE HERE
	return 0
}

// IsSquare returns true if the rectangle's width equals its height.
func (r Rectangle) IsSquare() bool {
	// YOUR CODE HERE
	return false
}

// Scale returns a NEW rectangle with width and height multiplied by factor.
// Do not modify the original rectangle.
func (r Rectangle) Scale(factor float64) Rectangle {
	// YOUR CODE HERE
	return Rectangle{}
}

// =========================================================================
// Exercise 3: BankAccount with Pointer Receiver Methods
// =========================================================================

// BankAccount represents a bank account with an owner and balance.
// The balance field should be unexported (lowercase) so it can only
// be modified through the methods.
type BankAccount struct {
	Owner   string
	balance float64
}

// NewBankAccount creates a new account with the given owner and initial balance.
// If initialBalance is negative, set balance to 0.
func NewBankAccount(owner string, initialBalance float64) *BankAccount {
	// YOUR CODE HERE
	return nil
}

// Deposit adds money to the account. The amount must be positive.
// Return an error message string if amount <= 0, or empty string on success.
// This MUST be a pointer receiver — we're modifying state.
func (a *BankAccount) Deposit(amount float64) string {
	// YOUR CODE HERE
	return ""
}

// Withdraw removes money from the account. The amount must be positive
// and must not exceed the current balance.
// Return an error message string if invalid, or empty string on success.
func (a *BankAccount) Withdraw(amount float64) string {
	// YOUR CODE HERE
	return ""
}

// Balance returns the current balance. Even though this doesn't mutate,
// use a pointer receiver for consistency with other BankAccount methods.
func (a *BankAccount) Balance() float64 {
	// YOUR CODE HERE
	return 0
}

// Transfer moves the given amount from this account to another account.
// Return an error message if the transfer fails (same rules as Withdraw).
func (a *BankAccount) Transfer(amount float64, to *BankAccount) string {
	// YOUR CODE HERE
	return ""
}

// =========================================================================
// Exercise 4: Struct Embedding — Admin embeds User
// =========================================================================

// Admin represents an administrator who is also a user. Embed the User
// struct and add these additional fields:
//   - Role        string
//   - Permissions []string
type Admin struct {
	User
	Role        string
	Permissions []string
}

// NewAdmin creates a new Admin with the given user info, role, and permissions.
func NewAdmin(firstName, lastName, email string, age int, role string, permissions []string) *Admin {
	// YOUR CODE HERE
	return nil
}

// HasPermission checks if the admin has the given permission in their
// Permissions slice. Return false if permissions is nil or empty.
func (a Admin) HasPermission(perm string) bool {
	// YOUR CODE HERE
	return false
}

// Promote adds a permission to the admin's permission list if it doesn't
// already exist.
func (a *Admin) Promote(perm string) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 5: Linked List Node
// =========================================================================

// ListNode represents a node in a singly linked list of integers.
// It should have a Value field (int) and a Next field (*ListNode).
type ListNode struct {
	Value int
	Next  *ListNode
}

// NewLinkedList creates a linked list from a slice of ints and returns
// the head node. If the slice is empty, return nil.
// The list should maintain the order of the input slice.
func NewLinkedList(values []int) *ListNode {
	// YOUR CODE HERE
	return nil
}

// ToSlice converts the linked list starting at this node back into a
// slice of ints. If the receiver is nil, return an empty (non-nil) slice.
func (n *ListNode) ToSlice() []int {
	// YOUR CODE HERE
	return nil
}

// Len returns the number of nodes in the list starting from this node.
func (n *ListNode) Len() int {
	// YOUR CODE HERE
	return 0
}

// Append adds a new node with the given value at the end of the list.
func (n *ListNode) Append(value int) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 6: HTTP-like Request Struct
// =========================================================================

// Header represents a collection of HTTP headers as key-value pairs.
// A header key can have multiple values (e.g., Accept: text/html, Accept: application/json).
type Header map[string][]string

// Request represents a simplified HTTP request.
type Request struct {
	Method  string
	Path    string
	Headers Header
	Body    string
}

// NewRequest creates a new Request with the given method and path.
// Initialize the Headers map so it's ready to use.
func NewRequest(method, path string) *Request {
	// YOUR CODE HERE
	return nil
}

// AddHeader adds a header value for the given key. Headers can have
// multiple values for the same key.
func (r *Request) AddHeader(key, value string) {
	// YOUR CODE HERE
}

// GetHeader returns the first value for the given header key, or empty
// string if the header doesn't exist.
func (r *Request) GetHeader(key string) string {
	// YOUR CODE HERE
	return ""
}

// GetAllHeaders returns all values for the given header key.
// Return nil if the key doesn't exist.
func (r *Request) GetAllHeaders(key string) []string {
	// YOUR CODE HERE
	return nil
}

// IsSecure returns true if the request path starts with "https://".
func (r *Request) IsSecure() bool {
	// YOUR CODE HERE
	return false
}

// =========================================================================
// Exercise 7: Implement fmt.Stringer
// =========================================================================

// Temperature represents a temperature value with a unit.
// The Unit field should be 'C' for Celsius or 'F' for Fahrenheit.
type Temperature struct {
	Degrees float64
	Unit    rune // 'C' or 'F'
}

// String implements the fmt.Stringer interface for Temperature.
// Format: "72.0°F" or "22.5°C"
// Use fmt.Sprintf with "%.1f°%c" format.
func (t Temperature) String() string {
	// YOUR CODE HERE
	return ""
}

// ToFahrenheit converts the temperature to Fahrenheit.
// If already Fahrenheit, return a copy unchanged.
// Formula: F = C*9/5 + 32
func (t Temperature) ToFahrenheit() Temperature {
	// YOUR CODE HERE
	return Temperature{}
}

// ToCelsius converts the temperature to Celsius.
// If already Celsius, return a copy unchanged.
// Formula: C = (F-32) * 5/9
func (t Temperature) ToCelsius() Temperature {
	// YOUR CODE HERE
	return Temperature{}
}

// =========================================================================
// Exercise 8: Builder Pattern with Pointer Receiver Chaining
// =========================================================================

// ServerConfig represents a server configuration built using the builder pattern.
type ServerConfig struct {
	Host         string
	Port         int
	TLSEnabled   bool
	ReadTimeout  int
	WriteTimeout int
	MaxConns     int
	LogLevel     string
}

// ServerBuilder is used to construct a ServerConfig step by step.
// Each method should return *ServerBuilder so calls can be chained.
type ServerBuilder struct {
	config ServerConfig
}

// NewServerBuilder creates a new builder with sensible defaults:
//   - Host: "localhost"
//   - Port: 8080
//   - TLSEnabled: false
//   - ReadTimeout: 30
//   - WriteTimeout: 30
//   - MaxConns: 100
//   - LogLevel: "info"
func NewServerBuilder() *ServerBuilder {
	// YOUR CODE HERE
	return nil
}

// WithHost sets the host and returns the builder for chaining.
func (b *ServerBuilder) WithHost(host string) *ServerBuilder {
	// YOUR CODE HERE
	return b
}

// WithPort sets the port and returns the builder for chaining.
func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
	// YOUR CODE HERE
	return b
}

// WithTLS enables or disables TLS and returns the builder for chaining.
func (b *ServerBuilder) WithTLS(enabled bool) *ServerBuilder {
	// YOUR CODE HERE
	return b
}

// WithTimeouts sets both read and write timeouts and returns the builder.
func (b *ServerBuilder) WithTimeouts(read, write int) *ServerBuilder {
	// YOUR CODE HERE
	return b
}

// WithMaxConns sets the maximum connections and returns the builder.
func (b *ServerBuilder) WithMaxConns(max int) *ServerBuilder {
	// YOUR CODE HERE
	return b
}

// WithLogLevel sets the log level and returns the builder.
func (b *ServerBuilder) WithLogLevel(level string) *ServerBuilder {
	// YOUR CODE HERE
	return b
}

// Build returns the final ServerConfig.
func (b *ServerBuilder) Build() ServerConfig {
	// YOUR CODE HERE
	return ServerConfig{}
}

// String implements fmt.Stringer for ServerConfig.
// Format: "host:port (TLS: enabled/disabled)"
func (c ServerConfig) String() string {
	// YOUR CODE HERE
	_ = fmt.Sprintf // hint: use fmt.Sprintf
	return ""
}
