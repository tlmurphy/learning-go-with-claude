package collections

import (
	"reflect"
	"sort"
	"testing"
)

// --- Exercise 1: Slice Operations ---

func TestPrepend(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		val      int
		expected []int
	}{
		{"prepend to non-empty", []int{2, 3, 4}, 1, []int{1, 2, 3, 4}},
		{"prepend to empty", []int{}, 1, []int{1}},
		{"prepend to single", []int{2}, 1, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Prepend(tt.s, tt.val)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Prepend(%v, %d) = %v, want %v.\n"+
					"  Hint: Use append([]int{val}, s...) to prepend.",
					tt.s, tt.val, got, tt.expected)
			}
		})
	}
}

func TestRemoveAt(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		index    int
		expected []int
	}{
		{"remove middle", []int{10, 20, 30, 40}, 1, []int{10, 30, 40}},
		{"remove first", []int{10, 20, 30}, 0, []int{20, 30}},
		{"remove last", []int{10, 20, 30}, 2, []int{10, 20}},
		{"out of bounds negative", []int{10, 20}, -1, []int{10, 20}},
		{"out of bounds too large", []int{10, 20}, 5, []int{10, 20}},
		{"single element", []int{42}, 0, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying test data.
			input := make([]int, len(tt.s))
			copy(input, tt.s)
			got := RemoveAt(input, tt.index)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RemoveAt(%v, %d) = %v, want %v.\n"+
					"  Hint: Use append(s[:index], s[index+1:]...) to remove. Check bounds first!",
					tt.s, tt.index, got, tt.expected)
			}
		})
	}
}

func TestInsertAt(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		index    int
		val      int
		expected []int
	}{
		{"insert middle", []int{10, 30, 40}, 1, 20, []int{10, 20, 30, 40}},
		{"insert at start", []int{20, 30}, 0, 10, []int{10, 20, 30}},
		{"insert at end", []int{10, 20}, 2, 30, []int{10, 20, 30}},
		{"out of bounds", []int{10, 20}, 5, 30, []int{10, 20}},
		{"negative index", []int{10, 20}, -1, 30, []int{10, 20}},
		{"insert into empty at 0", []int{}, 0, 1, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]int, len(tt.s))
			copy(input, tt.s)
			got := InsertAt(input, tt.index, tt.val)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("InsertAt(%v, %d, %d) = %v, want %v.\n"+
					"  Hint: Use append(s[:index], append([]int{val}, s[index:]...)...) to insert.",
					tt.s, tt.index, tt.val, got, tt.expected)
			}
		})
	}
}

// --- Exercise 2: Capacity Detective ---

func TestPredictCapacity(t *testing.T) {
	tests := []struct {
		name        string
		initLen     int
		initCap     int
		appendCount int
		wantLen     int
	}{
		{"no growth needed", 3, 5, 2, 5},
		{"exact fit", 3, 5, 2, 5},
		{"needs growth", 3, 3, 1, 4},
		{"zero append", 3, 5, 0, 3},
		{"empty start", 0, 0, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLen, gotCap := PredictCapacity(tt.initLen, tt.initCap, tt.appendCount)
			if gotLen != tt.wantLen {
				t.Errorf("PredictCapacity(%d, %d, %d) len = %d, want %d.\n"+
					"  Hint: Final len = initialLen + appendCount.",
					tt.initLen, tt.initCap, tt.appendCount, gotLen, tt.wantLen)
			}
			// Capacity should be at least as large as length.
			if gotCap < gotLen {
				t.Errorf("PredictCapacity(%d, %d, %d) cap = %d, which is less than len = %d.\n"+
					"  Capacity must always be >= length.",
					tt.initLen, tt.initCap, tt.appendCount, gotCap, gotLen)
			}
			// If we didn't exceed initial capacity, it should stay the same.
			if tt.initLen+tt.appendCount <= tt.initCap && gotCap != tt.initCap {
				t.Errorf("PredictCapacity(%d, %d, %d) cap = %d, want %d.\n"+
					"  When appending within capacity, the capacity shouldn't change.",
					tt.initLen, tt.initCap, tt.appendCount, gotCap, tt.initCap)
			}
		})
	}
}

// --- Exercise 3: Word Frequency ---

func TestWordFrequency(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected map[string]int
	}{
		{
			name:     "basic counting",
			words:    []string{"go", "is", "go", "great"},
			expected: map[string]int{"go": 2, "is": 1, "great": 1},
		},
		{
			name:     "single word repeated",
			words:    []string{"hello", "hello", "hello"},
			expected: map[string]int{"hello": 3},
		},
		{
			name:     "all unique",
			words:    []string{"a", "b", "c"},
			expected: map[string]int{"a": 1, "b": 1, "c": 1},
		},
		{
			name:     "empty input",
			words:    []string{},
			expected: map[string]int{},
		},
		{
			name:     "case sensitive",
			words:    []string{"Go", "go", "GO"},
			expected: map[string]int{"Go": 1, "go": 1, "GO": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordFrequency(tt.words)
			if got == nil {
				t.Fatal("WordFrequency returned nil. Return an initialized map.\n" +
					"  Hint: Use make(map[string]int) then increment: counts[word]++")
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("WordFrequency(%v) = %v, want %v.\n"+
					"  Hint: Range over words and use counts[word]++. The zero value for int (0) makes this work.",
					tt.words, got, tt.expected)
			}
		})
	}
}

// --- Exercise 4: Set Operations ---

func TestNewStringSet(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  int // expected number of unique items
	}{
		{"with duplicates", []string{"a", "b", "a", "c", "b"}, 3},
		{"all unique", []string{"x", "y", "z"}, 3},
		{"empty", []string{}, 0},
		{"single", []string{"only"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStringSet(tt.items)
			if got == nil {
				t.Fatal("NewStringSet returned nil. Return an initialized map.\n" +
					"  Hint: Use map[string]struct{}{} and add items with set[item] = struct{}{}")
			}
			if len(got) != tt.want {
				t.Errorf("NewStringSet(%v) has %d elements, want %d unique.\n"+
					"  Hint: The map automatically deduplicates — just add all items.",
					tt.items, len(got), tt.want)
			}
		})
	}
}

func TestSetContains(t *testing.T) {
	set := NewStringSet([]string{"apple", "banana", "cherry"})

	tests := []struct {
		item string
		want bool
	}{
		{"apple", true},
		{"banana", true},
		{"grape", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			got := SetContains(set, tt.item)
			if got != tt.want {
				t.Errorf("SetContains(set, %q) = %t, want %t.\n"+
					"  Hint: Use the comma-ok pattern: _, ok := set[item]; return ok",
					tt.item, got, tt.want)
			}
		})
	}
}

func TestSetUnion(t *testing.T) {
	a := NewStringSet([]string{"a", "b", "c"})
	b := NewStringSet([]string{"b", "c", "d"})

	got := SetUnion(a, b)
	if got == nil {
		t.Fatal("SetUnion returned nil.")
	}

	expected := []string{"a", "b", "c", "d"}
	if len(got) != len(expected) {
		t.Errorf("SetUnion has %d elements, want %d.\n"+
			"  Hint: Create a new set, add all elements from both a and b.",
			len(got), len(expected))
	}
	for _, item := range expected {
		if !SetContains(got, item) {
			t.Errorf("SetUnion missing %q", item)
		}
	}
}

func TestSetIntersection(t *testing.T) {
	a := NewStringSet([]string{"a", "b", "c"})
	b := NewStringSet([]string{"b", "c", "d"})

	got := SetIntersection(a, b)
	if got == nil {
		t.Fatal("SetIntersection returned nil.")
	}

	expected := []string{"b", "c"}
	if len(got) != len(expected) {
		t.Errorf("SetIntersection has %d elements, want %d.\n"+
			"  Hint: Iterate over one set and check if each element is in the other.",
			len(got), len(expected))
	}
	for _, item := range expected {
		if !SetContains(got, item) {
			t.Errorf("SetIntersection missing %q", item)
		}
	}
}

func TestSetDifference(t *testing.T) {
	a := NewStringSet([]string{"a", "b", "c"})
	b := NewStringSet([]string{"b", "c", "d"})

	got := SetDifference(a, b)
	if got == nil {
		t.Fatal("SetDifference returned nil.")
	}

	expected := []string{"a"}
	if len(got) != len(expected) {
		t.Errorf("SetDifference has %d elements, want %d.\n"+
			"  Hint: Iterate over a and include elements NOT in b.",
			len(got), len(expected))
	}
	for _, item := range expected {
		if !SetContains(got, item) {
			t.Errorf("SetDifference missing %q", item)
		}
	}
}

// --- Exercise 5: Matrix Operations ---

func TestNewMatrix(t *testing.T) {
	tests := []struct {
		name string
		rows int
		cols int
	}{
		{"3x3", 3, 3},
		{"2x4", 2, 4},
		{"1x1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMatrix(tt.rows, tt.cols)
			if got == nil {
				t.Fatal("NewMatrix returned nil.\n" +
					"  Hint: Create a [][]int with the right dimensions using make().")
			}
			if len(got) != tt.rows {
				t.Errorf("NewMatrix(%d, %d) has %d rows, want %d",
					tt.rows, tt.cols, len(got), tt.rows)
			}
			for i, row := range got {
				if len(row) != tt.cols {
					t.Errorf("Row %d has %d cols, want %d", i, len(row), tt.cols)
				}
				for j, val := range row {
					if val != 0 {
						t.Errorf("Matrix[%d][%d] = %d, want 0 (zero initialized)", i, j, val)
					}
				}
			}
		})
	}
}

func TestMatrixTranspose(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		expected [][]int
	}{
		{
			name: "3x3",
			matrix: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			expected: [][]int{
				{1, 4, 7},
				{2, 5, 8},
				{3, 6, 9},
			},
		},
		{
			name: "2x3 becomes 3x2",
			matrix: [][]int{
				{1, 2, 3},
				{4, 5, 6},
			},
			expected: [][]int{
				{1, 4},
				{2, 5},
				{3, 6},
			},
		},
		{
			name:     "empty",
			matrix:   [][]int{},
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatrixTranspose(tt.matrix)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MatrixTranspose = %v, want %v.\n"+
					"  Hint: Create a new matrix with rows and cols swapped.\n"+
					"  result[j][i] = matrix[i][j]",
					got, tt.expected)
			}
		})
	}
}

// --- Exercise 6: Deduplicate ---

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "with duplicates",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3},
			expected: []int{3, 1, 4, 5, 9, 2, 6},
		},
		{
			name:     "all same",
			input:    []int{7, 7, 7, 7},
			expected: []int{7},
		},
		{
			name:     "already unique",
			input:    []int{1, 2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "empty",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{42},
			expected: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Deduplicate(tt.input)
			if got == nil && len(tt.expected) == 0 {
				// nil is acceptable for empty result
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Deduplicate(%v) = %v, want %v.\n"+
					"  Hint: Use a map[int]struct{} to track seen values.\n"+
					"  Only append to result if the value hasn't been seen before.\n"+
					"  Order must be preserved (first occurrence).",
					tt.input, got, tt.expected)
			}
		})
	}
}

// --- Exercise 7: GroupBy ---

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		expected map[rune][]string
	}{
		{
			name:  "fruits",
			items: []string{"apple", "avocado", "banana", "blueberry", "cherry"},
			expected: map[rune][]string{
				'a': {"apple", "avocado"},
				'b': {"banana", "blueberry"},
				'c': {"cherry"},
			},
		},
		{
			name:     "with empty strings",
			items:    []string{"", "hello", "", "hi"},
			expected: map[rune][]string{'h': {"hello", "hi"}},
		},
		{
			name:     "empty input",
			items:    []string{},
			expected: map[rune][]string{},
		},
		{
			name:  "case sensitive",
			items: []string{"Apple", "apple", "Avocado"},
			expected: map[rune][]string{
				'A': {"Apple", "Avocado"},
				'a': {"apple"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupBy(tt.items)
			if got == nil {
				t.Fatal("GroupBy returned nil. Return an initialized map.")
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GroupBy(%v) = %v, want %v.\n"+
					"  Hint: Get the first rune with []rune(s)[0] or by ranging over s and breaking.\n"+
					"  Use append to build the slices: groups[firstRune] = append(groups[firstRune], s)",
					tt.items, got, tt.expected)
			}
		})
	}
}

// --- Exercise 8: Stack and Queue ---

func TestIntStack(t *testing.T) {
	s := &IntStack{}

	// Test empty stack.
	if s.Len() != 0 {
		t.Error("New stack should have Len() == 0")
	}
	if _, ok := s.Pop(); ok {
		t.Error("Pop on empty stack should return false.\n" +
			"  Hint: Check if len(s.data) == 0 before popping.")
	}
	if _, ok := s.Peek(); ok {
		t.Error("Peek on empty stack should return false.")
	}

	// Push and verify.
	s.Push(10)
	s.Push(20)
	s.Push(30)

	if s.Len() != 3 {
		t.Errorf("After 3 pushes, Len() = %d, want 3", s.Len())
	}

	// Peek should return top without removing.
	if val, ok := s.Peek(); !ok || val != 30 {
		t.Errorf("Peek() = (%d, %t), want (30, true).\n"+
			"  Hint: Return the last element of s.data without removing it.",
			val, ok)
	}
	if s.Len() != 3 {
		t.Error("Peek should not change the stack length.")
	}

	// Pop should return in LIFO order.
	vals := []int{}
	for s.Len() > 0 {
		val, ok := s.Pop()
		if !ok {
			t.Fatal("Pop returned false when stack was not empty")
		}
		vals = append(vals, val)
	}
	expected := []int{30, 20, 10}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Pop order = %v, want %v (LIFO).\n"+
			"  Hint: Pop from the end of the slice (last element).",
			vals, expected)
	}
}

func TestIntQueue(t *testing.T) {
	q := &IntQueue{}

	// Test empty queue.
	if q.Len() != 0 {
		t.Error("New queue should have Len() == 0")
	}
	if _, ok := q.Dequeue(); ok {
		t.Error("Dequeue on empty queue should return false.\n" +
			"  Hint: Check if len(q.data) == 0 before dequeuing.")
	}

	// Enqueue and verify.
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	if q.Len() != 3 {
		t.Errorf("After 3 enqueues, Len() = %d, want 3", q.Len())
	}

	// Dequeue should return in FIFO order.
	vals := []int{}
	for q.Len() > 0 {
		val, ok := q.Dequeue()
		if !ok {
			t.Fatal("Dequeue returned false when queue was not empty")
		}
		vals = append(vals, val)
	}
	expected := []int{10, 20, 30}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Dequeue order = %v, want %v (FIFO).\n"+
			"  Hint: Dequeue from the front of the slice (first element).\n"+
			"  Use q.data[0] to get the value, then q.data = q.data[1:] to remove it.",
			vals, expected)
	}
}

// --- Exercise 9: Merge Sorted ---

func TestMergeSorted(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{
			name:     "basic merge",
			a:        []int{1, 3, 5},
			b:        []int{2, 4, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "overlapping values",
			a:        []int{1, 3, 5},
			b:        []int{1, 2, 3},
			expected: []int{1, 1, 2, 3, 3, 5},
		},
		{
			name:     "first empty",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "second empty",
			a:        []int{1, 2, 3},
			b:        []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "both empty",
			a:        []int{},
			b:        []int{},
			expected: []int{},
		},
		{
			name:     "non-overlapping ranges",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "single elements",
			a:        []int{2},
			b:        []int{1},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeSorted(tt.a, tt.b)
			if got == nil && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MergeSorted(%v, %v) = %v, want %v.\n"+
					"  Hint: Use two index pointers (i for a, j for b).\n"+
					"  Compare a[i] and b[j], append the smaller one, advance that pointer.\n"+
					"  After one runs out, append the rest of the other.",
					tt.a, tt.b, got, tt.expected)
			}

			// Verify the result is actually sorted.
			if !sort.IntsAreSorted(got) {
				t.Errorf("Result is not sorted: %v.\n"+
					"  Hint: Don't concatenate and sort — use the merge algorithm with two pointers.",
					got)
			}
		})
	}
}
