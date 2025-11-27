// Package role provides HTTP REST API endpoints for role management operations.
// Supports creating, listing, retrieving, updating, and deleting roles.
package role

import (
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	roleService "nfcunha/aegis/domain/role"
)

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Description string `json:"description"`
}

type RoleResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

// RegisterApi registers all role-related HTTP routes with the Gin router.
// Endpoints include create, list, get, update, and delete.
//
// Parameters:
//   - router: The Gin engine to register routes with
func RegisterApi(router *gin.Engine) {
	roles := router.Group("/roles")
	{
		roles.POST("", createRole)
		roles.GET("", listRoles)
		roles.GET("/:name", getRole)
		roles.PUT("/:name", updateRole)
		roles.DELETE("/:name", deleteRole)
	}
}

func createRole(c *gin.Context) {
	log.Println("POST /roles - Create role request received")
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role already exists
	if roleService.ExistsRoleByName(req.Name) {
		log.Printf("Role already exists: %s", req.Name)
		c.JSON(http.StatusConflict, gin.H{"error": "role already exists"})
		return
	}

	// Create role
	role := roleService.CreateRole(req.Name, req.Description, "system")
	roleService.PersistRole(role)

	log.Printf("Role created successfully: %s", role.Name)
	c.JSON(http.StatusCreated, toRoleResponse(role))
}

func listRoles(c *gin.Context) {
	log.Println("GET /roles - List roles request received")
	roles := roleService.ListRoles()
	response := make([]RoleResponse, len(roles))
	for i, role := range roles {
		response[i] = toRoleResponse(role)
	}
	log.Printf("Returning %d roles", len(response))
	c.JSON(http.StatusOK, response)
}

func getRole(c *gin.Context) {
	name := c.Param("name")
	log.Printf("GET /roles/%s - Get role request received", name)

	role := roleService.GetRoleByName(name)
	if role == nil {
		log.Printf("Role not found: %s", name)
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	log.Printf("Returning role: %s", name)
	c.JSON(http.StatusOK, toRoleResponse(role))
}

func updateRole(c *gin.Context) {
	name := c.Param("name")

	role := roleService.GetRoleByName(name)
	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.Update(req.Description, "system")
	roleService.PersistRole(role)

	c.JSON(http.StatusOK, toRoleResponse(role))
}

func deleteRole(c *gin.Context) {
	name := c.Param("name")
	log.Printf("DELETE /roles/%s - Delete role request received", name)

	role := roleService.GetRoleByName(name)
	if role == nil {
		log.Printf("Role not found: %s", name)
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	roleService.DeleteRole(name)

	log.Printf("Role deleted: %s", name)
	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}

// toRoleResponse converts a domain Role model to an API RoleResponse.
//
// Parameters:
//   - role: The domain role to convert
//
// Returns:
//   - RoleResponse containing role data for API responses
func toRoleResponse(role *roleService.Role) RoleResponse {
	return RoleResponse{
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		CreatedBy:   role.CreatedBy,
		UpdatedAt:   role.UpdatedAt,
		UpdatedBy:   role.UpdatedBy,
	}
}
