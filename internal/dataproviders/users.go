package dataprovider

import (
	// "fmt"
	"time"

	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/gorm"
)

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

func EnrollUser(userID uint, classID uint) error {
	enrollment := models.User_Classes{
		UserID:  userID,
		ClassID: classID,
	}
	result := DB.Create(&enrollment)
	return result.Error
}

func StoreOTP(phone uint, otp string) error {
	otpRecord := models.OTPs{
		Phone: phone,
		OTP:   otp,
	}
	result := DB.Create(&otpRecord)
	return result.Error
}

func VerifyOTP(phone uint, otp string) (string, error) {
	var otpRecord models.OTPs
	err := DB.Where("phone = ? AND otp = ?", phone, otp).Last(&otpRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "Wrong!", nil // OTP not found
		}
		return "Failed!", err // Other error
	}

	// Check if OTP was created within 10 minutes
	if time.Since(otpRecord.CreatedAt) > 10*time.Minute {
		return "Expired!", nil // OTP expired
	}

	// Optionally, you can delete the OTP after verification
	DB.Delete(&otpRecord)

	return "Verified!", nil // OTP verified successfully
}

func IfAlreadyEnrolled(userID uint, classID uint, enrollment *models.User_Classes) error {
	return DB.Where("user_id = ? AND class_id = ?", userID, classID).First(enrollment).Error
}