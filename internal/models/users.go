package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	FirstName string    `gorm:"size:50;not null"`
	LastName  string    `gorm:"size:50;"`
	Email     string    `gorm:"size:100;not null;unique"`
	Phone     string    `gorm:"size:10;not null;unique"`
	UserName  string    `gorm:"size:50;not null;unique"`
	Password  string    `gorm:"not null"`
	Location  string    `gorm:"size:100;not null"`
	DOB       time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
