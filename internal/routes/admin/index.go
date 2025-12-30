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
		protected.POST("/resetPassword", admin_controller.ResetPassword)
		protected.GET("/quickSummary/:classId", middlewares.IsAdminClass(), admin_controller.QuickSummary)
		protected.GET("/todaySummary/:classId", middlewares.IsAdminClass(), admin_controller.TodaySummary)
		protected.GET("/calendar/:classId", middlewares.IsAdminClass(), admin_controller.Calendar)
		protected.POST("/markAttendance/:classId", middlewares.IsAdminClass(), admin_controller.MarkAttendance)
		protected.GET("/studentsList/:classId", middlewares.IsAdminClass(), admin_controller.StudentsList)
		protected.GET("/streak/:classId", middlewares.IsAdminClass(), admin_controller.Streak)
		protected.GET("/personalSummary/:classId", middlewares.IsAdminClass(), admin_controller.PersonalSummary)
		protected.GET("/report/:classId", middlewares.IsAdminClass(), admin_controller.Report)
		protected.GET("/personalReport/:classId", middlewares.IsAdminClass(), admin_controller.PersonalReport)
		protected.POST("/kickStudent/:classId", middlewares.IsAdminClass(), admin_controller.KickStudent)
		protected.POST("/banStudent/:classId", middlewares.IsAdminClass(), admin_controller.BanStudent)
	}
}
