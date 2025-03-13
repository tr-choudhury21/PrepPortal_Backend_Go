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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	documentCollection *mongo.Collection
	documentOnce       sync.Once
)

func getDocumentCollection() *mongo.Collection {
	documentOnce.Do(func() {
		documentCollection = config.GetCollection("documents")
	})
	return documentCollection
}

func CreateDocument(c *gin.Context) {

	documentCollection := getDocumentCollection()
	if documentCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Parse form data (Limit: 10 MB)
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large"})
		return
	}

	// Extract file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	// Upload file to Cloudinary
	fileURL, err := utils.UploadFile(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upload file"})
		return
	}

	// Create document instance
	doc := models.Document{
		ID:         primitive.NewObjectID(),
		Subject:    c.PostForm("subject"),
		Semester:   c.PostForm("semester"),
		Year:       c.PostForm("year"),
		Branch:     c.PostForm("branch"),
		Content:    c.PostForm("content"),
		FileUrl:    fileURL,
		FileName:   header.Filename,
		UploadedBy: c.PostForm("uploadedBy"),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert into MongoDB
	_, err = documentCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Document uploaded successfully", "document": doc})
}

// Get All Documents
func GetAllDocuments(c *gin.Context) {

	documentCollection := getDocumentCollection()
	if documentCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	cursor, err := documentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch documents"})
		return
	}
	defer cursor.Close(context.TODO())

	var documents []models.Document
	for cursor.Next(context.TODO()) {
		var doc models.Document
		if err := cursor.Decode(&doc); err != nil { // Handle decode errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding document"})
			return
		}
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// Get Documents by Branch
func GetDocumentsByBranch(c *gin.Context) {

	documentCollection := getDocumentCollection()
	if documentCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	branch := c.Param("branch")

	cursor, err := documentCollection.Find(context.TODO(), bson.M{"branch": branch})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch documents"})
		return
	}
	defer cursor.Close(context.TODO())

	var documents []models.Document
	for cursor.Next(context.TODO()) {
		var doc models.Document
		if err := cursor.Decode(&doc); err != nil { // Handle decode errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding document"})
			return
		}
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// Update Document
func UpdateDocument(c *gin.Context) {

	documentCollection := getDocumentCollection()
	if documentCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id) // Validate ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var updatedData models.Document
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedData.UpdatedAt = time.Now()

	update := bson.M{"$set": updatedData}
	result, err := documentCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// Delete Document
func DeleteDocument(c *gin.Context) {

	documentCollection := getDocumentCollection()
	if documentCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id) // Validate ID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	result, err := documentCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
