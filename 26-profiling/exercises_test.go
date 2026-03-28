package profiling

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Exercise 1: ReverseString
// ---------------------------------------------------------------------------

func TestReverseStringSimple(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"racecar", "racecar"},
		{"Go is fun", "nuf si oG"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ReverseStringSimple(tt.input)
			if got != tt.want {
				t.Errorf("ReverseStringSimple(%q) = %q, want %q. Build the result by iterating in reverse and concatenating with +=.", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverseStringFast(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"racecar", "racecar"},
		{"Go is fun", "nuf si oG"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ReverseStringFast(tt.input)
			if got != tt.want {
				t.Errorf("ReverseStringFast(%q) = %q, want %q. Use make([]byte, len(s)) and fill in reverse order.", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkReverseSimple(b *testing.B) {
	s := strings.Repeat("abcdefghij", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ReverseStringSimple(s)
	}
}

func BenchmarkReverseFast(b *testing.B) {
	s := strings.Repeat("abcdefghij", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ReverseStringFast(s)
	}
}

// ---------------------------------------------------------------------------
// Exercise 2: FindDuplicates
// ---------------------------------------------------------------------------

func TestFindDuplicatesMap(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"basic", []string{"a", "b", "a", "c", "b"}, []string{"a", "b"}},
		{"no dups", []string{"a", "b", "c"}, []string{}},
		{"all dups", []string{"a", "a", "a"}, []string{"a"}},
		{"empty", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDuplicatesMap(tt.input)
			if got == nil {
				got = []string{}
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if len(got) != len(tt.want) {
				t.Fatalf("FindDuplicatesMap(%v) returned %d items, want %d. Use a map[string]int to count occurrences.", tt.input, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("FindDuplicatesMap result[%d] = %q, want %q.", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFindDuplicatesSort(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"basic", []string{"a", "b", "a", "c", "b"}, []string{"a", "b"}},
		{"no dups", []string{"a", "b", "c"}, []string{}},
		{"all dups", []string{"a", "a", "a"}, []string{"a"}},
		{"empty", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDuplicatesSort(tt.input)
			if got == nil {
				got = []string{}
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if len(got) != len(tt.want) {
				t.Fatalf("FindDuplicatesSort(%v) returned %d items, want %d. Sort the slice, then scan for adjacent duplicates.", tt.input, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("FindDuplicatesSort result[%d] = %q, want %q.", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func BenchmarkFindDuplicatesMap(b *testing.B) {
	items := makeDupStrings(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FindDuplicatesMap(items)
	}
}

func BenchmarkFindDuplicatesSort(b *testing.B) {
	items := makeDupStrings(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FindDuplicatesSort(items)
	}
}

func makeDupStrings(n int) []string {
	items := make([]string, n)
	for i := range items {
		items[i] = fmt.Sprintf("item-%d", i%50) // 50 unique values → lots of dups
	}
	return items
}

// ---------------------------------------------------------------------------
// Exercise 3: MatrixMultiply
// ---------------------------------------------------------------------------

func TestMatrixMultiplyNaive(t *testing.T) {
	a := [][]int{
		{1, 2},
		{3, 4},
	}
	b := [][]int{
		{5, 6},
		{7, 8},
	}
	want := [][]int{
		{19, 22},
		{43, 50},
	}

	got := MatrixMultiplyNaive(a, b)
	if got == nil {
		t.Fatal("MatrixMultiplyNaive returned nil. Implement standard i,j,k matrix multiplication.")
	}
	assertMatrixEqual(t, "MatrixMultiplyNaive", got, want)
}

func TestMatrixMultiplyOptimized(t *testing.T) {
	a := [][]int{
		{1, 2},
		{3, 4},
	}
	b := [][]int{
		{5, 6},
		{7, 8},
	}
	want := [][]int{
		{19, 22},
		{43, 50},
	}

	got := MatrixMultiplyOptimized(a, b)
	if got == nil {
		t.Fatal("MatrixMultiplyOptimized returned nil. Implement i,k,j matrix multiplication for cache locality.")
	}
	assertMatrixEqual(t, "MatrixMultiplyOptimized", got, want)
}

func TestMatrixMultiply3x3(t *testing.T) {
	a := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	b := [][]int{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}
	want := [][]int{
		{30, 24, 18},
		{84, 69, 54},
		{138, 114, 90},
	}

	t.Run("naive", func(t *testing.T) {
		got := MatrixMultiplyNaive(a, b)
		if got == nil {
			t.Skip("MatrixMultiplyNaive not implemented yet")
		}
		assertMatrixEqual(t, "MatrixMultiplyNaive 3x3", got, want)
	})

	t.Run("optimized", func(t *testing.T) {
		got := MatrixMultiplyOptimized(a, b)
		if got == nil {
			t.Skip("MatrixMultiplyOptimized not implemented yet")
		}
		assertMatrixEqual(t, "MatrixMultiplyOptimized 3x3", got, want)
	})
}

func assertMatrixEqual(t *testing.T, name string, got, want [][]int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: got %d rows, want %d.", name, len(got), len(want))
	}
	for i := range want {
		if len(got[i]) != len(want[i]) {
			t.Fatalf("%s: row %d has %d cols, want %d.", name, i, len(got[i]), len(want[i]))
		}
		for j := range want[i] {
			if got[i][j] != want[i][j] {
				t.Errorf("%s: [%d][%d] = %d, want %d.", name, i, j, got[i][j], want[i][j])
			}
		}
	}
}

func BenchmarkMatrixMultiply(b *testing.B) {
	sizes := []int{10, 50, 100}
	for _, size := range sizes {
		a := makeMatrix(size)
		bm := makeMatrix(size)

		b.Run(fmt.Sprintf("naive/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				MatrixMultiplyNaive(a, bm)
			}
		})
		b.Run(fmt.Sprintf("optimized/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				MatrixMultiplyOptimized(a, bm)
			}
		})
	}
}

func makeMatrix(n int) [][]int {
	m := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
		for j := range m[i] {
			m[i][j] = rand.Intn(100)
		}
	}
	return m
}

// ---------------------------------------------------------------------------
// Exercise 4: RegisterPprof
// ---------------------------------------------------------------------------

func TestRegisterPprof(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPprof(mux)

	paths := []string{
		"/debug/pprof/",
		"/debug/pprof/cmdline",
		"/debug/pprof/profile",
		"/debug/pprof/symbol",
		"/debug/pprof/trace",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code == 404 {
				t.Errorf("Handler not registered for %s. Register a handler with mux.HandleFunc(%q, ...).", path, path)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Exercise 5: FastWordCount
// ---------------------------------------------------------------------------

func TestFastWordCount(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]int
	}{
		{
			name: "basic",
			text: "hello world hello",
			want: map[string]int{"hello": 2, "world": 1},
		},
		{
			name: "multiline",
			text: "Go is great\ngo is fast\nGO is fun",
			want: map[string]int{"go": 3, "is": 3, "great": 1, "fast": 1, "fun": 1},
		},
		{
			name: "empty",
			text: "",
			want: map[string]int{},
		},
		{
			name: "single word",
			text: "hello",
			want: map[string]int{"hello": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FastWordCount(tt.text)
			if got == nil {
				t.Fatal("FastWordCount returned nil. Use strings.Fields and strings.ToLower to count words.")
			}

			// Verify against SlowWordCount for correctness
			slow := SlowWordCount(tt.text)
			for k, v := range slow {
				if got[k] != v {
					t.Errorf("FastWordCount[%q] = %d, want %d (from SlowWordCount reference).", k, got[k], v)
				}
			}
			for k, v := range got {
				if slow[k] != v {
					t.Errorf("FastWordCount has extra key %q = %d not in reference.", k, v)
				}
			}
		})
	}
}

func BenchmarkSlowWordCount(b *testing.B) {
	text := strings.Repeat("the quick brown fox jumps over the lazy dog\n", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SlowWordCount(text)
	}
}

func BenchmarkFastWordCount(b *testing.B) {
	text := strings.Repeat("the quick brown fox jumps over the lazy dog\n", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FastWordCount(text)
	}
}

// ---------------------------------------------------------------------------
// Exercise 6: FormatInts — Sprintf vs direct conversion
// ---------------------------------------------------------------------------

func TestFormatIntsSprintf(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []string
	}{
		{"basic", []int{1, 2, 3}, []string{"1", "2", "3"}},
		{"negative", []int{-1, 0, 1}, []string{"-1", "0", "1"}},
		{"empty", []int{}, []string{}},
		{"large", []int{1000000}, []string{"1000000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatIntsSprintf(tt.nums)
			if got == nil && len(tt.want) > 0 {
				t.Fatal("FormatIntsSprintf returned nil. Use fmt.Sprintf to convert each int.")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("FormatIntsSprintf returned %d items, want %d.", len(got), len(tt.want))
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("FormatIntsSprintf[%d] = %q, want %q.", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestFormatIntsDirect(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []string
	}{
		{"basic", []int{1, 2, 3}, []string{"1", "2", "3"}},
		{"negative", []int{-1, 0, 1}, []string{"-1", "0", "1"}},
		{"empty", []int{}, []string{}},
		{"large", []int{1000000}, []string{"1000000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatIntsDirect(tt.nums)
			if got == nil && len(tt.want) > 0 {
				t.Fatal("FormatIntsDirect returned nil. Use strconv.Itoa or the itoa helper to convert each int.")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("FormatIntsDirect returned %d items, want %d.", len(got), len(tt.want))
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("FormatIntsDirect[%d] = %q, want %q.", i, v, tt.want[i])
				}
			}
		})
	}
}

func BenchmarkFormatIntsSprintf(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i * 7
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FormatIntsSprintf(nums)
	}
}

func BenchmarkFormatIntsDirect(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i * 7
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FormatIntsDirect(nums)
	}
}

// ---------------------------------------------------------------------------
// Exercise 7: ProcessRecords — allocation heavy vs light
// ---------------------------------------------------------------------------

func TestProcessRecordsHeavy(t *testing.T) {
	tests := []struct {
		name    string
		records []string
		want    []string
	}{
		{
			name:    "basic",
			records: []string{"name:Alice", "age:30"},
			want:    []string{"name=Alice", "age=30"},
		},
		{
			name:    "empty",
			records: []string{},
			want:    []string{},
		},
		{
			name:    "single",
			records: []string{"key:value"},
			want:    []string{"key=value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessRecordsHeavy(tt.records)
			if got == nil && len(tt.want) > 0 {
				t.Fatal("ProcessRecordsHeavy returned nil. Parse each 'name:value' into a Record, then format as 'name=value'.")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Got %d results, want %d.", len(got), len(tt.want))
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("result[%d] = %q, want %q.", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestProcessRecordsLight(t *testing.T) {
	tests := []struct {
		name    string
		records []string
		want    []string
	}{
		{
			name:    "basic",
			records: []string{"name:Alice", "age:30"},
			want:    []string{"name=Alice", "age=30"},
		},
		{
			name:    "empty",
			records: []string{},
			want:    []string{},
		},
		{
			name:    "single",
			records: []string{"key:value"},
			want:    []string{"key=value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessRecordsLight(tt.records)
			if got == nil && len(tt.want) > 0 {
				t.Fatal("ProcessRecordsLight returned nil. Use strings.Cut to parse, then format inline.")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Got %d results, want %d.", len(got), len(tt.want))
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("result[%d] = %q, want %q.", i, v, tt.want[i])
				}
			}
		})
	}
}

func BenchmarkProcessRecordsHeavy(b *testing.B) {
	records := makeRecords(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ProcessRecordsHeavy(records)
	}
}

func BenchmarkProcessRecordsLight(b *testing.B) {
	records := makeRecords(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ProcessRecordsLight(records)
	}
}

func makeRecords(n int) []string {
	records := make([]string, n)
	for i := range records {
		records[i] = fmt.Sprintf("key%d:value%d", i, i)
	}
	return records
}

// ---------------------------------------------------------------------------
// Exercise 8: FastHandler
// ---------------------------------------------------------------------------

func TestFastHandler(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic",
			input: "name: Alice\nage: 30",
			want:  "name=Alice&age=30",
		},
		{
			name:  "empty lines",
			input: "name: Alice\n\nage: 30\n",
			want:  "name=Alice&age=30",
		},
		{
			name:  "single pair",
			input: "key: value",
			want:  "key=value",
		},
		{
			name:  "empty",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FastHandler(tt.input)
			slow := SlowHandler(tt.input)
			if got != slow {
				t.Errorf("FastHandler result differs from SlowHandler.\n  FastHandler: %q\n  SlowHandler: %q\n  Use strings.Builder and strings.Cut for an efficient implementation.",
					got, slow)
			}
			if got != tt.want {
				t.Errorf("FastHandler(%q) = %q, want %q.", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkSlowHandler(b *testing.B) {
	input := strings.Repeat("key: value\nname: test\ncount: 42\n", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SlowHandler(input)
	}
}

func BenchmarkFastHandler(b *testing.B) {
	input := strings.Repeat("key: value\nname: test\ncount: 42\n", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FastHandler(input)
	}
}

// ---------------------------------------------------------------------------
// Bonus: Benchmarks from lesson.go functions
// ---------------------------------------------------------------------------

func BenchmarkFibonacci(b *testing.B) {
	b.Run("slow/n=20", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SlowFibonacci(20)
		}
	})
	b.Run("fast/n=20", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FastFibonacci(20)
		}
	})
	b.Run("fast/n=1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FastFibonacci(1000)
		}
	})
}

func BenchmarkSumSquares(b *testing.B) {
	b.Run("loop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SumSquaresLoop(10000)
		}
	})
	b.Run("formula", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SumSquaresFormula(10000)
		}
	})
}
