package routes

import (
	"github.com/gin-gonic/gin"
	"go-blockchain/internal/blockchain_info/handlers"
)

/**
 * @File: binance_routes.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 下午1:02
 * @Software: GoLand
 * @Version:  1.0
 */

func SetupRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/binance/price/:symbol", handlers.GetBinanceTickerPrice)
	routerGroup.GET("/binance/prices", handlers.GetAllBinanceTickers)
	routerGroup.GET("/binance/kLines/:symbol", handlers.GetBinanceKLines)
	// 新增獲取兩年 K 線數據並保存的路由
	routerGroup.GET("/binance/kLines/twoYears/:symbol", handlers.FetchTwoYearsKLinesAndSave)
}
