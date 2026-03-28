package main

import (
	"context"
	"sync"
)

// Result wraps a value with a possible error, so pipeline stages can
// propagate errors through channels without stopping the pipeline.
type Result[T any] struct {
	Value T
	Err   error
}

// StageFunc is the signature for a single pipeline stage. It takes an input
// value and produces an output value (or an error).
type StageFunc[In any, Out any] func(ctx context.Context, input In) (Out, error)

// RunStage starts `workers` goroutines that each read from `in`, apply `fn`,
// and send the result to the returned channel. The output channel is closed
// once all workers finish.
//
// TODO: Implement this function.
func RunStage[In any, Out any](
	ctx context.Context,
	in <-chan Result[In],
	fn StageFunc[In, Out],
	workers int,
) <-chan Result[Out] {
	out := make(chan Result[Out])

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ctx
			_ = fn
			// TODO: Range over in, call fn, send Result[Out] to out.
			// Skip items whose Err is already non-nil.
			// Respect ctx.Done().
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
