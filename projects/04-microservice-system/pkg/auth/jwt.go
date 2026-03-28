package auth

import (
	"errors"
	"time"
)

// JWT utilities for the microservice system.
//
// For this learning project, implement a simple JWT-like token system.
// You can use a real JWT library (github.com/golang-jwt/jwt/v5) or build
// a simplified version to understand the concepts.
//
// A JWT has three parts: header, payload, signature.
// The payload contains claims like user ID and expiration time.

// Claims holds the data encoded in a token.
type Claims struct {
	UserID    string    `json:"sub"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
}

// Common errors.
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// TokenGenerator creates and validates tokens.
//
// TODO: Implement these methods. You can either:
//   - Use a real JWT library (recommended for production)
//   - Build a simplified HMAC-based token for learning
type TokenGenerator struct {
	secret     string
	expiration time.Duration
}

// NewTokenGenerator creates a TokenGenerator with the given secret and
// token lifetime.
func NewTokenGenerator(secret string, expiration time.Duration) *TokenGenerator {
	return &TokenGenerator{
		secret:     secret,
		expiration: expiration,
	}
}

// Generate creates a new token for the given user ID.
//
// TODO: Implement token generation.
func (g *TokenGenerator) Generate(userID string) (string, error) {
	_ = userID
	return "", errors.New("not implemented")
}

// Validate parses a token string and returns the claims if valid.
//
// TODO: Implement token validation. Check signature and expiration.
func (g *TokenGenerator) Validate(tokenStr string) (*Claims, error) {
	_ = tokenStr
	return nil, errors.New("not implemented")
}
