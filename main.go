package main

import (
	"context"
	"github.com/mauleyzaola/validation/database"
	"github.com/mauleyzaola/validation/interfaces"
	"github.com/mauleyzaola/validation/migrations"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mauleyzaola/validation/assets"
	"github.com/mauleyzaola/validation/middlewares"
	"github.com/mauleyzaola/validation/rest"
	"github.com/mauleyzaola/validation/validators"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// connect to db
	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()

	// apply db migrations
	if err := migrations.MigrateDb(assets.MigrationFiles(),
		"files/migrations", dbUrl); err != nil {
		return err
	}

	// instantiate repository
	repository, err := database.NewPostgresDB(ctx, dbUrl)
	if err != nil {
		return err
	}
	defer func() { repository.Close() }()
	validator := validators.NewValidator()

	// define handlers
	http.HandleFunc("/users", createUserHandler(repository, validator))

	// start http server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}

func createUserHandler(repository interfaces.Repository, validator interfaces.Validator) http.HandlerFunc {
	handler := rest.NewServer(repository, validator)
	return middlewares.MethodChecker(
		http.MethodPost,
		middlewares.JSONValidator(
			handler.Validator(),
			handler.CreateUser,
		),
	)
}
