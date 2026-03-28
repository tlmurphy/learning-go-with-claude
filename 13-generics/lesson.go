package generics

/*
Module 13: Generics
====================

Generics (added in Go 1.18) let you write functions and types that work with
multiple types while maintaining type safety. Before generics, you had two
choices for polymorphic code: interfaces (runtime dispatch, potential type
assertions) or code generation. Generics give you a third option: compile-time
type parameterization.

Type Parameters and Constraints
-------------------------------
A generic function declares type parameters in square brackets:

  func Min[T constraints.Ordered](a, b T) T {
      if a < b {
          return a
      }
      return b
  }

  Min(3, 5)         // T is inferred as int
  Min("a", "b")     // T is inferred as string
  Min[float64](1, 2) // T is explicitly float64

Type parameters must satisfy a CONSTRAINT — an interface that specifies
what operations are allowed on the type.

Built-in Constraints
--------------------
  - any          — allows any type (alias for interface{})
  - comparable   — types that support == and != (map keys, for example)

The constraints package (golang.org/x/exp/constraints) provides:
  - Ordered    — types that support < > <= >= (integers, floats, strings)
  - Integer    — all integer types
  - Float      — all float types
  - Complex    — all complex types
  - Signed     — signed integer types
  - Unsigned   — unsigned integer types

Since we don't want to add an external dependency for this tutorial,
we'll define our own constraints inline. This is also totally normal in
production code.

Type Inference
--------------
Go can usually infer type arguments from function arguments:

  Min(3, 5)     // Go infers T = int
  Min("a", "b") // Go infers T = string

You only need explicit type arguments when:
  - The type can't be inferred from arguments
  - You want a specific type (e.g., Min[float64](1, 2))
  - You're instantiating a generic type (Set[string]{})

Generic Functions vs Generic Types
-----------------------------------
Generic functions: Parameterized behavior
  func Map[T, U any](slice []T, f func(T) U) []U

Generic types: Parameterized data structures
  type Stack[T any] struct { items []T }

Type Sets and Interface Constraints
------------------------------------
Constraints are interfaces with type sets. You can use unions:

  type Number interface {
      int | int8 | int16 | int32 | int64 |
      float32 | float64
  }

The ~ prefix means "any type whose underlying type is":

  type Stringish interface {
      ~string  // Includes string AND any type defined as `type MyStr string`
  }

This is important because named types (type UserID int) have different
types but the same underlying type.

When to Use Generics (and When NOT To)
--------------------------------------
USE generics for:
  - Container types (Stack, Queue, Set, Tree)
  - Utility functions over slices/maps (Map, Filter, Reduce, Contains)
  - When you need type safety across multiple types

DON'T use generics for:
  - Single-type code (just use the concrete type)
  - When interfaces work fine (behavior, not type, matters)
  - Premature abstraction ("I might need this to be generic someday")

The Go team's guidance: "Don't use generics until you find yourself writing
the same code three times for different types."

Practical Patterns
------------------
The slices package (standard library since Go 1.21) already has generic
helpers: slices.Contains, slices.Sort, slices.Map, etc. Before reaching
for custom generics, check if the standard library has what you need.
*/

import "fmt"

// ==========================================
// Defining Custom Constraints
// ==========================================

// Ordered is a constraint for types that support ordering operators.
// This is equivalent to constraints.Ordered from golang.org/x/exp.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Number is a constraint for numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// ==========================================
// Generic Functions
// ==========================================

// Min returns the smaller of two ordered values.
// Works with any type that supports the < operator.
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two ordered values.
func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Contains checks if a slice contains a value.
// The type must be comparable (supports ==).
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// ==========================================
// Map, Filter, Reduce
// ==========================================

// Map applies a function to every element of a slice, returning a new slice.
// This is a fundamental generic utility.
func Map[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// Filter returns a new slice containing only elements where the predicate
// returns true.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce reduces a slice to a single value by applying an accumulator function.
func Reduce[T any, U any](slice []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range slice {
		acc = f(acc, v)
	}
	return acc
}

// ==========================================
// Generic Data Structures
// ==========================================

// Stack is a generic LIFO data structure.
type Stack[T any] struct {
	items []T
}

// NewStack creates an empty stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item. Returns false if empty.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	last := len(s.items) - 1
	item := s.items[last]
	s.items = s.items[:last]
	return item, true
}

// Peek returns the top item without removing it. Returns false if empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the stack has no items.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// ==========================================
// Generic Set
// ==========================================

// Set is a generic unordered collection of unique values.
// T must be comparable because sets use a map internally.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates an empty set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// SetFrom creates a set from a slice of values.
func SetFrom[T comparable](values []T) *Set[T] {
	s := NewSet[T]()
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Add adds a value to the set.
func (s *Set[T]) Add(value T) {
	s.items[value] = struct{}{}
}

// Remove removes a value from the set.
func (s *Set[T]) Remove(value T) {
	delete(s.items, value)
}

// Contains checks if a value is in the set.
func (s *Set[T]) Contains(value T) bool {
	_, ok := s.items[value]
	return ok
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.items)
}

// Values returns all values as a slice (order is not guaranteed).
func (s *Set[T]) Values() []T {
	result := make([]T, 0, len(s.items))
	for v := range s.items {
		result = append(result, v)
	}
	return result
}

// Union returns a new set containing all elements from both sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for v := range s.items {
		result.Add(v)
	}
	for v := range other.items {
		result.Add(v)
	}
	return result
}

// Intersection returns a new set containing only elements present in both sets.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	// Iterate over the smaller set for efficiency
	smaller, larger := s, other
	if s.Len() > other.Len() {
		smaller, larger = other, s
	}
	for v := range smaller.items {
		if larger.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// ==========================================
// Type Inference Examples
// ==========================================

// DemoTypeInference shows cases where Go infers type arguments.
func DemoTypeInference() {
	// Type is inferred from arguments — no need to specify [int]
	_ = Min(3, 5)
	_ = Max("alpha", "beta")

	// Map infers both T and U from the function signature
	nums := []int{1, 2, 3}
	_ = Map(nums, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})

	// Sometimes you DO need explicit type arguments
	s := NewStack[int]() // Can't infer from no arguments
	s.Push(42)
}

// ==========================================
// Underlying Type Constraint (~)
// ==========================================

// UserID is a named type with int as its underlying type.
type UserID int

// OrderID is another named type based on int.
type OrderID int

// Sum works with any numeric type, including named types like UserID.
// The ~ prefix means "any type whose underlying type is..."
func Sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

// DemoUnderlyingTypes shows ~ constraint working with named types.
func DemoUnderlyingTypes() {
	// Works with plain int
	_ = Sum([]int{1, 2, 3})

	// Also works with UserID because ~int includes UserID
	_ = Sum([]UserID{UserID(1), UserID(2), UserID(3)})
}

// ==========================================
// Generic Result Type (Inspired by Rust)
// ==========================================

// Result represents either a successful value or an error.
// This provides a type-safe alternative to (T, error) tuples.
type Result[T any] struct {
	value T
	err   error
	ok    bool
}

// Ok creates a successful Result.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

// Err creates a failed Result.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// IsOk returns true if the Result contains a value.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr returns true if the Result contains an error.
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// Unwrap returns the value, panicking if it's an error.
// Use this only when you're sure the Result is Ok.
func (r Result[T]) Unwrap() T {
	if !r.ok {
		panic(fmt.Sprintf("called Unwrap on an error Result: %v", r.err))
	}
	return r.value
}

// UnwrapOr returns the value if Ok, or the provided default if Err.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.ok {
		return r.value
	}
	return defaultValue
}

// Error returns the error, or nil if Ok.
func (r Result[T]) Error() error {
	return r.err
}
