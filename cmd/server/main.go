package main

import (
	"fmt"
	"log"

	"github.com/hyphenXY/Streak-App/internal/dataproviders"
	"github.com/hyphenXY/Streak-App/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if you use one
	_ = godotenv.Load()

	if err := dataprovider.InitDB(); err != nil {
		log.Printf("⚠️  Could not connect to database: %v", err)
		// DB is nil; your app should handle it gracefully in handlers
	} else {
		defer dataprovider.DB.Close()
	}

	// Start router
	r := routes.SetupRouter()
	fmt.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
