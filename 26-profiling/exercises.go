package profiling

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

/*
=============================================================================
 EXERCISES: Profiling and Benchmarks
=============================================================================

 These exercises focus on writing benchmarks, comparing approaches, and
 understanding performance characteristics. Run the tests with:

   go test -v ./26-profiling/

 Run benchmarks with:

   go test -bench=. -benchmem ./26-profiling/

 For specific benchmarks:

   go test -bench=BenchmarkConcat -benchmem ./26-profiling/

 The exercises_test.go file contains BOTH regular tests AND benchmark
 functions. Implement the exercises here, then run benchmarks to compare.

=============================================================================
*/

// Exercise 1: ReverseString returns the input string reversed.
// Implement TWO versions: one using simple string concatenation,
// and one using a byte slice for efficiency.
//
// Example: ReverseStringSimple("hello") returns "olleh"
// Example: ReverseStringFast("hello") returns "olleh"
//
// After implementing, run: go test -bench=BenchmarkReverse -benchmem ./26-profiling/
// to see the performance difference.

// ReverseStringSimple reverses a string using += concatenation.
// This will be slow due to string immutability — each += creates a new string.
func ReverseStringSimple(s string) string {
	// YOUR CODE HERE
	// Loop through the string in reverse, building result with +=
	return ""
}

// ReverseStringFast reverses a string using a byte slice.
// Pre-allocate the byte slice, fill it in reverse, convert once at the end.
func ReverseStringFast(s string) string {
	// YOUR CODE HERE
	// Hint: make([]byte, len(s)), then fill from end to start
	return ""
}

// Exercise 2: FindDuplicates takes a slice of strings and returns a slice
// of strings that appear more than once. Implement two versions:
// one using a map[string]int counter, and one using sorting.
//
// Example: FindDuplicatesMap([]string{"a", "b", "a", "c", "b"}) returns ["a", "b"]
// The order of the result doesn't matter.
//
// After implementing, benchmark both to see which is faster for different sizes.

// FindDuplicatesMap finds duplicates using a map to count occurrences.
func FindDuplicatesMap(items []string) []string {
	// YOUR CODE HERE
	return nil
}

// FindDuplicatesSort finds duplicates by sorting first, then scanning for
// adjacent equal elements.
func FindDuplicatesSort(items []string) []string {
	// YOUR CODE HERE
	_ = sort.Strings
	return nil
}

// Exercise 3: MatrixMultiply multiplies two square matrices.
// Implement two versions: naive (i,j,k ordering) and cache-friendly
// (i,k,j ordering which improves cache locality).
//
// Both functions take two n×n matrices (as [][]int) and return the product.
// Assume matrices are square and non-empty.

// MatrixMultiplyNaive uses standard i,j,k loop ordering.
func MatrixMultiplyNaive(a, b [][]int) [][]int {
	// YOUR CODE HERE
	return nil
}

// MatrixMultiplyOptimized uses i,k,j loop ordering for better cache locality.
// In the i,j,k ordering, accessing b[k][j] causes cache misses because
// we stride down columns. In i,k,j ordering, we access b[k][j] sequentially
// across rows, which is much more cache-friendly.
func MatrixMultiplyOptimized(a, b [][]int) [][]int {
	// YOUR CODE HERE
	return nil
}

// Exercise 4: RegisterPprof registers pprof handlers on a given ServeMux.
//
// In production, you want pprof on a SEPARATE mux from your API, so it's
// not exposed to the internet. This function should register the standard
// pprof handlers on the given mux.
//
// Register these paths:
//   /debug/pprof/           — index page
//   /debug/pprof/cmdline    — command line arguments
//   /debug/pprof/profile    — CPU profile
//   /debug/pprof/symbol     — symbol lookup
//   /debug/pprof/trace      — execution trace
//
// Hint: use net/http/pprof package's handler functions.
// Since we can't import net/http/pprof in this file without side effects,
// just register a basic handler that returns a 200 with text indicating
// where pprof would be. The test just checks the routes exist.
func RegisterPprof(mux *http.ServeMux) {
	// YOUR CODE HERE
	// Register handlers for the 5 pprof endpoints listed above.
	// For this exercise, simple placeholder handlers are fine:
	//   mux.HandleFunc("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
	//       fmt.Fprint(w, "pprof index")
	//   })
	_ = mux
	_ = fmt.Fprint
}

// Exercise 5: OptimizeFunction takes an input and does some computation.
// The "slow" version is provided. Implement the "fast" version that
// produces the same result but more efficiently.
//
// SlowWordCount counts words in a text by splitting into lines, then
// splitting each line into words, creating intermediate slices everywhere.
func SlowWordCount(text string) map[string]int {
	counts := make(map[string]int)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			lower := strings.ToLower(word)
			counts[lower]++
		}
	}
	return counts
}

// FastWordCount should produce the same result as SlowWordCount but
// with fewer allocations. Hint: don't split into lines first —
// use strings.Fields on the entire text (it handles all whitespace).
// Also pre-allocate the map.
func FastWordCount(text string) map[string]int {
	// YOUR CODE HERE
	return nil
}

// Exercise 6: Demonstrate the cost of fmt.Sprintf vs direct conversion.
// FormatIntsSprintf formats a slice of ints as strings using fmt.Sprintf.
// FormatIntsDirect formats a slice of ints as strings using strconv.Itoa.
//
// Implement both, then benchmark to see the reflection overhead of fmt.

// FormatIntsSprintf converts each int to string using fmt.Sprintf("%d", n).
func FormatIntsSprintf(nums []int) []string {
	// YOUR CODE HERE
	_ = fmt.Sprintf
	return nil
}

// FormatIntsDirect converts each int to string using direct conversion
// (you can use the itoa function from lesson.go, or strconv.Itoa).
func FormatIntsDirect(nums []int) []string {
	// YOUR CODE HERE
	return nil
}

// Exercise 7: AllocationHeavy vs AllocationLight.
// ProcessRecords takes a slice of raw record strings (format: "name:value")
// and returns a slice of processed results.
//
// AllocationHeavy creates intermediate structs for each record.
// AllocationLight processes everything inline without intermediate allocations.

// Record is used by the allocation-heavy version.
type Record struct {
	Name  string
	Value string
}

// ProcessRecordsHeavy parses records into intermediate Record structs,
// then formats them. This creates more allocations.
func ProcessRecordsHeavy(records []string) []string {
	// YOUR CODE HERE
	// Step 1: Parse all records into []Record
	// Step 2: Format each Record into "Name=Value" string
	return nil
}

// ProcessRecordsLight parses and formats in a single pass,
// avoiding intermediate Record allocations.
func ProcessRecordsLight(records []string) []string {
	// YOUR CODE HERE
	// Parse and format in one step using strings.Cut
	return nil
}

// Exercise 8: OptimizeHandler simulates optimizing an HTTP-like handler.
// Given a JSON-like input (simplified as key:value lines), parse it,
// process it, and return a response string.
//
// The slow version uses multiple string operations and intermediate slices.
// The fast version minimizes allocations.

// SlowHandler parses input, processes it, and returns a formatted response.
func SlowHandler(input string) string {
	// Parse lines
	lines := strings.Split(input, "\n")
	var pairs []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			pairs = append(pairs, key+"="+val)
		}
	}
	return strings.Join(pairs, "&")
}

// FastHandler should produce the same result as SlowHandler but with
// fewer allocations. Use strings.Builder and strings.Cut.
func FastHandler(input string) string {
	// YOUR CODE HERE
	return ""
}
