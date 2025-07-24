package main

import (
    "fmt"
    "log"

    "github.com/hyphenXY/Streak-App/internal/dataproviders"
    "github.com/hyphenXY/Streak-App/internal/routes"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    _ = godotenv.Load()

    // Initialize DB (connect + migrate)
    if err := dataprovider.InitDB(); err != nil {
        log.Fatalf("❌ Could not initialize database: %v", err)
    }

    // Start the Gin router
    r := routes.SetupRouter()

    fmt.Println("🚀 Server running on http://localhost:8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("❌ Server failed: %v", err)
    }
}
