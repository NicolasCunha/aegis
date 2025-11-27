package permission

import (
	"testing"
	"time"
)

// TestCreatePermission tests permission creation
func TestCreatePermission(t *testing.T) {
	name := "users:read"
	description := "Read users permission"
	createdBy := "system"
	
	permission := CreatePermission(name, description, createdBy)
	
	if permission == nil {
		t.Fatal("CreatePermission should not return nil")
	}
	if permission.Name != name {
		t.Errorf("Expected name %s, got %s", name, permission.Name)
	}
	if permission.Description != description {
		t.Errorf("Expected description %s, got %s", description, permission.Description)
	}
	if permission.CreatedBy != createdBy {
		t.Errorf("Expected createdBy %s, got %s", createdBy, permission.CreatedBy)
	}
	if permission.UpdatedBy != createdBy {
		t.Errorf("Expected updatedBy %s, got %s", createdBy, permission.UpdatedBy)
	}
	if permission.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if permission.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

// TestCreatePermission_EmptyDescription tests creating a permission with empty description
func TestCreatePermission_EmptyDescription(t *testing.T) {
	permission := CreatePermission("users:read", "", "system")
	
	if permission.Description != "" {
		t.Errorf("Expected empty description, got %s", permission.Description)
	}
}

// TestCreatePermission_Timestamps tests that creation timestamps are properly set
func TestCreatePermission_Timestamps(t *testing.T) {
	before := time.Now()
	permission := CreatePermission("users:read", "Read users", "system")
	after := time.Now()
	
	if permission.CreatedAt.Before(before) || permission.CreatedAt.After(after) {
		t.Error("CreatedAt should be set to current time")
	}
	if permission.UpdatedAt.Before(before) || permission.UpdatedAt.After(after) {
		t.Error("UpdatedAt should be set to current time")
	}
}

// TestUpdate tests updating a permission's description
func TestUpdate(t *testing.T) {
	permission := CreatePermission("users:read", "Old description", "system")
	oldUpdatedAt := permission.UpdatedAt
	oldCreatedAt := permission.CreatedAt
	
	time.Sleep(1 * time.Millisecond) // Ensure timestamp difference
	newDescription := "New description"
	updatedBy := "admin"
	permission.Update(newDescription, updatedBy)
	
	if permission.Description != newDescription {
		t.Errorf("Expected description %s, got %s", newDescription, permission.Description)
	}
	if permission.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, permission.UpdatedBy)
	}
	if !permission.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
	if permission.CreatedAt != oldCreatedAt {
		t.Error("CreatedAt should not change on update")
	}
	if permission.CreatedBy != "system" {
		t.Error("CreatedBy should not change on update")
	}
}

// TestUpdate_EmptyDescription tests updating with empty description
func TestUpdate_EmptyDescription(t *testing.T) {
	permission := CreatePermission("users:read", "Original description", "system")
	
	permission.Update("", "admin")
	
	if permission.Description != "" {
		t.Errorf("Expected empty description, got %s", permission.Description)
	}
}

// TestUpdate_MultipleUpdates tests multiple consecutive updates
func TestUpdate_MultipleUpdates(t *testing.T) {
	permission := CreatePermission("users:read", "Description 1", "system")
	
	permission.Update("Description 2", "user1")
	if permission.Description != "Description 2" {
		t.Error("First update failed")
	}
	if permission.UpdatedBy != "user1" {
		t.Error("UpdatedBy not set on first update")
	}
	
	time.Sleep(1 * time.Millisecond)
	lastUpdate := permission.UpdatedAt
	
	permission.Update("Description 3", "user2")
	if permission.Description != "Description 3" {
		t.Error("Second update failed")
	}
	if permission.UpdatedBy != "user2" {
		t.Error("UpdatedBy not set on second update")
	}
	if !permission.UpdatedAt.After(lastUpdate) {
		t.Error("UpdatedAt should change on each update")
	}
}

// TestUpdate_SameDescription tests updating with same description
func TestUpdate_SameDescription(t *testing.T) {
	description := "Same description"
	permission := CreatePermission("users:read", description, "system")
	oldUpdatedAt := permission.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	permission.Update(description, "admin")
	
	// Should still update timestamps even if description is the same
	if !permission.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated even with same description")
	}
}

// TestPermission_NameImmutable tests that permission name cannot be changed after creation
func TestPermission_NameImmutable(t *testing.T) {
	permission := CreatePermission("users:read", "Read users", "system")
	originalName := permission.Name
	
	permission.Update("Updated description", "admin")
	
	if permission.Name != originalName {
		t.Error("Permission name should not change during updates")
	}
}
