package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tr-choudhury21/prepportal_backend/config"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := config.ConnectDB()
	defer client.Disconnect(context.Background())

	fmt.Println("MongoDB client connected and ready to use!")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to NITA PrepPortal API",
		})
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on Port:", port)

	router.Run(":" + port)
}
