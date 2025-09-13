package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Biswa-bob/bookstore/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateBookDTO struct {
	Title       *string                `json:"title,omitempty"`
	Author      *string                `json:"author,omitempty"`
	PriceCents  *int                   `json:"price_cents,omitempty"`
	PublishedAt *int64                 `json:"publishedAt,omitempty"`
	InStock     *bool                  `json:"inStock,omitempty"`
	ISBN        *string                `json:"isbn,omitempty"`
	Stock       int                    `json:"stock"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type MongoBooksStore struct {
	client *mongo.Client
}

func NewMongoBookStore(client *mongo.Client) *MongoBooksStore {
	return &MongoBooksStore{client: client}
}

type BooksStore interface {
	CreateBook(*models.Book) (*models.Book, error)
	GetBooks() ([]*models.BookResponse, error)
	GetBookById(id primitive.ObjectID) (*models.BookResponse, error)
	UpdateBookFieldsByID(id primitive.ObjectID, upd UpdateBookDTO) (*models.Book, error)
	DeleteBookById(id primitive.ObjectID) error
	GetBookOwner(id primitive.ObjectID) (string, error)
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

func (mc *MongoBooksStore) GetBooks() ([]*models.BookResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mc.client.Database("bookstore").Collection("books")

	// Empty filter matches all documents
	result, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer result.Close(ctx)

	var books []*models.BookResponse
	for result.Next(ctx) {
		var book models.Book
		if err := result.Decode(&book); err != nil {
			return nil, err
		}
		response := book.ToResponse()
		books = append(books, &response)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (mc *MongoBooksStore) GetBookById(objectID primitive.ObjectID) (*models.BookResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mc.client.Database("bookstore").Collection("books")

	// Empty filter matches all documents
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		return nil, err
	}
	response := book.ToResponse()
	return &response, nil
}

func (mc *MongoBooksStore) UpdateBookFieldsByID(id primitive.ObjectID, upd UpdateBookDTO) (*models.Book, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	set := bson.M{}
	if upd.Title != nil {
		set["title"] = *upd.Title
	}
	if upd.Author != nil {
		set["author"] = *upd.Author
	}
	if upd.PriceCents != nil {
		set["price_cents"] = *upd.PriceCents
	}
	if upd.PublishedAt != nil {
		set["publishedAt"] = *upd.PublishedAt
	}
	if upd.InStock != nil {
		set["inStock"] = *upd.InStock
	}

	if len(set) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": set}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After) // return the document AFTER update
	collection := mc.client.Database("bookstore").Collection("books")
	var book models.Book
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // not found
		}
		return nil, err
	}
	return &book, nil
}

func (mc *MongoBooksStore) DeleteBookById(objectID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mc.client.Database("bookstore").Collection("books")

	// Empty filter matches all documents
	var book models.Book
	err := collection.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MongoBooksStore) GetBookOwner(objectID primitive.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := mc.client.Database("bookstore").Collection("books")

	// Empty filter matches all documents
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		return "", err
	}

	return book.UserID, nil
}
