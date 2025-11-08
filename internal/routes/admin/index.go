package admin_routes

import (
	"github.com/gin-gonic/gin"
	admin_controller "github.com/hyphenXY/Streak-App/internal/controllers/admin"
	middlewares "github.com/hyphenXY/Streak-App/internal/middleware"
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
		protected.GET("/profile", admin_controller.Profile)
		protected.PATCH("/profile", admin_controller.UpdateProfile)
		protected.POST("/createClass", admin_controller.CreateClass)
		protected.POST("/logOutAdmin", admin_controller.LogOutAdmin)
	}
	
	protectedAdminClasses := r.Group("")
	protectedAdminClasses.Use(middlewares.AuthAdminMiddleware(), middlewares.IsAdminClass())
	{
		protected.GET("/quickSummary/:classId", admin_controller.QuickSummary)
		protected.POST("/markAttendance/:classId", admin_controller.MarkAttendance)
		protected.GET("/studentsList/:classId", admin_controller.StudentsList)
		protected.GET("/streak/:classId", admin_controller.Streak)
		protected.GET("/personalSummary/:classId", admin_controller.PersonalSummary)

	}
}
