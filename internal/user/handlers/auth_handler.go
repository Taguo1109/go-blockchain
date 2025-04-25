package handlers

import (
	"github.com/gin-gonic/gin"
	"go-blockchain/common/middlewares"
	"go-blockchain/common/models"
	"go-blockchain/internal/user/database"
	usermodel "go-blockchain/internal/user/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

/**
 * @File: auth_handler.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午11:21
 * @Software: GoLand
 * @Version:  1.0
 */

// Register handles user registration
func Register(c *gin.Context) {
	var userModel usermodel.User
	if err := c.ShouldBindJSON(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "註冊失敗",
			MsgDetail:  "請求體格式錯誤",
			Data:       err.Error(),
		})
		return
	}

	if userModel.Email == "" || userModel.Password == "" {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "註冊失敗",
			MsgDetail:  "電子郵件和密碼不能為空",
			Data:       nil,
		})
		return
	}

	existingUser, err := database.GetUserByEmail(userModel.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "註冊失敗",
			MsgDetail:  "查詢用戶時發生錯誤",
			Data:       err.Error(),
		})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, models.JsonResult{
			StatusCode: "409",
			Msg:        "註冊失敗",
			MsgDetail:  "該電子郵件已被註冊",
			Data:       nil,
		})
		return
	}

	if err := database.CreateUser(&userModel); err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "註冊失敗",
			MsgDetail:  "創建用戶時發生錯誤",
			Data:       err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.JsonResult{
		StatusCode: "201",
		Msg:        "成功",
		MsgDetail:  "註冊成功",
		Data: map[string]interface{}{
			"id":    userModel.ID,
			"email": userModel.Email,
		},
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.JsonResult{
			StatusCode: "400",
			Msg:        "登入失敗",
			MsgDetail:  "請求體格式錯誤",
			Data:       err.Error(),
		})
		return
	}

	user, err := database.GetUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "登入失敗",
			MsgDetail:  "查詢用戶時發生錯誤",
			Data:       err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, models.JsonResult{
			StatusCode: "401",
			Msg:        "登入失敗",
			MsgDetail:  "無效的憑證",
			Data:       nil,
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.JsonResult{
			StatusCode: "401",
			Msg:        "登入失敗",
			MsgDetail:  "無效的憑證",
			Data:       nil,
		})
		return
	}

	token, err := middlewares.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JsonResult{
			StatusCode: "500",
			Msg:        "登入失敗",
			MsgDetail:  "生成 token 失敗",
			Data:       err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        "成功",
		MsgDetail:  "登入成功",
		Data: map[string]string{
			"id":    user.ID,
			"token": token,
		},
	})
}

// Logout handles user logout (client-side clear token, server-side no specific action for stateless JWT)
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, models.JsonResult{
		StatusCode: "200",
		Msg:        "成功",
		MsgDetail:  "登出成功，客戶端應清除 token",
		Data:       nil,
	})
}
