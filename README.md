# Learning Go — An Interactive Tutorial

A hands-on, test-driven Go curriculum. Each module teaches concepts through
commented lesson files, then challenges you with exercises verified by tests.

## How to Use This Tutorial

### For Each Module

1. **Read** `lesson.go` — commented code explaining each concept with examples
2. **Implement** `exercises.go` — fill in the stubbed functions (look for `// YOUR CODE HERE`)
3. **Verify** — run the tests to check your work:

```bash
# Run tests for a specific module
go test ./01-variables-and-types/

# Run with verbose output to see individual test names
go test -v ./01-variables-and-types/

# Run a single test
go test -v -run TestExercise1 ./01-variables-and-types/
```

### For Projects

Projects are open-ended challenges that combine multiple modules. Each has a
README with requirements, hints, and stretch goals. Build them from scratch!

```bash
# Run project tests
go test ./projects/01-cli-task-manager/...
```

## Curriculum

### Foundations (Modules 01-08)

| Module | Topic | Key Concepts |
|--------|-------|-------------|
| 01 | Variables & Types | declarations, type system, zero values, constants, iota, type conversions |
| 02 | Control Flow | if/else, switch, for loops, range, labels, break/continue |
| 03 | Functions | multiple returns, named returns, variadic, closures, defer, function types |
| 04 | Collections | arrays, slices, maps, make vs new, slice internals, copy, delete |
| 05 | Structs & Methods | struct definition, embedding, methods, value vs pointer receivers |
| 06 | Interfaces | implicit implementation, empty interface, type assertions, type switches |
| 07 | Pointers | pointer basics, pass by value vs reference, nil pointers, when to use pointers |
| 08 | Error Handling | error interface, wrapping, errors.Is/As, custom errors, panic/recover |

**Project 1: CLI Task Manager** — Build a command-line task tracker using everything from the foundations.

---

### Intermediate (Modules 09-13)

| Module | Topic | Key Concepts |
|--------|-------|-------------|
| 09 | Packages & Modules | package organization, visibility, go.mod, dependencies, internal packages |
| 10 | Goroutines & Channels | goroutines, channels, buffered channels, directional channels, deadlocks |
| 11 | Concurrency Patterns | WaitGroup, Mutex, select, context, worker pools, fan-in/fan-out |
| 12 | Testing | unit tests, table-driven tests, subtests, test helpers, mocks, testdata |
| 13 | Generics | type parameters, constraints, generic functions, generic types |

**Project 2: Concurrent Pipeline** — Build a data processing pipeline with fan-out/fan-in concurrency.

---

### Web Services (Modules 14-21)

| Module | Topic | Key Concepts |
|--------|-------|-------------|
| 14 | net/http Fundamentals | Handler, HandlerFunc, ServeMux, Request, ResponseWriter, ListenAndServe |
| 15 | Routing & URL Patterns | path parameters, method routing, ServeMux patterns (Go 1.22+), subrouting |
| 16 | Building a REST API | CRUD operations, resource design, status codes, request validation |
| 17 | Middleware | middleware pattern, chaining, logging, recovery, CORS, rate limiting |
| 18 | JSON & Serialization | encoding/json, struct tags, custom marshalers, streaming, validation |
| 19 | Database Access | database/sql, connection pooling, prepared statements, transactions, migrations |
| 20 | Configuration | environment variables, config files, flag package, 12-factor app |
| 21 | Graceful Shutdown | os/signal, context cancellation, connection draining, health checks |

**Project 3: Bookstore API** — Build a complete REST API with database, middleware, and graceful shutdown.

---

### Advanced Web (Modules 22-24)

| Module | Topic | Key Concepts |
|--------|-------|-------------|
| 22 | gRPC Services | protobuf, service definitions, unary/streaming RPCs, interceptors |
| 23 | Authentication | JWT, middleware auth, API keys, bcrypt, session management |
| 24 | API Testing | httptest, integration tests, test fixtures, golden files, test containers |

**Project 4: Microservice System** — Build a gRPC service with a REST gateway and authentication.

---

### Performance & Quality (Modules 25-29)

| Module | Topic | Key Concepts |
|--------|-------|-------------|
| 25 | Memory & GC | stack vs heap, escape analysis, allocations, GC tuning, sync.Pool |
| 26 | Profiling & Benchmarks | pprof, CPU/memory profiles, benchmarks, tracing, optimization workflow |
| 27 | Code Smells | common anti-patterns, interface pollution, premature abstraction, Go idioms |
| 28 | Dependency Injection | constructor injection, wire, functional options, clean architecture |
| 29 | Production Patterns | structured logging, metrics, health checks, observability, circuit breakers |

**Project 5: Production-Ready Service (Capstone)** — Build a fully production-hardened microservice with observability, graceful degradation, and performance tuning.

---

## Tips

- **Don't skip the lessons.** The exercises assume you've read the lesson file first.
- **Read the test failures.** They're written to guide you toward the solution.
- **Experiment freely.** Add `fmt.Println` statements, write extra functions, break things on purpose.
- **Use `go doc`.** For example: `go doc fmt.Sprintf` or `go doc net/http.Handler`.
- **Check your work incrementally.** Run tests after each exercise, not all at once.
- **The projects are where real learning happens.** Don't rush past the modules to skip the projects.

## Go Version

This tutorial targets **Go 1.26**.
