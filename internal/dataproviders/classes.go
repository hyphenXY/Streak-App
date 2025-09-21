package dataprovider

import (
	// "gorm.io/gorm"
	"errors"

	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/gorm"
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

	if errors.Is(err, gorm.ErrRecordNotFound) {
		attendance = models.Attendance{
			ClassID:      classID,
			MarkedById:   userID,
			MarkedByRole: "user",
			Status:       "present",
		}
		return DB.Create(&attendance).Error
	}
	if err != nil {
		return err
	}
	return errors.New("already marked")
}

func MarkAttendanceByAdmin(classID uint, userID uint) error {
	// check in attendances table if record exists
	var attendance models.Attendance
	err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND marked_by_id = ? AND marked_by_role = ? AND DATE(created_at) = CURRENT_DATE", classID, userID, "admin").
		First(&attendance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		attendance = models.Attendance{
			ClassID:      classID,
			MarkedById:   userID,
			MarkedByRole: "admin",
			Status:       "present",
		}
		return DB.Create(&attendance).Error
	}
	if err != nil {
		return err
	}
	return errors.New("already marked")
}

func IsUserAdmin(userID uint, classID uint) (bool, error) {
	var count int64
	err := DB.Model(&models.Classes{}).Where("created_by_admin_id = ? AND id = ?", userID, classID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetStudentsByClassID(classID uint) ([]models.User, error) {
	var students []models.User
	err := DB.Joins("JOIN enrollments ON enrollments.user_id = users.id").
		Where("enrollments.class_id = ?", classID).
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func GetClassIDByCode(classCode string) (uint, error) {
	var class models.Classes
	err := DB.Where("class_code = ?", classCode).First(&class).Error
	if err != nil {
		return 0, err
	}
	return class.ID, nil
}
