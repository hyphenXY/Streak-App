package models

import "time"

type User_Classes struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:""`
	ClassID   uint      `gorm:""`
	CreatedAt time.Time
	UpdatedAt time.Time
}
