package store

import (
	"context"
	"time"

	"github.com/Biswa-bob/bookstore/internal/models"
	"go.mongodb.org/mongo-driver/bson"
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
	GetBooks() ([]*models.Book, error)
}

func (mc *MongoBooksStore) CreateBook(book *models.Book) (*models.Book, error) {

	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	collection := mc.client.Database("bookstore").Collection("books")

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

func (mc *MongoBooksStore) GetBooks() ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mc.client.Database("bookstore").Collection("books")

	// Empty filter matches all documents
	result, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer result.Close(ctx)

	var books []*models.Book
	for result.Next(ctx) {
		var book models.Book
		if err := result.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
