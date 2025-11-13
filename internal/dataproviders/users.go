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

func StoreOTP(phone uint, otp string) error {
	otpRecord := models.OTPs{
		Phone:      phone,
		OTP:        otp,
		Expiry:     time.Now().Add(10 * time.Minute),
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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
	// DB.Delete(&otpRecord)

	return "Verified!", nil // OTP verified successfully
}

func IsPhoneVerified(phone uint) (bool, error) {
	var otp models.OTPs
	err := DB.Where("phone = ? AND is_verified = ?", phone, true).First(&otp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Phone not verified
		}
		return false, err // Other error
	}
	return true, nil
}

func MarkPhoneVerified(phone uint) error {
	result := DB.Model(&models.OTPs{}).
		Where("phone = ?", phone).
		Update("is_verified", true)
	return result.Error
}

func RevokeUserRefreshToken(refreshToken string) error {
	result := DB.Model(&models.User{}).
		Where("refresh_token = ?", refreshToken).
		Updates(map[string]interface{}{
			"refresh_token":        nil,
			"refresh_token_expiry": nil,
		})
	return result.Error
}

func UpdateProfile(req map[string]interface{}, userId uint) error {
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

	result := DB.Model(&models.User{}).
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

func UpdateUserPassword(userId uint, hashedPassword string) error {
	result := DB.Model(&models.User{}).
		Where("id = ?", userId).
		Update("Password", hashedPassword)
	return result.Error
}

func IsOTPRecentlyVerified(otp uint, phone uint) (bool, error) {
	var otpRecord models.OTPs
	err := DB.Where("phone = ? AND otp = ?", phone, otp).Last(&otpRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // OTP not found
		}
		return false, err // Other error
	}

	// Check if OTP was created within 10 minutes
	if time.Since(otpRecord.CreatedAt) > 10*time.Minute {
		return false, nil // OTP expired
	}

	return true, nil // OTP verified successfully
}

func UpdateUserProfileImage(userId uint, imageURL string) error {
	result := DB.Model(&models.User{}).
		Where("id = ?", userId).
		Update("ProfileImageURL", imageURL)
	return result.Error
}