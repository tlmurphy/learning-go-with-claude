package authentication

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// =============================================================================
// Exercise 1: Password Hashing Tests
// =============================================================================

func TestPasswordHashing(t *testing.T) {
	hasher := NewSimpleHasher([]byte("test-secret-key-for-hashing"))

	t.Run("hash produces non-empty string", func(t *testing.T) {
		hash, err := hasher.Hash("mypassword")
		if err != nil {
			t.Fatalf("Hash returned error: %v", err)
		}
		if hash == "" {
			t.Error("Hash should return a non-empty string. " +
				"Use HMAC-SHA256 with the key, then encode as hex with fmt.Sprintf.")
		}
	})

	t.Run("same password produces same hash", func(t *testing.T) {
		hash1, _ := hasher.Hash("mypassword")
		hash2, _ := hasher.Hash("mypassword")
		if hash1 != hash2 {
			t.Error("Hashing the same password should produce the same hash. " +
				"Our simplified hasher is deterministic (unlike bcrypt which has random salt).")
		}
	})

	t.Run("different passwords produce different hashes", func(t *testing.T) {
		hash1, _ := hasher.Hash("password1")
		hash2, _ := hasher.Hash("password2")
		if hash1 == hash2 {
			t.Error("Different passwords should produce different hashes.")
		}
	})

	t.Run("verify correct password", func(t *testing.T) {
		hash, _ := hasher.Hash("mypassword")
		err := hasher.Verify("mypassword", hash)
		if err != nil {
			t.Errorf("Verify should return nil for correct password, got: %v. "+
				"Hash the password again and compare with TimingSafeEqual.", err)
		}
	})

	t.Run("verify wrong password", func(t *testing.T) {
		hash, _ := hasher.Hash("mypassword")
		err := hasher.Verify("wrongpassword", hash)
		if err == nil {
			t.Error("Verify should return ErrPasswordMismatch for wrong password.")
		}
		if err != ErrPasswordMismatch {
			t.Errorf("Expected ErrPasswordMismatch, got: %v", err)
		}
	})

	t.Run("different keys produce different hashes", func(t *testing.T) {
		hasher2 := NewSimpleHasher([]byte("different-key"))
		hash1, _ := hasher.Hash("samepassword")
		hash2, _ := hasher2.Hash("samepassword")
		if hash1 == hash2 {
			t.Error("Different keys should produce different hashes. " +
				"Make sure you're using h.key in the HMAC.")
		}
	})
}

// =============================================================================
// Exercise 2: HMAC Token Tests
// =============================================================================

func TestHMACToken(t *testing.T) {
	secret := []byte("my-test-secret-key")

	t.Run("create and verify token", func(t *testing.T) {
		token, err := CreateHMACToken("hello-world", secret)
		if err != nil {
			t.Fatalf("CreateHMACToken failed: %v", err)
		}
		if token == "" {
			t.Fatal("Token should not be empty. " +
				"Return base64url(payload) + \".\" + base64url(hmac).")
		}

		payload, err := VerifyHMACToken(token, secret)
		if err != nil {
			t.Fatalf("VerifyHMACToken failed: %v", err)
		}
		if payload != "hello-world" {
			t.Errorf("Expected payload 'hello-world', got %q. "+
				"Decode the base64url payload to get the original string.", payload)
		}
	})

	t.Run("token has two parts", func(t *testing.T) {
		token, _ := CreateHMACToken("test", secret)
		// Should be payload.signature
		parts := splitToken(token)
		if len(parts) != 2 {
			t.Errorf("Token should have 2 parts separated by '.', got %d parts. "+
				"Format: base64url(payload).base64url(hmac)", len(parts))
		}
	})

	t.Run("wrong secret fails verification", func(t *testing.T) {
		token, _ := CreateHMACToken("secret-data", secret)
		_, err := VerifyHMACToken(token, []byte("wrong-secret"))
		if err == nil {
			t.Error("Verification with wrong secret should fail. " +
				"Use hmac.Equal for timing-safe comparison.")
		}
	})

	t.Run("tampered token fails verification", func(t *testing.T) {
		token, _ := CreateHMACToken("original", secret)
		_, err := VerifyHMACToken(token+"tampered", secret)
		if err == nil {
			t.Error("Tampered token should fail verification.")
		}
	})

	t.Run("invalid format fails", func(t *testing.T) {
		_, err := VerifyHMACToken("no-dot-separator", secret)
		if err == nil {
			t.Error("Token without '.' separator should fail verification.")
		}
	})
}

// splitToken splits a token on '.'. Helper to avoid importing strings in test.
func splitToken(token string) []string {
	var parts []string
	current := ""
	for _, c := range token {
		if c == '.' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)
	return parts
}

// =============================================================================
// Exercise 3: Token Service Tests
// =============================================================================

func TestTokenService(t *testing.T) {
	ts := NewTokenService([]byte("test-secret-32-bytes-minimum!!!!"), "test-issuer")

	t.Run("issue and validate token", func(t *testing.T) {
		token, err := ts.IssueToken("user-123", []string{"admin"}, 1*time.Hour)
		if err != nil {
			t.Fatalf("IssueToken failed: %v", err)
		}
		if token == "" {
			t.Fatal("Token should not be empty.")
		}

		claims, err := ts.ValidateToken(token)
		if err != nil {
			t.Fatalf("ValidateToken failed: %v", err)
		}
		if claims == nil {
			t.Fatal("ValidateToken returned nil claims. Return the decoded Claims on success.")
		}
		if claims.Subject != "user-123" {
			t.Errorf("Expected subject 'user-123', got %q. "+
				"Set the Subject field in the Claims.", claims.Subject)
		}
		if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
			t.Errorf("Expected roles [admin], got %v", claims.Roles)
		}
		if claims.Issuer != "test-issuer" {
			t.Errorf("Expected issuer 'test-issuer', got %q. "+
				"Set the Issuer field to ts.issuer.", claims.Issuer)
		}
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		// Issue a token that's already expired
		token, err := ts.IssueToken("user-456", []string{"viewer"}, -1*time.Hour)
		if err != nil {
			t.Fatalf("IssueToken failed: %v", err)
		}

		_, err = ts.ValidateToken(token)
		if err == nil {
			t.Error("ValidateToken should reject expired tokens. " +
				"Check claims.IsExpired(time.Now()) during validation.")
		}
	})

	t.Run("wrong issuer is rejected", func(t *testing.T) {
		otherTS := NewTokenService([]byte("test-secret-32-bytes-minimum!!!!"), "other-issuer")
		token, _ := otherTS.IssueToken("user-789", nil, 1*time.Hour)

		_, err := ts.ValidateToken(token)
		if err == nil {
			t.Error("ValidateToken should reject tokens from a different issuer. " +
				"Compare claims.Issuer with ts.issuer.")
		}
	})

	t.Run("tampered token is rejected", func(t *testing.T) {
		token, _ := ts.IssueToken("user-000", nil, 1*time.Hour)
		_, err := ts.ValidateToken(token + "tampered")
		if err == nil {
			t.Error("ValidateToken should reject tampered tokens.")
		}
	})

	t.Run("token has correct timestamps", func(t *testing.T) {
		before := time.Now().Unix()
		token, _ := ts.IssueToken("user-ts", nil, 1*time.Hour)
		after := time.Now().Unix()

		claims, err := ts.ValidateToken(token)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if claims == nil {
			t.Fatal("ValidateToken returned nil claims. Return the decoded Claims on success.")
		}

		if claims.IssuedAt < before || claims.IssuedAt > after {
			t.Errorf("IssuedAt should be approximately now. Got %d, expected between %d and %d",
				claims.IssuedAt, before, after)
		}

		expectedExp := claims.IssuedAt + int64(time.Hour.Seconds())
		if claims.ExpiresAt != expectedExp {
			t.Errorf("ExpiresAt should be IssuedAt + TTL. Got %d, expected %d",
				claims.ExpiresAt, expectedExp)
		}
	})
}

// =============================================================================
// Exercise 4: Auth Middleware Tests
// =============================================================================

func TestExtractBearerToken(t *testing.T) {
	t.Run("valid bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer mytoken123")

		token, err := ExtractBearerToken(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if token != "mytoken123" {
			t.Errorf("Expected 'mytoken123', got %q. "+
				"Split on space and return the second part.", token)
		}
	})

	t.Run("case insensitive Bearer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "bearer mytoken")

		token, err := ExtractBearerToken(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if token != "mytoken" {
			t.Errorf("Expected 'mytoken', got %q. "+
				"Use strings.EqualFold for case-insensitive comparison.", token)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		_, err := ExtractBearerToken(req)
		if err == nil {
			t.Error("Should return error when Authorization header is missing.")
		}
	})

	t.Run("wrong scheme", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

		_, err := ExtractBearerToken(req)
		if err == nil {
			t.Error("Should return error for non-Bearer schemes.")
		}
	})
}

func TestBuildAuthMiddleware(t *testing.T) {
	ts := NewTokenService([]byte("test-secret-32-bytes-minimum!!!!"), "test-issuer")

	t.Run("valid token passes through", func(t *testing.T) {
		token, _ := ts.IssueToken("user-123", []string{"admin"}, 1*time.Hour)

		var capturedUser *AuthenticatedUser
		handler := BuildAuthMiddleware(ts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedUser = UserFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", rr.Code)
		}
		if capturedUser == nil {
			t.Fatal("Expected user in context. Use ContextWithUser to add the user.")
		}
		if capturedUser.ID != "user-123" {
			t.Errorf("Expected user ID 'user-123', got %q", capturedUser.ID)
		}
	})

	t.Run("missing token returns 401", func(t *testing.T) {
		handler := BuildAuthMiddleware(ts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called when token is missing.")
		}))

		req := httptest.NewRequest("GET", "/protected", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rr.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		handler := BuildAuthMiddleware(ts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called with invalid token.")
		}))

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rr.Code)
		}
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		token, _ := ts.IssueToken("user-expired", nil, -1*time.Hour)

		handler := BuildAuthMiddleware(ts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called with expired token.")
		}))

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rr.Code)
		}
	})
}

// =============================================================================
// Exercise 5: Role-Based Authorization Tests
// =============================================================================

func TestAuthorizer(t *testing.T) {
	perms := map[string][]Permission{
		"admin":  {PermCreate, PermRead, PermUpdate, PermDelete, PermManageUsers},
		"editor": {PermCreate, PermRead, PermUpdate},
		"viewer": {PermRead},
	}
	authz := NewAuthorizer(perms)

	t.Run("admin can do everything", func(t *testing.T) {
		ctx := ContextWithUser(context.Background(), &AuthenticatedUser{
			ID:    "alice",
			Roles: []string{"admin"},
		})

		for _, perm := range []Permission{PermCreate, PermRead, PermUpdate, PermDelete, PermManageUsers} {
			if !authz.CanPerform(ctx, perm) {
				t.Errorf("Admin should have %s permission. "+
					"Get user from context and check their roles against rolePerms.", perm)
			}
		}
	})

	t.Run("viewer can only read", func(t *testing.T) {
		ctx := ContextWithUser(context.Background(), &AuthenticatedUser{
			ID:    "carol",
			Roles: []string{"viewer"},
		})

		if !authz.CanPerform(ctx, PermRead) {
			t.Error("Viewer should have read permission.")
		}
		if authz.CanPerform(ctx, PermCreate) {
			t.Error("Viewer should NOT have create permission.")
		}
		if authz.CanPerform(ctx, PermDelete) {
			t.Error("Viewer should NOT have delete permission.")
		}
	})

	t.Run("no user returns false", func(t *testing.T) {
		ctx := context.Background() // No user in context
		if authz.CanPerform(ctx, PermRead) {
			t.Error("CanPerform should return false when no user is in context.")
		}
	})

	t.Run("unknown role returns false", func(t *testing.T) {
		ctx := ContextWithUser(context.Background(), &AuthenticatedUser{
			ID:    "unknown",
			Roles: []string{"superuser"}, // Not in rolePerms
		})
		if authz.CanPerform(ctx, PermRead) {
			t.Error("CanPerform should return false for unknown roles.")
		}
	})
}

func TestRequirePermissionMiddleware(t *testing.T) {
	perms := map[string][]Permission{
		"admin":  {PermCreate, PermRead, PermUpdate, PermDelete},
		"viewer": {PermRead},
	}
	authz := NewAuthorizer(perms)

	t.Run("authorized user passes through", func(t *testing.T) {
		called := false
		handler := authz.RequirePermission(PermRead)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/", nil)
		ctx := ContextWithUser(req.Context(), &AuthenticatedUser{ID: "alice", Roles: []string{"viewer"}})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !called {
			t.Error("Handler should be called for authorized users.")
		}
		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rr.Code)
		}
	})

	t.Run("unauthorized user gets 403", func(t *testing.T) {
		handler := authz.RequirePermission(PermDelete)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called for unauthorized users.")
		}))

		req := httptest.NewRequest("DELETE", "/resource", nil)
		ctx := ContextWithUser(req.Context(), &AuthenticatedUser{ID: "carol", Roles: []string{"viewer"}})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden, got %d. "+
				"Return 403 when user is authenticated but lacks permission.", rr.Code)
		}
	})

	t.Run("no user gets 401", func(t *testing.T) {
		handler := authz.RequirePermission(PermRead)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called without authentication.")
		}))

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d. "+
				"Return 401 when no user is in the context.", rr.Code)
		}
	})
}

// =============================================================================
// Exercise 6: API Key Validator Tests
// =============================================================================

func TestAPIKeyValidator(t *testing.T) {
	keys := map[string]string{
		"sk_live_key123": "service-a",
		"sk_live_key456": "service-b",
	}
	validator := NewAPIKeyValidator(keys)

	t.Run("valid key returns owner", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-API-Key", "sk_live_key123")

		owner, err := validator.ValidateRequest(req)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if owner != "service-a" {
			t.Errorf("Expected owner 'service-a', got %q. "+
				"Use TimingSafeEqual to compare keys and return the owner.", owner)
		}
	})

	t.Run("invalid key returns error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-API-Key", "sk_live_wrong")

		_, err := validator.ValidateRequest(req)
		if err == nil {
			t.Error("Should return error for invalid API key.")
		}
	})

	t.Run("missing key returns error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		_, err := validator.ValidateRequest(req)
		if err == nil {
			t.Error("Should return error when X-API-Key header is missing.")
		}
	})

	t.Run("middleware passes valid key", func(t *testing.T) {
		var capturedUser *AuthenticatedUser
		handler := validator.APIKeyMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedUser = UserFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("GET", "/api/data", nil)
		req.Header.Set("X-API-Key", "sk_live_key456")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", rr.Code)
		}
		if capturedUser == nil {
			t.Fatal("Expected AuthenticatedUser in context.")
		}
		if capturedUser.ID != "service-b" {
			t.Errorf("Expected user ID 'service-b', got %q", capturedUser.ID)
		}
	})

	t.Run("middleware rejects invalid key", func(t *testing.T) {
		handler := validator.APIKeyMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Handler should not be called with invalid key.")
		}))

		req := httptest.NewRequest("GET", "/api/data", nil)
		req.Header.Set("X-API-Key", "sk_live_invalid")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rr.Code)
		}
	})
}

// =============================================================================
// Exercise 7: Token Refresh Flow Tests
// =============================================================================

func TestRefreshService(t *testing.T) {
	tokenService := NewTokenService([]byte("test-secret-32-bytes-minimum!!!!"), "test-app")
	refreshStore := NewRefreshTokenStore()
	userRoles := map[string][]string{
		"user-1": {"admin"},
		"user-2": {"viewer"},
	}
	rs := NewRefreshService(tokenService, refreshStore, userRoles, 15*time.Minute, 7*24*time.Hour)

	t.Run("login returns both tokens", func(t *testing.T) {
		result, err := rs.Login("user-1", []string{"admin"})
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if result == nil {
			t.Fatal("Login should return a LoginResult.")
		}
		if result.AccessToken == "" {
			t.Error("Login should return an access token. " +
				"Use TokenService.IssueToken to create it.")
		}
		if result.RefreshToken == "" {
			t.Error("Login should return a refresh token. " +
				"Use RefreshTokenStore.GenerateToken to create it.")
		}

		// Verify access token is valid
		claims, err := tokenService.ValidateToken(result.AccessToken)
		if err != nil {
			t.Fatalf("Access token should be valid: %v", err)
		}
		if claims == nil {
			t.Fatal("ValidateToken returned nil claims.")
		}
		if claims.Subject != "user-1" {
			t.Errorf("Expected subject 'user-1', got %q", claims.Subject)
		}
	})

	t.Run("refresh issues new tokens", func(t *testing.T) {
		loginResult, _ := rs.Login("user-2", []string{"viewer"})
		if loginResult == nil {
			t.Fatal("Login should return a non-nil LoginResult.")
		}

		refreshResult, err := rs.Refresh(loginResult.RefreshToken)
		if err != nil {
			t.Fatalf("Refresh failed: %v", err)
		}
		if refreshResult == nil {
			t.Fatal("Refresh should return a LoginResult with new tokens.")
		}
		if refreshResult.AccessToken == "" {
			t.Error("Refresh should return a new access token.")
		}
		if refreshResult.RefreshToken == "" {
			t.Error("Refresh should return a new refresh token.")
		}

		// New access token should be valid
		claims, err := tokenService.ValidateToken(refreshResult.AccessToken)
		if err != nil {
			t.Fatalf("New access token should be valid: %v", err)
		}
		if claims == nil {
			t.Fatal("ValidateToken returned nil claims.")
		}
		if claims.Subject != "user-2" {
			t.Errorf("Expected subject 'user-2', got %q", claims.Subject)
		}
	})

	t.Run("refresh token is single-use", func(t *testing.T) {
		loginResult, _ := rs.Login("user-1", []string{"admin"})
		if loginResult == nil {
			t.Fatal("Login should return a non-nil LoginResult.")
		}

		// First refresh succeeds
		_, err := rs.Refresh(loginResult.RefreshToken)
		if err != nil {
			t.Fatalf("First refresh should succeed: %v", err)
		}

		// Second refresh with same token should fail
		_, err = rs.Refresh(loginResult.RefreshToken)
		if err == nil {
			t.Error("Second use of refresh token should fail. " +
				"RefreshTokenStore.ValidateAndConsume marks tokens as used.")
		}
	})

	t.Run("invalid refresh token fails", func(t *testing.T) {
		_, err := rs.Refresh("totally-invalid-token")
		if err == nil {
			t.Error("Refresh with invalid token should return an error.")
		}
	})
}

// =============================================================================
// Exercise 8: Auth Context Tests
// =============================================================================

func TestAuthContext(t *testing.T) {
	rolePerms := map[string][]Permission{
		"admin":  {PermCreate, PermRead, PermUpdate, PermDelete, PermManageUsers},
		"editor": {PermCreate, PermRead, PermUpdate},
		"viewer": {PermRead},
	}

	t.Run("create from claims", func(t *testing.T) {
		claims := &Claims{
			Subject:   "user-123",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "test",
			Roles:     []string{"editor"},
		}

		ac := NewAuthContext(claims, "the-token", rolePerms)
		if ac == nil {
			t.Fatal("NewAuthContext should return a non-nil AuthContext.")
		}
		if ac.UserID != "user-123" {
			t.Errorf("Expected UserID 'user-123', got %q. "+
				"Set UserID from claims.Subject.", ac.UserID)
		}
		if ac.Token != "the-token" {
			t.Errorf("Expected Token 'the-token', got %q", ac.Token)
		}
		if len(ac.Roles) != 1 || ac.Roles[0] != "editor" {
			t.Errorf("Expected roles [editor], got %v", ac.Roles)
		}

		// Editor should have create, read, update permissions
		expectedPerms := map[Permission]bool{
			PermCreate: true,
			PermRead:   true,
			PermUpdate: true,
		}
		for _, p := range ac.Permissions {
			if !expectedPerms[p] {
				t.Errorf("Unexpected permission %q for editor role", p)
			}
			delete(expectedPerms, p)
		}
		if len(expectedPerms) > 0 {
			t.Errorf("Missing permissions for editor: %v. "+
				"Resolve permissions from roles using the rolePerms map.", expectedPerms)
		}
	})

	t.Run("store and retrieve from context", func(t *testing.T) {
		ac := &AuthContext{
			UserID: "user-456",
			Roles:  []string{"viewer"},
		}

		ctx := WithAuthContext(context.Background(), ac)
		retrieved := GetAuthContext(ctx)

		if retrieved == nil {
			t.Fatal("GetAuthContext should return the stored AuthContext. " +
				"Use context.WithValue and type assertion.")
		}
		if retrieved.UserID != "user-456" {
			t.Errorf("Expected UserID 'user-456', got %q", retrieved.UserID)
		}
	})

	t.Run("IsAuthenticated", func(t *testing.T) {
		if IsAuthenticated(context.Background()) {
			t.Error("IsAuthenticated should return false for context without AuthContext.")
		}

		ac := &AuthContext{UserID: "user-789"}
		ctx := WithAuthContext(context.Background(), ac)
		if !IsAuthenticated(ctx) {
			t.Error("IsAuthenticated should return true when AuthContext is present.")
		}
	})

	t.Run("HasAnyRole", func(t *testing.T) {
		ac := &AuthContext{
			UserID: "user-1",
			Roles:  []string{"editor", "viewer"},
		}

		if !ac.HasAnyRole("editor") {
			t.Error("HasAnyRole should return true when user has the role.")
		}
		if !ac.HasAnyRole("admin", "editor") {
			t.Error("HasAnyRole should return true when user has any of the roles.")
		}
		if ac.HasAnyRole("admin", "superuser") {
			t.Error("HasAnyRole should return false when user has none of the roles.")
		}
	})

	t.Run("HasAllPermissions", func(t *testing.T) {
		ac := &AuthContext{
			UserID:      "user-1",
			Permissions: []Permission{PermCreate, PermRead, PermUpdate},
		}

		if !ac.HasAllPermissions(PermRead) {
			t.Error("HasAllPermissions should return true for a single matching permission.")
		}
		if !ac.HasAllPermissions(PermCreate, PermRead) {
			t.Error("HasAllPermissions should return true when all are present.")
		}
		if ac.HasAllPermissions(PermCreate, PermDelete) {
			t.Error("HasAllPermissions should return false when any permission is missing.")
		}
	})
}

// =============================================================================
// Test helpers
// =============================================================================

// parseJSONError extracts an error message from a JSON response body.
func parseJSONError(t *testing.T, body []byte) string {
	t.Helper()
	var result struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse JSON error response: %v", err)
	}
	return result.Error
}
