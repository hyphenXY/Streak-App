package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/hyphenXY/Streak-App/internal/middleware"
	"github.com/hyphenXY/Streak-App/internal/routes/admin"
	"github.com/hyphenXY/Streak-App/internal/routes/root"
	"github.com/hyphenXY/Streak-App/internal/routes/user"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add more custom middlewares
	r.Use(middlewares.CORSMiddleware())

	limiter := middlewares.NewClientLimiter(2, 5) // 2 req/sec per client, burst up to 5
	r.Use(limiter.LimitMiddleware())

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
