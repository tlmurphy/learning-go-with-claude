// Package collections covers Go's built-in collection types: arrays, slices,
// and maps — the data structures you'll use in almost every Go program.
package collections

import "fmt"

/*
=============================================================================
 COLLECTIONS IN GO
=============================================================================

Go has three built-in collection types:

  Arrays  — fixed size, value type (rarely used directly)
  Slices  — dynamic, reference type (the workhorse)
  Maps    — hash table (key-value pairs)

You'll use slices and maps constantly. Arrays exist mainly as the backing
store for slices — you'll rarely declare an array directly.

Understanding slices deeply is crucial. They're one of the most elegant
and most gotcha-prone features in Go. The slice header, backing array
sharing, append behavior, and capacity growth are all things you need to
internalize to write correct Go code.

=============================================================================
 ARRAYS
=============================================================================

Arrays in Go are fixed-size and are VALUE TYPES. This means:
  - The size is part of the type: [3]int and [4]int are different types!
  - Assigning an array copies all elements
  - Passing an array to a function copies all elements

Because of these properties, arrays are rarely used directly. Slices are
almost always preferred. But understanding arrays helps you understand
slices, since slices are backed by arrays.

=============================================================================
*/

// DemoArrays shows how arrays work in Go (and why you rarely use them).
func DemoArrays() {
	// Array declaration — size is part of the type.
	var a [3]int // [0, 0, 0] — zero values
	fmt.Println("Zero array:", a)

	// Array literal.
	b := [3]string{"Go", "is", "great"}
	fmt.Println("Literal array:", b)

	// Size inferred from elements.
	c := [...]int{10, 20, 30, 40, 50}
	fmt.Println("Inferred size:", c, "length:", len(c))

	// Arrays are VALUE TYPES — assignment copies!
	original := [3]int{1, 2, 3}
	copy := original // full copy
	copy[0] = 999    // doesn't affect original!
	fmt.Println("Original:", original, "Copy:", copy)

	// [3]int and [4]int are DIFFERENT types — you can't assign between them.
	// var d [4]int = a  // COMPILE ERROR: cannot use a (type [3]int) as [4]int
}

/*
=============================================================================
 SLICES
=============================================================================

Slices are Go's answer to dynamic arrays. They're by far the most commonly
used collection type. A slice is a lightweight descriptor (header) that
points to an underlying array:

  Slice header:
  ┌─────────┬─────┬──────────┐
  │ pointer │ len │ capacity │
  └─────────┴─────┴──────────┘
       │
       ▼
  ┌───┬───┬───┬───┬───┬───┬───┐
  │ 1 │ 2 │ 3 │ 4 │ . │ . │ . │  underlying array
  └───┴───┴───┴───┴───┴───┴───┘
  ←─── len ───→
  ←──────── cap ────────────→

Key concepts:
  - len(s) = number of elements currently in the slice
  - cap(s) = number of elements in the underlying array (from slice start)
  - append() may or may not allocate a new array
  - Multiple slices can share the same underlying array (aliasing!)

=============================================================================
*/

// DemoSlices shows slice creation, manipulation, and internals.
func DemoSlices() {
	// Slice literal (no size in brackets — that's what makes it a slice, not an array).
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("Slice: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// make() creates a slice with specified length and optional capacity.
	// This is how you pre-allocate when you know the size.
	m := make([]int, 3, 10) // length 3, capacity 10
	fmt.Printf("make([]int, 3, 10): %v, len=%d, cap=%d\n", m, len(m), cap(m))

	// Slicing creates a new slice header pointing to the SAME underlying array.
	sub := s[1:3] // elements at index 1, 2 (not 3!)
	fmt.Printf("s[1:3] = %v (shares underlying array with s)\n", sub)

	// GOTCHA: modifying sub also modifies s!
	sub[0] = 999
	fmt.Printf("After sub[0]=999: s=%v, sub=%v\n", s, sub)

	// Nil slice vs empty slice — both are valid, both have length 0.
	var nilSlice []int          // nil
	emptySlice := []int{}       // not nil, but empty
	makeSlice := make([]int, 0) // not nil, but empty
	fmt.Printf("nil slice: %v (nil? %t, len=%d)\n", nilSlice, nilSlice == nil, len(nilSlice))
	fmt.Printf("empty slice: %v (nil? %t, len=%d)\n", emptySlice, emptySlice == nil, len(emptySlice))
	fmt.Printf("make slice: %v (nil? %t, len=%d)\n", makeSlice, makeSlice == nil, len(makeSlice))

	// append works on nil slices — no need to initialize first!
	nilSlice = append(nilSlice, 1, 2, 3)
	fmt.Printf("nil slice after append: %v\n", nilSlice)
}

/*
=============================================================================
 APPEND AND CAPACITY GROWTH
=============================================================================

append() is how you add elements to a slice. It has important behavior:

  - If len < cap: append uses the existing underlying array (fast)
  - If len == cap: append allocates a new, larger array and copies (slow)

The growth strategy is roughly:
  - Double the capacity when small (< 256)
  - Grow by ~25% when larger

This means append() sometimes returns a slice with a different underlying
array. ALWAYS assign the result of append back to the variable:

  s = append(s, value)    // CORRECT
  append(s, value)         // BUG: result is discarded!

Performance tip: If you know the final size, use make() with a capacity
to avoid repeated allocations:

  result := make([]int, 0, expectedSize)
  for ... {
      result = append(result, value)
  }

=============================================================================
*/

// DemoAppendGrowth shows how append grows slice capacity.
func DemoAppendGrowth() {
	var s []int
	prevCap := cap(s)

	fmt.Println("Watching capacity growth:")
	for i := 0; i < 20; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("  len=%2d, cap changed: %d -> %d\n", len(s), prevCap, cap(s))
			prevCap = cap(s)
		}
	}

	// Pre-allocation avoids all the intermediate allocations.
	prealloc := make([]int, 0, 20) // capacity for 20 elements
	for i := 0; i < 20; i++ {
		prealloc = append(prealloc, i)
	}
	fmt.Printf("Pre-allocated: len=%d, cap=%d (no growth needed!)\n", len(prealloc), cap(prealloc))
}

/*
=============================================================================
 SLICE GOTCHAS
=============================================================================

Slices have several gotchas that catch even experienced Go developers:

1. SHARED BACKING ARRAY
   When you slice a slice, both share the same memory. Modifying one
   can affect the other.

2. STALE SLICES
   After append grows the array, old slices still point to the old array.
   They don't see new elements.

3. MEMORY LEAKS
   Slicing a huge slice to keep only a small part still retains the
   entire backing array in memory. Use copy() to avoid this.

4. RANGE VALUE IS A COPY
   for _, v := range slice — v is a COPY. Modifying v doesn't modify
   the slice. Use the index: slice[i] = newValue.

=============================================================================
*/

// DemoSliceGotchas shows common slice pitfalls.
func DemoSliceGotchas() {
	// Gotcha 1: Shared backing array.
	original := []int{1, 2, 3, 4, 5}
	slice1 := original[1:3] // [2, 3]
	slice1[0] = 999         // also modifies original!
	fmt.Println("Shared array gotcha:")
	fmt.Printf("  original: %v (modified!)\n", original)

	// Fix: use copy() to create an independent slice.
	original2 := []int{1, 2, 3, 4, 5}
	independent := make([]int, 2)
	copy(independent, original2[1:3])
	independent[0] = 999 // does NOT modify original2
	fmt.Printf("  With copy: original2=%v, independent=%v\n", original2, independent)

	// Gotcha 2: Range value is a copy.
	nums := []int{1, 2, 3}
	for _, v := range nums {
		v *= 2 // modifies the COPY, not the slice!
		_ = v
	}
	fmt.Printf("Range copy gotcha: %v (unchanged!)\n", nums)

	// Fix: use the index.
	for i := range nums {
		nums[i] *= 2
	}
	fmt.Printf("Using index: %v (modified correctly)\n", nums)
}

/*
=============================================================================
 MAPS
=============================================================================

Maps are Go's built-in hash tables. They provide O(1) average-case lookups,
insertions, and deletions.

Key rules:
  - Map keys must be comparable (==, !=). No slices, maps, or functions as keys.
  - The zero value of a map is nil. Reading from nil is OK (returns zero value).
    Writing to nil panics!
  - Always use make() or a literal to initialize before writing.
  - Iteration order is NOT guaranteed (Go randomizes it deliberately).
  - Maps are NOT safe for concurrent access (use sync.Mutex or sync.Map).

The comma-ok idiom is essential for maps:

  value, ok := myMap[key]
  if !ok {
      // key doesn't exist
  }

Without the comma-ok, you can't distinguish "key exists with zero value"
from "key doesn't exist."

=============================================================================
*/

// DemoMaps shows map creation, access, and the comma-ok pattern.
func DemoMaps() {
	// Map literal.
	ages := map[string]int{
		"Alice":   30,
		"Bob":     25,
		"Charlie": 35,
	}
	fmt.Println("Map literal:", ages)

	// make() for empty map with optional size hint.
	scores := make(map[string]int, 10) // hint: expect ~10 entries
	scores["Go"] = 95
	scores["Rust"] = 90
	fmt.Println("make map:", scores)

	// Access: returns zero value if key doesn't exist.
	fmt.Println("Alice's age:", ages["Alice"])   // 30
	fmt.Println("Unknown age:", ages["Unknown"]) // 0 (zero value)

	// Comma-ok pattern: distinguish "not found" from "zero value."
	age, ok := ages["Alice"]
	fmt.Printf("Alice: age=%d, exists=%t\n", age, ok)

	age, ok = ages["Unknown"]
	fmt.Printf("Unknown: age=%d, exists=%t\n", age, ok)

	// Delete a key.
	delete(ages, "Bob")
	fmt.Println("After delete:", ages)

	// Nil map gotcha:
	// var m map[string]int
	// m["key"] = 1  // PANIC! Can't write to nil map.
	// _ = m["key"]  // OK — returns 0. Reading nil map is safe.

	// Iteration order is random!
	fmt.Println("Iteration (random order):")
	for k, v := range ages {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

/*
=============================================================================
 COMMON MAP PATTERNS
=============================================================================

Maps are used for many things beyond simple key-value lookup:

  Set:          map[string]struct{}     (struct{} uses zero bytes!)
  Counter:      map[string]int
  Group by:     map[string][]Item
  Lookup table: map[int]string
  Cache:        map[Key]Value

The set pattern uses struct{} as the value type because it takes zero
bytes of memory. It's the idiomatic way to represent a set in Go:

  seen := map[string]struct{}{}
  seen["item"] = struct{}{}
  if _, exists := seen["item"]; exists { ... }

=============================================================================
*/

// DemoMapPatterns shows common map patterns.
func DemoMapPatterns() {
	// Pattern: Set using map[string]struct{}
	seen := make(map[string]struct{})
	words := []string{"hello", "world", "hello", "go", "world"}

	for _, w := range words {
		seen[w] = struct{}{}
	}
	fmt.Println("Unique words (set):")
	for word := range seen {
		fmt.Printf("  %s\n", word)
	}

	// Pattern: Counter
	counter := make(map[string]int)
	for _, w := range words {
		counter[w]++ // works even if key doesn't exist (zero value + 1)
	}
	fmt.Println("Word counts:")
	for word, count := range counter {
		fmt.Printf("  %s: %d\n", word, count)
	}

	// Pattern: Group by
	type Person struct {
		Name string
		City string
	}
	people := []Person{
		{"Alice", "NYC"}, {"Bob", "LA"}, {"Charlie", "NYC"}, {"Diana", "LA"},
	}
	byCity := make(map[string][]string)
	for _, p := range people {
		byCity[p.City] = append(byCity[p.City], p.Name)
	}
	fmt.Println("Group by city:")
	for city, names := range byCity {
		fmt.Printf("  %s: %v\n", city, names)
	}
}

/*
=============================================================================
 COPY AND SLICE MANIPULATION
=============================================================================

The built-in copy() function copies elements between slices:

  copy(dst, src)  // copies min(len(dst), len(src)) elements

It returns the number of elements copied. The destination must have
enough length (not just capacity!) to receive the elements.

Common slice operations that don't have built-in functions:
  - Insert at position
  - Remove at position
  - Prepend
  - Deduplicate

=============================================================================
*/

// DemoCopy shows how to use copy() and common slice operations.
func DemoCopy() {
	// Basic copy.
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3) // only length 3
	n := copy(dst, src)
	fmt.Printf("Copied %d elements: %v\n", n, dst)

	// Remove element at index i (order-preserving).
	s := []int{10, 20, 30, 40, 50}
	i := 2 // remove 30
	s = append(s[:i], s[i+1:]...)
	fmt.Printf("After removing index 2: %v\n", s)

	// Insert element at index i.
	s2 := []int{10, 20, 40, 50}
	i = 2 // insert 30 at index 2
	s2 = append(s2[:i], append([]int{30}, s2[i:]...)...)
	fmt.Printf("After inserting 30 at index 2: %v\n", s2)
}
