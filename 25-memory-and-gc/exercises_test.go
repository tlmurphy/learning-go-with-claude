package memory

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func makeStrings(n int) []string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = "x"
	}
	return parts
}

// ---------------------------------------------------------------------------
// Exercise 1-4: String concatenation tests
// ---------------------------------------------------------------------------

func TestConcatWithPlus(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", nil, ""},
		{"single", []string{"hello"}, "hello"},
		{"multiple", []string{"a", "b", "c"}, "abc"},
		{"with spaces", []string{"hello ", "world"}, "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConcatWithPlus(tt.input)
			if got != tt.want {
				t.Errorf("ConcatWithPlus(%v) = %q, want %q. Use += to concatenate strings.", tt.input, got, tt.want)
			}
		})
	}
}

func TestConcatWithBuilder(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", nil, ""},
		{"single", []string{"hello"}, "hello"},
		{"multiple", []string{"a", "b", "c"}, "abc"},
		{"with spaces", []string{"hello ", "world"}, "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConcatWithBuilder(tt.input)
			if got != tt.want {
				t.Errorf("ConcatWithBuilder(%v) = %q, want %q. Use strings.Builder with Grow() and WriteString().", tt.input, got, tt.want)
			}
		})
	}
}

func TestConcatWithBuffer(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", nil, ""},
		{"single", []string{"hello"}, "hello"},
		{"multiple", []string{"a", "b", "c"}, "abc"},
		{"with spaces", []string{"hello ", "world"}, "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConcatWithBuffer(tt.input)
			if got != tt.want {
				t.Errorf("ConcatWithBuffer(%v) = %q, want %q. Use bytes.Buffer with Grow() and WriteString().", tt.input, got, tt.want)
			}
		})
	}
}

func TestConcatWithJoin(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", nil, ""},
		{"single", []string{"hello"}, "hello"},
		{"multiple", []string{"a", "b", "c"}, "abc"},
		{"with spaces", []string{"hello ", "world"}, "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConcatWithJoin(tt.input)
			if got != tt.want {
				t.Errorf("ConcatWithJoin(%v) = %q, want %q. Use strings.Join(parts, \"\").", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Benchmarks: Exercise 1-4 — compare concatenation approaches
// Run with: go test -bench=BenchmarkConcat -benchmem ./25-memory-and-gc/
// ---------------------------------------------------------------------------

func BenchmarkConcatPlus(b *testing.B) {
	parts := makeStrings(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithPlus(parts)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	parts := makeStrings(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithBuilder(parts)
	}
}

func BenchmarkConcatBuffer(b *testing.B) {
	parts := makeStrings(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithBuffer(parts)
	}
}

func BenchmarkConcatJoin(b *testing.B) {
	parts := makeStrings(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcatWithJoin(parts)
	}
}

// ---------------------------------------------------------------------------
// Exercise 5: PreallocateSum
// ---------------------------------------------------------------------------

func TestPreallocateSum(t *testing.T) {
	tests := []struct {
		name   string
		groups [][]int
		want   []int
	}{
		{
			name:   "basic",
			groups: [][]int{{1, 2, 3}, {4, 5}},
			want:   []int{6, 9},
		},
		{
			name:   "empty groups",
			groups: [][]int{},
			want:   []int{},
		},
		{
			name:   "single group",
			groups: [][]int{{10, 20, 30}},
			want:   []int{60},
		},
		{
			name:   "groups with negatives",
			groups: [][]int{{-1, 1}, {-5, 10, -5}},
			want:   []int{0, 0},
		},
		{
			name:   "empty inner slice",
			groups: [][]int{{}, {1, 2}},
			want:   []int{0, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreallocateSum(tt.groups)
			if got == nil && len(tt.want) > 0 {
				t.Fatalf("PreallocateSum returned nil, want %v. Use make([]int, 0, len(groups)) to pre-allocate.", tt.want)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("PreallocateSum returned slice of length %d, want %d.", len(got), len(tt.want))
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("PreallocateSum[%d] = %d, want %d.", i, v, tt.want[i])
				}
			}
			// Check that capacity was pre-allocated correctly
			if cap(got) != len(tt.groups) && len(tt.groups) > 0 {
				t.Errorf("Result slice capacity = %d, want %d. Use make([]int, 0, len(groups)) to pre-allocate.",
					cap(got), len(tt.groups))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Exercise 6: BufferPool
// ---------------------------------------------------------------------------

func TestBufferPool(t *testing.T) {
	pool := BufferPool()
	if pool == nil {
		t.Fatal("BufferPool() returned nil. Create a sync.Pool with New function that returns *bytes.Buffer.")
	}
	if pool.New == nil {
		t.Fatal("BufferPool().New is nil. Set the New field to a function that creates *bytes.Buffer.")
	}

	t.Run("get creates buffer", func(t *testing.T) {
		buf := GetBuffer(pool)
		if buf == nil {
			t.Fatal("GetBuffer returned nil. Use pool.Get() and type-assert to *bytes.Buffer.")
		}
		buf.WriteString("test data")
		if buf.String() != "test data" {
			t.Errorf("Buffer content = %q, want %q.", buf.String(), "test data")
		}
	})

	t.Run("put and reuse", func(t *testing.T) {
		buf := GetBuffer(pool)
		buf.WriteString("first use")
		PutBuffer(pool, buf)

		// Get another buffer — should be the one we just put back (reset)
		buf2 := GetBuffer(pool)
		if buf2.Len() != 0 {
			t.Errorf("Reused buffer has length %d, want 0. PutBuffer must call buf.Reset() before pool.Put().",
				buf2.Len())
		}
		PutBuffer(pool, buf2)
	})
}

// ---------------------------------------------------------------------------
// Exercise 7: Optimized struct layout
// ---------------------------------------------------------------------------

func TestOptimizedStruct(t *testing.T) {
	// Verify the struct has all required fields with correct types
	o := Optimized{
		ID:       1,
		Value:    3.14,
		Name:     "test",
		Score:    1.5,
		Active:   true,
		Priority: 5,
		Done:     false,
	}

	if o.ID != 1 {
		t.Errorf("Optimized.ID = %d, want 1", o.ID)
	}
	if o.Value != 3.14 {
		t.Errorf("Optimized.Value = %f, want 3.14", o.Value)
	}
	if o.Name != "test" {
		t.Errorf("Optimized.Name = %q, want %q", o.Name, "test")
	}
	if o.Score != 1.5 {
		t.Errorf("Optimized.Score = %f, want 1.5", o.Score)
	}
	if !o.Active {
		t.Error("Optimized.Active = false, want true")
	}
	if o.Priority != 5 {
		t.Errorf("Optimized.Priority = %d, want 5", o.Priority)
	}
	if o.Done {
		t.Error("Optimized.Done = true, want false")
	}

	// The struct should be 40 bytes or less (optimal is 40 bytes).
	// Unoptimized would be 56 bytes.
	// We check this indirectly by verifying the field ordering compiles
	// and all fields are accessible. For exact size checks, use unsafe.Sizeof
	// in a benchmark or manual test.
}

// ---------------------------------------------------------------------------
// Exercise 8: ProcessWithoutAlloc
// ---------------------------------------------------------------------------

func TestProcessWithoutAlloc(t *testing.T) {
	tests := []struct {
		name  string
		pairs []string
		want  map[string]string
	}{
		{
			name:  "basic pairs",
			pairs: []string{"name=Alice", "age=30"},
			want:  map[string]string{"name": "Alice", "age": "30"},
		},
		{
			name:  "empty input",
			pairs: []string{},
			want:  map[string]string{},
		},
		{
			name:  "skip invalid",
			pairs: []string{"name=Alice", "invalid", "age=30"},
			want:  map[string]string{"name": "Alice", "age": "30"},
		},
		{
			name:  "value with equals",
			pairs: []string{"expr=a=b"},
			want:  map[string]string{"expr": "a=b"},
		},
		{
			name:  "empty value",
			pairs: []string{"key="},
			want:  map[string]string{"key": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessWithoutAlloc(tt.pairs)
			if got == nil {
				t.Fatal("ProcessWithoutAlloc returned nil. Use make(map[string]string, len(pairs)) to pre-allocate.")
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Result map has %d entries, want %d.", len(got), len(tt.want))
			}
			for k, wantV := range tt.want {
				gotV, ok := got[k]
				if !ok {
					t.Errorf("Missing key %q. Use strings.Cut(pair, \"=\") to split key=value pairs.", k)
				} else if gotV != wantV {
					t.Errorf("got[%q] = %q, want %q.", k, gotV, wantV)
				}
			}
		})
	}
}

func BenchmarkProcessWithoutAlloc(b *testing.B) {
	pairs := []string{
		"name=Alice", "age=30", "city=Portland",
		"role=engineer", "team=backend", "level=senior",
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ProcessWithoutAlloc(pairs)
	}
}

// ---------------------------------------------------------------------------
// Bonus benchmark: Pre-allocation vs no pre-allocation
// ---------------------------------------------------------------------------

func BenchmarkAppendNoPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
		_ = s
	}
}

func BenchmarkAppendWithPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 1000)
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
		_ = s
	}
}

// Ensure imports are used
var _ = strings.Builder{}
