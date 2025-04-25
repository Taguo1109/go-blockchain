package handlers

import (
	"github.com/gin-gonic/gin"
	"go-blockchain/common/models"
	"go-blockchain/internal/user/database"
	"net/http"
)

/**
 * @File: user_handler.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午11:33
 * @Software: GoLand
 * @Version:  1.0
 */

// GetUser handles fetching user information by ID
func GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "查詢失敗",
			MsgDetail:  "用戶 ID 不能為空",
			Data:       nil,
		})
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "查詢失敗",
			MsgDetail:  "查詢用戶時發生錯誤",
			Data:       err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, models.JsonResult{
			StatusCode: "404",
			Msg:        "查詢失敗",
			MsgDetail:  "找不到該用戶",
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        "成功",
		MsgDetail:  "查詢成功",
		Data: map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
