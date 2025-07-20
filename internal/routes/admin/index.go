package admin_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/admin"
)

func RegisterAdminRoutes(r *gin.RouterGroup) {
	r.POST("/signIn", admin_controller.SignIn)
	r.POST("/register", admin_controller.Register)
	r.GET("/homepage/:id", admin_controller.Homepage)
	r.GET("/personalHomepage/:id", admin_controller.PersonalHomepage)
	r.POST("/markAttendance/:id", admin_controller.MarkAttendance)
	r.GET("/profile/:id", admin_controller.Profile)
	r.PATCH("/profile/:id", admin_controller.UpdateProfile)
	r.POST("/sendOTP", admin_controller.SendOTP)
}
