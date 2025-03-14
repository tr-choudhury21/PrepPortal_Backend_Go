package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	ImageURL  string             `bson:"imageUrl,omitempty" json:"imageUrl"`
	Author    string             `bson:"author" json:"author"`
	AuthorID  primitive.ObjectID `bson:"author_id" json:"author_id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
