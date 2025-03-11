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
	err = dB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		return err
	}

	err = dB.Exec(`
		ALTER TABLE comments
		ADD CONSTRAINT fk_comments_post
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		ADD CONSTRAINT fk_comments_author
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE RESTRICT
	`).Error
	if err != nil {
		fmt.Println("error:", err)
	}

	err = dB.Exec(`
		ALTER TABLE posts
		ADD CONSTRAINT fk_posts_author
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	`).Error
	if err != nil {
		fmt.Println("error:", err)
	}

	return nil
}
