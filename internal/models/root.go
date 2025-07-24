package models

import "time"

type Root struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	FirstName string    `gorm:"size:50;"`
	LastName  string    `gorm:"size:50;"`
	Email     string    `gorm:"size:100;"`
	Phone     string    `gorm:"size:10;"`
	UserName  string    `gorm:"size:50;"`
	Password  string    `gorm:""`
	DOB       time.Time `gorm:""`
	CreatedAt time.Time
	UpdatedAt time.Time
}
