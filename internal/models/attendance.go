package models

import "time"

type Attendance struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	MarkedById   uint   `gorm:""`
	MarkedByRole string `gorm:"type:ENUM('admin', 'user');"`
	MarkedForId  uint   `gorm:""`
	Status       string `gorm:"type:ENUM('present', 'absent', 'unmarked');default:'unmarked';"`
	Reason       string `gorm:"size:255;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
