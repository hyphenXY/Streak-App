package dataprovider

import (
	// "gorm.io/gorm"
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