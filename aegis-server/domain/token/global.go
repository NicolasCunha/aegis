// Package token provides a global blacklist instance for token revocation.
package token

var (
	// GlobalBlacklist is the application-wide token blacklist instance.
	// It is initialized at application startup and used by all validation endpoints.
	GlobalBlacklist Blacklist
)

// InitializeBlacklist initializes the global blacklist with the specified implementation.
// This should be called once during application startup, before any HTTP handlers are registered.
//
// Parameters:
//   - blacklist: The blacklist implementation to use (e.g., MemoryBlacklist or Redis)
func InitializeBlacklist(blacklist Blacklist) {
	GlobalBlacklist = blacklist
}
