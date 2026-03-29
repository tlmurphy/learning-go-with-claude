package variables

/*
=============================================================================
 EXERCISES: Variables and Types
=============================================================================

 Work through these exercises in order. Each one builds on concepts from
 the lesson and from previous exercises. Run the tests with:

   make test 01

 Tip: Run a single test at a time while working:

   go test -v -run TestDeclareVariables ./01-variables-and-types/

=============================================================================
*/

// Exercise 1: DeclareVariables
//
// Return three values demonstrating different declaration styles:
//   - an int with value 42
//   - a string with value "hello"
//   - a bool with value true
//
// The point here is simple: get comfortable with Go's return syntax and
// basic types. In a web service, you'd return values like this from
// functions that parse request parameters.
func DeclareVariables() (int, string, bool) {
	// YOUR CODE HERE
	return 0, "", false
}

// Exercise 2: ZeroValues
//
// Return the zero values for each of these types WITHOUT explicitly
// assigning any values. Just declare variables with `var` and return them.
//
// Returns: int, float64, string, bool
//
// Understanding zero values is critical in Go. When you read a struct
// field that hasn't been set (like a JSON-decoded request body with
// missing fields), you'll get zero values. You need to know what they are.
func ZeroValues() (int, float64, string, bool) {
	// YOUR CODE HERE
	// Hint: declare variables with var (no assignment) and return them.
	return 0, 0, "", false
}

// Exercise 3: Define a set of constants using iota
//
// Create a function that returns the integer values of these HTTP status
// categories, which should be defined as constants with iota:
//
//   StatusInfo         = 1  (1xx informational)
//   StatusSuccess      = 2  (2xx success)
//   StatusRedirect     = 3  (3xx redirection)
//   StatusClientError  = 4  (4xx client errors)
//   StatusServerError  = 5  (5xx server errors)
//
// Define the constants using iota so that StatusInfo starts at 1.
// Return them all from the function.
//
// Hint: iota starts at 0, but you need to start at 1. There are two
// common approaches: use iota + 1 in the first constant, or use a blank
// identifier _ to skip 0.

// HTTPStatusCategory represents a category of HTTP status codes.
type HTTPStatusCategory int

const (
	// YOUR CODE HERE — define the constants using iota
	// StatusInfo HTTPStatusCategory = ...
	StatusInfo        HTTPStatusCategory = iota // placeholder, fix this
	StatusSuccess     HTTPStatusCategory = iota
	StatusRedirect    HTTPStatusCategory = iota
	StatusClientError HTTPStatusCategory = iota
	StatusServerError HTTPStatusCategory = iota
)

// GetStatusCategories returns all five HTTP status category values.
func GetStatusCategories() (int, int, int, int, int) {
	// YOUR CODE HERE
	// Return the integer values of each status category constant.
	return 0, 0, 0, 0, 0
}

// Exercise 4: Type Definitions — Temperature Converter
//
// Implement these conversion functions between Kelvin and Celsius.
// Kelvin = Celsius + 273.15
//
// Notice how the type system prevents you from accidentally passing a
// Celsius value where a Kelvin is expected, and vice versa.

// Kelvin is a temperature in Kelvin.
type Kelvin float64

// CelsiusToKelvin converts a Celsius temperature to Kelvin.
// Formula: K = C + 273.15
func CelsiusToKelvin(c Celsius) Kelvin {
	// YOUR CODE HERE
	return 0
}

// KelvinToCelsius converts a Kelvin temperature to Celsius.
// Formula: C = K - 273.15
func KelvinToCelsius(k Kelvin) Celsius {
	// YOUR CODE HERE
	return 0
}

// AbsoluteZeroCelsius returns absolute zero in Celsius.
// (This is a good test that your conversion is correct: 0 K = -273.15 C)
func AbsoluteZeroCelsius() Celsius {
	// YOUR CODE HERE
	return 0
}

// Exercise 5: SwapValues
//
// Go supports multiple assignment, which makes swapping elegant:
//
//	a, b = b, a
//
// Implement a function that takes two ints and returns them swapped.
// Then implement SwapStrings that does the same for strings.
//
// This is trivial in Go but requires a temp variable in many languages.
// Multiple return values + multiple assignment is a Go superpower.
func SwapInts(a, b int) (int, int) {
	// YOUR CODE HERE
	return 0, 0
}

// SwapStrings takes two strings and returns them in reversed order.
func SwapStrings(a, b string) (string, string) {
	// YOUR CODE HERE
	return "", ""
}

