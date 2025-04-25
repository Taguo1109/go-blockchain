package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-blockchain/common/database"
	"go-blockchain/common/middlewares"
	"go-blockchain/internal/blockchain_info/routes"
	"log"
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

	// 初始化資料庫
	_, err = database.InitDB("mysql")
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	apiGroup := router.Group("/api")
	routes.SetupRoutes(apiGroup) // 確保這裡傳遞了 routerGroup
	// 在這裡添加區塊鏈資訊服務的路由和處理函數

	port := os.Getenv("BLOCKCHAIN_INFO_SERVICE_PORT")
	if port == "" {
		port = "8083" // 區塊鏈資訊服務預設使用 8083 埠
	}

	log.Printf("Blockchain Info Service listening on port :%s", port) // 添加這行日誌

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run blockchain info service: %v", err) // 使用 log.Fatalf 輸出錯誤並退出
	}

}
