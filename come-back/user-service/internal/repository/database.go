package repository

import (
	"fmt"
	"log"
	"os"
	"user-service/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dB *gorm.DB

func InitDB(dsn string) error {
	var err error
	dB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = dB.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	addAdmin()

	sqlDB, _ := dB.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	return nil
}

func addAdmin() {
	var adminCount int64
	dB.Model(&model.User{}).Where("role = ?", model.RoleAdmin).Count(&adminCount)
	if adminCount != 0 {
		return
	}

	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash password, admin create abort")
		return
	}

	admin := model.User{
		Username: adminName,
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     model.RoleAdmin,
	}
	dB.Create(&admin)
	fmt.Println("Initial admin created:", adminName, adminEmail)
}
