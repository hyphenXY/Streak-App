package root_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/root"
	"github.com/hyphenXY/Streak-App/internal/middleware"
)

func RegisterRootRoutes(r *gin.RouterGroup) {
	r.POST("/signIn", root_controller.SignIn)
	r.POST("/register", root_controller.Register)
	r.GET("/health-check", root_controller.HealthCheck)

	protected := r.Group("")
	protected.Use(middlewares.AuthRootMiddleware())
	{
	protected.GET("/homepage/:id", root_controller.Homepage)
	protected.GET("/profile/:id", root_controller.Profile)
	protected.PATCH("/profile/:id", root_controller.UpdateProfile)
	protected.DELETE("/admin/:id", root_controller.DeleteAdmin)
	}
}
