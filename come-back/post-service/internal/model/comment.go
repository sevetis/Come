package model

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"postId"`
	Post      Post      `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
	AuthorID  uint      `gorm:"not null;index" json:"authorId"`
	Author    User      `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:RESTRICT"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
