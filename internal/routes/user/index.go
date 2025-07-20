package user_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/user"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	r.POST("/signIn", user_controller.SignIn)
	r.POST("/signUp", user_controller.SignUp)
	r.GET("/homepage/:id", user_controller.Homepage)
	r.POST("/markAttendance/:id", user_controller.MarkAttendance)
	r.GET("/profile/:id", user_controller.Profile)
	r.PATCH("/profile/:id", user_controller.UpdateProfile)
	r.POST("/sendOTP", user_controller.SendOTP)
}
