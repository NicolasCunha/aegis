package api

import (
	"os"
	"log"
	"github.com/gin-gonic/gin"
	authApi "nfcunha/aegis/api/auth"
	userApi "nfcunha/aegis/api/user"
	roleApi "nfcunha/aegis/api/role"
	permissionApi "nfcunha/aegis/api/permission"
)

const DEFAULT_SERVER_PORT = ":8080"

func getServerPort() string {
	envPort := os.Getenv("AEGIS_SERVER_PORT")
	if envPort != "" {
		log.Println("Using server port from environment: ", envPort)
		return ":" + envPort
	}
	log.Println("Using default server port: ", DEFAULT_SERVER_PORT)
	log.Println("To change the port, set the 'AEGIS_SERVER_PORT' environment variable.")
	return DEFAULT_SERVER_PORT
}

func RegisterApis() {
	router := gin.Default()
	
	// Create aegis context path group
	aegis := router.Group("/aegis")
	
	// Health check endpoint
	aegis.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "aegis",
			"message": "Service is up and running",
		})
	})
	
	// Register API routes under /aegis context path
	authApi.RegisterApi(aegis)
	userApi.RegisterApi(aegis)
	roleApi.RegisterApi(aegis)
	permissionApi.RegisterApi(aegis)
	
	err := router.Run(getServerPort())
	if err != nil {
		log.Println("Failed to start server:", err)
	}
}