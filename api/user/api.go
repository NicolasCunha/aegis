package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
	userService "nfcunha/aegis/domain/user"
)

func RegisterApi(router *gin.Engine) {
	router.GET("/", listUsers)
}

func listUsers(c *gin.Context) {
	users := userService.ListUsers()
	c.JSON(http.StatusOK, users)
}