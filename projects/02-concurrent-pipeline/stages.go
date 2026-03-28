package main

import (
	"context"
)

// This file contains the concrete stage functions for the text processing
// pipeline. Each function matches the StageFunc signature.

// WordCount holds the frequency count for each word in a single file.
type WordCount map[string]int

// ReadFile reads a file and returns its contents as a string.
// StageFunc[string, string] — input is a file path, output is file contents.
//
// TODO: Implement using os.ReadFile.
func ReadFile(ctx context.Context, path string) (string, error) {
	_ = ctx
	return "", nil
}

// Tokenize splits text into individual lowercased words.
// StageFunc[string, []string]
//
// TODO: Implement. Consider strings.Fields and strings.ToLower.
// Strip punctuation from words.
func Tokenize(ctx context.Context, text string) ([]string, error) {
	_ = ctx
	return nil, nil
}

// FilterStopWords removes common stop words from a word list.
// StageFunc[[]string, []string]
//
// TODO: Define a set of stop words and filter them out.
func FilterStopWords(ctx context.Context, words []string) ([]string, error) {
	_ = ctx
	return nil, nil
}

// CountWords counts the frequency of each word in the list.
// StageFunc[[]string, WordCount]
//
// TODO: Implement using a map.
func CountWords(ctx context.Context, words []string) (WordCount, error) {
	_ = ctx
	return nil, nil
}

// MergeWordCounts merges multiple WordCount maps into one.
// This is the final aggregation step — not a StageFunc but a standalone
// function that drains a channel.
//
// TODO: Implement. Read all WordCounts from the channel and merge.
func MergeWordCounts(ctx context.Context, counts <-chan Result[WordCount]) (WordCount, error) {
	_ = ctx
	return nil, nil
}
