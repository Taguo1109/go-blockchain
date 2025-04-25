package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-blockchain/common/database"
	"go-blockchain/common/middlewares"
	"go-blockchain/internal/user/routes"
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

	// 想關Debug資訊再打開
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// 載入 .env 檔案
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// 會有Gin框架日誌輸出，若是gin.New()則沒有
	router := gin.Default()

	// 註冊全域錯誤處理中間件
	router.Use(middlewares.GlobalErrorHandler())

	// 初始化資料庫
	_, err = database.InitDB("mysql")
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// 註冊路由
	routes.SetupAuthRoutes(router.Group("/auth"))
	routes.SetupUserRoutes(router.Group("/users"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		panic("Failed to run server: " + err.Error())
	}
}
