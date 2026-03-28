// Package variables introduces Go's type system, variable declarations,
// constants, and the foundational concept of zero values.
package variables

import "fmt"

/*
=============================================================================
 VARIABLES AND TYPES IN GO
=============================================================================

Go is a statically typed language, but it doesn't feel as heavy as Java or C++.
The type system is designed to be simple, explicit, and to catch bugs at compile
time rather than at 3 AM in production.

The key philosophy: Go never does implicit type conversions. If you want to mix
an int and a float64 in an expression, you must explicitly convert one. This
might feel annoying at first, but it prevents an entire class of subtle bugs
that plague languages with implicit conversions (looking at you, JavaScript).

=============================================================================
 VARIABLE DECLARATIONS
=============================================================================

Go gives you several ways to declare variables. Each has its place:

  var x int          // Explicit type, gets zero value (0)
  var x int = 42     // Explicit type with initial value
  var x = 42         // Type inferred from value (int)
  x := 42            // Short declaration — most common inside functions

The := operator is the workhorse of Go code. You'll use it constantly inside
functions. But it ONLY works inside functions — package-level variables must
use var.

Why does this matter for web services? When you're handling an HTTP request,
you'll declare variables constantly:

  name := r.URL.Query().Get("name")    // short declaration from query param
  age, err := strconv.Atoi(ageStr)     // multi-value short declaration

=============================================================================
*/

// DemoVarDeclarations shows the different ways to declare variables in Go.
func DemoVarDeclarations() {
	// Long form — useful when you want to be explicit about the type,
	// or when you want the zero value.
	var count int
	var name string
	var active bool

	fmt.Println("Long form (zero values):")
	fmt.Printf("  count=%d, name=%q, active=%t\n", count, name, active)

	// Long form with initialization
	var port int = 8080
	var host string = "localhost"
	fmt.Printf("  port=%d, host=%s\n", port, host)

	// Type inference with var — the compiler figures out the type from the value.
	// Useful at package level where := isn't allowed.
	var timeout = 30       // int
	var ratio = 3.14       // float64
	var greeting = "hello" // string
	fmt.Printf("  timeout=%d, ratio=%f, greeting=%s\n", timeout, ratio, greeting)

	// Short declaration with := (the Go workhorse)
	// This is what you'll use 90% of the time inside functions.
	age := 25
	price := 19.99
	message := "Go is great"
	fmt.Printf("  age=%d, price=%.2f, message=%s\n", age, price, message)

	// Multiple assignment — declare several variables at once.
	// You'll see this pattern with functions that return multiple values.
	x, y, z := 1, 2, 3
	fmt.Printf("  x=%d, y=%d, z=%d\n", x, y, z)

	// Block declaration — group related variables together.
	// Common at package level for configuration.
	var (
		maxRetries = 3
		baseURL    = "https://api.example.com"
		verbose    = false
	)
	fmt.Printf("  maxRetries=%d, baseURL=%s, verbose=%t\n", maxRetries, baseURL, verbose)
}

/*
=============================================================================
 BASIC TYPES
=============================================================================

Go's type system is straightforward. Here are the types you'll use most:

  Integers:  int, int8, int16, int32, int64
             uint, uint8, uint16, uint32, uint64
  Floats:    float32, float64
  Complex:   complex64, complex128 (yes, Go has built-in complex numbers!)
  String:    string (immutable sequence of bytes, usually UTF-8)
  Boolean:   bool
  Byte:      byte (alias for uint8)
  Rune:      rune (alias for int32, represents a Unicode code point)

When should you use int vs int64? Use plain `int` for most things — it's
platform-dependent (64-bit on modern systems) and it's what most standard
library functions expect. Use sized integers when you need a specific size
(binary protocols, memory-critical code, or matching an external API).

=============================================================================
*/

// DemoBasicTypes shows Go's built-in types and how they work.
func DemoBasicTypes() {
	// Integers — int is the default, platform-dependent (usually 64-bit).
	var i int = 42
	var i8 int8 = 127   // -128 to 127
	var i64 int64 = 9999 // when you need a specific size

	fmt.Println("Integer types:")
	fmt.Printf("  int=%d, int8=%d, int64=%d\n", i, i8, i64)

	// Floats — float64 is the default for decimal numbers.
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793

	fmt.Println("Float types:")
	fmt.Printf("  float32=%f, float64=%.15f\n", f32, f64)

	// Strings are immutable and UTF-8 encoded by default.
	// This is important: len(s) gives you BYTES, not characters!
	s := "Hello, 世界" // contains ASCII and Chinese characters
	fmt.Println("String internals:")
	fmt.Printf("  string=%q, len (bytes)=%d\n", s, len(s))
	fmt.Printf("  rune count=%d\n", len([]rune(s)))

	// byte is an alias for uint8, rune is an alias for int32.
	// This distinction matters when processing text.
	var b byte = 'A'  // single character in single quotes
	var r rune = '世'  // Unicode code point — needs int32 to represent
	fmt.Printf("  byte='%c' (%d), rune='%c' (%d)\n", b, b, r, r)

	// Booleans — no truthy/falsy in Go. Conditions must be explicitly bool.
	// if count { ... }  <-- COMPILE ERROR in Go. Must be: if count > 0 { ... }
	var isReady bool = true
	fmt.Printf("  bool=%t\n", isReady)
}

/*
=============================================================================
 ZERO VALUES
=============================================================================

This is one of Go's most important concepts. Every variable in Go is
initialized to a well-defined "zero value" if you don't give it one.
There is no "undefined" or uninitialized memory in Go.

  int, float64     -> 0, 0.0
  string           -> "" (empty string)
  bool             -> false
  pointer          -> nil
  slice, map, chan -> nil (but nil slices are still usable!)
  struct           -> all fields get their zero values

Why does this matter? It means you can declare variables and start using
them immediately without worrying about garbage values. Many Go APIs are
designed so that the zero value is useful — for example, a zero-value
sync.Mutex is an unlocked mutex, ready to use.

This design choice comes up constantly in web services: a zero-value
http.Server works out of the box, a zero-value bytes.Buffer is an empty
buffer ready for writing, etc.

=============================================================================
*/

// DemoZeroValues shows that every type in Go has a well-defined zero value.
func DemoZeroValues() {
	var i int
	var f float64
	var s string
	var b bool
	var p *int        // pointer
	var sl []int      // slice
	var m map[string]int // map

	fmt.Println("Zero values:")
	fmt.Printf("  int:     %d\n", i)
	fmt.Printf("  float64: %f\n", f)
	fmt.Printf("  string:  %q\n", s) // %q shows quotes, making empty string visible
	fmt.Printf("  bool:    %t\n", b)
	fmt.Printf("  pointer: %v\n", p)
	fmt.Printf("  slice:   %v (nil? %t, len=%d)\n", sl, sl == nil, len(sl))
	fmt.Printf("  map:     %v (nil? %t)\n", m, m == nil)

	// A nil slice is perfectly fine to append to!
	// This is a deliberate design choice — you don't need to initialize slices.
	sl = append(sl, 1, 2, 3)
	fmt.Printf("  slice after append: %v\n", sl)

	// But a nil map will PANIC if you try to write to it.
	// You must initialize maps before writing. This is a common gotcha!
	// m["key"] = 1  // PANIC: assignment to entry in nil map
	m = make(map[string]int) // initialize the map
	m["key"] = 1             // now this works
	fmt.Printf("  map after make: %v\n", m)
}

/*
=============================================================================
 CONSTANTS AND IOTA
=============================================================================

Constants in Go are computed at compile time. They can be "typed" or
"untyped" — and untyped constants are surprisingly flexible.

  const Pi = 3.14159          // untyped — adapts to context
  const MaxPort int = 65535   // typed — always int

Untyped constants are one of Go's clever features. The constant Pi above
has no fixed type — it adapts to whatever context it's used in. You can
use it as a float32, float64, or even in a complex number expression
without explicit conversion.

iota is Go's enum generator. It starts at 0 and increments for each
constant in a const block. It's simple but powerful — you can use
expressions with iota to create bit flags, powers of two, etc.

=============================================================================
*/

// LogLevel is a custom type for log levels, demonstrating iota enums.
type LogLevel int

const (
	// iota starts at 0 and increments by 1 for each constant.
	// This is how you create "enums" in Go.
	LogDebug LogLevel = iota // 0
	LogInfo                  // 1 — iota increments automatically
	LogWarn                  // 2
	LogError                 // 3
	LogFatal                 // 4
)

// ByteSize demonstrates iota with expressions for powers of 1024.
type ByteSize float64

const (
	_           = iota             // blank identifier discards 0
	KB ByteSize = 1 << (10 * iota) // 1 << 10 = 1024
	MB                              // 1 << 20
	GB                              // 1 << 30
	TB                              // 1 << 40
)

// DemoConstants shows how constants and iota work in Go.
func DemoConstants() {
	// Untyped constants adapt to their context.
	const Pi = 3.14159265358979323846
	const MaxConnections = 100

	// Pi can be used as float32, float64, etc. without conversion.
	var f32 float32 = Pi
	var f64 float64 = Pi
	fmt.Println("Untyped constant Pi:")
	fmt.Printf("  as float32: %f\n", f32)
	fmt.Printf("  as float64: %.15f\n", f64)

	// Typed constants are locked to their type.
	const port int = 8080
	// var p float64 = port  // COMPILE ERROR: cannot use port (type int) as float64

	fmt.Println("iota enum (LogLevel):")
	fmt.Printf("  Debug=%d, Info=%d, Warn=%d, Error=%d, Fatal=%d\n",
		LogDebug, LogInfo, LogWarn, LogError, LogFatal)

	fmt.Println("iota with expressions (ByteSize):")
	fmt.Printf("  KB=%.0f, MB=%.0f, GB=%.0f, TB=%.0f\n", KB, MB, GB, TB)

	// Constants must be determinable at compile time.
	// const now = time.Now()  // COMPILE ERROR: time.Now() is not a constant
	_ = MaxConnections // use it so the compiler doesn't complain
}

/*
=============================================================================
 TYPE CONVERSIONS
=============================================================================

Go has NO implicit type conversions. Period. This is a deliberate design
choice that prevents subtle bugs. If you want to combine an int and a
float64, you must explicitly convert one:

  var i int = 42
  var f float64 = float64(i)   // explicit conversion required

This also applies to seemingly compatible types:
  var i32 int32 = 42
  var i64 int64 = int64(i32)   // even int32 -> int64 needs explicit conversion!

String conversions have special behavior:
  string(65)         -> "A"  (converts int to Unicode code point!)
  string([]byte{...}) -> converts bytes to string
  []byte("hello")    -> converts string to byte slice

The strconv package is your friend for string <-> number conversions:
  strconv.Itoa(42)       -> "42"
  strconv.Atoi("42")     -> 42, nil  (note: returns an error!)

=============================================================================
*/

// DemoTypeConversions shows Go's explicit type conversion rules.
func DemoTypeConversions() {
	// Numeric conversions — always explicit.
	var i int = 42
	var f float64 = float64(i) // int -> float64
	var u uint = uint(i)       // int -> uint

	fmt.Println("Numeric conversions:")
	fmt.Printf("  int=%d -> float64=%f, uint=%d\n", i, f, u)

	// Be careful: narrowing conversions can lose data silently!
	var big int64 = 256
	var small int8 = int8(big) // 256 overflows int8 -> wraps to 0
	fmt.Printf("  int64(256) -> int8 = %d (overflow!)\n", small)

	// Float to int truncates (does NOT round).
	var pi float64 = 3.99
	var truncated int = int(pi) // 3, not 4
	fmt.Printf("  float64(3.99) -> int = %d (truncated!)\n", truncated)

	// String and byte/rune conversions.
	// string(int) converts to Unicode code point, NOT to decimal representation!
	fmt.Println("String conversions:")
	fmt.Printf("  string(65) = %q (Unicode code point, not \"65\"!)\n", string(rune(65)))

	// String <-> []byte for when you need mutable string data.
	s := "Hello"
	bytes := []byte(s) // string -> byte slice (copy!)
	bytes[0] = 'h'     // modify the copy
	s2 := string(bytes) // byte slice -> string (another copy!)
	fmt.Printf("  original=%q, modified=%q\n", s, s2)

	// String <-> []rune for Unicode-aware character manipulation.
	emoji := "Go is 🎉 fun"
	runes := []rune(emoji)
	fmt.Printf("  %q has %d bytes but %d runes\n", emoji, len(emoji), len(runes))
}

/*
=============================================================================
 TYPE DEFINITIONS AND ALIASES
=============================================================================

Go lets you create new types based on existing ones. There are two forms:

  type Celsius float64        // Type DEFINITION: creates a new, distinct type
  type Float = float64        // Type ALIAS: just another name for the same type

Type definitions are powerful because the new type is distinct — you can
add methods to it, and it won't accidentally mix with the underlying type:

  type UserID int64
  type OrderID int64
  // These are different types! You can't accidentally pass a UserID
  // where an OrderID is expected. The compiler catches it.

This is incredibly useful in web services. Instead of passing bare int64
values around (where you might mix up user IDs and order IDs), you create
distinct types that the compiler can check for you.

Type aliases (with =) are mainly used for gradual code migration and are
less common in everyday code.

=============================================================================
*/

// Celsius and Fahrenheit are distinct types — you can't accidentally mix them.
type Celsius float64
type Fahrenheit float64

// CelsiusToFahrenheit converts a Celsius temperature to Fahrenheit.
// Note: we need explicit conversion because these are different types!
func CelsiusToFahrenheit(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

// FahrenheitToCelsius converts a Fahrenheit temperature to Celsius.
func FahrenheitToCelsius(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// DemoTypeDefinitions shows how custom type definitions create type safety.
func DemoTypeDefinitions() {
	boiling := Celsius(100)
	converted := CelsiusToFahrenheit(boiling)

	fmt.Println("Type definitions:")
	fmt.Printf("  %v°C = %v°F\n", boiling, converted)

	// This would be a compile error — different types!
	// var temp Celsius = Fahrenheit(72)  // COMPILE ERROR

	// You CAN convert between the type and its underlying type:
	var raw float64 = float64(boiling) // Celsius -> float64 is OK
	fmt.Printf("  Celsius as float64: %f\n", raw)
}

/*
=============================================================================
 STRING INTERNALS
=============================================================================

Strings in Go are immutable sequences of bytes. They're almost always
UTF-8 encoded, but technically a string can hold arbitrary bytes.

Key things to know:
  - len(s) returns the number of BYTES, not characters
  - Indexing s[i] gives you a BYTE, not a character
  - Use []rune(s) or range loop to iterate over characters
  - Strings are immutable — to modify, convert to []byte or []rune
  - String concatenation with + creates a new string each time
    (use strings.Builder for building strings in loops)

This matters for web services when you're processing user input that
might contain non-ASCII characters (names, addresses, emoji, etc.).

=============================================================================
*/

// DemoStringInternals shows the difference between bytes and runes in strings.
func DemoStringInternals() {
	s := "Hello, 世界! 🌍"

	fmt.Println("String internals:")
	fmt.Printf("  string: %s\n", s)
	fmt.Printf("  byte length: %d\n", len(s))
	fmt.Printf("  rune count: %d\n", len([]rune(s)))

	// Iterating over bytes vs runes:
	fmt.Println("  Byte iteration (first 13 bytes):")
	for i := 0; i < 13 && i < len(s); i++ {
		fmt.Printf("    s[%d] = %d (%c)\n", i, s[i], s[i])
	}

	// Range over a string iterates by RUNE, not by byte.
	// This is the correct way to process characters.
	fmt.Println("  Rune iteration:")
	for i, r := range s {
		fmt.Printf("    index=%d, rune=%c (U+%04X)\n", i, r, r)
	}
	// Notice: the index jumps when multi-byte characters are encountered!
}
