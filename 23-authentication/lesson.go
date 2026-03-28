// Package authentication teaches authentication and authorization patterns
// in Go web services — from password hashing and token creation to middleware
// and role-based access control.
package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
=============================================================================
 AUTHENTICATION IN GO WEB SERVICES
=============================================================================

Authentication ("authn") answers: "Who are you?"
Authorization ("authz") answers: "What are you allowed to do?"

They're often confused, but the distinction matters:
  - Authentication: Verifying identity (login, token validation)
  - Authorization: Checking permissions (can this user delete posts?)

A request typically flows through both:
  1. Auth middleware extracts and validates the token (authn)
  2. Handler checks if the authenticated user has permission (authz)

The golden rules of authentication:
  1. NEVER store passwords in plain text (use bcrypt or argon2)
  2. NEVER roll your own crypto (use well-tested libraries)
  3. ALWAYS use HTTPS in production
  4. ALWAYS set token expiration times
  5. ALWAYS use timing-safe comparison for secrets

=============================================================================
 PASSWORD HASHING
=============================================================================

When a user creates an account, you need to store their password. But you
must NEVER store the actual password. Instead, store a hash — a one-way
transformation that can't be reversed.

Why not just SHA-256? Because SHA-256 is designed to be FAST. An attacker
with a GPU can try billions of SHA-256 hashes per second. Password hashing
algorithms like bcrypt are deliberately SLOW (configurable cost factor).

How bcrypt works:
  1. Take the password and a random salt
  2. Run the Blowfish cipher many times (cost factor controls how many)
  3. Output: $2a$10$salt_here_hash_here
     - $2a$ = algorithm version
     - $10$ = cost factor (2^10 = 1024 rounds)
     - The rest is salt + hash, base64 encoded

The salt prevents rainbow table attacks — even if two users have the same
password, their hashes will be different because each gets a random salt.

In Go, use golang.org/x/crypto/bcrypt:

  import "golang.org/x/crypto/bcrypt"

  // Hash a password (cost 12 is recommended for 2024+)
  hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

  // Verify a password against a hash
  err := bcrypt.CompareHashAndPassword(hash, []byte(password))
  if err != nil {
      // Password doesn't match
  }

Since this tutorial avoids external dependencies, our exercises use
interfaces to test the pattern without requiring the bcrypt package.

=============================================================================
*/

// PasswordHasher defines the interface for password hashing.
// In production, implement this with bcrypt from golang.org/x/crypto/bcrypt.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

// DemoPasswordHashing explains the password hashing pattern.
func DemoPasswordHashing() {
	fmt.Println("=== Password Hashing ===")
	fmt.Println()
	fmt.Println("Never store passwords in plain text!")
	fmt.Println()
	fmt.Println("Production code:")
	fmt.Println("  hash, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)")
	fmt.Println("  err := bcrypt.CompareHashAndPassword(hash, []byte(pw))")
	fmt.Println()
	fmt.Println("The bcrypt hash includes the salt and cost factor:")
	fmt.Println("  $2a$12$LJ4R6g0T3I5Y3r8F6J8bXOq7r5eZK8wY...")
	fmt.Println()
	fmt.Println("Key points:")
	fmt.Println("  - Cost factor of 12+ is recommended (higher = slower = more secure)")
	fmt.Println("  - Each hash has its own random salt built in")
	fmt.Println("  - Verification is timing-safe (constant time)")
	fmt.Println("  - Never log or expose password hashes")
}

/*
=============================================================================
 JWT (JSON WEB TOKENS)
=============================================================================

JWTs are the most common token format for API authentication. A JWT has
three parts, separated by dots:

  header.payload.signature

  Header:  {"alg": "HS256", "typ": "JWT"}  (base64url encoded)
  Payload: {"sub": "user-123", "exp": 1699999999, ...}  (base64url encoded)
  Signature: HMAC-SHA256(header + "." + payload, secret)

The key insight: JWTs are NOT encrypted. Anyone can decode the header and
payload. The signature only proves the token wasn't tampered with and was
issued by someone who knows the secret.

What goes in the payload (called "claims"):
  - sub: Subject (user ID) — who this token represents
  - exp: Expiration time — ALWAYS set this
  - iat: Issued at — when the token was created
  - iss: Issuer — who created the token (your service name)
  - aud: Audience — who the token is intended for
  - Custom claims: roles, permissions, etc.

JWT Best Practices:
  1. Keep tokens short-lived (15 min for access tokens)
  2. Use refresh tokens for long-lived sessions
  3. Never store sensitive data in the payload (it's not encrypted!)
  4. Use strong secrets (256+ bits of randomness)
  5. Validate ALL claims (expiry, issuer, audience)
  6. Consider token revocation for logout/security events

JWT vs Session Cookies:
  - JWTs are stateless (no server-side storage needed)
  - Sessions are stateful (server stores session data)
  - JWTs can't be revoked without extra infrastructure
  - Sessions can be invalidated instantly
  - JWTs work well for microservices (no shared session store)
  - Sessions work well for monoliths

=============================================================================
*/

// Claims represents the payload of a JWT-like token.
type Claims struct {
	Subject   string            `json:"sub"`           // User ID
	ExpiresAt int64             `json:"exp"`           // Expiration (Unix timestamp)
	IssuedAt  int64             `json:"iat"`           // Issued at (Unix timestamp)
	Issuer    string            `json:"iss,omitempty"` // Who issued the token
	Roles     []string          `json:"roles,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

// IsExpired checks if the claims have expired relative to the given time.
func (c *Claims) IsExpired(now time.Time) bool {
	return now.Unix() > c.ExpiresAt
}

// HasRole checks if the claims include a specific role.
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// TokenSigner handles creating and verifying HMAC-signed tokens.
// This is a simplified version of what a JWT library does internally.
type TokenSigner struct {
	secret []byte
}

// NewTokenSigner creates a signer with the given secret key.
// In production, use at least 256 bits (32 bytes) of cryptographic randomness.
func NewTokenSigner(secret []byte) *TokenSigner {
	return &TokenSigner{secret: secret}
}

// Sign creates a signed token from claims. The format is:
//
//	base64url(json(claims)).base64url(hmac-sha256(claims_json, secret))
//
// This is a simplified JWT — real JWTs also have a header part.
func (ts *TokenSigner) Sign(claims *Claims) (string, error) {
	// Marshal claims to JSON
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal claims: %w", err)
	}

	// Base64url encode the claims
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create HMAC signature
	mac := hmac.New(sha256.New, ts.secret)
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	// Combine: payload.signature
	return payload + "." + sig, nil
}

// Verify validates a token's signature and returns the claims.
// Returns an error if the signature is invalid or claims can't be decoded.
func (ts *TokenSigner) Verify(token string) (*Claims, error) {
	// Split into payload and signature
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid token format: expected payload.signature")
	}

	payload, sigEncoded := parts[0], parts[1]

	// Verify the HMAC signature
	mac := hmac.New(sha256.New, ts.secret)
	mac.Write([]byte(payload))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	// CRITICAL: Use hmac.Equal for timing-safe comparison!
	// A regular == comparison leaks information about how many bytes match,
	// which can be exploited in a timing attack.
	if !hmac.Equal([]byte(sigEncoded), []byte(expectedSig)) {
		return nil, errors.New("invalid token signature")
	}

	// Decode the claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	return &claims, nil
}

// DemoTokenSigning shows how HMAC-based token signing works.
func DemoTokenSigning() {
	fmt.Println("=== Token Signing (Simplified JWT) ===")
	fmt.Println()

	// Create a signer with a secret
	signer := NewTokenSigner([]byte("my-super-secret-key-at-least-32-bytes!!"))

	// Create claims
	claims := &Claims{
		Subject:   "user-123",
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "my-app",
		Roles:     []string{"editor", "viewer"},
	}

	// Sign the token
	token, err := signer.Sign(claims)
	if err != nil {
		fmt.Printf("Error signing: %v\n", err)
		return
	}
	fmt.Printf("Token: %s\n\n", token)

	// Verify the token
	verified, err := signer.Verify(token)
	if err != nil {
		fmt.Printf("Error verifying: %v\n", err)
		return
	}
	fmt.Printf("Subject: %s\n", verified.Subject)
	fmt.Printf("Roles: %v\n", verified.Roles)
	fmt.Printf("Expired: %v\n", verified.IsExpired(time.Now()))

	// Try tampering with the token
	tamperedToken := token + "x"
	_, err = signer.Verify(tamperedToken)
	fmt.Printf("\nTampered token: %v\n", err)
}

/*
=============================================================================
 API KEY AUTHENTICATION
=============================================================================

API keys are the simplest form of authentication: the client sends a
secret string with each request, and the server checks if it's valid.

  Authorization: ApiKey sk_live_abc123def456

Pros: Simple, easy to implement, easy for developers to use
Cons: Static (no expiry), hard to scope permissions, easy to leak

When to use API keys:
  - Server-to-server communication
  - Third-party integrations
  - When simplicity matters more than granular security

CRITICAL: When comparing API keys, you MUST use timing-safe comparison
(crypto/subtle.ConstantTimeCompare). A regular string comparison (==)
returns as soon as it finds a mismatch, which means comparing "aaa" vs
"bbb" is faster than comparing "aaa" vs "aab". An attacker can measure
this difference to guess the key one character at a time.

=============================================================================
*/

// APIKeyStore manages API keys. In production, store hashed keys in a database.
type APIKeyStore struct {
	mu   sync.RWMutex
	keys map[string]APIKeyInfo // key string -> info
}

// APIKeyInfo contains metadata about an API key.
type APIKeyInfo struct {
	Name      string    // Human-readable name ("Production API Key")
	Owner     string    // User or service that owns this key
	CreatedAt time.Time
	Roles     []string // What this key is authorized to do
}

// NewAPIKeyStore creates a new store.
func NewAPIKeyStore() *APIKeyStore {
	return &APIKeyStore{
		keys: make(map[string]APIKeyInfo),
	}
}

// AddKey stores a new API key with its metadata.
func (s *APIKeyStore) AddKey(key string, info APIKeyInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[key] = info
}

// ValidateKey checks an API key using timing-safe comparison.
// Returns the key info if valid, or an error if not.
func (s *APIKeyStore) ValidateKey(candidate string) (*APIKeyInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// We need to iterate ALL keys to maintain constant-time behavior.
	// If we just did map lookup, the timing would reveal whether the
	// key prefix exists in the map.
	var matchedInfo *APIKeyInfo
	for stored, info := range s.keys {
		// subtle.ConstantTimeCompare returns 1 if equal, 0 if not.
		// It always takes the same amount of time regardless of where
		// the strings differ (as long as they're the same length).
		if len(stored) == len(candidate) &&
			subtle.ConstantTimeCompare([]byte(stored), []byte(candidate)) == 1 {
			infoCopy := info
			matchedInfo = &infoCopy
		}
	}

	if matchedInfo == nil {
		return nil, errors.New("invalid API key")
	}
	return matchedInfo, nil
}

// DemoAPIKeyAuth shows API key authentication.
func DemoAPIKeyAuth() {
	fmt.Println("=== API Key Authentication ===")
	fmt.Println()

	store := NewAPIKeyStore()
	store.AddKey("sk_live_abc123", APIKeyInfo{
		Name:  "Production Key",
		Owner: "service-a",
		Roles: []string{"read", "write"},
	})

	// Valid key
	info, err := store.ValidateKey("sk_live_abc123")
	if err == nil {
		fmt.Printf("Valid key: owner=%s, roles=%v\n", info.Owner, info.Roles)
	}

	// Invalid key
	_, err = store.ValidateKey("sk_live_wrong")
	fmt.Printf("Invalid key: %v\n", err)

	fmt.Println()
	fmt.Println("CRITICAL: Always use crypto/subtle.ConstantTimeCompare")
	fmt.Println("for secret comparison. Regular == leaks timing information.")
}

/*
=============================================================================
 AUTHENTICATION MIDDLEWARE
=============================================================================

The standard pattern for authentication in Go web services is middleware
that runs before your handlers:

  1. Extract the token from the Authorization header
  2. Validate the token (check signature, expiry, etc.)
  3. If valid, add the user info to the request context
  4. If invalid, return 401 Unauthorized
  5. Call the next handler

This keeps auth logic out of individual handlers. Each handler just
pulls the user from the context and trusts that the middleware has
already validated them.

=============================================================================
*/

// contextKey is an unexported type for context keys, preventing collisions.
type contextKey string

const userContextKey contextKey = "authenticated_user"

// AuthenticatedUser represents a user extracted from a valid token.
type AuthenticatedUser struct {
	ID    string
	Roles []string
}

// UserFromContext extracts the authenticated user from a request context.
// Returns nil if no user is present (request wasn't authenticated).
func UserFromContext(ctx context.Context) *AuthenticatedUser {
	user, _ := ctx.Value(userContextKey).(*AuthenticatedUser)
	return user
}

// ContextWithUser adds an authenticated user to a context.
func ContextWithUser(ctx context.Context, user *AuthenticatedUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// AuthMiddleware returns HTTP middleware that validates bearer tokens.
// If the token is valid, it adds the AuthenticatedUser to the request context.
// If invalid or missing, it returns 401 Unauthorized.
func AuthMiddleware(signer *TokenSigner) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			// Expect "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, `{"error": "invalid authorization format"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := parts[1]

			// Verify the token
			claims, err := signer.Verify(tokenStr)
			if err != nil {
				http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Check expiration
			if claims.IsExpired(time.Now()) {
				http.Error(w, `{"error": "token expired"}`, http.StatusUnauthorized)
				return
			}

			// Add user to context
			user := &AuthenticatedUser{
				ID:    claims.Subject,
				Roles: claims.Roles,
			}
			ctx := ContextWithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// DemoAuthMiddleware explains the middleware authentication pattern.
func DemoAuthMiddleware() {
	fmt.Println("=== Authentication Middleware ===")
	fmt.Println()
	fmt.Println("The pattern:")
	fmt.Println("  1. Client sends: Authorization: Bearer <token>")
	fmt.Println("  2. Middleware validates the token")
	fmt.Println("  3. Middleware adds user to request context")
	fmt.Println("  4. Handler reads user from context")
	fmt.Println()
	fmt.Println("In your handler:")
	fmt.Println("  user := UserFromContext(r.Context())")
	fmt.Println("  if user == nil {")
	fmt.Println("      // This shouldn't happen if middleware is set up correctly")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("This separation means:")
	fmt.Println("  - Handlers don't duplicate auth logic")
	fmt.Println("  - Auth changes happen in one place")
	fmt.Println("  - Handlers can focus on business logic")
}

/*
=============================================================================
 ROLE-BASED ACCESS CONTROL (RBAC)
=============================================================================

RBAC is the most common authorization pattern. Instead of assigning
permissions directly to users, you assign roles, and roles have permissions.

  User → Roles → Permissions

Example:
  admin  → [create, read, update, delete, manage_users]
  editor → [create, read, update]
  viewer → [read]

Why roles instead of direct permissions?
  - Easier to manage (change a role, affect all users with that role)
  - Easier to audit (who has admin access?)
  - Easier to understand (is this user an admin or viewer?)

RBAC in a middleware chain:
  Request → AuthMiddleware (who?) → RoleMiddleware (allowed?) → Handler

=============================================================================
*/

// Permission represents an action that can be authorized.
type Permission string

const (
	PermCreate      Permission = "create"
	PermRead        Permission = "read"
	PermUpdate      Permission = "update"
	PermDelete      Permission = "delete"
	PermManageUsers Permission = "manage_users"
)

// RolePermissions defines what each role is allowed to do.
var RolePermissions = map[string][]Permission{
	"admin":  {PermCreate, PermRead, PermUpdate, PermDelete, PermManageUsers},
	"editor": {PermCreate, PermRead, PermUpdate},
	"viewer": {PermRead},
}

// HasPermission checks if any of the given roles grant the requested permission.
func HasPermission(roles []string, perm Permission) bool {
	for _, role := range roles {
		perms, ok := RolePermissions[role]
		if !ok {
			continue
		}
		for _, p := range perms {
			if p == perm {
				return true
			}
		}
	}
	return false
}

// RequireRole returns middleware that checks if the authenticated user
// has at least one of the required roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				http.Error(w, `{"error": "not authenticated"}`, http.StatusUnauthorized)
				return
			}

			for _, required := range roles {
				for _, userRole := range user.Roles {
					if required == userRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, `{"error": "insufficient permissions"}`, http.StatusForbidden)
		})
	}
}

// DemoRBAC shows role-based access control.
func DemoRBAC() {
	fmt.Println("=== Role-Based Access Control ===")
	fmt.Println()

	roles := map[string][]string{
		"alice": {"admin"},
		"bob":   {"editor"},
		"carol": {"viewer"},
	}

	for user, userRoles := range roles {
		perms := []string{}
		for _, perm := range []Permission{PermCreate, PermRead, PermUpdate, PermDelete, PermManageUsers} {
			if HasPermission(userRoles, perm) {
				perms = append(perms, string(perm))
			}
		}
		fmt.Printf("  %s (%v): %v\n", user, userRoles, perms)
	}
}

/*
=============================================================================
 REFRESH TOKEN PATTERN
=============================================================================

Short-lived access tokens (15 min) are more secure, but users don't want
to log in every 15 minutes. The solution: refresh tokens.

The flow:
  1. User logs in → gets access token (15 min) + refresh token (7 days)
  2. Client uses access token for API calls
  3. When access token expires, client sends refresh token to get a new pair
  4. Refresh token is single-use (new one issued with each refresh)

Why two tokens?
  - Access tokens are sent with every request (high exposure)
  - Refresh tokens are sent rarely (low exposure)
  - If an access token leaks, it expires quickly
  - If a refresh token leaks, it can be revoked server-side

Refresh token storage:
  - Access tokens: stateless (JWT, validated by signature)
  - Refresh tokens: stateful (stored in database, can be revoked)

=============================================================================
*/

// RefreshTokenStore manages refresh tokens. In production, use a database.
type RefreshTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]RefreshTokenInfo // token string -> info
}

// RefreshTokenInfo contains metadata about a refresh token.
type RefreshTokenInfo struct {
	UserID    string
	ExpiresAt time.Time
	Used      bool // Single-use: marked true after first use
}

// NewRefreshTokenStore creates a new store.
func NewRefreshTokenStore() *RefreshTokenStore {
	return &RefreshTokenStore{
		tokens: make(map[string]RefreshTokenInfo),
	}
}

// GenerateToken creates a new refresh token for a user.
func (s *RefreshTokenStore) GenerateToken(userID string, ttl time.Duration) (string, error) {
	// Generate random bytes for the token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = RefreshTokenInfo{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}
	return token, nil
}

// ValidateAndConsume validates a refresh token and marks it as used.
// Returns the user ID if valid. Returns an error if the token is
// invalid, expired, or already used.
func (s *RefreshTokenStore) ValidateAndConsume(token string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, exists := s.tokens[token]
	if !exists {
		return "", errors.New("invalid refresh token")
	}
	if info.Used {
		// Possible token theft! In production, revoke ALL tokens for this user.
		return "", errors.New("refresh token already used (possible theft)")
	}
	if time.Now().After(info.ExpiresAt) {
		return "", errors.New("refresh token expired")
	}

	// Mark as used (single-use)
	info.Used = true
	s.tokens[token] = info
	return info.UserID, nil
}

// RevokeAllForUser deletes all refresh tokens for a user (e.g., on logout).
func (s *RefreshTokenStore) RevokeAllForUser(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for token, info := range s.tokens {
		if info.UserID == userID {
			delete(s.tokens, token)
		}
	}
}

/*
=============================================================================
 OAUTH2 CONCEPTS
=============================================================================

OAuth2 is the standard for delegated authorization. Instead of users
giving you their Google/GitHub password, they authorize your app through
the provider's UI.

The Authorization Code Flow (most common for web apps):
  1. User clicks "Login with Google" → redirect to Google
  2. User logs in at Google, approves your app
  3. Google redirects back to your app with an authorization code
  4. Your server exchanges the code for access + refresh tokens
  5. Your server uses the access token to call Google APIs on user's behalf

Key terms:
  - Resource Owner: The user
  - Client: Your application
  - Authorization Server: Google, GitHub, etc.
  - Resource Server: The API being accessed
  - Scope: What your app can do (e.g., "read:email", "repo")

We won't implement OAuth2 here (it's complex and you should use a
library), but understanding the flow helps you design your auth system.

=============================================================================
 SECURITY CONSIDERATIONS
=============================================================================

Common vulnerabilities to watch for:

1. TIMING ATTACKS: Use crypto/subtle.ConstantTimeCompare for secrets.
   Regular == comparison leaks how many bytes match via timing.

2. TOKEN LEAKAGE: Never log tokens. Never include them in error messages.
   Never send them in URL query parameters (they end up in server logs).

3. WEAK SECRETS: Use crypto/rand for generating secrets, not math/rand.
   math/rand is deterministic and predictable.

4. MISSING HTTPS: Without TLS, tokens are sent in plain text.
   Always require HTTPS in production.

5. XSS (Cross-Site Scripting): If tokens are stored in localStorage,
   XSS can steal them. HttpOnly cookies are safer for web apps.

6. CSRF (Cross-Site Request Forgery): Cookie-based auth is vulnerable.
   Use SameSite cookies and CSRF tokens.

Security headers you should set:
  Strict-Transport-Security: max-age=31536000; includeSubDomains
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Content-Security-Policy: default-src 'self'

=============================================================================
*/

// GenerateRandomKey creates a cryptographically random key of the specified
// byte length. Use this for secrets, API keys, and token generation.
// Never use math/rand for security-sensitive randomness.
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// TimingSafeEqual compares two strings in constant time.
// This prevents timing attacks where an attacker measures how long
// comparison takes to guess secrets character by character.
func TimingSafeEqual(a, b string) bool {
	// subtle.ConstantTimeCompare requires equal lengths to be constant time.
	// If lengths differ, we still want constant-time behavior, so we
	// compare against a dummy to maintain timing consistency.
	if len(a) != len(b) {
		// Still do a comparison to keep timing consistent
		subtle.ConstantTimeCompare([]byte(a), []byte(a))
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// DemoSecurityPractices shows important security considerations.
func DemoSecurityPractices() {
	fmt.Println("=== Security Practices ===")
	fmt.Println()

	// Generate a random key (use crypto/rand, not math/rand!)
	key, _ := GenerateRandomKey(32)
	fmt.Printf("Random key (32 bytes): %x\n", key)

	// Timing-safe comparison
	secret := "my-api-key-abc123"
	fmt.Printf("Timing-safe equal (same): %v\n", TimingSafeEqual(secret, secret))
	fmt.Printf("Timing-safe equal (diff): %v\n", TimingSafeEqual(secret, "wrong-key"))

	fmt.Println()
	fmt.Println("Always use:")
	fmt.Println("  crypto/rand — for generating secrets")
	fmt.Println("  crypto/subtle — for comparing secrets")
	fmt.Println("  crypto/hmac — for message authentication")
	fmt.Println("  NEVER math/rand for security!")
}
