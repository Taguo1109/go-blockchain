package routes

import (
	"github.com/gin-gonic/gin"
	"go-blockchain/common/middlewares"
	"go-blockchain/internal/user/handlers"
)

/**
 * @File: user_routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午11:15
 * @Software: GoLand
 * @Version:  1.0
 */

func SetupUserRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/:id", middlewares.AuthMiddleware(), handlers.GetUser)
}
