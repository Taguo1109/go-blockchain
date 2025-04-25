# 微服務：區塊鏈資訊服務 (Blockchain Info Service) 文件

## 1. 服務目標

* 區塊鏈資訊服務負責從不同的區塊鏈網路或相關的加密貨幣交易所（目前是 Binance）獲取資訊，並以統一的 API 介面提供給其他服務或使用者。
* 目前主要功能是提供加密貨幣的價格資訊。

## 2. 功能詳述

### 2.1 獲取指定交易對的價格 (Get Ticker Price)

* 允許查詢特定交易對在 Binance 上的當前價格。
* **請求參數:**
    * `symbol` (路徑參數): 交易對代碼 (例如 BTCUSDT)。
* **流程:**
    1. 接收包含交易對代碼的請求。
    2. 使用 Binance API Go 語言套件 (`github.com/adshao/go-binance/v2`) 調用 Binance API 獲取該交易對的最新價格。
    3. 將 Binance API 的回應轉換為標準的 `JsonResult` 格式。
    4. 返回包含交易對代碼和價格的成功響應。
    5. 如果 Binance API 請求失敗或未找到該交易對，則返回錯誤響應。

### 2.2 獲取所有交易對的價格 (Get All Tickers)

* 允許查詢 Binance 上所有可用交易對的當前價格。
* **流程:**
    1. 接收請求。
    2. 使用 Binance API Go 語言套件 (`github.com/adshao/go-binance/v2`) 調用 Binance API 獲取所有交易對的最新價格列表。
    3. 將 Binance API 的回應轉換為一個包含所有交易對及其價格的 `map`，並封裝在標準的 `JsonResult` 格式中。
    4. 返回包含所有交易對和價格的成功響應。
    5. 如果 Binance API 請求失敗，則返回錯誤響應。

## 3. API 設計

所有 API 回應都將採用 `common/models/JsonResult` 結構。

### 3.1 `GET /api/binance/price/{symbol}` (獲取指定交易對的價格)

* **路徑參數:** `{symbol}` - 要查詢的交易對代碼 (例如 BTCUSDT)。
* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "獲取 Binance 價格成功",
      "msg_detail": "",
      "data": {
        "symbol": "BTCUSDT",
        "price": "93168.62000000"
      }
    }
    ```
* **回應 (失敗 400 Bad Request):** 請求中未提供 `symbol`。
* **回應 (失敗 500 Internal Server Error):** 獲取 Binance 價格失敗，查看 `msg_detail` 獲取詳細錯誤資訊。
* **回應 (失敗 404 Not Found):** Binance 上未找到該交易對。

### 3.2 `GET /api/binance/prices` (獲取所有交易對的價格)

* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "獲取 Binance 所有價格成功",
      "msg_detail": "",
      "data": {
        "1000CATBUSD": "0.00091113",
        "1000CATUSDT": "0.00076660",
        "1000CATEUR": "0.29093806",
        "1000CATBTC": "0.00007580",
        "1000DOGEUSDT": "0.00077000",
        // ... 其他交易對和價格
      }
    }
    ```
* **回應 (失敗 500 Internal Server Error):** 獲取 Binance 所有價格失敗，查看 `msg_detail` 獲取詳細錯誤資訊。

## 4. 資料模型 (Model)

### 4.1 `common/models/json_result.go`

```go
package models

type JsonResult struct {
    StatusCode string      `json:"status_code"`
    Msg        interface{} `json:"msg"`
    MsgDetail  string      `json:"msg_detail"`
    Data       interface{} `json:"data"`
}
```

## 5. 程式碼結構

```
common/
├── cache/
│   └── redis.go
├── database/
│   └── init_db.go
├── middleware/
│   ├── global_error_handler.go
│   └── jwt_auth.go
└── models/
    └── json_result.go

internal/
└── blockchain/
    ├── api/
    │   └── binance/
    │       └── ... (使用 [github.com/adshao/go-binance/v2](https://github.com/adshao/go-binance/v2) 套件)
    ├── handlers/
    │   └── binance_handler.go
    └── routes/
        └── routes.go

cmd/
└── blockchain-info-service/
    └── main.go
```

## 6.技術考量

* API 互動: 使用第三方 Go 語言套件 `github.com/adshao/go-binance/v2` 與 Binance API 進行通訊。
* 環境變數: 使用 `github.com/joho/godotenv` 載入 `.env` 檔案以管理 Binance API 金鑰和密鑰。
* 錯誤處理: 使用 Gin 的錯誤處理機制和全域錯誤處理中間件 (`common/middleware/global_error_handler.go`) 處理 API 請求和 Binance API 調用中可能發生的錯誤。
* 回應格式: 所有 API 回應都遵循 `common/models/JsonResult` 結構。

## 7.後續步驟

1. 繼續開發更多從 Binance API 獲取資訊的功能，例如 K 線數據、交易所資訊等。
2. 考慮整合其他區塊鏈網路的資訊來源。
3. 增加錯誤處理和日誌記錄的完善性。
4. 考慮添加緩存機制以減少對外部 API 的頻繁請求。