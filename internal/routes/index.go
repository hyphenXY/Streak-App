package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/hyphenXY/Streak-App/internal/routes/admin"
	"github.com/hyphenXY/Streak-App/internal/routes/root"
	"github.com/hyphenXY/Streak-App/internal/routes/user"
	"github.com/hyphenXY/Streak-App/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add more custom middlewares
	r.Use(middlewares.CORSMiddleware())
	// etc.

	// root-level routes (no prefix)
	rootGroup := r.Group("/root")
	root_routes.RegisterRootRoutes(rootGroup)

	// admin routes under /admin
	adminGroup := r.Group("/admin")
	admin_routes.RegisterAdminRoutes(adminGroup)

	// user routes under /user
	userGroup := r.Group("/user")
	user_routes.RegisterUserRoutes(userGroup)

	return r
}
