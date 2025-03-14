package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/config"
	"github.com/tr-choudhury21/prepportal_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	blogCollection *mongo.Collection
	blogOnce       sync.Once
)

func GetBlogCollection() *mongo.Collection {
	blogOnce.Do(func() {
		blogCollection = config.GetCollection("blogs")
	})

	return blogCollection
}

// CreateBlog allows registered users to add blogs
func CreateBlog(c *gin.Context) {

	blogCollection := GetBlogCollection()
	userCollection := getUserCollection()

	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//fetch user details
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": userEmail.(string)}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var blog models.Blog
	if err := c.BindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Set author details
	blog.ID = primitive.NewObjectID()
	blog.Author = userEmail.(string)
	blog.AuthorID = user.ID
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	// Insert into DB
	_, err = blogCollection.InsertOne(context.TODO(), blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving blog"})
		return
	}

	// Update user's contributions list
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$push": bson.M{"contributions": blog.ID}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user contributions"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog": blog})
}

// GetAllBlogs retrieves all blogs
func GetAllBlogs(c *gin.Context) {

	blogCollection := GetBlogCollection()
	cursor, err := blogCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching blogs"})
		return
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding blogs"})
		return
	}

	c.JSON(http.StatusOK, blogs)
}

// GetBlog retrieves a single blog by ID
func GetBlog(c *gin.Context) {

	blogCollection := GetBlogCollection()
	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	err = blogCollection.FindOne(context.TODO(), bson.M{"_id": blogID}).Decode(&blog)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// UpdateBlog allows authors to update their blogs
func UpdateBlog(c *gin.Context) {

	blogCollection := GetBlogCollection()
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blogUpdate models.Blog
	if err := c.BindJSON(&blogUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	filter := bson.M{"_id": blogID, "author": userEmail}
	update := bson.M{"$set": bson.M{"title": blogUpdate.Title, "content": blogUpdate.Content, "updatedAt": time.Now()}}

	result, err := blogCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
}

// DeleteBlog allows authors to delete their blogs
func DeleteBlog(c *gin.Context) {

	blogCollection := GetBlogCollection()
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	blogID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	filter := bson.M{"_id": blogID, "author": userEmail}
	result, err := blogCollection.DeleteOne(context.TODO(), filter)
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}
