package dataprovider

import (
	"time"

	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/gorm"
)

func CreateAdmin(admin *models.Admin) error {
	return DB.Create(admin).Error
}

func GetClassesByAdmin(adminID uint, classes *[]models.Classes) error {
	return DB.Where("created_by_admin_id = ?", adminID).Find(classes).Error
}

func UpdateAdminRefreshToken(adminID uint, refreshToken string, refreshTokenExpiry time.Time) error {
	result := DB.Model(&models.Admin{}).
		Where("id = ?", adminID).
		Updates(map[string]interface{}{
			"RefreshToken":       refreshToken,
			"RefreshTokenExpiry": refreshTokenExpiry,
		})
	return result.Error
}

func AdminNameById(adminID uint, admin *models.Admin) error {
	return DB.Where("id = ?", adminID).First(admin).Error
}

func RevokeAdminRefreshToken(refreshToken string) error {
	result := DB.Model(&models.Admin{}).
		Where("refresh_token = ?", refreshToken).
		Updates(map[string]interface{}{
			"refresh_token":        nil,
			"refresh_token_expiry": nil,
		})
	return result.Error
}

func GetAdminProfile(adminID uint) (*models.Admin, error) {
	admin := &models.Admin{}
	DB.Where("id = ?", adminID).First(admin)
	return admin, nil
}

func UpdateAdminProfile(req map[string]interface{}, userId uint) error {

	updates := map[string]interface{}{}

	if req["first_name"] != "" {
		updates["FirstName"] = req["first_name"]
	}
	if req["last_name"] != "" {
		updates["LastName"] = req["last_name"]
	}
	if req["Email"] != "" {
		updates["Email"] = req["Email"]
	}

	// Always update UpdatedAt
	updates["UpdatedAt"] = time.Now()

	if len(updates) == 0 {
		return nil
	}

	result := DB.Model(&models.Admin{}).
		Where("id = ?", userId).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateAdminPassword(adminID uint, hashedPassword string) error {
	result := DB.Model(&models.Admin{}).
		Where("id = ?", adminID).
		Update("Password", hashedPassword)
	return result.Error
}

func UpdateAdminProfileImage(adminId uint, imageURL string) error {
	result := DB.Model(&models.Admin{}).
		Where("id = ?", adminId).
		Update("ProfileImageURL", imageURL)
	return result.Error
}