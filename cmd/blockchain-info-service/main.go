package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-blockchain/common/middlewares"
	"go-blockchain/internal/blockchain_info/routes"
	"os"
)

/**
 * @File: main.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午10:23
 * @Software: GoLand
 * @Version:  1.0
 */

func main() {
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	router := gin.Default()
	router.Use(middlewares.GlobalErrorHandler())

	apiGroup := router.Group("/api")
	routes.SetupRoutes(apiGroup) // 確保這裡傳遞了 routerGroup
	// 在這裡添加區塊鏈資訊服務的路由和處理函數

	port := os.Getenv("BLOCKCHAIN_INFO_SERVICE_PORT")
	if port == "" {
		port = "8083" // 區塊鏈資訊服務預設使用 8083 埠
	}

	if err := router.Run(":" + port); err != nil {
		fmt.Println("Failed to run blockchain info service:", err)
	}

}
