package models

import "time"

type OTPs struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Role      string `gorm:"size:50;"`
	Phone     uint   `gorm:""`
	OTP       string `gorm:"size:6;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
