package controlflow

import (
	"strings"
)

/*
=============================================================================
 EXERCISES: Control Flow
=============================================================================

 Work through these exercises in order. Run tests with:

   go test -v ./02-control-flow/

 Run a single test:

   go test -v -run TestFizzBuzz ./02-control-flow/

=============================================================================
*/

// Exercise 1: FizzBuzz
//
// The classic — but write it in idiomatic Go.
//
// Given an integer n, return a string slice containing the FizzBuzz sequence
// from 1 to n (inclusive):
//   - If the number is divisible by 3, use "Fizz"
//   - If the number is divisible by 5, use "Buzz"
//   - If divisible by both 3 and 5, use "FizzBuzz"
//   - Otherwise, use the number as a string (e.g., "1", "2", "4")
//
// Use fmt.Sprintf or strconv.Itoa for number-to-string conversion.
// Think about the order of your conditions!
func FizzBuzz(n int) []string {
	// YOUR CODE HERE
	return nil
}

// Exercise 2: ClassifyAge
//
// Use an if statement with an init statement to classify an age string.
//
// Given an age as an int, return a category string:
//
//	age < 0:    "invalid"
//	age < 13:   "child"
//	age < 18:   "teenager"
//	age < 65:   "adult"
//	age >= 65:  "senior"
//
// Bonus insight: In a web handler, you'd often parse the age from a query
// parameter first, then classify:
//
//	if age, err := strconv.Atoi(ageParam); err != nil { ... }
func ClassifyAge(age int) string {
	// YOUR CODE HERE
	return ""
}

// Exercise 3: DayType
//
// Use a switch statement to categorize days of the week.
//
// Given a day name (case-insensitive), return:
//
//	"Monday" through "Friday" -> "weekday"
//	"Saturday", "Sunday"      -> "weekend"
//	anything else              -> "invalid"
//
// Hint: Normalize the input first. Use strings.ToLower or strings.Title.
// Switch cases can have multiple values: case "a", "b", "c":
func DayType(day string) string {
	// YOUR CODE HERE
	_ = strings.ToLower(day) // hint: you'll want to use this
	return ""
}

// Exercise 4: CountUnicodeCategories
//
// Range over a string and count characters by category.
//
// Given a string, return the count of:
//   - uppercase letters
//   - lowercase letters
//   - digits
//   - spaces (space, tab, newline)
//   - other characters (punctuation, symbols, etc.)
//
// Use a range loop (which iterates by rune) and a tagless switch
// with unicode package functions: unicode.IsUpper, unicode.IsLower,
// unicode.IsDigit, unicode.IsSpace.
//
// This is practical: you might validate password complexity this way.
func CountUnicodeCategories(s string) (upper, lower, digits, spaces, other int) {
	// YOUR CODE HERE
	return 0, 0, 0, 0, 0
}

// Exercise 5: FindInMatrix
//
// Search a 2D integer matrix for a target value using nested loops
// with a labeled break.
//
// Return the row index, column index, and whether the value was found.
// If not found, return -1, -1, false.
//
// You MUST use a labeled break to exit both loops when found. This is
// the primary pattern where labeled breaks shine.
func FindInMatrix(matrix [][]int, target int) (row, col int, found bool) {
	// YOUR CODE HERE
	return -1, -1, false
}

// Exercise 6: StateMachine
//
// Build a simple state machine that processes a string of commands.
//
// The machine has three states: "idle", "running", "stopped"
// Commands and transitions:
//
//	"idle"    + "start" -> "running"
//	"running" + "stop"  -> "stopped"
//	"running" + "pause" -> "idle"
//	"stopped" + "reset" -> "idle"
//	Any other combination -> state stays the same
//
// Process each command in order and return the final state.
// Start in the "idle" state.
//
// Hint: Use a switch inside a for loop. The switch can match on state,
// and inner switches or if statements can match on the command.
func StateMachine(commands []string) string {
	// YOUR CODE HERE
	return ""
}

// Exercise 7: CollatzSteps
//
// The Collatz conjecture: start with any positive integer n.
//   - If n is even, divide by 2
//   - If n is odd, multiply by 3 and add 1
//   - Repeat until you reach 1
//
// Return the number of steps it takes to reach 1.
// If n <= 0, return -1 (invalid input).
// If n == 1, return 0 (already at 1).
//
// Example: n=6 -> 6, 3, 10, 5, 16, 8, 4, 2, 1 = 8 steps
//
// Use a for loop (while-style). Think about which loop form is most natural.
func CollatzSteps(n int) int {
	// YOUR CODE HERE
	return 0
}

// Exercise 8: ProcessRecords
//
// This exercise combines multiple control flow patterns.
//
// Given a slice of string records in "key:value" format, process them
// according to these rules:
//
//  1. Skip empty strings (continue)
//  2. If a record doesn't contain ":", add it to the errors slice
//  3. If the key (part before ":") is "STOP", stop processing immediately
//     (break — do NOT include "STOP" in results or errors)
//  4. If the key is empty (record starts with ":"), add to errors
//  5. If the value (part after ":") is empty, use "default" as the value
//  6. Otherwise, add the key-value pair to the results map
//
// Return the results map and errors slice.
//
// This simulates parsing a simple configuration format — the kind of
// thing you'd do when processing form data or config files in a web service.
func ProcessRecords(records []string) (map[string]string, []string) {
	// YOUR CODE HERE
	return nil, nil
}
