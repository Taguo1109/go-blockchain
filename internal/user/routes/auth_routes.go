package routes

import (
	"github.com/gin-gonic/gin"
	"go-blockchain/internal/user/handlers"
)

/**
 * @File: auth_routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午11:13
 * @Software: GoLand
 * @Version:  1.0
 */

func SetupAuthRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/register", handlers.Register)
	routerGroup.POST("/login", handlers.Login)
	routerGroup.POST("/logout", handlers.Logout)
}
