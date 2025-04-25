package handlers

import (
	"github.com/adshao/go-binance/v2"
	"github.com/gin-gonic/gin"
	"go-blockchain/common/models"
	"net/http"
)

/**
 * @File: binance_handler.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 下午12:51
 * @Software: GoLand
 * @Version:  1.0
 */

// GetAllBinanceTickers 獲取 Binance 上所有交易對的當前價格(市場資料不需要KEY)
func GetAllBinanceTickers(c *gin.Context) {
	// 幣安API KEY 登入後到帳戶(側邊欄)的API管理申請
	apiKey := ""
	secretKey := ""

	client := binance.NewClient(apiKey, secretKey)

	prices, err := client.NewListPricesService().Do(c) // 不指定 Symbol
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "獲取 Binance 所有價格失敗",
			MsgDetail:  err.Error(),
			Data:       nil,
		})
		return
	}

	var data = make(map[string]string)
	for _, p := range prices {
		data[p.Symbol] = p.Price
	}

	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        "獲取 Binance 所有價格成功",
		MsgDetail:  "",
		Data:       data,
	})
}

// GetBinanceTickerPrice 獲取指定交易對在 Binance 上的當前價格(市場資料不需要KEY)
func GetBinanceTickerPrice(c *gin.Context) {
	// 幣安API KEY 登入後到帳戶(側邊欄)的API管理申請
	apiKey := ""
	secretKey := ""

	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "請求錯誤",
			MsgDetail:  "必須提供交易對代碼 (symbol)",
			Data:       nil,
		})
		return
	}

	client := binance.NewClient(apiKey, secretKey)

	prices, err := client.NewListPricesService().Symbol(symbol).Do(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "獲取 Binance 價格失敗",
			MsgDetail:  err.Error(),
			Data:       nil,
		})
		return
	}

	if len(prices) > 0 {
		c.JSON(http.StatusOK, models.JsonResult{
			StatusCode: "200",
			Msg:        "獲取 Binance 價格成功",
			MsgDetail:  "",
			Data: map[string]string{
				"symbol": prices[0].Symbol,
				"price":  prices[0].Price,
			},
		})
	} else {
		c.JSON(http.StatusNotFound, models.JsonResult{
			StatusCode: "404",
			Msg:        "未找到該交易對的價格",
			MsgDetail:  "",
			Data:       nil,
		})
	}
}
