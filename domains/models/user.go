package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key; column:ID" json:"id"`
	Name      string    `gorm:"column:Name" json:"name"`
	CreatedAt time.Time `gorm:"column:UpdatedAt" json:"created_at"`
}
