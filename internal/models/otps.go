package models

import "time"

type OTPs struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Phone      uint   `gorm:""`
	OTP        string `gorm:"size:6;"`
	IsVerified bool   `gorm:"default:false"`
	Expiry     time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
