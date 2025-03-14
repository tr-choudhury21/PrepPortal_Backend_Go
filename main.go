package main

import (
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
	config.ConnectDB()

	if config.DB == nil {
		log.Fatal("❌ MongoDB connection is still nil after calling ConnectDB!")
	} else {
		fmt.Println("✅ MongoDB connection is ready.")
	}

	//connect cloudinary
	config.InitCloudinary()

	router := gin.Default()

	//routes
	routes.AuthRoutes(router)
	routes.DocumentRoutes(router)
	routes.QnaRoutes(router)
	routes.BlogRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Server running on Port:", port)
	router.Run(":" + port)
}
