package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Subject    string             `bson:"subject" json:"subject"`
	Semester   string             `bson:"semester" json:"semester"`
	Year       string             `bson:"year" json:"year"`
	Branch     string             `bson:"branch" json:"branch"`
	Content    string             `bson:"content" json:"content"`
	FileUrl    string             `bson:"fileUrl" json:"fileUrl"`
	FileName   string             `bson:"fileName" json:"fileName"`
	UploadedBy string             `bson:"uploadedBy" json:"uploadedBy"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
