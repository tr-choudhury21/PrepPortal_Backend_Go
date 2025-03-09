package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tr-choudhury21/prepportal_backend/config"
	"github.com/tr-choudhury21/prepportal_backend/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//connect db
	client := config.ConnectDB()
	defer client.Disconnect(context.Background())
	fmt.Println("MongoDB client connected and ready to use!")

	router := gin.Default()

	//routes
	routes.AuthRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on Port:", port)
	router.Run(":" + port)
}
