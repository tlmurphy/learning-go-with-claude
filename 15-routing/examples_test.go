package routing

// All Demo* functions in this module return *http.ServeMux, which cannot
// be meaningfully exercised by calling them in a simple test wrapper.
//
// Skipped functions:
//   DemoMethodRouting   — returns *http.ServeMux
//   DemoPathParams      — returns *http.ServeMux
//   DemoWildcardRoutes  — returns *http.ServeMux
//   DemoPrecedence      — returns *http.ServeMux
//   DemoTrailingSlash   — returns *http.ServeMux
//   DemoSubrouting      — returns *http.ServeMux
//   DemoCustomNotFound  — returns *http.ServeMux
//   RegisterRoutes      — returns *http.ServeMux

import "testing"

// TestPlaceholder exists so `go test` has at least one test in the file.
func TestPlaceholder(t *testing.T) {
	t.Log("All Demo functions in 15-routing return *http.ServeMux; see lesson.go to explore them.")
}
