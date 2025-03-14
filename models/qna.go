package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Qna Model
type Qna struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string             `bson:"question" json:"question"`
	Answers   []Answer           `bson:"answers" json:"answers"`
	PostedBy  string             `bson:"postedBy" json:"postedBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Reports   []Report           `bson:"reports" json:"reports"`
	Upvotes   int                `bson:"upvotes" json:"upvotes"`
	Downvotes int                `bson:"downvotes" json:"downvotes"`
}

// Answer model
type Answer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text      string             `bson:"text" json:"text"`
	PostedBy  string             `bson:"postedBy" json:"postedBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Upvotes   int                `bson:"upvotes" json:"upvotes"`
	Downvotes int                `bson:"downvotes" json:"downvotes"`
}

// Report Model
type Report struct {
	ReportedBy primitive.ObjectID `bson:"reportedBy" json:"reportedBy"`
	Reason     string             `bson:"reason" json:"reason"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}
