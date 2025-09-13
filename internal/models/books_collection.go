package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID        string                 `bson:"user_id" json:"user_id"`
	Title         string                 `bson:"title" json:"title"`
	Author        string                 `bson:"author" json:"author"`
	ISBN          string                 `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Category      string                 `bson:"category,omitempty" json:"category,omitempty"`
	PublishedDate string                 `bson:"published_date,omitempty" json:"published_date,omitempty"`
	Publisher     string                 `bson:"publisher,omitempty" json:"publisher,omitempty"`
	Pages         int                    `bson:"pages,omitempty" json:"pages,omitempty"`
	Language      string                 `bson:"language,omitempty" json:"language,omitempty"`
	PriceCents    int                    `bson:"price_cents" json:"price_cents"`
	Stock         int                    `bson:"stock" json:"stock"`
	Description   string                 `bson:"description,omitempty" json:"description,omitempty"`
	Tags          []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	Metadata      map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt     time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at" json:"updated_at"`
}

type BookResponse struct {
	ID            primitive.ObjectID     `json:"id"`
	UserID        string                 `json:"user_id"`
	Title         string                 `json:"title"`
	Author        string                 `json:"author"`
	ISBN          string                 `json:"isbn,omitempty"`
	Category      string                 `json:"category,omitempty"`
	PublishedDate string                 `json:"published_date,omitempty"`
	Publisher     string                 `json:"publisher,omitempty"`
	Pages         int                    `json:"pages,omitempty"`
	Language      string                 `json:"language,omitempty"`
	PriceCents    int                    `json:"price_cents"`
	Stock         int                    `json:"stock"`
	InStock       bool                   `json:"inStock"`
	Description   string                 `json:"description,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
