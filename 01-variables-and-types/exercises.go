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

// Exercise 4: SafeIntToInt8
//
// Convert an int to int8 safely. If the value is outside the range of
// int8 (-128 to 127), return 0 and false. Otherwise return the converted
// value and true.
//
// This teaches two things:
//  1. Explicit type conversion (Go won't do it for you)
//  2. The comma-ok pattern (returning a success bool alongside a value)
//
// You'll see the comma-ok pattern everywhere in Go:
//
//	value, ok := myMap[key]
//	value, err := strconv.Atoi(s)
func SafeIntToInt8(n int) (int8, bool) {
	// YOUR CODE HERE
	return 0, false
}

// Exercise 5: StringByteRuneAnalysis
//
// Given a string, return:
//   - the number of bytes in the string
//   - the number of runes (characters) in the string
//   - the first rune (character) as a rune
//   - the last rune (character) as a rune
//
// For an empty string, return (0, 0, 0, 0).
//
// This is practical: when validating user input in a web form, you
// usually care about character count, not byte count. A username limit
// of "20 characters" should allow 20 emoji, not just 5 (since emoji
// can be 4 bytes each).
//
// Hint: convert the string to []rune for character-level operations.
func StringByteRuneAnalysis(s string) (int, int, rune, rune) {
	// YOUR CODE HERE
	return 0, 0, 0, 0
}

// Exercise 6: Type Definitions — Temperature Converter
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

// Exercise 7: SwapValues
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

// Exercise 8: ParseAndDescribe (Composite Exercise)
//
// This exercise combines several concepts from the module:
//   - Type definitions
//   - Type conversions
//   - Constants
//   - Multiple return values
//   - String operations
//
// Given a UserScore (a type definition over int), return:
//   - The score as a float64 percentage (score / MaxScore * 100)
//   - A rating string based on the percentage:
//       >= 90: "excellent"
//       >= 70: "good"
//       >= 50: "fair"
//       <  50: "needs improvement"
//   - Whether the score is passing (>= 50%)
//
// The MaxScore constant is 200.

// UserScore is a distinct type representing a user's score.
type UserScore int

// MaxScore is the maximum possible score.
const MaxScore UserScore = 200

// EvaluateScore analyzes a UserScore and returns a percentage, rating, and pass/fail.
func EvaluateScore(score UserScore) (percentage float64, rating string, passing bool) {
	// YOUR CODE HERE
	return 0, "", false
}
