package generics

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

// =============================================================================
// Exercise 1: MinSlice and MaxSlice
// =============================================================================

func TestMinSlice(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		tests := []struct {
			name  string
			input []int
			want  int
			found bool
		}{
			{"basic", []int{3, 1, 4, 1, 5}, 1, true},
			{"single", []int{42}, 42, true},
			{"negatives", []int{-3, -1, -4}, -4, true},
			{"empty", []int{}, 0, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, found := MinSlice(tt.input)
				if found != tt.found {
					t.Fatalf("MinSlice(%v) found = %v, want %v", tt.input, found, tt.found)
				}
				if found && got != tt.want {
					t.Errorf("MinSlice(%v) = %d, want %d", tt.input, got, tt.want)
				}
			})
		}
	})

	t.Run("strings", func(t *testing.T) {
		got, found := MinSlice([]string{"banana", "apple", "cherry"})
		if !found || got != "apple" {
			t.Errorf("MinSlice(strings) = (%q, %v), want (\"apple\", true)", got, found)
		}
	})

	t.Run("floats", func(t *testing.T) {
		got, found := MinSlice([]float64{3.14, 2.71, 1.41})
		if !found || got != 1.41 {
			t.Errorf("MinSlice(floats) = (%f, %v), want (1.41, true)", got, found)
		}
	})
}

func TestMaxSlice(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		tests := []struct {
			name  string
			input []int
			want  int
			found bool
		}{
			{"basic", []int{3, 1, 4, 1, 5}, 5, true},
			{"single", []int{42}, 42, true},
			{"negatives", []int{-3, -1, -4}, -1, true},
			{"empty", []int{}, 0, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, found := MaxSlice(tt.input)
				if found != tt.found {
					t.Fatalf("MaxSlice(%v) found = %v, want %v", tt.input, found, tt.found)
				}
				if found && got != tt.want {
					t.Errorf("MaxSlice(%v) = %d, want %d", tt.input, got, tt.want)
				}
			})
		}
	})

	t.Run("strings", func(t *testing.T) {
		got, found := MaxSlice([]string{"banana", "apple", "cherry"})
		if !found || got != "cherry" {
			t.Errorf("MaxSlice(strings) = (%q, %v), want (\"cherry\", true)", got, found)
		}
	})
}

// =============================================================================
// Exercise 2: MapSlice, FilterSlice, ReduceSlice
// =============================================================================

func TestMapSlice(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		got := MapSlice([]int{1, 2, 3}, func(n int) string {
			return fmt.Sprintf("#%d", n)
		})
		expected := []string{"#1", "#2", "#3"}
		if len(got) != len(expected) {
			t.Fatalf("MapSlice returned %d items, want %d", len(got), len(expected))
		}
		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("MapSlice result[%d] = %q, want %q", i, got[i], expected[i])
			}
		}
	})

	t.Run("double ints", func(t *testing.T) {
		got := MapSlice([]int{1, 2, 3}, func(n int) int { return n * 2 })
		expected := []int{2, 4, 6}
		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("result[%d] = %d, want %d", i, got[i], expected[i])
			}
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := MapSlice([]int{}, func(n int) int { return n })
		if len(got) != 0 {
			t.Errorf("expected empty result, got %v", got)
		}
	})
}

func TestFilterSlice(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		got := FilterSlice([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
		expected := []int{2, 4, 6}
		if len(got) != len(expected) {
			t.Fatalf("got %d items, want %d", len(got), len(expected))
		}
		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("result[%d] = %d, want %d", i, got[i], expected[i])
			}
		}
	})

	t.Run("filter strings by length", func(t *testing.T) {
		got := FilterSlice([]string{"a", "bb", "ccc", "dd"}, func(s string) bool { return len(s) > 1 })
		expected := []string{"bb", "ccc", "dd"}
		if len(got) != len(expected) {
			t.Fatalf("got %d items, want %d", len(got), len(expected))
		}
		for i := range got {
			if got[i] != expected[i] {
				t.Errorf("result[%d] = %q, want %q", i, got[i], expected[i])
			}
		}
	})

	t.Run("nothing passes filter", func(t *testing.T) {
		got := FilterSlice([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 })
		if len(got) != 0 {
			t.Errorf("expected empty result, got %v", got)
		}
	})
}

func TestReduceSlice(t *testing.T) {
	t.Run("sum integers", func(t *testing.T) {
		got := ReduceSlice([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int { return acc + n })
		if got != 15 {
			t.Errorf("sum = %d, want 15", got)
		}
	})

	t.Run("concatenate strings", func(t *testing.T) {
		got := ReduceSlice([]string{"a", "b", "c"}, "", func(acc string, s string) string { return acc + s })
		if got != "abc" {
			t.Errorf("concat = %q, want %q", got, "abc")
		}
	})

	t.Run("count elements", func(t *testing.T) {
		got := ReduceSlice([]string{"a", "b", "c"}, 0, func(acc int, s string) int { return acc + 1 })
		if got != 3 {
			t.Errorf("count = %d, want 3", got)
		}
	})

	t.Run("empty slice returns initial", func(t *testing.T) {
		got := ReduceSlice([]int{}, 42, func(acc, n int) int { return acc + n })
		if got != 42 {
			t.Errorf("reduce of empty slice = %d, want initial value 42", got)
		}
	})
}

// =============================================================================
// Exercise 3: ExerciseStack
// =============================================================================

func TestExerciseStack(t *testing.T) {
	t.Run("push and pop", func(t *testing.T) {
		s := NewExerciseStack[int]()
		if s == nil {
			t.Fatal("NewExerciseStack returned nil")
		}

		s.Push(1)
		s.Push(2)
		s.Push(3)

		if s.Size() != 3 {
			t.Errorf("Size() = %d, want 3", s.Size())
		}

		val, ok := s.Pop()
		if !ok || val != 3 {
			t.Errorf("Pop() = (%d, %v), want (3, true)", val, ok)
		}
		val, ok = s.Pop()
		if !ok || val != 2 {
			t.Errorf("Pop() = (%d, %v), want (2, true)", val, ok)
		}
		val, ok = s.Pop()
		if !ok || val != 1 {
			t.Errorf("Pop() = (%d, %v), want (1, true)", val, ok)
		}

		_, ok = s.Pop()
		if ok {
			t.Error("Pop() on empty stack should return false")
		}
	})

	t.Run("peek doesn't remove", func(t *testing.T) {
		s := NewExerciseStack[string]()
		if s == nil {
			t.Fatal("NewExerciseStack returned nil")
		}

		s.Push("hello")
		val, ok := s.Peek()
		if !ok || val != "hello" {
			t.Errorf("Peek() = (%q, %v), want (\"hello\", true)", val, ok)
		}
		if s.Size() != 1 {
			t.Errorf("Size after Peek = %d, want 1 (Peek should not remove)", s.Size())
		}
	})

	t.Run("empty stack", func(t *testing.T) {
		s := NewExerciseStack[int]()
		if s == nil {
			t.Fatal("NewExerciseStack returned nil")
		}

		if s.Size() != 0 {
			t.Errorf("new stack Size() = %d, want 0", s.Size())
		}
		_, ok := s.Peek()
		if ok {
			t.Error("Peek on empty stack should return false")
		}
	})

	t.Run("ToSlice", func(t *testing.T) {
		s := NewExerciseStack[int]()
		if s == nil {
			t.Fatal("NewExerciseStack returned nil")
		}

		s.Push(10)
		s.Push(20)
		s.Push(30)

		slice := s.ToSlice()
		expected := []int{10, 20, 30}
		if len(slice) != len(expected) {
			t.Fatalf("ToSlice() has %d items, want %d", len(slice), len(expected))
		}
		for i := range slice {
			if slice[i] != expected[i] {
				t.Errorf("ToSlice()[%d] = %d, want %d", i, slice[i], expected[i])
			}
		}
	})
}

// =============================================================================
// Exercise 4: ExerciseSet
// =============================================================================

func TestExerciseSet(t *testing.T) {
	t.Run("add contains remove", func(t *testing.T) {
		s := NewExerciseSet[string]()
		if s == nil {
			t.Fatal("NewExerciseSet returned nil")
		}

		s.Add("apple")
		s.Add("banana")
		s.Add("apple") // Duplicate — should not increase size

		if s.Size() != 2 {
			t.Errorf("Size() = %d, want 2 (duplicates should be ignored)", s.Size())
		}
		if !s.Contains("apple") {
			t.Error("set should contain apple")
		}
		if !s.Contains("banana") {
			t.Error("set should contain banana")
		}
		if s.Contains("cherry") {
			t.Error("set should not contain cherry")
		}

		s.Remove("apple")
		if s.Contains("apple") {
			t.Error("set should not contain apple after removal")
		}
		if s.Size() != 1 {
			t.Errorf("Size() = %d after removal, want 1", s.Size())
		}
	})

	t.Run("union", func(t *testing.T) {
		s1 := NewExerciseSet[int]()
		s2 := NewExerciseSet[int]()
		if s1 == nil || s2 == nil {
			t.Fatal("NewExerciseSet returned nil")
		}

		s1.Add(1)
		s1.Add(2)
		s1.Add(3)

		s2.Add(3)
		s2.Add(4)
		s2.Add(5)

		union := s1.Union(s2)
		if union == nil {
			t.Fatal("Union returned nil")
		}
		if union.Size() != 5 {
			t.Errorf("union Size() = %d, want 5", union.Size())
		}
		for _, v := range []int{1, 2, 3, 4, 5} {
			if !union.Contains(v) {
				t.Errorf("union should contain %d", v)
			}
		}
	})

	t.Run("intersection", func(t *testing.T) {
		s1 := NewExerciseSet[int]()
		s2 := NewExerciseSet[int]()
		if s1 == nil || s2 == nil {
			t.Fatal("NewExerciseSet returned nil")
		}

		s1.Add(1)
		s1.Add(2)
		s1.Add(3)

		s2.Add(2)
		s2.Add(3)
		s2.Add(4)

		inter := s1.Intersection(s2)
		if inter == nil {
			t.Fatal("Intersection returned nil")
		}
		if inter.Size() != 2 {
			t.Errorf("intersection Size() = %d, want 2", inter.Size())
		}
		if !inter.Contains(2) || !inter.Contains(3) {
			t.Error("intersection should contain 2 and 3")
		}
		if inter.Contains(1) || inter.Contains(4) {
			t.Error("intersection should not contain 1 or 4")
		}
	})

	t.Run("difference", func(t *testing.T) {
		s1 := NewExerciseSet[int]()
		s2 := NewExerciseSet[int]()
		if s1 == nil || s2 == nil {
			t.Fatal("NewExerciseSet returned nil")
		}

		s1.Add(1)
		s1.Add(2)
		s1.Add(3)

		s2.Add(2)
		s2.Add(4)

		diff := s1.Difference(s2)
		if diff == nil {
			t.Fatal("Difference returned nil")
		}
		if diff.Size() != 2 {
			t.Errorf("difference Size() = %d, want 2", diff.Size())
		}
		if !diff.Contains(1) || !diff.Contains(3) {
			t.Error("difference should contain 1 and 3")
		}
		if diff.Contains(2) {
			t.Error("difference should not contain 2 (it's in both sets)")
		}
	})
}

// =============================================================================
// Exercise 5: GenericCache
// =============================================================================

func TestGenericCache(t *testing.T) {
	t.Run("string cache", func(t *testing.T) {
		cache := NewGenericCache[string]()
		if cache == nil {
			t.Fatal("NewGenericCache returned nil")
		}

		cache.Set("greeting", "hello")
		cache.Set("farewell", "goodbye")

		val, ok := cache.Get("greeting")
		if !ok || val != "hello" {
			t.Errorf("Get(greeting) = (%q, %v), want (\"hello\", true)", val, ok)
		}

		_, ok = cache.Get("missing")
		if ok {
			t.Error("Get(missing) should return false")
		}

		if cache.Size() != 2 {
			t.Errorf("Size() = %d, want 2", cache.Size())
		}
	})

	t.Run("int cache", func(t *testing.T) {
		cache := NewGenericCache[int]()
		if cache == nil {
			t.Fatal("NewGenericCache returned nil")
		}

		cache.Set("age", 30)
		cache.Set("score", 100)

		val, ok := cache.Get("age")
		if !ok || val != 30 {
			t.Errorf("Get(age) = (%d, %v), want (30, true)", val, ok)
		}
	})

	t.Run("delete and keys", func(t *testing.T) {
		cache := NewGenericCache[float64]()
		if cache == nil {
			t.Fatal("NewGenericCache returned nil")
		}

		cache.Set("pi", 3.14)
		cache.Set("e", 2.71)
		cache.Set("phi", 1.618)

		cache.Delete("e")
		if cache.Size() != 2 {
			t.Errorf("Size after delete = %d, want 2", cache.Size())
		}
		if _, ok := cache.Get("e"); ok {
			t.Error("e should be deleted")
		}

		keys := cache.Keys()
		sort.Strings(keys)
		expected := []string{"phi", "pi"}
		if len(keys) != len(expected) {
			t.Fatalf("Keys() has %d entries, want %d", len(keys), len(expected))
		}
		for i := range keys {
			if keys[i] != expected[i] {
				t.Errorf("keys[%d] = %q, want %q", i, keys[i], expected[i])
			}
		}
	})
}

// =============================================================================
// Exercise 6: Custom Constraint (SumAll and Average)
// =============================================================================

func TestSumAll(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		got := SumAll([]int{1, 2, 3, 4, 5})
		if got != 15 {
			t.Errorf("SumAll(1..5) = %d, want 15", got)
		}
	})

	t.Run("floats", func(t *testing.T) {
		got := SumAll([]float64{1.5, 2.5, 3.0})
		if got != 7.0 {
			t.Errorf("SumAll(floats) = %f, want 7.0", got)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := SumAll([]int{})
		if got != 0 {
			t.Errorf("SumAll(empty) = %d, want 0", got)
		}
	})

	t.Run("named type", func(t *testing.T) {
		// Tests that ~int works with named types
		type Score int
		got := SumAll([]Score{10, 20, 30})
		if got != 60 {
			t.Errorf("SumAll(Score) = %d, want 60", got)
		}
	})
}

func TestAverage(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		got := Average([]int{10, 20, 30})
		if got != 20.0 {
			t.Errorf("Average(10,20,30) = %f, want 20.0", got)
		}
	})

	t.Run("floats", func(t *testing.T) {
		got := Average([]float64{1.0, 2.0, 3.0})
		if math.Abs(got-2.0) > 0.001 {
			t.Errorf("Average(1.0,2.0,3.0) = %f, want 2.0", got)
		}
	})

	t.Run("empty returns zero", func(t *testing.T) {
		got := Average([]int{})
		if got != 0 {
			t.Errorf("Average(empty) = %f, want 0", got)
		}
	})
}

// =============================================================================
// Exercise 7: ExerciseResult
// =============================================================================

func TestExerciseResult(t *testing.T) {
	t.Run("Ok result", func(t *testing.T) {
		r := NewOk(42)
		if !r.IsOk() {
			t.Error("IsOk() should be true for Ok result")
		}
		if r.IsErr() {
			t.Error("IsErr() should be false for Ok result")
		}
		val, ok := r.Value()
		if !ok || val != 42 {
			t.Errorf("Value() = (%d, %v), want (42, true)", val, ok)
		}
		if r.Error() != nil {
			t.Errorf("Error() = %v, want nil", r.Error())
		}
	})

	t.Run("Err result", func(t *testing.T) {
		r := NewErr[string](fmt.Errorf("something failed"))
		if r.IsOk() {
			t.Error("IsOk() should be false for Err result")
		}
		if !r.IsErr() {
			t.Error("IsErr() should be true for Err result")
		}
		_, ok := r.Value()
		if ok {
			t.Error("Value() should return false for Err result")
		}
		if r.Error() == nil {
			t.Error("Error() should not be nil for Err result")
		}
	})

	t.Run("UnwrapOrDefault", func(t *testing.T) {
		okResult := NewOk("success")
		if got := okResult.UnwrapOrDefault("default"); got != "success" {
			t.Errorf("UnwrapOrDefault = %q, want %q", got, "success")
		}

		errResult := NewErr[string](fmt.Errorf("fail"))
		if got := errResult.UnwrapOrDefault("default"); got != "default" {
			t.Errorf("UnwrapOrDefault = %q, want %q", got, "default")
		}
	})

	t.Run("MapResult success", func(t *testing.T) {
		r := NewOk(10)
		mapped := MapResult(r, func(n int) string {
			return fmt.Sprintf("value=%d", n)
		})
		if !mapped.IsOk() {
			t.Fatal("mapped result should be Ok")
		}
		val, _ := mapped.Value()
		if val != "value=10" {
			t.Errorf("mapped value = %q, want %q", val, "value=10")
		}
	})

	t.Run("MapResult error propagation", func(t *testing.T) {
		r := NewErr[int](fmt.Errorf("original error"))
		mapped := MapResult(r, func(n int) string {
			return "should not run"
		})
		if !mapped.IsErr() {
			t.Error("mapped Err result should still be Err")
		}
		if mapped.Error() == nil || mapped.Error().Error() != "original error" {
			t.Errorf("error should propagate, got: %v", mapped.Error())
		}
	})
}

// =============================================================================
// Exercise 8: Generic Linked List
// =============================================================================

func TestLinkedList(t *testing.T) {
	t.Run("prepend", func(t *testing.T) {
		ll := NewLinkedList[int]()
		if ll == nil {
			t.Fatal("NewLinkedList returned nil")
		}

		ll.Prepend(3)
		ll.Prepend(2)
		ll.Prepend(1)

		if ll.Len() != 3 {
			t.Errorf("Len() = %d, want 3", ll.Len())
		}

		head, ok := ll.Head()
		if !ok || head != 1 {
			t.Errorf("Head() = (%d, %v), want (1, true)", head, ok)
		}

		slice := ll.ToSlice()
		expected := []int{1, 2, 3}
		if len(slice) != len(expected) {
			t.Fatalf("ToSlice() has %d items, want %d", len(slice), len(expected))
		}
		for i := range slice {
			if slice[i] != expected[i] {
				t.Errorf("ToSlice()[%d] = %d, want %d", i, slice[i], expected[i])
			}
		}
	})

	t.Run("append", func(t *testing.T) {
		ll := NewLinkedList[string]()
		if ll == nil {
			t.Fatal("NewLinkedList returned nil")
		}

		ll.Append("a")
		ll.Append("b")
		ll.Append("c")

		slice := ll.ToSlice()
		expected := []string{"a", "b", "c"}
		if len(slice) != len(expected) {
			t.Fatalf("ToSlice() has %d items, want %d", len(slice), len(expected))
		}
		for i := range slice {
			if slice[i] != expected[i] {
				t.Errorf("ToSlice()[%d] = %q, want %q", i, slice[i], expected[i])
			}
		}
	})

	t.Run("empty list", func(t *testing.T) {
		ll := NewLinkedList[int]()
		if ll == nil {
			t.Fatal("NewLinkedList returned nil")
		}

		if ll.Len() != 0 {
			t.Errorf("Len() of empty list = %d, want 0", ll.Len())
		}
		_, ok := ll.Head()
		if ok {
			t.Error("Head() of empty list should return false")
		}
		if ll.String() != "[]" {
			t.Errorf("String() of empty list = %q, want %q", ll.String(), "[]")
		}
	})

	t.Run("ForEach", func(t *testing.T) {
		ll := NewLinkedList[int]()
		if ll == nil {
			t.Fatal("NewLinkedList returned nil")
		}

		ll.Append(10)
		ll.Append(20)
		ll.Append(30)

		var sum int
		ll.ForEach(func(v int) {
			sum += v
		})
		if sum != 60 {
			t.Errorf("sum of ForEach = %d, want 60", sum)
		}
	})

	t.Run("String representation", func(t *testing.T) {
		ll := NewLinkedList[int]()
		if ll == nil {
			t.Fatal("NewLinkedList returned nil")
		}

		ll.Append(1)
		ll.Append(2)
		ll.Append(3)

		got := ll.String()
		want := "[1 -> 2 -> 3]"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("mixed prepend and append", func(t *testing.T) {
		ll := NewLinkedList[int]()
		if ll == nil {
			t.Fatal("NewLinkedList returned nil")
		}

		ll.Append(2)
		ll.Prepend(1)
		ll.Append(3)
		ll.Prepend(0)

		slice := ll.ToSlice()
		expected := []int{0, 1, 2, 3}
		if len(slice) != len(expected) {
			t.Fatalf("ToSlice() has %d items, want %d", len(slice), len(expected))
		}
		for i := range slice {
			if slice[i] != expected[i] {
				t.Errorf("ToSlice()[%d] = %d, want %d", i, slice[i], expected[i])
			}
		}
	})
}

// =============================================================================
// Lesson Function Tests
// =============================================================================

func TestLessonGenerics(t *testing.T) {
	t.Run("Min and Max", func(t *testing.T) {
		if Min(3, 5) != 3 {
			t.Error("Min(3, 5) should be 3")
		}
		if Max(3, 5) != 5 {
			t.Error("Max(3, 5) should be 5")
		}
		if Min("apple", "banana") != "apple" {
			t.Error("Min(apple, banana) should be apple")
		}
	})

	t.Run("Contains", func(t *testing.T) {
		if !Contains([]int{1, 2, 3}, 2) {
			t.Error("Contains should find 2")
		}
		if Contains([]int{1, 2, 3}, 4) {
			t.Error("Contains should not find 4")
		}
	})

	t.Run("Map Filter Reduce", func(t *testing.T) {
		doubled := Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
		if len(doubled) != 3 || doubled[0] != 2 || doubled[1] != 4 || doubled[2] != 6 {
			t.Errorf("Map double = %v, want [2 4 6]", doubled)
		}

		evens := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
		if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 {
			t.Errorf("Filter even = %v, want [2 4]", evens)
		}

		sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
		if sum != 10 {
			t.Errorf("Reduce sum = %d, want 10", sum)
		}
	})

	t.Run("Stack", func(t *testing.T) {
		s := NewStack[int]()
		s.Push(1)
		s.Push(2)
		v, ok := s.Pop()
		if !ok || v != 2 {
			t.Errorf("Pop() = (%d, %v), want (2, true)", v, ok)
		}
	})

	t.Run("Set", func(t *testing.T) {
		s := SetFrom([]string{"a", "b", "c"})
		if !s.Contains("a") {
			t.Error("set should contain 'a'")
		}
		if s.Len() != 3 {
			t.Errorf("set Len() = %d, want 3", s.Len())
		}
	})

	t.Run("Result type", func(t *testing.T) {
		ok := Ok(42)
		if !ok.IsOk() || ok.Unwrap() != 42 {
			t.Error("Ok(42) should contain 42")
		}

		err := Err[int](fmt.Errorf("failed"))
		if !err.IsErr() {
			t.Error("Err should be an error")
		}
		if err.UnwrapOr(99) != 99 {
			t.Error("UnwrapOr should return default for Err")
		}
	})

	t.Run("Sum with named types", func(t *testing.T) {
		if Sum([]int{1, 2, 3}) != 6 {
			t.Error("Sum(1,2,3) should be 6")
		}
		ids := []UserID{UserID(10), UserID(20)}
		if Sum(ids) != 30 {
			t.Error("Sum(UserID 10, 20) should be 30")
		}
	})
}
