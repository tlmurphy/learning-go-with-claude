package generics

/*
Exercises: Generics
====================

These exercises build from simple generic functions to more complex generic
data structures. Each exercise reinforces a different aspect of Go's generics.
*/

import "fmt"

// =============================================================================
// Exercise 1: Generic Min/Max for Ordered Types
// =============================================================================

// MinSlice returns the minimum value in a slice of ordered values.
// Returns the zero value and false if the slice is empty.
func MinSlice[T Ordered](values []T) (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// MaxSlice returns the maximum value in a slice of ordered values.
// Returns the zero value and false if the slice is empty.
func MaxSlice[T Ordered](values []T) (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// =============================================================================
// Exercise 2: Generic Map, Filter, Reduce for Slices
// =============================================================================

// MapSlice applies a function to every element and returns the results.
// This is your own implementation — don't call the lesson's Map function.
func MapSlice[T any, U any](slice []T, f func(T) U) []U {
	// YOUR CODE HERE
	return nil
}

// FilterSlice returns elements where the predicate is true.
func FilterSlice[T any](slice []T, predicate func(T) bool) []T {
	// YOUR CODE HERE
	return nil
}

// ReduceSlice reduces a slice to a single value using an accumulator.
func ReduceSlice[T any, U any](slice []T, initial U, f func(U, T) U) U {
	// YOUR CODE HERE
	var zero U
	return zero
}

// =============================================================================
// Exercise 3: Generic Stack
// =============================================================================

// ExerciseStack is a generic LIFO stack.
// Implement it from scratch (don't use the lesson's Stack).
type ExerciseStack[T any] struct {
	elements []T
}

// NewExerciseStack creates an empty stack.
func NewExerciseStack[T any]() *ExerciseStack[T] {
	// YOUR CODE HERE
	return nil
}

// Push adds an item to the top of the stack.
func (s *ExerciseStack[T]) Push(item T) {
	// YOUR CODE HERE
}

// Pop removes and returns the top item.
// Returns (zero value, false) if empty.
func (s *ExerciseStack[T]) Pop() (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// Peek returns the top item without removing it.
// Returns (zero value, false) if empty.
func (s *ExerciseStack[T]) Peek() (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// Size returns the number of items in the stack.
func (s *ExerciseStack[T]) Size() int {
	// YOUR CODE HERE
	return 0
}

// ToSlice returns all items as a slice (bottom to top).
func (s *ExerciseStack[T]) ToSlice() []T {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 4: Generic Set with Union and Intersection
// =============================================================================

// ExerciseSet is a generic set of unique comparable values.
type ExerciseSet[T comparable] struct {
	data map[T]struct{}
}

// NewExerciseSet creates an empty set.
func NewExerciseSet[T comparable]() *ExerciseSet[T] {
	// YOUR CODE HERE
	return nil
}

// Add adds a value to the set.
func (s *ExerciseSet[T]) Add(value T) {
	// YOUR CODE HERE
}

// Remove removes a value from the set.
func (s *ExerciseSet[T]) Remove(value T) {
	// YOUR CODE HERE
}

// Contains checks if a value is in the set.
func (s *ExerciseSet[T]) Contains(value T) bool {
	// YOUR CODE HERE
	return false
}

// Size returns the number of elements.
func (s *ExerciseSet[T]) Size() int {
	// YOUR CODE HERE
	return 0
}

// Union returns a new set with elements from both sets.
func (s *ExerciseSet[T]) Union(other *ExerciseSet[T]) *ExerciseSet[T] {
	// YOUR CODE HERE
	return nil
}

// Intersection returns a new set with elements common to both sets.
func (s *ExerciseSet[T]) Intersection(other *ExerciseSet[T]) *ExerciseSet[T] {
	// YOUR CODE HERE
	return nil
}

// Difference returns a new set with elements in s but not in other.
func (s *ExerciseSet[T]) Difference(other *ExerciseSet[T]) *ExerciseSet[T] {
	// YOUR CODE HERE
	return nil
}

// =============================================================================
// Exercise 5: Generic Cache
// =============================================================================

// GenericCache is a type-safe cache with string keys and values of type T.
type GenericCache[T any] struct {
	data map[string]T
}

// NewGenericCache creates an empty cache.
func NewGenericCache[T any]() *GenericCache[T] {
	// YOUR CODE HERE
	return nil
}

// Set stores a value with the given key.
func (c *GenericCache[T]) Set(key string, value T) {
	// YOUR CODE HERE
}

// Get retrieves a value by key. Returns (value, true) if found,
// (zero value, false) if not.
func (c *GenericCache[T]) Get(key string) (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// Delete removes a key from the cache.
func (c *GenericCache[T]) Delete(key string) {
	// YOUR CODE HERE
}

// Keys returns all keys in the cache.
func (c *GenericCache[T]) Keys() []string {
	// YOUR CODE HERE
	return nil
}

// Size returns the number of entries.
func (c *GenericCache[T]) Size() int {
	// YOUR CODE HERE
	return 0
}

// =============================================================================
// Exercise 6: Custom Type Constraint
// =============================================================================

// Summable is a constraint for types that support the + operator.
// Define it to include all integer and float types.
type Summable interface {
	// YOUR CODE HERE — define the type set
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// SumAll returns the sum of all values in a slice.
// Uses the Summable constraint.
func SumAll[T Summable](values []T) T {
	// YOUR CODE HERE
	var zero T
	return zero
}

// Average returns the arithmetic mean of a slice of values.
// Returns 0 if the slice is empty.
func Average[T Summable](values []T) float64 {
	// YOUR CODE HERE
	return 0
}

// =============================================================================
// Exercise 7: Generic Result Type
// =============================================================================

// ExerciseResult represents a computation that either succeeded with a
// value of type T, or failed with an error.
type ExerciseResult[T any] struct {
	value T
	err   error
	isOk  bool
}

// NewOk creates a successful result containing the given value.
func NewOk[T any](value T) ExerciseResult[T] {
	// YOUR CODE HERE
	return ExerciseResult[T]{}
}

// NewErr creates a failed result containing the given error.
func NewErr[T any](err error) ExerciseResult[T] {
	// YOUR CODE HERE
	return ExerciseResult[T]{}
}

// IsOk returns true if the result is successful.
func (r ExerciseResult[T]) IsOk() bool {
	// YOUR CODE HERE
	return false
}

// IsErr returns true if the result is an error.
func (r ExerciseResult[T]) IsErr() bool {
	// YOUR CODE HERE
	return false
}

// Value returns the value and a boolean indicating success.
// Returns (zero, false) if the result is an error.
func (r ExerciseResult[T]) Value() (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// Error returns the error, or nil if successful.
func (r ExerciseResult[T]) Error() error {
	// YOUR CODE HERE
	return nil
}

// UnwrapOrDefault returns the value if Ok, or the provided default if Err.
func (r ExerciseResult[T]) UnwrapOrDefault(defaultVal T) T {
	// YOUR CODE HERE
	var zero T
	return zero
}

// MapResult transforms the value inside a Result using the given function.
// If the Result is Err, the error is propagated without calling f.
func MapResult[T any, U any](r ExerciseResult[T], f func(T) U) ExerciseResult[U] {
	// YOUR CODE HERE
	return ExerciseResult[U]{}
}

// =============================================================================
// Exercise 8: Generic Linked List
// =============================================================================

// Node is a generic linked list node.
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// LinkedList is a generic singly linked list.
type LinkedList[T any] struct {
	head *Node[T]
	size int
}

// NewLinkedList creates an empty linked list.
func NewLinkedList[T any]() *LinkedList[T] {
	// YOUR CODE HERE
	return nil
}

// Prepend adds a value to the front of the list.
func (l *LinkedList[T]) Prepend(value T) {
	// YOUR CODE HERE
}

// Append adds a value to the end of the list.
func (l *LinkedList[T]) Append(value T) {
	// YOUR CODE HERE
}

// Head returns the first value in the list.
// Returns (zero, false) if the list is empty.
func (l *LinkedList[T]) Head() (T, bool) {
	// YOUR CODE HERE
	var zero T
	return zero, false
}

// Len returns the number of elements.
func (l *LinkedList[T]) Len() int {
	// YOUR CODE HERE
	return 0
}

// ToSlice converts the linked list to a slice, preserving order.
func (l *LinkedList[T]) ToSlice() []T {
	// YOUR CODE HERE
	return nil
}

// ForEach applies a function to every element in the list.
func (l *LinkedList[T]) ForEach(f func(T)) {
	// YOUR CODE HERE
}

// String returns a string representation like "[1 -> 2 -> 3]".
// Returns "[]" for an empty list.
func (l *LinkedList[T]) String() string {
	if l.head == nil {
		return "[]"
	}
	result := "["
	current := l.head
	for current != nil {
		result += fmt.Sprintf("%v", current.Value)
		if current.Next != nil {
			result += " -> "
		}
		current = current.Next
	}
	result += "]"
	return result
}
