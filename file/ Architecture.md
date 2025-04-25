# 區塊鏈 Web 專案 - 微服務架構 SA 文件 (草稿)

## 1. 專案目標

* 建立一個基於微服務架構的 Web 平台，讓使用者能夠與區塊鏈技術進行初步互動，目標是幫助開發者（你）熟悉 Web3 的相關概念和技術，並學習微服務的設計與實作。

## 2. 目標用戶

* 對區塊鏈技術感興趣的開發者和學習者。

## 3. 核心功能與微服務劃分

### 3.1 用戶服務 (User Service)

* **功能:** 處理用戶註冊、登入、登出、帳戶管理等。
* **API 端點 (範例):**
    * `POST /users`: 用於用戶註冊，接收電子郵件和密碼。
    * `POST /auth/login`: 用於用戶登入，接收電子郵件和密碼。
    * `POST /auth/logout`: 用於用戶登出。
    * `GET /users/{userId}`: 查詢用戶資訊。
* **資料儲存:** 用戶資料庫 (例如 MySQL)。
* **技術選型:** Golang, Gin。

### 3.2 錢包服務 (Wallet Service)

* **功能:** 處理錢包的生成、地址的儲存與管理，以及與用戶的關聯。
* **API 端點 (範例):**
    * `POST /wallets`: 為指定用戶生成新的錢包地址。
        * **請求體:** `{"userId": "..."}`
        * **回應:** 返回生成的錢包地址。
    * `GET /wallets/{walletAddress}`: 查詢錢包的詳細資訊 (可能包含關聯的用戶 ID)。
* **資料儲存:** 錢包地址資料庫 (例如 MySQL)。
* **技術選型:** Golang, Gin。

### 3.3 區塊鏈資訊服務 (Blockchain Info Service)

* **功能:** 與不同的區塊鏈網路互動，提供查詢區塊高度、交易資訊、錢包餘額等功能。
* **API 端點 (範例):**
    * `GET /blockchain/{network}/latest_block`: 查詢指定網路的最新區塊高度。
        * **路徑參數:** `network` (例如 "ethereum-goerli")
        * **回應:** 返回最新區塊高度。
    * `GET /blockchain/{network}/transaction/{txHash}`: 查詢指定網路中交易哈希的詳細資訊。
        * **路徑參數:** `network`, `txHash`
        * **回應:** 返回交易詳細資訊。
    * `GET /blockchain/{network}/balance/{walletAddress}`: 查詢指定網路中錢包地址的餘額。
        * **路徑參數:** `network`, `walletAddress`
        * **回應:** 返回錢包餘額。
* **區塊鏈互動:**
    * 可能需要針對不同的區塊鏈網路使用不同的 SDK 或 API 進行互動 (例如，Go-Ethereum 庫用於以太坊，或第三方區塊鏈數據服務 API)。
* **技術選型:** Golang, Gin。

## 4. 技術架構 (更新)

* **後端:** 採用微服務架構，主要包含用戶服務、錢包服務和區塊鏈資訊服務。每個服務使用 Golang 和 Gin 框架開發。
* **資料庫:**
    * 用戶服務使用用戶資料庫 (MySQL)。
    * 錢包服務使用錢包地址資料庫 (MySQL)。
    * 區塊鏈資訊服務本身可能不直接儲存業務資料，而是依賴與區塊鏈節點或第三方服務的互動。
* **服務發現與協調 (Service Discovery & Coordination):** (在這個初步階段，我們可以先簡化，不引入複雜的服務發現機制。後續如果服務數量增多，可以考慮 Consul, etcd 或 Kubernetes 等方案。)
* **API 網關 (API Gateway):** (作為所有外部請求的入口點，將請求路由到對應的微服務。可以考慮使用 Kong, Tyk 或 Nginx 等。)
* **通訊方式:** 微服務之間可以採用 RESTful API 或 gRPC 進行同步通訊。由於你對 gRPC 有所了解，後續可以考慮使用 gRPC 來提高服務間的通訊效率。

## 5. API 設計 (更新 - 針對 API 網關)

來自前端的請求將首先到達 API 網關，然後由網關路由到相應的微服務。以下是通過 API 網關暴露的端點 (範例)：

### 5.1 用戶管理

* `POST /api/users`: 路由到用戶服務的註冊接口。
* `POST /api/auth/login`: 路由到用戶服務的登錄接口。
* `POST /api/auth/logout`: 路由到用戶服務的登出接口。
* `GET /api/users/{userId}`: 路由到用戶服務的查詢用戶信息接口。

### 5.2 錢包管理

* `POST /api/wallets`: 路由到錢包服務的生成錢包接口 (需要在請求中包含用戶 ID)。
* `GET /api/wallets/{walletAddress}`: 路由到錢包服務的查詢錢包信息接口。

### 5.3 區塊鏈資訊查詢

* `GET /api/blockchain/{network}/latest_block`: 路由到區塊鏈資訊服務的查詢最新區塊高度接口。
* `GET /api/blockchain/{network}/transaction/{txHash}`: 路由到區塊鏈資訊服務的查詢交易信息接口。
* `GET /api/blockchain/{network}/balance/{walletAddress}`: 路由到區塊鏈資訊服務的查詢錢包餘額接口。

## 6. 後續步驟 (更新)

* 搭建 Golang 開發環境並初始化各個微服務的 Gin 專案。
* 選擇並配置 API 網關。
* 實現用戶服務的 API (註冊、登入、登出)。
* 設計用戶服務的資料庫模型。
* 逐步實現錢包服務和區塊鏈資訊服務的 API，並考慮它們之間以及與用戶服務的交互方式。
* 研究如何與選定的區塊鏈網路互動。
* 設計錢包服務的資料庫模型。

## 7. 目錄結構 (建議)
```
├── cmd/
│   ├── api-gateway/     # API 網關應用程式
│   │   └── main.go
│   ├── user-service/    # 用戶服務應用程式
│   │   └── main.go
│   ├── wallet-service/  # 錢包服務應用程式
│   │   └── main.go
│   └── blockchain-info-service/ # 區塊鏈資訊服務應用程式
│       └── main.go
├── internal/              # 包含應用程式的內部程式碼，不會被其他專案直接引用
│   ├── user/              # 用戶服務的內部邏輯
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── database/
│   │   ├── routes/
│   │   └── config/
│   ├── wallet/            # 錢包服務的內部邏輯
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── database/
│   │   ├── routes/
│   │   └── config/
│   └── blockchain_info/   # 區塊鏈資訊服務的內部邏輯
│       ├── handlers/
│       ├── clients/
│       ├── routes/
│       └── config/
├── common/                  # 可被多個服務共用的程式碼 (例如工具函數、錯誤定義)
│   ├── utils/
│   ├── errors/
│   └── ...
├── docker-compose.yml       # (後續可選) 用於本地開發環境的 Docker Compose 配置
└── README.md
``` 