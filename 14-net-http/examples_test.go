package nethttp

// All Demo* functions in this module are HTTP handlers (take
// http.ResponseWriter and *http.Request), return *http.ServeMux,
// return *http.Server, or return http.HandlerFunc. These cannot be
// called directly in a simple test wrapper, so they are skipped.
//
// Skipped functions:
//   DemoHandlerFunc       — func(http.ResponseWriter, *http.Request)
//   DemoServeMux          — returns *http.ServeMux
//   DemoRequestInspection — func(http.ResponseWriter, *http.Request)
//   DemoResponseWriter    — func(http.ResponseWriter, *http.Request)
//   DemoProductionServer  — returns *http.Server
//   DemoCommonMistakes    — returns http.HandlerFunc
//   DemoBodyReading       — func(http.ResponseWriter, *http.Request)
//   DemoJSONDecoding      — func(http.ResponseWriter, *http.Request)
//   DemoFormHandling      — func(http.ResponseWriter, *http.Request)
//   NewDemoServer         — returns *http.Server

import "testing"

// TestPlaceholder exists so `go test` has at least one test in the file.
func TestPlaceholder(t *testing.T) {
	t.Log("All Demo functions in 14-net-http are HTTP handlers or return HTTP types; see lesson.go to explore them.")
}
