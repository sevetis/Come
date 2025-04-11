package repository

import (
	"come-back/user-service/internal/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dB *gorm.DB

func InitDB(dsn string) {
	var err error
	dB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dB.AutoMigrate(&model.User{})

	sqlDB, _ := dB.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
}
