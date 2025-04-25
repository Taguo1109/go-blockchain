package database

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-blockchain/common/database"
	"go-blockchain/internal/user/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

/**
 * @File: user_db.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午10:51
 * @Software: GoLand
 * @Version:  1.0
 */

func CreateUser(user *models.User) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = string(hashedPassword)

	result := db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to insert user: %w", result.Error)
	}

	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to query user by email: %w", result.Error)
	}

	return &user, nil
}

func GetUserByID(id string) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var user models.User
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to query user by ID: %w", result.Error)
	}

	// 為了安全，查詢 ID 時不返回密碼
	user.Password = ""

	return &user, nil
}
