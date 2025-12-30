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
	protectedUserClasses := r.Group("")
	protectedUserClasses.Use(middlewares.AuthUserMiddleware())
	{
		protectedUserClasses.POST("/markAttendance/:classID", middlewares.IsUserClass(), user_controller.MarkAttendance)
		protectedUserClasses.GET("/classDetails/:classID", middlewares.IsUserClass(), user_controller.ClassDetails)
		protectedUserClasses.GET("/calendar/:classID", middlewares.IsUserClass(), user_controller.Calendar)
		protectedUserClasses.GET("/streak/:classID", middlewares.IsUserClass(), user_controller.Streak)
		protectedUserClasses.GET("/quickSummary/:classID", middlewares.IsUserClass(), user_controller.QuickSummary)
		protectedUserClasses.GET("/report/:classID", middlewares.IsUserClass(), user_controller.Report)
		protectedUserClasses.POST("/enroll/:classCode", middlewares.IsAllowedToEnroll(), user_controller.Enroll)
		protectedUserClasses.GET("/classList", user_controller.ClassList)
		protectedUserClasses.POST("/logOutUser", user_controller.LogOutUser)
		protectedUserClasses.PATCH("/profile/:id", user_controller.UpdateProfile)
		protectedUserClasses.GET("/profile", user_controller.Profile)
		protectedUserClasses.POST("/resetPassword", user_controller.ResetPassword)
	}
}
