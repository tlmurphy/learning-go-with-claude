# Project 02: Concurrent Pipeline

**Prerequisite modules:** 09-13 (Intermediate)

## Overview

Build a concurrent text processing pipeline that reads files from a directory,
processes them through multiple stages, and aggregates results. This project
puts goroutines, channels, sync primitives, context, and generics to work in
a real-world pattern.

Pipelines are everywhere in production Go: log processing, ETL, data enrichment,
image transformation. Learning this pattern well will pay off repeatedly.

## Requirements

### Core Pipeline

Build a pipeline with these stages, each running as one or more goroutines:

```
files on disk → Read → Tokenize → Filter → Count → Aggregate → result
```

1. **Read** — read each text file into a string
2. **Tokenize** — split the text into individual words
3. **Filter** — remove stop words (the, a, an, is, etc.)
4. **Count** — count word frequencies per file
5. **Aggregate** — merge per-file counts into a single result

### Concurrency Requirements
- Each stage runs in its own goroutine(s), connected by channels
- Support a configurable number of workers per stage (fan-out)
- Merge results from multiple workers back into one channel (fan-in)
- Use `sync.WaitGroup` to know when all workers in a stage are done
- Use `context.Context` for cancellation and timeout
- Collect and report errors without stopping the whole pipeline

### Technical Requirements
- Use generics to define a reusable pipeline stage type
- Channels should be typed and directional (`<-chan`, `chan<-`) in function signatures
- Close channels properly — no goroutine leaks
- Handle the case where the input directory is empty or missing

## Hints

<details>
<summary>Generic stage type</summary>

Define a generic function type for a pipeline stage:

```go
type StageFunc[In any, Out any] func(ctx context.Context, input In) (Out, error)
```

Then write a `RunStage` function that spins up N workers, each reading from
an input channel, applying the StageFunc, and sending results to an output channel.

</details>

<details>
<summary>Fan-out / fan-in pattern</summary>

Fan-out: start N goroutines all reading from the same input channel.
Fan-in: have all N goroutines write to the same output channel, and close it
when all N are done (use a WaitGroup).

```go
func FanOut[In, Out any](ctx context.Context, in <-chan In, fn StageFunc[In, Out], workers int) <-chan Out {
    out := make(chan Out)
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range in {
                // process and send to out
            }
        }()
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

</details>

<details>
<summary>Error handling without stopping</summary>

Use a `Result[T]` type that carries either a value or an error:

```go
type Result[T any] struct {
    Value T
    Err   error
}
```

Send `Result[T]` through channels instead of bare values.

</details>

<details>
<summary>errgroup</summary>

The `golang.org/x/sync/errgroup` package simplifies managing a group of
goroutines that can fail. It combines a WaitGroup with error collection and
integrates with context cancellation.

</details>

## Stretch Goals

- **Monitoring goroutine** — a separate goroutine that periodically logs items/sec throughput for each stage
- **Backpressure** — use bounded (buffered) channels and observe how it affects throughput
- **Pipeline visualization** — print an ASCII diagram showing stages, worker counts, and channel buffer sizes
- **Dynamic scaling** — monitor channel buffer utilization and spin up/down workers in response

## Test Data

Create a `testdata/` directory with a few `.txt` files to process. Any plain
text works — try pasting paragraphs from public-domain books.
