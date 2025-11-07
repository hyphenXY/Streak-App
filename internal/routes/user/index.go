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
	protectedUserClasses.Use(middlewares.AuthUserMiddleware(), middlewares.IsUserClass())
	{
		protectedUserClasses.POST("/markAttendance/:classID", user_controller.MarkAttendance)
		protectedUserClasses.GET("/classDetails/:classID", user_controller.ClassDetails)
		protectedUserClasses.GET("/calendar/:classID", user_controller.Calendar)
		protectedUserClasses.GET("/streak/:classID", user_controller.Streak)
		protectedUserClasses.GET("/quickSummary/:classID", user_controller.QuickSummary)
	}

	protectedUser := r.Group("")
	protectedUser.Use(middlewares.AuthUserMiddleware())
	{
		protectedUserClasses.POST("/enroll/:classCode", user_controller.Enroll)
		protectedUserClasses.GET("/classList", user_controller.ClassList)
		protectedUser.POST("/logOutUser", user_controller.LogOutUser)
		protectedUser.PATCH("/profile/:id", user_controller.UpdateProfile)
		protectedUser.GET("/profile", user_controller.Profile)

	}
}
