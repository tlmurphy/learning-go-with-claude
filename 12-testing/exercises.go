package testingmod

/*
=============================================================================
 EXERCISES: Testing
=============================================================================

 Work through these exercises in order. Run tests with:

   make test 12

 Run a single test:

   go test -v -run TestTitleCase ./12-testing/

=============================================================================
*/

import (
	"fmt"
	"strings"
	"unicode"
)

// =============================================================================
// Exercise 1: Function to Test with Table-Driven Tests
// =============================================================================

// TitleCase converts a string to title case. Each word's first letter is
// capitalized, and the rest are lowercased.
//
// Words are separated by spaces. Leading/trailing spaces are preserved.
// Multiple consecutive spaces are preserved.
//
// Examples:
//
//	TitleCase("hello world")     → "Hello World"
//	TitleCase("gO is AWESOME")  → "Go Is Awesome"
//	TitleCase("")               → ""
//	TitleCase("  hello  ")      → "  Hello  "
func TitleCase(s string) string {
	if s == "" {
		return ""
	}

	words := strings.Split(s, " ")
	for i, word := range words {
		if word == "" {
			continue
		}
		runes := []rune(strings.ToLower(word))
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// =============================================================================
// Exercise 2: Functions to Test Edge Cases
// =============================================================================

// Truncate shortens a string to maxLen characters, adding "..." if truncated.
// If the string is shorter than or equal to maxLen, it's returned unchanged.
// If maxLen is less than 3, the result is just the first maxLen characters
// (no room for "...").
//
// Examples:
//
//	Truncate("hello world", 5)  → "he..."
//	Truncate("hi", 5)           → "hi"
//	Truncate("hello", 2)        → "he"
//	Truncate("", 5)             → ""
//	Truncate("abc", 3)          → "abc"
func Truncate(s string, maxLen int) string {
	if maxLen < 0 {
		maxLen = 0
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// =============================================================================
// Exercise 3: Two Implementations for Benchmarking
// =============================================================================

// CountVowelsLoop counts vowels (a, e, i, o, u) in a string using a loop.
func CountVowelsLoop(s string) int {
	count := 0
	for _, ch := range strings.ToLower(s) {
		switch ch {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// CountVowelsReplace counts vowels by removing all non-vowels and
// measuring what's left. This is an alternative approach.
func CountVowelsReplace(s string) int {
	lower := strings.ToLower(s)
	vowelsOnly := strings.Map(func(r rune) rune {
		switch r {
		case 'a', 'e', 'i', 'o', 'u':
			return r
		}
		return -1 // -1 means drop this rune
	}, lower)
	return len([]rune(vowelsOnly))
}

// =============================================================================
// Exercise 4: Helper Function Target
// =============================================================================

// Divide performs integer division and returns (quotient, remainder, error).
// Returns an error if divisor is zero.
func Divide(dividend, divisor int) (int, int, error) {
	if divisor == 0 {
		return 0, 0, fmt.Errorf("division by zero")
	}
	return dividend / divisor, dividend % divisor, nil
}

// =============================================================================
// Exercise 5: Function for Fixture-Based Testing
// =============================================================================

// ParseCSVLine parses a single CSV line into fields.
// Handles quoted fields (double-quotes) and escaped quotes ("").
// This is a simplified parser — not fully RFC 4180 compliant.
//
// Examples:
//
//	ParseCSVLine("a,b,c")           → ["a", "b", "c"]
//	ParseCSVLine(`"hello","world"`) → ["hello", "world"]
//	ParseCSVLine("")                → [""]
func ParseCSVLine(line string) []string {
	var fields []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		switch {
		case ch == '"' && !inQuotes:
			inQuotes = true
		case ch == '"' && inQuotes:
			if i+1 < len(line) && line[i+1] == '"' {
				current.WriteByte('"')
				i++ // Skip escaped quote
			} else {
				inQuotes = false
			}
		case ch == ',' && !inQuotes:
			fields = append(fields, current.String())
			current.Reset()
		default:
			current.WriteByte(ch)
		}
	}
	fields = append(fields, current.String())
	return fields
}

// =============================================================================
// Exercise 6: The UserService (for Mocking) — see lesson.go
// =============================================================================
// The UserService, UserStore interface, and UserRecord are defined in lesson.go.
// Write tests that mock the UserStore interface.

// =============================================================================
// Exercise 7: Function for Golden File Testing
// =============================================================================

// GenerateInventoryReport creates a formatted inventory report.
// This output format should be tested using golden files.
func GenerateInventoryReport(storeName string, items map[string]int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Inventory Report: %s\n", storeName))
	b.WriteString(strings.Repeat("-", 40) + "\n")

	total := 0
	// Note: map iteration is non-deterministic, so we need sorted keys
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	// Sort for deterministic output
	sortStrings(keys)

	for _, name := range keys {
		qty := items[name]
		b.WriteString(fmt.Sprintf("  %-20s %5d\n", name, qty))
		total += qty
	}
	b.WriteString(strings.Repeat("-", 40) + "\n")
	b.WriteString(fmt.Sprintf("  %-20s %5d\n", "TOTAL", total))
	return b.String()
}

// sortStrings is a simple insertion sort to avoid importing "sort" just for this.
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		key := ss[i]
		j := i - 1
		for j >= 0 && ss[j] > key {
			ss[j+1] = ss[j]
			j--
		}
		ss[j+1] = key
	}
}

// =============================================================================
// Exercise 8: Function for Parallel Tests
// =============================================================================

// WordFrequency counts the frequency of each word in the input text.
// Words are converted to lowercase and split on whitespace.
// Punctuation is NOT stripped (it's part of the word).
func WordFrequency(text string) map[string]int {
	freq := make(map[string]int)
	words := strings.Fields(strings.ToLower(text))
	for _, w := range words {
		freq[w]++
	}
	return freq
}
