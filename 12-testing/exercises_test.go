package testingmod

/*
Exercises: Write the Tests!
============================

This module is inverted — the code is implemented in exercises.go and lesson.go.
YOUR JOB is to fill in the test stubs below.

Each exercise tells you WHAT to test. You write the test logic.

Run your tests with:
  go test -v ./12-testing/
  go test -race ./12-testing/
  go test -bench=. ./12-testing/
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// =============================================================================
// Exercise 1: Table-Driven Tests for TitleCase
// =============================================================================

// TestTitleCase should test the TitleCase function using table-driven tests.
//
// Write a test table with at least these cases:
//   - Simple two-word string
//   - Already title-cased string
//   - ALL CAPS input
//   - all lowercase input
//   - Empty string
//   - Single word
//   - String with extra spaces
//
// Use t.Run for each test case.
// Use descriptive failure messages that include input, got, and want.
func TestTitleCase(t *testing.T) {
	// YOUR CODE HERE
	//
	// Example structure:
	// tests := []struct {
	//     name  string
	//     input string
	//     want  string
	// }{
	//     {"simple", "hello world", "Hello World"},
	//     // Add more cases...
	// }
	// for _, tt := range tests {
	//     t.Run(tt.name, func(t *testing.T) {
	//         got := TitleCase(tt.input)
	//         if got != tt.want {
	//             t.Errorf("TitleCase(%q) = %q, want %q", tt.input, got, tt.want)
	//         }
	//     })
	// }
}

// =============================================================================
// Exercise 2: Edge Case Tests for Truncate
// =============================================================================

// TestTruncate should test the Truncate function with special attention to
// edge cases.
//
// Test cases should include:
//   - String shorter than maxLen (no truncation)
//   - String exactly at maxLen
//   - String longer than maxLen (truncation with "...")
//   - Empty string
//   - maxLen of 0
//   - maxLen of 1, 2, 3 (boundary around the "..." behavior)
//   - Unicode strings (e.g., emoji, CJK characters)
//   - Very long strings
//
// Use subtests to organize by category (e.g., "normal cases", "edge cases",
// "unicode").
func TestTruncate(t *testing.T) {
	// YOUR CODE HERE
}

// =============================================================================
// Exercise 3: Benchmark Two Implementations
// =============================================================================

// BenchmarkCountVowelsLoop should benchmark CountVowelsLoop.
//
// Use b.N for the loop count. Test with a reasonably long string
// to get meaningful results.
//
// Run with: go test -bench=BenchmarkCountVowels -benchmem ./12-testing/
func BenchmarkCountVowelsLoop(b *testing.B) {
	// YOUR CODE HERE
	//
	// Hint:
	// input := "The quick brown fox jumps over the lazy dog"
	// for i := 0; i < b.N; i++ {
	//     CountVowelsLoop(input)
	// }
}

// BenchmarkCountVowelsReplace should benchmark CountVowelsReplace.
//
// Use the SAME input string as BenchmarkCountVowelsLoop so the
// benchmarks are directly comparable.
func BenchmarkCountVowelsReplace(b *testing.B) {
	// YOUR CODE HERE
}

// =============================================================================
// Exercise 4: Test Helper Function
// =============================================================================

// assertDivision is a test helper for the Divide function.
//
// Implement this helper so it:
//  1. Calls t.Helper() (so failure line numbers point to the caller)
//  2. Calls Divide(dividend, divisor)
//  3. Asserts no error was returned
//  4. Asserts quotient matches wantQ
//  5. Asserts remainder matches wantR
//  6. Uses descriptive error messages
func assertDivision(t *testing.T, dividend, divisor, wantQ, wantR int) {
	// YOUR CODE HERE
}

// assertDivisionError is a test helper that asserts Divide returns an error.
//
// Implement this helper so it:
//  1. Calls t.Helper()
//  2. Calls Divide(dividend, divisor)
//  3. Asserts an error WAS returned
func assertDivisionError(t *testing.T, dividend, divisor int) {
	// YOUR CODE HERE
}

// TestDivide should test the Divide function using the helpers above.
//
// Test cases:
//   - Basic division (10 / 3 = 3 remainder 1)
//   - Even division (10 / 5 = 2 remainder 0)
//   - Division by 1
//   - Division of 0
//   - Division by zero (should error)
//   - Negative numbers
func TestDivide(t *testing.T) {
	// YOUR CODE HERE
	//
	// Use assertDivision and assertDivisionError helpers
}

// =============================================================================
// Exercise 5: Test with Fixtures from testdata/
// =============================================================================

// TestParseCSVLineFromFixture should test ParseCSVLine using test data
// loaded from the testdata/ directory.
//
// Steps:
//  1. Read the file testdata/csv_samples.txt
//  2. Each line is a test input for ParseCSVLine
//  3. Parse each line and verify the result has the expected number of fields
//
// The testdata/csv_samples.txt file has been provided with sample CSV lines.
//
// Hint: Use os.ReadFile to load the fixture.
func TestParseCSVLineFromFixture(t *testing.T) {
	// YOUR CODE HERE
	//
	// Hint:
	// data, err := os.ReadFile(filepath.Join("testdata", "csv_samples.txt"))
	// if err != nil {
	//     t.Fatalf("failed to read fixture: %v", err)
	// }
	// Then split by newlines and test each line.

	// Suppress unused import warnings — remove these when you implement
	_ = os.ReadFile
	_ = filepath.Join
}

// =============================================================================
// Exercise 6: Mock a Dependency
// =============================================================================

// MockUserStore is a mock implementation of the UserStore interface.
//
// Implement this mock so that:
//   - It stores users in memory (a map or slice)
//   - GetUser returns the user if found, or an error
//   - ListUsers returns all stored users
//   - CreateUser creates a new user with an auto-incremented ID
//
// Tip: You can pre-populate it with test data in each test.
type MockUserStore struct {
	// YOUR CODE HERE — define fields
	// Hint: users map[int]*UserRecord, nextID int, etc.
}

// Implement the UserStore interface on MockUserStore:

func (m *MockUserStore) GetUser(id int) (*UserRecord, error) {
	// YOUR CODE HERE
	return nil, fmt.Errorf("not implemented")
}

func (m *MockUserStore) ListUsers() ([]*UserRecord, error) {
	// YOUR CODE HERE
	return nil, fmt.Errorf("not implemented")
}

func (m *MockUserStore) CreateUser(name, email string) (*UserRecord, error) {
	// YOUR CODE HERE
	return nil, fmt.Errorf("not implemented")
}

// TestUserServiceWithMock should test UserService using MockUserStore.
//
// Test:
//   - GetUserDisplayName with a valid user
//   - GetUserDisplayName with a non-existent user (should error)
//   - ListUserNames with multiple users
//   - CreateAndGreet with valid input
func TestUserServiceWithMock(t *testing.T) {
	// YOUR CODE HERE
}

// =============================================================================
// Exercise 7: Golden File Test
// =============================================================================

// TestGenerateInventoryReportGolden should test GenerateInventoryReport
// using the golden file pattern.
//
// Steps:
//  1. Generate a report with known input
//  2. Compare the output against testdata/inventory_report.golden
//  3. If the -update flag is set, write the actual output as the new golden file
//
// For simplicity, you can skip the -update flag and just compare against
// the golden file. The golden file has been pre-created for you.
//
// Hint:
//
//	expected, err := os.ReadFile(filepath.Join("testdata", "inventory_report.golden"))
//	actual := GenerateInventoryReport(...)
//	if actual != string(expected) { t.Errorf(...) }
func TestGenerateInventoryReportGolden(t *testing.T) {
	// YOUR CODE HERE
}

// =============================================================================
// Exercise 8: Parallel Tests
// =============================================================================

// TestWordFrequencyParallel should test WordFrequency using t.Parallel()
// to run test cases concurrently.
//
// Requirements:
//   - Use table-driven tests
//   - Call t.Parallel() in each subtest
//   - Remember to capture the loop variable (tt := tt)
//   - Test at least 4 different inputs
//
// Gotcha to be aware of: with t.Parallel(), all subtests start at roughly
// the same time. If they share any mutable state, you'll get races.
// WordFrequency is a pure function, so this should be safe.
func TestWordFrequencyParallel(t *testing.T) {
	// YOUR CODE HERE
	//
	// Example structure:
	// tests := []struct {
	//     name  string
	//     input string
	//     want  map[string]int
	// }{
	//     {"simple", "hello hello world", map[string]int{"hello": 2, "world": 1}},
	//     // Add more cases...
	// }
	// for _, tt := range tests {
	//     tt := tt  // IMPORTANT: capture loop variable for parallel tests
	//     t.Run(tt.name, func(t *testing.T) {
	//         t.Parallel()
	//         got := WordFrequency(tt.input)
	//         // Compare got with tt.want...
	//     })
	// }
}

// =============================================================================
// Verify lesson code compiles and works
// =============================================================================

func TestLessonFunctions(t *testing.T) {
	t.Run("Reverse", func(t *testing.T) {
		if got := Reverse("hello"); got != "olleh" {
			t.Errorf("Reverse(hello) = %q, want %q", got, "olleh")
		}
	})

	t.Run("IsPalindrome", func(t *testing.T) {
		if !IsPalindrome("racecar") {
			t.Error("racecar should be a palindrome")
		}
		if IsPalindrome("hello") {
			t.Error("hello should not be a palindrome")
		}
	})

	t.Run("Abs", func(t *testing.T) {
		if got := Abs(-5); got != 5 {
			t.Errorf("Abs(-5) = %d, want 5", got)
		}
	})

	t.Run("Clamp", func(t *testing.T) {
		if got := Clamp(15, 0, 10); got != 10 {
			t.Errorf("Clamp(15, 0, 10) = %d, want 10", got)
		}
	})

	t.Run("FormatReport", func(t *testing.T) {
		report := FormatReport("Test", []string{"a", "b"})
		if report == "" {
			t.Error("FormatReport returned empty string")
		}
	})

	t.Run("SlowHash", func(t *testing.T) {
		h1 := SlowHash("test")
		h2 := SlowHash("test")
		if h1 != h2 {
			t.Error("SlowHash should be deterministic")
		}
	})
}
