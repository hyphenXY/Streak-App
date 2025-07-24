package dataprovider

import (
	// "fmt"
	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/gorm"
)

var UserDB *gorm.DB

func CreateUser(user *models.User) error {
	result := UserDB.Create(user)
	return result.Error
}
