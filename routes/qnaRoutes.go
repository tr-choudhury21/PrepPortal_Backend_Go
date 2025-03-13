package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/controllers"
	"github.com/tr-choudhury21/prepportal_backend/middleware"
)

func QnaRoutes(router *gin.Engine) {
	qnaGroup := router.Group("/qna")
	{
		qnaGroup.POST("/ask", middleware.AuthMiddleware(), controllers.AskQuestion)
		qnaGroup.POST("/answer/:id", middleware.AuthMiddleware(), controllers.AnswerQuestion)
		qnaGroup.GET("/all", controllers.GetPaginatedQnA)
		//qnaGroup.GET("/paginated", controllers.GetPaginatedQnA)
		qnaGroup.POST("/vote/:id", middleware.AuthMiddleware(), controllers.VoteQuestion)
		qnaGroup.POST("/answer/vote/:id", middleware.AuthMiddleware(), controllers.VoteAnswer)
	}
}
