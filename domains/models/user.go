package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key; column:ID" json:"id"`
	Name      string    `gorm:"column:Name" json:"name"`
	CreatedAt time.Time `gorm:"column:UpdatedAt" json:"created_at"`
}

type Users struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}
