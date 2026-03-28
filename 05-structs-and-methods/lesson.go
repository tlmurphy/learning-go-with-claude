package structs

/*
=============================================================================
 Module 05: Structs and Methods
=============================================================================

 Structs are Go's primary mechanism for grouping related data together.
 If you're coming from an object-oriented language, structs fill the role
 of classes — but without inheritance. This isn't a limitation; it's a
 deliberate design choice that leads to simpler, more composable code.

 Go's philosophy: composition over inheritance. You'll see why as we go.

 WHY STRUCTS MATTER FOR WEB SERVICES:
 Every HTTP request you handle, every database row you read, every JSON
 payload you parse — these all become structs. Structs are the backbone
 of your data model. Understanding them deeply is essential.

=============================================================================
*/

import "fmt"

// -------------------------------------------------------------------------
// Struct Definition and Initialization
// -------------------------------------------------------------------------

// A struct groups named fields together. Each field has a name and a type.
// By convention, exported structs and fields start with an uppercase letter.
type Server struct {
	Host     string
	Port     int
	Protocol string
	TLS      bool
}

// There are several ways to create a struct value:
func DemoStructCreation() {
	// 1. Named field initialization (preferred — clear and order-independent)
	s1 := Server{
		Host:     "localhost",
		Port:     8080,
		Protocol: "http",
		TLS:      false,
	}

	// 2. Positional initialization (fragile — breaks if fields are reordered)
	// Avoid this in production code. It's a maintenance hazard.
	s2 := Server{"localhost", 443, "https", true}

	// 3. Partial initialization — unspecified fields get zero values
	s3 := Server{Host: "example.com"}
	// s3.Port == 0, s3.Protocol == "", s3.TLS == false

	// 4. Zero value struct — every field is its zero value
	var s4 Server
	// s4.Host == "", s4.Port == 0, etc.

	fmt.Println(s1, s2, s3, s4)
}

// -------------------------------------------------------------------------
// Zero Value Structs
// -------------------------------------------------------------------------

/*
 One of Go's best features: the zero value is useful. A well-designed struct
 should be usable without explicit initialization. The standard library is
 full of examples:

   var buf bytes.Buffer    // ready to use, no constructor needed
   var mu sync.Mutex       // ready to use
   var wg sync.WaitGroup   // ready to use

 When designing your own structs, aim for useful zero values. This makes
 your API easier to use and reduces the chance of nil pointer panics.
*/

// Counter is usable immediately without initialization.
// Its zero value (count: 0) is perfectly valid.
type Counter struct {
	count int
}

func DemoZeroValues() {
	var c Counter        // count is 0 — a perfectly valid counter
	fmt.Println(c.count) // 0
}

// -------------------------------------------------------------------------
// Anonymous Structs
// -------------------------------------------------------------------------

/*
 Anonymous structs are struct types defined inline, without giving them a
 name. They're surprisingly useful in two scenarios:

 1. Table-driven tests (you'll see this constantly in Go codebases)
 2. One-off JSON responses where defining a named type is overkill

 Don't overuse them — if you use the same shape more than once, give it
 a name. But for ephemeral, single-use structures, they're great.
*/

func DemoAnonymousStructs() {
	// Great for test cases
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"zero", 0, 0},
		{"positive", 5, 25},
		{"negative", -3, 9},
	}

	for _, tt := range tests {
		fmt.Printf("Test %s: input=%d, expected=%d\n", tt.name, tt.input, tt.expected)
	}

	// Great for quick JSON shaping in HTTP handlers
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "ok",
		Message: "request processed",
	}
	fmt.Println(response)
}

// -------------------------------------------------------------------------
// Struct Embedding (Composition Over Inheritance)
// -------------------------------------------------------------------------

/*
 Go doesn't have inheritance. Instead, it has embedding — you can include
 one struct type inside another, and the inner type's fields and methods
 get "promoted" to the outer type. This is composition, not inheritance:

 - There's no "is-a" relationship. An Admin is not a "kind of" User.
 - The embedded type is just a field (you can still access it directly).
 - Method promotion is syntactic sugar, not polymorphism.

 This keeps things simple. You never have to reason about complex
 inheritance hierarchies or worry about the fragile base class problem.
*/

type Person struct {
	FirstName string
	LastName  string
	Email     string
}

func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// Employee embeds Person — it gets all of Person's fields and methods.
type Employee struct {
	Person            // Embedded (not named "Person Person" — that's different)
	EmployeeID string
	Department string
}

func DemoEmbedding() {
	emp := Employee{
		Person: Person{
			FirstName: "Ada",
			LastName:  "Lovelace",
			Email:     "ada@example.com",
		},
		EmployeeID: "EMP001",
		Department: "Engineering",
	}

	// Promoted fields — access Person's fields directly
	fmt.Println(emp.FirstName) // "Ada" (promoted from Person)
	fmt.Println(emp.Email)     // "ada@example.com" (promoted)

	// Promoted methods — call Person's methods directly
	fmt.Println(emp.FullName()) // "Ada Lovelace" (promoted)

	// You can still access the embedded struct explicitly
	fmt.Println(emp.Person.FirstName) // "Ada" (explicit access)
}

// -------------------------------------------------------------------------
// Methods: Value Receivers vs Pointer Receivers
// -------------------------------------------------------------------------

/*
 THIS IS ONE OF THE MOST IMPORTANT CONCEPTS IN GO.

 A method is just a function with a special "receiver" argument. The
 receiver determines which type the method belongs to.

 Value receiver:  func (t Type) Method()   — gets a COPY of t
 Pointer receiver: func (t *Type) Method() — gets a POINTER to t

 The choice between them has real consequences:

 VALUE RECEIVER:
  - Cannot modify the original struct (works on a copy)
  - Safe for concurrent use (each goroutine gets its own copy)
  - Good for small, immutable types (coordinates, colors, money)

 POINTER RECEIVER:
  - CAN modify the original struct
  - Avoids copying large structs (performance)
  - Required when the method needs to mutate state
  - If any method needs a pointer receiver, ALL methods should use one
    (for consistency and to satisfy interfaces correctly)

 RULE OF THUMB: When in doubt, use a pointer receiver. It's almost always
 what you want for structs that represent stateful objects.
*/

// Point is small and immutable — value receivers are fine.
type Point struct {
	X, Y float64
}

// Value receiver — doesn't modify the point, returns a new one.
func (p Point) Translate(dx, dy float64) Point {
	return Point{X: p.X + dx, Y: p.Y + dy}
}

// Account has mutable state — pointer receivers are necessary.
type Account struct {
	Owner   string
	balance float64 // unexported — can only be modified through methods
}

// Pointer receiver — modifies the account's balance in place.
func (a *Account) Deposit(amount float64) {
	if amount > 0 {
		a.balance += amount
	}
}

// Pointer receiver — must be consistent since Deposit uses one.
func (a *Account) Balance() float64 {
	return a.balance
}

func DemoMethods() {
	// Value receiver: the original point is unchanged
	p := Point{1, 2}
	p2 := p.Translate(3, 4)
	fmt.Println(p)  // {1 2} — unchanged
	fmt.Println(p2) // {4 6} — new point

	// Pointer receiver: the original account IS changed
	acct := &Account{Owner: "Alice"}
	acct.Deposit(100)
	fmt.Printf("%s's balance: $%.2f\n", acct.Owner, acct.Balance()) // $100.00

	// Go automatically takes the address when calling pointer methods
	// on an addressable value — so this also works:
	acct2 := Account{Owner: "Bob"}
	acct2.Deposit(50) // Go translates this to (&acct2).Deposit(50)
	fmt.Printf("%s's balance: $%.2f\n", acct2.Owner, acct2.Balance())
}

// -------------------------------------------------------------------------
// Constructor Functions (The NewXxx Pattern)
// -------------------------------------------------------------------------

/*
 Go doesn't have constructors. Instead, the convention is to write a
 function called NewXxx that returns an initialized value. This is
 especially useful when:

 - Fields need validation
 - Unexported fields need to be set
 - The zero value isn't useful for this type
 - You want to return a pointer (which is the common case)
*/

type Config struct {
	Host         string
	Port         int
	ReadTimeout  int // seconds
	WriteTimeout int // seconds
	MaxRetries   int
}

// NewConfig creates a Config with sensible defaults.
// It returns a pointer — the caller can modify it further.
func NewConfig(host string, port int) *Config {
	return &Config{
		Host:         host,
		Port:         port,
		ReadTimeout:  30,
		WriteTimeout: 30,
		MaxRetries:   3,
	}
}

func DemoConstructor() {
	cfg := NewConfig("api.example.com", 443)
	fmt.Printf("Config: %s:%d (timeout: %ds)\n", cfg.Host, cfg.Port, cfg.ReadTimeout)

	// Override defaults as needed
	cfg.ReadTimeout = 60
	cfg.MaxRetries = 5
}

// -------------------------------------------------------------------------
// Struct Tags (Preview)
// -------------------------------------------------------------------------

/*
 Struct tags are string annotations attached to struct fields. They're
 metadata that other packages can read at runtime using reflection. You'll
 use them constantly for:

 - JSON serialization:  `json:"field_name"`
 - Database mapping:    `db:"column_name"`
 - Validation:          `validate:"required,email"`
 - Form binding:        `form:"field_name"`

 We're just previewing them here — they become essential when we cover
 JSON handling and web services.
*/

// APIResponse shows common struct tags you'll encounter in web services.
type APIResponse struct {
	ID        int    `json:"id"`
	UserName  string `json:"user_name"`
	Email     string `json:"email,omitempty"` // omit from JSON if empty
	Password  string `json:"-"`               // never include in JSON output
	CreatedAt string `json:"created_at"`
}

// -------------------------------------------------------------------------
// Comparable Structs
// -------------------------------------------------------------------------

/*
 A struct is comparable (can use == and !=) if ALL of its fields are
 comparable types. This means no slices, maps, or function fields.

 Comparable structs can be used as map keys, which is surprisingly useful.
 For example, a (row, col) coordinate can be a map key for a grid.
*/

type Coordinate struct {
	Row, Col int
}

func DemoComparableStructs() {
	c1 := Coordinate{1, 2}
	c2 := Coordinate{1, 2}
	c3 := Coordinate{3, 4}

	fmt.Println(c1 == c2) // true — all fields are equal
	fmt.Println(c1 == c3) // false

	// Use as map keys!
	grid := map[Coordinate]string{
		{0, 0}: "origin",
		{1, 0}: "right",
		{0, 1}: "up",
	}
	fmt.Println(grid[Coordinate{0, 0}]) // "origin"
}

// -------------------------------------------------------------------------
// Method Sets and Interface Preview
// -------------------------------------------------------------------------

/*
 A type's "method set" determines which interfaces it satisfies. This is
 a subtle but important concept:

 - Type T's method set includes methods with value receivers only.
 - Type *T's method set includes methods with BOTH value and pointer receivers.

 This means if you define a method with a pointer receiver, only *T (not T)
 satisfies interfaces requiring that method. We'll explore this fully in
 the interfaces module, but keep it in mind as you design your methods.

 Practical implication: if your type will need to satisfy an interface,
 be deliberate about your receiver types.
*/

// Stringer is Go's equivalent of toString(). Implementing fmt.Stringer
// lets your type control how it appears in fmt.Println, Printf %v, etc.
type Color struct {
	R, G, B uint8
}

// String implements fmt.Stringer. Now Color prints nicely everywhere.
func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func DemoStringer() {
	red := Color{255, 0, 0}
	fmt.Println(red) // "#ff0000" — our String() method is called
}
