package pointers

/*
=============================================================================
 Module 07: Pointers
=============================================================================

 Pointers are one of Go's most important concepts, and they're simpler
 than you might fear. A pointer is just a variable that holds the memory
 address of another variable. That's it.

 If you've used C/C++, Go's pointers will feel familiar but safer — there's
 no pointer arithmetic, no manual memory management, and the garbage
 collector handles cleanup. If you're coming from Python or JavaScript,
 think of pointers as making "pass by reference" explicit.

 WHY POINTERS MATTER:
 Go is a pass-by-value language. When you pass a struct to a function, the
 function gets a COPY. If the function modifies it, the original is
 unchanged. Pointers let you share data without copying, and let functions
 modify the caller's data.

 WHY POINTERS MATTER FOR WEB SERVICES:
 - HTTP request/response bodies are typically accessed through pointers
 - JSON optional fields use pointer types (*string, *int) to distinguish
   "field is absent" from "field is zero"
 - Database connections, loggers, and services are passed as pointers
 - Middleware chains pass shared context through pointer-based structures

=============================================================================
*/

import "fmt"

// -------------------------------------------------------------------------
// What Pointers Are: & (address-of) and * (dereference)
// -------------------------------------------------------------------------

/*
 Two operators are all you need:

 &  (address-of):  Get the memory address of a variable
                    &x gives you a pointer to x

 *  (dereference):  Follow a pointer to the value it points to
                    *p gives you the value that p points to

 The * symbol does double duty:
 - In a TYPE declaration, *int means "pointer to int"
 - In an EXPRESSION, *p means "get the value p points to"
*/

func DemoPointerBasics() {
	x := 42
	p := &x // p is a pointer to x (type: *int)

	fmt.Println(x)  // 42 — the value
	fmt.Println(p)  // 0xc0000b4008 — the memory address (varies)
	fmt.Println(*p) // 42 — dereferencing the pointer gives the value

	// Modifying through the pointer changes the original
	*p = 100
	fmt.Println(x) // 100 — x was changed through the pointer!

	// You can also create pointers with new()
	q := new(int)   // allocates an int, returns *int, value is 0
	fmt.Println(*q) // 0 — zero value

	// But using & is more common and more readable
	name := "Alice"
	namePtr := &name
	fmt.Println(*namePtr) // "Alice"
}

// -------------------------------------------------------------------------
// Pass by Value in Go
// -------------------------------------------------------------------------

/*
 THIS IS CRITICAL TO UNDERSTAND:

 Everything in Go is passed by value. When you call a function, the
 argument is COPIED. This applies to:

 - Basic types (int, string, bool)
 - Structs (the entire struct is copied)
 - Arrays (the entire array is copied — this is unusual!)
 - Pointers (the pointer itself is copied, but it still points to the
   same data — this is how you share data)

 Slices and maps seem like exceptions, but they're not: the slice HEADER
 (pointer + length + capacity) is copied, and the map HEADER (pointer to
 internal structure) is copied. The underlying data isn't copied, which
 is why modifications to slice elements are visible to the caller.
*/

type Coordinate struct {
	X, Y float64
}

// This CANNOT modify the original — it receives a copy
func translateValue(c Coordinate, dx, dy float64) {
	c.X += dx // modifying the copy — original is unchanged
	c.Y += dy
}

// This CAN modify the original — it receives a pointer
func translatePointer(c *Coordinate, dx, dy float64) {
	c.X += dx // modifying through the pointer — original IS changed
	c.Y += dy
}

func DemoPassByValue() {
	c := Coordinate{X: 1, Y: 2}

	translateValue(c, 10, 20)
	fmt.Println(c) // {1, 2} — unchanged! The function got a copy.

	translatePointer(&c, 10, 20)
	fmt.Println(c) // {11, 22} — changed! The function got a pointer.
}

// -------------------------------------------------------------------------
// When to Use Pointers vs Values
// -------------------------------------------------------------------------

/*
 USE POINTERS WHEN:
 1. The function needs to modify the caller's data
 2. The struct is large (avoid expensive copies)
 3. You need to represent "no value" / "absent" (nil pointer)
 4. The type has mutable internal state (counters, caches, connections)
 5. You need consistency — if some methods use pointer receivers, all should

 USE VALUES WHEN:
 1. The type is small and immutable (coordinates, colors, time.Time)
 2. You want safety — the caller can't accidentally modify your data
 3. The type is a basic type (int, string, bool) — always pass by value
 4. You're creating a functional-style API (returns new values)

 COMMON PATTERNS IN WEB SERVICES:
 - Services and handlers: always pointers (they have state/dependencies)
 - Request/response bodies: typically pointers (large, often modified)
 - Configuration: pointer if mutable, value if immutable
 - Value objects (money, coordinates): values (small, immutable)
*/

// -------------------------------------------------------------------------
// Nil Pointers and Safe Handling
// -------------------------------------------------------------------------

/*
 A pointer that doesn't point to anything has the value nil. Dereferencing
 a nil pointer causes a runtime panic — this is Go's equivalent of a null
 pointer exception.

 Always check for nil when:
 - A function might return a nil pointer (like a lookup that finds nothing)
 - You receive a pointer from outside your package
 - You're working with optional fields

 You can even call methods on nil receivers — this is a useful pattern
 for providing safe defaults.
*/

type SafeConfig struct {
	Host    string
	Port    int
	Verbose bool
}

// GetHost handles nil receiver gracefully — returns a default.
// This is a legitimate Go pattern. The method checks if the receiver
// is nil and returns a sensible default instead of panicking.
func (c *SafeConfig) GetHost() string {
	if c == nil {
		return "localhost" // safe default
	}
	return c.Host
}

func (c *SafeConfig) GetPort() int {
	if c == nil {
		return 8080 // safe default
	}
	return c.Port
}

func DemoNilPointers() {
	// Normal usage
	cfg := &SafeConfig{Host: "example.com", Port: 443}
	fmt.Println(cfg.GetHost()) // "example.com"

	// Nil pointer — doesn't panic because the method handles it
	var nilCfg *SafeConfig
	fmt.Println(nilCfg.GetHost()) // "localhost" — safe default
	fmt.Println(nilCfg.GetPort()) // 8080 — safe default
}

// -------------------------------------------------------------------------
// Pointers and Slices
// -------------------------------------------------------------------------

/*
 Slices are already "reference-like" because a slice header contains a
 pointer to the underlying array. When you pass a slice to a function:

 - The HEADER is copied (pointer, length, capacity)
 - But the pointer still points to the same array
 - So modifications to ELEMENTS are visible to the caller
 - But append() might create a new array, invisible to the caller

 This is why you often see functions return a slice even when they
 modify it — the append might have changed the underlying array.
*/

func DemoSlicesAndPointers() {
	s := []int{1, 2, 3}

	// Modifying elements is visible (shared underlying array)
	modifyElements(s)
	fmt.Println(s) // [100 2 3] — element 0 was changed

	// But append might not be visible (could create new array)
	appendInvisible(s)
	fmt.Println(s) // [100 2 3] — append happened on a copy of the header
}

func modifyElements(s []int) {
	s[0] = 100 // modifies the shared underlying array
}

func appendInvisible(s []int) {
	s = append(s, 4) // this might create a new array — caller doesn't see it
	_ = s
}

// -------------------------------------------------------------------------
// new() vs & — Both Create Pointers
// -------------------------------------------------------------------------

/*
 There are two ways to create a pointer to a newly allocated value:

   p := new(int)          // allocates, returns *int pointing to 0
   p := &MyStruct{...}    // allocates, returns *MyStruct with fields set

 The & form is almost always preferred because:
 1. You can set field values at creation time
 2. It's more readable
 3. It works the same way

 new() is mainly useful for basic types when you need a pointer:
   p := new(int)    // can't do &int{} — that's not valid syntax
*/

func DemoNewVsAddress() {
	// new() — useful for basic types
	p := new(int) // *int pointing to 0
	*p = 42
	fmt.Println(*p) // 42

	// & — preferred for structs
	cfg := &SafeConfig{Host: "example.com", Port: 443}
	fmt.Println(cfg.Host) // "example.com"

	// For basic types, you can use a helper if you need a pointer literal:
	s := stringPtr("hello") // common pattern in codebases
	fmt.Println(*s)         // "hello"
}

// stringPtr is a common helper — returns a pointer to a string value.
// You'll see these all over Go codebases, especially for optional JSON fields.
func stringPtr(s string) *string {
	return &s
}

// -------------------------------------------------------------------------
// Stack vs Heap (Escape Analysis Preview)
// -------------------------------------------------------------------------

/*
 In many languages, you have to think about stack vs heap allocation.
 In Go, the compiler decides for you through "escape analysis":

 - If a variable doesn't leave the function, it goes on the stack (fast)
 - If a variable "escapes" (returned via pointer, captured in closure),
   it goes on the heap (slightly slower, GC'd)

 You don't need to worry about this for correctness — it's an optimization.
 But it's good to know:
 - Returning a pointer to a local variable is SAFE in Go (it escapes to heap)
 - The garbage collector handles heap cleanup
 - You can check with: go build -gcflags="-m" (shows escape analysis)

 Compare with C: returning a pointer to a local variable is UNDEFINED
 BEHAVIOR. In Go, it's perfectly fine and common.
*/

func createOnHeap() *Coordinate {
	c := Coordinate{X: 1, Y: 2} // This "escapes" because we return a pointer
	return &c                    // Allocated on heap — safe in Go, dangerous in C
}

func DemoEscapeAnalysis() {
	c := createOnHeap()
	fmt.Println(c) // &{1 2} — works perfectly, c lives on the heap
}

// -------------------------------------------------------------------------
// Common Pointer Patterns in Web Services
// -------------------------------------------------------------------------

/*
 OPTIONAL JSON FIELDS:
 In JSON APIs, there's a difference between "field is absent" and "field
 is its zero value." For example:
   {"name": ""}     — name is present but empty
   {}               — name is absent
   {"name": null}   — name is explicitly null

 Using *string lets you distinguish these cases:
   nil     → field was absent / null
   &""     → field was present but empty
   &"Bob"  → field was present with a value

 This pattern is everywhere in API development.
*/

// UserUpdate represents a JSON PATCH request. Pointer fields allow us
// to distinguish "not provided" (nil) from "set to empty" (&"").
type UserUpdate struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Age   *int    `json:"age,omitempty"`
}

func DemoOptionalFields() {
	// Name provided, email not provided, age explicitly set to 0
	name := "Alice"
	age := 0
	update := UserUpdate{
		Name:  &name, // provided
		Email: nil,   // not provided — don't update this field
		Age:   &age,  // explicitly set to 0
	}

	if update.Name != nil {
		fmt.Printf("Update name to: %s\n", *update.Name)
	}
	if update.Email != nil {
		fmt.Printf("Update email to: %s\n", *update.Email)
	} else {
		fmt.Println("Email not provided — don't change it")
	}
	if update.Age != nil {
		fmt.Printf("Update age to: %d\n", *update.Age)
	}
}
