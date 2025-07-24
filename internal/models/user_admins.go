package models

import "time"

type User_Admins struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	AdminID   uint      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
