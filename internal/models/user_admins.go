package models

import "time"

type User_Admins struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:""`
	AdminID   uint      `gorm:""`
	CreatedAt time.Time
	UpdatedAt time.Time
}
