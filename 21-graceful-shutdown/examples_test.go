package shutdown

// DemoSignalContext is the only Demo* function in this module.
// It blocks waiting for an OS signal (SIGINT/SIGTERM), so it cannot
// be called in a test wrapper without hanging indefinitely.
//
// Skipped functions:
//   DemoSignalContext — blocks on signal.NotifyContext

import "testing"

// TestPlaceholder exists so `go test` has at least one test in the file.
func TestPlaceholder(t *testing.T) {
	t.Log("DemoSignalContext blocks on OS signals and cannot run in a test; see lesson.go to explore it.")
}
