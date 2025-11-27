package role

import (
	"testing"
	"time"
)

// TestCreateRole tests role creation
func TestCreateRole(t *testing.T) {
	name := "admin"
	description := "Administrator role"
	createdBy := "system"
	
	role := CreateRole(name, description, createdBy)
	
	if role == nil {
		t.Fatal("CreateRole should not return nil")
	}
	if role.Name != name {
		t.Errorf("Expected name %s, got %s", name, role.Name)
	}
	if role.Description != description {
		t.Errorf("Expected description %s, got %s", description, role.Description)
	}
	if role.CreatedBy != createdBy {
		t.Errorf("Expected createdBy %s, got %s", createdBy, role.CreatedBy)
	}
	if role.UpdatedBy != createdBy {
		t.Errorf("Expected updatedBy %s, got %s", createdBy, role.UpdatedBy)
	}
	if role.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if role.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

// TestCreateRole_EmptyDescription tests creating a role with empty description
func TestCreateRole_EmptyDescription(t *testing.T) {
	role := CreateRole("admin", "", "system")
	
	if role.Description != "" {
		t.Errorf("Expected empty description, got %s", role.Description)
	}
}

// TestCreateRole_Timestamps tests that creation timestamps are properly set
func TestCreateRole_Timestamps(t *testing.T) {
	before := time.Now()
	role := CreateRole("admin", "Admin role", "system")
	after := time.Now()
	
	if role.CreatedAt.Before(before) || role.CreatedAt.After(after) {
		t.Error("CreatedAt should be set to current time")
	}
	if role.UpdatedAt.Before(before) || role.UpdatedAt.After(after) {
		t.Error("UpdatedAt should be set to current time")
	}
}

// TestUpdate tests updating a role's description
func TestUpdate(t *testing.T) {
	role := CreateRole("admin", "Old description", "system")
	oldUpdatedAt := role.UpdatedAt
	oldCreatedAt := role.CreatedAt
	
	time.Sleep(1 * time.Millisecond) // Ensure timestamp difference
	newDescription := "New description"
	updatedBy := "admin"
	role.Update(newDescription, updatedBy)
	
	if role.Description != newDescription {
		t.Errorf("Expected description %s, got %s", newDescription, role.Description)
	}
	if role.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, role.UpdatedBy)
	}
	if !role.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
	if role.CreatedAt != oldCreatedAt {
		t.Error("CreatedAt should not change on update")
	}
	if role.CreatedBy != "system" {
		t.Error("CreatedBy should not change on update")
	}
}

// TestUpdate_EmptyDescription tests updating with empty description
func TestUpdate_EmptyDescription(t *testing.T) {
	role := CreateRole("admin", "Original description", "system")
	
	role.Update("", "admin")
	
	if role.Description != "" {
		t.Errorf("Expected empty description, got %s", role.Description)
	}
}

// TestUpdate_MultipleUpdates tests multiple consecutive updates
func TestUpdate_MultipleUpdates(t *testing.T) {
	role := CreateRole("admin", "Description 1", "system")
	
	role.Update("Description 2", "user1")
	if role.Description != "Description 2" {
		t.Error("First update failed")
	}
	if role.UpdatedBy != "user1" {
		t.Error("UpdatedBy not set on first update")
	}
	
	time.Sleep(1 * time.Millisecond)
	lastUpdate := role.UpdatedAt
	
	role.Update("Description 3", "user2")
	if role.Description != "Description 3" {
		t.Error("Second update failed")
	}
	if role.UpdatedBy != "user2" {
		t.Error("UpdatedBy not set on second update")
	}
	if !role.UpdatedAt.After(lastUpdate) {
		t.Error("UpdatedAt should change on each update")
	}
}

// TestUpdate_SameDescription tests updating with same description
func TestUpdate_SameDescription(t *testing.T) {
	description := "Same description"
	role := CreateRole("admin", description, "system")
	oldUpdatedAt := role.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	role.Update(description, "admin")
	
	// Should still update timestamps even if description is the same
	if !role.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated even with same description")
	}
}

// TestRole_NameImmutable tests that role name cannot be changed after creation
func TestRole_NameImmutable(t *testing.T) {
	role := CreateRole("admin", "Admin role", "system")
	originalName := role.Name
	
	role.Update("Updated description", "admin")
	
	if role.Name != originalName {
		t.Error("Role name should not change during updates")
	}
}
