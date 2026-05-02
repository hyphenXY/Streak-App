package models

import "time"

type Classes struct {
	ID                        uint   `gorm:"primaryKey;autoIncrement"`
	Name                      string `gorm:"size:50;"`
	Email                     string `gorm:"size:100;"`
	Phone                     string `gorm:"size:10;"`
	CreatedByAdminId          uint   `gorm:"size:50;"`
	ClassCode                 string `gorm:"size:10;uniqueIndex"`
	NumberOfWorkingDaysInWeek uint8  `gorm:"check:number_of_working_days_in_week_range,number_of_working_days_in_week >= 1 AND number_of_working_days_in_week <= 7"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
