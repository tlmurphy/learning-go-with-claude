// Package profiling covers Go's benchmarking tools, CPU and memory profiling
// with pprof, and the optimization workflow used in production systems.
package profiling

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

/*
=============================================================================
 PROFILING AND BENCHMARKS IN GO
=============================================================================

Go has world-class profiling tools built right into the standard library.
No third-party APM agent, no bytecode instrumentation, no separate profiler
binary. It's all there: benchmarks in the test framework, pprof for CPU and
memory profiling, and the trace tool for execution visualization.

The golden rule of optimization:

  "Premature optimization is the root of all evil." — Donald Knuth

But the FULL quote continues: "Yet we should not pass up our opportunities
in that critical 3%."

The real skill isn't knowing how to optimize — it's knowing WHEN and WHERE.
Profile first. Always profile first. Your intuition about what's slow is
almost always wrong.

=============================================================================
 BENCHMARKS: func BenchmarkXxx(b *testing.B)
=============================================================================

Go benchmarks live alongside your tests. The test framework handles
everything: warm-up, iteration count, timing.

  func BenchmarkFoo(b *testing.B) {
      for i := 0; i < b.N; i++ {
          Foo()
      }
  }

Go 1.24+ Modern Alternative — b.Loop():

  func BenchmarkFoo(b *testing.B) {
      for b.Loop() {
          Foo()
      }
  }

b.Loop() is simpler and also prevents the compiler from optimizing
away the loop body — a common benchmarking pitfall with the b.N pattern.

Key points:
  - b.N is set by the framework — it runs enough iterations for stable timing
  - The framework automatically increases b.N until the benchmark runs long
    enough (at least 1 second by default)
  - NEVER use a fixed loop count — always use b.N

Run benchmarks:
  go test -bench=.                    # run all benchmarks
  go test -bench=BenchmarkFoo         # run specific benchmark
  go test -bench=. -benchmem          # include allocation stats
  go test -bench=. -count=5           # run 5 times for statistical accuracy
  go test -bench=. -benchtime=5s      # run each benchmark for 5 seconds

=============================================================================
 b.ResetTimer, b.StopTimer, b.StartTimer
=============================================================================

Sometimes you need setup that shouldn't be timed:

  func BenchmarkProcess(b *testing.B) {
      data := expensiveSetup() // don't count this
      b.ResetTimer()           // reset the clock
      for i := 0; i < b.N; i++ {
          Process(data)
      }
  }

For per-iteration setup:

  func BenchmarkWithSetup(b *testing.B) {
      for i := 0; i < b.N; i++ {
          b.StopTimer()
          data := prepareTestData()
          b.StartTimer()
          Process(data)
      }
  }

Warning: StopTimer/StartTimer in a loop adds overhead. Prefer ResetTimer
with one-time setup when possible.

=============================================================================
 b.ReportAllocs()
=============================================================================

Call b.ReportAllocs() to include allocation statistics:

  func BenchmarkFoo(b *testing.B) {
      b.ReportAllocs()
      for i := 0; i < b.N; i++ {
          Foo()
      }
  }

Output:
  BenchmarkFoo-8   1000000   1234 ns/op   256 B/op   3 allocs/op

This tells you each call to Foo takes ~1234ns, allocates 256 bytes across
3 separate heap allocations. Reducing allocs/op often has more impact than
micro-optimizing CPU.

=============================================================================
 SUB-BENCHMARKS
=============================================================================

Compare approaches side by side:

  func BenchmarkSort(b *testing.B) {
      sizes := []int{10, 100, 1000, 10000}
      for _, size := range sizes {
          b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
              data := makeData(size)
              b.ResetTimer()
              for i := 0; i < b.N; i++ {
                  sort.Ints(data)
              }
          })
      }
  }

Output:
  BenchmarkSort/size=10-8      10000000    120 ns/op
  BenchmarkSort/size=100-8      1000000   1050 ns/op
  BenchmarkSort/size=1000-8      100000  15000 ns/op

=============================================================================
 b.RunParallel — CONCURRENT BENCHMARKS
=============================================================================

For benchmarking concurrent code:

  func BenchmarkConcurrent(b *testing.B) {
      b.RunParallel(func(pb *testing.PB) {
          for pb.Next() {
              DoWork()
          }
      })
  }

This runs DoWork() across multiple goroutines simultaneously, which is
essential for benchmarking things like sync.Pool, sync.Map, or any
concurrent data structure.

=============================================================================
 CPU PROFILING WITH pprof
=============================================================================

Generate a CPU profile:
  go test -cpuprofile=cpu.out -bench=.

Analyze it:
  go tool pprof cpu.out

Interactive commands:
  top          — functions using the most CPU
  top -cum     — functions with most cumulative time (including callees)
  list FuncName — annotated source code showing time per line
  web          — open call graph in browser (needs graphviz)
  png          — save call graph as PNG

The web UI (recommended):
  go tool pprof -http=:8080 cpu.out
  Opens a browser with flame graphs, call graphs, source view, etc.

=============================================================================
 MEMORY PROFILING
=============================================================================

Generate a memory profile:
  go test -memprofile=mem.out -bench=.

Two key metrics:
  - alloc_space: total bytes allocated over time (shows allocation rate)
  - inuse_space: bytes currently in use (shows memory leaks)

  go tool pprof -alloc_space mem.out   # where are allocations happening?
  go tool pprof -inuse_space mem.out   # what's currently live?

For finding memory leaks: use inuse_space
For reducing GC pressure: use alloc_space

=============================================================================
 GOROUTINE PROFILING
=============================================================================

To find goroutine leaks in a running service, expose pprof over HTTP:

  import _ "net/http/pprof"  // registers handlers on DefaultServeMux

  go func() {
      log.Println(http.ListenAndServe("localhost:6060", nil))
  }()

Then: go tool pprof http://localhost:6060/debug/pprof/goroutine

This shows you every goroutine's stack trace. If you see thousands of
goroutines stuck in the same place, you have a leak.

=============================================================================
 BLOCK PROFILING
=============================================================================

Block profiling shows where goroutines block on synchronization:
  - Channel sends/receives
  - Mutex Lock calls
  - Select statements

Enable in tests:
  go test -blockprofile=block.out -bench=.

Useful for finding contention in concurrent code.

=============================================================================
 THE TRACE TOOL
=============================================================================

go tool trace gives you a timeline view of your program:
  go test -trace=trace.out -bench=.
  go tool trace trace.out

Shows:
  - Goroutine scheduling (when each goroutine runs/blocks)
  - GC events and their duration
  - System calls
  - Network I/O

Go 1.25 adds runtime/trace.FlightRecorder for continuous trace recording in
production — it keeps a rolling buffer of trace data you can dump on demand
when something goes wrong.

=============================================================================
 net/http/pprof — PROFILING LIVE SERVICES
=============================================================================

For production services, import net/http/pprof to expose profiling
endpoints:

  import _ "net/http/pprof"

This registers handlers at /debug/pprof/ on the default mux. In production,
serve pprof on a SEPARATE port that's not exposed to the internet:

  // Main API server
  go http.ListenAndServe(":8080", apiRouter)

  // Admin/debug server (internal network only)
  go http.ListenAndServe("localhost:6060", nil)

Then from your machine:
  go tool pprof http://your-server:6060/debug/pprof/profile?seconds=30

=============================================================================
 BENCHSTAT — COMPARING RESULTS
=============================================================================

benchstat compares benchmark results statistically:

  go test -bench=. -count=10 > old.txt
  # make changes
  go test -bench=. -count=10 > new.txt
  benchstat old.txt new.txt

Output:
  name     old time/op   new time/op   delta
  Foo-8    1.23µs ± 2%   0.89µs ± 1%  -27.64% (p=0.000 n=10+10)

The ± shows variance, p-value shows statistical significance, and delta
shows the improvement. Use -count=10 minimum for meaningful results.

Install: go install golang.org/x/perf/cmd/benchstat@latest

=============================================================================
 THE OPTIMIZATION WORKFLOW
=============================================================================

  1. Write correct code first. Make it work.
  2. Write benchmarks for the hot path.
  3. Profile to find the actual bottleneck (don't guess!).
  4. Optimize the bottleneck.
  5. Benchmark again to verify improvement.
  6. Repeat from step 3.

Common mistakes:
  - Optimizing code that isn't the bottleneck (95% of optimization attempts)
  - Benchmarking without -count (results vary between runs)
  - Not using b.ResetTimer (setup time pollutes results)
  - Compiler optimizing away your benchmark (see "compiler tricks" below)

=============================================================================
 PREVENTING COMPILER OPTIMIZATIONS IN BENCHMARKS
=============================================================================

The Go compiler can optimize away function calls whose results aren't used.
To prevent this:

  var result int  // package-level variable

  func BenchmarkFoo(b *testing.B) {
      var r int
      for i := 0; i < b.N; i++ {
          r = Foo()
      }
      result = r  // prevent the compiler from eliminating Foo()
  }

Without assigning to a package-level variable, the compiler might notice
that r is never used outside the benchmark and eliminate the Foo() call.

Note: If you use b.Loop() (Go 1.24+) instead of the manual b.N loop,
the compiler cannot eliminate the loop body, so this sink trick is unnecessary.

=============================================================================
*/

// DemoBenchmarkPattern shows the standard benchmark loop pattern.
// This isn't runnable as a benchmark itself (that's in _test.go files),
// but it demonstrates the pattern.
func DemoBenchmarkPattern() {
	fmt.Println("Standard benchmark pattern:")
	fmt.Println(`
  func BenchmarkMyFunc(b *testing.B) {
      // Setup (not timed)
      data := prepareData()
      b.ResetTimer()

      for i := 0; i < b.N; i++ {
          MyFunc(data)
      }
  }`)
}

// SlowFibonacci computes fibonacci numbers using naive recursion.
// This is deliberately slow — O(2^n) — to use as a profiling target.
func SlowFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return SlowFibonacci(n-1) + SlowFibonacci(n-2)
}

// FastFibonacci computes fibonacci numbers using iteration.
// O(n) time, O(1) space — the correct approach.
func FastFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// SortInts sorts a copy of the input slice (doesn't modify original).
// Used to demonstrate sub-benchmarks with different input sizes.
func SortInts(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	sort.Ints(result)
	return result
}

// IsPrime checks if a number is prime. Two implementations for comparison.
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n < 4 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// IsPrimeNaive checks if a number is prime using the naive approach.
// Checks all numbers up to n-1 (terrible but good for benchmarking).
func IsPrimeNaive(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// CountPrimesBelow counts primes below n using the given primality check.
func CountPrimesBelow(n int, isPrime func(int) bool) int {
	count := 0
	for i := 2; i < n; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

// BuildStringSlowly creates a string by repeated concatenation (slow).
func BuildStringSlowly(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "x"
	}
	return s
}

// BuildStringFast creates a string using strings.Builder (fast).
func BuildStringFast(n int) string {
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		b.WriteString("x")
	}
	return b.String()
}

// SumSquares computes the sum of squares from 1 to n.
// Two approaches: loop vs formula. Used to show that algorithmic
// improvements beat micro-optimization.
func SumSquaresLoop(n int) int64 {
	var sum int64
	for i := int64(1); i <= int64(n); i++ {
		sum += i * i
	}
	return sum
}

// SumSquaresFormula uses the mathematical formula: n(n+1)(2n+1)/6
func SumSquaresFormula(n int) int64 {
	n64 := int64(n)
	return n64 * (n64 + 1) * (2*n64 + 1) / 6
}

// ProcessDataWithReflection uses fmt.Sprintf (which uses reflection) to
// convert values. This is the "slow" approach.
func ProcessDataWithReflection(values []int) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = fmt.Sprintf("%d", v)
	}
	return result
}

// ProcessDataDirect uses direct string building without reflection.
func ProcessDataDirect(values []int) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = itoa(v)
	}
	return result
}

// itoa is a simple integer-to-string conversion without fmt (no reflection).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	// Max int64 has 19 digits
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// Ensure math import is used
var _ = math.MaxFloat64
