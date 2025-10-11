package admin_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/admin"
	"github.com/hyphenXY/Streak-App/internal/middleware"
)

func RegisterAdminRoutes(r *gin.RouterGroup) {
	r.POST("/signIn", admin_controller.SignIn)
	r.POST("/signUp", admin_controller.SignUp)
	r.POST("/sendOTP", admin_controller.SendOTP)
	r.POST("/verifyOTP", admin_controller.VerifyOTP)
	r.POST("/refreshToken", admin_controller.RefreshTokenUser)
	
	protected := r.Group("")
	protected.Use(middlewares.AuthAdminMiddleware())
	{
		protected.GET("/classList", admin_controller.ClassList)
		protected.GET("/daySummary", admin_controller.DaySummary)
		protected.POST("/markAttendance/:classId", admin_controller.MarkAttendance)
		protected.GET("/profile/:id", admin_controller.Profile)
		protected.PATCH("/profile/:id", admin_controller.UpdateProfile)
		protected.POST("/createClass", admin_controller.CreateClass)
		protected.GET("/studentsList", admin_controller.StudentsList)
		protected.POST("/logOutAdmin", admin_controller.LogOutAdmin)
		protected.GET("/streak", admin_controller.Streak)
	}
}
