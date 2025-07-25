package models

import "time"

type Classes struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	Name             string `gorm:"size:50;"`
	Email            string `gorm:"size:100;"`
	Phone            string `gorm:"size:10;"`
	CreatedByAdminId uint   `gorm:"size:50;"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
