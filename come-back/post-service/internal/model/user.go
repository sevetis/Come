package model

import "time"

type UserRole int

const (
	RoleUser UserRole = iota
	RoleAdmin
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(255);not null" json:"username"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Avatar    string    `gorm:"type:varchar(255)" json:"avatar,omitempty"`
	Role      UserRole  `gorm:"not null;default:0" json:"role"`
	Banned    bool      `gorm:"default:false" json:"banned"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
