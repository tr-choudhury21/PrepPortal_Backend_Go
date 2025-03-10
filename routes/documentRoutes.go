package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/controllers"
	"github.com/tr-choudhury21/prepportal_backend/middleware"
)

func DocumentRoutes(router *gin.Engine) {
	docs := router.Group("/documents")
	{
		docs.POST("/", middleware.AuthMiddleware(), controllers.CreateDocument) // Protected
		docs.GET("/", controllers.GetAllDocuments)
		docs.GET("/:branch", controllers.GetDocumentsByBranch)
		docs.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateDocument)    // Protected
		docs.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteDocument) // Protected
	}
}
