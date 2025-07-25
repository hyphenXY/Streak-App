package dataprovider

import (
    "fmt"
    "log"
    "os"

    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "github.com/hyphenXY/Streak-App/internal/models"
    
)

var DB *gorm.DB

func InitDB() error {
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        user, password, host, port, dbname)

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("error connecting to DB: %w", err)
    }

    log.Println("✅ Connected to MySQL (via GORM)!")

    // 👉 Load your models here and auto-migrate:
    err = DB.AutoMigrate(
        &models.User{},       // import from your models package
        &models.Admin{},
        &models.Root{},
        &models.Attendance{},
        &models.User_Classes{},
        &models.Classes{},
    )
    if err != nil {
        return fmt.Errorf("auto migration failed: %w", err)
    }

    log.Println("✅ Tables migrated successfully!")
    return nil
}
