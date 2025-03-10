package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/config"
	"github.com/tr-choudhury21/prepportal_backend/models"
	"github.com/tr-choudhury21/prepportal_backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

var documentCollection = config.GetCollection("documents")

func CreateDocument(c *gin.Context) {
	// Parse form data
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large"})
		return
	}

	// Extract file from the request
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

	// Insert document into MongoDB
	_, err = documentCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save document"})
		return
	}

	// Send success response
	c.JSON(http.StatusCreated, gin.H{"message": "Document uploaded successfully", "document": doc})
}

//get all docs

func GetAllDocuments(c *gin.Context) {
	cursor, err := documentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch documents"})
		return
	}
	defer cursor.Close(context.TODO())

	var documents []models.Document
	for cursor.Next(context.TODO()) {
		var doc models.Document
		cursor.Decode(&doc)
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// Get Documents by Branch
func GetDocumentsByBranch(c *gin.Context) {
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
		cursor.Decode(&doc)
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// Update Document
func UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	var updatedData models.Document
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedData.UpdatedAt = time.Now()

	update := bson.M{"$set": updatedData}
	_, err := documentCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// Delete Document
func DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := documentCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
