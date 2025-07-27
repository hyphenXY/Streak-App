package models

import "time"

type OTPs struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Phone     uint   `gorm:""`
	OTP       string `gorm:"size:6;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
