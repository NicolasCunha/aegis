// Package permission provides HTTP REST API endpoints for permission management operations.
// Supports creating, listing, retrieving, updating, and deleting permissions.
package permission

import (
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	permissionService "nfcunha/aegis/domain/permission"
)

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Description string `json:"description"`
}

type PermissionResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

// RegisterApi registers all permission-related HTTP routes with the Gin router.
// Endpoints include create, list, get, update, and delete.
//
// Parameters:
//   - router: The Gin RouterGroup to register routes with (already under /aegis)
func RegisterApi(router gin.IRouter) {
	permissions := router.Group("/permissions")
	{
		permissions.POST("", createPermission)
		permissions.GET("", listPermissions)
		permissions.GET("/:name", getPermission)
		permissions.PUT("/:name", updatePermission)
		permissions.DELETE("/:name", deletePermission)
	}
}

func createPermission(c *gin.Context) {
	log.Println("POST /permissions - Create permission request received")
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if permission already exists
	if permissionService.ExistsPermissionByName(req.Name) {
		log.Printf("Permission already exists: %s", req.Name)
		c.JSON(http.StatusConflict, gin.H{"error": "permission already exists"})
		return
	}

	// Create permission
	permission := permissionService.CreatePermission(req.Name, req.Description, "system")
	permissionService.PersistPermission(permission)

	log.Printf("Permission created successfully: %s", permission.Name)
	c.JSON(http.StatusCreated, toPermissionResponse(permission))
}

func listPermissions(c *gin.Context) {
	log.Println("GET /permissions - List permissions request received")
	permissions := permissionService.ListPermissions()
	response := make([]PermissionResponse, len(permissions))
	for i, permission := range permissions {
		response[i] = toPermissionResponse(permission)
	}
	log.Printf("Returning %d permissions", len(response))
	c.JSON(http.StatusOK, response)
}

func getPermission(c *gin.Context) {
	name := c.Param("name")
	log.Printf("GET /aegis/permissions/%s - Get permission request received", name)

	permission := permissionService.GetPermissionByName(name)
	if permission == nil {
		log.Printf("Permission not found: %s", name)
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	log.Printf("Returning permission: %s", name)
	c.JSON(http.StatusOK, toPermissionResponse(permission))
}

func updatePermission(c *gin.Context) {
	name := c.Param("name")

	permission := permissionService.GetPermissionByName(name)
	if permission == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission.Update(req.Description, "system")
	permissionService.PersistPermission(permission)

	c.JSON(http.StatusOK, toPermissionResponse(permission))
}

func deletePermission(c *gin.Context) {
	name := c.Param("name")
	log.Printf("DELETE /aegis/permissions/%s - Delete permission request received", name)

	permission := permissionService.GetPermissionByName(name)
	if permission == nil {
		log.Printf("Permission not found: %s", name)
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	permissionService.DeletePermission(name)

	log.Printf("Permission deleted: %s", name)
	c.JSON(http.StatusOK, gin.H{"message": "permission deleted successfully"})
}

// toPermissionResponse converts a domain Permission model to an API PermissionResponse.
//
// Parameters:
//   - permission: The domain permission to convert
//
// Returns:
//   - PermissionResponse containing permission data for API responses
func toPermissionResponse(permission *permissionService.Permission) PermissionResponse {
	return PermissionResponse{
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		CreatedBy:   permission.CreatedBy,
		UpdatedAt:   permission.UpdatedAt,
		UpdatedBy:   permission.UpdatedBy,
	}
}
