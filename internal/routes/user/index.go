package user_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/user"
	"github.com/hyphenXY/Streak-App/internal/middleware"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
    // Public routes
    r.POST("/signIn", user_controller.SignIn)
    r.POST("/signUp", user_controller.SignUp)
    r.POST("/sendOTP", user_controller.SendOTP)

    // Protected routes
    protected := r.Group("")
    protected.Use(middlewares.AuthUserMiddleware())
    {
        protected.GET("/homepage/:id", user_controller.Homepage)
        protected.POST("/markAttendance/:id", user_controller.MarkAttendance)
        protected.GET("/profile/:id", user_controller.Profile)
        protected.PATCH("/profile/:id", user_controller.UpdateProfile)
    }
}
