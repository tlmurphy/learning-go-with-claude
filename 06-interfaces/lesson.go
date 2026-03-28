package interfaces

/*
=============================================================================
 Module 06: Interfaces
=============================================================================

 Interfaces are the single most important concept in Go for writing
 flexible, testable, maintainable code. They define behavior — what a
 type can DO — rather than what a type IS.

 The key insight: Go interfaces are satisfied IMPLICITLY. There is no
 "implements" keyword. If your type has the right methods, it satisfies
 the interface. Period. The type doesn't even need to know the interface
 exists.

 This is sometimes called "structural typing" or "duck typing" — if it
 walks like a duck and quacks like a duck, it's a duck.

 WHY THIS MATTERS FOR WEB SERVICES:
 Interfaces are everywhere in Go web development:
  - http.Handler interface powers all HTTP routing
  - io.Reader/io.Writer are how you process request/response bodies
  - database/sql uses interfaces for driver abstraction
  - You'll use interfaces for dependency injection, middleware, and testing

=============================================================================
*/

import (
	"fmt"
	"math"
	"strings"
)

// -------------------------------------------------------------------------
// Implicit Interface Satisfaction
// -------------------------------------------------------------------------

/*
 In Java or C#, you'd write: class Dog implements Animal { ... }
 In Go, you just... implement the methods. No declaration needed.

 This has profound implications:
 - You can define interfaces AFTER the concrete types exist
 - You can define interfaces in a DIFFERENT package from the types
 - Third-party types can satisfy your interfaces without modification

 This is why Go's standard library interfaces are so powerful — your
 custom types can satisfy io.Reader without importing or knowing about io.
*/

// Speaker is a simple interface requiring one method.
type Speaker interface {
	Speak() string
}

// Dog satisfies Speaker by having a Speak() method. Note: no "implements".
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + " says: Woof!"
}

// Cat also satisfies Speaker — completely independently of Dog.
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return c.Name + " says: Meow!"
}

// MakeThemSpeak accepts any Speaker — it doesn't care about the concrete type.
func MakeThemSpeak(speakers []Speaker) []string {
	results := make([]string, len(speakers))
	for i, s := range speakers {
		results[i] = s.Speak()
	}
	return results
}

func DemoImplicitInterfaces() {
	// Both Dog and Cat satisfy Speaker without any explicit declaration
	animals := []Speaker{
		Dog{Name: "Rex"},
		Cat{Name: "Whiskers"},
	}
	for _, msg := range MakeThemSpeak(animals) {
		fmt.Println(msg)
	}
}

// -------------------------------------------------------------------------
// The Power of Small Interfaces
// -------------------------------------------------------------------------

/*
 Go's most powerful interfaces are tiny:

   type Reader interface { Read(p []byte) (n int, err error) }
   type Writer interface { Write(p []byte) (n int, err error) }
   type Stringer interface { String() string }
   type Error interface { Error() string }

 One or two methods. That's it. Small interfaces are:
 - Easy to implement (low barrier for new types)
 - Easy to compose (combine small interfaces into larger ones)
 - Easy to mock in tests (less to fake)

 The Go proverb: "The bigger the interface, the weaker the abstraction."

 In web services, you'll see this constantly:
 - A function that takes io.Reader can read from files, HTTP bodies,
   strings, compressed streams, encrypted data — anything.
 - A function that takes io.Writer can write to files, HTTP responses,
   buffers, network connections — anything.
*/

// Describer is a small interface — just one method.
type Describer interface {
	Describe() string
}

// Circle satisfies Describer
type Circle struct {
	Radius float64
}

func (c Circle) Describe() string {
	return fmt.Sprintf("Circle with radius %.2f (area: %.2f)", c.Radius, math.Pi*c.Radius*c.Radius)
}

// Square satisfies Describer
type Square struct {
	Side float64
}

func (s Square) Describe() string {
	return fmt.Sprintf("Square with side %.2f (area: %.2f)", s.Side, s.Side*s.Side)
}

// -------------------------------------------------------------------------
// Interface Composition (Embedding Interfaces)
// -------------------------------------------------------------------------

/*
 Just as you can embed structs in structs, you can embed interfaces in
 interfaces. This lets you build complex interfaces from simple ones.

 The standard library does this extensively:

   type ReadWriter interface {
       Reader
       Writer
   }

   type ReadWriteCloser interface {
       Reader
       Writer
       Closer
   }

 Start small, compose as needed. Don't define a big interface upfront.
*/

// Sizer has one method
type Sizer interface {
	Size() float64
}

// Namer has one method
type Namer interface {
	Name() string
}

// NamedSizer composes both — requires both Name() and Size().
type NamedSizer interface {
	Namer
	Sizer
}

// File satisfies NamedSizer because it has both Name() and Size()
type File struct {
	name string
	size float64
}

func (f File) Name() string  { return f.name }
func (f File) Size() float64 { return f.size }

func DemoComposition() {
	f := File{name: "data.json", size: 1024}

	// f satisfies Namer, Sizer, AND NamedSizer — all implicitly
	var n Namer = f
	var s Sizer = f
	var ns NamedSizer = f

	fmt.Println(n.Name(), s.Size(), ns.Name(), ns.Size())
}

// -------------------------------------------------------------------------
// The Empty Interface: any / interface{}
// -------------------------------------------------------------------------

/*
 The empty interface (interface{}) has zero methods, so EVERY type
 satisfies it. Go 1.18 added "any" as an alias — they're identical:

   var x interface{} = 42  // old style
   var y any = "hello"     // new style (preferred)

 When to use it:
 - fmt.Println takes ...any (it has to accept anything)
 - JSON unmarshaling into unknown structures
 - Generic containers (before Go had generics)

 When to AVOID it:
 - Almost everywhere else! Using any throws away type safety.
 - If you know the type, use the type. If you know the behavior, use
   an interface. Only use any as a last resort.

 The introduction of generics in Go 1.18 eliminated many use cases for
 empty interfaces. Prefer generics when you need type flexibility with
 type safety.
*/

func DemoEmptyInterface() {
	// any can hold any value
	var things []any
	things = append(things, 42, "hello", true, 3.14, []int{1, 2, 3})

	for _, thing := range things {
		fmt.Printf("Type: %T, Value: %v\n", thing, thing)
	}
}

// -------------------------------------------------------------------------
// Type Assertions and Type Switches
// -------------------------------------------------------------------------

/*
 When you have an interface value and need to access the concrete type
 underneath, you use type assertions or type switches.

 Type assertion:  value, ok := iface.(ConcreteType)
 Type switch:     switch v := iface.(type) { case ConcreteType: ... }

 Always use the two-value form of type assertions (value, ok). The single-
 value form panics if the assertion fails — that's almost never what you want.
*/

func DescribeValue(i any) string {
	// Type switch — the idiomatic way to handle multiple types
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("integer: %d", v)
	case string:
		return fmt.Sprintf("string: %q (length %d)", v, len(v))
	case bool:
		if v {
			return "boolean: true"
		}
		return "boolean: false"
	case []int:
		return fmt.Sprintf("int slice with %d elements", len(v))
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

func DemoTypeAssertions() {
	var s Speaker = Dog{Name: "Rex"}

	// Two-value assertion — safe, returns ok=false on failure
	dog, ok := s.(Dog)
	if ok {
		fmt.Println("It's a dog named", dog.Name)
	}

	// This would be false (it's a Dog, not a Cat)
	_, ok = s.(Cat)
	fmt.Println("Is it a cat?", ok) // false

	// Type switch for multiple types
	fmt.Println(DescribeValue(42))
	fmt.Println(DescribeValue("hello"))
	fmt.Println(DescribeValue(true))
}

// -------------------------------------------------------------------------
// Interface Values: The (type, value) Pair — Nil Gotcha
// -------------------------------------------------------------------------

/*
 THIS IS ONE OF GO'S MOST CONFUSING ASPECTS. Read carefully.

 An interface value is internally a pair: (concrete type, concrete value).
 An interface is nil ONLY when BOTH the type and value are nil.

 This leads to a notorious gotcha:

   var p *Dog = nil          // a nil pointer to Dog
   var s Speaker = p         // s holds (type=*Dog, value=nil)
   fmt.Println(s == nil)     // FALSE! s has a type, even though value is nil

 This trips up even experienced Go developers. The fix: don't assign
 typed nil pointers to interfaces. Return the interface type directly:

   func GetSpeaker() Speaker {
       return nil  // this is truly nil — no type, no value
   }
*/

func DemoNilInterface() {
	// Truly nil interface — both type and value are nil
	var s Speaker
	fmt.Println(s == nil) // true

	// NON-nil interface holding a nil pointer — the gotcha
	var d *Dog = nil
	s = d
	fmt.Println(s == nil) // FALSE! type is *Dog, value is nil
	fmt.Println(d == nil) // true — the pointer itself is nil

	// This is why you should check the concrete type if you're unsure
}

// -------------------------------------------------------------------------
// Accept Interfaces, Return Structs
// -------------------------------------------------------------------------

/*
 This is one of Go's most important design principles:

 ACCEPT INTERFACES: Make your function parameters interfaces so callers
 can pass any type that satisfies the contract. This makes your code
 flexible and testable.

 RETURN STRUCTS: Return concrete types so callers get the full API of
 the type. Don't hide concrete capabilities behind an interface.

 Example from the standard library:
   func Copy(dst Writer, src Reader) (int64, error)
   // Accepts interfaces — works with any reader/writer

   func NewBuffer(buf []byte) *Buffer
   // Returns concrete *Buffer — caller gets full Buffer API

 In web services, this principle drives dependency injection:
   func NewUserService(db Database) *UserService
   // db is an interface — easy to pass a mock for testing
   // returns concrete *UserService — caller gets all methods
*/

// Formatter is an interface — accept this in function parameters
type Formatter interface {
	Format(s string) string
}

// UpperFormatter is a concrete type — return this from constructors
type UpperFormatter struct{}

func (uf UpperFormatter) Format(s string) string {
	return strings.ToUpper(s)
}

// FormatAll accepts the interface — any Formatter works
func FormatAll(f Formatter, items []string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = f.Format(item)
	}
	return result
}

// NewUpperFormatter returns the concrete type, not the interface
func NewUpperFormatter() *UpperFormatter {
	return &UpperFormatter{}
}

// -------------------------------------------------------------------------
// Common Standard Library Interfaces
// -------------------------------------------------------------------------

/*
 You'll encounter these interfaces constantly in Go:

 fmt.Stringer:
   type Stringer interface { String() string }
   Controls how your type appears in fmt.Println, Printf %v, etc.

 error:
   type error interface { Error() string }
   The foundation of all error handling in Go.

 io.Reader:
   type Reader interface { Read(p []byte) (n int, err error) }
   The universal "I can be read from" interface. Files, HTTP bodies,
   network connections, strings — all are Readers.

 io.Writer:
   type Writer interface { Write(p []byte) (n int, err error) }
   The universal "I can be written to" interface.

 io.Closer:
   type Closer interface { Close() error }
   Anything that needs cleanup (files, connections, etc.)

 http.Handler:
   type Handler interface { ServeHTTP(ResponseWriter, *Request) }
   The interface that powers all Go HTTP handling.

 sort.Interface:
   type Interface interface { Len() int; Less(i, j int) bool; Swap(i, j int) }
   Makes any collection sortable. (Generics have made this less common.)
*/

// -------------------------------------------------------------------------
// Interface Pollution Anti-Pattern
// -------------------------------------------------------------------------

/*
 Don't create interfaces prematurely. A common mistake is to define an
 interface for every type "just in case." This is interface pollution.

 BAD (premature abstraction):
   type UserRepository interface { ... }
   type userRepository struct { ... }
   func NewUserRepository() UserRepository { return &userRepository{} }

 GOOD (interface defined by the consumer, when needed):
   // In the user service package:
   type UserStore interface { FindByID(id int) (*User, error) }
   // In the repository package:
   type UserRepository struct { ... }
   func (r *UserRepository) FindByID(id int) (*User, error) { ... }

 Define interfaces where they're USED, not where they're implemented.
 This is the opposite of Java's approach, and it works beautifully because
 Go interfaces are implicit.
*/
