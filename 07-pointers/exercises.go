package pointers

import "fmt"

// =========================================================================
// Exercise 1: Swap Function Using Pointers
// =========================================================================

// Swap swaps the values of two integers using pointers.
// After Swap(&a, &b), a should have b's original value and vice versa.
func Swap(a, b *int) {
	// YOUR CODE HERE
}

// SwapStrings swaps two string values using pointers.
func SwapStrings(a, b *string) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 2: Modify a Struct Through a Pointer
// =========================================================================

// Player represents a game player with stats.
type Player struct {
	Name   string
	Health int
	Score  int
	Level  int
}

// Heal increases the player's health by the given amount, up to a max of 100.
// If the player is nil, do nothing (don't panic).
func Heal(p *Player, amount int) {
	// YOUR CODE HERE
}

// TakeDamage decreases the player's health by the given amount, minimum 0.
// If the player is nil, do nothing.
func TakeDamage(p *Player, amount int) {
	// YOUR CODE HERE
}

// LevelUp increments the player's level and adds 10 to their score.
// If the player is nil, do nothing.
func LevelUp(p *Player) {
	// YOUR CODE HERE
}

// ResetPlayer resets a player to starting stats: Health=100, Score=0, Level=1.
// Name stays the same. If the player is nil, do nothing.
func ResetPlayer(p *Player) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 3: Linked List Operations Using Pointers
// =========================================================================

// Node represents a node in a doubly linked list.
type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

// DoublyLinkedList represents a doubly linked list with head and tail pointers.
type DoublyLinkedList struct {
	Head *Node
	Tail *Node
	Len  int
}

// NewDoublyLinkedList creates an empty doubly linked list.
func NewDoublyLinkedList() *DoublyLinkedList {
	// YOUR CODE HERE
	return nil
}

// PushBack adds a value to the end of the list.
func (dl *DoublyLinkedList) PushBack(value int) {
	// YOUR CODE HERE
}

// PushFront adds a value to the beginning of the list.
func (dl *DoublyLinkedList) PushFront(value int) {
	// YOUR CODE HERE
}

// PopFront removes and returns the first value from the list.
// Returns the value and true if successful, or 0 and false if the list is empty.
func (dl *DoublyLinkedList) PopFront() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// PopBack removes and returns the last value from the list.
// Returns the value and true if successful, or 0 and false if the list is empty.
func (dl *DoublyLinkedList) PopBack() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// ToSlice converts the list to a slice (front to back).
func (dl *DoublyLinkedList) ToSlice() []int {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 4: Optional Fields Pattern
// =========================================================================

// Profile represents a user profile where fields are optional.
// nil means "not set" — different from empty string or zero.
type Profile struct {
	DisplayName *string
	Bio         *string
	Age         *int
	Website     *string
}

// StringPtr returns a pointer to the given string.
// This is a helper function you'll see in many Go codebases.
func StringPtr(s string) *string {
	// YOUR CODE HERE
	return nil
}

// IntPtr returns a pointer to the given int.
func IntPtr(n int) *int {
	// YOUR CODE HERE
	return nil
}

// GetDisplayName returns the display name if set, or "Anonymous" if nil.
func (p *Profile) GetDisplayName() string {
	// YOUR CODE HERE
	return ""
}

// GetBio returns the bio if set, or "No bio provided" if nil.
func (p *Profile) GetBio() string {
	// YOUR CODE HERE
	return ""
}

// GetAge returns the age as a string if set, or "Not specified" if nil.
// Format the age as a plain number string, e.g., "25".
func (p *Profile) GetAge() string {
	// YOUR CODE HERE
	return ""
}

// SetFields updates only the non-nil fields from the update profile.
// If a field in the update is nil, leave the original unchanged.
// This is exactly how JSON PATCH operations work in APIs.
func (p *Profile) SetFields(update Profile) {
	// YOUR CODE HERE
}

// =========================================================================
// Exercise 5: Binary Search Tree
// =========================================================================

// BST represents a binary search tree node.
type BST struct {
	Value int
	Left  *BST
	Right *BST
}

// Insert adds a value to the BST. If the value already exists, do nothing.
// Return the (possibly new) root of the tree.
// - Values less than current go left
// - Values greater than current go right
//
// Why does this return *BST? Because if the tree is initially nil (empty),
// we need to return the new root node. The caller pattern is:
//
//	root = root.Insert(5)
func (t *BST) Insert(value int) *BST {
	// YOUR CODE HERE
	return nil
}

// Search returns true if the value exists in the BST.
func (t *BST) Search(value int) bool {
	// YOUR CODE HERE
	return false
}

// InOrder returns the values of the BST in sorted order (in-order traversal).
// If the tree is nil, return an empty (non-nil) slice.
func (t *BST) InOrder() []int {
	// YOUR CODE HERE
	return nil
}

// Min returns the minimum value in the BST and true, or 0 and false if empty.
func (t *BST) Min() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// Max returns the maximum value in the BST and true, or 0 and false if empty.
func (t *BST) Max() (int, bool) {
	// YOUR CODE HERE
	return 0, false
}

// =========================================================================
// Exercise 6: Nil Receiver Pattern
// =========================================================================

// Logger represents a simple logger. If nil, logging is silently disabled.
// This is a common Go pattern — nil means "no-op."
type Logger struct {
	Prefix   string
	Messages []string
}

// Log appends a message to the logger's Messages slice.
// If the logger is nil, do nothing (silently discard the message).
// Format: "[prefix] message"
func (l *Logger) Log(message string) {
	// YOUR CODE HERE
}

// LastMessage returns the most recent log message.
// If the logger is nil or has no messages, return "".
func (l *Logger) LastMessage() string {
	// YOUR CODE HERE
	return ""
}

// Count returns the number of messages logged.
// If the logger is nil, return 0.
func (l *Logger) Count() int {
	// YOUR CODE HERE
	return 0
}

// =========================================================================
// Exercise 7: Pointer Aliasing
// =========================================================================

// Counter holds a count value.
type Counter struct {
	count int
}

// Increment adds 1 to the counter.
func (c *Counter) Increment() {
	// YOUR CODE HERE
}

// Value returns the current count.
func (c *Counter) Value() int {
	// YOUR CODE HERE
	return 0
}

// ShareCounter takes a Counter and returns two pointers that both
// point to the SAME counter. Incrementing through either pointer
// should affect the other.
//
// This demonstrates pointer aliasing — two variables referencing the
// same underlying data. This is important to understand because it
// happens naturally with slices, maps, and any shared pointer.
func ShareCounter(c *Counter) (*Counter, *Counter) {
	// YOUR CODE HERE
	return nil, nil
}

// CopyCounter takes a Counter pointer and returns a pointer to a NEW
// Counter with the same count value. The original and copy should be
// independent — changing one doesn't affect the other.
func CopyCounter(c *Counter) *Counter {
	// YOUR CODE HERE
	return nil
}

// =========================================================================
// Exercise 8: Reference-Counted Cache
// =========================================================================

// CacheEntry holds a value and tracks how many times it's been accessed.
type CacheEntry struct {
	Key      string
	Value    string
	RefCount int
}

// Cache is a simple reference-counted cache. Each time an entry is
// accessed (via Get), its RefCount increases. Entries can be evicted
// based on their reference count.
type Cache struct {
	entries map[string]*CacheEntry
}

// NewCache creates a new empty Cache.
func NewCache() *Cache {
	// YOUR CODE HERE
	return nil
}

// Put adds or updates a cache entry. If the key already exists,
// update the value but keep the existing RefCount.
// New entries start with RefCount = 0.
func (c *Cache) Put(key, value string) {
	// YOUR CODE HERE
}

// Get retrieves a cache entry by key. If found, increment its RefCount
// and return the entry pointer and true. If not found, return nil and false.
//
// Returning a pointer to the entry is intentional — it lets callers see
// the live RefCount.
func (c *Cache) Get(key string) (*CacheEntry, bool) {
	// YOUR CODE HERE
	return nil, false
}

// EvictLeastUsed removes the entry with the lowest RefCount from the cache.
// If there are ties, remove any one of the tied entries.
// If the cache is empty, do nothing.
func (c *Cache) EvictLeastUsed() {
	// YOUR CODE HERE
}

// Size returns the number of entries in the cache.
func (c *Cache) Size() int {
	// YOUR CODE HERE
	return 0
}

// Entries returns all cache entries as a slice (order not guaranteed).
func (c *Cache) Entries() []*CacheEntry {
	// YOUR CODE HERE
	return nil
}

// Ensure unused import doesn't cause compile error
var _ = fmt.Sprintf
