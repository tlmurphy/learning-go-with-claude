package pointers

import (
	"fmt"
	"testing"
)

// =========================================================================
// Exercise 1 Tests: Swap
// =========================================================================

func TestSwap(t *testing.T) {
	tests := []struct {
		a, b     int
		wantA, wantB int
	}{
		{1, 2, 2, 1},
		{0, 0, 0, 0},
		{-1, 1, 1, -1},
		{42, 99, 99, 42},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_%d", tt.a, tt.b), func(t *testing.T) {
			a, b := tt.a, tt.b
			Swap(&a, &b)
			if a != tt.wantA || b != tt.wantB {
				t.Errorf("After Swap(%d, %d): got a=%d, b=%d; want a=%d, b=%d",
					tt.a, tt.b, a, b, tt.wantA, tt.wantB)
			}
		})
	}
}

func TestSwapStrings(t *testing.T) {
	a, b := "hello", "world"
	SwapStrings(&a, &b)
	if a != "world" || b != "hello" {
		t.Errorf("After SwapStrings: a=%q, b=%q; want a=%q, b=%q",
			a, b, "world", "hello")
	}
}

// =========================================================================
// Exercise 2 Tests: Player Struct Modification
// =========================================================================

func TestHeal(t *testing.T) {
	tests := []struct {
		name       string
		health     int
		amount     int
		wantHealth int
	}{
		{"normal heal", 50, 20, 70},
		{"heal to max", 90, 20, 100},
		{"over heal capped at 100", 80, 50, 100},
		{"already full", 100, 10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Name: "Test", Health: tt.health}
			Heal(p, tt.amount)
			if p.Health != tt.wantHealth {
				t.Errorf("After Heal(%d, %d): Health = %d, want %d",
					tt.health, tt.amount, p.Health, tt.wantHealth)
			}
		})
	}

	t.Run("nil player", func(t *testing.T) {
		// Should not panic
		Heal(nil, 10)
	})
}

func TestTakeDamage(t *testing.T) {
	tests := []struct {
		name       string
		health     int
		amount     int
		wantHealth int
	}{
		{"normal damage", 100, 30, 70},
		{"lethal damage", 50, 60, 0},
		{"exact kill", 50, 50, 0},
		{"overkill capped at 0", 10, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Player{Name: "Test", Health: tt.health}
			TakeDamage(p, tt.amount)
			if p.Health != tt.wantHealth {
				t.Errorf("After TakeDamage(%d, %d): Health = %d, want %d",
					tt.health, tt.amount, p.Health, tt.wantHealth)
			}
		})
	}

	t.Run("nil player", func(t *testing.T) {
		TakeDamage(nil, 10)
	})
}

func TestLevelUp(t *testing.T) {
	p := &Player{Name: "Hero", Health: 100, Score: 50, Level: 3}
	LevelUp(p)

	if p.Level != 4 {
		t.Errorf("Level = %d, want 4", p.Level)
	}
	if p.Score != 60 {
		t.Errorf("Score = %d, want 60 (should add 10)", p.Score)
	}

	// Nil safety
	LevelUp(nil)
}

func TestResetPlayer(t *testing.T) {
	p := &Player{Name: "Hero", Health: 30, Score: 500, Level: 10}
	ResetPlayer(p)

	if p.Name != "Hero" {
		t.Errorf("Name = %q, want %q (name should not change)", p.Name, "Hero")
	}
	if p.Health != 100 {
		t.Errorf("Health = %d, want 100", p.Health)
	}
	if p.Score != 0 {
		t.Errorf("Score = %d, want 0", p.Score)
	}
	if p.Level != 1 {
		t.Errorf("Level = %d, want 1", p.Level)
	}

	ResetPlayer(nil)
}

// =========================================================================
// Exercise 3 Tests: Doubly Linked List
// =========================================================================

func TestDoublyLinkedListPushBack(t *testing.T) {
	dl := NewDoublyLinkedList()
	if dl == nil {
		t.Fatal("NewDoublyLinkedList returned nil")
	}

	dl.PushBack(1)
	dl.PushBack(2)
	dl.PushBack(3)

	got := dl.ToSlice()
	want := []int{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("ToSlice() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ToSlice()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
	if dl.Len != 3 {
		t.Errorf("Len = %d, want 3", dl.Len)
	}
}

func TestDoublyLinkedListPushFront(t *testing.T) {
	dl := NewDoublyLinkedList()
	if dl == nil {
		t.Fatal("NewDoublyLinkedList returned nil")
	}

	dl.PushFront(3)
	dl.PushFront(2)
	dl.PushFront(1)

	got := dl.ToSlice()
	want := []int{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("ToSlice() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ToSlice()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestDoublyLinkedListPopFront(t *testing.T) {
	dl := NewDoublyLinkedList()
	if dl == nil {
		t.Fatal("NewDoublyLinkedList returned nil")
	}

	dl.PushBack(1)
	dl.PushBack(2)
	dl.PushBack(3)

	val, ok := dl.PopFront()
	if !ok || val != 1 {
		t.Errorf("PopFront() = (%d, %v), want (1, true)", val, ok)
	}
	if dl.Len != 2 {
		t.Errorf("Len after PopFront = %d, want 2", dl.Len)
	}

	val, ok = dl.PopFront()
	if !ok || val != 2 {
		t.Errorf("PopFront() = (%d, %v), want (2, true)", val, ok)
	}

	val, ok = dl.PopFront()
	if !ok || val != 3 {
		t.Errorf("PopFront() = (%d, %v), want (3, true)", val, ok)
	}

	// Empty list
	_, ok = dl.PopFront()
	if ok {
		t.Error("PopFront on empty list should return ok=false")
	}
}

func TestDoublyLinkedListPopBack(t *testing.T) {
	dl := NewDoublyLinkedList()
	if dl == nil {
		t.Fatal("NewDoublyLinkedList returned nil")
	}

	dl.PushBack(1)
	dl.PushBack(2)
	dl.PushBack(3)

	val, ok := dl.PopBack()
	if !ok || val != 3 {
		t.Errorf("PopBack() = (%d, %v), want (3, true)", val, ok)
	}

	val, ok = dl.PopBack()
	if !ok || val != 2 {
		t.Errorf("PopBack() = (%d, %v), want (2, true)", val, ok)
	}

	val, ok = dl.PopBack()
	if !ok || val != 1 {
		t.Errorf("PopBack() = (%d, %v), want (1, true)", val, ok)
	}

	_, ok = dl.PopBack()
	if ok {
		t.Error("PopBack on empty list should return ok=false")
	}
}

func TestDoublyLinkedListPrevLinks(t *testing.T) {
	dl := NewDoublyLinkedList()
	if dl == nil {
		t.Fatal("NewDoublyLinkedList returned nil")
	}

	dl.PushBack(1)
	dl.PushBack(2)
	dl.PushBack(3)

	// Traverse backwards from tail to verify Prev links
	node := dl.Tail
	if node == nil {
		t.Fatal("Tail is nil after PushBack")
	}
	var backward []int
	for node != nil {
		backward = append(backward, node.Value)
		node = node.Prev
	}
	want := []int{3, 2, 1}
	if len(backward) != len(want) {
		t.Fatalf("Backward traversal length = %d, want %d", len(backward), len(want))
	}
	for i := range want {
		if backward[i] != want[i] {
			t.Errorf("Backward[%d] = %d, want %d", i, backward[i], want[i])
		}
	}
}

// =========================================================================
// Exercise 4 Tests: Optional Fields
// =========================================================================

func TestStringPtr(t *testing.T) {
	p := StringPtr("hello")
	if p == nil {
		t.Fatal("StringPtr returned nil")
	}
	if *p != "hello" {
		t.Errorf("*StringPtr(\"hello\") = %q, want %q", *p, "hello")
	}
}

func TestIntPtr(t *testing.T) {
	p := IntPtr(42)
	if p == nil {
		t.Fatal("IntPtr returned nil")
	}
	if *p != 42 {
		t.Errorf("*IntPtr(42) = %d, want 42", *p)
	}
}

func TestProfileGetters(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		p := &Profile{}
		if p.GetDisplayName() != "Anonymous" {
			t.Errorf("GetDisplayName() = %q, want %q", p.GetDisplayName(), "Anonymous")
		}
		if p.GetBio() != "No bio provided" {
			t.Errorf("GetBio() = %q, want %q", p.GetBio(), "No bio provided")
		}
		if p.GetAge() != "Not specified" {
			t.Errorf("GetAge() = %q, want %q", p.GetAge(), "Not specified")
		}
	})

	t.Run("all set", func(t *testing.T) {
		p := &Profile{
			DisplayName: StringPtr("Alice"),
			Bio:         StringPtr("Go developer"),
			Age:         IntPtr(30),
		}
		if p.GetDisplayName() != "Alice" {
			t.Errorf("GetDisplayName() = %q, want %q", p.GetDisplayName(), "Alice")
		}
		if p.GetBio() != "Go developer" {
			t.Errorf("GetBio() = %q, want %q", p.GetBio(), "Go developer")
		}
		if p.GetAge() != "30" {
			t.Errorf("GetAge() = %q, want %q", p.GetAge(), "30")
		}
	})

	t.Run("empty strings are valid", func(t *testing.T) {
		p := &Profile{
			DisplayName: StringPtr(""),
		}
		if p.GetDisplayName() != "" {
			t.Errorf("GetDisplayName() = %q, want empty string (not default)", p.GetDisplayName())
		}
	})
}

func TestProfileSetFields(t *testing.T) {
	namePtr := StringPtr("Alice")
	bioPtr := StringPtr("Original bio")
	agePtr := IntPtr(25)
	if namePtr == nil || bioPtr == nil || agePtr == nil {
		t.Fatal("StringPtr/IntPtr returned nil — implement these helpers first")
	}

	p := &Profile{
		DisplayName: namePtr,
		Bio:         bioPtr,
		Age:         agePtr,
	}

	updatedBio := StringPtr("Updated bio")
	if updatedBio == nil {
		t.Fatal("StringPtr returned nil")
	}

	// Only update Bio — leave everything else alone
	p.SetFields(Profile{
		Bio: updatedBio,
	})

	if p.DisplayName == nil || *p.DisplayName != "Alice" {
		t.Errorf("DisplayName should be unchanged as \"Alice\"")
	}
	if p.Bio == nil || *p.Bio != "Updated bio" {
		t.Errorf("Bio should be \"Updated bio\"")
	}
	if p.Age == nil || *p.Age != 25 {
		t.Errorf("Age should be unchanged as 25")
	}
}

// =========================================================================
// Exercise 5 Tests: BST
// =========================================================================

func TestBSTInsertAndSearch(t *testing.T) {
	var root *BST
	root = root.Insert(5)
	root = root.Insert(3)
	root = root.Insert(7)
	root = root.Insert(1)
	root = root.Insert(4)

	if root == nil {
		t.Fatal("BST root is nil after Insert")
	}

	tests := []struct {
		value int
		want  bool
	}{
		{5, true},
		{3, true},
		{7, true},
		{1, true},
		{4, true},
		{2, false},
		{6, false},
		{0, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("search_%d", tt.value), func(t *testing.T) {
			got := root.Search(tt.value)
			if got != tt.want {
				t.Errorf("Search(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestBSTSearchNil(t *testing.T) {
	var root *BST
	if root.Search(42) {
		t.Error("Search on nil tree should return false")
	}
}

func TestBSTInOrder(t *testing.T) {
	var root *BST
	for _, v := range []int{5, 3, 7, 1, 4, 6, 8} {
		root = root.Insert(v)
	}

	got := root.InOrder()
	want := []int{1, 3, 4, 5, 6, 7, 8}
	if len(got) != len(want) {
		t.Fatalf("InOrder() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("InOrder()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestBSTInOrderNil(t *testing.T) {
	var root *BST
	got := root.InOrder()
	if got == nil {
		t.Error("InOrder() on nil tree should return empty (non-nil) slice")
	}
	if len(got) != 0 {
		t.Errorf("InOrder() on nil tree length = %d, want 0", len(got))
	}
}

func TestBSTMinMax(t *testing.T) {
	var root *BST
	for _, v := range []int{5, 3, 7, 1, 4, 6, 8} {
		root = root.Insert(v)
	}

	t.Run("min", func(t *testing.T) {
		val, ok := root.Min()
		if !ok || val != 1 {
			t.Errorf("Min() = (%d, %v), want (1, true)", val, ok)
		}
	})

	t.Run("max", func(t *testing.T) {
		val, ok := root.Max()
		if !ok || val != 8 {
			t.Errorf("Max() = (%d, %v), want (8, true)", val, ok)
		}
	})

	t.Run("nil tree", func(t *testing.T) {
		var empty *BST
		_, ok := empty.Min()
		if ok {
			t.Error("Min() on nil tree should return ok=false")
		}
		_, ok = empty.Max()
		if ok {
			t.Error("Max() on nil tree should return ok=false")
		}
	})
}

func TestBSTDuplicates(t *testing.T) {
	var root *BST
	root = root.Insert(5)
	root = root.Insert(5) // duplicate — should be ignored
	root = root.Insert(5) // duplicate — should be ignored

	got := root.InOrder()
	if len(got) != 1 {
		t.Errorf("InOrder() after duplicate inserts length = %d, want 1", len(got))
	}
}

// =========================================================================
// Exercise 6 Tests: Nil Receiver Logger
// =========================================================================

func TestLogger(t *testing.T) {
	l := &Logger{Prefix: "APP"}
	l.Log("starting up")
	l.Log("ready")

	if l.Count() != 2 {
		t.Errorf("Count() = %d, want 2", l.Count())
	}

	last := l.LastMessage()
	want := "[APP] ready"
	if last != want {
		t.Errorf("LastMessage() = %q, want %q", last, want)
	}
}

func TestLoggerNilReceiver(t *testing.T) {
	var l *Logger

	// None of these should panic
	l.Log("this should not panic")

	if l.Count() != 0 {
		t.Errorf("nil Logger.Count() = %d, want 0", l.Count())
	}
	if l.LastMessage() != "" {
		t.Errorf("nil Logger.LastMessage() = %q, want empty", l.LastMessage())
	}
}

// =========================================================================
// Exercise 7 Tests: Pointer Aliasing
// =========================================================================

func TestShareCounter(t *testing.T) {
	c := &Counter{count: 0}
	a, b := ShareCounter(c)

	if a == nil || b == nil {
		t.Fatal("ShareCounter returned nil")
	}

	a.Increment()
	a.Increment()
	b.Increment()

	// All three should show the same count because they're aliases
	if c.Value() != 3 {
		t.Errorf("Original counter = %d, want 3", c.Value())
	}
	if a.Value() != 3 {
		t.Errorf("Alias a = %d, want 3", a.Value())
	}
	if b.Value() != 3 {
		t.Errorf("Alias b = %d, want 3 — both aliases should share the same counter", b.Value())
	}
}

func TestCopyCounter(t *testing.T) {
	c := &Counter{count: 5}
	copy := CopyCounter(c)

	if copy == nil {
		t.Fatal("CopyCounter returned nil")
	}

	if copy.Value() != 5 {
		t.Errorf("Copy value = %d, want 5", copy.Value())
	}

	// Modifying the copy should NOT affect the original
	copy.Increment()
	copy.Increment()

	if c.Value() != 5 {
		t.Errorf("Original value = %d, want 5 (copy should be independent)", c.Value())
	}
	if copy.Value() != 7 {
		t.Errorf("Copy value = %d, want 7", copy.Value())
	}
}

// =========================================================================
// Exercise 8 Tests: Reference-Counted Cache
// =========================================================================

func TestCachePutAndGet(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("NewCache returned nil")
	}

	c.Put("key1", "value1")
	c.Put("key2", "value2")

	t.Run("get existing", func(t *testing.T) {
		entry, ok := c.Get("key1")
		if !ok {
			t.Fatal("Get(\"key1\") returned ok=false")
		}
		if entry.Value != "value1" {
			t.Errorf("Value = %q, want %q", entry.Value, "value1")
		}
		if entry.RefCount != 1 {
			t.Errorf("RefCount = %d, want 1 (first Get)", entry.RefCount)
		}
	})

	t.Run("get increments refcount", func(t *testing.T) {
		c.Get("key1")
		c.Get("key1")
		entry, _ := c.Get("key1")
		if entry.RefCount != 4 {
			t.Errorf("RefCount after 4 Gets = %d, want 4", entry.RefCount)
		}
	})

	t.Run("get missing", func(t *testing.T) {
		_, ok := c.Get("nonexistent")
		if ok {
			t.Error("Get(\"nonexistent\") should return ok=false")
		}
	})
}

func TestCachePutUpdate(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("NewCache returned nil")
	}

	c.Put("key1", "original")
	c.Get("key1") // RefCount = 1
	c.Get("key1") // RefCount = 2

	c.Put("key1", "updated")
	entry, ok := c.Get("key1")
	if !ok {
		t.Fatal("Get after update returned ok=false")
	}
	if entry.Value != "updated" {
		t.Errorf("Value = %q, want %q", entry.Value, "updated")
	}
	if entry.RefCount != 3 {
		t.Errorf("RefCount = %d, want 3 (should keep old count + this Get)", entry.RefCount)
	}
}

func TestCacheEvictLeastUsed(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("NewCache returned nil")
	}

	c.Put("popular", "data1")
	c.Put("unpopular", "data2")

	// Access "popular" several times
	c.Get("popular")
	c.Get("popular")
	c.Get("popular")

	// Access "unpopular" once
	c.Get("unpopular")

	c.EvictLeastUsed()

	if c.Size() != 1 {
		t.Errorf("Size after eviction = %d, want 1", c.Size())
	}

	_, ok := c.Get("unpopular")
	if ok {
		t.Error("\"unpopular\" should have been evicted (lowest RefCount)")
	}

	_, ok = c.Get("popular")
	if !ok {
		t.Error("\"popular\" should still be in cache")
	}
}

func TestCacheSize(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("NewCache returned nil")
	}

	if c.Size() != 0 {
		t.Errorf("Empty cache Size() = %d, want 0", c.Size())
	}

	c.Put("a", "1")
	c.Put("b", "2")
	if c.Size() != 2 {
		t.Errorf("Size() = %d, want 2", c.Size())
	}
}

func TestCacheEvictEmpty(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("NewCache returned nil")
	}

	// Should not panic
	c.EvictLeastUsed()
}
