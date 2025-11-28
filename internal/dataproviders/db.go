package dataprovider

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Check for missing environment variables
	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		return fmt.Errorf("missing required database environment variables: user=%s, password=%s, host=%s, port=%s, dbname=%s", user, password, host, port, dbname)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to DB: %w", err)
	}

	log.Println("✅ Connected to MySQL (via GORM)!")

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// 👉 Load your models here and auto-migrate:
	err = DB.AutoMigrate(
		&models.User{}, // import from your models package
		&models.Admin{},
		&models.Root{},
		&models.Attendance{},
		&models.User_Classes{},
		&models.Classes{},
		&models.OTPs{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("✅ Tables migrated successfully!")
	return nil
}
