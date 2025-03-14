package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/controllers"
	"github.com/tr-choudhury21/prepportal_backend/middleware"
)

func AuthRoutes(router *gin.Engine) {

	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Signup)
		auth.POST("/login", controllers.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetUserProfile)
		auth.PUT("/profile", middleware.AuthMiddleware(), controllers.UpdateUserProfile)
		auth.GET("/leaderboard", controllers.GetLeaderboard)
	}
}
