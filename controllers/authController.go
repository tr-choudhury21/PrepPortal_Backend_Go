package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/config"
	"github.com/tr-choudhury21/prepportal_backend/models"
	"github.com/tr-choudhury21/prepportal_backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var (
	userCollection *mongo.Collection
	once           sync.Once
)

func getUserCollection() *mongo.Collection {
	once.Do(func() {
		userCollection = config.GetCollection("users")
	})
	return userCollection
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func Signup(c *gin.Context) {

	userCollection := getUserCollection()
	var user models.User

	// Parse JSON body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request",
		})

		return
	}

	if user.FullName == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Check if user exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user.Password = hashedPassword
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	// Insert user into database
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "details": err.Error()})
		return
	}

	// Generate JWT Token (Now using utils package)
	token, err := utils.GenerateToken(user.Email, user.FullName, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "token": token})
}

// Login Controller
func Login(c *gin.Context) {

	userCollection := getUserCollection()

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Find user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare Passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT Token (Now using utils package)
	token, err := utils.GenerateToken(user.Email, user.FullName, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// GetUserProfile retrieves the authenticated user's profile
func GetUserProfile(c *gin.Context) {

	userCollection := getUserCollection()

	// Extract user email from context
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Find user in the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Return user profile
	// c.JSON(http.StatusOK, gin.H{
	// 	"id":        user.ID.Hex(),
	// 	"fullName":  user.FullName,
	// 	"email":     user.Email,
	// 	"createdAt": user.CreatedAt,
	// })

	c.JSON(http.StatusOK, gin.H{"user": user})
}

//update profile

func UpdateUserProfile(c *gin.Context) {
	userCollection := getUserCollection()

	// Get user email from JWT token
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch user from DB
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": userEmail.(string)}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse request body
	var updateRequest struct {
		FullName string `json:"fullName,omitempty"`
		Bio      string `json:"bio,omitempty"`
	}

	if err := c.BindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Prepare update data
	update := bson.M{}
	if updateRequest.FullName != "" {
		update["fullName"] = updateRequest.FullName
	}
	if updateRequest.Bio != "" {
		update["bio"] = updateRequest.Bio
	}
	update["updatedAt"] = time.Now()

	// Update user profile
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": userEmail.(string)}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

//GetLeaderBoard of Contributors

func GetLeaderboard(c *gin.Context) {
	opts := options.Find().SetSort(bson.M{"reputation": -1}).SetLimit(10)
	cursor, _ := userCollection.Find(context.TODO(), bson.M{}, opts)

	var users []models.User
	cursor.All(context.TODO(), &users)

	c.JSON(http.StatusOK, users)
}
