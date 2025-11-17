package api

import (
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	userApi "nfcunha/aegis/api/user"
)

func getServerPort() string {
	envPort := os.Getenv("SERVER_PORT")
	if envPort != "" {
		fmt.Println("Using server port from environment:", envPort)
		return ":" + envPort
	}
	fmt.Println("Using default server port: 8080")
	return ":8080"
}

func RegisterApis() {
	router := gin.Default()
	userApi.RegisterApi(router)
	err := router.Run(getServerPort())
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}