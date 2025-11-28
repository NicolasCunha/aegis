package token

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryBlacklist_Add(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	jti := "test-jti-123"
	expiresAt := time.Now().Add(1 * time.Hour)
	
	bl.Add(jti, expiresAt)
	
	if !bl.IsBlacklisted(jti) {
		t.Errorf("Expected token to be blacklisted after Add")
	}
}

func TestMemoryBlacklist_Add_UpdatesExpiration(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	jti := "test-jti-456"
	firstExpiration := time.Now().Add(1 * time.Hour)
	secondExpiration := time.Now().Add(2 * time.Hour)
	
	// Add token with first expiration
	bl.Add(jti, firstExpiration)
	
	// Add same token with different expiration
	bl.Add(jti, secondExpiration)
	
	// Token should still be blacklisted
	if !bl.IsBlacklisted(jti) {
		t.Errorf("Expected token to remain blacklisted after update")
	}
	
	// Verify only one entry exists
	if bl.Size() != 1 {
		t.Errorf("Expected size 1, got %d", bl.Size())
	}
}

func TestMemoryBlacklist_IsBlacklisted_False(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	jti := "non-existent-token"
	
	if bl.IsBlacklisted(jti) {
		t.Errorf("Expected non-existent token to not be blacklisted")
	}
}

func TestMemoryBlacklist_Cleanup_RemovesExpired(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	// Add expired token
	expiredJti := "expired-token"
	expiredTime := time.Now().Add(-1 * time.Hour)
	bl.Add(expiredJti, expiredTime)
	
	// Add valid token
	validJti := "valid-token"
	validTime := time.Now().Add(1 * time.Hour)
	bl.Add(validJti, validTime)
	
	// Verify both tokens are blacklisted
	if !bl.IsBlacklisted(expiredJti) {
		t.Errorf("Expected expired token to be blacklisted before cleanup")
	}
	if !bl.IsBlacklisted(validJti) {
		t.Errorf("Expected valid token to be blacklisted before cleanup")
	}
	
	// Run cleanup
	bl.Cleanup()
	
	// Verify expired token is removed
	if bl.IsBlacklisted(expiredJti) {
		t.Errorf("Expected expired token to be removed after cleanup")
	}
	
	// Verify valid token remains
	if !bl.IsBlacklisted(validJti) {
		t.Errorf("Expected valid token to remain after cleanup")
	}
	
	// Verify size is correct
	if bl.Size() != 1 {
		t.Errorf("Expected size 1 after cleanup, got %d", bl.Size())
	}
}

func TestMemoryBlacklist_Cleanup_EmptyBlacklist(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	// Run cleanup on empty blacklist (should not panic)
	bl.Cleanup()
	
	if bl.Size() != 0 {
		t.Errorf("Expected size 0, got %d", bl.Size())
	}
}

func TestMemoryBlacklist_Size(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	if bl.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", bl.Size())
	}
	
	// Add multiple tokens
	bl.Add("token-1", time.Now().Add(1*time.Hour))
	bl.Add("token-2", time.Now().Add(1*time.Hour))
	bl.Add("token-3", time.Now().Add(1*time.Hour))
	
	if bl.Size() != 3 {
		t.Errorf("Expected size 3, got %d", bl.Size())
	}
	
	// Add duplicate (should not increase size)
	bl.Add("token-1", time.Now().Add(2*time.Hour))
	
	if bl.Size() != 3 {
		t.Errorf("Expected size to remain 3 after duplicate, got %d", bl.Size())
	}
}

func TestMemoryBlacklist_Concurrency(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	// Number of concurrent operations
	numGoroutines := 100
	numOperations := 50
	
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations
	
	// Concurrent Add operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				jti := time.Now().Format("add-token-%d-%d")
				bl.Add(jti, time.Now().Add(1*time.Hour))
			}
		}(i)
	}
	
	// Concurrent IsBlacklisted operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				jti := time.Now().Format("check-token-%d-%d")
				bl.IsBlacklisted(jti)
			}
		}(i)
	}
	
	// Concurrent Size and Cleanup operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				bl.Size()
				if j%10 == 0 {
					bl.Cleanup()
				}
			}
		}(i)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	
	// Test passes if no race conditions occurred
	t.Logf("Concurrent operations completed successfully. Final size: %d", bl.Size())
}

func TestMemoryBlacklist_Cleanup_PartialExpiration(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	now := time.Now()
	
	// Add tokens with various expiration times
	bl.Add("token-expired-1", now.Add(-2*time.Hour))
	bl.Add("token-expired-2", now.Add(-1*time.Hour))
	bl.Add("token-valid-1", now.Add(1*time.Hour))
	bl.Add("token-valid-2", now.Add(2*time.Hour))
	bl.Add("token-valid-3", now.Add(3*time.Hour))
	
	// Verify all are blacklisted
	if bl.Size() != 5 {
		t.Errorf("Expected size 5 before cleanup, got %d", bl.Size())
	}
	
	// Run cleanup
	bl.Cleanup()
	
	// Verify expired tokens are removed
	if bl.IsBlacklisted("token-expired-1") {
		t.Errorf("Expected token-expired-1 to be removed")
	}
	if bl.IsBlacklisted("token-expired-2") {
		t.Errorf("Expected token-expired-2 to be removed")
	}
	
	// Verify valid tokens remain
	if !bl.IsBlacklisted("token-valid-1") {
		t.Errorf("Expected token-valid-1 to remain")
	}
	if !bl.IsBlacklisted("token-valid-2") {
		t.Errorf("Expected token-valid-2 to remain")
	}
	if !bl.IsBlacklisted("token-valid-3") {
		t.Errorf("Expected token-valid-3 to remain")
	}
	
	// Verify final size
	if bl.Size() != 3 {
		t.Errorf("Expected size 3 after cleanup, got %d", bl.Size())
	}
}

func TestMemoryBlacklist_ExpirationBoundary(t *testing.T) {
	bl := NewMemoryBlacklist()
	
	now := time.Now()
	
	// Add token expiring exactly now
	jti := "boundary-token"
	bl.Add(jti, now)
	
	// Token should be blacklisted before cleanup
	if !bl.IsBlacklisted(jti) {
		t.Errorf("Expected token to be blacklisted before cleanup")
	}
	
	// Small delay to ensure time passes
	time.Sleep(10 * time.Millisecond)
	
	// Run cleanup
	bl.Cleanup()
	
	// Token should be removed as it's expired
	if bl.IsBlacklisted(jti) {
		t.Errorf("Expected token to be removed after cleanup (expired at boundary)")
	}
}

func TestGlobalBlacklist_Initialization(t *testing.T) {
	// Save original global blacklist
	original := GlobalBlacklist
	defer func() {
		GlobalBlacklist = original
	}()
	
	// Test initialization
	bl := NewMemoryBlacklist()
	InitializeBlacklist(bl)
	
	if GlobalBlacklist == nil {
		t.Errorf("Expected GlobalBlacklist to be initialized")
	}
	
	// Verify it's the same instance
	jti := "global-test-token"
	GlobalBlacklist.Add(jti, time.Now().Add(1*time.Hour))
	
	if !bl.IsBlacklisted(jti) {
		t.Errorf("Expected global blacklist to be same instance as initialized blacklist")
	}
}
