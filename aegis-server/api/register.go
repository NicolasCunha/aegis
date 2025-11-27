package api

import (
	"os"
	"log"
	"github.com/gin-gonic/gin"
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
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "aegis",
			"message": "Service is up and running",
		})
	})
	
	// Register API routes
	userApi.RegisterApi(router)
	roleApi.RegisterApi(router)
	permissionApi.RegisterApi(router)
	
	err := router.Run(getServerPort())
	if err != nil {
		log.Println("Failed to start server:", err)
	}
}