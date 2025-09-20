package models

import "time"

type User struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	FirstName          string     `gorm:"size:50;"`
	LastName           string     `gorm:"size:50;"`
	Email              string     `gorm:"size:100;"`
	Phone              string     `gorm:"size:10;"`
	UserName           string     `gorm:"size:50;"`
	Password           string     `gorm:""`
	Location           string     `gorm:"size:100;"`
	DOB                time.Time  `gorm:""`
	RefreshToken       *string    `gorm:"size:255"`
	RefreshTokenExpiry *time.Time `gorm:""`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
