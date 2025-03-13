package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/controllers"
	"github.com/tr-choudhury21/prepportal_backend/middleware"
)

func BlogRoutes(router *gin.Engine) {
	blogGroup := router.Group("/blogs")
	{
		blogGroup.GET("/", controllers.GetAllBlogs)
		blogGroup.GET("/:id", controllers.GetBlog)
		blogGroup.POST("/", middleware.AuthMiddleware(), controllers.CreateBlog)
		blogGroup.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateBlog)
		blogGroup.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteBlog)
	}
}
