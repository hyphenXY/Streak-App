package dataprovider

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        user, password, host, port, dbname)

    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        return fmt.Errorf("error creating DB handle: %w", err)
    }

    if err = DB.Ping(); err != nil {
        // Close the handle since Ping failed
        DB.Close()
        DB = nil
        return fmt.Errorf("error pinging DB: %w", err)
    }

    log.Println("✅ Connected to MySQL!")
    return nil
}
