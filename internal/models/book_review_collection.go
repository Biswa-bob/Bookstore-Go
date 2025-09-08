package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookReview struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID    primitive.ObjectID `bson:"book_id,omitempty" json:"book_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Rating    float64            `bson:"rating" json:"rating"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
