// Package token provides token blacklist functionality for JWT token revocation.
// Tokens can be revoked before their natural expiration by adding them to the blacklist.
package token

import "time"

// Blacklist defines the interface for managing revoked tokens.
// Implementations must be thread-safe for concurrent access.
//
// The blacklist stores JWT IDs (JTI claims) of revoked tokens along with their
// expiration times. Once a token expires naturally, it can be removed from the
// blacklist to conserve memory/storage.
type Blacklist interface {
	// Add adds a token to the blacklist by its JTI (JWT ID).
	// The expiresAt parameter should match the token's original expiration time.
	//
	// Parameters:
	//   - jti: The unique JWT ID claim from the token
	//   - expiresAt: When the token expires naturally (used for cleanup)
	//
	// Returns:
	//   - Error if the operation fails
	Add(jti string, expiresAt time.Time) error

	// IsBlacklisted checks if a token is currently on the blacklist.
	//
	// Parameters:
	//   - jti: The unique JWT ID claim to check
	//
	// Returns:
	//   - true if the token is blacklisted, false otherwise
	IsBlacklisted(jti string) bool

	// Cleanup removes expired entries from the blacklist.
	// This should be called periodically (e.g., hourly) to prevent memory growth.
	// Tokens that have expired naturally no longer need to be tracked.
	//
	// Returns:
	//   - Number of entries removed
	Cleanup() int

	// Size returns the current number of entries in the blacklist.
	// Useful for monitoring and metrics.
	//
	// Returns:
	//   - Number of blacklisted tokens
	Size() int
}

// BlacklistEntry represents a single entry in the token blacklist.
type BlacklistEntry struct {
	JTI       string    // JWT ID from the token's claims
	ExpiresAt time.Time // When the token expires naturally
	RevokedAt time.Time // When the token was revoked
}
