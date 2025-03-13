package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Client
	once   sync.Once         // Ensures ConnectDB() is only called once
	dbName = "prepportal_Go" // Define database name globally
)

// ConnectDB initializes the MongoDB connection (singleton)
func ConnectDB() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Warning: .env file not found, using environment variables")
		}

		mongoURI := os.Getenv("MONGO_URI")
		if mongoURI == "" {
			log.Fatal("❌ MongoDB URI not found in environment variables")
		}

		fmt.Println("🔗 Connecting to MongoDB at:", mongoURI)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
		}

		// Ping the MongoDB server to verify the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("❌ Failed to ping MongoDB: %v", err)
		}

		DB = client
		fmt.Println("✅ Connected to MongoDB successfully!")
	})
}

// GetCollection returns a MongoDB collection safely
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("❌ MongoDB connection is not initialized. Call ConnectDB() first.")
	}
	return DB.Database(dbName).Collection(collectionName)
}
