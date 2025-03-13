package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tr-choudhury21/prepportal_backend/config"
	"github.com/tr-choudhury21/prepportal_backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	qnaCollection *mongo.Collection
	qnaOnce       sync.Once
)

func getQnaCollection() *mongo.Collection {
	qnaOnce.Do(func() {
		qnaCollection = config.GetCollection("qna")
	})
	return qnaCollection
}

// Ask a Question
func AskQuestion(c *gin.Context) {
	var qna models.Qna

	qnaCollection := getQnaCollection()

	// Extract user information from JWT token

	fullName, nameExists := c.Get("fullName")

	if !nameExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.BindJSON(&qna); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	qna.ID = primitive.NewObjectID()
	qna.CreatedAt = time.Now()
	qna.Upvotes = 0
	qna.Downvotes = 0

	fmt.Println("Retrieved Full Name:", fullName)

	if nameStr, ok := fullName.(string); ok {
		qna.PostedBy = nameStr
	} else {
		qna.PostedBy = "Anonymous"
	}

	_, err := qnaCollection.InsertOne(context.TODO(), qna)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question posted successfully", "question": qna})
}

// Answer a Question
func AnswerQuestion(c *gin.Context) {
	id := c.Param("id")
	var answer models.Answer

	qnaCollection := getQnaCollection()

	fmt.Println("Incoming Answer Request:", answer)

	if err := c.BindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fmt.Println("Question ID:", id)

	qnaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	fullName, nameExists := c.Get("fullName")

	if !nameExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if nameStr, ok := fullName.(string); ok {
		answer.PostedBy = nameStr
	} else {
		answer.PostedBy = "Anonymous"
	}

	answer.ID = primitive.NewObjectID()
	answer.CreatedAt = time.Now()
	answer.Upvotes = 0
	answer.Downvotes = 0

	filter := bson.M{"_id": qnaID}

	// Ensure answers field is an array
	setIfNull := bson.M{
		"$set": bson.M{"answers": bson.M{"$ifNull": []interface{}{bson.M{"$type": "array"}, []models.Answer{}}}},
	}

	// Update MongoDB
	_, err = qnaCollection.UpdateOne(context.TODO(), filter, setIfNull)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize answers array"})
		return
	}

	update := bson.M{"$push": bson.M{"answers": answer}}

	// Perform update
	result, err := qnaCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("MongoDB UpdateOne Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post answer"})
		return
	}

	// Debug: Check MongoDB update result
	fmt.Println("MongoDB Update Result:", result)

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Answer posted successfully"})
}

// Get All QnA (Paginated)
func GetPaginatedQnA(c *gin.Context) {
	qnaCollection := getQnaCollection()
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)
	var qnaList []models.Qna

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	cursor, err := qnaCollection.Find(context.TODO(), bson.M{}, findOptions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var qna models.Qna
		if err := cursor.Decode(&qna); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
			return
		}
		qnaList = append(qnaList, qna)
	}

	c.JSON(http.StatusOK, qnaList)
}

// Upvote/Downvote Question
func VoteQuestion(c *gin.Context) {
	qnaCollection := getQnaCollection()
	id := c.Param("id")
	var request struct {
		VoteType string `json:"voteType"` // "upvote" or "downvote"
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	qnaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	update := bson.M{}
	if request.VoteType == "upvote" {
		update = bson.M{"$inc": bson.M{"upvotes": 1}}
	} else if request.VoteType == "downvote" {
		update = bson.M{"$inc": bson.M{"downvotes": 1}}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote type"})
		return
	}

	_, err = qnaCollection.UpdateOne(context.TODO(), bson.M{"_id": qnaID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

// Upvote/Downvote Answer
func VoteAnswer(c *gin.Context) {
	qnaCollection := getQnaCollection()
	id := c.Param("id")
	var request struct {
		VoteType string `json:"voteType"` // "upvote" or "downvote"
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	answerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	filter := bson.M{"answers._id": answerID}
	update := bson.M{"$inc": bson.M{"answers.$.upvotes": 1}}
	if request.VoteType == "downvote" {
		update = bson.M{"$inc": bson.M{"answers.$.downvotes": 1}}
	}

	_, err = qnaCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}
