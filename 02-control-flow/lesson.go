// Package controlflow covers Go's control structures: if/else, switch, for loops,
// range, labels, and the patterns that make Go control flow distinctive.
package controlflow

import (
	"fmt"
	"strings"
	"unicode"
)

/*
=============================================================================
 CONTROL FLOW IN GO
=============================================================================

Go's control flow is intentionally simple. There are no while loops, no
do-while loops, no ternary operators. This isn't a limitation — it's a
design choice. When there's only one loop keyword (for) and no ternary,
every Go developer reads control flow the same way.

The key patterns you need to know:
  - if/else with init statements (the error-handling workhorse)
  - switch (more powerful than in most languages)
  - for (the only loop, but it wears many hats)
  - range (the idiomatic way to iterate)

=============================================================================
 IF/ELSE WITH INIT STATEMENTS
=============================================================================

Go's if statement has a superpower: the init statement. You can declare
a variable that's scoped to the if/else block:

  if err := doSomething(); err != nil {
      // handle error
  }
  // err doesn't exist here — it's scoped to the if block

This pattern is EVERYWHERE in Go. You'll see it hundreds of times in any
Go codebase, especially for error handling:

  if user, err := db.FindUser(id); err != nil {
      http.Error(w, "user not found", 404)
      return
  } else {
      json.NewEncoder(w).Encode(user)
  }

The init statement keeps variables tightly scoped, preventing them from
polluting the surrounding function scope. This is a deliberate design
choice that makes code easier to reason about.

=============================================================================
*/

// DemoIfElse shows Go's if/else, including init statements.
func DemoIfElse() {
	// Basic if/else — note: no parentheses around the condition!
	// Coming from C/Java/JS, this feels weird at first.
	x := 42
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	}

	// If with init statement — the variable is scoped to the if/else block.
	// This is the idiomatic way to handle operations that can fail.
	if length := len("hello"); length > 3 {
		fmt.Printf("  long string (length=%d)\n", length)
	} else {
		fmt.Printf("  short string (length=%d)\n", length)
	}
	// length doesn't exist here — clean scope!

	// Real-world pattern: parsing with validation.
	data := "key=value"
	if idx := strings.Index(data, "="); idx >= 0 {
		key := data[:idx]
		value := data[idx+1:]
		fmt.Printf("  parsed: key=%q, value=%q\n", key, value)
	} else {
		fmt.Println("  no '=' found in data")
	}

	// Go has NO ternary operator. This is intentional.
	// Instead of:  result = condition ? a : b
	// You write:
	status := "inactive"
	if x > 0 {
		status = "active"
	}
	fmt.Printf("  status: %s\n", status)
}

/*
=============================================================================
 SWITCH STATEMENTS
=============================================================================

Go's switch is much more powerful than C's:
  - No fall-through by default (no need for break!)
  - Cases can be expressions, not just constants
  - Tagless switch acts like if/else chains but reads cleaner
  - Type switch (preview for interfaces module)

The lack of fall-through prevents a massive class of bugs. If you DO want
fall-through (rare), use the explicit `fallthrough` keyword.

=============================================================================
*/

// DemoSwitch shows Go's switch statement variations.
func DemoSwitch() {
	// Expression switch — the most common form.
	day := "Wednesday"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println("  weekday")
	case "Saturday", "Sunday":
		fmt.Println("  weekend")
	default:
		fmt.Println("  unknown day")
	}

	// No fall-through by default! Each case breaks automatically.
	// This prevents the common bug in C/Java where you forget `break`.
	// If you actually WANT fall-through (rare), use `fallthrough`:
	n := 1
	switch n {
	case 1:
		fmt.Println("  one")
		fallthrough // explicitly fall through to next case
	case 2:
		fmt.Println("  two (or fell through from one)")
	case 3:
		fmt.Println("  three")
	}

	// Tagless switch — like an if/else chain but cleaner.
	// Great for classifying values by range.
	score := 85
	switch {
	case score >= 90:
		fmt.Println("  A grade")
	case score >= 80:
		fmt.Println("  B grade")
	case score >= 70:
		fmt.Println("  C grade")
	default:
		fmt.Println("  below C")
	}

	// Switch with init statement — same pattern as if.
	switch lang := strings.ToLower("GO"); lang {
	case "go":
		fmt.Println("  great choice!")
	case "rust":
		fmt.Println("  also great!")
	default:
		fmt.Printf("  %s is fine too\n", lang)
	}
}

/*
=============================================================================
 FOR LOOPS
=============================================================================

Go has exactly one loop keyword: `for`. But it handles every case:

  for i := 0; i < 10; i++ { }     // C-style for loop
  for condition { }                 // while loop
  for { }                           // infinite loop
  for i, v := range collection { } // iterate over collections

The range keyword is how you idiomatically iterate in Go. It works with
strings, arrays, slices, maps, and channels.

IMPORTANT: range over a string iterates by RUNE (Unicode code point),
not by byte. This is different from indexing with s[i] which gives bytes.

=============================================================================
*/

// DemoForLoops shows the different forms of Go's for loop.
func DemoForLoops() {
	// C-style for loop — familiar if you know C, Java, or JavaScript.
	fmt.Println("C-style for:")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %d", i)
	}
	fmt.Println()

	// While-style — just a for with a condition.
	// Go doesn't have a while keyword. for IS while.
	fmt.Println("While-style for:")
	n := 1
	for n < 32 {
		fmt.Printf("  %d", n)
		n *= 2
	}
	fmt.Println()

	// Infinite loop — use break to exit.
	// Common for server main loops, event processing, etc.
	fmt.Println("Infinite loop with break:")
	count := 0
	for {
		if count >= 3 {
			break
		}
		fmt.Printf("  iteration %d", count)
		count++
	}
	fmt.Println()

	// continue — skip the rest of this iteration.
	fmt.Println("Continue (skip evens):")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("  %d", i)
	}
	fmt.Println()
}

/*
=============================================================================
 RANGE
=============================================================================

Range is the idiomatic way to iterate over collections in Go. It returns
up to two values depending on the collection type:

  for i, v := range slice   { }  // index, value (copy!)
  for k, v := range map     { }  // key, value (random order!)
  for i, r := range string  { }  // byte index, rune
  for v := range channel    { }  // value (blocks until closed)

CRITICAL GOTCHA: The value from range is a COPY. Modifying v in the loop
does NOT modify the original collection. If you need to modify elements,
use the index:

  for i := range slice {
      slice[i] = newValue  // this modifies the original
  }

=============================================================================
*/

// DemoRange shows how range works with different collection types.
func DemoRange() {
	// Range over a slice — index and value.
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Println("Range over slice:")
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}

	// Use _ to ignore index or value.
	fmt.Println("Range ignoring index:")
	for _, fruit := range fruits {
		fmt.Printf("  %s\n", fruit)
	}

	// Range over a string — iterates by RUNE, not byte!
	// The index is the byte position, the value is the rune.
	fmt.Println("Range over string (byte index, rune):")
	for i, r := range "Go世界" {
		fmt.Printf("  byte[%d] = '%c' (U+%04X)\n", i, r, r)
	}
	// Output shows: byte[0]='G', byte[1]='o', byte[2]='世', byte[5]='界'
	// Notice the index jumps from 2 to 5 because '世' is 3 bytes in UTF-8.

	// Range over a map — order is NOT guaranteed!
	// Go randomizes map iteration order on purpose to prevent code from
	// depending on a specific order.
	ages := map[string]int{"Alice": 30, "Bob": 25, "Charlie": 35}
	fmt.Println("Range over map (random order):")
	for name, age := range ages {
		fmt.Printf("  %s: %d\n", name, age)
	}
}

/*
=============================================================================
 LABELS, BREAK, AND CONTINUE
=============================================================================

Go supports labeled statements, which let you break or continue an outer
loop from inside an inner loop. This is cleaner than using boolean flags.

  outer:
  for i := 0; i < 10; i++ {
      for j := 0; j < 10; j++ {
          if someCondition {
              break outer  // breaks the OUTER loop
          }
      }
  }

Labels are occasionally useful but shouldn't be overused. If you find
yourself needing complex label logic, consider extracting a function.

Go also has goto, but it's rarely used in practice. It exists mainly for
machine-generated code. You'll almost never write goto in normal Go code.

=============================================================================
*/

// DemoLabels shows how to use labeled break and continue.
func DemoLabels() {
	// Without labels, break only exits the innermost loop.
	// With labels, you can break out of any enclosing loop.
	fmt.Println("Labeled break — finding a value in a 2D grid:")

	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	target := 5

search:
	for row, cols := range grid {
		for col, val := range cols {
			if val == target {
				fmt.Printf("  Found %d at [%d][%d]\n", target, row, col)
				break search // exits BOTH loops
			}
		}
	}

	// Labeled continue — skip to the next iteration of the outer loop.
	fmt.Println("Labeled continue — skip rows with negative values:")
	data := [][]int{
		{1, 2, 3},
		{4, -1, 6},
		{7, 8, 9},
	}

outer:
	for i, row := range data {
		for _, val := range row {
			if val < 0 {
				fmt.Printf("  Skipping row %d (contains negative)\n", i)
				continue outer
			}
		}
		fmt.Printf("  Row %d is all positive: %v\n", i, row)
	}
}

/*
=============================================================================
 PUTTING IT ALL TOGETHER
=============================================================================

Here's a more realistic example that combines several control flow patterns.
This function processes a string of user input, demonstrating how these
patterns work together in practice.

=============================================================================
*/

// AnalyzeText demonstrates combining control flow patterns on a string.
// It categorizes each character and prints a summary.
func AnalyzeText(text string) (letters, digits, spaces, other int) {
	for _, r := range text {
		switch {
		case unicode.IsLetter(r):
			letters++
		case unicode.IsDigit(r):
			digits++
		case unicode.IsSpace(r):
			spaces++
		default:
			other++
		}
	}
	return
}
