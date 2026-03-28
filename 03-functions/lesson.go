// Package functions covers Go's function features: multiple returns, variadic
// functions, closures, defer, and functions as first-class values.
package functions

import "fmt"

/*
=============================================================================
 FUNCTIONS IN GO
=============================================================================

Functions are the building blocks of Go programs. Go's approach to functions
has some distinctive features:

  - Multiple return values (not tuples — actually multiple values)
  - Named return values (useful for documentation, dangerous if overused)
  - Variadic functions (accepting variable numbers of arguments)
  - Functions as first-class values (pass them around, return them)
  - Closures (functions that capture their environment)
  - defer (guaranteed cleanup, even on panic)

These features combine to create Go's signature patterns: the error return
pattern, middleware chains, resource cleanup with defer, and functional
options for configuration.

=============================================================================
 MULTIPLE RETURN VALUES
=============================================================================

Go functions can return multiple values. This isn't syntactic sugar for
returning a tuple or struct — they're genuinely separate return values.

The most important use: the (value, error) pattern. Nearly every Go
function that can fail returns an error as its last value:

  file, err := os.Open("data.txt")
  user, err := db.FindUser(id)
  body, err := io.ReadAll(resp.Body)

This is Go's alternative to exceptions. It forces you to handle errors
at every call site. It might feel verbose, but it makes error handling
visible and explicit — you can see every error path by reading the code.

=============================================================================
*/

// Divide demonstrates the (value, error) pattern with multiple returns.
// In Go, errors are values, not exceptions. This function returns both
// the result and an error. The caller MUST check the error.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide %f by zero", a)
	}
	return a / b, nil
}

// DemoMultipleReturns shows how to call and handle multiple return values.
func DemoMultipleReturns() {
	// The standard pattern: call, check error, then use value.
	result, err := Divide(10, 3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result)
	}

	// When the error case occurs:
	result, err = Divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 0 = %.4f\n", result)
	}

	// You can use _ to discard a return value you don't need.
	// But NEVER discard errors in production code!
	result, _ = Divide(10, 2) // only OK if you're absolutely sure it won't fail
	fmt.Printf("10 / 2 = %.4f (error discarded)\n", result)
}

/*
=============================================================================
 NAMED RETURN VALUES
=============================================================================

Go lets you name your return values. This serves two purposes:

  1. Documentation: named returns document what each value means
  2. "Naked" return: you can use `return` without arguments

Named returns are great for documentation but the "naked return" feature
is controversial. Naked returns make it hard to see what's being returned,
especially in long functions. The Go community generally recommends:

  - DO name returns for documentation in complex functions
  - DO use naked returns ONLY in very short functions (< 10 lines)
  - DON'T use naked returns in long functions — be explicit

=============================================================================
*/

// MinMax uses named returns for documentation.
// In a short function like this, named returns are clear and helpful.
func MinMax(values []int) (min, max int) {
	if len(values) == 0 {
		return 0, 0 // explicit return is always fine
	}

	min, max = values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return // naked return — returns min and max. OK here because function is short.
}

/*
=============================================================================
 VARIADIC FUNCTIONS
=============================================================================

Variadic functions accept a variable number of arguments. The variadic
parameter must be the last parameter and it receives the arguments as a
slice.

  func Sum(nums ...int) int { }

You can call it as:
  Sum(1, 2, 3)           // individual arguments
  Sum(mySlice...)         // spread a slice (note the ...)

The fmt.Println function you've been using is variadic — that's why you
can pass any number of arguments to it!

=============================================================================
*/

// Sum demonstrates a variadic function.
// nums receives the arguments as an []int slice.
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// DemoVariadic shows variadic function calls.
func DemoVariadic() {
	// Call with individual arguments.
	fmt.Println("Sum(1,2,3) =", Sum(1, 2, 3))

	// Call with no arguments — nums is an empty (nil) slice.
	fmt.Println("Sum() =", Sum())

	// Spread a slice into a variadic function.
	numbers := []int{10, 20, 30}
	fmt.Println("Sum(numbers...) =", Sum(numbers...))

	// You can mix fixed and variadic parameters:
	// func Printf(format string, args ...interface{}) — first arg is fixed
}

/*
=============================================================================
 FUNCTIONS AS FIRST-CLASS VALUES
=============================================================================

In Go, functions are first-class values. You can:
  - Assign them to variables
  - Pass them as arguments to other functions
  - Return them from functions

This enables powerful patterns like callbacks, middleware, and strategies.

  var handler func(int) string    // variable holding a function
  func apply(f func(int) int, n int) int { return f(n) }  // function as parameter

This is foundational for Go web development. HTTP handlers, middleware,
and route configurations all use functions as values:

  http.HandleFunc("/api/users", handleUsers)  // passing a function
  handler = loggingMiddleware(handler)         // function returning a function

=============================================================================
*/

// ApplyToSlice takes a function and applies it to every element of a slice.
// This is a basic higher-order function — it takes a function as a parameter.
func ApplyToSlice(nums []int, f func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = f(n)
	}
	return result
}

// DemoFunctionsAsValues shows functions being passed around as values.
func DemoFunctionsAsValues() {
	double := func(n int) int { return n * 2 }
	square := func(n int) int { return n * n }

	nums := []int{1, 2, 3, 4, 5}

	fmt.Println("Original:", nums)
	fmt.Println("Doubled:", ApplyToSlice(nums, double))
	fmt.Println("Squared:", ApplyToSlice(nums, square))

	// Functions can also be declared inline (anonymous functions):
	fmt.Println("Plus 10:", ApplyToSlice(nums, func(n int) int { return n + 10 }))
}

/*
=============================================================================
 CLOSURES
=============================================================================

A closure is a function that references variables from its enclosing scope.
The function "closes over" those variables — it captures them by reference,
not by value.

This means:
  1. The closure can read and modify the captured variables
  2. The variables survive as long as the closure exists
  3. Multiple closures can share the same captured variables

Closures are used constantly in Go:
  - Counter functions that maintain state
  - Middleware that wraps handlers with extra behavior
  - Iterators and generators
  - Configuration with functional options

GOTCHA: Closures in loops capture the variable, not its value at that
moment. This is a common source of bugs:

  for i := 0; i < 5; i++ {
      go func() { fmt.Println(i) }()  // BUG: all goroutines see i=5
  }

  for i := 0; i < 5; i++ {
      i := i  // shadow with new variable (common fix)
      go func() { fmt.Println(i) }()  // OK: each goroutine has its own i
  }

=============================================================================
*/

// MakeCounter returns a closure that counts up from 0.
// Each call to the returned function increments and returns the count.
// The count variable is captured by the closure and persists between calls.
func MakeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// MakeAccumulator returns a closure that accumulates a running total.
func MakeAccumulator() func(int) int {
	total := 0
	return func(n int) int {
		total += n
		return total
	}
}

// DemoClosures shows closures in action.
func DemoClosures() {
	// Each call to MakeCounter creates a NEW count variable.
	counter1 := MakeCounter()
	counter2 := MakeCounter()

	fmt.Println("Counter 1:", counter1(), counter1(), counter1()) // 1, 2, 3
	fmt.Println("Counter 2:", counter2(), counter2())             // 1, 2 (independent!)

	acc := MakeAccumulator()
	fmt.Println("Accumulator: add 10 ->", acc(10))  // 10
	fmt.Println("Accumulator: add 20 ->", acc(20))  // 30
	fmt.Println("Accumulator: add 5  ->", acc(5))   // 35
}

/*
=============================================================================
 DEFER
=============================================================================

defer schedules a function call to run when the enclosing function returns.
Deferred calls are executed in LIFO (Last In, First Out) order — like a
stack.

Key rules:
  1. Arguments to deferred calls are evaluated immediately (when defer runs)
  2. The deferred function itself executes later (when the function returns)
  3. Multiple defers execute in reverse order (LIFO)
  4. Deferred functions can read and modify named return values

defer is most commonly used for cleanup:
  - Closing files: defer f.Close()
  - Releasing locks: defer mu.Unlock()
  - Closing database connections: defer db.Close()
  - Closing HTTP response bodies: defer resp.Body.Close()

The guarantee that defer runs even if the function panics makes it
essential for resource management. It's Go's answer to try/finally.

=============================================================================
*/

// DemoDefer shows how defer works, including LIFO ordering.
func DemoDefer() {
	// Defers execute in LIFO order — last defer runs first.
	fmt.Println("Defer order (watch the numbers):")
	for i := 1; i <= 3; i++ {
		defer fmt.Printf("  deferred: %d\n", i) // arguments evaluated NOW
	}
	fmt.Println("  after loop, before function returns")
	// Output will be:
	//   after loop, before function returns
	//   deferred: 3
	//   deferred: 2
	//   deferred: 1
}

// DemoDeferArgEvaluation shows that defer arguments are evaluated immediately.
func DemoDeferArgEvaluation() {
	x := 10
	defer fmt.Printf("  deferred x = %d\n", x) // x is evaluated NOW (10)
	x = 20
	fmt.Printf("  current x = %d\n", x) // prints 20
	// deferred print will show 10, not 20!
}

/*
=============================================================================
 INIT FUNCTIONS
=============================================================================

Every package can have one or more init() functions. These run automatically
when the package is loaded, before main() runs. They're used for:

  - Registering database drivers
  - Initializing package-level state
  - Validating configuration

Rules:
  - init() takes no arguments and returns no values
  - A file can have multiple init() functions (but don't overdo it)
  - init() functions run in the order the files are presented to the compiler
  - Package init() runs after all variable declarations in that package

You shouldn't overuse init() — it can make code harder to test and reason
about. Prefer explicit initialization in main() when possible.

We won't demonstrate init() here because it has side effects at import
time, but it's good to know about.

=============================================================================
*/

/*
=============================================================================
 METHOD EXPRESSIONS AND METHOD VALUES (PREVIEW)
=============================================================================

Go lets you use methods as function values in two ways. This is a preview
for when we cover methods and interfaces:

  // Method value — bound to a specific receiver
  greeter := myGreeter.Greet  // captures the receiver
  greeter("World")            // calls myGreeter.Greet("World")

  // Method expression — requires passing the receiver explicitly
  greet := Greeter.Greet      // unbound method
  greet(myGreeter, "World")   // must pass the receiver

Method values are used constantly with http.Handler:
  mux.HandleFunc("/users", server.handleUsers)  // method value

We'll explore this fully in the methods and interfaces module.

=============================================================================
*/

// DemoFunctionComposition shows building new functions from existing ones.
func DemoFunctionComposition() {
	// Compose takes two functions and returns a new function that
	// applies them in sequence: compose(f, g)(x) = f(g(x))
	compose := func(f, g func(int) int) func(int) int {
		return func(x int) int {
			return f(g(x))
		}
	}

	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }

	doubleAndAddOne := compose(addOne, double) // first double, then add one
	addOneAndDouble := compose(double, addOne) // first add one, then double

	fmt.Println("Compose double then addOne:")
	fmt.Printf("  f(3) = %d (3*2=6, 6+1=7)\n", doubleAndAddOne(3))

	fmt.Println("Compose addOne then double:")
	fmt.Printf("  f(3) = %d (3+1=4, 4*2=8)\n", addOneAndDouble(3))
}
