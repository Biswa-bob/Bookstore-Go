package store

import (
	"context"
	"time"

	"github.com/Biswa-bob/bookstore/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBooksStore struct {
	client *mongo.Client
}

func NewMongoBookStore(client *mongo.Client) *MongoBooksStore {
	return &MongoBooksStore{client: client}
}

type BooksStore interface {
	CreateBook(*models.Book) (*models.Book, error)
}

func (mc *MongoBooksStore) CreateBook(book *models.Book) (*models.Book, error) {
	// set timestamps
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	// choose database + collection
	collection := mc.client.Database("bookstore").Collection("books")

	// insert book
	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		return nil, err
	}

	// assert InsertedID to ObjectID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		book.ID = oid
	}

	return book, nil
}
