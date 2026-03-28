# Project 05: Production-Ready Service (Capstone)

**Prerequisite modules:** 25-29 (Performance & Quality)

## Overview

This project is different from the others. Instead of building from scratch,
you are handed a working but messy "order processing" service and tasked with
transforming it into production-quality code.

The service handles orders for a fictional e-commerce system. It accepts orders
via HTTP, validates them against an external inventory service (simulated),
processes payments (simulated), and stores the results. It works. It also has
a lot of problems.

## The Starting Code

The code you have been given is **functional** — it compiles, it runs, it
processes orders. But it was written quickly with no regard for production
concerns. Your job is to identify and fix every issue.

### What is wrong with it (this list is not exhaustive)

The codebase has problems in several categories. Part of the exercise is
finding them all, so this list is intentionally vague:

- **Concurrency bugs** — there are race conditions and goroutine leaks hiding
  in the code. Run it with `-race` and under load to find them.
- **Error handling** — errors are swallowed, ignored, or returned without
  context. Some error paths leave the system in inconsistent states.
- **Architecture** — the code uses global mutable state and tight coupling.
  It is hard to test any single piece in isolation.
- **Observability** — there is no structured logging, no health checks, no
  way to know if the service is healthy or struggling.
- **Reliability** — there is no graceful shutdown, no retry logic for external
  calls, no circuit breaking when the inventory service is down.
- **Performance** — there are memory inefficiencies in hot paths: unnecessary
  allocations, string concatenation in loops, missing pre-allocation.
- **API design** — error responses are inconsistent, status codes are wrong
  in places, and the JSON handling has subtle issues.

### This works but...

- The goroutine that processes orders in the background? It leaks if the
  server shuts down. And it doesn't have any error recovery.
- The global `orders` map? Multiple goroutines read and write it without
  synchronization.
- The `processOrder` function that calls the inventory service? If the
  inventory service is slow, it holds up everything with no timeout.
- The handler that builds a response string by concatenating in a loop?
  Under load, that is going to create a lot of garbage.
- The error handling that returns `"something went wrong"` for every error?
  Good luck debugging that in production.
- The configuration that is hardcoded? Have fun redeploying to change the port.

## Your Mission

### Phase 1: Understand
Read all the code. Make notes. Run it. Hit it with requests. Run it with
`go run -race .` and see what happens.

### Phase 2: Test
Write tests for the **current** behavior before changing anything. These tests
are your safety net — they prove that your refactoring does not break functionality.

### Phase 3: Fix Architecture
- Remove global state
- Introduce dependency injection
- Define clean interfaces between layers
- Make the code testable in isolation

### Phase 4: Fix Bugs
- Fix race conditions
- Fix goroutine leaks
- Fix error handling
- Fix incorrect HTTP status codes

### Phase 5: Add Production Concerns
- Structured logging with `log/slog`
- Health checks (liveness + readiness)
- Graceful shutdown
- Context propagation
- Circuit breaker for external service calls
- Retry with exponential backoff
- Configuration from environment variables

### Phase 6: Optimize
- Find and fix memory inefficiencies
- Pre-allocate slices where sizes are known
- Use `sync.Pool` where appropriate
- Use `strings.Builder` instead of concatenation
- Write benchmarks that prove your optimizations work

## What Success Looks Like

- `go vet ./...` passes
- `go run -race .` shows no data races under load
- Tests cover core behavior
- Benchmarks prove performance improvements
- The service starts up cleanly, handles requests, and shuts down gracefully
- Logs are structured and useful
- Health checks report actual readiness
- External service failures are handled gracefully (circuit breaker, retry)

## Stretch Goals

- **OpenTelemetry-style tracing** — add spans to trace a request through the system
- **Metrics dashboard** — expose a `/debug/metrics` endpoint with request counts, latencies, error rates
- **Integration tests** — test the full HTTP flow with a test server
- **Performance budget** — set allocation targets for key operations and enforce them in benchmarks
