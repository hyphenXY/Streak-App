package main

import (
	"fmt"

	"github.com/hyphenXY/Streak-App/internal/routes"
)

func main() {
	r := routes.SetupRouter()

	fmt.Println("🚀 Server running on http://localhost:8080")
	// Gin automatically listens and serves
	r.Run(":8080")
}
