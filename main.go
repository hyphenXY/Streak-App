package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hyphenXY/Streak-App/internal/cache"
	dataprovider "github.com/hyphenXY/Streak-App/internal/dataproviders"
	"github.com/hyphenXY/Streak-App/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (only used locally)
	_ = godotenv.Load("./.env")

	// Initialize DB (connect + migrate)
	if err := dataprovider.InitDB(); err != nil {
		log.Fatalf("❌ Could not initialize database: %v", err)
	}

	cache.InitRedis()

	// Start the Gin router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Render's default fallback port
	}
	fmt.Printf("🚀 Server running on http://0.0.0.0:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
