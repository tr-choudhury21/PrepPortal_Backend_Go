package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var CLD *cloudinary.Cloudinary

func InitCloudinary() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Fatal("Missing Cloudinary environment variables")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}

	CLD = cld
	log.Println("✅ Cloudinary initialized successfully!")
}
