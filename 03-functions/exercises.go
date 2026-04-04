package functions

import "fmt"

/*
=============================================================================
 EXERCISES: Functions
=============================================================================

 Work through these exercises in order. Run tests with:

   make test 03

 Run a single test:

   go test -v -run TestSafeDivide ./03-functions/

=============================================================================
*/

// Exercise 1: SafeDivide
//
// Implement a safe division function that returns:
//   - quotient (integer division)
//   - remainder
//   - error (if dividing by zero)
//
// For zero divisor, return 0, 0, and an error with the message:
//
//	"division by zero"
//
// This is the fundamental Go pattern: (result, error). You'll write
// hundreds of functions like this in Go web services.
//
// Example: SafeDivide(17, 5) returns 3, 2, nil
// Example: SafeDivide(10, 0) returns 0, 0, error("division by zero")
func SafeDivide(a, b int) (quotient, remainder int, err error) {
	// YOUR CODE HERE
	return 0, 0, nil
}

// Exercise 2: VariadicSum and VariadicAverage
//
// VariadicSum returns the sum of all provided integers.
// If no arguments are provided, return 0.
//
// VariadicAverage returns the average of all provided float64 values.
// If no arguments are provided, return 0 and an error with message
// "no values provided". Otherwise return the average and nil.
//
// Think about edge cases: what if someone calls Average()? Division by
// zero! This is why returning an error is important.
func VariadicSum(nums ...int) int {
	// YOUR CODE HERE
	return 0
}

// VariadicAverage returns the average of the given float64 values.
// Returns 0 and an error if no values are provided.
func VariadicAverage(nums ...float64) (float64, error) {
	// YOUR CODE HERE
	return 0, nil
}

// Exercise 3: Higher-Order Functions
//
// Implement Filter, Map, and Reduce for int slices.
//
// These are fundamental functional programming patterns that Go supports
// through first-class functions. You'll see similar patterns in Go web
// frameworks for request processing pipelines.

// Filter returns a new slice containing only elements where the predicate
// returns true. The original slice must not be modified.
//
// Example: Filter([]int{1,2,3,4,5}, func(n int) bool { return n > 3 })
//
//	returns []int{4, 5}
func Filter(nums []int, predicate func(int) bool) []int {
	// YOUR CODE HERE
	return nil
}

// Map applies a transform function to each element and returns a new slice.
// The original slice must not be modified.
//
// Example: Map([]int{1,2,3}, func(n int) int { return n * 2 })
//
//	returns []int{2, 4, 6}
func Map(nums []int, transform func(int) int) []int {
	// YOUR CODE HERE
	return nil
}

// Reduce reduces a slice to a single value by applying an accumulator
// function. The initial value is the starting value for accumulation.
//
// Example: Reduce([]int{1,2,3,4}, 0, func(acc, n int) int { return acc + n })
//
//	returns 10
func Reduce(nums []int, initial int, accumulator func(int, int) int) int {
	// YOUR CODE HERE
	return 0
}

// Exercise 4: MakeCounter
//
// Create a closure-based counter. NewCounter takes a starting value and
// returns three functions:
//   - increment: adds 1 to the counter and returns the new value
//   - decrement: subtracts 1 from the counter and returns the new value
//   - value: returns the current counter value without changing it
//
// All three functions must share the same counter state.
// This demonstrates how closures capture variables by reference.
//
// Example:
//
//	inc, dec, val := NewCounter(10)
//	inc()  // returns 11
//	inc()  // returns 12
//	dec()  // returns 11
//	val()  // returns 11
func NewCounter(start int) (increment, decrement, value func() int) {
	// YOUR CODE HERE
	return nil, nil, nil
}

// Exercise 5: Middleware Pattern
//
// This is a preview of how Go web frameworks work.
//
// Implement a Logger "middleware" that wraps a function. Given a function
// that takes a string and returns a string, return a new function that:
//  1. Records that it was called (append to the log slice)
//  2. Calls the original function
//  3. Records the result (append to the log slice)
//  4. Returns the original result
//
// The log entries should be:
//
//	"calling with: <input>"
//	"returned: <result>"
//
// The log slice is provided as a pointer so the wrapper can append to it.
// To append to a pointer to a slice, dereference both sides:
//
//	*log = append(*log, "new entry")
//
// In real Go web servers, this exact pattern is used for logging, auth,
// rate limiting, etc:
//
//	handler = loggingMiddleware(handler)
//	handler = authMiddleware(handler)
func Logger(fn func(string) string, log *[]string) func(string) string {
	// YOUR CODE HERE
	return nil
}

// Exercise 6: SafeCall
//
// This exercise teaches defer's most important real-world use: recovering
// from panics. In Go web servers, a panic in a handler would crash the
// entire server — middleware uses exactly this pattern to catch panics
// and return an error response instead.
//
// Implement SafeCall: it calls fn() and returns its result. If fn panics,
// recover the panic and return an empty string and an error containing
// the panic message.
//
// The error message should be: "panic recovered: <panic value>"
//
// Example:
//
//	SafeCall(func() string { return "hello" })       // returns "hello", nil
//	SafeCall(func() string { panic("oh no") })       // returns "", error("panic recovered: oh no")
//
// Hint: defer runs even when a function panics. Use recover() inside a
// deferred function to catch the panic. You'll need named return values
// so the deferred function can modify the return values.
func SafeCall(fn func() string) (result string, err error) {
	// YOUR CODE HERE
	return
}

// Exercise 7: Compose
//
// Implement function composition. Compose takes two functions f and g,
// and returns a new function that computes f(g(x)).
//
// The mathematical notation is (f . g)(x) = f(g(x)).
// g is applied first, then f is applied to g's result.
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	addOne := func(x int) int { return x + 1 }
//	doubleThenAdd := Compose(addOne, double)
//	doubleThenAdd(3) // = addOne(double(3)) = addOne(6) = 7
func Compose(f, g func(int) int) func(int) int {
	// YOUR CODE HERE
	return nil
}

// Exercise 8: Memoize
//
// Implement a memoized version of a function that takes an int and
// returns an int. The first time a particular input is seen, compute
// and cache the result. On subsequent calls with the same input, return
// the cached result.
//
// Return the memoized function AND a function that returns how many
// times the original function was actually called (cache misses).
//
// This is a practical pattern — in web services, you'd memoize expensive
// database queries or API calls.
//
// Example:
//
//	square := func(n int) int { return n * n }
//	memoSquare, callCount := Memoize(square)
//	memoSquare(4)  // computes 16, callCount() returns 1
//	memoSquare(4)  // returns cached 16, callCount() still returns 1
//	memoSquare(5)  // computes 25, callCount() returns 2
func Memoize(fn func(int) int) (memoized func(int) int, callCount func() int) {
	// YOUR CODE HERE
	_ = fmt.Sprintf // hint: you won't need fmt, but the import is here for other functions
	return nil, nil
}
