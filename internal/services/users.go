package services

import (
	"errors"

	"github.com/hyphenXY/Streak-App/internal/constants"
	"github.com/hyphenXY/Streak-App/internal/dataproviders"
	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/gorm"
)

func IsUserEnrolledInClass(userID uint, classID uint) (bool, error) {
	var enrollment models.User_Classes
	err := dataprovider.DB.Where("user_id = ? AND class_id = ?", userID, classID).First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return enrollment.Status == constants.UserEnrollment.Enrolled, nil
}
