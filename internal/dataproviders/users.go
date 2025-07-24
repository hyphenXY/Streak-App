package dataprovider

import (
	// "fmt"
	"time"

	"github.com/hyphenXY/Streak-App/internal/models"
)

// var DB *gorm.DB

func CreateUser(user *models.User) error {
	result := DB.Create(user)
	return result.Error
}

func UpdateUserRefreshToken(userID uint, refreshToken string, refreshTokenExpiry time.Time) error {
	result := DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"RefreshToken":       refreshToken,
			"RefreshTokenExpiry": refreshTokenExpiry,
		})
	return result.Error
}
