package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// 在這裡添加錢包服務的路由和處理函數

	port := os.Getenv("WALLET_SERVICE_PORT")
	if port == "" {
		port = "8082" // 錢包服務預設使用 8082 埠
	}

	if err := router.Run(":" + port); err != nil {
		fmt.Println("Failed to run wallet service:", err)
	}
}
