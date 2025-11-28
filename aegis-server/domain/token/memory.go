package token

import (
	"sync"
	"time"
)

// MemoryBlacklist implements the Blacklist interface using an in-memory map.
// It provides thread-safe token revocation using sync.RWMutex for concurrent access.
//
// This implementation is suitable for development and single-instance deployments.
// For production with multiple instances, consider using a Redis-backed implementation.
type MemoryBlacklist struct {
	entries map[string]*BlacklistEntry // Map of JTI -> BlacklistEntry
	mu      sync.RWMutex                // Protects concurrent access to entries
}

// NewMemoryBlacklist creates a new in-memory blacklist instance.
//
// Returns:
//   - A new MemoryBlacklist ready for use
func NewMemoryBlacklist() *MemoryBlacklist {
	return &MemoryBlacklist{
		entries: make(map[string]*BlacklistEntry),
	}
}

// Add adds a token to the blacklist by its JTI.
// Thread-safe for concurrent writes.
//
// Parameters:
//   - jti: The unique JWT ID from the token's claims
//   - expiresAt: When the token expires naturally
//
// Returns:
//   - Always returns nil (error interface for future implementations)
func (b *MemoryBlacklist) Add(jti string, expiresAt time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries[jti] = &BlacklistEntry{
		JTI:       jti,
		ExpiresAt: expiresAt,
		RevokedAt: time.Now(),
	}

	return nil
}

// IsBlacklisted checks if a token is currently on the blacklist.
// Thread-safe for concurrent reads.
//
// Parameters:
//   - jti: The JWT ID to check
//
// Returns:
//   - true if the token is blacklisted, false otherwise
func (b *MemoryBlacklist) IsBlacklisted(jti string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.entries[jti]
	return exists
}

// Cleanup removes expired entries from the blacklist.
// Tokens that have expired naturally no longer need to be tracked.
// Thread-safe for concurrent cleanup operations.
//
// Returns:
//   - Number of entries removed
func (b *MemoryBlacklist) Cleanup() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	removed := 0

	// Iterate and remove expired entries
	for jti, entry := range b.entries {
		if entry.ExpiresAt.Before(now) {
			delete(b.entries, jti)
			removed++
		}
	}

	return removed
}

// Size returns the current number of blacklisted tokens.
// Thread-safe for concurrent reads.
//
// Returns:
//   - Number of entries in the blacklist
func (b *MemoryBlacklist) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.entries)
}
