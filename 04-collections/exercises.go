package collections

/*
=============================================================================
 EXERCISES: Collections
=============================================================================

 Work through these exercises in order. Run tests with:

   go test -v ./04-collections/

 Run a single test:

   go test -v -run TestSliceOps ./04-collections/

=============================================================================
*/

// Exercise 1: Slice Operations
//
// Implement basic slice manipulation functions. These operations don't
// have built-in functions in Go, so you need to know how to do them
// with append and slicing.

// Prepend adds an element to the beginning of a slice and returns the new slice.
// Example: Prepend([]int{2, 3, 4}, 1) returns [1, 2, 3, 4]
func Prepend(s []int, val int) []int {
	// YOUR CODE HERE
	return nil
}

// RemoveAt removes the element at the given index and returns the new slice.
// If the index is out of bounds, return the original slice unchanged.
// The order of remaining elements must be preserved.
// Example: RemoveAt([]int{10, 20, 30, 40}, 1) returns [10, 30, 40]
func RemoveAt(s []int, index int) []int {
	// YOUR CODE HERE
	return nil
}

// InsertAt inserts a value at the given index and returns the new slice.
// If index == len(s), append to the end. If index is out of bounds, return
// the original slice unchanged.
// Example: InsertAt([]int{10, 30, 40}, 1, 20) returns [10, 20, 30, 40]
func InsertAt(s []int, index int, val int) []int {
	// YOUR CODE HERE
	return nil
}

// Exercise 2: Capacity Detective
//
// This exercise tests your understanding of slice capacity and growth.

// PredictCapacity creates a slice using make with the given length and capacity,
// then appends n elements. Return the final length and capacity.
//
// This helps you understand when append allocates a new backing array.
//
// Example: PredictCapacity(3, 5, 2) creates make([]int, 3, 5), appends 2 elements.
//
//	Final: len=5, cap=5 (fits in existing capacity)
//
// Example: PredictCapacity(3, 5, 4) creates make([]int, 3, 5), appends 4 elements.
//
//	Final: len=7, cap=10 (exceeded capacity, had to grow)
func PredictCapacity(initialLen, initialCap, appendCount int) (finalLen, finalCap int) {
	// YOUR CODE HERE
	return 0, 0
}

// Exercise 3: WordFrequency
//
// Count the frequency of each word in a slice of strings.
// Return a map of word -> count.
// Words are case-sensitive (don't normalize).
//
// Example: WordFrequency([]string{"go", "is", "go"}) returns map[go:2, is:1]
//
// This is one of the most common map patterns. In web services, you'd use
// this for analytics, counting API hits per endpoint, etc.
func WordFrequency(words []string) map[string]int {
	// YOUR CODE HERE
	return nil
}

// Exercise 4: Set Operations
//
// Implement a set using map[string]struct{} — the idiomatic Go set pattern.
// struct{} takes zero bytes, making it memory-efficient.

// NewStringSet creates a set from the given strings (deduplicated).
func NewStringSet(items []string) map[string]struct{} {
	// YOUR CODE HERE
	return nil
}

// SetContains returns true if the set contains the item.
func SetContains(set map[string]struct{}, item string) bool {
	// YOUR CODE HERE
	return false
}

// SetUnion returns a new set containing all elements from both sets.
func SetUnion(a, b map[string]struct{}) map[string]struct{} {
	// YOUR CODE HERE
	return nil
}

// SetIntersection returns a new set containing only elements present in both sets.
func SetIntersection(a, b map[string]struct{}) map[string]struct{} {
	// YOUR CODE HERE
	return nil
}

// SetDifference returns a new set containing elements in a but not in b.
func SetDifference(a, b map[string]struct{}) map[string]struct{} {
	// YOUR CODE HERE
	return nil
}

// Exercise 5: Matrix Operations
//
// Work with slices of slices (2D slices).

// NewMatrix creates a rows x cols matrix initialized to zero.
// Return a [][]int where each inner slice has length cols.
func NewMatrix(rows, cols int) [][]int {
	// YOUR CODE HERE
	return nil
}

// MatrixTranspose returns the transpose of a matrix.
// The transpose flips rows and columns: element [i][j] becomes [j][i].
// If the matrix is empty, return an empty matrix.
func MatrixTranspose(matrix [][]int) [][]int {
	// YOUR CODE HERE
	return nil
}

// Exercise 6: Deduplicate
//
// Remove duplicate values from a slice while preserving the order of first
// occurrence.
//
// Example: Deduplicate([]int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}) returns [3, 1, 4, 5, 9, 2, 6]
//
// Hint: Use a map to track which values you've already seen.
func Deduplicate(s []int) []int {
	// YOUR CODE HERE
	return nil
}

// Exercise 7: GroupBy
//
// Group a slice of strings by their first character.
// Return a map where keys are the first rune of each string, and values
// are slices of strings that start with that character.
//
// Empty strings should be ignored.
// The grouping should be case-sensitive.
//
// Example: GroupBy([]string{"apple", "avocado", "banana", "blueberry", "cherry"})
//
//	returns: map['a':["apple","avocado"], 'b':["banana","blueberry"], 'c':["cherry"]]
func GroupBy(items []string) map[rune][]string {
	// YOUR CODE HERE
	return nil
}

// Exercise 8: Stack and Queue
//
// Implement a simple stack (LIFO) and queue (FIFO) using slices.
// These are returned as structs with methods, previewing the next module.

// IntStack is a LIFO stack of integers.
type IntStack struct {
	data []int
}

// Push adds a value to the top of the stack.
func (s *IntStack) Push(val int) {
	// YOUR CODE HERE
}

// Pop removes and returns the top value from the stack.
// Returns 0 and false if the stack is empty.
func (s *IntStack) Pop() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// Peek returns the top value without removing it.
// Returns 0 and false if the stack is empty.
func (s *IntStack) Peek() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// Len returns the number of elements in the stack.
func (s *IntStack) Len() int {
	// YOUR CODE HERE
	return 0
}

// IntQueue is a FIFO queue of integers.
type IntQueue struct {
	data []int
}

// Enqueue adds a value to the back of the queue.
func (q *IntQueue) Enqueue(val int) {
	// YOUR CODE HERE
}

// Dequeue removes and returns the front value from the queue.
// Returns 0 and false if the queue is empty.
func (q *IntQueue) Dequeue() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// Len returns the number of elements in the queue.
func (q *IntQueue) Len() int {
	// YOUR CODE HERE
	return 0
}

// Exercise 9: MergeSorted
//
// Merge two sorted (ascending) integer slices into a single sorted slice.
// The input slices are guaranteed to be sorted. Do NOT sort the result —
// use the merge algorithm (similar to merge sort's merge step).
//
// This is an important algorithm to know. It runs in O(n+m) time, which
// is better than concatenating and sorting O((n+m) * log(n+m)).
//
// Example: MergeSorted([]int{1, 3, 5}, []int{2, 4, 6}) returns [1, 2, 3, 4, 5, 6]
// Example: MergeSorted([]int{1, 2, 3}, []int{}) returns [1, 2, 3]
func MergeSorted(a, b []int) []int {
	// YOUR CODE HERE
	return nil
}
