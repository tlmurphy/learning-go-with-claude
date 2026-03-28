package functions

import (
	"strings"
	"testing"
)

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name          string
		a, b          int
		wantQuotient  int
		wantRemainder int
		wantErr       bool
	}{
		{"basic division", 17, 5, 3, 2, false},
		{"even division", 10, 2, 5, 0, false},
		{"zero numerator", 0, 5, 0, 0, false},
		{"negative division", -17, 5, -3, -2, false},
		{"divide by zero", 10, 0, 0, 0, true},
		{"one", 1, 1, 1, 0, false},
		{"larger remainder", 3, 7, 0, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, r, err := SafeDivide(tt.a, tt.b)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SafeDivide(%d, %d): expected error for division by zero, got nil.\n"+
						"  Hint: Check if b == 0 and return an error using fmt.Errorf().",
						tt.a, tt.b)
				}
				return
			}

			if err != nil {
				t.Errorf("SafeDivide(%d, %d): unexpected error: %v", tt.a, tt.b, err)
				return
			}

			if q != tt.wantQuotient {
				t.Errorf("SafeDivide(%d, %d) quotient = %d, want %d.\n"+
					"  Hint: Use the / operator for integer division.",
					tt.a, tt.b, q, tt.wantQuotient)
			}
			if r != tt.wantRemainder {
				t.Errorf("SafeDivide(%d, %d) remainder = %d, want %d.\n"+
					"  Hint: Use the %% operator for remainder.",
					tt.a, tt.b, r, tt.wantRemainder)
			}
		})
	}
}

func TestVariadicSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{"multiple numbers", []int{1, 2, 3, 4, 5}, 15},
		{"single number", []int{42}, 42},
		{"no numbers", []int{}, 0},
		{"with negatives", []int{-1, 1, -2, 2}, 0},
		{"large sum", []int{100, 200, 300}, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VariadicSum(tt.nums...)
			if got != tt.expected {
				t.Errorf("VariadicSum(%v) = %d, want %d.\n"+
					"  Hint: Range over the variadic parameter (it's a slice) and sum the values.",
					tt.nums, got, tt.expected)
			}
		})
	}
}

func TestVariadicAverage(t *testing.T) {
	tests := []struct {
		name     string
		nums     []float64
		expected float64
		wantErr  bool
	}{
		{"multiple numbers", []float64{1, 2, 3, 4, 5}, 3.0, false},
		{"single number", []float64{42.0}, 42.0, false},
		{"no numbers", []float64{}, 0, true},
		{"decimal result", []float64{1, 2}, 1.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VariadicAverage(tt.nums...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("VariadicAverage() with no args: expected error, got nil.\n" +
						"  Hint: Return an error when no values are provided (division by zero!).",
					)
				}
				return
			}

			if err != nil {
				t.Errorf("VariadicAverage(%v): unexpected error: %v", tt.nums, err)
				return
			}

			if diff := got - tt.expected; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("VariadicAverage(%v) = %f, want %f.\n"+
					"  Hint: Sum all values and divide by float64(len(nums)).",
					tt.nums, got, tt.expected)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		nums      []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "filter positives",
			nums:      []int{-2, -1, 0, 1, 2},
			predicate: func(n int) bool { return n > 0 },
			expected:  []int{1, 2},
		},
		{
			name:      "filter evens",
			nums:      []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "filter none match",
			nums:      []int{1, 2, 3},
			predicate: func(n int) bool { return n > 10 },
			expected:  []int{},
		},
		{
			name:      "filter all match",
			nums:      []int{1, 2, 3},
			predicate: func(n int) bool { return n > 0 },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "empty input",
			nums:      []int{},
			predicate: func(n int) bool { return true },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.nums, tt.predicate)
			if got == nil {
				t.Fatal("Filter returned nil. Return an empty slice, not nil.\n" +
					"  Hint: Initialize result with make([]int, 0) or var result []int followed by append.")
			}
			if len(got) != len(tt.expected) {
				t.Fatalf("Filter returned %d elements, want %d.\n  got:  %v\n  want: %v",
					len(got), len(tt.expected), got, tt.expected)
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("Filter result[%d] = %d, want %d", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name      string
		nums      []int
		transform func(int) int
		expected  []int
	}{
		{
			name:      "double",
			nums:      []int{1, 2, 3},
			transform: func(n int) int { return n * 2 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "square",
			nums:      []int{1, 2, 3, 4},
			transform: func(n int) int { return n * n },
			expected:  []int{1, 4, 9, 16},
		},
		{
			name:      "negate",
			nums:      []int{1, -2, 3},
			transform: func(n int) int { return -n },
			expected:  []int{-1, 2, -3},
		},
		{
			name:      "empty input",
			nums:      []int{},
			transform: func(n int) int { return n },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.nums, tt.transform)
			if got == nil {
				t.Fatal("Map returned nil. Return an empty slice for empty input.\n" +
					"  Hint: Use make([]int, len(nums)) to pre-allocate the result slice.")
			}
			if len(got) != len(tt.expected) {
				t.Fatalf("Map returned %d elements, want %d", len(got), len(tt.expected))
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("Map result[%d] = %d, want %d", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name        string
		nums        []int
		initial     int
		accumulator func(int, int) int
		expected    int
	}{
		{
			name:        "sum",
			nums:        []int{1, 2, 3, 4, 5},
			initial:     0,
			accumulator: func(acc, n int) int { return acc + n },
			expected:    15,
		},
		{
			name:        "product",
			nums:        []int{1, 2, 3, 4},
			initial:     1,
			accumulator: func(acc, n int) int { return acc * n },
			expected:    24,
		},
		{
			name:    "max",
			nums:    []int{3, 1, 4, 1, 5, 9},
			initial: 0,
			accumulator: func(acc, n int) int {
				if n > acc {
					return n
				}
				return acc
			},
			expected: 9,
		},
		{
			name:        "empty with initial",
			nums:        []int{},
			initial:     42,
			accumulator: func(acc, n int) int { return acc + n },
			expected:    42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reduce(tt.nums, tt.initial, tt.accumulator)
			if got != tt.expected {
				t.Errorf("Reduce(%v, %d, fn) = %d, want %d.\n"+
					"  Hint: Start with the initial value, then apply accumulator(result, element) for each element.",
					tt.nums, tt.initial, got, tt.expected)
			}
		})
	}
}

func TestNewCounter(t *testing.T) {
	inc, dec, val := NewCounter(10)

	if inc == nil || dec == nil || val == nil {
		t.Fatal("NewCounter returned nil functions. Return three closures that share a counter variable.\n" +
			"  Hint: Declare a count variable, then return three anonymous functions that capture it.")
	}

	tests := []struct {
		name     string
		action   func() int
		expected int
	}{
		{"initial value", val, 10},
		{"increment 1", inc, 11},
		{"increment 2", inc, 12},
		{"value after increments", val, 12},
		{"decrement 1", dec, 11},
		{"decrement 2", dec, 10},
		{"value after decrements", val, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.action()
			if got != tt.expected {
				t.Errorf("counter %s: got %d, want %d.\n"+
					"  Hint: All three closures must capture the same variable.",
					tt.name, got, tt.expected)
			}
		})
	}

	// Verify independence: a second counter shouldn't affect the first.
	inc2, _, val2 := NewCounter(0)
	inc2()
	inc2()

	if val2() != 2 {
		t.Error("Second counter should be at 2")
	}
	if val() != 10 {
		t.Error("First counter should still be at 10 — counters must be independent.\n" +
			"  Hint: Each call to NewCounter should create a new count variable.")
	}
}

func TestLogger(t *testing.T) {
	upper := func(s string) string {
		return strings.ToUpper(s)
	}

	var log []string
	wrapped := Logger(upper, &log)

	if wrapped == nil {
		t.Fatal("Logger returned nil. Return a function that wraps the original.\n" +
			"  Hint: Return func(s string) string { ... } that logs before and after calling fn.")
	}

	result := wrapped("hello")

	if result != "HELLO" {
		t.Errorf("Logger-wrapped function returned %q, want %q.\n"+
			"  The wrapped function must return the same result as the original.",
			result, "HELLO")
	}

	if len(log) != 2 {
		t.Fatalf("Expected 2 log entries, got %d: %v.\n"+
			"  Hint: Append \"calling with: <input>\" before calling fn,\n"+
			"  and \"returned: <result>\" after.",
			len(log), log)
	}

	if log[0] != "calling with: hello" {
		t.Errorf("log[0] = %q, want %q", log[0], "calling with: hello")
	}
	if log[1] != "returned: HELLO" {
		t.Errorf("log[1] = %q, want %q", log[1], "returned: HELLO")
	}

	// Call again to verify log accumulates.
	wrapped("world")
	if len(log) != 4 {
		t.Errorf("After second call, expected 4 log entries, got %d.\n"+
			"  The log should accumulate across calls.",
			len(log))
	}
}

func TestDeferOrder(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected string
	}{
		{"n=0", 0, ""},
		{"n=1", 1, "1"},
		{"n=3", 3, "3,2,1"},
		{"n=5", 5, "5,4,3,2,1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeferOrder(tt.n)
			if got != tt.expected {
				t.Errorf("DeferOrder(%d) = %q, want %q.\n"+
					"  Hint: Defer is LIFO — if you defer 1,2,3, they execute as 3,2,1.\n"+
					"  Build the string in reverse order: n, n-1, ..., 2, 1.",
					tt.n, got, tt.expected)
			}
		})
	}
}

func TestCompose(t *testing.T) {
	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }
	square := func(x int) int { return x * x }
	negate := func(x int) int { return -x }

	tests := []struct {
		name     string
		f, g     func(int) int
		input    int
		expected int
	}{
		{"addOne after double", addOne, double, 3, 7},  // double(3)=6, addOne(6)=7
		{"double after addOne", double, addOne, 3, 8},  // addOne(3)=4, double(4)=8
		{"square after double", square, double, 3, 36}, // double(3)=6, square(6)=36
		{"negate after addOne", negate, addOne, 5, -6}, // addOne(5)=6, negate(6)=-6
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			composed := Compose(tt.f, tt.g)
			if composed == nil {
				t.Fatal("Compose returned nil. Return a function that applies g first, then f.\n" +
					"  Hint: return func(x int) int { return f(g(x)) }")
			}
			got := composed(tt.input)
			if got != tt.expected {
				t.Errorf("Compose(f, g)(%d) = %d, want %d.\n"+
					"  Remember: Compose(f, g)(x) = f(g(x)) — g is applied first.",
					tt.input, got, tt.expected)
			}
		})
	}
}

func TestMemoize(t *testing.T) {
	expensive := func(n int) int {
		return n * n
	}

	memo, getMemoCallCount := Memoize(expensive)

	if memo == nil || getMemoCallCount == nil {
		t.Fatal("Memoize returned nil functions. Return a memoized function and a call counter.\n" +
			"  Hint: Use a map[int]int as a cache inside the closure.")
	}

	// First call — should compute.
	result := memo(4)
	if result != 16 {
		t.Errorf("memo(4) = %d, want 16", result)
	}
	if getMemoCallCount() != 1 {
		t.Errorf("After first call, callCount = %d, want 1.\n"+
			"  The original function should be called once for a new input.",
			getMemoCallCount())
	}

	// Second call with same input — should use cache.
	result = memo(4)
	if result != 16 {
		t.Errorf("memo(4) cached = %d, want 16", result)
	}
	if getMemoCallCount() != 1 {
		t.Errorf("After cached call, callCount = %d, want 1.\n"+
			"  The original function should NOT be called again for cached input.",
			getMemoCallCount())
	}

	// New input — should compute again.
	result = memo(5)
	if result != 25 {
		t.Errorf("memo(5) = %d, want 25", result)
	}
	if getMemoCallCount() != 2 {
		t.Errorf("After new input, callCount = %d, want 2.\n"+
			"  Each unique input should call the original function exactly once.",
			getMemoCallCount())
	}

	// Verify cached values persist.
	if memo(4) != 16 || memo(5) != 25 {
		t.Error("Cached values should persist across calls.")
	}
	if getMemoCallCount() != 2 {
		t.Errorf("Call count should still be 2 after accessing cached values, got %d.",
			getMemoCallCount())
	}
}
