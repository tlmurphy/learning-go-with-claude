package apitesting

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
=============================================================================
 EXERCISES: API Testing
=============================================================================

 Work through these exercises in order. Each one builds on concepts from
 the lesson. Run the tests with:

   go test -v ./24-api-testing/

 Tip: Run a single test at a time while working:

   go test -v -run TestTableDrivenCRUD ./24-api-testing/

 NOTE: These exercises are unusual — you're writing TEST HELPERS and
 TEST INFRASTRUCTURE, not application code. The exercises_test.go file
 tests that your helpers work correctly.

=============================================================================
*/

// Exercise 1: Table-Driven API Test Cases
//
// Define a struct type and constructor for table-driven API tests.
// This will be the foundation for many of the later exercises.
//
// APITestCase represents a single API test case. It should have fields for:
//   - Name: descriptive test name (string)
//   - Method: HTTP method (string)
//   - Path: URL path (string)
//   - Body: request body as string (empty for GET/DELETE)
//   - Headers: map[string]string of request headers
//   - ExpectedStatus: expected HTTP status code (int)
//   - ExpectedBodyContains: string that should appear in response body
//
// RunAPITest executes a single test case against a handler.
// It should:
//   1. Create a request using httptest.NewRequest (use strings.NewReader for body)
//   2. Set all headers from the test case
//   3. Create a ResponseRecorder
//   4. Call the handler
//   5. Check the status code matches ExpectedStatus
//   6. If ExpectedBodyContains is non-empty, check the body contains it
//
// RunAPITests runs multiple test cases as subtests using t.Run.

type APITestCase struct {
	// YOUR CODE HERE — add fields
	Name                 string
	Method               string
	Path                 string
	Body                 string
	Headers              map[string]string
	ExpectedStatus       int
	ExpectedBodyContains string
}

func RunAPITest(t *testing.T, handler http.Handler, tc APITestCase) {
	t.Helper()
	// YOUR CODE HERE
}

func RunAPITests(t *testing.T, handler http.Handler, tests []APITestCase) {
	t.Helper()
	// YOUR CODE HERE
}

// Exercise 2: Test Middleware in Isolation
//
// Implement TestableMiddleware — a helper that makes it easy to test
// middleware without building a full handler chain.
//
// MiddlewareResult captures what happened when middleware was applied:
//   - NextCalled: whether the middleware called the next handler
//   - StatusCode: the response status code
//   - Headers: response headers
//   - Body: response body as string
//   - RequestHeaders: headers that the next handler received
//     (middleware might add/modify headers before passing to next)
//
// TestMiddleware applies the middleware to a simple recording handler,
// sends the request, and captures the results.
//
// Parameters:
//   - middleware: the middleware to test (func(http.Handler) http.Handler)
//   - req: the request to send through the middleware
//
// Returns a MiddlewareResult.

type MiddlewareResult struct {
	NextCalled     bool
	StatusCode     int
	Headers        http.Header
	Body           string
	RequestHeaders http.Header // What the next handler received
}

func TestMiddleware(middleware func(http.Handler) http.Handler, req *http.Request) *MiddlewareResult {
	// YOUR CODE HERE
	return &MiddlewareResult{}
}

// Exercise 3: Test Helper Functions
//
// Build commonly used assertion helpers for API testing. Every helper
// MUST call t.Helper() so failure messages point to the caller.
//
// AssertStatus: Check that the response recorder has the expected status code.
// AssertJSON: Check that the response has Content-Type: application/json.
// AssertBodyContains: Check that the response body contains the expected string.
// AssertHeader: Check that a specific response header has the expected value.
//
// MakeRequest: A convenience function that creates a request, sends it to
// the handler, and returns the response recorder. Reduces boilerplate.
//
// MakeAuthRequest: Like MakeRequest but also sets the Authorization header
// with a Bearer token.

// TB is the interface for testing helpers — satisfied by both *testing.T
// and *testing.B. Using this instead of *testing.T lets our helpers work
// in benchmarks too, and makes them easier to test.
type TB interface {
	Helper()
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func AssertStatus(t TB, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()
	// YOUR CODE HERE
}

func AssertJSON(t TB, rr *httptest.ResponseRecorder) {
	t.Helper()
	// YOUR CODE HERE
}

func AssertBodyContains(t TB, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	// YOUR CODE HERE
}

func AssertHeader(t TB, rr *httptest.ResponseRecorder, key, expected string) {
	t.Helper()
	// YOUR CODE HERE
}

func MakeRequest(t TB, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	// YOUR CODE HERE
	return httptest.NewRecorder()
}

func MakeAuthRequest(t TB, handler http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	// YOUR CODE HERE
	return httptest.NewRecorder()
}

// Exercise 4: Integration Test with httptest.Server
//
// Implement SetupTestServer which creates a fully configured test server
// with routing and middleware. This simulates what your real server looks
// like in production, but for tests.
//
// Requirements:
//   - Create a new ServeMux
//   - Register handlers:
//     GET /users/{id}  → HandleGetUser
//     GET /users       → HandleListUsers
//     POST /users      → HandleCreateUser
//     DELETE /users/{id} → HandleDeleteUser
//   - Wrap the mux with the provided middleware (applied in order)
//   - Return an httptest.Server (caller will defer server.Close())
//   - Also return the UserStore so tests can verify state changes
//
// Note: Use http.NewServeMux() and HandleFunc for routing. For the path
// matching, use the pattern "/users/" to match both /users and /users/id.
// In Go 1.22+, you can use "GET /users/{id}" patterns.

type TestServer struct {
	Server *httptest.Server
	Store  *UserStore
}

func SetupTestServer(middlewares ...func(http.Handler) http.Handler) *TestServer {
	// YOUR CODE HERE
	return &TestServer{}
}

// Exercise 5: Test Fixture Factory
//
// Build a factory for creating test data. In real applications, setting
// up test data is one of the most tedious parts of testing. Factories
// make it easy to create valid objects with sensible defaults, overriding
// only the fields that matter for each test.
//
// Requirements:
//   - UserFactory creates User objects with auto-generated defaults.
//   - NewUserFactory creates a factory with a counter starting at 1.
//   - MakeUser creates a user with defaults:
//     ID: "user-N" (N is the counter, auto-incremented)
//     Name: "Test User N"
//     Email: "user-N@test.com"
//   - WithName, WithEmail, WithID return options that override defaults.
//     Use the functional options pattern.
//   - MakeUsers creates N users with defaults.

type UserOption func(*User)

type UserFactory struct {
	// YOUR CODE HERE — add fields
	counter int
}

func NewUserFactory() *UserFactory {
	// YOUR CODE HERE
	return &UserFactory{}
}

func WithName(name string) UserOption {
	// YOUR CODE HERE
	return func(u *User) {}
}

func WithEmail(email string) UserOption {
	// YOUR CODE HERE
	return func(u *User) {}
}

func WithID(id string) UserOption {
	// YOUR CODE HERE
	return func(u *User) {}
}

func (f *UserFactory) MakeUser(opts ...UserOption) User {
	// YOUR CODE HERE
	return User{}
}

func (f *UserFactory) MakeUsers(count int) []User {
	// YOUR CODE HERE
	return nil
}

// Exercise 6: Golden File Testing
//
// Implement golden file comparison helpers.
//
// CompareWithGolden:
//   - Takes the test, a golden file name, and the actual response bytes
//   - Loads the golden file using LoadGoldenFile (from lesson.go)
//   - Normalizes both actual and expected JSON using NormalizeJSON (from lesson.go)
//   - Compares them and fails the test if they differ
//   - On failure, show both expected and actual for easy debugging
//
// CompareResponseWithGolden:
//   - Takes the test, a golden file name, and a ResponseRecorder
//   - Extracts the body from the recorder
//   - Delegates to CompareWithGolden

func CompareWithGolden(t *testing.T, goldenFile string, actual []byte) {
	t.Helper()
	// YOUR CODE HERE
}

func CompareResponseWithGolden(t *testing.T, goldenFile string, rr *httptest.ResponseRecorder) {
	t.Helper()
	// YOUR CODE HERE
}

// Exercise 7: Error Scenario Testing
//
// Implement a comprehensive error test runner that tests various error
// conditions for an API endpoint.
//
// ErrorTestCase describes an error scenario to test:
//   - Name: descriptive test name
//   - Method: HTTP method
//   - Path: URL path
//   - Body: request body (may be malformed JSON)
//   - Headers: request headers
//   - ExpectedStatus: expected HTTP status code
//   - ExpectedError: expected error message in the response JSON's "error" field
//
// RunErrorTests runs all error test cases as subtests, checking both
// the status code and the error message.

type ErrorTestCase struct {
	Name           string
	Method         string
	Path           string
	Body           string
	Headers        map[string]string
	ExpectedStatus int
	ExpectedError  string
}

func RunErrorTests(t *testing.T, handler http.Handler, tests []ErrorTestCase) {
	t.Helper()
	// YOUR CODE HERE
}

// Exercise 8: Complete API Test Suite
//
// Build a complete test suite builder that ties everything together.
// This is the "bring it all together" exercise.
//
// APITestSuite provides a convenient API for testing a CRUD service:
//   - Setup creates the test server and factory
//   - TestCreate verifies POST creates a resource and returns 201
//   - TestGet verifies GET returns a resource by ID
//   - TestGetNotFound verifies GET returns 404 for missing resources
//   - TestList verifies GET returns a list of resources
//   - TestDelete verifies DELETE removes a resource and returns 204
//   - TestDeleteNotFound verifies DELETE returns 404 for missing resources
//   - RunAll runs all the above tests as subtests
//
// The suite uses SetupTestServer internally and cleans up after itself.

type APITestSuite struct {
	// YOUR CODE HERE — add fields
	server  *TestServer
	factory *UserFactory
}

func NewAPITestSuite() *APITestSuite {
	// YOUR CODE HERE
	return &APITestSuite{}
}

func (s *APITestSuite) Setup(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) Teardown(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) TestCreate(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) TestGet(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) TestGetNotFound(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) TestList(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) TestDelete(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) TestDeleteNotFound(t *testing.T) {
	t.Helper()
	// YOUR CODE HERE
}

func (s *APITestSuite) RunAll(t *testing.T) {
	// YOUR CODE HERE
}
