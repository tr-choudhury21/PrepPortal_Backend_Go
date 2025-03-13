package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName  string             `bson:"fullName" json:"fullName"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password,omitempty" json:"password"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
