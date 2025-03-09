package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/controllers"
)

func AuthRoutes(router *gin.Engine) {

	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}
}
