# 微服務：用戶服務 (User Service) 文件

## 1. 服務目標

* 用戶服務負責處理平台用戶的註冊、登入、登出以及用戶帳戶的基本管理功能。它是整個區塊鏈 Web 應用程式的基礎服務，為其他服務提供用戶身份驗證和授權的能力。

## 2. 功能詳述

### 2.1 用戶註冊 (User Registration)

* 允許新用戶透過電子郵件和密碼創建帳戶。
* **驗證:**
    * 電子郵件格式的有效性（使用正規表達式或其他驗證方法）。
    * 密碼的強度（例如，最小長度、包含大小寫字母、數字和特殊字元等）。
    * 檢查資料庫中是否已存在相同電子郵件的用戶。
* **流程:**
    1. 接收包含電子郵件和密碼的請求。
    2. 驗證輸入資料。
    3. 對密碼進行加鹽哈希處理（例如使用 bcrypt）。
    4. 將用戶資訊（包含哈希後的密碼）儲存到用戶資料庫。
    5. 返回包含用戶基本資訊的成功響應。

### 2.2 用戶登入 (User Login)

* 允許已註冊的用戶使用其電子郵件和密碼登入系統。
* **流程:**
    1. 接收包含電子郵件和密碼的請求。
    2. 根據提供的電子郵件查詢用戶資料庫。
    3. 如果找到用戶，則比較輸入的密碼與資料庫中儲存的哈希密碼。
    4. 如果驗證成功，則生成一個身份驗證憑證（JWT）。
    5. 返回包含 JWT 的成功響應。
    6. 如果驗證失敗，則返回錯誤響應。

### 2.3 用戶登出 (User Logout)

* 允許已登入的用戶安全地終止其會話。
* **流程:**
    1. 接收登出請求（可能需要身份驗證）。
    2. 在客戶端，通常會清除儲存的 JWT。
    3. 在伺服器端，可以選擇性地將 JWT 加入黑名單或使 Session 失效（如果使用 Session）。
    4. 返回成功響應。

### 2.4 查詢用戶資訊 (Get User Information)

* 允許經過身份驗證的用戶或具有特定權限的其他服務查詢特定用戶的帳戶資訊（不包含敏感資訊如密碼）。
* **流程:**
    1. 接收包含用戶 ID 的請求。
    2. 驗證請求的身份（確保請求者有權查詢該用戶資訊）。
    3. 從資料庫中檢索用戶資訊。
    4. 返回包含用戶基本資訊的成功響應。
    5. 如果找不到用戶或請求者無權訪問，則返回錯誤響應。

## 3. API 設計

所有 API 回應都將採用 `common/models/JsonResult` 結構。

### 3.1 `POST /users` (用戶註冊)

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
        "email": "user@example.com",
        "created_at": "2025-04-25T10:35:00Z"
      }
    }
    ```
* **回應 (失敗 400 Bad Request):**
    ```json
    {
      "status_code": "400",
      "msg": "註冊失敗",
      "msg_detail": "無效的電子郵件格式或密碼強度不足",
      "data": null
    }
    ```
* **回應 (失敗 409 Conflict):**
    ```json
    {
      "status_code": "409",
      "msg": "註冊失敗",
      "msg_detail": "電子郵件地址已存在",
      "data": null
    }
    ```

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
* **回應 (失敗 401 Unauthorized):**
    ```json
    {
      "status_code": "401",
      "msg": "登入失敗",
      "msg_detail": "無效的憑證",
      "data": null
    }
    ```

### 3.3 `POST /auth/logout` (用戶登出)

* **請求體:** (通常為空)
* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "登出成功",
      "msg_detail": "",
      "data": null
    }
    ```
* **回應 (失敗 401 Unauthorized):**
    ```json
    {
      "status_code": "401",
      "msg": "登出失敗",
      "msg_detail": "未經授權",
      "data": null
    }
    ```

### 3.4 `GET /users/{userId}` (查詢用戶資訊)

* **路徑參數:** `{userId}` - 要查詢的用戶 ID。
* **請求頭:** 需要包含身份驗證憑證 (例如 `Authorization: Bearer <token>`).
* **回應 (成功 200 OK):**
    ```json
    {
      "status_code": "200",
      "msg": "查詢成功",
      "msg_detail": "",
      "data": {
        "id": "user-uuid-123",
        "email": "user@example.com",
        "created_at": "2025-04-25T10:35:00Z"
      }
    }
    ```
* **回應 (失敗 401 Unauthorized):**
    ```json
    {
      "status_code": "401",
      "msg": "查詢失敗",
      "msg_detail": "未經授權",
      "data": null
    }
    ```
* **回應 (失敗 404 Not Found):**
    ```json
    {
      "status_code": "404",
      "msg": "查詢失敗",
      "msg_detail": "找不到該用戶",
      "data": null
    }
    ```

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