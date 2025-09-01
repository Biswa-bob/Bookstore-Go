package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Title       string                 `bson:"title" json:"title"`
	Author      string                 `bson:"author" json:"author"`
	ISBN        string                 `bson:"isbn,omitempty" json:"isbn,omitempty"`
	PriceCents  int                    `bson:"price_cents" json:"price_cents"`
	Stock       int                    `bson:"stock" json:"stock"`
	Description string                 `bson:"description,omitempty" json:"description,omitempty"`
	Tags        []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}
