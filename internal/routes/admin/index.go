package admin_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/admin"
	"github.com/hyphenXY/Streak-App/internal/middleware"
)

func RegisterAdminRoutes(r *gin.RouterGroup) {
	r.POST("/signIn", admin_controller.SignIn)
	r.POST("/register", admin_controller.Register)
	r.POST("/sendOTP", admin_controller.SendOTP)
	
	protected := r.Group("")
	protected.Use(middlewares.AuthAdminMiddleware())
	{
		protected.GET("/homepage/:id", admin_controller.Homepage)
		protected.GET("/personalHomepage/:id", admin_controller.PersonalHomepage)
		protected.POST("/markAttendance/:id", admin_controller.MarkAttendance)
		protected.GET("/profile/:id", admin_controller.Profile)
		protected.PATCH("/profile/:id", admin_controller.UpdateProfile)
		protected.POST("/createClass", admin_controller.CreateClass)
	}
}
