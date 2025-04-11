package repository

import (
	"come-back/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dB *gorm.DB

func InitMySQL(dsn string) error {
	var err error
	dB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("MySQL connection failed: %v", err)
	}
	fmt.Println("MySQL connected")
	err = dB.AutoMigrate(&model.Post{}, &model.Comment{})
	if err != nil {
		return err
	}

	return nil
}
