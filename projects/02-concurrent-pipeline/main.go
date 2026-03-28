package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Concurrent Text Processing Pipeline
//
// Steps to build:
//  1. Define your generic pipeline stage types in pipeline.go
//  2. Implement each concrete stage function in stages.go
//  3. Wire the stages together here in main()
//  4. Add context timeout and cancellation
//  5. Print the aggregated word counts at the end
//
// Usage:
//   pipeline <directory>
//
// The program reads all .txt files from the directory and processes them
// through the pipeline, printing the top N most frequent words.

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	dir := os.Args[1]

	// Create a context with a timeout so the pipeline doesn't run forever.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO: Build and run the pipeline:
	//  1. List .txt files in dir
	//  2. Feed file paths into the Read stage
	//  3. Pipe through Tokenize → Filter → Count → Aggregate
	//  4. Print results

	_ = ctx
	_ = dir

	fmt.Println("Pipeline - not yet implemented")
}
