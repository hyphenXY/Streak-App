package dataprovider

import (
	// "gorm.io/gorm"
	"errors"
	"github.com/hyphenXY/Streak-App/internal/models"
)

func IfClassExists(classID uint) (bool, error) {
	var count int64
	err := DB.Model(&models.Classes{}).Where("id = ?", classID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateClass(class *models.Classes) error {
	return DB.Create(class).Error
}

func MarkAttendanceByUser(classID uint, userID uint) error {
	// check in attendances table if record exists
	var attendance models.Attendance
	err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND marked_by_id = ? AND marked_by_role = ? AND DATE(created_at) = CURRENT_DATE", classID, userID, "user").
		First(&attendance).Error
	if err != nil {
		return err
	}
	if attendance.Status == "present" {
		return errors.New("already marked")

	}

	return DB.Model(&models.Attendance{}).Where("id = ?", attendance.ID).Update("present", true).Error
}
