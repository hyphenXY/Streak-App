package models

import "time"

type Attendance struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	MarkedById uint      `gorm:"not null"`
	MarkedByRole string    `gorm:"type:ENUM('admin', 'user');not null"`
	MarkedForId uint      `gorm:"not null"`
	Status    string    `gorm:"type:ENUM('present', 'absent', 'unmarked');default:'unmarked';not null"`
	Reason	string    `gorm:"size:255;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
