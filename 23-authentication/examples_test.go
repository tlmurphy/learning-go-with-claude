package authentication

import "testing"

func TestDemoPasswordHashing(t *testing.T) {
	DemoPasswordHashing()
}

func TestDemoTokenSigning(t *testing.T) {
	DemoTokenSigning()
}

func TestDemoAPIKeyAuth(t *testing.T) {
	DemoAPIKeyAuth()
}

// AuthMiddleware returns func(http.Handler) http.Handler — skipped.

func TestDemoAuthMiddleware(t *testing.T) {
	DemoAuthMiddleware()
}

// RequireRole returns func(http.Handler) http.Handler — skipped.

func TestDemoRBAC(t *testing.T) {
	DemoRBAC()
}

func TestDemoSecurityPractices(t *testing.T) {
	DemoSecurityPractices()
}
