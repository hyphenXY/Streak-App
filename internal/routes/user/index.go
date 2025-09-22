package user_routes

import (
	"github.com/gin-gonic/gin"
	user_controller "github.com/hyphenXY/Streak-App/internal/controllers/user"
	middlewares "github.com/hyphenXY/Streak-App/internal/middleware"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	// Public routes
	r.POST("/signIn", user_controller.SignIn)
	r.POST("/signUp", user_controller.SignUp)
	r.POST("/sendOTP", user_controller.SendOTP)
	r.POST("/verifyOTP", user_controller.VerifyOTP)
	r.POST("/refreshToken", user_controller.RefreshTokenUser)

	// Protected routes
	protected := r.Group("")
	protected.Use(middlewares.AuthUserMiddleware())
	{
		protected.GET("/classList", user_controller.ClassList)
		protected.POST("/markAttendance/:id", user_controller.MarkAttendance)
		protected.GET("/profile/:id", user_controller.Profile)
		protected.PATCH("/profile/:id", user_controller.UpdateProfile)
		protected.POST("/enroll/:classCode", user_controller.Enroll)
		protected.POST("/classDetails", user_controller.ClassDetails)
		protected.GET("/myClasses/:id", user_controller.QuickSummary)
		protected.GET("/calendar", user_controller.Calendar)
		protected.POST("/logOutUser", user_controller.LogOutUser)
	}
}
