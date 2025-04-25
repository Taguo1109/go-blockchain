package models

import "time"

/**
 * @File: user.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午10:34
 * @Software: GoLand
 * @Version:  1.0
 */

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // 不希望在 JSON 響應中返回密碼
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
