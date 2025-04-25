package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/gin-gonic/gin"
	"go-blockchain/common/models"
	"net/http"
	"os"
	"strconv"
	"time"
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

	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        "獲取 Binance 所有價格成功",
		MsgDetail:  "",
		Data:       prices,
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
			Msg:        "獲取 " + symbol + " 價格成功",
			MsgDetail:  "",
			Data:       prices[0],
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

// GetBinanceKLines 獲取指定交易對的 K 線數據
func GetBinanceKLines(c *gin.Context) {
	symbol := c.Param("symbol")
	interval := c.Query("interval") // 查詢參數，例如 1m, 5m, 1h, 1d 等
	limitStr := c.Query("limit")    // 查詢參數，限制返回的 K 線數量

	if symbol == "" || interval == "" {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "請求錯誤",
			MsgDetail:  "必須提供交易對代碼 (symbol) 和時間間隔 (interval)",
			Data:       nil,
		})
		return
	}

	limit := 100 // 預設限制
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	client := binance.NewClient(apiKey, secretKey)

	kLines, err := client.NewKlinesService().Symbol(symbol).Interval(interval).Limit(limit).Do(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "獲取 Binance-" + symbol + "K 線數據失敗",
			MsgDetail:  err.Error(),
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        "獲取 Binance K 線數據成功",
		MsgDetail:  "",
		Data:       kLines, // 直接返回 kLines 陣列
	})
}

// FetchTwoYearsKLinesAndSave 獲取指定交易對近兩年的 K 線數據並保存到 JSON 檔案
func FetchTwoYearsKLinesAndSave(c *gin.Context) {
	symbol := c.Param("symbol")
	interval := c.Query("interval") // 查詢參數，例如 1m, 5m, 1h, 1d 等

	if symbol == "" || interval == "" {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "請求錯誤",
			MsgDetail:  "必須提供交易對代碼 (symbol) 和時間間隔 (interval)",
			Data:       nil,
		})
		return
	}

	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	client := binance.NewClient(apiKey, secretKey)

	currentTime := time.Now()
	twoYearsAgo := currentTime.AddDate(-2, 0, 0)
	startTime := twoYearsAgo.Unix() * 1000 // 毫秒
	endTime := currentTime.Unix() * 1000   // 毫秒
	limit := 1000

	var allKLines []binance.Kline
	for {
		kLines, err := client.NewKlinesService().Symbol(symbol).Interval(interval).StartTime(startTime).EndTime(endTime).Limit(limit).Do(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.JsonResult{
				StatusCode: "500",
				Msg:        fmt.Sprintf("獲取 %s %s K 線數據失敗", symbol, interval),
				MsgDetail:  err.Error(),
				Data:       nil,
			})
			return
		}

		if len(kLines) == 0 {
			break
		}

		for _, kline := range kLines {
			allKLines = append(allKLines, binance.Kline{
				OpenTime:                 kline.OpenTime,
				Open:                     kline.Open,
				High:                     kline.High,
				Low:                      kline.Low,
				Close:                    kline.Close,
				Volume:                   kline.Volume,
				CloseTime:                kline.CloseTime,
				QuoteAssetVolume:         kline.QuoteAssetVolume,
				TradeNum:                 kline.TradeNum,
				TakerBuyBaseAssetVolume:  kline.TakerBuyBaseAssetVolume,
				TakerBuyQuoteAssetVolume: kline.TakerBuyQuoteAssetVolume,
			})
		}

		// 設定下一次請求的起始時間為當前批次最後一筆 K 線的結束時間
		startTime = kLines[len(kLines)-1].CloseTime + 1
		if startTime >= endTime {
			break
		}

		time.Sleep(time.Millisecond * 100) // 避免過於頻繁的請求
	}

	// 將數據保存到 JSON 檔案
	filename := fmt.Sprintf("./%s_%s_two_years.json", symbol, interval)
	file, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "創建 JSON 檔案失敗",
			MsgDetail:  err.Error(),
			Data:       nil,
		})
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(allKLines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "將數據寫入 JSON 檔案失敗",
			MsgDetail:  err.Error(),
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        fmt.Sprintf("成功獲取 %s %s 近兩年 K 線數據並保存到 %s", symbol, interval, filename),
		MsgDetail:  "",
		Data:       filename,
	})
}
