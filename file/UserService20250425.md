# 微服務：用戶服務 (User Service) 文件

## 1. 服務目標

* 用戶服務負責處理平台用戶的註冊、登入、登出以及用戶帳戶的基本管理功能。它是整個區塊鏈 Web 應用程式的基礎服務，為其他服務提供用戶身份驗證和授權的能力。

## 2. 功能詳述

### 2.1 用戶註冊 (User Registration)

* 允許新用戶透過電子郵件和密碼創建帳戶。
* **驗證:**
  * 電子郵件和密碼不得為空。
  * 檢查資料庫中是否已存在相同電子郵件的用戶。
* **流程:**
  1. 接收包含電子郵件和密碼的請求。
  2. 驗證輸入資料。
  3. 對密碼進行加鹽哈希處理（使用 bcrypt）。
  4. 將用戶資訊（包含哈希後的密碼）儲存到用戶資料庫（MySQL，使用 GORM）。
  5. 返回包含用戶基本資訊的成功響應。

### 2.2 用戶登入 (User Login)

* 允許已註冊的用戶使用其電子郵件和密碼登入系統。
* **流程:**
  1. 接收包含電子郵件和密碼的請求。
  2. 根據提供的電子郵件查詢用戶資料庫（MySQL，使用 GORM）。
  3. 如果找到用戶，則比較輸入的密碼與資料庫中儲存的哈希密碼。
  4. 如果驗證成功，則生成一個 JWT。
  5. 返回包含 JWT 的成功響應。
  6. 如果驗證失敗，則返回錯誤響應。

### 2.3 用戶登出 (User Logout)

* 允許已登入的用戶安全地終止其會話。
* **流程:**
  * 客戶端通常會清除儲存的 JWT。
  * 伺服器端對於無狀態的 JWT 通常不執行特定操作。
  * 返回成功響應。

### 2.4 查詢用戶資訊 (Get User Information)

* 允許經過身份驗證的用戶查詢特定用戶的帳戶資訊（不包含敏感資訊如密碼）。
* **流程:**
  1. 接收包含用戶 ID 的請求，請求頭中需包含有效的 JWT。
  2. 使用 JWT 中間件驗證請求的身份。
  3. 從資料庫中檢索用戶資訊（MySQL，使用 GORM）。
  4. 返回包含用戶基本資訊的成功響應。
  5. 如果找不到用戶或身份驗證失敗，則返回錯誤響應。

## 3. API 設計

所有 API 回應都將採用 `common/models/JsonResult` 結構。

### 3.1 `POST /auth/register` (用戶註冊)

* **請求體 (JSON):**
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123"
    }
    ```
* **回應 (成功 201 Created):**
    ```json
    {
      "status_code": "201",
      "msg": "註冊成功",
      "msg_detail": "",
      "data": {
        "id": "user-uuid-123",
        "email": "user@example.com"
      }
    }
    ```
* **回應 (失敗 400 Bad Request):** 請求體格式錯誤或電子郵件/密碼為空。
* **回應 (失敗 409 Conflict):** 該電子郵件已被註冊。
* **回應 (失敗 500 Internal Server Error):** 伺服器內部錯誤，例如資料庫操作失敗。

### 3.2 `POST /auth/login` (用戶登入)

* **請求體 (JSON):**
    ```json
    {
      "email": "user@example.com",
      "password": "securePassword123"
    }
    ```
* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "登入成功",
      "msg_detail": "",
      "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
      }
    }
    ```
* **回應 (失敗 400 Bad Request):** 請求體格式錯誤。
* **回應 (失敗 401 Unauthorized):** 無效的憑證（電子郵件或密碼錯誤）。
* **回應 (失敗 500 Internal Server Error):** 伺服器內部錯誤，例如資料庫查詢失敗或生成 token 失敗。

### 3.3 `POST /auth/logout` (用戶登出)

* **請求體:** (通常為空)
* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "登出成功",
      "msg_detail": "客戶端應清除 token",
      "data": null
    }
    ```

### 3.4 `GET /users/{id}` (查詢用戶資訊)

* **路徑參數:** `{id}` - 要查詢的用戶 ID。
* **請求頭:** 需要包含有效的 JWT (例如 `Authorization: Bearer <token>`).
* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "查詢成功",
      "msg_detail": "",
      "data": {
        "id": "user-uuid-123",
        "email": "user@example.com",
        "created_at": "2025-04-25T12:00:00+08:00",
        "updated_at": "2025-04-25T12:00:00+08:00"
      }
    }
    ```
* **回應 (失敗 400 Bad Request):** 用戶 ID 格式錯誤。
* **回應 (失敗 401 Unauthorized):** 無效的授權憑證。
* **回應 (失敗 404 Not Found):** 找不到該用戶。
* **回應 (失敗 500 Internal Server Error):** 伺服器內部錯誤，例如資料庫查詢失敗。

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

### 4.2 `internal/user/models/user.go`

```go
package models

import "time"

type User struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"` // 不希望在 JSON 響應中返回密碼
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## 5. 程式碼結構

```
common/
├── cache/
│   └── redis.go        # Redis 初始化 (未來可能使用)
├── database/
│   └── init_db.go      # 資料庫 (MySQL) 初始化 (使用 GORM)
├── middleware/
│   ├── global_error_handler.go # 全域錯誤攔截器
│   └── jwt_auth.go           # JWT 相關邏輯 (生成、驗證、Claims)
└── models/
    └── response.go         # 通用 JSON 回應結構

internal/
└── user/
    ├── database/
    │   └── user_db.go      # User 模型的資料庫操作 (使用 GORM)
    ├── handlers/
    │   ├── auth_handler.go # 用戶註冊、登入、登出處理
    │   └── user_handler.go # 查詢用戶資訊處理
    ├── models/
    │   └── user.go         # User 資料模型
    └── routes/
        ├── auth_routes.go  # 用戶驗證相關路由
        └── user_routes.go  # 用戶資訊查詢相關路由

cmd/
└── user-service/
    └── main.go           # UserService 應用程式入口
```

## 6. 技術考量

* 資料庫: MySQL (使用 `gorm.io/gorm` 和 `gorm.io/driver/mysql` 進行操作)。
* 身份驗證: JWT (JSON Web Tokens) 基於 `github.com/golang-jwt/jwt/v5` 庫，相關邏輯放在 `common/middleware/jwt_auth.go`。
* 密碼儲存: 使用 `golang.org/x/crypto/bcrypt` 庫進行加鹽哈希處理。
* 唯一 ID 生成: 使用 `github.com/google/uuid` 庫生成用戶 ID。
* 環境變數管理: 使用 `github.com/joho/godotenv` 庫載入 .env 檔案。
* 全域錯誤處理: 使用自定義的 Gin 中間件 (`common/middleware/global_error_handler.go`) 捕獲並處理 panic 錯誤。

## 7. 後續步驟

1. 確保 MySQL 資料庫已建立，並且 users 表格已按照提供的 SQL 語句創建。
2. 在 `.env` 檔案中配置正確的資料庫連線資訊 (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`) 和 JWT 簽名金鑰 (`JWT_SECRET`)。
3. 運行 `cmd/user-service/main.go` 啟動服務。
4. 根據 API 設計部分提供的端點和請求體進行 API 測試。