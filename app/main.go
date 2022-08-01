package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/hashicorp/immutable-infrastructure/db"
	"github.com/hashicorp/immutable-infrastructure/handler"
	"github.com/hashicorp/immutable-infrastructure/tmdbApi"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func goDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	// For example: POSTGRES_URL="postgres://postgres:mysecretpassword@localhost:5432/postgres"

	dbHost := goDotEnvVar("POSTGRES_URL")
	fmt.Printf("godotenv : %s = %s \n", "db Host", dbHost)

	database, err := sql.Open("pgx", dbHost)
	if err != nil {
		log.Fatal("oops, db connection failed", err)
	}

	err = validateSchema(database)
	if err != nil {
		log.Fatal("oops, db migration failed", err)
	}

	log.Println("Serving on 0.0.0.0:8080")

	http.ListenAndServe(":8080", &handler.Handlers{
		MovieGetter: &tmdbApi.Client{
			APIKey: os.Getenv("TMDB"),
		},
		DB: db.New(database),
	})
}

//go:embed db/migrations/*.sql
var fs embed.FS

// Migrate migrates the Postgres schema to the current version.

func validateSchema(db *sql.DB) (retErr error) {
	sourceInstance, err := iofs.New(fs, "db/migrations")
	if err != nil {
		return err
	}
	defer func() {
		err := sourceInstance.Close()
		if retErr == nil {
			retErr = err
		}
	}()
	driverInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceInstance, "postgres", driverInstance)
	if err != nil {
		return err
	}
	err = m.Up() // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
