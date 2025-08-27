package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Biswa-bob/bookstore/internal/api"
	"github.com/Biswa-bob/bookstore/internal/middleware"
	"github.com/Biswa-bob/bookstore/internal/store"
	"github.com/Biswa-bob/bookstore/migrations"
)

type Application struct {
	Logger       *log.Logger
	UserHandler  *api.UserHandler
	TokenHandler *api.TokenHandler
	Middleware   middleware.UserMiddleware
	DB           *sql.DB
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
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will go here
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// our handlers will go here
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:       logger,
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
		Middleware:   middlewareHandler,
		DB:           pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
