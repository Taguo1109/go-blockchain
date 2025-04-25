package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blockchain/common/models"
	"net/http"
)

/**
 * @File: global_error_handler.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午11:07
 * @Software: GoLand
 * @Version:  1.0
 */

// GlobalErrorHandler handles unexpected panics and returns JSON instead of crashing the server.
// 可作為全域錯誤攔截器使用。
func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 捕捉 panic，回傳統一格式
				c.JSON(http.StatusInternalServerError, models.JsonResult{
					StatusCode: "500",
					Msg:        "全域攔截：伺服器錯誤",
					MsgDetail:  fmt.Sprintf("%v", err),
					Data:       nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
