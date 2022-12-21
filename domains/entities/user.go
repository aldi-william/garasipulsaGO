package entities

import "time"

type User struct {
	ID        uint      `gorm:"primary_key; column:ID" json:"id"`
	Name      string    `gorm:"column:Name" json:"name"`
	CreatedAt time.Time `gorm:"column:CreatedAt" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt" json:"updated_at"`
}
