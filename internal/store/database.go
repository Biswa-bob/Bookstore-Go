package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "io/fs"
	_ "github.com/jackc/pgx/v4/stdlib"
	// "github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Connected to Database...")
	return db, nil
}

// OpenMongo opens a connection to MongoDB and returns a client + error
func OpenMongo() (*mongo.Client, error) {
	// MongoDB URI (change "localhost" to "mongo" if running inside docker-compose network)
	uri := "mongodb://mongo:mongo@localhost:27017"

	clientOptions := options.Client().ApplyURI(uri)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongo: connect %w", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo: ping %w", err)
	}

	fmt.Println("Connected to MongoDB...")
	return client, nil
}

// func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
// 	goose.SetBaseFS(migrationsFS)
// 	defer func() {
// 		goose.SetBaseFS(nil)
// 	}()
// 	return Migrate(db, dir)
// }

// func Migrate(db *sql.DB, dir string) error {
// 	err := goose.SetDialect("postgres")
// 	if err != nil {
// 		return fmt.Errorf("migrate: %w", err)
// 	}

// 	err = goose.Up(db, dir)
// 	if err != nil {
// 		return fmt.Errorf("goose up: %w", err)
// 	}
// 	return nil
// }
