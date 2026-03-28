// Package memory explores Go's memory model, garbage collector, and techniques
// for writing allocation-efficient code in production systems.
package memory

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

/*
=============================================================================
 MEMORY AND GARBAGE COLLECTION IN GO
=============================================================================

Understanding memory in Go is the difference between a service that handles
10,000 requests/second and one that handles 100,000. Most Go code doesn't
need to worry about memory — the GC is excellent. But when you're in the
hot path of a high-throughput service, every allocation counts.

The key insight: in Go, the COMPILER decides whether memory lives on the
stack or the heap. You don't control this directly like in C/C++. But you
can influence it by understanding what causes "escape to heap."

=============================================================================
 STACK VS HEAP
=============================================================================

Stack memory is fast and free (in terms of GC cost):
  - Allocated and freed automatically when functions return
  - Each goroutine has its own stack (starts at ~2-8 KB, grows as needed)
  - No garbage collection needed — just move the stack pointer

Heap memory is where the GC comes in:
  - Shared across all goroutines
  - Must be garbage collected when no longer referenced
  - Allocating on the heap means the GC has more work to do

Rule of thumb: if the compiler can PROVE a variable doesn't outlive its
function, it goes on the stack. Otherwise, it "escapes" to the heap.

=============================================================================
 ESCAPE ANALYSIS
=============================================================================

Run: go build -gcflags="-m" ./yourpackage/
This shows you exactly what escapes to the heap and why.

Common causes of heap allocation:
  1. Returning a pointer to a local variable
  2. Storing a value in an interface (boxing)
  3. Closures that capture local variables
  4. Sending a pointer to a channel
  5. Slices that grow beyond their initial capacity
  6. Maps (always heap-allocated)
  7. Values too large for the stack

Example:
  func newUser() *User { return &User{Name: "Alice"} }
  // ^^ User escapes to heap because we return a pointer to it.
  // The pointer might outlive the function, so Go must heap-allocate.

  func process(u User) { ... }
  // ^^ User stays on stack if process() doesn't cause it to escape.
  // Passing by value is often CHEAPER than passing by pointer for small structs
  // because it avoids heap allocation and GC pressure.

=============================================================================
 THE GARBAGE COLLECTOR
=============================================================================

Go's GC is a concurrent, tri-color mark-and-sweep collector:

  1. Mark phase: Starting from "roots" (goroutine stacks, globals), the GC
     walks all reachable objects. It uses three colors:
     - White: not yet visited (candidates for collection)
     - Gray: visited but children not yet scanned
     - Black: visited and all children scanned

  2. Sweep phase: White objects (unreachable) are freed.

The magic: this runs CONCURRENTLY with your program. Go's GC prioritizes
low latency over throughput — typical pause times are under 1ms, even with
large heaps. This is why Go is so popular for web services.

Go 1.26 defaults to the new "Green Tea" GC with 10-40% less overhead. It
refines the pacer (which decides when to trigger GC) and reduces the amount
of scanning work needed per cycle.

=============================================================================
 GOGC AND GOMEMLIMIT
=============================================================================

GOGC controls how aggressively the GC runs. Default is 100, meaning the GC
triggers when heap size doubles since last collection.
  - GOGC=50  → GC runs more often, uses less memory, more CPU on GC
  - GOGC=200 → GC runs less often, uses more memory, less CPU on GC
  - GOGC=off → disables GC entirely (useful for short-lived CLI tools)

GOMEMLIMIT (Go 1.19+) sets a soft memory limit. The GC will work harder to
stay under this limit. This is incredibly useful in containers:
  - Container has 512MB → set GOMEMLIMIT=450MiB (leave headroom for stacks, etc.)
  - The GC will increase its pace as memory approaches the limit
  - Much better than GOGC alone for containerized services

Pro tip: In production, set GOMEMLIMIT to ~80% of your container's memory
limit and leave GOGC at the default. This gives you predictable memory usage
without constant GC overhead.

=============================================================================
 sync.Pool — REDUCING GC PRESSURE
=============================================================================

sync.Pool is a concurrent-safe pool of temporary objects. Objects in a pool
may be removed at any time (during GC), so never rely on them persisting.

Use sync.Pool when:
  - You allocate and discard the same type of object frequently
  - Objects are expensive to create (large buffers, parsed templates)
  - You're in a hot path and want to reduce allocation rate

The classic use case is buffer reuse in HTTP handlers:

  var bufPool = sync.Pool{
      New: func() interface{} { return new(bytes.Buffer) },
  }

  func handleRequest(data []byte) string {
      buf := bufPool.Get().(*bytes.Buffer)
      buf.Reset()            // Always reset before use!
      defer bufPool.Put(buf) // Return to pool when done
      buf.Write(data)
      return buf.String()
  }

=============================================================================
 REDUCING ALLOCATIONS — PRACTICAL TECHNIQUES
=============================================================================

1. Pre-allocate slices when you know the size:
     make([]User, 0, len(userIDs))  // not just []User{}

2. Use strings.Builder for concatenation:
     var b strings.Builder
     b.Grow(estimatedSize)  // pre-allocate internal buffer
     for _, s := range parts { b.WriteString(s) }

3. Reuse buffers with sync.Pool (as above)

4. Avoid unnecessary pointer indirection for small structs

5. Use arrays instead of slices for small, fixed-size collections:
     var cache [8]Entry  // lives on stack, no heap allocation

6. Watch out for interface boxing — storing a small concrete value in an
   interface{} (or any) can cause an allocation

=============================================================================
 strings.Builder VS STRING CONCATENATION
=============================================================================

String concatenation with += is O(n^2) because strings are immutable.
Each += creates a new string and copies everything:

  s := ""
  for i := 0; i < 1000; i++ {
      s += "x"   // 1000 allocations, each copying all previous data
  }

strings.Builder grows a byte buffer and only creates the final string once:

  var b strings.Builder
  for i := 0; i < 1000; i++ {
      b.WriteString("x")  // appends to internal buffer, occasional regrowth
  }
  result := b.String()  // one allocation for final string

bytes.Buffer is similar but optimized for []byte. strings.Builder is
optimized for building strings (avoids a copy in String()).

strings.Join is also efficient — it calculates total length, allocates once,
and copies all strings. Great when you have a slice of strings.

=============================================================================
 INTERFACE BOXING COST
=============================================================================

When you store a concrete value in an interface variable, Go may need to
allocate. For example:

  var x interface{} = 42  // may allocate to put int on heap

The compiler is smart about this for small values and common cases, but in
hot loops, interface boxing can add up. This is one reason why generics
(Go 1.18+) can be faster than interface-based polymorphism — generics are
monomorphized at compile time, avoiding boxing entirely.

=============================================================================
 MEMORY PROFILES
=============================================================================

To understand where your allocations come from:
  go test -memprofile=mem.out -bench=.
  go tool pprof mem.out

Key pprof commands:
  - top: shows functions with most allocations
  - list FuncName: shows annotated source code
  - web: opens a call graph in your browser

Two allocation metrics:
  - alloc_space: total bytes allocated (includes freed objects)
  - inuse_space: bytes currently in use (what's live right now)

Use alloc_space to find where allocations happen (reduce GC pressure).
Use inuse_space to find memory leaks.

=============================================================================
 STRUCT LAYOUT AND PADDING
=============================================================================

Go aligns struct fields in memory. A bool (1 byte) followed by an int64
(8 bytes) wastes 7 bytes of padding. Reordering fields from largest to
smallest can reduce struct size:

  // Bad: 24 bytes (with padding)
  type Bad struct {
      a bool    // 1 byte + 7 padding
      b int64   // 8 bytes
      c bool    // 1 byte + 7 padding
  }

  // Good: 16 bytes
  type Good struct {
      b int64   // 8 bytes
      a bool    // 1 byte
      c bool    // 1 byte + 6 padding
  }

This matters when you have millions of these structs in memory.

=============================================================================
 FINALIZERS — USUALLY A CODE SMELL
=============================================================================

runtime.SetFinalizer lets you attach a cleanup function to an object.
The GC calls it before collecting the object. Sounds useful, right?

Problems:
  - No guarantee WHEN it runs (or if it runs at all before program exit)
  - Finalizers delay garbage collection (object survives one extra GC cycle)
  - Order of finalization is unpredictable
  - They make the GC's job harder

Better alternatives:
  - Explicit Close() methods (the Go idiom)
  - defer for scoped cleanup
  - context.Context for request-scoped resources

The one semi-legitimate use: safety nets to detect resource leaks during
development (log a warning if Close wasn't called).

=============================================================================
*/

// DemoStackVsHeap shows variables that stay on the stack vs escape to heap.
// Run: go build -gcflags="-m" ./25-memory-and-gc/
// to see which variables escape.
func DemoStackVsHeap() {
	// This stays on the stack — it doesn't escape the function.
	x := 42
	_ = x

	// This escapes to the heap — we're passing it to fmt.Println,
	// which accepts interface{}, causing boxing and escape.
	y := 42
	fmt.Println("y escapes to heap because of interface boxing:", y)
}

// DemoEscapeAnalysis demonstrates common escape scenarios.
func DemoEscapeAnalysis() {
	// Case 1: Returning a pointer causes escape
	p := newPointer()
	fmt.Println("Pointer from heap:", *p)

	// Case 2: Value stays on stack (passed by value)
	v := newValue()
	fmt.Println("Value from stack:", v)

	// Case 3: Closure captures variable, may cause escape
	counter := 0
	increment := func() {
		counter++ // counter escapes because closure captures it
	}
	increment()
	fmt.Println("Counter:", counter)
}

// newPointer returns a pointer — the value MUST escape to heap
// because the caller needs it after this function returns.
func newPointer() *int {
	x := 42
	return &x // x escapes to heap
}

// newValue returns a value — it's copied to the caller's stack frame.
// No heap allocation needed.
func newValue() int {
	x := 42
	return x // x stays on stack, value is copied
}

// DemoStringConcatenation shows why += is expensive for building strings.
func DemoStringConcatenation() {
	// BAD: O(n^2) — each += allocates a new string
	s := ""
	for i := 0; i < 100; i++ {
		s += "x"
	}
	fmt.Println("Concat length:", len(s))

	// GOOD: strings.Builder — amortized O(n)
	var b strings.Builder
	b.Grow(100) // optional but helpful: pre-allocate
	for i := 0; i < 100; i++ {
		b.WriteString("x")
	}
	fmt.Println("Builder length:", len(b.String()))

	// ALSO GOOD: bytes.Buffer — similar performance
	var buf bytes.Buffer
	buf.Grow(100)
	for i := 0; i < 100; i++ {
		buf.WriteString("x")
	}
	fmt.Println("Buffer length:", buf.Len())

	// GOOD for slices: strings.Join — allocates once
	parts := make([]string, 100)
	for i := range parts {
		parts[i] = "x"
	}
	joined := strings.Join(parts, "")
	fmt.Println("Join length:", len(joined))
}

// DemoSyncPool shows how to use sync.Pool to reuse buffers.
func DemoSyncPool() {
	// Create a pool of bytes.Buffer objects
	pool := &sync.Pool{
		New: func() any {
			fmt.Println("  Creating new buffer")
			return new(bytes.Buffer)
		},
	}

	// First Get: pool is empty, calls New
	buf := pool.Get().(*bytes.Buffer)
	buf.WriteString("hello")
	fmt.Println("Got from pool:", buf.String())

	// Return to pool (always Reset before Put!)
	buf.Reset()
	pool.Put(buf)

	// Second Get: reuses the buffer we put back
	buf2 := pool.Get().(*bytes.Buffer)
	buf2.WriteString("world")
	fmt.Println("Reused from pool:", buf2.String())
	buf2.Reset()
	pool.Put(buf2)
}

// DemoPreAllocation shows the difference between growing and pre-allocating.
func DemoPreAllocation() {
	// BAD: starts with capacity 0, grows multiple times
	var s1 []int
	for i := 0; i < 1000; i++ {
		s1 = append(s1, i) // Multiple allocations as slice grows
	}

	// GOOD: pre-allocate with known capacity
	s2 := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		s2 = append(s2, i) // No reallocation — capacity was set
	}

	fmt.Printf("Both have length %d\n", len(s1))
	_ = s2
}

// DemoMemStats shows how to read runtime memory statistics.
func DemoMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc (current heap):   %d KB\n", m.Alloc/1024)
	fmt.Printf("TotalAlloc (lifetime):  %d KB\n", m.TotalAlloc/1024)
	fmt.Printf("Sys (OS memory):        %d KB\n", m.Sys/1024)
	fmt.Printf("NumGC (GC cycles):      %d\n", m.NumGC)
	fmt.Printf("HeapObjects (live):     %d\n", m.HeapObjects)
}

// DemoStructPadding shows how field ordering affects struct size.
func DemoStructPadding() {
	// Poorly ordered: wastes space on padding
	type Padded struct {
		a bool  // 1 byte + 7 padding
		b int64 // 8 bytes
		c bool  // 1 byte + 7 padding
		d int64 // 8 bytes
	}

	// Well ordered: minimizes padding
	type Compact struct {
		b int64 // 8 bytes
		d int64 // 8 bytes
		a bool  // 1 byte
		c bool  // 1 byte + 6 padding
	}

	type Padded2 = Padded
	type Compact2 = Compact

	// unsafe.Sizeof would show the actual sizes, but we can demonstrate
	// the concept without importing unsafe:
	fmt.Printf("Padded struct has fields: bool, int64, bool, int64\n")
	fmt.Printf("  Expected size: 32 bytes (with padding)\n")
	fmt.Printf("Compact struct has fields: int64, int64, bool, bool\n")
	fmt.Printf("  Expected size: 24 bytes (fields grouped by alignment)\n")
}
