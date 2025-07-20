package root_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/controllers/root"
)

func RegisterRootRoutes(r *gin.RouterGroup) {
	r.POST("/signIn", root_controller.SignIn)
	r.POST("/register", root_controller.Register)
	r.GET("/homepage/:id", root_controller.Homepage)
	r.GET("/profile/:id", root_controller.Profile)
	r.PATCH("/profile/:id", root_controller.UpdateProfile)
	r.DELETE("/admin/:id", root_controller.DeleteAdmin)
}
