package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

/**
 * @File: init_db.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午10:52
 * @Software: GoLand
 * @Version:  1.0
 */

var db *gorm.DB

func InitDB(dbType string) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	var dsn string
	var err error

	switch dbType {
	case "mysql":
		// 從環境變數中讀取 MySQL 資料庫連線資訊
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")

		if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
			return nil, fmt.Errorf("MySQL database connection environment variables are not set")
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName,
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to open MySQL database with GORM: %w", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get generic database object: %w", err)
		}

		sqlDB.SetMaxIdleConns(5)  // 設定資料庫連接池的最大閒置連接數為 5。即使沒有活動的連接，也會保持最多 5 個閒置連接處於開啟狀態，以便快速響應後續請求。
		sqlDB.SetMaxOpenConns(10) // 設定資料庫連接池的最大開啟連接數為 10。
		// sqlDB.SetConnMaxLifetime(time.Hour) // GORM 預設會處理連接的生命週期

	// 可以根據需要添加其他資料庫類型的 case
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	log.Printf("%s database connected successfully with GORM", dbType)
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}
