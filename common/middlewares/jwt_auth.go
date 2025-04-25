package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-blockchain/common/models"
	"net/http"
	"os"
	"strings"
	"time"
)

/**
 * @File: jwt_auth.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午11:17
 * @Software: GoLand
 * @Version:  1.0
 */

// JwtKey Secret key for JWT signing
var JwtKey = []byte(os.Getenv("JWT_SECRET"))

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token
func GenerateToken(userID string, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// AuthMiddleware is the JWT authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.JsonResult{
				StatusCode: "401",
				Msg:        "驗證失敗",
				MsgDetail:  "未提供授權憑證",
				Data:       nil,
			})
			return
		}

		parts := strings.Split(tokenString, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.JsonResult{
				StatusCode: "401",
				Msg:        "驗證失敗",
				MsgDetail:  "無效的授權憑證格式",
				Data:       nil,
			})
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return JwtKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.JsonResult{
				StatusCode: "401",
				Msg:        "驗證失敗",
				MsgDetail:  fmt.Sprintf("無效的 token: %v", err),
				Data:       nil,
			})
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// 將用戶 ID 存儲到 Gin 上下文中，以便後續的處理函數可以使用
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.JsonResult{
				StatusCode: "401",
				Msg:        "驗證失敗",
				MsgDetail:  "無效的 token claims",
				Data:       nil,
			})
			return
		}
	}
}
