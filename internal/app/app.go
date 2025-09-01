package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Biswa-bob/bookstore/internal/api"
	"go.mongodb.org/mongo-driver/mongo"

	// "github.com/Biswa-bob/bookstore/internal/middleware"
	"github.com/Biswa-bob/bookstore/internal/store"
	"github.com/Biswa-bob/bookstore/migrations"
)

type Application struct {
	Logger      *log.Logger
	BookHandler *api.BooksHandler
	// UserHandler  *api.UserHandler
	// TokenHandler *api.TokenHandler
	// Middleware   middleware.UserMiddleware
	DB     *sql.DB
	CLIENT *mongo.Client
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()

	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// MongoDB
	mongoClient, err := store.OpenMongo()
	if err != nil {
		log.Fatalf("Mongo error: %v", err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will go here
	booksStore := store.NewMongoBookStore(mongoClient)
	// userStore := store.NewPostgresUserStore(pgDB)
	// tokenStore := store.NewPostgresTokenStore(pgDB)

	// our handlers will go here
	bookHandler := api.NewBooksHandler(booksStore, logger)
	// userHandler := api.NewUserHandler(userStore, logger)
	// tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	// middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:      logger,
		BookHandler: bookHandler,
		// UserHandler:  userHandler,
		// TokenHandler: tokenHandler,
		// Middleware:   middlewareHandler,
		DB:     pgDB,
		CLIENT: mongoClient,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}

func (a *Application) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
	if a.CLIENT != nil {
		ctx := context.TODO()
		if err := a.CLIENT.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB client: %v", err)
		} else {
			log.Println("Disconnected from MongoDB")
		}
	}
}
