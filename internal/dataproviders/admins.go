package dataprovider

import (
	"time"

	"github.com/hyphenXY/Streak-App/internal/models"
)

func CreateAdmin(admin *models.Admin) error {
	return DB.Create(admin).Error
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
